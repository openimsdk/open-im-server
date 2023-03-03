package msg

import (
	"OpenIM/internal/common/check"
	"OpenIM/pkg/common/db/cache"
	"OpenIM/pkg/common/db/controller"
	"OpenIM/pkg/common/db/localcache"
	"OpenIM/pkg/common/db/relation"
	relationTb "OpenIM/pkg/common/db/table/relation"
	"OpenIM/pkg/common/db/unrelation"
	"OpenIM/pkg/common/prome"
	"OpenIM/pkg/discoveryregistry"
	"OpenIM/pkg/proto/msg"
	"google.golang.org/grpc"
)

type msgServer struct {
	RegisterCenter    discoveryregistry.SvcDiscoveryRegistry
	MsgDatabase       controller.MsgDatabase
	ExtendMsgDatabase controller.ExtendMsgDatabase
	Group             *check.GroupChecker
	User              *check.UserCheck
	Conversation      *check.ConversationChecker
	friend            *check.FriendChecker
	*localcache.GroupLocalCache
	black         *check.BlackChecker
	MessageLocker MessageLocker
}

func Start(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error {
	rdb, err := cache.NewRedis()
	if err != nil {
		return err
	}
	mongo, err := unrelation.NewMongo()
	if err != nil {
		return err
	}

	cacheModel := cache.NewCacheModel(rdb)
	msgDocModel := unrelation.NewMsgMongoDriver(mongo.GetDatabase())
	extendMsgModel := unrelation.NewExtendMsgSetMongoDriver(mongo.GetDatabase())

	extendMsgDatabase := controller.NewExtendMsgDatabase(extendMsgModel)
	msgDatabase := controller.NewMsgDatabase(msgDocModel, cacheModel)

	s := &msgServer{
		Conversation:      check.NewConversationChecker(client),
		User:              check.NewUserCheck(client),
		Group:             check.NewGroupChecker(client),
		MsgDatabase:       msgDatabase,
		ExtendMsgDatabase: extendMsgDatabase,
		RegisterCenter:    client,
		GroupLocalCache:   localcache.NewGroupLocalCache(client),
		black:             check.NewBlackChecker(client),
		friend:            check.NewFriendChecker(client),
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
