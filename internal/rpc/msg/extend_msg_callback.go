package msg

import (
	cbapi "OpenIM/pkg/callbackstruct"
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/http"
	"OpenIM/pkg/common/tracelog"
	"OpenIM/pkg/proto/msg"
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
		OpUserID:              tracelog.GetOpUserID(ctx),
		SessionType:           setReq.SessionType,
		ReactionExtensionList: setReq.ReactionExtensions,
		ClientMsgID:           setReq.ClientMsgID,
		IsReact:               setReq.IsReact,
		IsExternalExtensions:  setReq.IsExternalExtensions,
		MsgFirstModifyTime:    setReq.MsgFirstModifyTime,
	}
	resp := &cbapi.CallbackBeforeSetMessageReactionExtResp{}
	if err := http.CallBackPostReturn(cbURL(), req, resp, config.Config.Callback.CallbackAfterSendGroupMsg); err != nil {
		return err
	}
	setReq.MsgFirstModifyTime = resp.MsgFirstModifyTime
	return nil
}

func CallbackDeleteMessageReactionExtensions(setReq *msg.DeleteMessagesReactionExtensionsReq) error {
	if !config.Config.Callback.CallbackAfterSendGroupMsg.Enable {
		return nil
	}
	req := &cbapi.CallbackDeleteMessageReactionExtReq{
		OperationID:           setReq.OperationID,
		CallbackCommand:       constant.CallbackBeforeDeleteMessageReactionExtensionsCommand,
		SourceID:              setReq.SourceID,
		OpUserID:              setReq.OpUserID,
		SessionType:           setReq.SessionType,
		ReactionExtensionList: setReq.ReactionExtensions,
		ClientMsgID:           setReq.ClientMsgID,
		IsExternalExtensions:  setReq.IsExternalExtensions,
		MsgFirstModifyTime:    setReq.MsgFirstModifyTime,
	}
	resp := &cbapi.CallbackDeleteMessageReactionExtResp{}
	return http.CallBackPostReturn(cbURL(), req, resp, config.Config.Callback.CallbackAfterSendGroupMsg)
}

func CallbackGetMessageListReactionExtensions(ctx context.Context, getReq *msg.GetMessagesReactionExtensionsReq) error {
	if !config.Config.Callback.CallbackAfterSendGroupMsg.Enable {
		return nil
	}
	req := &cbapi.CallbackGetMessageListReactionExtReq{
		OperationID:     tracelog.GetOperationID(ctx),
		CallbackCommand: constant.CallbackGetMessageListReactionExtensionsCommand,
		SourceID:        getReq.SourceID,
		OpUserID:        tracelog.GetOperationID(ctx),
		SessionType:     getReq.SessionType,
		TypeKeyList:     getReq.TypeKeys,
		MessageKeyList:  getReq.MessageReactionKeys,
	}
	resp := &cbapi.CallbackGetMessageListReactionExtResp{}
	return http.CallBackPostReturn(cbURL(), req, resp, config.Config.Callback.CallbackAfterSendGroupMsg)
}

func CallbackAddMessageReactionExtensions(ctx context.Context, setReq *msg.ModifyMessageReactionExtensionsReq) error {
	req := &cbapi.CallbackAddMessageReactionExtReq{
		OperationID:           tracelog.GetOperationID(ctx),
		CallbackCommand:       constant.CallbackAddMessageListReactionExtensionsCommand,
		SourceID:              setReq.SourceID,
		OpUserID:              tracelog.GetOperationID(ctx),
		SessionType:           setReq.SessionType,
		ReactionExtensionList: setReq.ReactionExtensions,
		ClientMsgID:           setReq.ClientMsgID,
		IsReact:               setReq.IsReact,
		IsExternalExtensions:  setReq.IsExternalExtensions,
		MsgFirstModifyTime:    setReq.MsgFirstModifyTime,
	}
	resp := &cbapi.CallbackAddMessageReactionExtResp{}
	return http.CallBackPostReturn(cbURL(), req, resp, config.Config.Callback.CallbackAfterSendGroupMsg)
}
