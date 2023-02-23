package msg

import (
	"OpenIM/pkg/common/constant"
	promePkg "OpenIM/pkg/common/prome"
	pbConversation "OpenIM/pkg/proto/conversation"
	"OpenIM/pkg/proto/msg"
	"OpenIM/pkg/proto/sdkws"
	"OpenIM/pkg/utils"
	"context"
	"github.com/golang/protobuf/proto"
	"sync"
)

func (m *msgServer) sendMsgSuperGroupChat(ctx context.Context, req *msg.SendMsgReq) (resp *msg.SendMsgResp, err error) {
	resp = &msg.SendMsgResp{}
	promePkg.PromeInc(promePkg.WorkSuperGroupChatMsgRecvSuccessCounter)
	// callback
	if err = CallbackBeforeSendGroupMsg(ctx, req); err != nil && err != constant.ErrCallbackContinue {
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
	if err = CallbackAfterSendGroupMsg(ctx, req); err != nil {
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
	if err = CallbackBeforeSendSingleMsg(ctx, req); err != nil && err != constant.ErrCallbackContinue {
		return nil, err
	}
	_, err = m.messageVerification(ctx, req)
	if err != nil {
		return nil, err
	}
	isSend, err := m.modifyMessageByUserMessageReceiveOpt(ctx, req.MsgData.RecvID, req.MsgData.SendID, constant.SingleChatType, req)
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
	err = CallbackAfterSendSingleMsg(ctx, req)
	if err != nil && err != constant.ErrCallbackContinue {
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
	err = CallbackBeforeSendGroupMsg(ctx, req)
	if err != nil && err != constant.ErrCallbackContinue {
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
			return nil, err
		}
		err = proto.Unmarshal(tips.Detail, &memberKickedTips)
		if err != nil {
			return nil, err
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
	err = CallbackAfterSendGroupMsg(ctx, req)
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

			err := m.Conversation.ModifyConversationField(ctx, &conversationReq)
			if err != nil {
				return
			}

			if tag {
				conversationReq.UserIDList = utils.DifferenceString(atUserID, memberUserIDList)
				conversation.GroupAtType = constant.AtAll
				err := m.Conversation.ModifyConversationField(ctx, &conversationReq)
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

func (m *msgServer) SendMsg(ctx context.Context, req *msg.SendMsgReq) (resp *msg.SendMsgResp, error error) {
	resp = &msg.SendMsgResp{}
	flag := isMessageHasReadEnabled(req.MsgData)
	if !flag {
		return nil, constant.ErrMessageHasReadDisable.Wrap()
	}
	m.encapsulateMsgData(req.MsgData)
	if err := CallbackMsgModify(ctx, req); err != nil && err != constant.ErrCallbackContinue {
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

func (m *msgServer) GetMaxAndMinSeq(ctx context.Context, req *sdkws.GetMaxAndMinSeqReq) (*sdkws.GetMaxAndMinSeqResp, error) {
	resp := new(sdkws.GetMaxAndMinSeqResp)
	m2 := make(map[string]*sdkws.MaxAndMinSeq)
	maxSeq, err := m.MsgInterface.GetUserMaxSeq(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	minSeq, err := m.MsgInterface.GetUserMinSeq(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	resp.MaxSeq = maxSeq
	resp.MinSeq = minSeq
	if len(req.GroupIDList) > 0 {
		resp.GroupMaxAndMinSeq = make(map[string]*sdkws.MaxAndMinSeq)
		for _, groupID := range req.GroupIDList {
			maxSeq, err := m.MsgInterface.GetGroupMaxSeq(ctx, groupID)
			if err != nil {
				return nil, err
			}
			minSeq, err := m.MsgInterface.GetGroupMinSeq(ctx, groupID)
			if err != nil {
				return nil, err
			}
			m2[groupID] = &sdkws.MaxAndMinSeq{
				MaxSeq: maxSeq,
				MinSeq: minSeq,
			}
		}
	}
	return resp, nil
}

func (m *msgServer) PullMessageBySeqList(ctx context.Context, req *sdkws.PullMessageBySeqListReq) (*sdkws.PullMessageBySeqListResp, error) {
	resp := &sdkws.PullMessageBySeqListResp{GroupMsgDataList: make(map[string]*sdkws.MsgDataList)}
	msgs, err := m.MsgInterface.GetMessagesBySeq(ctx, req.UserID, req.Seqs)
	if err != nil {
		return nil, err
	}
	resp.List = msgs
	for userID, list := range req.GroupSeqList {
		msgs, err := m.MsgInterface.GetMessagesBySeq(ctx, userID, req.Seqs)
		if err != nil {
			return nil, err
		}
		resp.GroupMsgDataList[userID] = &sdkws.MsgDataList{
			MsgDataList: msgs,
		}
	}
	return resp, nil
}
