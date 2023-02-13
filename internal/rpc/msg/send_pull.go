package msg

import (
	"Open_IM/pkg/common/constant"
	promePkg "Open_IM/pkg/common/prometheus"
	pbConversation "Open_IM/pkg/proto/conversation"
	"Open_IM/pkg/proto/msg"
	"Open_IM/pkg/proto/sdkws"
	"Open_IM/pkg/utils"
	"context"
	go_redis "github.com/go-redis/redis/v8"
	"github.com/golang/protobuf/proto"
	"sync"
)

func (m *msgServer) sendMsgSuperGroupChat(ctx context.Context, req *msg.SendMsgReq) (resp *msg.SendMsgResp, err error) {
	promePkg.PromeInc(promePkg.WorkSuperGroupChatMsgRecvSuccessCounter)
	// callback
	if err = callbackBeforeSendGroupMsg(req); err != nil {
		return nil, err
	}

	if _, err = m.messageVerification(ctx, req); err != nil {
		promePkg.PromeInc(promePkg.WorkSuperGroupChatMsgProcessFailedCounter)
		return nil, err
	}
	msgToMQSingle := msg.MsgDataToMQ{MsgData: req.MsgData}
	err = m.MsgInterface.MsgToMQ(ctx, msgToMQSingle.MsgData.GroupID, &msgToMQSingle)
	if err != nil {
		return nil, err
	}
	// callback
	if err = callbackAfterSendGroupMsg(req); err != nil {
		return nil, err
	}

	promePkg.PromeInc(promePkg.WorkSuperGroupChatMsgProcessSuccessCounter)
	resp.SendTime = msgToMQSingle.MsgData.SendTime
	resp.ServerMsgID = msgToMQSingle.MsgData.ServerMsgID
	resp.ClientMsgID = msgToMQSingle.MsgData.ClientMsgID
	return resp, nil
}
func (m *msgServer) sendMsgNotification(ctx context.Context, req *msg.SendMsgReq) (resp *msg.SendMsgResp, err error) {
	msgToMQSingle := msg.MsgDataToMQ{MsgData: req.MsgData}
	err = m.MsgInterface.MsgToMQ(ctx, msgToMQSingle.MsgData.RecvID, &msgToMQSingle)
	if err != nil {
		return nil, err
	}
	if msgToMQSingle.MsgData.SendID != msgToMQSingle.MsgData.RecvID { //Filter messages sent to yourself
		err = m.MsgInterface.MsgToMQ(ctx, msgToMQSingle.MsgData.SendID, &msgToMQSingle)
		if err != nil {
			return nil, err
		}
	}

	resp.SendTime = msgToMQSingle.MsgData.SendTime
	resp.ServerMsgID = msgToMQSingle.MsgData.ServerMsgID
	resp.ClientMsgID = msgToMQSingle.MsgData.ClientMsgID
	return resp, nil
}

func (m *msgServer) sendMsgSingleChat(ctx context.Context, req *msg.SendMsgReq) (resp *msg.SendMsgResp, err error) {
	promePkg.PromeInc(promePkg.SingleChatMsgRecvSuccessCounter)
	if err = callbackBeforeSendSingleMsg(req); err != nil {
		return nil, err
	}
	_, err = m.messageVerification(ctx, req)
	if err != nil {
		return nil, err
	}
	isSend, err := modifyMessageByUserMessageReceiveOpt(req.MsgData.RecvID, req.MsgData.SendID, constant.SingleChatType, req)
	if err != nil {
		return nil, err
	}
	msgToMQSingle := msg.MsgDataToMQ{MsgData: req.MsgData}
	if isSend {
		err = m.MsgInterface.MsgToMQ(ctx, req.MsgData.RecvID, &msgToMQSingle)
		if err != nil {
			return nil, constant.ErrInternalServer.Wrap("insert to mq")
		}
	}
	if msgToMQSingle.MsgData.SendID != msgToMQSingle.MsgData.RecvID { //Filter messages sent to yourself
		err = m.MsgInterface.MsgToMQ(ctx, req.MsgData.SendID, &msgToMQSingle)
		if err != nil {
			return nil, constant.ErrInternalServer.Wrap("insert to mq")
		}
	}
	err = callbackAfterSendSingleMsg(req)
	if err != nil {
		return nil, err
	}
	promePkg.PromeInc(promePkg.SingleChatMsgProcessSuccessCounter)
	resp.SendTime = msgToMQSingle.MsgData.SendTime
	resp.ServerMsgID = msgToMQSingle.MsgData.ServerMsgID
	resp.ClientMsgID = msgToMQSingle.MsgData.ClientMsgID
	return resp, nil
}

