package msg

import (
	cbApi "Open_IM/pkg/call_back_struct"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/http"
	"Open_IM/pkg/common/log"
	pbChat "Open_IM/pkg/proto/chat"
	"Open_IM/pkg/utils"
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
	}
}

func callbackBeforeSendSingleMsg(msg *pbChat.SendMsgReq) (canSend bool, err error) {
	if !config.Config.Callback.CallbackBeforeSendSingleMsg.Enable {
		return true, nil
	}
	log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), msg)
	commonCallbackReq := copyCallbackCommonReqStruct(msg)
	commonCallbackReq.CallbackCommand = constant.CallbackBeforeSendSingleMsgCommand
	req := cbApi.CallbackBeforeSendSingleMsgReq{
		CommonCallbackReq: commonCallbackReq,
		RecvID:            msg.MsgData.RecvID,
	}
	resp := &cbApi.CallbackBeforeSendSingleMsgResp{
		CommonCallbackResp: cbApi.CommonCallbackResp{},
	}
	//utils.CopyStructFields(req, msg.MsgData)
	defer log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), req, *resp)
	if err := http.PostReturn(config.Config.Callback.CallbackUrl, req, resp, config.Config.Callback.CallbackBeforeSendSingleMsg.CallbackTimeOut); err != nil {
		if !config.Config.Callback.CallbackBeforeSendSingleMsg.CallbackFailedContinue {
			return false, err
		} else {
			return true, err
		}
	} else {
		if resp.ActionCode == constant.ActionForbidden && resp.ErrCode == constant.CallbackHandleSuccess {
			return false, nil
		}
	}
	return true, err
}

func callbackAfterSendSingleMsg(msg *pbChat.SendMsgReq) error {
	if !config.Config.Callback.CallbackAfterSendSingleMsg.Enable {
		return nil
	}
	log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), msg)
	commonCallbackReq := copyCallbackCommonReqStruct(msg)
	commonCallbackReq.CallbackCommand = constant.CallbackAfterSendSingleMsgCommand
	req := cbApi.CallbackAfterSendSingleMsgReq{
		CommonCallbackReq: commonCallbackReq,
		RecvID:            msg.MsgData.RecvID,
	}
	resp := &cbApi.CallbackAfterSendSingleMsgResp{CommonCallbackResp: cbApi.CommonCallbackResp{}}
	//utils.CopyStructFields(req, msg.MsgData)
	defer log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), req, *resp)
	if err := http.PostReturn(config.Config.Callback.CallbackUrl, req, resp, config.Config.Callback.CallbackAfterSendSingleMsg.CallbackTimeOut); err != nil {
		return err
	}
	return nil
}

func callbackBeforeSendGroupMsg(msg *pbChat.SendMsgReq) (canSend bool, err error) {
	if !config.Config.Callback.CallbackBeforeSendGroupMsg.Enable {
		return true, nil
	}
	log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), msg)
	commonCallbackReq := copyCallbackCommonReqStruct(msg)
	commonCallbackReq.CallbackCommand = constant.CallbackBeforeSendGroupMsgCommand
	req := cbApi.CallbackAfterSendGroupMsgReq{
		CommonCallbackReq: commonCallbackReq,
		GroupID:           msg.MsgData.GroupID,
	}
	resp := &cbApi.CallbackBeforeSendGroupMsgResp{CommonCallbackResp: cbApi.CommonCallbackResp{}}
	//utils.CopyStructFields(req, msg.MsgData)
	defer log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), req, *resp)
	if err := http.PostReturn(config.Config.Callback.CallbackUrl, req, resp, config.Config.Callback.CallbackBeforeSendGroupMsg.CallbackTimeOut); err != nil {
		if !config.Config.Callback.CallbackBeforeSendGroupMsg.CallbackFailedContinue {
			return false, err
		} else {
			return true, err
		}
	} else {
		if resp.ActionCode == constant.ActionForbidden && resp.ErrCode == constant.CallbackHandleSuccess {
			return false, nil
		}
	}
	return true, err
}

func callbackAfterSendGroupMsg(msg *pbChat.SendMsgReq) error {
	if !config.Config.Callback.CallbackAfterSendGroupMsg.Enable {
		return nil
	}
	log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), msg)
	commonCallbackReq := copyCallbackCommonReqStruct(msg)
	commonCallbackReq.CallbackCommand = constant.CallbackAfterSendGroupMsgCommand
	req := cbApi.CallbackAfterSendGroupMsgReq{
		CommonCallbackReq: commonCallbackReq,
		GroupID:           msg.MsgData.GroupID,
	}
	resp := &cbApi.CallbackAfterSendGroupMsgResp{CommonCallbackResp: cbApi.CommonCallbackResp{}}

	//utils.CopyStructFields(req, msg.MsgData)
	defer log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), req, *resp)
	if err := http.PostReturn(config.Config.Callback.CallbackUrl, req, resp, config.Config.Callback.CallbackAfterSendGroupMsg.CallbackTimeOut); err != nil {
		return err
	}
	return nil
}

func callbackWordFilter(msg *pbChat.SendMsgReq) (canSend bool, err error) {
	if !config.Config.Callback.CallbackWordFilter.Enable || msg.MsgData.ContentType != constant.Text {
		return true, nil
	}
	log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), msg)
	commonCallbackReq := copyCallbackCommonReqStruct(msg)
	commonCallbackReq.CallbackCommand = constant.CallbackWordFilterCommand
	req := cbApi.CallbackWordFilterReq{
		CommonCallbackReq: commonCallbackReq,
	}
	resp := &cbApi.CallbackWordFilterResp{CommonCallbackResp: cbApi.CommonCallbackResp{}}
	//utils.CopyStructFields(&req., msg.MsgData)
	defer log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), req, *resp)
	if err := http.PostReturn(config.Config.Callback.CallbackUrl, req, resp, config.Config.Callback.CallbackWordFilter.CallbackTimeOut); err != nil {
		if !config.Config.Callback.CallbackWordFilter.CallbackFailedContinue {
			log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), "callback failed and config disable, stop this operation")
			return false, err
		} else {
			return true, err
		}
	} else {
		if resp.ActionCode == constant.ActionForbidden && resp.ErrCode == constant.CallbackHandleSuccess {
			return false, nil
		}
		if resp.ErrCode == constant.CallbackHandleSuccess {
			msg.MsgData.Content = []byte(resp.Content)
		}
		log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), string(msg.MsgData.Content))
	}
	return true, err
}
