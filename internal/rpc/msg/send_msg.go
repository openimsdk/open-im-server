package msg

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	http2 "Open_IM/pkg/common/http"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbChat "Open_IM/pkg/proto/chat"
	pbGroup "Open_IM/pkg/proto/group"
	sdk_ws "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type MsgCallBackReq struct {
	SendID       string `json:"sendID"`
	RecvID       string `json:"recvID"`
	Content      string `json:"content"`
	SendTime     int64  `json:"sendTime"`
	MsgFrom      int32  `json:"msgFrom"`
	ContentType  int32  `json:"contentType"`
	SessionType  int32  `json:"sessionType"`
	PlatformID   int32  `json:"senderPlatformID"`
	MsgID        string `json:"msgID"`
	IsOnlineOnly bool   `json:"isOnlineOnly"`
}
type MsgCallBackResp struct {
	ErrCode         int32  `json:"errCode"`
	ErrMsg          string `json:"errMsg"`
	ResponseErrCode int32  `json:"responseErrCode"`
	ResponseResult  struct {
		ModifiedMsg string `json:"modifiedMsg"`
		Ext         string `json:"ext"`
	}
}

func (rpc *rpcChat) encapsulateMsgData(msg *sdk_ws.MsgData) {
	msg.ServerMsgID = GetMsgID(msg.SendID)
	if msg.SendTime == 0 {
		msg.SendTime = utils.GetCurrentTimestampByMill()
	}
	switch msg.ContentType {
	case constant.Text:
		fallthrough
	case constant.Picture:
		fallthrough
	case constant.Voice:
		fallthrough
	case constant.Video:
		fallthrough
	case constant.File:
		fallthrough
	case constant.AtText:
		fallthrough
	case constant.Merger:
		fallthrough
	case constant.Card:
		fallthrough
	case constant.Location:
		fallthrough
	case constant.Custom:
		fallthrough
	case constant.Quote:
		utils.SetSwitchFromOptions(msg.Options, constant.IsConversationUpdate, true)
		utils.SetSwitchFromOptions(msg.Options, constant.IsUnreadCount, true)
		utils.SetSwitchFromOptions(msg.Options, constant.IsSenderSync, true)
	case constant.Revoke:
		utils.SetSwitchFromOptions(msg.Options, constant.IsUnreadCount, false)
		utils.SetSwitchFromOptions(msg.Options, constant.IsOfflinePush, false)
	case constant.HasReadReceipt:
		utils.SetSwitchFromOptions(msg.Options, constant.IsConversationUpdate, false)
		utils.SetSwitchFromOptions(msg.Options, constant.IsUnreadCount, false)
		utils.SetSwitchFromOptions(msg.Options, constant.IsOfflinePush, false)
	case constant.Typing:
		utils.SetSwitchFromOptions(msg.Options, constant.IsHistory, false)
		utils.SetSwitchFromOptions(msg.Options, constant.IsPersistent, false)
		utils.SetSwitchFromOptions(msg.Options, constant.IsSenderSync, false)
		utils.SetSwitchFromOptions(msg.Options, constant.IsConversationUpdate, false)
		utils.SetSwitchFromOptions(msg.Options, constant.IsUnreadCount, false)
		utils.SetSwitchFromOptions(msg.Options, constant.IsOfflinePush, false)

	}
}
func (rpc *rpcChat) SendMsg(_ context.Context, pb *pbChat.SendMsgReq) (*pbChat.SendMsgResp, error) {
	replay := pbChat.SendMsgResp{}
	log.NewDebug(pb.OperationID, "rpc sendMsg come here", pb.String())
	//if !utils.VerifyToken(pb.Token, pb.SendID) {
	//	return returnMsg(&replay, pb, http.StatusUnauthorized, "token validate err,not authorized", "", 0)
	rpc.encapsulateMsgData(pb.MsgData)
	msgToMQ := pbChat.MsgDataToMQ{Token: pb.Token, OperationID: pb.OperationID}
	//options := utils.JsonStringToMap(pbData.Options)
	isHistory := utils.GetSwitchFromOptions(pb.MsgData.Options, constant.IsHistory)
	mReq := MsgCallBackReq{
		SendID:      pb.MsgData.SendID,
		RecvID:      pb.MsgData.RecvID,
		Content:     string(pb.MsgData.Content),
		SendTime:    pb.MsgData.SendTime,
		MsgFrom:     pb.MsgData.MsgFrom,
		ContentType: pb.MsgData.ContentType,
		SessionType: pb.MsgData.SessionType,
		PlatformID:  pb.MsgData.SenderPlatformID,
		MsgID:       pb.MsgData.ClientMsgID,
	}
	if !isHistory {
		mReq.IsOnlineOnly = true
	}
	mResp := MsgCallBackResp{}
	if config.Config.MessageCallBack.CallbackSwitch {
		bMsg, err := http2.Post(config.Config.MessageCallBack.CallbackUrl, mReq, config.Config.MessageCallBack.CallBackTimeOut)
		if err != nil {
			log.ErrorByKv("callback to Business server err", pb.OperationID, "args", pb.String(), "err", err.Error())
			return returnMsg(&replay, pb, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "", 0)
		} else if err = json.Unmarshal(bMsg, &mResp); err != nil {
			log.ErrorByKv("ws json Unmarshal err", pb.OperationID, "args", pb.String(), "err", err.Error())
			return returnMsg(&replay, pb, 200, err.Error(), "", 0)
		} else {
			if mResp.ErrCode != 0 {
				return returnMsg(&replay, pb, mResp.ResponseErrCode, mResp.ErrMsg, "", 0)
			} else {
				pb.MsgData.Content = []byte(mResp.ResponseResult.ModifiedMsg)
			}
		}
	}
	switch pb.MsgData.SessionType {
	case constant.SingleChatType:
		isSend := modifyMessageByUserMessageReceiveOpt(pb.MsgData.RecvID, pb.MsgData.SendID, constant.SingleChatType, pb)
		if isSend {
			msgToMQ.MsgData = pb.MsgData
			err1 := rpc.sendMsgToKafka(&msgToMQ, msgToMQ.MsgData.RecvID)
			if err1 != nil {
				log.NewError(msgToMQ.OperationID, "kafka send msg err:RecvID", msgToMQ.MsgData.RecvID, msgToMQ.String())
				return returnMsg(&replay, pb, 201, "kafka send msg err", "", 0)
			}
		}
		if msgToMQ.MsgData.SendID != msgToMQ.MsgData.RecvID { //Filter messages sent to yourself
			err2 := rpc.sendMsgToKafka(&msgToMQ, msgToMQ.MsgData.SendID)
			if err2 != nil {
				log.NewError(msgToMQ.OperationID, "kafka send msg err:SendID", msgToMQ.MsgData.SendID, msgToMQ.String())
				return returnMsg(&replay, pb, 201, "kafka send msg err", "", 0)
			}
		}
		return returnMsg(&replay, pb, 0, "", msgToMQ.MsgData.ServerMsgID, msgToMQ.MsgData.SendTime)
	case constant.GroupChatType:
		etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
		client := pbGroup.NewGroupClient(etcdConn)
		req := &pbGroup.GetGroupAllMemberReq{
			GroupID:     pb.MsgData.GroupID,
			OperationID: pb.OperationID,
		}
		reply, err := client.GetGroupAllMember(context.Background(), req)
		if err != nil {
			log.Error(pb.Token, pb.OperationID, "rpc send_msg getGroupInfo failed, err = %s", err.Error())
			return returnMsg(&replay, pb, 201, err.Error(), "", 0)
		}
		if reply.ErrCode != 0 {
			log.Error(pb.Token, pb.OperationID, "rpc send_msg getGroupInfo failed, err = %s", reply.ErrMsg)
			return returnMsg(&replay, pb, reply.ErrCode, reply.ErrMsg, "", 0)
		}
		groupID := pb.MsgData.GroupID
		for _, v := range reply.MemberList {
			pb.MsgData.RecvID = v.UserID
			isSend := modifyMessageByUserMessageReceiveOpt(v.UserID, groupID, constant.GroupChatType, pb)
			if isSend {
				msgToMQ.MsgData = pb.MsgData
				err := rpc.sendMsgToKafka(&msgToMQ, v.UserID)
				if err != nil {
					log.NewError(msgToMQ.OperationID, "kafka send msg err:UserId", v.UserID, msgToMQ.String())
					return returnMsg(&replay, pb, 201, "kafka send msg err", "", 0)
				}
			}

		}
		return returnMsg(&replay, pb, 0, "", msgToMQ.MsgData.ServerMsgID, msgToMQ.MsgData.SendTime)
	default:
		return returnMsg(&replay, pb, 203, "unkonwn sessionType", "", 0)
	}
}