func (m *msgServer) sendMsgGroupChat(ctx context.Context, req *msg.SendMsgReq) (resp *msg.SendMsgResp, err error) {
	// callback
	promePkg.PromeInc(promePkg.GroupChatMsgRecvSuccessCounter)
	err = callbackBeforeSendGroupMsg(req)
	if err != nil {
		return nil, err
	}

	var memberUserIDList []string
	if memberUserIDList, err = m.messageVerification(ctx, req); err != nil {
		promePkg.PromeInc(promePkg.GroupChatMsgProcessFailedCounter)
		return nil, err
	}

	var addUidList []string
	switch req.MsgData.ContentType {
	case constant.MemberKickedNotification:
		var tips sdkws.TipsComm
		var memberKickedTips sdkws.MemberKickedTips
		err := proto.Unmarshal(req.MsgData.Content, &tips)
		if err != nil {

		}
		err = proto.Unmarshal(tips.Detail, &memberKickedTips)
		if err != nil {

		}

		for _, v := range memberKickedTips.KickedUserList {
			addUidList = append(addUidList, v.UserID)
		}
	case constant.MemberQuitNotification:
		addUidList = append(addUidList, req.MsgData.SendID)

	default:
	}
	if len(addUidList) > 0 {
		memberUserIDList = append(memberUserIDList, addUidList...)
	}

	//split  parallel send
	var wg sync.WaitGroup
	var split = 20
	msgToMQSingle := msg.MsgDataToMQ{MsgData: req.MsgData}
	mErr := make([]error, 0)
	var mutex sync.RWMutex
	remain := len(memberUserIDList) % split
	for i := 0; i < len(memberUserIDList)/split; i++ {
		wg.Add(1)
		tmp := valueCopy(req)
		go func() {
			err := m.sendMsgToGroupOptimization(ctx, memberUserIDList[i*split:(i+1)*split], tmp, &wg)
			if err != nil {
				mutex.Lock()
				mErr = append(mErr, err)
				mutex.Unlock()
			}

		}()
	}
	if remain > 0 {
		wg.Add(1)
		tmp := valueCopy(req)
		go m.sendMsgToGroupOptimization(ctx, memberUserIDList[split*(len(memberUserIDList)/split):], tmp, &wg)
	}

	wg.Wait()

	// callback
	err = callbackAfterSendGroupMsg(req)
	if err != nil {
		return nil, err
	}

	for _, v := range mErr {
		if v != nil {
			return nil, v
		}
	}

	if req.MsgData.ContentType == constant.AtText {
		go func() {
			var conversationReq pbConversation.ModifyConversationFieldReq
			var tag bool
			var atUserID []string
			conversation := pbConversation.Conversation{
				OwnerUserID:      req.MsgData.SendID,
				ConversationID:   utils.GetConversationIDBySessionType(req.MsgData.GroupID, constant.GroupChatType),
				ConversationType: constant.GroupChatType,
				GroupID:          req.MsgData.GroupID,
			}
			conversationReq.Conversation = &conversation
			conversationReq.FieldType = constant.FieldGroupAtType
			tagAll := utils.IsContain(constant.AtAllString, req.MsgData.AtUserIDList)
			if tagAll {
				atUserID = utils.DifferenceString([]string{constant.AtAllString}, req.MsgData.AtUserIDList)
				if len(atUserID) == 0 { //just @everyone
					conversationReq.UserIDList = memberUserIDList
					conversation.GroupAtType = constant.AtAll
				} else { //@Everyone and @other people
					conversationReq.UserIDList = atUserID
					conversation.GroupAtType = constant.AtAllAtMe
					tag = true
				}
			} else {
				conversationReq.UserIDList = req.MsgData.AtUserIDList
				conversation.GroupAtType = constant.AtMe
			}

			_, err := m.ModifyConversationField(context.Background(), &conversationReq)
			if err != nil {
				return
			}

			if tag {
				conversationReq.UserIDList = utils.DifferenceString(atUserID, memberUserIDList)
				conversation.GroupAtType = constant.AtAll
				_, err := m.ModifyConversationField(context.Background(), &conversationReq)
				if err != nil {
					return
				}
			}
		}()
	}
	//

	promePkg.PromeInc(promePkg.GroupChatMsgProcessSuccessCounter)
	resp.SendTime = msgToMQSingle.MsgData.SendTime
	resp.ServerMsgID = msgToMQSingle.MsgData.ServerMsgID
	resp.ClientMsgID = msgToMQSingle.MsgData.ClientMsgID
	return resp, nil
}

func (m *msgServer) ModifyConversationField(ctx context.Context, req *pbConversation.ModifyConversationFieldReq) (*pbConversation.ModifyConversationFieldResp, error) {

}

func (m *msgServer) SendMsg(ctx context.Context, req *msg.SendMsgReq) (resp *msg.SendMsgResp, error error) {
	resp = &msg.SendMsgResp{}
	flag := isMessageHasReadEnabled(req.MsgData)
	if !flag {
		return nil, constant.ErrMessageHasReadDisable.Wrap()
	}
	m.encapsulateMsgData(req.MsgData)
	if err := callbackMsgModify(req); err != nil {
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
		return nil, constant.ErrArgs.Wrap("unknown sessionType")
	}
}

