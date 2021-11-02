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
	"Open_IM/src/grpc-etcdv3/getcdv3"
	pbChat "Open_IM/src/proto/chat"
	pbGroup "Open_IM/src/proto/group"
	pbRelay "Open_IM/src/proto/relay"
	push "Open_IM/src/push/jpush"
	rpcChat "Open_IM/src/rpc/chat/chat"
	"Open_IM/src/utils"
	"context"
	"encoding/json"
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
	MOptions := utils.JsonStringToMap(Options) //Control whether to push message to sender's other terminal
	//isSenderSync := utils.GetSwitchFromOptions(MOptions, "senderSync")
	isOfflinePush := utils.GetSwitchFromOptions(MOptions, "offlinePush")
	log.InfoByKv("Get chat from msg_transfer And push chat", sendPbData.OperationID, "PushData", sendPbData)
	grpcCons := getcdv3.GetConn4Unique(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOnlineMessageRelayName)
	//Online push message
	log.InfoByKv("test", sendPbData.OperationID, "len  grpc", len(grpcCons), "data", sendPbData)
	for _, v := range grpcCons {
		msgClient := pbRelay.NewOnlineMessageRelayServiceClient(v)
		reply, err := msgClient.MsgToUser(context.Background(), sendPbData)
		if err != nil {
			log.InfoByKv("push data to client rpc err", sendPbData.OperationID, "err", err)
		}
		if reply != nil && reply.Resp != nil && err == nil {
			wsResult = append(wsResult, reply.Resp...)
		}
	}
	log.InfoByKv("push_result", sendPbData.OperationID, "result", wsResult, "sendData", sendPbData)
	if isOfflinePush {
		for _, v := range wsResult {
			if v.ResultCode == 0 {
				continue
			}
			//supported terminal
			for _, t := range pushTerminal {
				if v.RecvPlatFormID == t {
					//Use offline push messaging
					var UIDList []string
					UIDList = append(UIDList, sendPbData.RecvID)
					customContent := EChatContent{
						SessionType: int(sendPbData.SessionType),
						From:        sendPbData.SendID,
						To:          sendPbData.RecvID,
						Seq:         sendPbData.RecvSeq,
					}
					bCustomContent, _ := json.Marshal(customContent)
					jsonCustomContent := string(bCustomContent)
					push.JGAccountListPush(UIDList, jsonCustomContent, utils.PlatformIDToName(t))

				}
			}
		}

	}
}

func SendMsgByWS(m *pbChat.WSToMsgSvrChatMsg) {
	m.MsgID = rpcChat.GetMsgID(m.SendID)
	m.ClientMsgID = m.MsgID
	switch m.SessionType {
	case constant.SingleChatType:
		sendMsgToKafka(m, m.SendID, "msgKey--sendID")
		sendMsgToKafka(m, m.RecvID, "msgKey--recvID")
	case constant.GroupChatType:
		etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
		client := pbGroup.NewGroupClient(etcdConn)
		req := &pbGroup.GetGroupAllMemberReq{
			GroupID:     m.RecvID,
			Token:       config.Config.Secret,
			OperationID: m.OperationID,
		}
		reply, err := client.GetGroupAllMember(context.Background(), req)
		if err != nil {
			log.Error(m.Token, m.OperationID, "rpc  getGroupInfo failed, err = %s", err.Error())
			return
		}
		if reply.ErrorCode != 0 {
			log.Error(m.Token, m.OperationID, "rpc  getGroupInfo failed, err = %s", reply.ErrorMsg)
			return
		}
		groupID := m.RecvID
		for i, v := range reply.MemberList {
			m.RecvID = v.UserId + " " + groupID
			sendMsgToKafka(m, utils.IntToString(i), "msgKey--recvID+\" \"+groupID")
		}
	default:

	}

}
func sendMsgToKafka(m *pbChat.WSToMsgSvrChatMsg, key string, flag string) {
	pid, offset, err := producer.SendMessage(m, key)
	if err != nil {
		log.ErrorByKv("kafka send failed", m.OperationID, "send data", m.String(), "pid", pid, "offset", offset, "err", err.Error(), flag, key)
	}

}
