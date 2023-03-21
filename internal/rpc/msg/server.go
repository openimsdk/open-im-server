package msg

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/controller"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/localcache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/unrelation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/prome"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/tokenverify"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient/check"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
)

type MessageInterceptorChain []MessageInterceptorFunc
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
	Handlers      MessageInterceptorChain
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

func (m *msgServer) GetMaxAndMinSeq(ctx context.Context, req *sdkws.GetMaxAndMinSeqReq) (*sdkws.GetMaxAndMinSeqResp, error) {
	if err := tokenverify.CheckAccessV3(ctx, req.UserID); err != nil {
		return nil, err
	}
	resp := new(sdkws.GetMaxAndMinSeqResp)
	m2 := make(map[string]*sdkws.MaxAndMinSeq)
	maxSeq, err := m.MsgDatabase.GetUserMaxSeq(ctx, req.UserID)
	if err != nil && errs.Unwrap(err) != redis.Nil {
		return nil, err
	}
	minSeq, err := m.MsgDatabase.GetUserMinSeq(ctx, req.UserID)
	if err != nil && errs.Unwrap(err) != redis.Nil {
		return nil, err
	}
	resp.MaxSeq = maxSeq
	resp.MinSeq = minSeq
	if len(req.GroupIDs) > 0 {
		for _, groupID := range req.GroupIDs {
			maxSeq, err := m.MsgDatabase.GetGroupMaxSeq(ctx, groupID)
			if err != nil && errs.Unwrap(err) != redis.Nil {
				return nil, err
			}
			minSeq, err := m.MsgDatabase.GetGroupMinSeq(ctx, groupID)
			if err != nil && errs.Unwrap(err) != redis.Nil {
				return nil, err
			}
			m2[groupID] = &sdkws.MaxAndMinSeq{
				MaxSeq: maxSeq,
				MinSeq: minSeq,
			}
		}
	}
	resp.GroupMaxAndMinSeq = m2
	return resp, nil
}

func (m *msgServer) PullMessageBySeqs(ctx context.Context, req *sdkws.PullMessageBySeqsReq) (*sdkws.PullMessageBySeqsResp, error) {
	resp := &sdkws.PullMessageBySeqsResp{GroupMsgDataList: make(map[string]*sdkws.MsgDataList)}
	msgs, err := m.MsgDatabase.GetMsgBySeqs(ctx, req.UserID, req.Seqs)
	if err != nil {
		return nil, err
	}
	resp.List = msgs
	for groupID, list := range req.GroupSeqs {
		msgs, err := m.MsgDatabase.GetSuperGroupMsgBySeqs(ctx, groupID, list.Seqs)
		if err != nil {
			return nil, err
		}
		resp.GroupMsgDataList[groupID] = &sdkws.MsgDataList{
			MsgDataList: msgs,
		}
	}
	return resp, nil
}
