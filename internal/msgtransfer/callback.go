package msgtransfer

import (
	"context"
	"encoding/base64"

	"github.com/openimsdk/open-im-server/v3/pkg/apistruct"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/webhook"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/openimsdk/tools/utils/stringutil"
	"google.golang.org/protobuf/proto"

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
	if msg.ContentType == constant.Typing {
		return
	}

	if !filterAfterMsg(msg, after) {
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

	mc.webhookClient.AsyncPostWithQuery(ctx, cbReq.GetCallbackCommand(), cbReq, &cbapi.CallbackAfterMsgSaveDBResp{}, after, buildKeyMsgDataQuery(msg))
}

func buildKeyMsgDataQuery(msg *sdkws.MsgData) map[string]string {
	keyMsgData := apistruct.KeyMsgData{
		SendID:  msg.SendID,
		RecvID:  msg.RecvID,
		GroupID: msg.GroupID,
	}

	return map[string]string{
		webhook.Key: base64.StdEncoding.EncodeToString(stringutil.StructToJsonBytes(keyMsgData)),
	}
}

func filterAfterMsg(msg *sdkws.MsgData, after *config.AfterConfig) bool {
	return filterMsg(msg, after.AttentionIds, after.DeniedTypes)
}

func filterMsg(msg *sdkws.MsgData, attentionIds []string, deniedTypes []int32) bool {
	// According to the attentionIds configuration, only some users are sent
	if len(attentionIds) != 0 && msg.ContentType == constant.SingleChatType && !datautil.Contain(msg.RecvID, attentionIds...) {
		return false
	}

	if len(attentionIds) != 0 && msg.ContentType == constant.ReadGroupChatType && !datautil.Contain(msg.GroupID, attentionIds...) {
		return false
	}

	if defaultDeniedTypes(msg.ContentType) {
		return false
	}

	if len(deniedTypes) != 0 && datautil.Contain(msg.ContentType, deniedTypes...) {
		return false
	}

	return true
}

func defaultDeniedTypes(contentType int32) bool {
	if contentType >= constant.NotificationBegin && contentType <= constant.NotificationEnd {
		return true
	}
	if contentType == constant.Typing {
		return true
	}
	return false
}
