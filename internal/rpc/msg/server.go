package msg

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/controller"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/localcache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/tx"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/unrelation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/prome"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/tokenverify"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"google.golang.org/grpc"
)

type MessageInterceptorChain []MessageInterceptorFunc
type msgServer struct {
	RegisterCenter    discoveryregistry.SvcDiscoveryRegistry
	MsgDatabase       controller.CommonMsgDatabase
	ExtendMsgDatabase controller.ExtendMsgDatabase
	Group             *rpcclient.GroupClient
	User              *rpcclient.UserClient
	Conversation      *rpcclient.ConversationClient
	friend            *rpcclient.FriendClient
	black             *rpcclient.BlackClient
	GroupLocalCache   *localcache.GroupLocalCache
	MessageLocker     MessageLocker
	Handlers          MessageInterceptorChain
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
	cacheModel := cache.NewMsgCacheModel(rdb)
	msgDocModel := unrelation.NewMsgMongoDriver(mongo.GetDatabase())
	extendMsgModel := unrelation.NewExtendMsgSetMongoDriver(mongo.GetDatabase())
	extendMsgCacheModel := cache.NewExtendMsgSetCacheRedis(rdb, extendMsgModel, cache.GetDefaultOpt())
	extendMsgDatabase := controller.NewExtendMsgDatabase(extendMsgModel, extendMsgCacheModel, tx.NewMongo(mongo.GetClient()))
	msgDatabase := controller.NewCommonMsgDatabase(msgDocModel, cacheModel)

	s := &msgServer{
		Conversation:      rpcclient.NewConversationClient(client),
		User:              rpcclient.NewUserClient(client),
		Group:             rpcclient.NewGroupClient(client),
		MsgDatabase:       msgDatabase,
		ExtendMsgDatabase: extendMsgDatabase,
		RegisterCenter:    client,
		GroupLocalCache:   localcache.NewGroupLocalCache(client),
		black:             rpcclient.NewBlackClient(client),
		friend:            rpcclient.NewFriendClient(client),
		MessageLocker:     NewLockerMessage(cacheModel),
	}
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

func (m *msgServer) SendMsg(ctx context.Context, req *msg.SendMsgReq) (resp *msg.SendMsgResp, error error) {
	resp = &msg.SendMsgResp{}
	flag := isMessageHasReadEnabled(req.MsgData)
	if !flag {
		return nil, errs.ErrMessageHasReadDisable.Wrap()
	}
	m.encapsulateMsgData(req.MsgData)
	if err := CallbackMsgModify(ctx, req); err != nil && err != errs.ErrCallbackContinue {
		return nil, err
	}
	switch req.MsgData.SessionType {
	case constant.SingleChatType:
		return m.sendMsgSingleChat(ctx, req)
	case constant.GroupChatType:
		return m.sendMsgGroupChat(ctx, req)
	case constant.NotificationChatType:
		return m.sendMsgNotification(ctx, req)
	case constant.SuperGroupChatType:
		return m.sendMsgSuperGroupChat(ctx, req)
	default:
		return nil, errs.ErrArgs.Wrap("unknown sessionType")
	}
}

func (m *msgServer) GetMaxSeq(ctx context.Context, req *sdkws.GetMaxSeqReq) (*sdkws.GetMaxSeqResp, error) {
	if err := tokenverify.CheckAccessV3(ctx, req.UserID); err != nil {
		return nil, err
	}
	conversationIDs, err := m.Conversation.GetConversationIDs(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	for _, conversationID := range conversationIDs {
		conversationIDs = append(conversationIDs, utils.GetNotificationConvetstionID(conversationID))
	}
	log.ZDebug(ctx, "GetMaxSeq", "conversationIDs", conversationIDs)
	maxSeqs, err := m.MsgDatabase.GetMaxSeqs(ctx, conversationIDs)
	if err != nil {
		log.ZWarn(ctx, "GetMaxSeqs error", err, "conversationIDs", conversationIDs, "maxSeqs", maxSeqs)
		return nil, err
	}
	resp := new(sdkws.GetMaxSeqResp)
	resp.MaxSeqs = maxSeqs
	return resp, nil
}

func (m *msgServer) PullMessageBySeqs(ctx context.Context, req *sdkws.PullMessageBySeqsReq) (*sdkws.PullMessageBySeqsResp, error) {
	resp := &sdkws.PullMessageBySeqsResp{}
	resp.Msgs = make(map[string]*sdkws.PullMsgs)
	resp.NotificationMsgs = make(map[string]*sdkws.PullMsgs)
	for _, seq := range req.SeqRanges {
		if !utils.IsNotification(seq.ConversationID) {
			msgs, err := m.MsgDatabase.GetMsgBySeqsRange(ctx, seq.ConversationID, seq.Begin, seq.End, seq.Num)
			if err != nil {
				return nil, err
			}
			resp.Msgs[seq.ConversationID] = &sdkws.PullMsgs{Msgs: msgs}
		} else {
			var seqs []int64
			for i := seq.Begin; i <= seq.End; i++ {
				seqs = append(seqs, i)
			}
			notificationMsgs, err := m.MsgDatabase.GetMsgBySeqs(ctx, seq.ConversationID, seqs)
			if err != nil {
				return nil, err
			}
			resp.Msgs[seq.ConversationID] = &sdkws.PullMsgs{Msgs: notificationMsgs}

		}
	}
	return resp, nil
}
