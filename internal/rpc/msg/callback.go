package msg

import (
	cbApi "Open_IM/pkg/call_back_struct"
	"Open_IM/pkg/common/callback"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/http"
	"Open_IM/pkg/common/log"
	pbChat "Open_IM/pkg/proto/msg"
	"Open_IM/pkg/utils"
	http2 "net/http"
)

func copyCallbackCommonReqStruct(msg *pbChat.SendMsgReq) cbApi.CommonCallbackReq {
	req := cbApi.CommonCallbackReq{
		SendID:           msg.MsgData.SendID,
		ServerMsgID:      msg.MsgData.ServerMsgID,
		ClientMsgID:      msg.MsgData.ClientMsgID,
		OperationID:      msg.OperationID,
		SenderPlatformID: msg.MsgData.SenderPlatformID,
		SenderNickname:   msg.MsgData.SenderNickname,
		SessionType:      msg.MsgData.SessionType,
		MsgFrom:          msg.MsgData.MsgFrom,
		ContentType:      msg.MsgData.ContentType,
		Status:           msg.MsgData.Status,
		CreateTime:       msg.MsgData.CreateTime,
		AtUserIDList:     msg.MsgData.AtUserIDList,
		SenderFaceURL:    msg.MsgData.SenderFaceURL,
		Content:          callback.GetContent(msg.MsgData),
		Seq:              msg.MsgData.Seq,
		Ex:               msg.MsgData.Ex,
	}
	return req
}

func callbackBeforeSendSingleMsg(msg *pbChat.SendMsgReq) cbApi.CommonCallbackResp {
	callbackResp := cbApi.CommonCallbackResp{OperationID: msg.OperationID}
	if !config.Config.Callback.CallbackBeforeSendSingleMsg.Enable {
		return callbackResp
	}
	log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), msg)
	commonCallbackReq := copyCallbackCommonReqStruct(msg)
	commonCallbackReq.CallbackCommand = constant.CallbackBeforeSendSingleMsgCommand
	req := cbApi.CallbackBeforeSendSingleMsgReq{
		CommonCallbackReq: commonCallbackReq,
		RecvID:            msg.MsgData.RecvID,
	}
	resp := &cbApi.CallbackBeforeSendSingleMsgResp{
		CommonCallbackResp: &callbackResp,
	}
	//utils.CopyStructFields(req, msg.MsgData)
	defer log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), req, *resp)
	if err := http.CallBackPostReturn(config.Config.Callback.CallbackUrl, constant.CallbackBeforeSendSingleMsgCommand, req, resp, config.Config.Callback.CallbackBeforeSendSingleMsg.CallbackTimeOut); err != nil {
		callbackResp.ErrCode = http2.StatusInternalServerError
		callbackResp.ErrMsg = err.Error()
		if !config.Config.Callback.CallbackBeforeSendSingleMsg.CallbackFailedContinue {
			callbackResp.ActionCode = constant.ActionForbidden
			return callbackResp
		} else {
			callbackResp.ActionCode = constant.ActionAllow
			return callbackResp
		}
	}
	return callbackResp
}

func callbackAfterSendSingleMsg(msg *pbChat.SendMsgReq) cbApi.CommonCallbackResp {
	callbackResp := cbApi.CommonCallbackResp{OperationID: msg.OperationID}
	if !config.Config.Callback.CallbackAfterSendSingleMsg.Enable {
		return callbackResp
	}
	log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), msg)
	commonCallbackReq := copyCallbackCommonReqStruct(msg)
	commonCallbackReq.CallbackCommand = constant.CallbackAfterSendSingleMsgCommand
	req := cbApi.CallbackAfterSendSingleMsgReq{
		CommonCallbackReq: commonCallbackReq,
		RecvID:            msg.MsgData.RecvID,
	}
	resp := &cbApi.CallbackAfterSendSingleMsgResp{CommonCallbackResp: &callbackResp}
	defer log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), req, *resp)
	if err := http.CallBackPostReturn(config.Config.Callback.CallbackUrl, constant.CallbackAfterSendSingleMsgCommand, req, resp, config.Config.Callback.CallbackAfterSendSingleMsg.CallbackTimeOut); err != nil {
		callbackResp.ErrCode = http2.StatusInternalServerError
		callbackResp.ErrMsg = err.Error()
		return callbackResp
	}
	return callbackResp
}

func callbackBeforeSendGroupMsg(msg *pbChat.SendMsgReq) cbApi.CommonCallbackResp {
	callbackResp := cbApi.CommonCallbackResp{OperationID: msg.OperationID}
	if !config.Config.Callback.CallbackBeforeSendGroupMsg.Enable {
		return callbackResp
	}
	log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), msg)
	commonCallbackReq := copyCallbackCommonReqStruct(msg)
	commonCallbackReq.CallbackCommand = constant.CallbackBeforeSendGroupMsgCommand
	req := cbApi.CallbackAfterSendGroupMsgReq{
		CommonCallbackReq: commonCallbackReq,
		GroupID:           msg.MsgData.GroupID,
	}
	resp := &cbApi.CallbackBeforeSendGroupMsgResp{CommonCallbackResp: &callbackResp}
	defer log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), req, *resp)
	if err := http.CallBackPostReturn(config.Config.Callback.CallbackUrl, constant.CallbackBeforeSendGroupMsgCommand, req, resp, config.Config.Callback.CallbackBeforeSendGroupMsg.CallbackTimeOut); err != nil {
		callbackResp.ErrCode = http2.StatusInternalServerError
		callbackResp.ErrMsg = err.Error()
		if !config.Config.Callback.CallbackBeforeSendGroupMsg.CallbackFailedContinue {
			callbackResp.ActionCode = constant.ActionForbidden
			return callbackResp
		} else {
			callbackResp.ActionCode = constant.ActionAllow
			return callbackResp
		}
	}
	return callbackResp
}

