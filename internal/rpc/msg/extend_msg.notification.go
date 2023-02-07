package msg

import (
	"Open_IM/pkg/api_struct"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/getcdv3"
	"Open_IM/pkg/proto/msg"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
)

func ExtendMessageUpdatedNotification(operationID, sendID string, sourceID string, sessionType int32,
	req *msg.SetMessageReactionExtensionsReq, resp *msg.SetMessageReactionExtensionsResp, isHistory bool, isReactionFromCache bool) {
	var m api_struct.ReactionMessageModifierNotification
	m.SourceID = req.SourceID
	m.OpUserID = req.OpUserID
	m.SessionType = req.SessionType
	keyMap := make(map[string]*open_im_sdk.KeyValue)
	for _, valueResp := range resp.Result {
		if valueResp.ErrCode == 0 {
			keyMap[valueResp.KeyValue.TypeKey] = valueResp.KeyValue
		}
	}
	if len(keyMap) == 0 {
		log.NewWarn(operationID, "all key set failed can not send notification", *req)
		return
	}
	m.SuccessReactionExtensionList = keyMap
	m.ClientMsgID = req.ClientMsgID
	m.IsReact = resp.IsReact
	m.IsExternalExtensions = req.IsExternalExtensions
	m.MsgFirstModifyTime = resp.MsgFirstModifyTime
	messageReactionSender(operationID, sendID, sourceID, sessionType, constant.ReactionMessageModifier, utils.StructToJsonString(m), isHistory, isReactionFromCache)
}
func ExtendMessageDeleteNotification(operationID, sendID string, sourceID string, sessionType int32,
	req *msg.DeleteMessageListReactionExtensionsReq, resp *msg.DeleteMessageListReactionExtensionsResp, isHistory bool, isReactionFromCache bool) {
	var m api_struct.ReactionMessageDeleteNotification
	m.SourceID = req.SourceID
	m.OpUserID = req.OpUserID
	m.SessionType = req.SessionType
	keyMap := make(map[string]*open_im_sdk.KeyValue)
	for _, valueResp := range resp.Result {
		if valueResp.ErrCode == 0 {
			keyMap[valueResp.KeyValue.TypeKey] = valueResp.KeyValue
		}
	}
	if len(keyMap) == 0 {
		log.NewWarn(operationID, "all key set failed can not send notification", *req)
		return
	}
	m.SuccessReactionExtensionList = keyMap
	m.ClientMsgID = req.ClientMsgID
	m.MsgFirstModifyTime = req.MsgFirstModifyTime

	messageReactionSender(operationID, sendID, sourceID, sessionType, constant.ReactionMessageDeleter, utils.StructToJsonString(m), isHistory, isReactionFromCache)
}
func messageReactionSender(operationID, sendID string, sourceID string, sessionType, contentType int32, content string, isHistory bool, isReactionFromCache bool) {
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
		OperationID: operationID,
		MsgData: &open_im_sdk.MsgData{
			SendID:      sendID,
			ClientMsgID: utils.GetMsgID(sendID),
			SessionType: sessionType,
			MsgFrom:     constant.SysMsgType,
			ContentType: contentType,
			Content:     []byte(content),
			//	ForceList:        params.ForceList,
			CreateTime: utils.GetCurrentTimestampByMill(),
			Options:    options,
		},
	}
	switch sessionType {
	case constant.SingleChatType, constant.NotificationChatType:
		pbData.MsgData.RecvID = sourceID
	case constant.GroupChatType, constant.SuperGroupChatType:
		pbData.MsgData.GroupID = sourceID
	}
	etcdConn, err := rpc.GetConn(context.Background(), config.Config.RpcRegisterName.OpenImMsgName)
	if err != nil {
		return
	}
	client := msg.NewMsgClient(etcdConn)
	reply, err := client.SendMsg(context.Background(), &pbData)
	if err != nil {
		log.NewError(operationID, "SendMsg rpc failed, ", pbData.String(), err.Error())
	} else if reply.ErrCode != 0 {
		log.NewError(operationID, "SendMsg rpc failed, ", pbData.String(), reply.ErrCode, reply.ErrMsg)
	}
}
