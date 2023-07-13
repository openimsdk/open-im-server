package msg

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"

	"google.golang.org/grpc"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/controller"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/localcache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/unrelation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/prome"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/conversation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
)

type MessageInterceptorChain []MessageInterceptorFunc
type msgServer struct {
	RegisterCenter         discoveryregistry.SvcDiscoveryRegistry
	MsgDatabase            controller.CommonMsgDatabase
	Group                  *rpcclient.GroupRPCClient
	User                   *rpcclient.UserRPCClient
	Conversation           *rpcclient.ConversationRPCClient
	friend                 *rpcclient.FriendRpcClient
	GroupLocalCache        *localcache.GroupLocalCache
	ConversationLocalCache *localcache.ConversationLocalCache
	MessageLocker          MessageLocker
	Handlers               MessageInterceptorChain
	notificationSender     *rpcclient.NotificationSender
}

func (m *msgServer) GetMaxSeq(ctx context.Context, req *sdkws.GetMaxSeqReq) (*sdkws.GetMaxSeqResp, error) {
	// TODO implement me
	panic("implement me")
}

func (m *msgServer) PullMessageBySeqs(ctx context.Context, req *sdkws.PullMessageBySeqsReq) (*sdkws.PullMessageBySeqsResp, error) {
	// TODO implement me
	panic("implement me")
}

func (m *msgServer) SearchMessage(ctx context.Context, req *msg.SearchMessageReq) (*msg.SearchMessageResp, error) {
	// TODO implement me
	panic("implement me")
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
	conversationClient := rpcclient.NewConversationRPCClient(client)
	userRPCClient := rpcclient.NewUserRPCClient(client)
	groupRpcClient := rpcclient.NewGroupRPCClient(client)
	friendRpcClient := rpcclient.NewFriendRpcClient(client)
	mysql, err := relation.NewGormDB()
	msgMysModel := relation.NewChatLogGorm(mysql)
	msgDatabase := controller.NewCommonMsgDatabase(msgDocModel, cacheModel, msgMysModel)
	s := &msgServer{
		Conversation:           &conversationClient,
		User:                   &userRPCClient,
		Group:                  &groupRpcClient,
		MsgDatabase:            msgDatabase,
		RegisterCenter:         client,
		GroupLocalCache:        localcache.NewGroupLocalCache(&groupRpcClient),
		ConversationLocalCache: localcache.NewConversationLocalCache(&conversationClient),
		friend:                 &friendRpcClient,
		MessageLocker:          NewLockerMessage(cacheModel),
	}
	s.notificationSender = rpcclient.NewNotificationSender(rpcclient.WithLocalSendMsg(s.SendMsg))
	s.addInterceptorHandler(MessageHasReadEnabled)
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

func (m *msgServer) conversationAndGetRecvID(conversation *conversation.Conversation, userID string) (recvID string) {
	if conversation.ConversationType == constant.SingleChatType ||
		conversation.ConversationType == constant.NotificationChatType {
		if userID == conversation.OwnerUserID {
			recvID = conversation.UserID
		} else {
			recvID = conversation.OwnerUserID
		}
	} else if conversation.ConversationType == constant.SuperGroupChatType {
		recvID = conversation.GroupID
	}
	return
}
