/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@open-im.io").
** time(2021/3/5 14:31).
 */
package logic

import (
	"Open_IM/internal/push"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbPush "Open_IM/pkg/proto/push"
	pbRelay "Open_IM/pkg/proto/relay"
	pbRtc "Open_IM/pkg/proto/rtc"
	"Open_IM/pkg/utils"
	"context"
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"strings"
)

type OpenIMContent struct {
	SessionType int    `json:"sessionType"`
	From        string `json:"from"`
	To          string `json:"to"`
	Seq         uint32 `json:"seq"`
}
type AtContent struct {
	Text       string   `json:"text"`
	AtUserList []string `json:"atUserList"`
	IsAtSelf   bool     `json:"isAtSelf"`
}

var grpcCons []*grpc.ClientConn

func MsgToUser(pushMsg *pbPush.PushMsgReq) {
	var wsResult []*pbRelay.SingleMsgToUser
	isOfflinePush := utils.GetSwitchFromOptions(pushMsg.MsgData.Options, constant.IsOfflinePush)
	log.Debug(pushMsg.OperationID, "Get msg from msg_transfer And push msg", pushMsg.String())
	if len(grpcCons) == 0 {
		log.NewWarn(pushMsg.OperationID, "first GetConn4Unique ")
		grpcCons = getcdv3.GetConn4Unique(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOnlineMessageRelayName)
	}
	//Online push message
	log.Debug("test", pushMsg.OperationID, "len  grpc", len(grpcCons), "data", pushMsg.String())
	for _, v := range grpcCons {
		msgClient := pbRelay.NewOnlineMessageRelayServiceClient(v)
		reply, err := msgClient.OnlinePushMsg(context.Background(), &pbRelay.OnlinePushMsgReq{OperationID: pushMsg.OperationID, MsgData: pushMsg.MsgData, PushToUserID: pushMsg.PushToUserID})
		if err != nil {
			log.NewError("push data to client rpc err", pushMsg.OperationID, "err", err)
			continue
		}
		if reply != nil && reply.Resp != nil {
			wsResult = append(wsResult, reply.Resp...)
		}
	}
	log.NewInfo(pushMsg.OperationID, "push_result", wsResult, "sendData", pushMsg.MsgData)
	successCount++
	if isOfflinePush && pushMsg.PushToUserID != pushMsg.MsgData.SendID {
		for _, v := range wsResult {
			if v.ResultCode == 0 {
				if utils.IsContainInt32(v.RecvPlatFormID, pushTerminal) {
					break
				}
				continue
			}
			if utils.IsContainInt32(v.RecvPlatFormID, pushTerminal) {
				//Use offline push messaging
				var UIDList []string
				UIDList = append(UIDList, v.RecvID)
				customContent := OpenIMContent{
					SessionType: int(pushMsg.MsgData.SessionType),
					From:        pushMsg.MsgData.SendID,
					To:          pushMsg.MsgData.RecvID,
					Seq:         pushMsg.MsgData.Seq,
				}
				bCustomContent, _ := json.Marshal(customContent)
				jsonCustomContent := string(bCustomContent)
				var content string
				if pushMsg.MsgData.OfflinePushInfo != nil {
					content = pushMsg.MsgData.OfflinePushInfo.Title

				} else {
					switch pushMsg.MsgData.ContentType {
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
						_ = utils.JsonStringToStruct(string(pushMsg.MsgData.Content), &a)
						if utils.IsContain(v.RecvID, a.AtUserList) {
							content = constant.ContentType2PushContent[constant.AtText] + constant.ContentType2PushContent[constant.Common]
						} else {
							content = constant.ContentType2PushContent[constant.GroupMsg]
						}
					case constant.SignalingNotification:
						content = constant.ContentType2PushContent[constant.SignalMsg]
					default:
						content = constant.ContentType2PushContent[constant.Common]
					}
				}
				callbackResp := callbackOfflinePush(pushMsg.OperationID, UIDList[0], pushMsg.MsgData)
				log.NewDebug(pushMsg.OperationID, utils.GetSelfFuncName(), "offline callback Resp")
				if callbackResp.ErrCode != 0 {
					log.NewError(pushMsg.OperationID, utils.GetSelfFuncName(), "callbackOfflinePush result: ", callbackResp)
				}
				if callbackResp.ActionCode != constant.ActionAllow {
					log.NewDebug(pushMsg.OperationID, utils.GetSelfFuncName(), "offlinePush stop")
					break
				}

				if offlinePusher == nil {
					break
				}
				opts, err := GetOfflinePushOpts(pushMsg)
				if err != nil {
					log.NewError(pushMsg.OperationID, utils.GetSelfFuncName(), "GetOfflinePushOpts failed", pushMsg, err.Error())
				}
				log.NewInfo(pushMsg.OperationID, utils.GetSelfFuncName(), "opts:", opts)
				pushResult, err := offlinePusher.Push(UIDList, content, jsonCustomContent, pushMsg.OperationID, opts)
				if err != nil {
					log.NewError(pushMsg.OperationID, "offline push error", pushMsg.String(), err.Error())
				} else {
					log.NewDebug(pushMsg.OperationID, "offline push return result is ", pushResult, pushMsg.MsgData)
				}
				break
			}

		}

	}
}

