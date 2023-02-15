package msg

import (
	cbapi "Open_IM/pkg/callbackstruct"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/http"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/proto/msg"
	"Open_IM/pkg/utils"
	"context"
)

func CallbackSetMessageReactionExtensions(ctx context.Context, setReq *msg.SetMessageReactionExtensionsReq) error {
	if !config.Config.Callback.CallbackAfterSendGroupMsg.Enable {
		return nil
	}
	req := &cbapi.CallbackBeforeSetMessageReactionExtReq{
		OperationID:           tracelog.GetOperationID(ctx),
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
	resp := &cbapi.CallbackBeforeSetMessageReactionExtResp{}
	return http.CallBackPostReturn(config.Config.Callback.CallbackUrl, req, resp, config.Config.Callback.CallbackAfterSendGroupMsg)
}

func CallbackDeleteMessageReactionExtensions(setReq *msg.DeleteMessageListReactionExtensionsReq) error {
	if !config.Config.Callback.CallbackAfterSendGroupMsg.Enable {
		return nil
	}
	req := &cbapi.CallbackDeleteMessageReactionExtReq{
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
	resp := &cbapi.CallbackDeleteMessageReactionExtResp{}
	defer log.NewDebug(setReq.OperationID, utils.GetSelfFuncName(), req, *resp)
	return http.CallBackPostReturn(config.Config.Callback.CallbackUrl, req, resp, config.Config.Callback.CallbackAfterSendGroupMsg)
}

//func CallbackGetMessageListReactionExtensions(getReq *msg.GetMessageListReactionExtensionsReq) error {
//	if !config.Config.Callback.CallbackAfterSendGroupMsg.Enable {
//		return nil
//	}
//	req := cbapi.CallbackGetMessageListReactionExtReq{
//		OperationID:     getReq.OperationID,
//		CallbackCommand: constant.CallbackGetMessageListReactionExtensionsCommand,
//		SourceID:        getReq.SourceID,
//		OpUserID:        getReq.OpUserID,
//		SessionType:     getReq.SessionType,
//		TypeKeyList:     getReq.TypeKeyList,
//		MessageKeyList:  getReq.MessageReactionKeyList,
//	}
//	resp := &cbApi.CallbackGetMessageListReactionExtResp{CommonCallbackResp: &callbackResp}
//	defer log.NewDebug(getReq.OperationID, utils.GetSelfFuncName(), req, *resp)
//	if err := http.CallBackPostReturn(config.Config.Callback.CallbackUrl, constant.CallbackGetMessageListReactionExtensionsCommand, req, resp, config.Config.Callback.CallbackAfterSendGroupMsg.CallbackTimeOut); err != nil {
//		callbackResp.ErrCode = http2.StatusInternalServerError
//		callbackResp.ErrMsg = err.Error()
//	}
//	return resp
//}
//func callbackAddMessageReactionExtensions(setReq *msg.AddMessageReactionExtensionsReq) *cb.CallbackAddMessageReactionExtResp {
//	callbackResp := cbapi.CommonCallbackResp{OperationID: setReq.OperationID}
//	log.NewDebug(setReq.OperationID, utils.GetSelfFuncName(), setReq.String())
//	req := cbapi.CallbackAddMessageReactionExtReq{
//		OperationID:           setReq.OperationID,
//		CallbackCommand:       constant.CallbackAddMessageListReactionExtensionsCommand,
//		SourceID:              setReq.SourceID,
//		OpUserID:              setReq.OpUserID,
//		SessionType:           setReq.SessionType,
//		ReactionExtensionList: setReq.ReactionExtensionList,
//		ClientMsgID:           setReq.ClientMsgID,
//		IsReact:               setReq.IsReact,
//		IsExternalExtensions:  setReq.IsExternalExtensions,
//		MsgFirstModifyTime:    setReq.MsgFirstModifyTime,
//	}
//	resp := &cbapi.CallbackAddMessageReactionExtResp{CommonCallbackResp: &callbackResp}
//	defer log.NewDebug(setReq.OperationID, utils.GetSelfFuncName(), req, *resp, *resp.CommonCallbackResp, resp.IsReact, resp.MsgFirstModifyTime)
//	if err := http.CallBackPostReturn(config.Config.Callback.CallbackUrl, constant.CallbackAddMessageListReactionExtensionsCommand, req, resp, config.Config.Callback.CallbackAfterSendGroupMsg.CallbackTimeOut); err != nil {
//		callbackResp.ErrCode = http2.StatusInternalServerError
//		callbackResp.ErrMsg = err.Error()
//	}
//	return resp
//
//}
