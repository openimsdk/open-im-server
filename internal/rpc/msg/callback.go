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
	log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), msg)
	if !config.Config.Callback.CallbackbeforeSendSingleMsg.Enable || msg.MsgData.ContentType != constant.Text {
		return true, nil
	}
	req := cbApi.CallbackBeforeSendSingleMsgReq{CommonCallbackReq:cbApi.CommonCallbackReq{
	}}
	resp := &cbApi.CallbackBeforeSendSingleMsgResp{CommonCallbackResp:cbApi.CommonCallbackResp{
	}}
	utils.CopyStructFields(req, msg.MsgData)
	if err := http.PostReturn(config.Config.Callback.CallbackUrl, req, resp, config.Config.Callback.CallbackbeforeSendSingleMsg.CallbackTimeOut); err != nil && !config.Config.Callback.CallbackbeforeSendSingleMsg.CallbackFailedContinue{
		return false, err
	}
	if resp.ActionCode == constant.ActionForbidden {
		return false, nil
	}
	return true, nil
}


func callbackAfterSendSingleMsg(msg *pbChat.SendMsgReq) error {
	log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), msg)
	if !config.Config.Callback.CallbackAfterSendSingleMsg.Enable || msg.MsgData.ContentType != constant.Text {
		return nil
	}
	req := cbApi.CallbackAfterSendSingleMsgReq{CommonCallbackReq: cbApi.CommonCallbackReq{}}
	resp := &cbApi.CallbackAfterSendSingleMsgResp{CommonCallbackResp: cbApi.CommonCallbackResp{}}
	utils.CopyStructFields(req, msg.MsgData)
	if err := http.PostReturn(config.Config.Callback.CallbackUrl, req, resp, config.Config.Callback.CallbackAfterSendSingleMsg.CallbackTimeOut); err != nil && config.Config.Callback.CallbackAfterSendSingleMsg.CallbackFailedContinue{
		return err
	}
	return nil
}


func callbackBeforeSendGroupMsg(msg *pbChat.SendMsgReq) (canSend bool, err error) {
	log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), msg)
	if !config.Config.Callback.CallbackBeforeSendGroupMsg.Enable || msg.MsgData.ContentType != constant.Text {
		return true, nil
	}
	req := cbApi.CallbackBeforeSendSingleMsgReq{CommonCallbackReq: cbApi.CommonCallbackReq{}}
	resp := &cbApi.CallbackBeforeSendGroupMsgResp{CommonCallbackResp: cbApi.CommonCallbackResp{}}
	utils.CopyStructFields(req, msg.MsgData)
	if err := http.PostReturn(config.Config.Callback.CallbackUrl, req, resp, config.Config.Callback.CallbackBeforeSendGroupMsg.CallbackTimeOut); err != nil && !config.Config.Callback.CallbackBeforeSendGroupMsg.CallbackFailedContinue{
		return false, nil
	}
	if resp.ActionCode == constant.ActionForbidden {
		return false, nil
	}
	return true, nil
}

func callbackAfterSendGroupMsg(msg *pbChat.SendMsgReq) error {
	log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), msg)
	if !config.Config.Callback.CallbackAfterSendGroupMsg.Enable || msg.MsgData.ContentType != constant.Text {
		return nil
	}
	return nil
}


func callBackWordFilter(msg *pbChat.SendMsgReq) (canSend bool, err error) {
	log.NewDebug(msg.OperationID, utils.GetSelfFuncName(), msg)
	if !config.Config.Callback.CallbackWordFilter.Enable || msg.MsgData.ContentType != constant.Text {
		return true, nil
	}
	req := cbApi.CallBackWordFilterReq{
		CommonCallbackReq: cbApi.CommonCallbackReq{},
	}
	resp := cbApi.CallBackWordFilterResp{CommonCallbackResp: cbApi.CommonCallbackResp{}}
	utils.CopyStructFields(&req, msg.MsgData)
	if err := http.PostReturn(msg.OperationID, req, &resp, config.Config.Callback.CallbackWordFilter.CallbackTimeOut); err != nil  {
		return false, err
	}
	msg.MsgData.Content = resp.Content
	return true, nil
}





