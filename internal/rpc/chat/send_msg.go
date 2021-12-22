package chat

import (
	"Open_IM/internal/api/group"
	"Open_IM/internal/push/content_struct"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	immysql "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
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

type NotificationMsg struct {
	SendID      string
	RecvID      string
	Content     []byte
	MsgFrom     int32
	ContentType int32
	SessionType int32
	OperationID string
}

func Notification(n *NotificationMsg, onlineUserOnly bool) {

}

//message GroupCreatedTips{
//  GroupInfo Group = 1;
//  GroupMemberFullInfo Creator = 2;
//  repeated GroupMemberFullInfo MemberList = 3;
//  uint64 OperationTime = 4;
//}
func GroupCreatedNotification(operationID string, creator *immysql.User, group *immysql.Group, memberList []immysql.GroupMember) {
	var n NotificationMsg
	n.SendID = creator.UserID
	n.RecvID = group.GroupID
	n.ContentType = constant.CreateGroupTip
	n.SessionType = constant.GroupChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = operationID

	var groupCreated open_im_sdk.GroupCreatedTips
	groupCreated.Group = &open_im_sdk.GroupInfo{}
	utils.CopyStructFields(groupCreated.Group, group)
	groupCreated.Creator = &open_im_sdk.GroupMemberFullInfo{}
	utils.CopyStructFields(groupCreated.Creator, creator)
	for _, v := range memberList {
		var groupMemberInfo open_im_sdk.GroupMemberFullInfo
		utils.CopyStructFields(&groupMemberInfo, v)
		groupCreated.MemberList = append(groupCreated.MemberList, &groupMemberInfo)
	}
	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(groupCreated)
	tips.DefaultTips = config.Config.Notification.GroupCreated.DefaultTips.Tips
	n.Content, _ = json.Marshal(tips)
	Notification(&n, false)
}

//message ReceiveJoinApplicationTips{
//  GroupInfo Group = 1;
//  PublicUserInfo Applicant  = 2;
//  string 	Reason = 3;
//}
func ReceiveJoinApplicationNotification(operationID, RecvID string, applicant *immysql.User, group *immysql.Group) {
	var n NotificationMsg
	n.SendID = applicant.UserID
	n.RecvID = RecvID
	n.ContentType = constant.ApplyJoinGroupTip
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = operationID

	var joniGroup open_im_sdk.ReceiveJoinApplicationTips
	joniGroup.Group = &open_im_sdk.GroupInfo{}
	utils.CopyStructFields(joniGroup.Group, group)
	joniGroup.Applicant = &open_im_sdk.PublicUserInfo{}
	utils.CopyStructFields(joniGroup.Applicant, applicant)

	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(joniGroup)
	tips.DefaultTips = config.Config.Notification.ApplyJoinGroup.DefaultTips.Tips
	n.Content, _ = json.Marshal(tips)
	Notification(&n, false)
}

//message ApplicationProcessedTips{
//  GroupInfo Group = 1;
//  GroupMemberFullInfo OpUser = 2;
//  int32 Result = 3;
//  string 	Reason = 4;
//}
func ApplicationProcessedNotification(operationID, RecvID string, group immysql.Group, opUser immysql.GroupMember, result int32, Reason string) {

}

//message MemberInvitedTips{
//  GroupInfo Group = 1;
//  GroupMemberFullInfo OpUser = 2;
//  GroupMemberFullInfo InvitedUser = 3;
//  uint64 OperationTime = 4;
//}
func MemberInvitedNotification(operationID string, group immysql.Group, opUser immysql.GroupMember, invitedUser immysql.GroupMember) {

}

//message MemberKickedTips{
//  GroupInfo Group = 1;
//  GroupMemberFullInfo OpUser = 2;
//  GroupMemberFullInfo KickedUser = 3;
//  uint64 OperationTime = 4;
//}
func MemberKickedNotification(operationID string, group immysql.Group, opUser immysql.GroupMember, KickedUser immysql.GroupMember) {

}

//message GroupInfoChangedTips{
//  int32 ChangedType = 1; //bitwise operators: 1:groupName; 10:Notification  100:Introduction; 1000:FaceUrl
//  GroupInfo Group = 2;
//  GroupMemberFullInfo OpUser = 3;
//}
func GroupInfoChangedNotification(operationID string, changedType int32, group *immysql.Group, opUser *immysql.GroupMember) {
	var n NotificationMsg
	n.SendID = opUser.UserID
	n.RecvID = group.GroupID
	n.ContentType = constant.ChangeGroupInfoTip
	n.SessionType = constant.GroupChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = operationID

	var groupInfoChanged open_im_sdk.GroupInfoChangedTips
	groupInfoChanged.Group = &open_im_sdk.GroupInfo{}
	utils.CopyStructFields(groupInfoChanged.Group, group)
	groupInfoChanged.OpUser = &open_im_sdk.GroupMemberFullInfo{}
	utils.CopyStructFields(groupInfoChanged.OpUser, opUser)
	groupInfoChanged.ChangedType = changedType

	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(groupInfoChanged)
	tips.DefaultTips = config.Config.Notification.GroupInfoChanged.DefaultTips.Tips
	n.Content, _ = json.Marshal(tips)
	Notification(&n, false)
}

//message MemberLeaveTips{
//  GroupInfo Group = 1;
//  GroupMemberFullInfo LeaverUser = 2;
//  uint64 OperationTime = 3;
//}
func MemberLeaveNotification(operationID string, group *immysql.Group, leaverUser *immysql.GroupMember) {

}

//message MemberEnterTips{
//  GroupInfo Group = 1;
//  GroupMemberFullInfo EntrantUser = 2;
//  uint64 OperationTime = 3;
//}
func MemberEnterNotification(operationID string, group *immysql.Group, entrantUser *immysql.GroupMember) {

}

//message MemberInfoChangedTips{
//  int32 ChangeType = 1; //1:info changed; 2:mute
//  GroupMemberFullInfo OpUser = 2; //who do this
//  GroupMemberFullInfo FinalInfo = 3; //
//  uint64 MuteTime = 4;
//  GroupInfo Group = 5;
//}
func MemberInfoChangedNotification(operationID string, group *immysql.Group, opUser *immysql.GroupMember, userFinalInfo *immysql.GroupMember) {

}

//message FriendApplicationAddedTips{
//  PublicUserInfo OpUser = 1; //user1
//  FriendApplication Application = 2;
//  PublicUserInfo  OpedUser = 3; //user2
//}
func FriendApplicationAddedNotification(operationID string, opUser *immysql.User, opedUser *immysql.User, application *immysql.FriendRequest) {

}

//message FriendApplicationProcessedTips{
//  PublicUserInfo     OpUser = 1;  //user2
//  PublicUserInfo     OpedUser = 2; //user1
//  int32 result = 3; //1: accept; -1: reject
//}
func FriendApplicationProcessedNotification(operationID string, opUser *immysql.User, OpedUser *immysql.User, result int32) {

}

//message FriendAddedTips{
//  FriendInfo Friend = 1;
//}
//message FriendInfo{
//  UserInfo OwnerUser = 1;
//  string Remark = 2;
//  uint64 CreateTime = 3;
//  UserInfo FriendUser = 4;
//}

func FriendAddedNotification(operationID string, opUser *immysql.User, friendUser *immysql.Friend) {

}

//message FriendDeletedTips{
//  FriendInfo Friend = 1;
//}
func FriendDeletedNotification(operationID string, opUser *immysql.User, friendUser *immysql.Friend) {

}

//message FriendInfoChangedTips{
//  FriendInfo Friend = 1;
//  PublicUserInfo OpUser = 2;
//  uint64 OperationTime = 3;
//}
func FriendInfoChangedNotification(operationID string, opUser *immysql.User, friendUser *immysql.Friend) {

}

//message BlackAddedTips{
//    BlackInfo Black = 1;
//}
//message BlackInfo{
//  PublicUserInfo OwnerUser = 1;
//  string Remark = 2;
//  uint64 CreateTime = 3;
//  PublicUserInfo BlackUser = 4;
//}
func BlackAddedNotification(operationID string, opUser *immysql.User, blackUser *immysql.User) {

}

//message BlackDeletedTips{
//  BlackInfo Black = 1;
//}
func BlackDeletedNotification(operationID string, opUser *immysql.User, blackUser *immysql.User) {

}

//message SelfInfoUpdatedTips{
//  UserInfo SelfUserInfo = 1;
//  PublicUserInfo OpUser = 2;
//  uint64 OperationTime = 3;
//}
func SelfInfoUpdatedNotification(operationID string, opUser *immysql.User, selfUser *immysql.User) {

}
