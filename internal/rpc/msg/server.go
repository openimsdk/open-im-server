package msg

import (
	"OpenIM/internal/common/check"
	"OpenIM/pkg/common/db/controller"
	"OpenIM/pkg/common/db/localcache"
	"OpenIM/pkg/common/db/relation"
	relationTb "OpenIM/pkg/common/db/table/relation"
	"OpenIM/pkg/common/prome"
	"OpenIM/pkg/discoveryregistry"
	"OpenIM/pkg/proto/msg"
	"google.golang.org/grpc"
)

type msgServer struct {
	RegisterCenter discoveryregistry.SvcDiscoveryRegistry
	MsgDatabase    controller.MsgDatabase
	Group          *check.GroupChecker
	User           *check.UserCheck
	Conversation   *check.ConversationChecker
	friend         *check.FriendChecker
	*localcache.GroupLocalCache
	black         *check.BlackChecker
	MessageLocker MessageLocker
}

func Start(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error {
	mysql, err := relation.NewGormDB()
	if err != nil {
		return err
	}
	if err := mysql.AutoMigrate(&relationTb.UserModel{}); err != nil {
		return err
	}
	s := &msgServer{
		Conversation: check.NewConversationChecker(client),
		User:         check.NewUserCheck(client),
		Group:        check.NewGroupChecker(client),
		//MsgDatabase: controller.MsgDatabase(),
		RegisterCenter:  client,
		GroupLocalCache: localcache.NewGroupMemberIDsLocalCache(client),
		black:           check.NewBlackChecker(client),
		friend:          check.NewFriendChecker(client),
	}
	s.initPrometheus()
	msg.RegisterMsgServer(server, s)
	return nil
}

func (m *msgServer) initPrometheus() {
	prome.NewMsgPullFromRedisSuccessCounter()
	prome.NewMsgPullFromRedisFailedCounter()
	prome.NewMsgPullFromMongoSuccessCounter()
	prome.NewMsgPullFromMongoFailedCounter()
	prome.NewSingleChatMsgRecvSuccessCounter()
	prome.NewGroupChatMsgRecvSuccessCounter()
	prome.NewWorkSuperGroupChatMsgRecvSuccessCounter()
	prome.NewSingleChatMsgProcessSuccessCounter()
	prome.NewSingleChatMsgProcessFailedCounter()
	prome.NewGroupChatMsgProcessSuccessCounter()
	prome.NewGroupChatMsgProcessFailedCounter()
	prome.NewWorkSuperGroupChatMsgProcessSuccessCounter()
	prome.NewWorkSuperGroupChatMsgProcessFailedCounter()
}
