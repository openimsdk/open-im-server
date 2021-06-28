/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@open-im.io").
** time(2021/3/5 14:31).
 */
package logic

import (
	"Open_IM/src/common/config"
	"Open_IM/src/common/constant"
	"Open_IM/src/common/log"
	pbChat "Open_IM/src/proto/chat"
	pbRelay "Open_IM/src/proto/relay"
	pbGetInfo "Open_IM/src/proto/user"
	rpcChat "Open_IM/src/rpc/chat/chat"
	"Open_IM/src/rpc/user/internal_service"
	"Open_IM/src/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/skiffer-git/grpc-etcdv3/getcdv3"
	"strings"
)

type EChatContent struct {
	SessionType int    `json:"chatType"`
	From        string `json:"from"`
	To          string `json:"to"`
	Seq         int64  `json:"seq"`
}

func MsgToUser(sendPbData *pbRelay.MsgToUserReq, OfflineInfo, Options string) {
	var wsResult []*pbRelay.SingleMsgToUser
	isShouldOfflinePush := true
	MOptions := utils.JsonStringToMap(Options)
	isOfflinePush := utils.GetSwitchFromOptions(MOptions, "offlinePush")
	log.InfoByKv("Get chat from msg_transfer And push chat", sendPbData.OperationID, "PushData", sendPbData)
	grpcCons := getcdv3.GetConn4Unique(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOnlineMessageRelayName)
	//Online push message
	for _, v := range grpcCons {
		msgClient := pbRelay.NewOnlineMessageRelayServiceClient(v)
		reply, err := msgClient.MsgToUser(context.Background(), sendPbData)
		if reply != nil && reply.Resp != nil && err == nil {
			wsResult = append(wsResult, reply.Resp...)
		}
	}
	if isOfflinePush && sendPbData.ContentType != constant.SyncSenderMsg {

		for _, t := range pushTerminal {
			for _, v := range wsResult {
				if v.RecvPlatFormID == t && v.ResultCode == 0 {
					isShouldOfflinePush = false
					break
				}
			}
			if isShouldOfflinePush {
				//Use offline push messaging
				var UIDList []string
				UIDList = append(UIDList, sendPbData.RecvID)
				var sendUIDList []string
				sendUIDList = append(sendUIDList, sendPbData.SendID)
				userInfo, err := internal_service.GetUserInfoClient(&pbGetInfo.GetUserInfoReq{UserIDList: sendUIDList, OperationID: sendPbData.OperationID})
				if err != nil {
					log.ErrorByArgs(fmt.Sprintf("err=%v,call GetUserInfoClient rpc server failed", err))
					return
				}

				customContent := EChatContent{
					SessionType: int(sendPbData.SessionType),
					From:        sendPbData.SendID,
					To:          sendPbData.RecvID,
					Seq:         sendPbData.RecvSeq,
				}
				bCustomContent, _ := json.Marshal(customContent)

				jsonCustomContent := string(bCustomContent)
				switch sendPbData.ContentType {
				case constant.Text:
					IOSAccountListPush(UIDList, userInfo.Data[0].Name, sendPbData.Content, jsonCustomContent)
				case constant.Picture:
					IOSAccountListPush(UIDList, userInfo.Data[0].Name, constant.ContentType2PushContent[constant.Picture], jsonCustomContent)
				case constant.Voice:
					IOSAccountListPush(UIDList, userInfo.Data[0].Name, constant.ContentType2PushContent[constant.Voice], jsonCustomContent)
				case constant.Video:
					IOSAccountListPush(UIDList, userInfo.Data[0].Name, constant.ContentType2PushContent[constant.Video], jsonCustomContent)
				case constant.File:
					IOSAccountListPush(UIDList, userInfo.Data[0].Name, constant.ContentType2PushContent[constant.File], jsonCustomContent)
				default:

				}

			} else {
				isShouldOfflinePush = true
			}
		}

	}

}

func SendMsgByWS(m *pbChat.WSToMsgSvrChatMsg) {
	m.MsgID = rpcChat.GetMsgID(m.SendID)
	pid, offset, err := producer.SendMessage(m, m.SendID)
	if err != nil {
		log.ErrorByKv("sys send msg to kafka  failed", m.OperationID, "send data", m.String(), "pid", pid, "offset", offset, "err", err.Error(), "msgKey--sendID", m.SendID)
	}
	pid, offset, err = producer.SendMessage(m, m.RecvID)
	if err != nil {
		log.ErrorByKv("kafka send failed", m.OperationID, "send data", m.String(), "pid", pid, "offset", offset, "err", err.Error(), "msgKey--recvID", m.RecvID)
	}
}
