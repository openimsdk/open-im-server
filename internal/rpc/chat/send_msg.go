package chat

import (
	"Open_IM/internal/api/group"
	"Open_IM/internal/push/content_struct"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	http2 "Open_IM/pkg/common/http"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbChat "Open_IM/pkg/proto/chat"
	pbGroup "Open_IM/pkg/proto/group"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"encoding/json"
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

func (rpc *rpcChat) UserSendMsg(_ context.Context, pb *pbChat.UserSendMsgReq) (*pbChat.UserSendMsgResp, error) {
	replay := pbChat.UserSendMsgResp{}
	log.NewDebug(pb.OperationID, "rpc sendMsg come here", pb.String())
	//if !utils.VerifyToken(pb.Token, pb.SendID) {
	//	return returnMsg(&replay, pb, http.StatusUnauthorized, "token validate err,not authorized", "", 0)
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
	pbData.ClientMsgID = pb.ClientMsgID
	pbData.SendID = pb.SendID
	pbData.SenderNickName = pb.SenderNickName
	pbData.SenderFaceURL = pb.SenderFaceURL
	pbData.MsgID = serverMsgID
	pbData.OperationID = pb.OperationID
	pbData.Token = pb.Token
	if pb.SendTime == 0 {
		pbData.SendTime = utils.GetCurrentTimestampByNano()
	} else {
		pbData.SendTime = pb.SendTime
	}
	options := utils.JsonStringToMap(pbData.Options)
	isHistory := utils.GetSwitchFromOptions(options, "history")
	mReq := MsgCallBackReq{
		SendID:      pb.SendID,
		RecvID:      pb.RecvID,
		Content:     pb.Content,
		SendTime:    pbData.SendTime,
		MsgFrom:     pbData.MsgFrom,
		ContentType: pb.ContentType,
		SessionType: pb.SessionType,
		PlatformID:  pb.PlatformID,
		MsgID:       pb.ClientMsgID,
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
				pbData.Content = mResp.ResponseResult.ModifiedMsg
			}
		}
	}
	switch pbData.SessionType {
	case constant.SingleChatType:
		isSend := modifyMessageByUserMessageReceiveOpt(pbData.RecvID, pbData.SendID, constant.SingleChatType, &pbData)
		if isSend {
			err1 := rpc.sendMsgToKafka(&pbData, pbData.RecvID)
			if err1 != nil {
				log.NewError(pbData.OperationID, "kafka send msg err:RecvID", pbData.RecvID, pbData.String())
				return returnMsg(&replay, pb, 201, "kafka send msg err", "", 0)
			}
		}
		err2 := rpc.sendMsgToKafka(&pbData, pbData.SendID)
		if err2 != nil {
			log.NewError(pbData.OperationID, "kafka send msg err:SendID", pbData.SendID, pbData.String())
			return returnMsg(&replay, pb, 201, "kafka send msg err", "", 0)
		}
		return returnMsg(&replay, pb, 0, "", serverMsgID, pbData.SendTime)
	case constant.GroupChatType:
		etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
		client := pbGroup.NewGroupClient(etcdConn)
		req := &pbGroup.GetGroupAllMemberReq{
			GroupID:     pbData.RecvID,
			Token:       pbData.Token,
			OperationID: pbData.OperationID,
		}
		reply, err := client.GetGroupAllMember(context.Background(), req)
		if err != nil {
			log.Error(pbData.Token, pbData.OperationID, "rpc send_msg getGroupInfo failed, err = %s", err.Error())
			return returnMsg(&replay, pb, 201, err.Error(), "", 0)
		}
		if reply.ErrorCode != 0 {
			log.Error(pbData.Token, pbData.OperationID, "rpc send_msg getGroupInfo failed, err = %s", reply.ErrorMsg)
			return returnMsg(&replay, pb, reply.ErrorCode, reply.ErrorMsg, "", 0)
		}
		var addUidList []string
		switch pbData.ContentType {
		case constant.KickGroupMemberTip:
			var notification content_struct.NotificationContent
			var kickContent group.KickGroupMemberReq
			err := utils.JsonStringToStruct(pbData.Content, &notification)
			if err != nil {
				log.ErrorByKv("json unmarshall err", pbData.OperationID, "err", err.Error())
				return returnMsg(&replay, pb, 200, err.Error(), "", 0)
			} else {
				err := utils.JsonStringToStruct(notification.Detail, &kickContent)
				if err != nil {
					log.ErrorByKv("json unmarshall err", pbData.OperationID, "err", err.Error())
					return returnMsg(&replay, pb, 200, err.Error(), "", 0)
				}
				for _, v := range kickContent.UidListInfo {
					addUidList = append(addUidList, v.UserId)
				}
			}
		case constant.QuitGroupTip:
			addUidList = append(addUidList, pbData.SendID)
		default:
		}
		groupID := pbData.RecvID
		for i, v := range reply.MemberList {
			pbData.RecvID = v.UserId + " " + groupID
			isSend := modifyMessageByUserMessageReceiveOpt(v.UserId, groupID, constant.GroupChatType, &pbData)
			if isSend {
				err := rpc.sendMsgToKafka(&pbData, utils.IntToString(i))
				if err != nil {
					log.NewError(pbData.OperationID, "kafka send msg err:UserId", v.UserId, pbData.String())
					return returnMsg(&replay, pb, 201, "kafka send msg err", "", 0)
				}
			}

		}
		for i, v := range addUidList {
			pbData.RecvID = v + " " + groupID
			isSend := modifyMessageByUserMessageReceiveOpt(v, groupID, constant.GroupChatType, &pbData)
			if isSend {
				err := rpc.sendMsgToKafka(&pbData, utils.IntToString(i+1))
				if err != nil {
					log.NewError(pbData.OperationID, "kafka send msg err:UserId", v, pbData.String())
					return returnMsg(&replay, pb, 201, "kafka send msg err", "", 0)
				}
			}
		}
		return returnMsg(&replay, pb, 0, "", serverMsgID, pbData.SendTime)
	default:
		return returnMsg(&replay, pb, 203, "unkonwn sessionType", "", 0)

	}

}

