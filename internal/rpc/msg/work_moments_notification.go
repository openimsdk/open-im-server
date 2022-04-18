package msg

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbChat "Open_IM/pkg/proto/chat"
	pbCommon "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"encoding/json"
	"strings"
)

func CommentOneWorkMomentNotification(operationID, recvID string, comment db.CommentMsg, user db.User) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", recvID, user, comment)
	var req pbChat.SendMsgReq
	var msgData pbCommon.MsgData
	msgData.SendID = user.UserID
	msgData.RecvID = recvID
	msgData.ContentType = constant.WorkMomentNewCommentNotification
	msgData.SessionType = constant.SingleChatType
	msgData.MsgFrom = constant.UserMsgType
	bytes, err := json.Marshal(comment)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "marshal failed", err.Error())
	}
	msgData.Content = bytes
	msgData.SenderFaceURL = user.FaceURL
	msgData.SenderNickname = user.Nickname
	msgData.CreateTime = utils.GetCurrentTimestampByMill()
	msgData.ClientMsgID = utils.GetMsgID(user.UserID)
	req.MsgData = &msgData
	req.OperationID = operationID
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfflineMessageName)
	client := pbChat.NewChatClient(etcdConn)
	respPb, err := client.SendMsg(context.Background(), &req)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "send msg failed", err.Error())
		return
	}
	if respPb.ErrCode != 0 {
		log.NewError(operationID, utils.GetSelfFuncName(), "send tag msg failed ", respPb)
	}
}