func callbackAfterSendGroupMsg(msg *pbChat.SendMsgReq) cbApi.CommonCallbackResp {
	callbackResp := cbApi.CommonCallbackResp{OperationID: msg.OperationID}
	if !config.Config.Callback.CallbackAfterSendGroupMsg.Enable {
		return callbackResp
	}
	log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), msg)
	commonCallbackReq := copyCallbackCommonReqStruct(msg)
	commonCallbackReq.CallbackCommand = constant.CallbackAfterSendGroupMsgCommand
	req := cbApi.CallbackAfterSendGroupMsgReq{
		CommonCallbackReq: commonCallbackReq,
		GroupID:           msg.MsgData.GroupID,
	}
	resp := &cbApi.CallbackAfterSendGroupMsgResp{CommonCallbackResp: &callbackResp}
	defer log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), req, *resp)
	if err := http.CallBackPostReturn(config.Config.Callback.CallbackUrl, constant.CallbackAfterSendGroupMsgCommand, req, resp, config.Config.Callback.CallbackAfterSendGroupMsg.CallbackTimeOut); err != nil {
		callbackResp.ErrCode = http2.StatusInternalServerError
		callbackResp.ErrMsg = err.Error()
		return callbackResp
	}
	return callbackResp
}

func callbackMsgModify(msg *pbChat.SendMsgReq) cbApi.CommonCallbackResp {
	log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), msg)
	callbackResp := cbApi.CommonCallbackResp{OperationID: msg.OperationID}
	if !config.Config.Callback.CallbackMsgModify.Enable {
		return callbackResp
	}
	commonCallbackReq := copyCallbackCommonReqStruct(msg)
	commonCallbackReq.CallbackCommand = constant.CallbackMsgModifyCommand
	req := cbApi.CallbackMsgModifyCommandReq{
		CommonCallbackReq: commonCallbackReq,
		GroupID:           msg.MsgData.GroupID,
		RecvID:            msg.MsgData.RecvID,
	}
	resp := &cbApi.CallbackMsgModifyCommandResp{CommonCallbackResp: &callbackResp}
	defer log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), req, *resp)
	if err := http.CallBackPostReturn(config.Config.Callback.CallbackUrl, constant.CallbackMsgModifyCommand, req, resp, config.Config.Callback.CallbackMsgModify.CallbackTimeOut); err != nil {
		callbackResp.ErrCode = http2.StatusInternalServerError
		callbackResp.ErrMsg = err.Error()
		if !config.Config.Callback.CallbackMsgModify.CallbackFailedContinue {
			callbackResp.ActionCode = constant.ActionForbidden
			return callbackResp
		} else {
			callbackResp.ActionCode = constant.ActionAllow
			return callbackResp
		}
	}
	if resp.ErrCode == constant.CallbackHandleSuccess && resp.ActionCode == constant.ActionAllow {
		if resp.Content != nil {
			msg.MsgData.Content = []byte(*resp.Content)
		}
		if resp.RecvID != nil {
			msg.MsgData.RecvID = *resp.RecvID
		}
		if resp.GroupID != nil {
			msg.MsgData.GroupID = *resp.GroupID
		}
		if resp.ClientMsgID != nil {
			msg.MsgData.ClientMsgID = *resp.ClientMsgID
		}
		if resp.ServerMsgID != nil {
			msg.MsgData.ServerMsgID = *resp.ServerMsgID
		}
		if resp.SenderPlatformID != nil {
			msg.MsgData.SenderPlatformID = *resp.SenderPlatformID
		}
		if resp.SenderNickname != nil {
			msg.MsgData.SenderNickname = *resp.SenderNickname
		}
		if resp.SenderFaceURL != nil {
			msg.MsgData.SenderFaceURL = *resp.SenderFaceURL
		}
		if resp.SessionType != nil {
			msg.MsgData.SessionType = *resp.SessionType
		}
		if resp.MsgFrom != nil {
			msg.MsgData.MsgFrom = *resp.MsgFrom
		}
		if resp.ContentType != nil {
			msg.MsgData.ContentType = *resp.ContentType
		}
		if resp.Status != nil {
			msg.MsgData.Status = *resp.Status
		}
		if resp.Options != nil {
			msg.MsgData.Options = *resp.Options
		}
		if resp.OfflinePushInfo != nil {
			msg.MsgData.OfflinePushInfo = resp.OfflinePushInfo
		}
		if resp.AtUserIDList != nil {
			msg.MsgData.AtUserIDList = *resp.AtUserIDList
		}
		if resp.MsgDataList != nil {
			msg.MsgData.MsgDataList = *resp.MsgDataList
		}
		if resp.AttachedInfo != nil {
			msg.MsgData.AttachedInfo = *resp.AttachedInfo
		}
		if resp.Ex != nil {
			msg.MsgData.Ex = *resp.Ex
		}

	}
	log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), string(msg.MsgData.Content))
	return callbackResp
}
