package msg

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/controller"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/localcache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/tx"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/unrelation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/prome"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
	"google.golang.org/grpc"
)

type MessageInterceptorChain []MessageInterceptorFunc
type msgServer struct {
	RegisterCenter         discoveryregistry.SvcDiscoveryRegistry
	MsgDatabase            controller.CommonMsgDatabase
	ExtendMsgDatabase      controller.ExtendMsgDatabase
	Group                  *rpcclient.GroupClient
	User                   *rpcclient.UserClient
	Conversation           *rpcclient.ConversationClient
	friend                 *rpcclient.FriendClient
	black                  *rpcclient.BlackClient
	GroupLocalCache        *localcache.GroupLocalCache
	ConversationLocalCache *localcache.ConversationLocalCache
	MessageLocker          MessageLocker
	Handlers               MessageInterceptorChain
	notificationSender     *rpcclient.NotificationSender
}

func (m *msgServer) addInterceptorHandler(interceptorFunc ...MessageInterceptorFunc) {
	m.Handlers = append(m.Handlers, interceptorFunc...)
}

func (m *msgServer) execInterceptorHandler(ctx context.Context, req *msg.SendMsgReq) error {
	for _, handler := range m.Handlers {
		msgData, err := handler(ctx, req)
		if err != nil {
			return err
		}
		req.MsgData = msgData
	}
	return nil
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
	if err := mongo.CreateMsgIndex(); err != nil {
		return err
	}
	cacheModel := cache.NewMsgCacheModel(rdb)
	msgDocModel := unrelation.NewMsgMongoDriver(mongo.GetDatabase())
	extendMsgModel := unrelation.NewExtendMsgSetMongoDriver(mongo.GetDatabase())
	extendMsgCacheModel := cache.NewExtendMsgSetCacheRedis(rdb, extendMsgModel, cache.GetDefaultOpt())
	extendMsgDatabase := controller.NewExtendMsgDatabase(extendMsgModel, extendMsgCacheModel, tx.NewMongo(mongo.GetClient()))
	msgDatabase := controller.NewCommonMsgDatabase(msgDocModel, cacheModel)
	s := &msgServer{
		Conversation:           rpcclient.NewConversationClient(client),
		User:                   rpcclient.NewUserClient(client),
		Group:                  rpcclient.NewGroupClient(client),
		MsgDatabase:            msgDatabase,
		ExtendMsgDatabase:      extendMsgDatabase,
		RegisterCenter:         client,
		GroupLocalCache:        localcache.NewGroupLocalCache(client),
		ConversationLocalCache: localcache.NewConversationLocalCache(client),
		black:                  rpcclient.NewBlackClient(client),
		friend:                 rpcclient.NewFriendClient(client),
		MessageLocker:          NewLockerMessage(cacheModel),
	}
	s.notificationSender = rpcclient.NewNotificationSender(rpcclient.WithLocalSendMsg(s.SendMsg))
	s.addInterceptorHandler(MessageHasReadEnabled, MessageModifyCallback)
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
