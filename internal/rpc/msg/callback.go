package msg

import (
	cbApi "Open_IM/pkg/call_back_struct"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/http"
	"Open_IM/pkg/common/log"
	pbChat "Open_IM/pkg/proto/msg"
	"Open_IM/pkg/utils"
	http2 "net/http"
)

func copyCallbackCommonReqStruct(msg *pbChat.SendMsgReq) cbApi.CommonCallbackReq {
	return cbApi.CommonCallbackReq{
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
		Content:          string(msg.MsgData.Content),
		AtUserIDList:     msg.MsgData.AtUserIDList,
		SenderFaceURL:    msg.MsgData.SenderFaceURL,
	}
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
	if err := http.PostReturn(config.Config.Callback.CallbackUrl, req, resp, config.Config.Callback.CallbackBeforeSendSingleMsg.CallbackTimeOut); err != nil {
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
	if err := http.PostReturn(config.Config.Callback.CallbackUrl, req, resp, config.Config.Callback.CallbackAfterSendSingleMsg.CallbackTimeOut); err != nil {
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
	if err := http.PostReturn(config.Config.Callback.CallbackUrl, req, resp, config.Config.Callback.CallbackBeforeSendGroupMsg.CallbackTimeOut); err != nil {
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
	if err := http.PostReturn(config.Config.Callback.CallbackUrl, req, resp, config.Config.Callback.CallbackAfterSendGroupMsg.CallbackTimeOut); err != nil {
		callbackResp.ErrCode = http2.StatusInternalServerError
		callbackResp.ErrMsg = err.Error()
		return callbackResp
	}
	return callbackResp
}

func callbackWordFilter(msg *pbChat.SendMsgReq) cbApi.CommonCallbackResp {
	callbackResp := cbApi.CommonCallbackResp{OperationID: msg.OperationID}
	if !config.Config.Callback.CallbackWordFilter.Enable || msg.MsgData.ContentType != constant.Text {
		return callbackResp
	}
	log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), msg)
	commonCallbackReq := copyCallbackCommonReqStruct(msg)
	commonCallbackReq.CallbackCommand = constant.CallbackWordFilterCommand
	req := cbApi.CallbackWordFilterReq{
		CommonCallbackReq: commonCallbackReq,
	}
	resp := &cbApi.CallbackWordFilterResp{CommonCallbackResp: &callbackResp}
	defer log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), req, *resp)
	if err := http.PostReturn(config.Config.Callback.CallbackUrl, req, resp, config.Config.Callback.CallbackWordFilter.CallbackTimeOut); err != nil {
		callbackResp.ErrCode = http2.StatusInternalServerError
		callbackResp.ErrMsg = err.Error()
		if !config.Config.Callback.CallbackWordFilter.CallbackFailedContinue {
			callbackResp.ActionCode = constant.ActionForbidden
			return callbackResp
		} else {
			callbackResp.ActionCode = constant.ActionAllow
			return callbackResp
		}
	}
	if resp.ErrCode == constant.CallbackHandleSuccess && resp.ActionCode == constant.ActionAllow && resp.Content != "" {
		msg.MsgData.Content = []byte(resp.Content)
	}
	log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), string(msg.MsgData.Content))
	return callbackResp
}
