package msg

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbChat "Open_IM/pkg/proto/msg"
	pbCommon "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"strings"
)

func TagSendMessage(operationID string, user *db.User, recvID, content string, senderPlatformID int32) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", user.UserID, recvID, content)
	var req pbChat.SendMsgReq
	var msgData pbCommon.MsgData
	msgData.SendID = user.UserID
	msgData.RecvID = recvID
	msgData.ContentType = constant.Custom
	msgData.SessionType = constant.SingleChatType
	msgData.MsgFrom = constant.UserMsgType
	msgData.Content = []byte(content)
	msgData.SenderFaceURL = user.FaceURL
	msgData.SenderNickname = user.Nickname
	msgData.Options = map[string]bool{}
	msgData.Options[constant.IsSenderConversationUpdate] = false
	msgData.Options[constant.IsSenderNotificationPush] = false
	msgData.CreateTime = utils.GetCurrentTimestampByMill()
	msgData.ClientMsgID = utils.GetMsgID(user.UserID)
	msgData.SenderPlatformID = senderPlatformID
	req.MsgData = &msgData
	req.OperationID = operationID
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMsgName, operationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		return
	}

	client := pbChat.NewMsgClient(etcdConn)
	respPb, err := client.SendMsg(context.Background(), &req)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "send msg failed", err.Error())
		return
	}
	if respPb.ErrCode != 0 {
		log.NewError(operationID, utils.GetSelfFuncName(), "send tag msg failed ", respPb)
	}
}
