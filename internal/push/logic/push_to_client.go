/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@open-im.io").
** time(2021/3/5 14:31).
 */
package logic

import (
	"Open_IM/internal/push"
	utils2 "Open_IM/internal/utils"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbPush "Open_IM/pkg/proto/push"
	pbRelay "Open_IM/pkg/proto/relay"
	pbRtc "Open_IM/pkg/proto/rtc"
	"Open_IM/pkg/utils"
	"context"
	"strings"

	"github.com/golang/protobuf/proto"
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

//var grpcCons []*grpc.ClientConn

func MsgToUser(pushMsg *pbPush.PushMsgReq) {
	var wsResult []*pbRelay.SingelMsgToUserResultList
	isOfflinePush := utils.GetSwitchFromOptions(pushMsg.MsgData.Options, constant.IsOfflinePush)
	log.Debug(pushMsg.OperationID, "Get msg from msg_transfer And push msg", pushMsg.String())
	grpcCons := getcdv3.GetDefaultGatewayConn4Unique(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), pushMsg.OperationID)

	var UIDList = []string{pushMsg.PushToUserID}
	callbackResp := callbackOnlinePush(pushMsg.OperationID, UIDList, pushMsg.MsgData)
	log.NewDebug(pushMsg.OperationID, utils.GetSelfFuncName(), "OnlinePush callback Resp")
	if callbackResp.ErrCode != 0 {
		log.NewError(pushMsg.OperationID, utils.GetSelfFuncName(), "callbackOnlinePush result: ", callbackResp)
	}
	if callbackResp.ActionCode != constant.ActionAllow {
		log.NewDebug(pushMsg.OperationID, utils.GetSelfFuncName(), "OnlinePush stop")
		return
	}

	//Online push message
	log.Debug(pushMsg.OperationID, "len  grpc", len(grpcCons), "data", pushMsg.String())
	for _, v := range grpcCons {
		msgClient := pbRelay.NewRelayClient(v)
		reply, err := msgClient.SuperGroupOnlineBatchPushOneMsg(context.Background(), &pbRelay.OnlineBatchPushOneMsgReq{OperationID: pushMsg.OperationID, MsgData: pushMsg.MsgData, PushToUserIDList: []string{pushMsg.PushToUserID}})
		if err != nil {
			log.NewError("SuperGroupOnlineBatchPushOneMsg push data to client rpc err", pushMsg.OperationID, "err", err)
			continue
		}
		if reply != nil && reply.SinglePushResult != nil {
			wsResult = append(wsResult, reply.SinglePushResult...)
		}
	}
	log.NewInfo(pushMsg.OperationID, "push_result", wsResult, "sendData", pushMsg.MsgData)
	successCount++
	if isOfflinePush && pushMsg.PushToUserID != pushMsg.MsgData.SendID {
		// save invitation info for offline push
		for _, v := range wsResult {
			if v.OnlinePush {
				return
			}
		}
		if pushMsg.MsgData.ContentType == constant.SignalingNotification {
			if err := db.DB.HandleSignalInfo(pushMsg.OperationID, pushMsg.MsgData); err != nil {
				log.NewError(pushMsg.OperationID, utils.GetSelfFuncName(), err.Error(), pushMsg.MsgData)
				return
			}
		}
		var title, detailContent string
		callbackResp := callbackOfflinePush(pushMsg.OperationID, UIDList, pushMsg.MsgData, &[]string{})
		log.NewDebug(pushMsg.OperationID, utils.GetSelfFuncName(), "offline callback Resp")
		if callbackResp.ErrCode != 0 {
			log.NewError(pushMsg.OperationID, utils.GetSelfFuncName(), "callbackOfflinePush result: ", callbackResp)
		}
		if callbackResp.ActionCode != constant.ActionAllow {
			log.NewDebug(pushMsg.OperationID, utils.GetSelfFuncName(), "offlinePush stop")
			return
		}
		if pushMsg.MsgData.OfflinePushInfo != nil {
			title = pushMsg.MsgData.OfflinePushInfo.Title
			detailContent = pushMsg.MsgData.OfflinePushInfo.Desc
		}

		if offlinePusher == nil {
			return
		}
		opts, err := GetOfflinePushOpts(pushMsg)
		if err != nil {
			log.NewError(pushMsg.OperationID, utils.GetSelfFuncName(), "GetOfflinePushOpts failed", pushMsg, err.Error())
		}
		log.NewInfo(pushMsg.OperationID, utils.GetSelfFuncName(), UIDList, title, detailContent, "opts:", opts)
		if title == "" {
			switch pushMsg.MsgData.ContentType {
			case constant.Text:
				fallthrough
			case constant.Picture:
				fallthrough
			case constant.Voice:
				fallthrough
			case constant.Video:
				fallthrough
			case constant.File:
				title = constant.ContentType2PushContent[int64(pushMsg.MsgData.ContentType)]
			case constant.AtText:
				a := AtContent{}
				_ = utils.JsonStringToStruct(string(pushMsg.MsgData.Content), &a)
				if utils.IsContain(pushMsg.PushToUserID, a.AtUserList) {
					title = constant.ContentType2PushContent[constant.AtText] + constant.ContentType2PushContent[constant.Common]
				} else {
					title = constant.ContentType2PushContent[constant.GroupMsg]
				}
			case constant.SignalingNotification:
				title = constant.ContentType2PushContent[constant.SignalMsg]
			default:
				title = constant.ContentType2PushContent[constant.Common]

			}
			detailContent = title
		}
		pushResult, err := offlinePusher.Push(UIDList, title, detailContent, pushMsg.OperationID, opts)
		if err != nil {
			log.NewError(pushMsg.OperationID, "offline push error", pushMsg.String(), err.Error())
		} else {
			log.NewDebug(pushMsg.OperationID, "offline push return result is ", pushResult, pushMsg.MsgData)
		}
	}
}

