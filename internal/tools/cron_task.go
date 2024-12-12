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
	"github.com/openimsdk/tools/utils/runtimeenv"
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

	runTimeEnv string
}

func Start(ctx context.Context, config *CronTaskConfig) error {
	config.runTimeEnv = runtimeenv.PrintRuntimeEnvironment()

	log.CInfo(ctx, "CRON-TASK server is initializing", "runTimeEnv", config.runTimeEnv, "chatRecordsClearTime", config.CronTask.CronExecuteTime, "msgDestructTime", config.CronTask.RetainChatRecords)
	if config.CronTask.RetainChatRecords < 1 {
		return errs.New("msg destruct time must be greater than 1").Wrap()
	}
	client, err := kdisc.NewDiscoveryRegister(&config.Discovery, config.runTimeEnv)
	if err != nil {
		return errs.WrapMsg(err, "failed to register discovery service")
	}
	client.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	ctx = mcontext.SetOpUserID(ctx, config.Share.IMAdminUserID[0])

	msgConn, err := client.GetConn(ctx, config.Discovery.RpcService.Msg)
	if err != nil {
		return err
	}

	thirdConn, err := client.GetConn(ctx, config.Discovery.RpcService.Third)
	if err != nil {
		return err
	}

	conversationConn, err := client.GetConn(ctx, config.Discovery.RpcService.Conversation)
	if err != nil {
		return err
	}

	msgClient := msg.NewMsgClient(msgConn)
	conversationClient := pbconversation.NewConversationClient(conversationConn)
	thirdClient := third.NewThirdClient(thirdConn)

	crontab := cron.New()

	// scheduled hard delete outdated Msgs in specific time.
	destructMsgsFunc := func() {
		now := time.Now()
		deltime := now.Add(-time.Hour * 24 * time.Duration(config.CronTask.RetainChatRecords))
		ctx := mcontext.SetOperationID(ctx, fmt.Sprintf("cron_%d_%d", os.Getpid(), deltime.UnixMilli()))
		log.ZDebug(ctx, "Destruct chat records", "deltime", deltime, "timestamp", deltime.UnixMilli())

		if _, err := msgClient.DestructMsgs(ctx, &msg.DestructMsgsReq{Timestamp: deltime.UnixMilli()}); err != nil {
			log.ZError(ctx, "cron destruct chat records failed", err, "deltime", deltime, "cont", time.Since(now))
			return
		}
		log.ZDebug(ctx, "cron destruct chat records success", "deltime", deltime, "cont", time.Since(now))
	}
	if _, err := crontab.AddFunc(config.CronTask.CronExecuteTime, destructMsgsFunc); err != nil {
		return errs.Wrap(err)
	}

	// scheduled soft delete outdated Msgs in specific time when user set `is_msg_destruct` feature.
	clearMsgFunc := func() {
		now := time.Now()
		ctx := mcontext.SetOperationID(ctx, fmt.Sprintf("cron_%d_%d", os.Getpid(), now.UnixMilli()))
		log.ZDebug(ctx, "clear msg cron start", "now", now)

		conversations, err := conversationClient.GetConversationsNeedClearMsg(ctx, &pbconversation.GetConversationsNeedClearMsgReq{})
		if err != nil {
			log.ZError(ctx, "Get conversation need Destruct msgs failed.", err)
			return
		}

		_, err = msgClient.ClearMsg(ctx, &msg.ClearMsgReq{Conversations: conversations.Conversations})
		if err != nil {
			log.ZError(ctx, "Clear Msg failed.", err)
			return
		}

		log.ZDebug(ctx, "clear msg cron task completed", "cont", time.Since(now))
	}
	if _, err := crontab.AddFunc(config.CronTask.CronExecuteTime, clearMsgFunc); err != nil {
		return errs.Wrap(err)
	}

	// scheduled delete outdated file Objects and their datas in specific time.
	deleteObjectFunc := func() {
		now := time.Now()
		executeNum := 5
		// number of pagination. if need modify, need update value in third.DeleteOutdatedData
		pageShowNumber := 500
		deleteTime := now.Add(-time.Hour * 24 * time.Duration(config.CronTask.FileExpireTime))
		ctx := mcontext.SetOperationID(ctx, fmt.Sprintf("cron_%d_%d", os.Getpid(), deleteTime.UnixMilli()))
		log.ZDebug(ctx, "deleteoutDatedData", "deletetime", deleteTime, "timestamp", deleteTime.UnixMilli())

		if len(config.CronTask.DeleteObjectType) == 0 {
			log.ZDebug(ctx, "cron deleteoutDatedData not type need delete", "deletetime", deleteTime, "DeleteObjectType", config.CronTask.DeleteObjectType, "cont", time.Since(now))
			return
		}

		for i := 0; i < executeNum; i++ {
			resp, err := thirdClient.DeleteOutdatedData(ctx, &third.DeleteOutdatedDataReq{ExpireTime: deleteTime.UnixMilli(), ObjectGroup: config.CronTask.DeleteObjectType})
			if err != nil {
				log.ZError(ctx, "cron deleteoutDatedData failed", err, "deleteTime", deleteTime, "cont", time.Since(now))
				return
			}
			if resp.Count == 0 || resp.Count < int32(pageShowNumber) {
				break
			}
		}

		log.ZDebug(ctx, "cron deleteoutDatedData success", "deltime", deleteTime, "cont", time.Since(now))
	}
	if _, err := crontab.AddFunc(config.CronTask.CronExecuteTime, deleteObjectFunc); err != nil {
		return errs.Wrap(err)
	}

	log.ZDebug(ctx, "start cron task", "CronExecuteTime", config.CronTask.CronExecuteTime)
	crontab.Start()
	<-ctx.Done()
	return nil
}
