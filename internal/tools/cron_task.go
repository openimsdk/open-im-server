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
	"fmt"
	"os"
	"time"

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

	msgClient := msg.NewMsgClient(msgConn)
	conversationClient := pbconversation.NewConversationClient(conversationConn)
	thirdClient := third.NewThirdClient(thirdConn)

	crontab := cron.New()

	// scheduled hard delete outdated Msgs in specific time.
	clearMsgFunc := func() {
		now := time.Now()
		deltime := now.Add(-time.Hour * 24 * time.Duration(config.CronTask.RetainChatRecords))
		ctx := mcontext.SetOperationID(ctx, fmt.Sprintf("cron_%d_%d", os.Getpid(), deltime.UnixMilli()))
		log.ZInfo(ctx, "clear chat records", "deltime", deltime, "timestamp", deltime.UnixMilli())
		if _, err := msgClient.ClearMsg(ctx, &msg.ClearMsgReq{Timestamp: deltime.UnixMilli()}); err != nil {
			log.ZError(ctx, "cron clear chat records failed", err, "deltime", deltime, "cont", time.Since(now))
			return
		}
		log.ZInfo(ctx, "cron clear chat records success", "deltime", deltime, "cont", time.Since(now))
	}
	if _, err := crontab.AddFunc(config.CronTask.CronExecuteTime, clearMsgFunc); err != nil {
		return errs.Wrap(err)
	}

	// scheduled soft delete outdated Msgs in specific time when user set `is_msg_destruct` feature.
	msgDestructFunc := func() {
		now := time.Now()
		ctx := mcontext.SetOperationID(ctx, fmt.Sprintf("cron_%d_%d", os.Getpid(), now.UnixMilli()))
		log.ZInfo(ctx, "msg destruct cron start", "now", now)

		conversations, err := conversationClient.GetConversationsNeedDestructMsgs(ctx, &pbconversation.GetConversationsNeedDestructMsgsReq{})
		if err != nil {
			log.ZError(ctx, "Get conversation need Destruct msgs failed.", err)
			return
		} else {
			_, err := msgClient.DestructMsgs(ctx, &msg.DestructMsgsReq{Conversations: conversations.Conversations})
			if err != nil {
				log.ZError(ctx, "Destruct Msgs failed.", err)
				return
			}
		}
		log.ZInfo(ctx, "msg destruct cron task completed", "cont", time.Since(now))
	}
	if _, err := crontab.AddFunc(config.CronTask.CronExecuteTime, msgDestructFunc); err != nil {
		return errs.Wrap(err)
	}

	// scheduled delete outdated file Objects and their datas in specific time.
	deleteObjectFunc := func() {
		now := time.Now()
		deleteTime := now.Add(-time.Hour * 24 * time.Duration(config.CronTask.FileExpireTime))
		ctx := mcontext.SetOperationID(ctx, fmt.Sprintf("cron_%d_%d", os.Getpid(), deleteTime.UnixMilli()))
		log.ZInfo(ctx, "deleteoutDatedData ", "deletetime", deleteTime, "timestamp", deleteTime.UnixMilli())
		if _, err := thirdClient.DeleteOutdatedData(ctx, &third.DeleteOutdatedDataReq{ExpireTime: deleteTime.UnixMilli()}); err != nil {
			log.ZError(ctx, "cron deleteoutDatedData failed", err, "deleteTime", deleteTime, "cont", time.Since(now))
			return
		}
		log.ZInfo(ctx, "cron deleteoutDatedData success", "deltime", deleteTime, "cont", time.Since(now))
	}
	if _, err := crontab.AddFunc(config.CronTask.CronExecuteTime, deleteObjectFunc); err != nil {
		return errs.Wrap(err)
	}

	log.ZInfo(ctx, "start cron task", "CronExecuteTime", config.CronTask.CronExecuteTime)
	crontab.Start()
	<-ctx.Done()
	return nil
}