type WSToMsgSvrChatMsg struct {
	SendID      string `protobuf:"bytes,1,opt,name=SendID" json:"SendID,omitempty"`
	RecvID      string `protobuf:"bytes,2,opt,name=RecvID" json:"RecvID,omitempty"`
	Content     string `protobuf:"bytes,3,opt,name=Content" json:"Content,omitempty"`
	MsgFrom     int32  `protobuf:"varint,5,opt,name=MsgFrom" json:"MsgFrom,omitempty"`
	ContentType int32  `protobuf:"varint,8,opt,name=ContentType" json:"ContentType,omitempty"`
	SessionType int32  `protobuf:"varint,9,opt,name=SessionType" json:"SessionType,omitempty"`
	OperationID string `protobuf:"bytes,10,opt,name=OperationID" json:"OperationID,omitempty"`
}

func CreateGroupNotification(sendID string, creator im_mysql_model.User, group im_mysql_model.Group, memberList []im_mysql_model.GroupMember) {
	var msg WSToMsgSvrChatMsg
	msg.OperationID = utils.OperationIDGenerator()
	msg.SendID = sendID
	msg.RecvID = group.GroupId
	msg.ContentType = constant.CreateGroupTip
	msg.SessionType = constant.GroupChatType
	msg.MsgFrom = constant.SysMsgType

	var tip open_im_sdk.CreateGroupTip
	tip.Group = &open_im_sdk.GroupInfoTip{}
	utils.CopyStructFields(tip.Group, group)
	tip.Creator = &open_im_sdk.UserInfoTip{}
	utils.CopyStructFields(tip.Creator, creator)
	for _, v := range memberList {
		var groupMemberInfo open_im_sdk.GroupMemberFullInfoTip
		utils.CopyStructFields(&groupMemberInfo, v)
		tip.MemberList = append(tip.MemberList, &groupMemberInfo)
	}

	msg.Content = utils.StructToJsonString(tip)
	var offlinePushInfo open_im_sdk.OfflinePushInfo
	offlinePushInfo.Title = "create group title"
	offlinePushInfo.Desc = "create group desc"
	offlinePushInfo.Ext = "create group ext"
	Notification(&msg, false, offlinePushInfo)
}

func Notification(m *WSToMsgSvrChatMsg, onlineUserOnly bool, offlinePushInfo open_im_sdk.OfflinePushInfo) {

}

func (rpc *rpcChat) sendMsgToKafka(m *pbChat.WSToMsgSvrChatMsg, key string) error {
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
func returnMsg(replay *pbChat.UserSendMsgResp, pb *pbChat.UserSendMsgReq, errCode int32, errMsg, serverMsgID string, sendTime int64) (*pbChat.UserSendMsgResp, error) {
	replay.ErrCode = errCode
	replay.ErrMsg = errMsg
	replay.ReqIdentifier = pb.ReqIdentifier
	replay.ClientMsgID = pb.ClientMsgID
	replay.ServerMsgID = serverMsgID
	replay.SendTime = sendTime
	return replay, nil
}
func modifyMessageByUserMessageReceiveOpt(userID, sourceID string, sessionType int, msg *pbChat.WSToMsgSvrChatMsg) bool {
	conversationID := utils.GetConversationIDBySessionType(sourceID, sessionType)
	opt, err := db.DB.GetSingleConversationMsgOpt(userID, conversationID)
	if err != nil {
		log.NewError(msg.OperationID, "GetSingleConversationMsgOpt from redis err", msg.String())
		return true
	}
	switch opt {
	case constant.ReceiveMessage:
		return true
	case constant.NotReceiveMessage:
		return false
	case constant.ReceiveNotNotifyMessage:
		options := utils.JsonStringToMap(msg.Options)
		if options == nil {
			options = make(map[string]int32, 2)
		}
		utils.SetSwitchFromOptions(options, "offlinePush", 0)
		msg.Options = utils.MapIntToJsonString(options)
		return true
	}

	return true
}
