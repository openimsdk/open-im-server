package msg

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	"Open_IM/pkg/proto/msg"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"strings"
)

func ExtendMessageUpdatedNotification(operationID, sendID string, sourceID string, sessionType int32,
	req *msg.ModifyMessageReactionExtensionsReq, resp *msg.ModifyMessageReactionExtensionsResp, isHistory bool) {
	m := make(map[string]interface{})
	m["rep"] = req
	m["resp"] = resp
	messageReactionSender(operationID, sendID, sourceID, sessionType, constant.ReactionMessageModifier, utils.StructToJsonString(m), isHistory)
}
func messageReactionSender(operationID, sendID string, sourceID string, sessionType, contentType int32, content string, isHistory bool) {
	options := make(map[string]bool, 5)
	utils.SetSwitchFromOptions(options, constant.IsOfflinePush, false)
	utils.SetSwitchFromOptions(options, constant.IsConversationUpdate, false)
	utils.SetSwitchFromOptions(options, constant.IsSenderConversationUpdate, false)
	utils.SetSwitchFromOptions(options, constant.IsUnreadCount, false)
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
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMsgName, operationID)
	if etcdConn == nil {
		errMsg := operationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(operationID, errMsg)
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
