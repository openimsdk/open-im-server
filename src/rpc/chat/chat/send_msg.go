package rpcChat

import (
	"Open_IM/src/common/log"
	pbChat "Open_IM/src/proto/chat"
	"Open_IM/src/utils"
	"context"
	"math/rand"
	"strconv"
	"time"
)

func (rpc *rpcChat) UserSendMsg(_ context.Context, pb *pbChat.UserSendMsgReq) (*pbChat.UserSendMsgResp, error) {

	serverMsgID := GetMsgID(pb.SendID)
	pbData := pbChat.WSToMsgSvrChatMsg{}
	pbData.MsgFrom = pb.MsgFrom
	pbData.SessionType = pb.SessionType
	pbData.ContentType = pb.ContentType
	pbData.Content = pb.Content
	pbData.RecvID = pb.RecvID
	pbData.ForceList = pb.ForceList
	pbData.OfflineInfo = pb.OffLineInfo
	pbData.Options = pb.Options
	pbData.PlatformID = pb.PlatformID
	pbData.SendID = pb.SendID
	pbData.MsgID = serverMsgID
	pbData.OperationID = pb.OperationID
	pbData.Token = pb.Token
	pbData.SendTime = utils.GetCurrentTimestampBySecond()
	rpc.sendMsgToKafka(&pbData, pbData.RecvID)
	rpc.sendMsgToKafka(&pbData, pbData.SendID)
	replay := pbChat.UserSendMsgResp{}
	replay.ReqIdentifier = pb.ReqIdentifier
	replay.MsgIncr = pb.MsgIncr
	replay.ClientMsgID = pb.ClientMsgID
	replay.ServerMsgID = serverMsgID

	return &replay, nil
}
func (rpc *rpcChat) sendMsgToKafka(m *pbChat.WSToMsgSvrChatMsg, key string) {
	pid, offset, err := rpc.producer.SendMessage(m, key)
	if err != nil {
		log.ErrorByKv("kafka send failed", m.OperationID, "send data", m.String(), "pid", pid, "offset", offset, "err", err.Error())
	}
}
func GetMsgID(sendID string) string {
	t := time.Now().Format("2006-01-02 15:04:05")
	return t + "-" + sendID + "-" + strconv.Itoa(rand.Int())
}
