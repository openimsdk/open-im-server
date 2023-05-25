package msg

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	promePkg "github.com/OpenIMSDK/Open-IM-Server/pkg/common/prome"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	pbMsg "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
)

func (m *msgServer) SendMsg(ctx context.Context, req *msg.SendMsgReq) (resp *msg.SendMsgResp, error error) {
	resp = &msg.SendMsgResp{}
	flag := isMessageHasReadEnabled(req.MsgData)
	if !flag {
		return nil, errs.ErrMessageHasReadDisable.Wrap()
	}
	m.encapsulateMsgData(req.MsgData)
	if err := callbackMsgModify(ctx, req); err != nil && err != errs.ErrCallbackContinue {
		return nil, err
	}
	switch req.MsgData.SessionType {
	case constant.SingleChatType:
		return m.sendMsgSingleChat(ctx, req)
	case constant.NotificationChatType:
		return m.sendMsgNotification(ctx, req)
	case constant.SuperGroupChatType:
		return m.sendMsgSuperGroupChat(ctx, req)
	default:
		return nil, errs.ErrArgs.Wrap("unknown sessionType")
	}
}

func (m *msgServer) sendMsgSuperGroupChat(ctx context.Context, req *pbMsg.SendMsgReq) (resp *pbMsg.SendMsgResp, err error) {
	resp = &pbMsg.SendMsgResp{}
	promePkg.Inc(promePkg.WorkSuperGroupChatMsgRecvSuccessCounter)
	if _, err = m.messageVerification(ctx, req); err != nil {
		promePkg.Inc(promePkg.WorkSuperGroupChatMsgProcessFailedCounter)
		return nil, err
	}
	err = m.MsgDatabase.MsgToMQ(ctx, utils.GenConversationUniqueKeyForGroup(req.MsgData.GroupID), req.MsgData)
	if err != nil {
		return nil, err
	}
	if err = callbackAfterSendGroupMsg(ctx, req); err != nil {
		log.ZError(ctx, "CallbackAfterSendGroupMsg", err)
	}
	promePkg.Inc(promePkg.WorkSuperGroupChatMsgProcessSuccessCounter)
	resp.SendTime = req.MsgData.SendTime
	resp.ServerMsgID = req.MsgData.ServerMsgID
	resp.ClientMsgID = req.MsgData.ClientMsgID
	return resp, nil
}

func (m *msgServer) sendMsgNotification(ctx context.Context, req *pbMsg.SendMsgReq) (resp *pbMsg.SendMsgResp, err error) {
	promePkg.Inc(promePkg.SingleChatMsgRecvSuccessCounter)
	if err := m.MsgDatabase.MsgToMQ(ctx, utils.GenConversationUniqueKeyForSingle(req.MsgData.SendID, req.MsgData.RecvID), req.MsgData); err != nil {
		promePkg.Inc(promePkg.SingleChatMsgProcessFailedCounter)
		return nil, err
	}
	resp = &pbMsg.SendMsgResp{
		ServerMsgID: req.MsgData.ServerMsgID,
		ClientMsgID: req.MsgData.ClientMsgID,
		SendTime:    req.MsgData.SendTime,
	}
	return resp, nil
}

func (m *msgServer) sendMsgSingleChat(ctx context.Context, req *pbMsg.SendMsgReq) (resp *pbMsg.SendMsgResp, err error) {
	promePkg.Inc(promePkg.SingleChatMsgRecvSuccessCounter)
	_, err = m.messageVerification(ctx, req)
	if err != nil {
		return nil, err
	}
	var isSend bool = true
	isNotification := utils.IsNotificationByMsg(req.MsgData)
	if !isNotification {
		isSend, err = m.modifyMessageByUserMessageReceiveOpt(ctx, req.MsgData.RecvID, utils.GenConversationIDForSingle(req.MsgData.SendID, req.MsgData.RecvID), constant.SingleChatType, req)
		if err != nil {
			return nil, err
		}
	}
	if !isSend {
		promePkg.Inc(promePkg.SingleChatMsgProcessFailedCounter)
		return nil, errs.ErrUserNotRecvMsg
	} else {
		if err := m.MsgDatabase.MsgToMQ(ctx, utils.GenConversationUniqueKeyForSingle(req.MsgData.SendID, req.MsgData.RecvID), req.MsgData); err != nil {
			promePkg.Inc(promePkg.SingleChatMsgProcessFailedCounter)
			return nil, err
		}
		err = callbackAfterSendSingleMsg(ctx, req)
		if err != nil && err != errs.ErrCallbackContinue {
			return nil, err
		}
		resp = &msg.SendMsgResp{
			ServerMsgID: req.MsgData.ServerMsgID,
			ClientMsgID: req.MsgData.ClientMsgID,
			SendTime:    req.MsgData.SendTime,
		}
		promePkg.Inc(promePkg.SingleChatMsgProcessSuccessCounter)
		return resp, nil
	}
}
