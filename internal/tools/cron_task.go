// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tools

import (
	"context"
	"strings"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	kdisc "github.com/openimsdk/open-im-server/v3/pkg/common/discoveryregister"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database/mgo"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcli"
	pbconversation "github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/third"
	"github.com/openimsdk/tools/db/mongoutil"

	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/mw"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/robfig/cron/v3"
)

type CronTaskConfig struct {
	CronTask      config.CronTask
	Share         config.Share
	Discovery     config.Discovery
	MongodbConfig config.Mongo
}

func Start(ctx context.Context, cfg *CronTaskConfig) error {
	config.FillCronTaskDefaults(&cfg.CronTask)
	log.CInfo(ctx, "CRON-TASK server is initializing", "chatRecordsClearTime", cfg.CronTask.CronExecuteTime, "msgDestructTime", cfg.CronTask.RetainChatRecords)
	if cfg.CronTask.RetainChatRecords < 1 {
		return errs.New("msg destruct time must be greater than 1").Wrap()
	}
	client, err := kdisc.NewDiscoveryRegister(&cfg.Discovery, &cfg.Share, nil)
	if err != nil {
		return errs.WrapMsg(err, "failed to register discovery service")
	}
	client.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	ctx = mcontext.SetOpUserID(ctx, cfg.Share.IMAdminUserID[0])

	msgConn, err := client.GetConn(ctx, cfg.Share.RpcRegisterName.Msg)
	if err != nil {
		return err
	}
	thirdConn, err := client.GetConn(ctx, cfg.Share.RpcRegisterName.Third)
	if err != nil {
		return err
	}
	conversationConn, err := client.GetConn(ctx, cfg.Share.RpcRegisterName.Conversation)
	if err != nil {
		return err
	}
	authConn, err := client.GetConn(ctx, cfg.Share.RpcRegisterName.Auth)
	if err != nil {
		return err
	}

	mgocli, err := mongoutil.NewMongoDB(ctx, cfg.MongodbConfig.Build())
	if err != nil {
		return errs.WrapMsg(err, "crontask: connect mongodb failed")
	}
	db := mgocli.GetDB()

	userOfflineRecordDB, err := mgo.NewUserOfflineRecordMongo(db)
	if err != nil {
		return errs.WrapMsg(err, "crontask: init user_offline_record collection failed")
	}

	srv := &cronServer{
		ctx:                 ctx,
		config:              cfg,
		cron:                cron.New(),
		msgClient:           msg.NewMsgClient(msgConn),
		conversationClient:  pbconversation.NewConversationClient(conversationConn),
		thirdClient:         third.NewThirdClient(thirdConn),
		authClient:          rpcli.NewAuthClient(authConn),
		userOfflineRecordDB: userOfflineRecordDB,
		chatAPIAddress:      cfg.CronTask.ChatAPI.Address,
	}

	if err := srv.registerClearS3(); err != nil {
		return err
	}
	if err := srv.registerDeleteMsg(); err != nil {
		return err
	}
	if err := srv.registerClearUserMsg(); err != nil {
		return err
	}
	if err := srv.registerClearBurnExpiredMsgs(); err != nil {
		return err
	}
	if err := srv.registerClearGroupBurnExpiredMsgs(); err != nil {
		return err
	}
	if err := srv.registerDeleteExpiredOfflineUsers(); err != nil {
		return err
	}
	log.ZDebug(ctx, "start cron task", "CronExecuteTime", cfg.CronTask.CronExecuteTime)
	srv.cron.Start()
	<-ctx.Done()
	return nil
}

type cronServer struct {
	ctx                 context.Context
	config              *CronTaskConfig
	cron                *cron.Cron
	msgClient           msg.MsgClient
	conversationClient  pbconversation.ConversationClient
	thirdClient         third.ThirdClient
	authClient          *rpcli.AuthClient
	userOfflineRecordDB database.UserOfflineRecord
	chatAPIAddress      string
}

func (c *cronServer) registerClearS3() error {
	if c.config.CronTask.FileExpireTime <= 0 || len(c.config.CronTask.DeleteObjectType) == 0 {
		log.ZInfo(c.ctx, "disable scheduled cleanup of s3", "fileExpireTime", c.config.CronTask.FileExpireTime, "deleteObjectType", c.config.CronTask.DeleteObjectType)
		return nil
	}
	_, err := c.cron.AddFunc(c.config.CronTask.CronExecuteTime, c.clearS3)
	return errs.WrapMsg(err, "failed to register clear s3 cron task")
}

func (c *cronServer) registerDeleteMsg() error {
	if c.config.CronTask.RetainChatRecords <= 0 {
		log.ZInfo(c.ctx, "disable scheduled cleanup of chat records", "retainChatRecords", c.config.CronTask.RetainChatRecords)
		return nil
	}
	_, err := c.cron.AddFunc(c.config.CronTask.CronExecuteTime, c.deleteMsg)
	return errs.WrapMsg(err, "failed to register delete msg cron task")
}

func (c *cronServer) registerClearUserMsg() error {
	_, err := c.cron.AddFunc(c.config.CronTask.CronExecuteTime, c.clearUserMsg)
	return errs.WrapMsg(err, "failed to register clear user msg cron task")
}

func (c *cronServer) registerClearBurnExpiredMsgs() error {
	schedule := strings.TrimSpace(c.config.CronTask.BurnCronExecuteTime)
	if schedule == "" {
		schedule = c.config.CronTask.CronExecuteTime
	}
	_, err := c.cron.AddFunc(schedule, c.clearBurnExpiredMsgs)
	return errs.WrapMsg(err, "failed to register clear burn expired msgs cron task")
}

func (c *cronServer) registerClearGroupBurnExpiredMsgs() error {
	_, err := c.cron.AddFunc(c.config.CronTask.CronExecuteTime, c.clearGroupBurnExpiredMsgs)
	return errs.WrapMsg(err, "failed to register clear group burn expired msgs cron task")
}

// registerDeleteExpiredOfflineUsers 注册每小时执行一次的用户自动删除任务。
// 固定使用 "@hourly" 表达式，与其他任务使用的 CronExecuteTime 独立。
// chatAPI.address 未配置时跳过注册。
func (c *cronServer) registerDeleteExpiredOfflineUsers() error {
	if c.chatAPIAddress == "" {
		log.ZInfo(c.ctx, "disable auto delete expired offline users: chatAPI.address not configured")
		return nil
	}
	_, err := c.cron.AddFunc("@hourly", c.deleteExpiredOfflineUsers)
	return errs.WrapMsg(err, "failed to register delete expired offline users cron task")
}
