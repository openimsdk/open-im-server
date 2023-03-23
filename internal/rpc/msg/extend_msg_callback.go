package msg

import (
	cbApi "Open_IM/pkg/call_back_struct"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/http"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/proto/msg"
	"Open_IM/pkg/utils"
	http2 "net/http"
)

func callbackSetMessageReactionExtensions(setReq *msg.SetMessageReactionExtensionsReq) *cbApi.CallbackBeforeSetMessageReactionExtResp {
	callbackResp := cbApi.CommonCallbackResp{OperationID: setReq.OperationID}
	log.NewDebug(setReq.OperationID, utils.GetSelfFuncName(), setReq.String())
	req := cbApi.CallbackBeforeSetMessageReactionExtReq{
		OperationID:           setReq.OperationID,
		CallbackCommand:       constant.CallbackBeforeSetMessageReactionExtensionCommand,
		SourceID:              setReq.SourceID,
		OpUserID:              setReq.OpUserID,
		SessionType:           setReq.SessionType,
		ReactionExtensionList: setReq.ReactionExtensionList,
		ClientMsgID:           setReq.ClientMsgID,
		IsReact:               setReq.IsReact,
		IsExternalExtensions:  setReq.IsExternalExtensions,
		MsgFirstModifyTime:    setReq.MsgFirstModifyTime,
	}
	resp := &cbApi.CallbackBeforeSetMessageReactionExtResp{CommonCallbackResp: &callbackResp}
	defer log.NewDebug(setReq.OperationID, utils.GetSelfFuncName(), req, *resp)
	if err := http.CallBackPostReturn(config.Config.Callback.CallbackUrl, constant.CallbackBeforeSetMessageReactionExtensionCommand, req, resp, config.Config.Callback.CallbackAfterSendGroupMsg.CallbackTimeOut); err != nil {
		callbackResp.ErrCode = http2.StatusInternalServerError
		callbackResp.ErrMsg = err.Error()
	}
	return resp

}

func callbackDeleteMessageReactionExtensions(setReq *msg.DeleteMessageListReactionExtensionsReq) *cbApi.CallbackDeleteMessageReactionExtResp {
	callbackResp := cbApi.CommonCallbackResp{OperationID: setReq.OperationID}
	log.NewDebug(setReq.OperationID, utils.GetSelfFuncName(), setReq.String())
	req := cbApi.CallbackDeleteMessageReactionExtReq{
		OperationID:           setReq.OperationID,
		CallbackCommand:       constant.CallbackBeforeDeleteMessageReactionExtensionsCommand,
		SourceID:              setReq.SourceID,
		OpUserID:              setReq.OpUserID,
		SessionType:           setReq.SessionType,
		ReactionExtensionList: setReq.ReactionExtensionList,
		ClientMsgID:           setReq.ClientMsgID,
		IsExternalExtensions:  setReq.IsExternalExtensions,
		MsgFirstModifyTime:    setReq.MsgFirstModifyTime,
	}
	resp := &cbApi.CallbackDeleteMessageReactionExtResp{CommonCallbackResp: &callbackResp}
	defer log.NewDebug(setReq.OperationID, utils.GetSelfFuncName(), req, *resp)
	if err := http.CallBackPostReturn(config.Config.Callback.CallbackUrl, constant.CallbackBeforeDeleteMessageReactionExtensionsCommand, req, resp, config.Config.Callback.CallbackAfterSendGroupMsg.CallbackTimeOut); err != nil {
		callbackResp.ErrCode = http2.StatusInternalServerError
		callbackResp.ErrMsg = err.Error()
	}
	return resp
}
func callbackGetMessageListReactionExtensions(getReq *msg.GetMessageListReactionExtensionsReq) *cbApi.CallbackGetMessageListReactionExtResp {
	callbackResp := cbApi.CommonCallbackResp{OperationID: getReq.OperationID}
	log.NewDebug(getReq.OperationID, utils.GetSelfFuncName(), getReq.String())
	req := cbApi.CallbackGetMessageListReactionExtReq{
		OperationID:     getReq.OperationID,
		CallbackCommand: constant.CallbackGetMessageListReactionExtensionsCommand,
		SourceID:        getReq.SourceID,
		OpUserID:        getReq.OpUserID,
		SessionType:     getReq.SessionType,
		TypeKeyList:     getReq.TypeKeyList,
		MessageKeyList:  getReq.MessageReactionKeyList,
	}
	resp := &cbApi.CallbackGetMessageListReactionExtResp{CommonCallbackResp: &callbackResp}
	defer log.NewDebug(getReq.OperationID, utils.GetSelfFuncName(), req, *resp)
	if err := http.CallBackPostReturn(config.Config.Callback.CallbackUrl, constant.CallbackGetMessageListReactionExtensionsCommand, req, resp, config.Config.Callback.CallbackAfterSendGroupMsg.CallbackTimeOut); err != nil {
		callbackResp.ErrCode = http2.StatusInternalServerError
		callbackResp.ErrMsg = err.Error()
	}
	return resp
}
func callbackAddMessageReactionExtensions(setReq *msg.AddMessageReactionExtensionsReq) *cbApi.CallbackAddMessageReactionExtResp {
	callbackResp := cbApi.CommonCallbackResp{OperationID: setReq.OperationID}
	log.NewDebug(setReq.OperationID, utils.GetSelfFuncName(), setReq.String())
	req := cbApi.CallbackAddMessageReactionExtReq{
		OperationID:           setReq.OperationID,
		CallbackCommand:       constant.CallbackAddMessageListReactionExtensionsCommand,
		SourceID:              setReq.SourceID,
		OpUserID:              setReq.OpUserID,
		SessionType:           setReq.SessionType,
		ReactionExtensionList: setReq.ReactionExtensionList,
		ClientMsgID:           setReq.ClientMsgID,
		IsReact:               setReq.IsReact,
		IsExternalExtensions:  setReq.IsExternalExtensions,
		MsgFirstModifyTime:    setReq.MsgFirstModifyTime,
	}
	resp := &cbApi.CallbackAddMessageReactionExtResp{CommonCallbackResp: &callbackResp}
	defer log.NewDebug(setReq.OperationID, utils.GetSelfFuncName(), req, *resp, *resp.CommonCallbackResp, resp.IsReact, resp.MsgFirstModifyTime)
	if err := http.CallBackPostReturn(config.Config.Callback.CallbackUrl, constant.CallbackAddMessageListReactionExtensionsCommand, req, resp, config.Config.Callback.CallbackAfterSendGroupMsg.CallbackTimeOut); err != nil {
		callbackResp.ErrCode = http2.StatusInternalServerError
		callbackResp.ErrMsg = err.Error()
	}
	return resp

}