func (rpc *rpcChat) sendMsgToKafka(m *pbChat.MsgDataToMQ, key string) error {
	pid, offset, err := rpc.producer.SendMessage(m, key)
	if err != nil {
		log.ErrorByKv("kafka send failed", m.OperationID, "send data", m.String(), "pid", pid, "offset", offset, "err", err.Error(), "key", key)
	}
	return err
}
func GetMsgID(sendID string) string {
	t := time.Now().Format("2006-01-02 15:04:05")
	return t + "-" + sendID + "-" + strconv.Itoa(rand.Int())
}

func returnMsg(replay *pbChat.SendMsgResp, pb *pbChat.SendMsgReq, errCode int32, errMsg, serverMsgID string, sendTime int64) (*pbChat.SendMsgResp, error) {
	replay.ErrCode = errCode
	replay.ErrMsg = errMsg
	replay.ServerMsgID = serverMsgID
	replay.ClientMsgID = pb.MsgData.ClientMsgID
	replay.SendTime = sendTime
	return replay, nil
}

func modifyMessageByUserMessageReceiveOpt(userID, sourceID string, sessionType int, pb *pbChat.SendMsgReq) bool {
	conversationID := utils.GetConversationIDBySessionType(sourceID, sessionType)
	opt, err := db.DB.GetSingleConversationMsgOpt(userID, conversationID)
	if err != nil || err != redis.ErrNil {
		log.NewError(pb.OperationID, "GetSingleConversationMsgOpt from redis err", pb.String())
		return true
	}
	switch opt {
	case constant.ReceiveMessage:
		return true
	case constant.NotReceiveMessage:
		return false
	case constant.ReceiveNotNotifyMessage:
		if pb.MsgData.Options == nil {
			pb.MsgData.Options = make(map[string]bool, 10)
		}
		utils.SetSwitchFromOptions(pb.MsgData.Options, constant.IsOfflinePush, false)
		return true
	}

	return true
}