func MsgToSuperGroupUser(pushMsg *pbPush.PushMsgReq) {
	var wsResult []*pbRelay.SingelMsgToUserResultList
	isOfflinePush := utils.GetSwitchFromOptions(pushMsg.MsgData.Options, constant.IsOfflinePush)
	log.Debug(pushMsg.OperationID, "Get super group msg from msg_transfer And push msg", pushMsg.String(), config.Config.Callback.CallbackBeforeSuperGroupOnlinePush.Enable)
	var pushToUserIDList []string
	if config.Config.Callback.CallbackBeforeSuperGroupOnlinePush.Enable {
		callbackResp := callbackBeforeSuperGroupOnlinePush(pushMsg.OperationID, pushMsg.PushToUserID, pushMsg.MsgData, &pushToUserIDList)
		log.NewDebug(pushMsg.OperationID, utils.GetSelfFuncName(), "offline callback Resp")
		if callbackResp.ErrCode != 0 {
			log.NewError(pushMsg.OperationID, utils.GetSelfFuncName(), "callbackOfflinePush result: ", callbackResp)
		}
		if callbackResp.ActionCode != constant.ActionAllow {
			log.NewDebug(pushMsg.OperationID, utils.GetSelfFuncName(), "onlinePush stop")
			return
		}
		log.NewDebug(pushMsg.OperationID, utils.GetSelfFuncName(), "callback userIDList Resp", pushToUserIDList)
	}
	if len(pushToUserIDList) == 0 {
		userIDList, err := utils2.GetGroupMemberUserIDList(pushMsg.MsgData.GroupID, pushMsg.OperationID)
		if err != nil {
			log.Error(pushMsg.OperationID, "GetGroupMemberUserIDList failed ", err.Error(), pushMsg.MsgData.GroupID)
			return
		}
		pushToUserIDList = userIDList
		//getGroupMemberIDListFromCacheReq := &pbCache.GetGroupMemberIDListFromCacheReq{OperationID: pushMsg.OperationID, GroupID: pushMsg.MsgData.GroupID}
		//etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCacheName, pushMsg.OperationID)
		//if etcdConn == nil {
		//	errMsg := pushMsg.OperationID + "getcdv3.GetDefaultConn == nil"
		//	log.NewError(pushMsg.OperationID, errMsg)
		//	return
		//}
		//client := pbCache.NewCacheClient(etcdConn)
		//cacheResp, err := client.GetGroupMemberIDListFromCache(context.Background(), getGroupMemberIDListFromCacheReq)
		//if err != nil {
		//	log.NewError(pushMsg.OperationID, "GetGroupMemberIDListFromCache rpc call failed ", err.Error())
		//	return
		//}
		//if cacheResp.CommonResp.ErrCode != 0 {
		//	log.NewError(pushMsg.OperationID, "GetGroupMemberIDListFromCache rpc logic call failed ", cacheResp.String())
		//	return
		//}
		//pushToUserIDList = cacheResp.UserIDList
	}

	grpcCons := getcdv3.GetDefaultGatewayConn4Unique(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), pushMsg.OperationID)

	//Online push message
	log.Debug(pushMsg.OperationID, "len  grpc", len(grpcCons), "data", pushMsg.String())
	for _, v := range grpcCons {
		msgClient := pbRelay.NewRelayClient(v)
		reply, err := msgClient.SuperGroupOnlineBatchPushOneMsg(context.Background(), &pbRelay.OnlineBatchPushOneMsgReq{OperationID: pushMsg.OperationID, MsgData: pushMsg.MsgData, PushToUserIDList: pushToUserIDList})
		if err != nil {
			log.NewError("push data to client rpc err", pushMsg.OperationID, "err", err)
			continue
		}
		if reply != nil && reply.SinglePushResult != nil {
			wsResult = append(wsResult, reply.SinglePushResult...)
		}
	}
	log.Debug(pushMsg.OperationID, "push_result", wsResult, "sendData", pushMsg.MsgData)
	successCount++
	if isOfflinePush {
		var onlineSuccessUserIDList []string
		onlineSuccessUserIDList = append(onlineSuccessUserIDList, pushMsg.MsgData.SendID)
		for _, v := range wsResult {
			if v.OnlinePush && v.UserID != pushMsg.MsgData.SendID {
				onlineSuccessUserIDList = append(onlineSuccessUserIDList, v.UserID)
			}
		}
		onlineFailedUserIDList := utils.DifferenceString(onlineSuccessUserIDList, pushToUserIDList)
		//Use offline push messaging
		var title, detailContent string
		if len(onlineFailedUserIDList) > 0 {
			var offlinePushUserIDList []string
			var needOfflinePushUserIDList []string
			callbackResp := callbackOfflinePush(pushMsg.OperationID, onlineFailedUserIDList, pushMsg.MsgData, &offlinePushUserIDList)
			log.NewDebug(pushMsg.OperationID, utils.GetSelfFuncName(), "offline callback Resp")
			if callbackResp.ErrCode != 0 {
				log.NewError(pushMsg.OperationID, utils.GetSelfFuncName(), "callbackOfflinePush result: ", callbackResp)
			}
			if callbackResp.ActionCode != constant.ActionAllow {
				log.NewDebug(pushMsg.OperationID, utils.GetSelfFuncName(), "offlinePush stop")
				return
			}
			if pushMsg.MsgData.OfflinePushInfo != nil {
				title = pushMsg.MsgData.OfflinePushInfo.Title
				detailContent = pushMsg.MsgData.OfflinePushInfo.Desc
			}
			if len(offlinePushUserIDList) > 0 {
				needOfflinePushUserIDList = offlinePushUserIDList
			} else {
				needOfflinePushUserIDList = onlineFailedUserIDList
			}

			if offlinePusher == nil {
				return
			}
			opts, err := GetOfflinePushOpts(pushMsg)
			if err != nil {
				log.NewError(pushMsg.OperationID, utils.GetSelfFuncName(), "GetOfflinePushOpts failed", pushMsg, err.Error())
			}
			log.NewInfo(pushMsg.OperationID, utils.GetSelfFuncName(), onlineFailedUserIDList, title, detailContent, "opts:", opts)
			if title == "" {
				switch pushMsg.MsgData.ContentType {
				case constant.Text:
					fallthrough
				case constant.Picture:
					fallthrough
				case constant.Voice:
					fallthrough
				case constant.Video:
					fallthrough
				case constant.File:
					title = constant.ContentType2PushContent[int64(pushMsg.MsgData.ContentType)]
				case constant.AtText:
					a := AtContent{}
					_ = utils.JsonStringToStruct(string(pushMsg.MsgData.Content), &a)
					if utils.IsContain(pushMsg.PushToUserID, a.AtUserList) {
						title = constant.ContentType2PushContent[constant.AtText] + constant.ContentType2PushContent[constant.Common]
					} else {
						title = constant.ContentType2PushContent[constant.GroupMsg]
					}
				case constant.SignalingNotification:
					title = constant.ContentType2PushContent[constant.SignalMsg]
				default:
					title = constant.ContentType2PushContent[constant.Common]

				}
				detailContent = title
			}
			pushResult, err := offlinePusher.Push(needOfflinePushUserIDList, title, detailContent, pushMsg.OperationID, opts)
			if err != nil {
				log.NewError(pushMsg.OperationID, "offline push error", pushMsg.String(), err.Error())
			} else {
				log.NewDebug(pushMsg.OperationID, "offline push return result is ", pushResult, pushMsg.MsgData)
			}
		}

	}
}

func GetOfflinePushOpts(pushMsg *pbPush.PushMsgReq) (opts push.PushOpts, err error) {
	if pushMsg.MsgData.ContentType < constant.SignalingNotificationEnd && pushMsg.MsgData.ContentType > constant.SignalingNotificationBegin {
		req := &pbRtc.SignalReq{}
		if err := proto.Unmarshal(pushMsg.MsgData.Content, req); err != nil {
			return opts, utils.Wrap(err, "")
		}
		log.NewDebug(pushMsg.OperationID, utils.GetSelfFuncName(), "SignalReq: ", req.String())
		switch req.Payload.(type) {
		case *pbRtc.SignalReq_Invite, *pbRtc.SignalReq_InviteInGroup:
			opts.Signal.ClientMsgID = pushMsg.MsgData.ClientMsgID
			log.NewDebug(pushMsg.OperationID, opts)
		}
	}
	if pushMsg.MsgData.OfflinePushInfo != nil {
		opts.IOSBadgeCount = pushMsg.MsgData.OfflinePushInfo.IOSBadgeCount
		opts.IOSPushSound = pushMsg.MsgData.OfflinePushInfo.IOSPushSound
		opts.Data = pushMsg.MsgData.OfflinePushInfo.Ex
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
//		etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
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
