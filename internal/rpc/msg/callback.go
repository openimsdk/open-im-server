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

func callbackBeforeSendSingleMsg(msg *pbChat.SendMsgReq) (canSend bool, err error) {
	if !config.Config.Callback.CallbackbeforeSendSingleMsg.Enable {
		return true, nil
	}
	log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), msg)
	req := cbApi.CallbackBeforeSendSingleMsgReq{CommonCallbackReq:cbApi.CommonCallbackReq{
		CallbackCommand: constant.CallbackBeforeSendSingleMsgCommand,
	}}
	resp := &cbApi.CallbackBeforeSendSingleMsgResp{CommonCallbackResp:cbApi.CommonCallbackResp{
	}}
	defer log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), req, *resp)
	utils.CopyStructFields(req, msg.MsgData)
	req.Content = string(msg.MsgData.Content)
	if err := http.PostReturn(config.Config.Callback.CallbackUrl, req, resp, config.Config.Callback.CallbackbeforeSendSingleMsg.CallbackTimeOut); err != nil{
		if !config.Config.Callback.CallbackbeforeSendSingleMsg.CallbackFailedContinue {
			return false, err
		}
	} else {
		if resp.ActionCode == constant.ActionForbidden && resp.ErrCode == constant.CallbackHandleSuccess {
			return false, nil
		}
	}
	return true, nil
}


func callbackAfterSendSingleMsg(msg *pbChat.SendMsgReq) error {
	if !config.Config.Callback.CallbackAfterSendSingleMsg.Enable {
		return nil
	}
	log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), msg)
	req := cbApi.CallbackAfterSendSingleMsgReq{CommonCallbackReq: cbApi.CommonCallbackReq{CallbackCommand:constant.CallbackAfterSendSingleMsgCommand}}
	resp := &cbApi.CallbackAfterSendSingleMsgResp{CommonCallbackResp: cbApi.CommonCallbackResp{}}
	defer log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), req, *resp)
	utils.CopyStructFields(req, msg.MsgData)
	req.Content = string(msg.MsgData.Content)
	if err := http.PostReturn(config.Config.Callback.CallbackUrl, req, resp, config.Config.Callback.CallbackAfterSendSingleMsg.CallbackTimeOut); err != nil{
		return err
	}
	return nil
}


func callbackBeforeSendGroupMsg(msg *pbChat.SendMsgReq) (canSend bool, err error) {
	if !config.Config.Callback.CallbackBeforeSendGroupMsg.Enable {
		return true, nil
	}
	log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), msg)
	req := cbApi.CallbackBeforeSendSingleMsgReq{CommonCallbackReq: cbApi.CommonCallbackReq{CallbackCommand:constant.CallbackBeforeSendGroupMsgCommand}}
	resp := &cbApi.CallbackBeforeSendGroupMsgResp{CommonCallbackResp: cbApi.CommonCallbackResp{}}
	defer log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), req, *resp)
	utils.CopyStructFields(req, msg.MsgData)
	req.Content = string(msg.MsgData.Content)
	if err := http.PostReturn(config.Config.Callback.CallbackUrl, req, resp, config.Config.Callback.CallbackBeforeSendGroupMsg.CallbackTimeOut); err != nil {
		if !config.Config.Callback.CallbackBeforeSendGroupMsg.CallbackFailedContinue {
			return false, nil
		}
	} else {
		if resp.ActionCode == constant.ActionForbidden && resp.ErrCode == constant.CallbackHandleSuccess {
			return false, nil
		}
	}
	return true, nil
}

func callbackAfterSendGroupMsg(msg *pbChat.SendMsgReq) error {
	if !config.Config.Callback.CallbackAfterSendGroupMsg.Enable {
		return nil
	}
	log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), msg)
	req := cbApi.CallbackAfterSendGroupMsgReq{CommonCallbackReq: cbApi.CommonCallbackReq{CallbackCommand:constant.CallbackAfterSendGroupMsgCommand}}
	resp := &cbApi.CallbackAfterSendGroupMsgResp{CommonCallbackResp: cbApi.CommonCallbackResp{}}
	defer log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), req, *resp)
	utils.CopyStructFields(req, msg.MsgData)
	req.Content = string(msg.MsgData.Content)
	if err := http.PostReturn(config.Config.Callback.CallbackUrl, req, resp, config.Config.Callback.CallbackAfterSendGroupMsg.CallbackTimeOut); err != nil {
		return err
	}
	return nil
}


func callBackWordFilter(msg *pbChat.SendMsgReq) (canSend bool, err error) {
	if !config.Config.Callback.CallbackWordFilter.Enable || msg.MsgData.ContentType != constant.Text {
		return true, nil
	}
	log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), msg)
	req := cbApi.CallbackWordFilterReq{CommonCallbackReq: cbApi.CommonCallbackReq{CallbackCommand:constant.CallbackWordFilterCommand}}
	resp := &cbApi.CallbackWordFilterResp{CommonCallbackResp: cbApi.CommonCallbackResp{}}
	utils.CopyStructFields(&req, msg.MsgData)
	req.Content = string(msg.MsgData.Content)
	defer log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), req, *resp)
	if err := http.PostReturn(config.Config.Callback.CallbackUrl, req, resp, config.Config.Callback.CallbackWordFilter.CallbackTimeOut); err != nil {
		if !config.Config.Callback.CallbackWordFilter.CallbackFailedContinue {
			log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), "config disable, stop this operation")
			return false, err
		}
	} else {
		if resp.ActionCode == constant.ActionForbidden && resp.ErrCode == constant.CallbackHandleSuccess {
			return false, nil
		}
		msg.MsgData.Content = []byte(resp.Content)
		log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), msg.MsgData.Content)
	}
	return true, nil
}