func (m *msgServer) GetMaxAndMinSeq(_ context.Context, in *sdkws.GetMaxAndMinSeqReq) (*sdkws.GetMaxAndMinSeqResp, error) {
	log.NewInfo(in.OperationID, "rpc getMaxAndMinSeq is arriving", in.String())
	resp := new(sdkws.GetMaxAndMinSeqResp)
	m := make(map[string]*sdkws.MaxAndMinSeq)
	var maxSeq, minSeq uint64
	var err1, err2 error
	maxSeq, err1 = commonDB.DB.GetUserMaxSeq(in.UserID)
	minSeq, err2 = commonDB.DB.GetUserMinSeq(in.UserID)
	if (err1 != nil && err1 != go_redis.Nil) || (err2 != nil && err2 != go_redis.Nil) {
		log.NewError(in.OperationID, "getMaxSeq from redis error", in.String())
		if err1 != nil {
			log.NewError(in.OperationID, utils.GetSelfFuncName(), err1.Error())
		}
		if err2 != nil {
			log.NewError(in.OperationID, utils.GetSelfFuncName(), err2.Error())
		}
		resp.ErrCode = 200
		resp.ErrMsg = "redis get err"
		return resp, nil
	}
	resp.MaxSeq = uint32(maxSeq)
	resp.MinSeq = uint32(minSeq)
	for _, groupID := range in.GroupIDList {
		x := new(sdkws.MaxAndMinSeq)
		maxSeq, _ := commonDB.DB.GetGroupMaxSeq(groupID)
		minSeq, _ := commonDB.DB.GetGroupUserMinSeq(groupID, in.UserID)
		x.MaxSeq = uint32(maxSeq)
		x.MinSeq = uint32(minSeq)
		m[groupID] = x
	}
	resp.GroupMaxAndMinSeq = m
	return resp, nil
}

func (rpc *rpcChat) PullMessageBySeqList(_ context.Context, in *sdkws.PullMessageBySeqListReq) (*sdkws.PullMessageBySeqListResp, error) {
	log.NewInfo(in.OperationID, "rpc PullMessageBySeqList is arriving", in.String())
	resp := new(sdkws.PullMessageBySeqListResp)
	m := make(map[string]*sdkws.MsgDataList)
	redisMsgList, failedSeqList, err := commonDB.DB.GetMessageListBySeq(in.UserID, in.SeqList, in.OperationID)
	if err != nil {
		if err != go_redis.Nil {
			promePkg.PromeAdd(promePkg.MsgPullFromRedisFailedCounter, len(failedSeqList))
			log.Error(in.OperationID, "get message from redis exception", err.Error(), failedSeqList)
		} else {
			log.Debug(in.OperationID, "get message from redis is nil", failedSeqList)
		}
		msgList, err1 := commonDB.DB.GetMsgBySeqListMongo2(in.UserID, failedSeqList, in.OperationID)
		if err1 != nil {
			promePkg.PromeAdd(promePkg.MsgPullFromMongoFailedCounter, len(failedSeqList))
			log.Error(in.OperationID, "PullMessageBySeqList data error", in.String(), err1.Error())
			resp.ErrCode = 201
			resp.ErrMsg = err1.Error()
			return resp, nil
		} else {
			promePkg.PromeAdd(promePkg.MsgPullFromMongoSuccessCounter, len(msgList))
			redisMsgList = append(redisMsgList, msgList...)
			resp.List = redisMsgList
		}
	} else {
		promePkg.PromeAdd(promePkg.MsgPullFromRedisSuccessCounter, len(redisMsgList))
		resp.List = redisMsgList
	}

	for k, v := range in.GroupSeqList {
		x := new(sdkws.MsgDataList)
		redisMsgList, failedSeqList, err := commonDB.DB.GetMessageListBySeq(k, v.SeqList, in.OperationID)
		if err != nil {
			if err != go_redis.Nil {
				promePkg.PromeAdd(promePkg.MsgPullFromRedisFailedCounter, len(failedSeqList))
				log.Error(in.OperationID, "get message from redis exception", err.Error(), failedSeqList)
			} else {
				log.Debug(in.OperationID, "get message from redis is nil", failedSeqList)
			}
			msgList, err1 := commonDB.DB.GetSuperGroupMsgBySeqListMongo(k, failedSeqList, in.OperationID)
			if err1 != nil {
				promePkg.PromeAdd(promePkg.MsgPullFromMongoFailedCounter, len(failedSeqList))
				log.Error(in.OperationID, "PullMessageBySeqList data error", in.String(), err1.Error())
				resp.ErrCode = 201
				resp.ErrMsg = err1.Error()
				return resp, nil
			} else {
				promePkg.PromeAdd(promePkg.MsgPullFromMongoSuccessCounter, len(msgList))
				redisMsgList = append(redisMsgList, msgList...)
				x.MsgDataList = redisMsgList
				m[k] = x
			}
		} else {
			promePkg.PromeAdd(promePkg.MsgPullFromRedisSuccessCounter, len(redisMsgList))
			x.MsgDataList = redisMsgList
			m[k] = x
		}
	}
	resp.GroupMsgDataList = m
	return resp, nil
}
