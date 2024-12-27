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
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	kdisc "github.com/openimsdk/open-im-server/v3/pkg/common/discoveryregister"
	pbconversation "github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/third"

	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/mw"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/robfig/cron/v3"
)

type CronTaskConfig struct {
	CronTask  config.CronTask
	Share     config.Share
	Discovery config.Discovery
}

func Start(ctx context.Context, config *CronTaskConfig) error {
	log.CInfo(ctx, "CRON-TASK server is initializing", "chatRecordsClearTime", config.CronTask.CronExecuteTime, "msgDestructTime", config.CronTask.RetainChatRecords)
	if config.CronTask.RetainChatRecords < 1 {
		return errs.New("msg destruct time must be greater than 1").Wrap()
	}
	client, err := kdisc.NewDiscoveryRegister(&config.Discovery, &config.Share)
	if err != nil {
		return errs.WrapMsg(err, "failed to register discovery service")
	}
	client.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	ctx = mcontext.SetOpUserID(ctx, config.Share.IMAdminUserID[0])

	msgConn, err := client.GetConn(ctx, config.Share.RpcRegisterName.Msg)
	if err != nil {
		return err
	}

	thirdConn, err := client.GetConn(ctx, config.Share.RpcRegisterName.Third)
	if err != nil {
		return err
	}

	conversationConn, err := client.GetConn(ctx, config.Share.RpcRegisterName.Conversation)
	if err != nil {
		return err
	}

	srv := &cronServer{
		ctx:                ctx,
		config:             config,
		cron:               cron.New(),
		msgClient:          msg.NewMsgClient(msgConn),
		conversationClient: pbconversation.NewConversationClient(conversationConn),
		thirdClient:        third.NewThirdClient(thirdConn),
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
	log.ZDebug(ctx, "start cron task", "CronExecuteTime", config.CronTask.CronExecuteTime)
	srv.cron.Start()
	<-ctx.Done()
	return nil
}

type cronServer struct {
	ctx                context.Context
	config             *CronTaskConfig
	cron               *cron.Cron
	msgClient          msg.MsgClient
	conversationClient pbconversation.ConversationClient
	thirdClient        third.ThirdClient
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
