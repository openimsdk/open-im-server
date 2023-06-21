package notification

import (
	"context"
	"encoding/json"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"google.golang.org/protobuf/proto"
)

type ExtendMsgNotificationSender struct {
	*rpcclient.MessageRpcClient
}

func NewExtendMsgNotificationSender(client discoveryregistry.SvcDiscoveryRegistry) *ExtendMsgNotificationSender {
	return &ExtendMsgNotificationSender{}
}

func (e *ExtendMsgNotificationSender) ExtendMessageUpdatedNotification(ctx context.Context, sendID string, conversationID string, sessionType int32,
	req *msg.SetMessageReactionExtensionsReq, resp *msg.SetMessageReactionExtensionsResp, isHistory bool, isReactionFromCache bool) {
	var content sdkws.ReactionMessageModifierNotification
	content.ConversationID = req.ConversationID
	content.OpUserID = mcontext.GetOpUserID(ctx)
	content.SessionType = req.SessionType
	keyMap := make(map[string]*sdkws.KeyValue)
	for _, valueResp := range resp.Result {
		if valueResp.ErrCode == 0 {
			keyMap[valueResp.KeyValue.TypeKey] = valueResp.KeyValue
		}
	}
	if len(keyMap) == 0 {
		return
	}
	content.SuccessReactionExtensions = keyMap
	content.ClientMsgID = req.ClientMsgID
	content.IsReact = resp.IsReact
	content.IsExternalExtensions = req.IsExternalExtensions
	content.MsgFirstModifyTime = resp.MsgFirstModifyTime
	e.messageReactionSender(ctx, sendID, conversationID, sessionType, constant.ReactionMessageModifier, &content, isHistory, isReactionFromCache)
}
func (e *ExtendMsgNotificationSender) ExtendMessageDeleteNotification(ctx context.Context, sendID string, conversationID string, sessionType int32,
	req *msg.DeleteMessagesReactionExtensionsReq, resp *msg.DeleteMessagesReactionExtensionsResp, isHistory bool, isReactionFromCache bool) {
	var content sdkws.ReactionMessageDeleteNotification
	content.ConversationID = req.ConversationID
	content.OpUserID = req.OpUserID
	content.SessionType = req.SessionType
	keyMap := make(map[string]*sdkws.KeyValue)
	for _, valueResp := range resp.Result {
		if valueResp.ErrCode == 0 {
			keyMap[valueResp.KeyValue.TypeKey] = valueResp.KeyValue
		}
	}
	if len(keyMap) == 0 {
		return
	}
	content.SuccessReactionExtensions = keyMap
	content.ClientMsgID = req.ClientMsgID
	content.MsgFirstModifyTime = req.MsgFirstModifyTime
	e.messageReactionSender(ctx, sendID, conversationID, sessionType, constant.ReactionMessageDeleter, &content, isHistory, isReactionFromCache)
}
func (e *ExtendMsgNotificationSender) messageReactionSender(ctx context.Context, sendID string, conversationID string, sessionType, contentType int32, m proto.Message, isHistory bool, isReactionFromCache bool) error {
	options := make(map[string]bool, 5)
	utils.SetSwitchFromOptions(options, constant.IsOfflinePush, false)
	utils.SetSwitchFromOptions(options, constant.IsConversationUpdate, false)
	utils.SetSwitchFromOptions(options, constant.IsSenderConversationUpdate, false)
	utils.SetSwitchFromOptions(options, constant.IsUnreadCount, false)
	utils.SetSwitchFromOptions(options, constant.IsReactionFromCache, isReactionFromCache)
	if !isHistory {
		utils.SetSwitchFromOptions(options, constant.IsHistory, false)
		utils.SetSwitchFromOptions(options, constant.IsPersistent, false)
	}
	bytes, err := json.Marshal(m)
	if err != nil {
		return errs.ErrData.Wrap(err.Error())
	}
	pbData := msg.SendMsgReq{
		MsgData: &sdkws.MsgData{
			SendID:      sendID,
			ClientMsgID: utils.GetMsgID(sendID),
			SessionType: sessionType,
			MsgFrom:     constant.SysMsgType,
			ContentType: contentType,
			Content:     bytes,
			CreateTime:  utils.GetCurrentTimestampByMill(),
			Options:     options,
		},
	}
	switch sessionType {
	case constant.SingleChatType, constant.NotificationChatType:
		pbData.MsgData.RecvID = conversationID
	case constant.GroupChatType, constant.SuperGroupChatType:
		pbData.MsgData.GroupID = conversationID
	}
	_, err = e.SendMsg(ctx, &pbData)
	return err
}
