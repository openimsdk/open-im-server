package msg

import (
	"context"
	"sync"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	promePkg "github.com/OpenIMSDK/Open-IM-Server/pkg/common/prome"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	pbConversation "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/conversation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"google.golang.org/protobuf/proto"
)

func (m *msgServer) sendMsgSuperGroupChat(ctx context.Context, req *msg.SendMsgReq) (resp *msg.SendMsgResp, err error) {
	resp = &msg.SendMsgResp{}
	promePkg.Inc(promePkg.WorkSuperGroupChatMsgRecvSuccessCounter)
	if _, err = m.messageVerification(ctx, req); err != nil {
		promePkg.Inc(promePkg.WorkSuperGroupChatMsgProcessFailedCounter)
		return nil, err
	}
	err = m.MsgDatabase.MsgToMQ(ctx, req.MsgData.GroupID, req.MsgData)
	if err != nil {
		return nil, err
	}
	// callback
	if err = CallbackAfterSendGroupMsg(ctx, req); err != nil && err != errs.ErrCallbackContinue {
		return nil, err
	}

	promePkg.Inc(promePkg.WorkSuperGroupChatMsgProcessSuccessCounter)
	resp.SendTime = req.MsgData.SendTime
	resp.ServerMsgID = req.MsgData.ServerMsgID
	resp.ClientMsgID = req.MsgData.ClientMsgID
	return resp, nil
}

func (m *msgServer) sendMsgNotification(ctx context.Context, req *msg.SendMsgReq) (resp *msg.SendMsgResp, err error) {
	err = m.MsgDatabase.MsgToMQ(ctx, req.MsgData.RecvID, req.MsgData)
	if err != nil {
		return nil, err
	}
	if req.MsgData.SendID != req.MsgData.RecvID { //Filter messages sent to yourself
		err = m.MsgDatabase.MsgToMQ(ctx, req.MsgData.SendID, req.MsgData)
		if err != nil {
			return nil, err
		}
	}
	resp = &msg.SendMsgResp{
		ServerMsgID: req.MsgData.ServerMsgID,
		ClientMsgID: req.MsgData.ClientMsgID,
		SendTime:    req.MsgData.SendTime,
	}
	return resp, nil
}

func (m *msgServer) sendMsgSingleChat(ctx context.Context, req *msg.SendMsgReq) (resp *msg.SendMsgResp, err error) {
	promePkg.Inc(promePkg.SingleChatMsgRecvSuccessCounter)
	_, err = m.messageVerification(ctx, req)
	if err != nil {
		return nil, err
	}

	var isSend bool = true
	conversationID := utils.GetConversationIDByMsg(req.MsgData)

	if utils.MsgIsNotification(req.MsgData) {
		isSend, err = m.modifyMessageByUserMessageReceiveOpt(ctx, req.MsgData.RecvID, conversationID, constant.SingleChatType, req)
		if err != nil {
			return nil, err
		}
	}
	if isSend {
		err = m.MsgDatabase.MsgToMQ(ctx, conversationID, req.MsgData)
		if err != nil {
			return nil, errs.ErrInternalServer.Wrap("insert to mq")
		}
	}
	err = CallbackAfterSendSingleMsg(ctx, req)
	if err != nil && err != errs.ErrCallbackContinue {
		return nil, err
	}
	promePkg.Inc(promePkg.SingleChatMsgProcessSuccessCounter)
	resp = &msg.SendMsgResp{
		ServerMsgID: req.MsgData.ServerMsgID,
		ClientMsgID: req.MsgData.ClientMsgID,
		SendTime:    req.MsgData.SendTime,
	}
	return resp, nil
}

func (m *msgServer) sendMsgGroupChat(ctx context.Context, req *msg.SendMsgReq) (resp *msg.SendMsgResp, err error) {
	// callback
	promePkg.Inc(promePkg.GroupChatMsgRecvSuccessCounter)

	var memberUserIDList []string
	if memberUserIDList, err = m.messageVerification(ctx, req); err != nil {
		promePkg.Inc(promePkg.GroupChatMsgProcessFailedCounter)
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
		go func(i int) {
			err := m.sendMsgToGroupOptimization(ctx, memberUserIDList[i*split:(i+1)*split], tmp, &wg)
			if err != nil {
				mutex.Lock()
				mErr = append(mErr, err)
				mutex.Unlock()
			}

		}(i)
	}
	if remain > 0 {
		wg.Add(1)
		tmp := valueCopy(req)
		go m.sendMsgToGroupOptimization(ctx, memberUserIDList[split*(len(memberUserIDList)/split):], tmp, &wg)
	}

	wg.Wait()

	// callback
	err = CallbackAfterSendGroupMsg(ctx, req)
	if err != nil && err != errs.ErrCallbackContinue {
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
				ConversationID:   utils.GetConversationIDBySessionType(constant.GroupChatType, req.MsgData.GroupID),
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

	promePkg.Inc(promePkg.GroupChatMsgProcessSuccessCounter)
	resp.SendTime = msgToMQSingle.MsgData.SendTime
	resp.ServerMsgID = msgToMQSingle.MsgData.ServerMsgID
	resp.ClientMsgID = msgToMQSingle.MsgData.ClientMsgID
	return resp, nil
}
