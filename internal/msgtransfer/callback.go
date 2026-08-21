package msgtransfer

import (
	"context"

	"google.golang.org/protobuf/proto"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/webhook"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/mcontext"

	cbapi "github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
)

func toCommonCallback(ctx context.Context, msg *sdkws.MsgData, command string) cbapi.CommonCallbackReq {
	return cbapi.CommonCallbackReq{
		SendID:           msg.SendID,
		ServerMsgID:      msg.ServerMsgID,
		CallbackCommand:  command,
		ClientMsgID:      msg.ClientMsgID,
		OperationID:      mcontext.GetOperationID(ctx),
		SenderPlatformID: msg.SenderPlatformID,
		SenderNickname:   msg.SenderNickname,
		SessionType:      msg.SessionType,
		MsgFrom:          msg.MsgFrom,
		ContentType:      msg.ContentType,
		Status:           msg.Status,
		SendTime:         msg.SendTime,
		CreateTime:       msg.CreateTime,
		AtUserIDList:     msg.AtUserIDList,
		SenderFaceURL:    msg.SenderFaceURL,
		Content:          GetContent(msg),
		Seq:              uint32(msg.Seq),
		Ex:               msg.Ex,
	}
}

func GetContent(msg *sdkws.MsgData) string {
	if msg.ContentType >= constant.NotificationBegin && msg.ContentType <= constant.NotificationEnd {
		var tips sdkws.TipsComm
		_ = proto.Unmarshal(msg.Content, &tips)
		content := tips.JsonDetail
		return content
	} else {
		return string(msg.Content)
	}
}

func (mc *OnlineHistoryMongoConsumerHandler) webhookAfterMsgSaveDB(ctx context.Context, after *config.AfterConfig, msg *sdkws.MsgData) {
	target := msg.RecvID
	if msg.SessionType == constant.ReadGroupChatType {
		target = msg.GroupID
	}
	if !webhook.FilterAfterMsg(msg, after, target) {
		return
	}

	cbReq := &cbapi.CallbackAfterMsgSaveDBReq{
		CommonCallbackReq: toCommonCallback(ctx, msg, cbapi.CallbackAfterMsgSaveDBCommand),
	}

	switch msg.SessionType {
	case constant.SingleChatType, constant.NotificationChatType:
		cbReq.RecvID = msg.RecvID
	case constant.ReadGroupChatType:
		cbReq.GroupID = msg.GroupID
	default:
	}

	mc.webhookClient.AsyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, &cbapi.CallbackAfterMsgSaveDBResp{}, after)
}
