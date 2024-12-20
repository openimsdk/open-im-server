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
	"github.com/robfig/cron/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
)

func TestName(t *testing.T) {
	conf := &config.Discovery{
		Enable: config.ETCD,
		Etcd: config.Etcd{
			RootDirectory: "openim",
			Address:       []string{"localhost:12379"},
		},
	}
	client, err := kdisc.NewDiscoveryRegister(conf, "source")
	if err != nil {
		panic(err)
	}
	client.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	ctx := mcontext.SetOpUserID(context.Background(), "imAdmin")
	msgConn, err := client.GetConn(ctx, "msg-rpc-service")
	if err != nil {
		panic(err)
	}
	thirdConn, err := client.GetConn(ctx, "third-rpc-service")
	if err != nil {
		panic(err)
	}

	conversationConn, err := client.GetConn(ctx, "conversation-rpc-service")
	if err != nil {
		panic(err)
	}

	srv := &cronServer{
		ctx: ctx,
		config: &CronTaskConfig{
			CronTask: config.CronTask{
				RetainChatRecords: 1,
				FileExpireTime:    1,
				DeleteObjectType:  []string{"msg-picture", "msg-file", "msg-voice", "msg-video", "msg-video-snapshot", "sdklog", ""},
			},
		},
		cron:               cron.New(),
		msgClient:          msg.NewMsgClient(msgConn),
		conversationClient: pbconversation.NewConversationClient(conversationConn),
		thirdClient:        third.NewThirdClient(thirdConn),
	}
	srv.deleteMsg()
	//srv.clearS3()
	//srv.clearUserMsg()
}
