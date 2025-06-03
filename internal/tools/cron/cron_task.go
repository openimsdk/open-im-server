package cron

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	disetcd "github.com/openimsdk/open-im-server/v3/pkg/common/discovery/etcd"
	pbconversation "github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/third"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/discovery/etcd"

	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/utils/runtimeenv"
	"github.com/robfig/cron/v3"
	"google.golang.org/grpc"
)

type Config struct {
	CronTask  config.CronTask
	Share     config.Share
	Discovery config.Discovery
}

func Start(ctx context.Context, conf *Config, client discovery.Conn, service grpc.ServiceRegistrar) error {
	log.CInfo(ctx, "CRON-TASK server is initializing", "runTimeEnv", runtimeenv.RuntimeEnvironment(), "chatRecordsClearTime", conf.CronTask.CronExecuteTime, "msgDestructTime", conf.CronTask.RetainChatRecords)
	if conf.CronTask.RetainChatRecords < 1 {
		log.ZInfo(ctx, "disable cron")
		<-ctx.Done()
		return nil
	}
	ctx = mcontext.SetOpUserID(ctx, conf.Share.IMAdminUserID[0])

	msgConn, err := client.GetConn(ctx, conf.Discovery.RpcService.Msg)
	if err != nil {
		return err
	}

	thirdConn, err := client.GetConn(ctx, conf.Discovery.RpcService.Third)
	if err != nil {
		return err
	}

	conversationConn, err := client.GetConn(ctx, conf.Discovery.RpcService.Conversation)
	if err != nil {
		return err
	}

	if conf.Discovery.Enable == config.ETCD {
		cm := disetcd.NewConfigManager(client.(*etcd.SvcDiscoveryRegistryImpl).GetClient(), []string{
			conf.CronTask.GetConfigFileName(),
			conf.Share.GetConfigFileName(),
			conf.Discovery.GetConfigFileName(),
		})
		cm.Watch(ctx)
	}

	locker, err := NewEtcdLocker(client.(*etcd.SvcDiscoveryRegistryImpl).GetClient())
	if err != nil {
		return err
	}

	if err := locker.Start(ctx); err != nil {
		return err
	}

	srv := &cronServer{
		ctx:                ctx,
		config:             conf,
		cron:               cron.New(),
		msgClient:          msg.NewMsgClient(msgConn),
		conversationClient: pbconversation.NewConversationClient(conversationConn),
		thirdClient:        third.NewThirdClient(thirdConn),
		locker:             locker,
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
	log.ZDebug(ctx, "start cron task", "CronExecuteTime", conf.CronTask.CronExecuteTime)
	srv.cron.Start()
	log.ZDebug(ctx, "cron task server is running")
	<-ctx.Done()
	log.ZDebug(ctx, "cron task server is shutting down")

	locker.Stop() // release distributed lock

	return nil
}

type cronServer struct {
	ctx                context.Context
	config             *Config
	cron               *cron.Cron
	msgClient          msg.MsgClient
	conversationClient pbconversation.ConversationClient
	thirdClient        third.ThirdClient
	locker             *EtcdLocker
}

func (c *cronServer) registerClearS3() error {
	if c.config.CronTask.FileExpireTime <= 0 || len(c.config.CronTask.DeleteObjectType) == 0 {
		log.ZInfo(c.ctx, "disable scheduled cleanup of s3", "fileExpireTime", c.config.CronTask.FileExpireTime, "deleteObjectType", c.config.CronTask.DeleteObjectType)
		return nil
	}
	_, err := c.cron.AddFunc(c.config.CronTask.CronExecuteTime, func() {
		c.locker.ExecuteWithLock(c.ctx, c.clearS3)
	})
	return errs.WrapMsg(err, "failed to register clear s3 cron task")
}

func (c *cronServer) registerDeleteMsg() error {
	if c.config.CronTask.RetainChatRecords <= 0 {
		log.ZInfo(c.ctx, "disable scheduled cleanup of chat records", "retainChatRecords", c.config.CronTask.RetainChatRecords)
		return nil
	}
	_, err := c.cron.AddFunc(c.config.CronTask.CronExecuteTime, func() {
		c.locker.ExecuteWithLock(c.ctx, c.deleteMsg)
	})
	return errs.WrapMsg(err, "failed to register delete msg cron task")
}

func (c *cronServer) registerClearUserMsg() error {
	_, err := c.cron.AddFunc(c.config.CronTask.CronExecuteTime, func() {
		c.locker.ExecuteWithLock(c.ctx, c.clearUserMsg)
	})
	return errs.WrapMsg(err, "failed to register clear user msg cron task")
}
