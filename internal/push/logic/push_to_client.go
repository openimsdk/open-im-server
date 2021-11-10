/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@open-im.io").
** time(2021/3/5 14:31).
 */
package logic

import (
	push "Open_IM/internal/push/jpush"
	rpcChat "Open_IM/internal/rpc/chat"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbChat "Open_IM/pkg/proto/chat"
	pbGroup "Open_IM/pkg/proto/group"
	pbRelay "Open_IM/pkg/proto/relay"
	"Open_IM/pkg/utils"
	"context"
	"encoding/json"
	"strings"
)

type OpenIMContent struct {
	SessionType int    `json:"sessionType"`
	From        string `json:"from"`
	To          string `json:"to"`
	Seq         int64  `json:"seq"`
}
type AtContent struct {
	Text       string   `json:"text"`
	AtUserList []string `json:"atUserList"`
	IsAtSelf   bool     `json:"isAtSelf"`
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
	if sendPbData.ContentType != constant.Typing && sendPbData.ContentType != constant.HasReadReceipt {
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
						UIDList = append(UIDList, v.RecvID)
						customContent := OpenIMContent{
							SessionType: int(sendPbData.SessionType),
							From:        sendPbData.SendID,
							To:          sendPbData.RecvID,
							Seq:         sendPbData.RecvSeq,
						}
						bCustomContent, _ := json.Marshal(customContent)
						jsonCustomContent := string(bCustomContent)
						var content string
						switch sendPbData.ContentType {
						case constant.Text:
							content = constant.ContentType2PushContent[constant.Text]
						case constant.Picture:
							content = constant.ContentType2PushContent[constant.Picture]
						case constant.Voice:
							content = constant.ContentType2PushContent[constant.Voice]
						case constant.Video:
							content = constant.ContentType2PushContent[constant.Video]
						case constant.File:
							content = constant.ContentType2PushContent[constant.File]
						case constant.AtText:
							a := AtContent{}
							_ = utils.JsonStringToStruct(sendPbData.Content, &a)
							if utils.IsContain(v.RecvID, a.AtUserList) {
								content = constant.ContentType2PushContent[constant.AtText] + constant.ContentType2PushContent[constant.Common]
							} else {
								content = constant.ContentType2PushContent[constant.GroupMsg]
							}
						default:
						}
						pushResult, err := push.JGAccountListPush(UIDList, content, jsonCustomContent, utils.PlatformIDToName(t))
						if err != nil {
							log.NewError(sendPbData.OperationID, "offline push error", sendPbData.String(), err.Error(), t)
						} else {
							log.NewDebug(sendPbData.OperationID, "offline push return result is ", string(pushResult), sendPbData, t)
						}

					}
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