type NotificationMsg struct {
	SendID      string
	RecvID      string
	Content     []byte //  open_im_sdk.TipsComm
	MsgFrom     int32
	ContentType int32
	SessionType int32
	OperationID string
}

func Notification(n *NotificationMsg) {
	return
	var req pbChat.SendMsgReq
	var msg sdk_ws.MsgData
	var offlineInfo sdk_ws.OfflinePushInfo
	var title, desc, ex string
	var pushSwitch bool
	req.OperationID = n.OperationID
	msg.SendID = n.SendID
	msg.RecvID = n.RecvID
	msg.Content = n.Content
	msg.MsgFrom = n.MsgFrom
	msg.ContentType = n.ContentType
	msg.SessionType = n.SessionType
	msg.CreateTime = utils.GetCurrentTimestampByMill()
	msg.ClientMsgID = utils.GetMsgID(n.SendID)
	switch n.SessionType {
	case constant.GroupChatType:
		msg.RecvID = ""
		msg.GroupID = n.RecvID
	}
	if true {
		msg.Options = make(map[string]bool, 10)
		//utils.SetSwitchFromOptions(msg.Options, constant.IsOfflinePush, false)
		utils.SetSwitchFromOptions(msg.Options, constant.IsHistory, false)
		utils.SetSwitchFromOptions(msg.Options, constant.IsPersistent, false)
	}
	offlineInfo.IOSBadgeCount = config.Config.IOSPush.BadgeCount
	offlineInfo.IOSPushSound = config.Config.IOSPush.PushSound
	//switch msg.ContentType {
	//case constant.GroupCreatedNotification:
	//	pushSwitch = config.Config.Notification.GroupCreated.OfflinePush.PushSwitch
	//	title = config.Config.Notification.GroupCreated.OfflinePush.Title
	//	desc = config.Config.Notification.GroupCreated.OfflinePush.Desc
	//	ex = config.Config.Notification.GroupCreated.OfflinePush.Ext
	//case constant.GroupInfoChangedNotification:
	//	pushSwitch = config.Config.Notification.GroupInfoChanged.OfflinePush.PushSwitch
	//	title = config.Config.Notification.GroupInfoChanged.OfflinePush.Title
	//	desc = config.Config.Notification.GroupInfoChanged.OfflinePush.Desc
	//	ex = config.Config.Notification.GroupInfoChanged.OfflinePush.Ext
	//case constant.JoinApplicationNotification:
	//	pushSwitch = config.Config.Notification.ApplyJoinGroup.OfflinePush.PushSwitch
	//	title = config.Config.Notification.ApplyJoinGroup.OfflinePush.Title
	//	desc = config.Config.Notification.ApplyJoinGroup.OfflinePush.Desc
	//	ex = config.Config.Notification.ApplyJoinGroup.OfflinePush.Ext
	//}
	utils.SetSwitchFromOptions(msg.Options, constant.IsOfflinePush, pushSwitch)
	offlineInfo.Title = title
	offlineInfo.Desc = desc
	offlineInfo.Ex = ex
	msg.OfflinePushInfo = &offlineInfo
	req.MsgData = &msg
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfflineMessageName)
	client := pbChat.NewChatClient(etcdConn)
	reply, err := client.SendMsg(context.Background(), &req)
	if err != nil {
		log.NewError(req.OperationID, "SendMsg rpc failed, ", req.String(), err.Error())
	} else if reply.ErrCode != 0 {
		log.NewError(req.OperationID, "SendMsg rpc failed, ", req.String())
	}
}
