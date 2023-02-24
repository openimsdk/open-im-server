package msg

import (
	"OpenIM/internal/common/check"
	"OpenIM/pkg/common/db/controller"
	"OpenIM/pkg/common/db/localcache"
	"OpenIM/pkg/common/db/relation"
	tablerelation "OpenIM/pkg/common/db/table/relation"
	discoveryRegistry "OpenIM/pkg/discoveryregistry"
	"github.com/OpenIMSDK/openKeeper"

	promePkg "OpenIM/pkg/common/prome"
	"OpenIM/pkg/proto/msg"
	"google.golang.org/grpc"
)

type msgServer struct {
	RegisterCenter discoveryRegistry.SvcDiscoveryRegistry
	MsgDatabase    controller.MsgDatabase
	Group          *check.GroupChecker
	User           *check.UserCheck
	Conversation   *check.ConversationChecker
	friend         *check.FriendChecker
	*localcache.GroupLocalCache
	black         *check.BlackChecker
	MessageLocker MessageLocker
}

type deleteMsg struct {
	UserID      string
	OpUserID    string
	SeqList     []uint32
	OperationID string
}

func Start(client *openKeeper.ZkClient, server *grpc.Server) error {
	mysql, err := relation.NewGormDB()
	if err != nil {
		return err
	}
	if err := mysql.AutoMigrate(&tablerelation.UserModel{}); err != nil {
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
	promePkg.NewMsgPullFromRedisSuccessCounter()
	promePkg.NewMsgPullFromRedisFailedCounter()
	promePkg.NewMsgPullFromMongoSuccessCounter()
	promePkg.NewMsgPullFromMongoFailedCounter()
	promePkg.NewSingleChatMsgRecvSuccessCounter()
	promePkg.NewGroupChatMsgRecvSuccessCounter()
	promePkg.NewWorkSuperGroupChatMsgRecvSuccessCounter()
	promePkg.NewSingleChatMsgProcessSuccessCounter()
	promePkg.NewSingleChatMsgProcessFailedCounter()
	promePkg.NewGroupChatMsgProcessSuccessCounter()
	promePkg.NewGroupChatMsgProcessFailedCounter()
	promePkg.NewWorkSuperGroupChatMsgProcessSuccessCounter()
	promePkg.NewWorkSuperGroupChatMsgProcessFailedCounter()
}
