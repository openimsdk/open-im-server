package rpcclient

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	pbauth "github.com/openimsdk/protocol/auth"
	pbconversation "github.com/openimsdk/protocol/conversation"
	pbgroup "github.com/openimsdk/protocol/group"
	pbmsg "github.com/openimsdk/protocol/msg"
	pbmsggateway "github.com/openimsdk/protocol/msggateway"
	pbpush "github.com/openimsdk/protocol/push"
	pbrelation "github.com/openimsdk/protocol/relation"
	pbthird "github.com/openimsdk/protocol/third"
	pbuser "github.com/openimsdk/protocol/user"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/system/program"
	"google.golang.org/grpc"
)

func InitRpcCaller(discov discovery.SvcDiscoveryRegistry, service config.RpcService) error {
	initConn := func(discov discovery.SvcDiscoveryRegistry, name string, initFunc func(conn *grpc.ClientConn)) error {
		conn, err := discov.GetConn(context.Background(), name)
		if err != nil {
			program.ExitWithError(err)
			return err
		}
		initFunc(conn)
		return nil
	}
	if err := initConn(discov, service.Auth, pbauth.InitAuth); err != nil {
		return err
	}
	if err := initConn(discov, service.Conversation, pbconversation.InitConversation); err != nil {
		return err
	}
	if err := initConn(discov, service.Group, pbgroup.InitGroup); err != nil {
		return err
	}
	if err := initConn(discov, service.Msg, pbmsg.InitMsg); err != nil {
		return err
	}
	if err := initConn(discov, service.MessageGateway, pbmsggateway.InitMsgGateway); err != nil {
		return err
	}
	if err := initConn(discov, service.Push, pbpush.InitPushMsgService); err != nil {
		return err
	}
	if err := initConn(discov, service.Friend, pbrelation.InitFriend); err != nil {
		return err
	}
	if err := initConn(discov, service.Third, pbthird.InitThird); err != nil {
		return err
	}
	if err := initConn(discov, service.User, pbuser.InitUser); err != nil {
		return err
	}

	return nil
}
