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
	pbCache "Open_IM/pkg/proto/cache"
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
	var wsResult []*pbRelay.SingelMsgToUserResultList
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
		reply, err := msgClient.OnlineBatchPushOneMsg(context.Background(), &pbRelay.OnlineBatchPushOneMsgReq{OperationID: pushMsg.OperationID, MsgData: pushMsg.MsgData, PushToUserIDList: []string{pushMsg.PushToUserID}})
		if err != nil {
			log.NewError("push data to client rpc err", pushMsg.OperationID, "err", err)
			continue
		}
		if reply != nil && reply.SinglePushResult != nil {
			wsResult = append(wsResult, reply.SinglePushResult...)
		}
	}
	log.NewInfo(pushMsg.OperationID, "push_result", wsResult, "sendData", pushMsg.MsgData)
	successCount++
	if isOfflinePush && pushMsg.PushToUserID != pushMsg.MsgData.SendID {
		for _, v := range wsResult {
			if v.OnlinePush {
				return
			}
		}
		//Use offline push messaging
		var UIDList []string
		UIDList = append(UIDList, pushMsg.PushToUserID)
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
				if utils.IsContain(pushMsg.PushToUserID, a.AtUserList) {
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
			return
		}
		if offlinePusher == nil {
			return
		}
		opts, err := GetOfflinePushOpts(pushMsg)
		if err != nil {
			log.NewError(pushMsg.OperationID, utils.GetSelfFuncName(), "GetOfflinePushOpts failed", pushMsg, err.Error())
		}
		log.NewInfo(pushMsg.OperationID, utils.GetSelfFuncName(), UIDList, content, jsonCustomContent, "opts:", opts)
		pushResult, err := offlinePusher.Push(UIDList, content, jsonCustomContent, pushMsg.OperationID, opts)
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
	log.Debug(pushMsg.OperationID, "Get super group msg from msg_transfer And push msg", pushMsg.String())
	getGroupMemberIDListFromCacheReq := &pbCache.GetGroupMemberIDListFromCacheReq{OperationID: pushMsg.OperationID, GroupID: pushMsg.MsgData.GroupID}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCacheName)
	client := pbCache.NewCacheClient(etcdConn)
	cacheResp, err := client.GetGroupMemberIDListFromCache(context.Background(), getGroupMemberIDListFromCacheReq)
	if err != nil {
		log.NewError(pushMsg.OperationID, "GetGroupMemberIDListFromCache rpc call failed ", err.Error())
		return
	}
	if cacheResp.CommonResp.ErrCode != 0 {
		log.NewError(pushMsg.OperationID, "GetGroupMemberIDListFromCache rpc logic call failed ", cacheResp.String())
		return
	}
	if len(grpcCons) == 0 {
		log.NewWarn(pushMsg.OperationID, "first GetConn4Unique ")
		grpcCons = getcdv3.GetConn4Unique(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOnlineMessageRelayName)
	}
	//Online push message
	log.Debug("test", pushMsg.OperationID, "len  grpc", len(grpcCons), "data", pushMsg.String())
	for _, v := range grpcCons {
		msgClient := pbRelay.NewOnlineMessageRelayServiceClient(v)
		reply, err := msgClient.SuperGroupOnlineBatchPushOneMsg(context.Background(), &pbRelay.OnlineBatchPushOneMsgReq{OperationID: pushMsg.OperationID, MsgData: pushMsg.MsgData, PushToUserIDList: cacheResp.UserIDList})
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
		onlineFailedUserIDList := utils.DifferenceString(onlineSuccessUserIDList, cacheResp.UserIDList)
		//Use offline push messaging
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
				if utils.IsContain(pushMsg.PushToUserID, a.AtUserList) {
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

		callbackResp := callbackOfflinePush(pushMsg.OperationID, onlineFailedUserIDList[0], pushMsg.MsgData)
		log.NewDebug(pushMsg.OperationID, utils.GetSelfFuncName(), "offline callback Resp")
		if callbackResp.ErrCode != 0 {
			log.NewError(pushMsg.OperationID, utils.GetSelfFuncName(), "callbackOfflinePush result: ", callbackResp)
		}
		if callbackResp.ActionCode != constant.ActionAllow {
			log.NewDebug(pushMsg.OperationID, utils.GetSelfFuncName(), "offlinePush stop")
			return
		}
		if offlinePusher == nil {
			return
		}
		opts, err := GetOfflinePushOpts(pushMsg)
		if err != nil {
			log.NewError(pushMsg.OperationID, utils.GetSelfFuncName(), "GetOfflinePushOpts failed", pushMsg, err.Error())
		}
		log.NewInfo(pushMsg.OperationID, utils.GetSelfFuncName(), onlineFailedUserIDList, content, jsonCustomContent, "opts:", opts)
		pushResult, err := offlinePusher.Push(onlineFailedUserIDList, content, jsonCustomContent, pushMsg.OperationID, opts)
		if err != nil {
			log.NewError(pushMsg.OperationID, "offline push error", pushMsg.String(), err.Error())
		} else {
			log.NewDebug(pushMsg.OperationID, "offline push return result is ", pushResult, pushMsg.MsgData)
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
	return opts, nil
}
