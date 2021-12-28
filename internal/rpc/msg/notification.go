package msg

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	pbFriend "Open_IM/pkg/proto/friend"
	pbGroup "Open_IM/pkg/proto/group"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"encoding/json"
)

//message GroupCreatedTips{
//  GroupInfo Group = 1;
//  GroupMemberFullInfo Creator = 2;
//  repeated GroupMemberFullInfo MemberList = 3;
//  uint64 OperationTime = 4;
//} creator->group

func setOpUserInfo(operationID, opUserID, groupID string, groupMemberInfo *open_im_sdk.GroupMemberFullInfo) {
	return
	if token_verify.IsMangerUserID(opUserID) {
		u, err := imdb.GetUserByUserID(opUserID)
		if err != nil {
			log.NewError(operationID, "FindUserByUID failed ", err.Error(), opUserID)
			return
		}
		utils.CopyStructFields(groupMemberInfo, u)
		groupMemberInfo.AppMangerLevel = 1
	} else {
		u, err := imdb.GetGroupMemberInfoByGroupIDAndUserID(groupID, opUserID)
		if err != nil {
			log.NewError(operationID, "FindGroupMemberInfoByGroupIdAndUserId failed ", err.Error(), groupID, opUserID)
			return
		}
		utils.CopyStructFields(groupMemberInfo, u)
	}
}

func setGroupInfo(operationID, groupID string, groupInfo *open_im_sdk.GroupInfo, ownerUserID string) {
	return
	group, err := imdb.GetGroupInfoByGroupID(groupID)
	if err != nil {
		log.NewError(operationID, "FindGroupInfoByGroupId failed ", err.Error(), groupID)
		return
	}
	utils.CopyStructFields(groupInfo, group)

	if ownerUserID != "" {
		groupInfo.OwnerUserID = ownerUserID
		//	setGroupPublicUserInfo(operationID, groupID, ownerUserID, groupInfo.Owner)
	}
}

func setGroupMemberInfo(operationID, groupID, userID string, groupMemberInfo *open_im_sdk.GroupMemberFullInfo) {
	return
	group, err := imdb.GetGroupMemberInfoByGroupIDAndUserID(groupID, userID)
	if err != nil {
		log.NewError(operationID, "FindGroupMemberInfoByGroupIdAndUserId failed ", err.Error(), groupID, userID)
		return
	}
	utils.CopyStructFields(groupMemberInfo, group)
}

//func setGroupPublicUserInfo(operationID, groupID, userID string, publicUserInfo *open_im_sdk.PublicUserInfo) {
//	group, err := imdb.GetGroupMemberInfoByGroupIDAndUserID(groupID, userID)
//	if err != nil {
//		log.NewError(operationID, "FindGroupMemberInfoByGroupIdAndUserId failed ", err.Error(), groupID, userID)
//		return
//	}
//	utils.CopyStructFields(publicUserInfo, group)
//}

//创建群后调用
func GroupCreatedNotification(operationID, opUserID, OwnerUserID, groupID string, initMemberList []string) {
	return
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
	return
	managerList, err := imdb.GetOwnerManagerByGroupID(req.GroupID)
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

	apply, err := imdb.GetUserByUserID(req.OpUserID)
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
	return
	var n NotificationMsg
	n.SendID = req.OpUserID
	n.ContentType = constant.ApplicationProcessedNotification
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = req.OperationID
	n.RecvID = req.FromUserID

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
	return
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
	return
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
	return
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
	return
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
	return
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
	return
	from, err1 := imdb.GetUserByUserID(fromUserID)
	to, err2 := imdb.GetUserByUserID(toUserID)
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
	return
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
	return
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
	return
	var n NotificationMsg
	n.SendID = fromUserID
	n.RecvID = toUserID
	n.ContentType = constant.FriendAddedNotification
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = operationID

	var FriendAddedTips open_im_sdk.FriendAddedTips

	user, err := imdb.GetUserByUserID(opUserID)
	if err != nil {
		log.NewError(operationID, "FindUserByUID failed ", err.Error(), opUserID)

	} else {
		utils.CopyStructFields(FriendAddedTips.OpUser, user)
	}

	friend, err := imdb.GetFriendRelationshipFromFriend(fromUserID, toUserID)
	if err != nil {
		log.NewError(operationID, "FindUserByUID failed ", err.Error(), fromUserID, toUserID)
	} else {
		FriendAddedTips.Friend.Remark = friend.Remark
	}

	from, err := imdb.GetUserByUserID(fromUserID)
	if err != nil {
		log.NewError(operationID, "FindUserByUID failed ", err.Error(), fromUserID)

	} else {
		utils.CopyStructFields(FriendAddedTips.Friend, from)
	}

	to, err := imdb.GetUserByUserID(toUserID)
	if err != nil {
		log.NewError(operationID, "FindUserByUID failed ", err.Error(), toUserID)

	} else {
		utils.CopyStructFields(FriendAddedTips.Friend.FriendUser, to)
	}

	fromUserNickname, toUserNickname := from.Nickname, to.Nickname
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
	return
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
func FriendInfoChangedNotification(operationID, opUserID, fromUserID, toUserID string) {
	return
	var n NotificationMsg
	n.SendID = fromUserID
	n.RecvID = toUserID
	n.ContentType = constant.FriendInfoChangedNotification
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = operationID

	var FriendInfoChangedTips open_im_sdk.FriendInfoChangedTips
	FriendInfoChangedTips.FromToUserID.FromUserID = fromUserID
	FriendInfoChangedTips.FromToUserID.ToUserID = toUserID
	fromUserNickname, toUserNickname := getFromToUserNickname(operationID, fromUserID, toUserID)
	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(FriendInfoChangedTips)
	tips.DefaultTips = fromUserNickname + " FriendDeletedNotification " + toUserNickname
	n.Content, _ = json.Marshal(tips)
	Notification(&n, true)
}

func BlackAddedNotification(req *pbFriend.AddBlacklistReq) {
	return
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
	return
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
	return
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
	u, err := imdb.GetUserByUserID(userID)
	if err != nil {
		log.NewError(operationID, "FindUserByUID failed ", err.Error(), userID)
	}

	tips.Detail, _ = json.Marshal(SelfInfoUpdatedTips)
	tips.DefaultTips = u.Nickname + " SelfInfoUpdatedNotification "
	n.Content, _ = json.Marshal(tips)
	Notification(&n, true)
}
