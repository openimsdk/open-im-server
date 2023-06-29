package msg

import (
	"context"

	cbapi "github.com/OpenIMSDK/Open-IM-Server/pkg/callbackstruct"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/http"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
)

func callbackSetMessageReactionExtensions(ctx context.Context, setReq *msg.SetMessageReactionExtensionsReq) error {
	if !config.Config.Callback.CallbackAfterSendGroupMsg.Enable {
		return nil
	}
	req := &cbapi.CallbackBeforeSetMessageReactionExtReq{
		OperationID:           mcontext.GetOperationID(ctx),
		CallbackCommand:       constant.CallbackBeforeSetMessageReactionExtensionCommand,
		ConversationID:        setReq.ConversationID,
		OpUserID:              mcontext.GetOpUserID(ctx),
		SessionType:           setReq.SessionType,
		ReactionExtensionList: setReq.ReactionExtensions,
		ClientMsgID:           setReq.ClientMsgID,
		IsReact:               setReq.IsReact,
		IsExternalExtensions:  setReq.IsExternalExtensions,
		MsgFirstModifyTime:    setReq.MsgFirstModifyTime,
	}
	resp := &cbapi.CallbackBeforeSetMessageReactionExtResp{}
	if err := http.CallBackPostReturn(ctx, cbURL(), req, resp, config.Config.Callback.CallbackAfterSendGroupMsg); err != nil {
		return err
	}
	setReq.MsgFirstModifyTime = resp.MsgFirstModifyTime
	return nil
}

func callbackDeleteMessageReactionExtensions(ctx context.Context, setReq *msg.DeleteMessagesReactionExtensionsReq) error {
	if !config.Config.Callback.CallbackAfterSendGroupMsg.Enable {
		return nil
	}
	req := &cbapi.CallbackDeleteMessageReactionExtReq{
		OperationID:           setReq.OperationID,
		CallbackCommand:       constant.CallbackBeforeDeleteMessageReactionExtensionsCommand,
		ConversationID:        setReq.ConversationID,
		OpUserID:              setReq.OpUserID,
		SessionType:           setReq.SessionType,
		ReactionExtensionList: setReq.ReactionExtensions,
		ClientMsgID:           setReq.ClientMsgID,
		IsExternalExtensions:  setReq.IsExternalExtensions,
		MsgFirstModifyTime:    setReq.MsgFirstModifyTime,
	}
	resp := &cbapi.CallbackDeleteMessageReactionExtResp{}
	return http.CallBackPostReturn(ctx, cbURL(), req, resp, config.Config.Callback.CallbackAfterSendGroupMsg)
}

func callbackGetMessageListReactionExtensions(ctx context.Context, getReq *msg.GetMessagesReactionExtensionsReq) error {
	if !config.Config.Callback.CallbackAfterSendGroupMsg.Enable {
		return nil
	}
	req := &cbapi.CallbackGetMessageListReactionExtReq{
		OperationID:     mcontext.GetOperationID(ctx),
		CallbackCommand: constant.CallbackGetMessageListReactionExtensionsCommand,
		ConversationID:  getReq.ConversationID,
		OpUserID:        mcontext.GetOperationID(ctx),
		SessionType:     getReq.SessionType,
		TypeKeyList:     getReq.TypeKeys,
	}
	resp := &cbapi.CallbackGetMessageListReactionExtResp{}
	return http.CallBackPostReturn(ctx, cbURL(), req, resp, config.Config.Callback.CallbackAfterSendGroupMsg)
}

func callbackAddMessageReactionExtensions(ctx context.Context, setReq *msg.ModifyMessageReactionExtensionsReq) error {
	req := &cbapi.CallbackAddMessageReactionExtReq{
		OperationID:           mcontext.GetOperationID(ctx),
		CallbackCommand:       constant.CallbackAddMessageListReactionExtensionsCommand,
		ConversationID:        setReq.ConversationID,
		OpUserID:              mcontext.GetOperationID(ctx),
		SessionType:           setReq.SessionType,
		ReactionExtensionList: setReq.ReactionExtensions,
		ClientMsgID:           setReq.ClientMsgID,
		IsReact:               setReq.IsReact,
		IsExternalExtensions:  setReq.IsExternalExtensions,
		MsgFirstModifyTime:    setReq.MsgFirstModifyTime,
	}
	resp := &cbapi.CallbackAddMessageReactionExtResp{}
	return http.CallBackPostReturn(ctx, cbURL(), req, resp, config.Config.Callback.CallbackAfterSendGroupMsg)
}
