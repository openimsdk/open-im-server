package notification

import (
	"OpenIM/pkg/apistruct"
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/tracelog"
	"OpenIM/pkg/proto/msg"
	sdkws "OpenIM/pkg/proto/sdkws"
	"OpenIM/pkg/utils"
	"context"
)

func (c *Check) ExtendMessageUpdatedNotification(ctx context.Context, sendID string, sourceID string, sessionType int32,
	req *msg.SetMessageReactionExtensionsReq, resp *msg.SetMessageReactionExtensionsResp, isHistory bool, isReactionFromCache bool) {
	var m apistruct.ReactionMessageModifierNotification
	m.SourceID = req.SourceID
	m.OpUserID = tracelog.GetOpUserID(ctx)
	m.SessionType = req.SessionType
	keyMap := make(map[string]*sdkws.KeyValue)
	for _, valueResp := range resp.Result {
		if valueResp.ErrCode == 0 {
			keyMap[valueResp.KeyValue.TypeKey] = valueResp.KeyValue
		}
	}
	if len(keyMap) == 0 {
		return
	}
	m.SuccessReactionExtensions = keyMap
	m.ClientMsgID = req.ClientMsgID
	m.IsReact = resp.IsReact
	m.IsExternalExtensions = req.IsExternalExtensions
	m.MsgFirstModifyTime = resp.MsgFirstModifyTime
	c.messageReactionSender(ctx, sendID, sourceID, sessionType, constant.ReactionMessageModifier, utils.StructToJsonString(m), isHistory, isReactionFromCache)
}
func (c *Check) ExtendMessageDeleteNotification(ctx context.Context, sendID string, sourceID string, sessionType int32,
	req *msg.DeleteMessagesReactionExtensionsReq, resp *msg.DeleteMessagesReactionExtensionsResp, isHistory bool, isReactionFromCache bool) {
	var m apistruct.ReactionMessageDeleteNotification
	m.SourceID = req.SourceID
	m.OpUserID = req.OpUserID
	m.SessionType = req.SessionType
	keyMap := make(map[string]*sdkws.KeyValue)
	for _, valueResp := range resp.Result {
		if valueResp.ErrCode == 0 {
			keyMap[valueResp.KeyValue.TypeKey] = valueResp.KeyValue
		}
	}
	if len(keyMap) == 0 {
		return
	}
	m.SuccessReactionExtensions = keyMap
	m.ClientMsgID = req.ClientMsgID
	m.MsgFirstModifyTime = req.MsgFirstModifyTime

	c.messageReactionSender(ctx, sendID, sourceID, sessionType, constant.ReactionMessageDeleter, utils.StructToJsonString(m), isHistory, isReactionFromCache)
}
func (c *Check) messageReactionSender(ctx context.Context, sendID string, sourceID string, sessionType, contentType int32, content string, isHistory bool, isReactionFromCache bool) error {
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
	pbData := msg.SendMsgReq{
		MsgData: &sdkws.MsgData{
			SendID:      sendID,
			ClientMsgID: utils.GetMsgID(sendID),
			SessionType: sessionType,
			MsgFrom:     constant.SysMsgType,
			ContentType: contentType,
			Content:     []byte(content),
			CreateTime:  utils.GetCurrentTimestampByMill(),
			Options:     options,
		},
	}
	switch sessionType {
	case constant.SingleChatType, constant.NotificationChatType:
		pbData.MsgData.RecvID = sourceID
	case constant.GroupChatType, constant.SuperGroupChatType:
		pbData.MsgData.GroupID = sourceID
	}
	_, err := c.Msg.SendMsg(ctx, &pbData)
	return err
}