func GetOfflinePushOpts(pushMsg *pbPush.PushMsgReq) (opts push.PushOpts, err error) {
	if pushMsg.MsgData.ContentType < constant.SignalingNotificationEnd && pushMsg.MsgData.ContentType > constant.SignalingNotification {
		req := &pbRtc.SignalReq{}
		if err := proto.Unmarshal(pushMsg.MsgData.Content, req); err != nil {
			return opts, err
		}
		log.NewInfo("", utils.GetSelfFuncName(), "SignalReq: ", req.String())
		switch req.Payload.(type) {
		case *pbRtc.SignalReq_Invite, *pbRtc.SignalReq_InviteInGroup:
			opts.Signal.ClientMsgID = pushMsg.MsgData.ClientMsgID
		}

	}
	return opts, nil
}

//func SendMsgByWS(m *pbChat.WSToMsgSvrChatMsg) {
//	m.MsgID = rpcChat.GetMsgID(m.SendID)
//	m.ClientMsgID = m.MsgID
//	switch m.SessionType {
//	case constant.SingleChatType:
//		sendMsgToKafka(m, m.SendID, "msgKey--sendID")
//		sendMsgToKafka(m, m.RecvID, "msgKey--recvID")
//	case constant.GroupChatType:
//		etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
//		client := pbGroup.NewGroupClient(etcdConn)
//		req := &pbGroup.Req{
//			GroupID:     m.RecvID,
//			Token:       config.Config.Secret,
//			OperationID: m.OperationID,
//		}
//		reply, err := client.(context.Background(), req)
//		if err != nil {
//			log.Error(m.Token, m.OperationID, "rpc  getGroupInfo failed, err = %s", err.Error())
//			return
//		}
//		if reply.ErrorCode != 0 {
//			log.Error(m.Token, m.OperationID, "rpc  getGroupInfo failed, err = %s", reply.ErrorMsg)
//			return
//		}
//		groupID := m.RecvID
//		for i, v := range reply.MemberList {
//			m.RecvID = v.UserId + " " + groupID
//			sendMsgToKafka(m, utils.IntToString(i), "msgKey--recvID+\" \"+groupID")
//		}
//	default:
//
//	}
//}
//
//func sendMsgToKafka(m *pbChat.WSToMsgSvrChatMsg, key string, flag string) {
//	pid, offset, err := producer.SendMessage(m, key)
//	if err != nil {
//		log.ErrorByKv("kafka send failed", m.OperationID, "send data", m.String(), "pid", pid, "offset", offset, "err", err.Error(), flag, key)
//	}
//
//}
