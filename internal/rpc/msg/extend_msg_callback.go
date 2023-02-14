package msg

import (
	cb "Open_IM/pkg/callbackstruct"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/http"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/proto/msg"
	"Open_IM/pkg/utils"
	http2 "net/http"
)

func callbackSetMessageReactionExtensions(setReq *msg.SetMessageReactionExtensionsReq) *cb.CallbackBeforeSetMessageReactionExtResp {
	callbackResp := cb.CommonCallbackResp{OperationID: setReq.OperationID}
	log.NewDebug(setReq.OperationID, utils.GetSelfFuncName(), setReq.String())
	req := cb.CallbackBeforeSetMessageReactionExtReq{
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
	resp := &cb.CallbackBeforeSetMessageReactionExtResp{CommonCallbackResp: &callbackResp}
	defer log.NewDebug(setReq.OperationID, utils.GetSelfFuncName(), req, *resp)
	if err := http.CallBackPostReturn(config.Config.Callback.CallbackUrl, constant.CallbackBeforeSetMessageReactionExtensionCommand, req, resp, config.Config.Callback.CallbackAfterSendGroupMsg.CallbackTimeOut); err != nil {
		callbackResp.ErrCode = http2.StatusInternalServerError
		callbackResp.ErrMsg = err.Error()
	}
	return resp

}

func callbackDeleteMessageReactionExtensions(setReq *msg.DeleteMessageListReactionExtensionsReq) *cb.CallbackDeleteMessageReactionExtResp {
	callbackResp := cb.CommonCallbackResp{OperationID: setReq.OperationID}
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
func callbackGetMessageListReactionExtensions(getReq *msg.GetMessageListReactionExtensionsReq) *cb.CallbackGetMessageListReactionExtResp {
	callbackResp := cb.CommonCallbackResp{OperationID: getReq.OperationID}
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
func callbackAddMessageReactionExtensions(setReq *msg.AddMessageReactionExtensionsReq) *cb.CallbackAddMessageReactionExtResp {
	callbackResp := cb.CommonCallbackResp{OperationID: setReq.OperationID}
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
	resp := &cb.CallbackAddMessageReactionExtResp{CommonCallbackResp: &callbackResp}
	defer log.NewDebug(setReq.OperationID, utils.GetSelfFuncName(), req, *resp, *resp.CommonCallbackResp, resp.IsReact, resp.MsgFirstModifyTime)
	if err := http.CallBackPostReturn(config.Config.Callback.CallbackUrl, constant.CallbackAddMessageListReactionExtensionsCommand, req, resp, config.Config.Callback.CallbackAfterSendGroupMsg.CallbackTimeOut); err != nil {
		callbackResp.ErrCode = http2.StatusInternalServerError
		callbackResp.ErrMsg = err.Error()
	}
	return resp

}
