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
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	kdisc "github.com/openimsdk/open-im-server/v3/pkg/common/discoveryregister"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/mw"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"time"

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
	log.CInfo(ctx, "CRON-TASK server is initializing", "chatRecordsClearTime", config.CronTask.ChatRecordsClearTime, "msgDestructTime", config.CronTask.RetainChatRecords)
	if config.CronTask.RetainChatRecords < 1 {
		return errs.New("msg destruct time must be greater than 1").Wrap()
	}
	client, err := kdisc.NewDiscoveryRegister(&config.Discovery, &config.Share)
	if err != nil {
		return errs.WrapMsg(err, "failed to register discovery service")
	}
	client.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	ctx = mcontext.SetOpUserID(ctx, config.Share.IMAdminUserID[0])
	conn, err := client.GetConn(ctx, config.Share.RpcRegisterName.Msg)
	if err != nil {
		return err
	}
	cli := msg.NewMsgClient(conn)
	crontab := cron.New()
	clearFunc := func() {
		now := time.Now()
		deltime := now.Add(-time.Hour * 24 * time.Duration(config.CronTask.RetainChatRecords))
		ctx := mcontext.SetOperationID(ctx, fmt.Sprintf("cron_%d_%d", os.Getpid(), deltime.UnixMilli()))
		log.ZInfo(ctx, "clear chat records", "deltime", deltime, "timestamp", deltime.UnixMilli())
		if _, err := cli.ClearMsg(ctx, &msg.ClearMsgReq{Timestamp: deltime.UnixMilli()}); err != nil {
			log.ZError(ctx, "cron clear chat records failed", err, "deltime", deltime, "cont", time.Since(now))
			return
		}
		log.ZInfo(ctx, "cron clear chat records success", "deltime", deltime, "cont", time.Since(now))
	}
	if _, err := crontab.AddFunc(config.CronTask.ChatRecordsClearTime, clearFunc); err != nil {
		return errs.Wrap(err)
	}
	log.ZInfo(ctx, "start cron task", "chatRecordsClearTime", config.CronTask.ChatRecordsClearTime)
	crontab.Start()
	<-ctx.Done()
	return nil
}
