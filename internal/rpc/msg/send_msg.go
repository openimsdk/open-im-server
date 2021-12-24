package msg

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	http2 "Open_IM/pkg/common/http"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbChat "Open_IM/pkg/proto/chat"
	pbFriend "Open_IM/pkg/proto/friend"
	pbGroup "Open_IM/pkg/proto/group"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
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
		msg.SendTime = utils.GetCurrentTimestampByNano()
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
		err2 := rpc.sendMsgToKafka(&msgToMQ, msgToMQ.MsgData.SendID)
		if err2 != nil {
			log.NewError(msgToMQ.OperationID, "kafka send msg err:SendID", msgToMQ.MsgData.SendID, msgToMQ.String())
			return returnMsg(&replay, pb, 201, "kafka send msg err", "", 0)
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
	Content     []byte
	MsgFrom     int32
	ContentType int32
	SessionType int32
	OperationID string
}

func Notification(n *NotificationMsg, onlineUserOnly bool) {
	var req pbChat.SendMsgReq
	var msg sdk_ws.MsgData
	var offlineInfo sdk_ws.OfflinePushInfo
	var title, desc, ext string
	var pushSwitch bool
	req.OperationID = n.OperationID
	msg.SendID = n.SendID
	msg.RecvID = n.RecvID
	msg.Content = n.Content
	msg.MsgFrom = n.MsgFrom
	msg.ContentType = n.ContentType
	msg.SessionType = n.SessionType
	msg.CreateTime = utils.GetCurrentTimestampByNano()
	msg.ClientMsgID = utils.GetMsgID(n.SendID)
	switch n.SessionType {
	case constant.GroupChatType:
		msg.RecvID = ""
		msg.GroupID = n.RecvID
	}
	if onlineUserOnly {
		msg.Options = make(map[string]bool, 10)
		utils.SetSwitchFromOptions(msg.Options, constant.IsOfflinePush, false)
		utils.SetSwitchFromOptions(msg.Options, constant.IsHistory, false)
		utils.SetSwitchFromOptions(msg.Options, constant.IsPersistent, false)
	}
	offlineInfo.IOSBadgeCount = config.Config.IOSPush.BadgeCount
	offlineInfo.IOSPushSound = config.Config.IOSPush.PushSound
	switch msg.ContentType {
	case constant.CreateGroupTip:
		pushSwitch = config.Config.Notification.GroupCreated.OfflinePush.PushSwitch
		title = config.Config.Notification.GroupCreated.OfflinePush.Title
		desc = config.Config.Notification.GroupCreated.OfflinePush.Desc
		ext = config.Config.Notification.GroupCreated.OfflinePush.Ext
	case constant.ChangeGroupInfoTip:
		pushSwitch = config.Config.Notification.GroupInfoChanged.OfflinePush.PushSwitch
		title = config.Config.Notification.GroupInfoChanged.OfflinePush.Title
		desc = config.Config.Notification.GroupInfoChanged.OfflinePush.Desc
		ext = config.Config.Notification.GroupInfoChanged.OfflinePush.Ext
	case constant.ApplyJoinGroupTip:
		pushSwitch = config.Config.Notification.ApplyJoinGroup.OfflinePush.PushSwitch
		title = config.Config.Notification.ApplyJoinGroup.OfflinePush.Title
		desc = config.Config.Notification.ApplyJoinGroup.OfflinePush.Desc
		ext = config.Config.Notification.ApplyJoinGroup.OfflinePush.Ext
	}
	utils.SetSwitchFromOptions(msg.Options, constant.IsOfflinePush, pushSwitch)
	offlineInfo.Title = title
	offlineInfo.Desc = desc
	offlineInfo.Ext = ext
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

//message GroupCreatedTips{
//  GroupInfo Group = 1;
//  GroupMemberFullInfo Creator = 2;
//  repeated GroupMemberFullInfo MemberList = 3;
//  uint64 OperationTime = 4;
//} creator->group

func setOpUserInfo(operationID, opUserID, groupID string, groupMemberInfo *open_im_sdk.GroupMemberFullInfo) {
	if token_verify.IsMangerUserID(opUserID) {
		u, err := imdb.FindUserByUID(opUserID)
		if err != nil {
			log.NewError(operationID, "FindUserByUID failed ", err.Error(), opUserID)
			return
		}
		utils.CopyStructFields(groupMemberInfo, u)
		groupMemberInfo.AppMangerLevel = 1
	} else {
		u, err := imdb.FindGroupMemberInfoByGroupIdAndUserId(groupID, opUserID)
		if err != nil {
			log.NewError(operationID, "FindGroupMemberInfoByGroupIdAndUserId failed ", err.Error(), groupID, opUserID)
			return
		}
		utils.CopyStructFields(groupMemberInfo, u)
	}
}

func setGroupInfo(operationID, groupID string, groupInfo *open_im_sdk.GroupInfo, ownerUserID string) {
	group, err := imdb.FindGroupInfoByGroupId(groupID)
	if err != nil {
		log.NewError(operationID, "FindGroupInfoByGroupId failed ", err.Error(), groupID)
		return
	}
	utils.CopyStructFields(groupInfo, group)

	if ownerUserID != "" {
		groupInfo.Owner = &open_im_sdk.PublicUserInfo{}
		setGroupPublicUserInfo(operationID, groupID, ownerUserID, groupInfo.Owner)
	}
}

func setGroupMemberInfo(operationID, groupID, userID string, groupMemberInfo *open_im_sdk.GroupMemberFullInfo) {
	group, err := imdb.FindGroupMemberInfoByGroupIdAndUserId(groupID, userID)
	if err != nil {
		log.NewError(operationID, "FindGroupMemberInfoByGroupIdAndUserId failed ", err.Error(), groupID, userID)
		return
	}
	utils.CopyStructFields(groupMemberInfo, group)
}

func setGroupPublicUserInfo(operationID, groupID, userID string, publicUserInfo *open_im_sdk.PublicUserInfo) {
	group, err := imdb.FindGroupMemberInfoByGroupIdAndUserId(groupID, userID)
	if err != nil {
		log.NewError(operationID, "FindGroupMemberInfoByGroupIdAndUserId failed ", err.Error(), groupID, userID)
		return
	}
	utils.CopyStructFields(publicUserInfo, group)
}

//创建群后调用
func GroupCreatedNotification(operationID, opUserID, OwnerUserID, groupID string, initMemberList []string) {
	var n NotificationMsg
	n.SendID = opUserID
	n.RecvID = groupID
	n.ContentType = constant.GroupCreatedNotification
	n.SessionType = constant.GroupChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = operationID

	GroupCreatedTips := open_im_sdk.GroupCreatedTips{Group: &open_im_sdk.GroupInfo{},
		Creator: &open_im_sdk.GroupMemberFullInfo{}}
	setOpUserInfo(operationID, GroupCreatedTips.Creator.UserID, groupID, GroupCreatedTips.Creator)

	setGroupInfo(operationID, groupID, GroupCreatedTips.Group, OwnerUserID)

	for _, v := range initMemberList {
		var groupMemberInfo open_im_sdk.GroupMemberFullInfo
		setGroupMemberInfo(operationID, groupID, v, &groupMemberInfo)
		GroupCreatedTips.MemberList = append(GroupCreatedTips.MemberList, &groupMemberInfo)
	}

	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(GroupCreatedTips)
	tips.DefaultTips = config.Config.Notification.GroupCreated.DefaultTips.Tips
	n.Content, _ = json.Marshal(tips)
	Notification(&n, false)
}

//message ReceiveJoinApplicationTips{
//  GroupInfo Group = 1;
//  PublicUserInfo Applicant  = 2;
//  string 	Reason = 3;
//}  apply->all managers GroupID              string   `protobuf:"bytes,1,opt,name=GroupID" json:"GroupID,omitempty"`
//	ReqMessage           string   `protobuf:"bytes,2,opt,name=ReqMessage" json:"ReqMessage,omitempty"`
//	OpUserID             string   `protobuf:"bytes,3,opt,name=OpUserID" json:"OpUserID,omitempty"`
//	OperationID          string   `protobuf:"bytes,4,opt,name=OperationID" json:"OperationID,omitempty"`
//申请进群后调用
func JoinApplicationNotification(req *pbGroup.JoinGroupReq) {
	managerList, err := imdb.GetOwnerManagerByGroupId(req.GroupID)
	if err != nil {
		log.NewError(req.OperationID, "GetOwnerManagerByGroupId failed ", err.Error(), req.GroupID)
		return
	}

	var n NotificationMsg
	n.SendID = req.OpUserID
	n.ContentType = constant.JoinApplicationNotification
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = req.OperationID

	JoinGroupApplicationTips := open_im_sdk.JoinGroupApplicationTips{Group: &open_im_sdk.GroupInfo{}, Applicant: &open_im_sdk.PublicUserInfo{}}
	setGroupInfo(req.OperationID, req.GroupID, JoinGroupApplicationTips.Group, "")

	apply, err := imdb.FindUserByUID(req.OpUserID)
	if err != nil {
		log.NewError(req.OperationID, "FindUserByUID failed ", err.Error(), req.OpUserID)
		return
	}
	utils.CopyStructFields(JoinGroupApplicationTips.Applicant, apply)
	JoinGroupApplicationTips.Reason = req.ReqMessage

	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(JoinGroupApplicationTips)
	tips.DefaultTips = "JoinGroupApplicationTips"
	n.Content, _ = json.Marshal(tips)
	for _, v := range managerList {
		n.RecvID = v.UserID
		Notification(&n, true)
	}
}

//message ApplicationProcessedTips{
//  GroupInfo Group = 1;
//  GroupMemberFullInfo OpUser = 2;
//  int32 Result = 3;
//  string 	Reason = 4;
//}
//处理进群请求后调用
func ApplicationProcessedNotification(req *pbGroup.GroupApplicationResponseReq) {
	var n NotificationMsg
	n.SendID = req.OpUserID
	n.ContentType = constant.ApplicationProcessedNotification
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = req.OperationID
	n.RecvID = req.ToUserID

	ApplicationProcessedTips := open_im_sdk.ApplicationProcessedTips{Group: &open_im_sdk.GroupInfo{}, OpUser: &open_im_sdk.GroupMemberFullInfo{}}
	setGroupInfo(req.OperationID, req.GroupID, ApplicationProcessedTips.Group, "")
	setOpUserInfo(req.OperationID, req.OpUserID, req.GroupID, ApplicationProcessedTips.OpUser)
	ApplicationProcessedTips.Reason = req.HandledMsg
	ApplicationProcessedTips.Result = req.HandleResult

	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(ApplicationProcessedTips)
	tips.DefaultTips = "ApplicationProcessedNotification"
	n.Content, _ = json.Marshal(tips)

	Notification(&n, true)
}

//message MemberInvitedTips{
//  GroupInfo Group = 1;
//  GroupMemberFullInfo OpUser = 2;
//  GroupMemberFullInfo InvitedUser = 3;
//  uint64 OperationTime = 4;
//}
//被邀请进群后调用
func MemberInvitedNotification(operationID, groupID, opUserID, reason string, invitedUserIDList []string) {
	var n NotificationMsg
	n.SendID = opUserID
	n.ContentType = constant.MemberInvitedNotification
	n.SessionType = constant.GroupChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = operationID

	ApplicationProcessedTips := open_im_sdk.MemberInvitedTips{Group: &open_im_sdk.GroupInfo{}, OpUser: &open_im_sdk.GroupMemberFullInfo{}}
	setGroupInfo(operationID, groupID, ApplicationProcessedTips.Group, "")
	setOpUserInfo(operationID, opUserID, groupID, ApplicationProcessedTips.OpUser)
	for _, v := range invitedUserIDList {
		var groupMemberInfo open_im_sdk.GroupMemberFullInfo
		setGroupMemberInfo(operationID, groupID, v, &groupMemberInfo)
		ApplicationProcessedTips.InvitedUserList = append(ApplicationProcessedTips.InvitedUserList, &groupMemberInfo)
	}
	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(ApplicationProcessedTips)
	tips.DefaultTips = "MemberInvitedNotification"
	n.Content, _ = json.Marshal(tips)
	n.RecvID = groupID
	Notification(&n, true)
}

//message MemberKickedTips{
//  GroupInfo Group = 1;
//  GroupMemberFullInfo OpUser = 2;
//  GroupMemberFullInfo KickedUser = 3;
//  uint64 OperationTime = 4;
//}
//被踢后调用
func MemberKickedNotification(req *pbGroup.KickGroupMemberReq, kickedUserIDList []string) {
	var n NotificationMsg
	n.SendID = req.OpUserID
	n.ContentType = constant.MemberKickedNotification
	n.SessionType = constant.GroupChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = req.OperationID

	MemberKickedTips := open_im_sdk.MemberKickedTips{Group: &open_im_sdk.GroupInfo{}, OpUser: &open_im_sdk.GroupMemberFullInfo{}}
	setGroupInfo(req.OperationID, req.GroupID, MemberKickedTips.Group, "")
	setOpUserInfo(req.OperationID, req.OpUserID, req.GroupID, MemberKickedTips.OpUser)
	for _, v := range kickedUserIDList {
		var groupMemberInfo open_im_sdk.GroupMemberFullInfo
		setGroupMemberInfo(req.OperationID, req.GroupID, v, &groupMemberInfo)
		MemberKickedTips.KickedUserList = append(MemberKickedTips.KickedUserList, &groupMemberInfo)
	}
	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(MemberKickedTips)
	tips.DefaultTips = "MemberKickedNotification"
	n.Content, _ = json.Marshal(tips)
	n.RecvID = req.GroupID
	Notification(&n, true)

	for _, v := range kickedUserIDList {
		n.SessionType = constant.SingleChatType
		n.RecvID = v
		Notification(&n, true)
	}
}

//message GroupInfoChangedTips{
//  int32 ChangedType = 1; //bitwise operators: 1:groupName; 10:Notification  100:Introduction; 1000:FaceUrl
//  GroupInfo Group = 2;
//  GroupMemberFullInfo OpUser = 3;
//}

//群信息改变后掉用
func GroupInfoChangedNotification(operationID, opUserID, groupID string, changedType int32) {
	var n NotificationMsg
	n.SendID = opUserID
	n.ContentType = constant.GroupInfoChangedNotification
	n.SessionType = constant.GroupChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = operationID

	GroupInfoChangedTips := open_im_sdk.GroupInfoChangedTips{Group: &open_im_sdk.GroupInfo{}, OpUser: &open_im_sdk.GroupMemberFullInfo{}}
	setGroupInfo(operationID, groupID, GroupInfoChangedTips.Group, opUserID)
	setOpUserInfo(operationID, opUserID, groupID, GroupInfoChangedTips.OpUser)
	GroupInfoChangedTips.ChangedType = changedType
	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(GroupInfoChangedTips)
	tips.DefaultTips = "GroupInfoChangedNotification"
	n.Content, _ = json.Marshal(tips)
	n.RecvID = groupID
	Notification(&n, false)
}

/*
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
*/

//message MemberLeaveTips{
//  GroupInfo Group = 1;
//  GroupMemberFullInfo LeaverUser = 2;
//  uint64 OperationTime = 3;
//}

//群成员退群后调用
func MemberLeaveNotification(req *pbGroup.QuitGroupReq) {
	var n NotificationMsg
	n.SendID = req.OpUserID
	n.ContentType = constant.MemberLeaveNotification
	n.SessionType = constant.GroupChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = req.OperationID

	MemberLeaveTips := open_im_sdk.MemberLeaveTips{Group: &open_im_sdk.GroupInfo{}, LeaverUser: &open_im_sdk.GroupMemberFullInfo{}}
	setGroupInfo(req.OperationID, req.GroupID, MemberLeaveTips.Group, "")
	setOpUserInfo(req.OperationID, req.OpUserID, req.GroupID, MemberLeaveTips.LeaverUser)

	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(MemberLeaveTips)
	tips.DefaultTips = "MemberLeaveNotification"
	n.Content, _ = json.Marshal(tips)
	n.RecvID = req.GroupID
	Notification(&n, true)

	n.SessionType = constant.SingleChatType
	n.RecvID = req.OpUserID
	Notification(&n, true)
}

//message MemberEnterTips{
//  GroupInfo Group = 1;
//  GroupMemberFullInfo EntrantUser = 2;
//  uint64 OperationTime = 3;
//}
//群成员主动申请进群，管理员同意后调用，
func MemberEnterNotification(req *pbGroup.GroupApplicationResponseReq) {
	var n NotificationMsg
	n.SendID = req.OpUserID
	n.ContentType = constant.MemberEnterNotification
	n.SessionType = constant.GroupChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = req.OperationID

	MemberLeaveTips := open_im_sdk.MemberEnterTips{Group: &open_im_sdk.GroupInfo{}, EntrantUser: &open_im_sdk.GroupMemberFullInfo{}}
	setGroupInfo(req.OperationID, req.GroupID, MemberLeaveTips.Group, "")
	setOpUserInfo(req.OperationID, req.OpUserID, req.GroupID, MemberLeaveTips.EntrantUser)

	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(MemberLeaveTips)
	tips.DefaultTips = "MemberEnterNotification"
	n.Content, _ = json.Marshal(tips)
	n.RecvID = req.GroupID
	Notification(&n, true)

}

//message MemberInfoChangedTips{
//  int32 ChangeType = 1; //1:info changed; 2:mute
//  GroupMemberFullInfo OpUser = 2; //who do this
//  GroupMemberFullInfo FinalInfo = 3; //
//  uint64 MuteTime = 4;
//  GroupInfo Group = 5;
//}
//func MemberInfoChangedNotification(operationID string, group *immysql.Group, opUser *immysql.GroupMember, userFinalInfo *immysql.GroupMember) {

//}

//message FriendApplicationAddedTips{
//  PublicUserInfo OpUser = 1; //user1
//  FriendApplication Application = 2;
//  PublicUserInfo  OpedUser = 3; //user2
//}

func getFromToUserNickname(operationID, fromUserID, toUserID string) (string, string) {
	from, err1 := imdb.FindUserByUID(fromUserID)
	to, err2 := imdb.FindUserByUID(toUserID)
	if err1 != nil || err2 != nil {
		log.NewError("FindUserByUID failed ", err1, err2, fromUserID, toUserID)
	}
	fromNickname, toNickname := "", ""
	if from != nil {
		fromNickname = from.Nickname
	}
	if to != nil {
		toNickname = to.Nickname
	}
	return fromNickname, toNickname
}

func FriendApplicationAddedNotification(req *pbFriend.AddFriendReq) {
	var n NotificationMsg
	n.SendID = req.CommID.FromUserID
	n.RecvID = req.CommID.ToUserID
	n.ContentType = constant.FriendApplicationAddedNotification
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = req.CommID.OperationID

	var FriendApplicationAddedTips open_im_sdk.FriendApplicationAddedTips
	FriendApplicationAddedTips.FromToUserID.FromUserID = req.CommID.FromUserID
	FriendApplicationAddedTips.FromToUserID.ToUserID = req.CommID.ToUserID
	fromUserNickname, toUserNickname := getFromToUserNickname(req.CommID.OperationID, req.CommID.FromUserID, req.CommID.ToUserID)
	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(FriendApplicationAddedTips)
	tips.DefaultTips = fromUserNickname + " FriendApplicationAddedNotification " + toUserNickname
	n.Content, _ = json.Marshal(tips)
	Notification(&n, true)
}

func FriendApplicationProcessedNotification(req *pbFriend.AddFriendResponseReq) {
	var n NotificationMsg
	n.SendID = req.CommID.FromUserID
	n.RecvID = req.CommID.ToUserID
	n.ContentType = constant.FriendApplicationProcessedNotification
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = req.CommID.OperationID

	var FriendApplicationProcessedTips open_im_sdk.FriendApplicationProcessedTips
	FriendApplicationProcessedTips.FromToUserID.FromUserID = req.CommID.FromUserID
	FriendApplicationProcessedTips.FromToUserID.ToUserID = req.CommID.ToUserID
	fromUserNickname, toUserNickname := getFromToUserNickname(req.CommID.OperationID, req.CommID.FromUserID, req.CommID.ToUserID)
	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(FriendApplicationProcessedTips)
	tips.DefaultTips = fromUserNickname + " FriendApplicationProcessedNotification " + toUserNickname
	n.Content, _ = json.Marshal(tips)
	Notification(&n, true)
}

func FriendAddedNotification(operationID, opUserID, fromUserID, toUserID string) {
	var n NotificationMsg
	n.SendID = fromUserID
	n.RecvID = toUserID
	n.ContentType = constant.FriendAddedNotification
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = operationID

	var FriendAddedTips open_im_sdk.FriendAddedTips

	user, err := imdb.FindUserByUID(opUserID)
	if err != nil {
		log.NewError(operationID, "FindUserByUID failed ", err.Error(), opUserID)

	} else {
		utils.CopyStructFields(FriendAddedTips.OpUser, user)
	}

	friend, err := imdb.FindFriendRelationshipFromFriend(fromUserID, toUserID)
	if err != nil {
		log.NewError(operationID, "FindUserByUID failed ", err.Error(), fromUserID, toUserID)
	} else {
		FriendAddedTips.Friend.Remark = friend.Remark
	}

	from, err := imdb.FindUserByUID(fromUserID)
	if err != nil {
		log.NewError(operationID, "FindUserByUID failed ", err.Error(), fromUserID)

	} else {
		utils.CopyStructFields(FriendAddedTips.Friend.OwnerUser, from)
	}

	to, err := imdb.FindUserByUID(toUserID)
	if err != nil {
		log.NewError(operationID, "FindUserByUID failed ", err.Error(), toUserID)

	} else {
		utils.CopyStructFields(FriendAddedTips.Friend.FriendUser, to)
	}

	fromUserNickname, toUserNickname := FriendAddedTips.Friend.OwnerUser.Nickname, FriendAddedTips.Friend.FriendUser.Nickname
	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(FriendAddedTips)
	tips.DefaultTips = fromUserNickname + " FriendAddedNotification " + toUserNickname
	n.Content, _ = json.Marshal(tips)
	Notification(&n, true)
}

//message FriendDeletedTips{
//  FriendInfo Friend = 1;
//}
func FriendDeletedNotification(req *pbFriend.DeleteFriendReq) {
	var n NotificationMsg
	n.SendID = req.CommID.FromUserID
	n.RecvID = req.CommID.ToUserID
	n.ContentType = constant.FriendDeletedNotification
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = req.CommID.OperationID

	var FriendDeletedTips open_im_sdk.FriendDeletedTips
	FriendDeletedTips.FromToUserID.FromUserID = req.CommID.FromUserID
	FriendDeletedTips.FromToUserID.ToUserID = req.CommID.ToUserID
	fromUserNickname, toUserNickname := getFromToUserNickname(req.CommID.OperationID, req.CommID.FromUserID, req.CommID.ToUserID)
	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(FriendDeletedTips)
	tips.DefaultTips = fromUserNickname + " FriendDeletedNotification " + toUserNickname
	n.Content, _ = json.Marshal(tips)
	Notification(&n, true)
}

//message FriendInfoChangedTips{
//  FriendInfo Friend = 1;
//  PublicUserInfo OpUser = 2;
//  uint64 OperationTime = 3;
//}
func FriendInfoChangedNotification(req *pbFriend.SetFriendCommentReq) {
	var n NotificationMsg
	n.SendID = req.CommID.FromUserID
	n.RecvID = req.CommID.ToUserID
	n.ContentType = constant.FriendInfoChangedNotification
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = req.CommID.OperationID

	var FriendInfoChangedTips open_im_sdk.FriendInfoChangedTips
	FriendInfoChangedTips.FromToUserID.FromUserID = req.CommID.FromUserID
	FriendInfoChangedTips.FromToUserID.ToUserID = req.CommID.ToUserID
	fromUserNickname, toUserNickname := getFromToUserNickname(req.CommID.OperationID, req.CommID.FromUserID, req.CommID.ToUserID)
	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(FriendInfoChangedTips)
	tips.DefaultTips = fromUserNickname + " FriendDeletedNotification " + toUserNickname
	n.Content, _ = json.Marshal(tips)
	Notification(&n, true)
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
func BlackAddedNotification(req *pbFriend.AddBlacklistReq) {
	var n NotificationMsg
	n.SendID = req.CommID.FromUserID
	n.RecvID = req.CommID.ToUserID
	n.ContentType = constant.BlackAddedNotification
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = req.CommID.OperationID

	var BlackAddedTips open_im_sdk.BlackAddedTips
	BlackAddedTips.FromToUserID.FromUserID = req.CommID.FromUserID
	BlackAddedTips.FromToUserID.ToUserID = req.CommID.ToUserID
	fromUserNickname, toUserNickname := getFromToUserNickname(req.CommID.OperationID, req.CommID.FromUserID, req.CommID.ToUserID)
	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(BlackAddedTips)
	tips.DefaultTips = fromUserNickname + " BlackAddedNotification " + toUserNickname
	n.Content, _ = json.Marshal(tips)
	Notification(&n, true)
}

//message BlackDeletedTips{
//  BlackInfo Black = 1;
//}
func BlackDeletedNotification(req *pbFriend.RemoveBlacklistReq) {
	var n NotificationMsg
	n.SendID = req.CommID.FromUserID
	n.RecvID = req.CommID.ToUserID
	n.ContentType = constant.BlackDeletedNotification
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = req.CommID.OperationID

	var BlackDeletedTips open_im_sdk.BlackDeletedTips
	BlackDeletedTips.FromToUserID.FromUserID = req.CommID.FromUserID
	BlackDeletedTips.FromToUserID.ToUserID = req.CommID.ToUserID
	fromUserNickname, toUserNickname := getFromToUserNickname(req.CommID.OperationID, req.CommID.FromUserID, req.CommID.ToUserID)
	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(BlackDeletedTips)
	tips.DefaultTips = fromUserNickname + " BlackDeletedNotification " + toUserNickname
	n.Content, _ = json.Marshal(tips)
	Notification(&n, true)
}

//message SelfInfoUpdatedTips{
//  UserInfo SelfUserInfo = 1;
//  PublicUserInfo OpUser = 2;
//  uint64 OperationTime = 3;
//}
func SelfInfoUpdatedNotification(operationID, userID string) {
	var n NotificationMsg
	n.SendID = userID
	n.RecvID = userID
	n.ContentType = constant.SelfInfoUpdatedNotification
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = operationID

	var SelfInfoUpdatedTips open_im_sdk.SelfInfoUpdatedTips
	SelfInfoUpdatedTips.UserID = userID

	var tips open_im_sdk.TipsComm
	u, err := imdb.FindUserByUID(userID)
	if err != nil {
		log.NewError(operationID, "FindUserByUID failed ", err.Error(), userID)
	}

	tips.Detail, _ = json.Marshal(SelfInfoUpdatedTips)
	tips.DefaultTips = u.Nickname + " SelfInfoUpdatedNotification "
	n.Content, _ = json.Marshal(tips)
	Notification(&n, true)
}
