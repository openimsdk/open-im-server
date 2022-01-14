package msg

import (
	"Open_IM/pkg/common/constant"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	utils2 "Open_IM/pkg/common/utils"
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

func setOpUserInfo(opUserID, groupID string, groupMemberInfo *open_im_sdk.GroupMemberFullInfo) error {
	if token_verify.IsMangerUserID(opUserID) {
		u, err := imdb.GetUserByUserID(opUserID)
		if err != nil {
			return utils.Wrap(err, "GetUserByUserID failed")
		}
		utils.CopyStructFields(groupMemberInfo, u)
		groupMemberInfo.GroupID = groupID
	} else {
		u, err := imdb.GetGroupMemberInfoByGroupIDAndUserID(groupID, opUserID)
		if err != nil {
			return utils.Wrap(err, "GetGroupMemberInfoByGroupIDAndUserID failed")
		}
		if err = utils2.GroupMemberDBCopyOpenIM(groupMemberInfo, u); err != nil {
			return utils.Wrap(err, "")
		}
	}
	return nil
}

func setGroupInfo(groupID string, groupInfo *open_im_sdk.GroupInfo) error {
	group, err := imdb.GetGroupInfoByGroupID(groupID)
	if err != nil {
		return utils.Wrap(err, "GetGroupInfoByGroupID failed")
	}
	err = utils2.GroupDBCopyOpenIM(groupInfo, group)
	if err != nil {
		return utils.Wrap(err, "GetGroupMemberNumByGroupID failed")
	}
	return nil
}

func setGroupMemberInfo(groupID, userID string, groupMemberInfo *open_im_sdk.GroupMemberFullInfo) error {
	groupMember, err := imdb.GetGroupMemberInfoByGroupIDAndUserID(groupID, userID)
	if err != nil {
		return utils.Wrap(err, "")
	}
	if err = utils2.GroupMemberDBCopyOpenIM(groupMemberInfo, groupMember); err != nil {
		return utils.Wrap(err, "")
	}
	return nil
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
	var n NotificationMsg
	n.SendID = opUserID
	n.RecvID = groupID
	n.ContentType = constant.GroupCreatedNotification
	n.SessionType = constant.GroupChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = operationID

	GroupCreatedTips := open_im_sdk.GroupCreatedTips{Group: &open_im_sdk.GroupInfo{},
		Creator: &open_im_sdk.GroupMemberFullInfo{}}
	if err := setOpUserInfo(GroupCreatedTips.Creator.UserID, groupID, GroupCreatedTips.Creator); err != nil {
		log.NewError(operationID, "setOpUserInfo failed ", err.Error(), GroupCreatedTips.Creator.UserID, groupID, GroupCreatedTips.Creator)
	}
	err := setGroupInfo(groupID, GroupCreatedTips.Group)
	if err != nil {
		log.NewError(operationID, "setGroupInfo failed ", groupID, GroupCreatedTips.Group)
		return
	}
	for _, v := range initMemberList {
		var groupMemberInfo open_im_sdk.GroupMemberFullInfo
		setGroupMemberInfo(groupID, v, &groupMemberInfo)
		GroupCreatedTips.MemberList = append(GroupCreatedTips.MemberList, &groupMemberInfo)
	}
	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(GroupCreatedTips)
	tips.DefaultTips = config.Config.Notification.GroupCreated.DefaultTips.Tips
	n.Content, _ = json.Marshal(tips)
	log.NewInfo(operationID, "Notification ", n)
	Notification(&n)
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
	var n NotificationMsg
	n.SendID = req.OpUserID
	n.ContentType = constant.JoinApplicationNotification
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = req.OperationID

	JoinGroupApplicationTips := open_im_sdk.JoinGroupApplicationTips{Group: &open_im_sdk.GroupInfo{}, Applicant: &open_im_sdk.PublicUserInfo{}}
	err := setGroupInfo(req.GroupID, JoinGroupApplicationTips.Group)
	if err != nil {
		log.NewError(req.OperationID, "setGroupInfo failed ", err.Error(), req.GroupID, JoinGroupApplicationTips.Group)
		return
	}

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
	managerList, err := imdb.GetOwnerManagerByGroupID(req.GroupID)
	if err != nil {
		log.NewError(req.OperationID, "GetOwnerManagerByGroupId failed ", err.Error(), req.GroupID)
		return
	}
	for _, v := range managerList {
		n.RecvID = v.UserID
		log.NewInfo(req.OperationID, "Notification ", n)
		Notification(&n)
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
	n.RecvID = req.FromUserID
	ApplicationProcessedTips := open_im_sdk.ApplicationProcessedTips{Group: &open_im_sdk.GroupInfo{}, OpUser: &open_im_sdk.GroupMemberFullInfo{}}
	setGroupInfo(req.GroupID, ApplicationProcessedTips.Group)
	setOpUserInfo(req.OpUserID, req.GroupID, ApplicationProcessedTips.OpUser)
	ApplicationProcessedTips.Reason = req.HandledMsg
	ApplicationProcessedTips.Result = req.HandleResult

	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(ApplicationProcessedTips)
	tips.DefaultTips = "ApplicationProcessedNotification"
	n.Content, _ = json.Marshal(tips)

	Notification(&n)
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
	setGroupInfo(groupID, ApplicationProcessedTips.Group)
	setOpUserInfo(opUserID, groupID, ApplicationProcessedTips.OpUser)
	for _, v := range invitedUserIDList {
		var groupMemberInfo open_im_sdk.GroupMemberFullInfo
		setGroupMemberInfo(groupID, v, &groupMemberInfo)
		ApplicationProcessedTips.InvitedUserList = append(ApplicationProcessedTips.InvitedUserList, &groupMemberInfo)
	}
	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(ApplicationProcessedTips)
	tips.DefaultTips = "MemberInvitedNotification"
	n.Content, _ = json.Marshal(tips)
	n.RecvID = groupID
	Notification(&n)
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
	setGroupInfo(req.GroupID, MemberKickedTips.Group)
	setOpUserInfo(req.OpUserID, req.GroupID, MemberKickedTips.OpUser)
	for _, v := range kickedUserIDList {
		var groupMemberInfo open_im_sdk.GroupMemberFullInfo
		setGroupMemberInfo(req.GroupID, v, &groupMemberInfo)
		MemberKickedTips.KickedUserList = append(MemberKickedTips.KickedUserList, &groupMemberInfo)
	}
	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(MemberKickedTips)
	tips.DefaultTips = "MemberKickedNotification"
	n.Content, _ = json.Marshal(tips)
	n.RecvID = req.GroupID
	Notification(&n)

	for _, v := range kickedUserIDList {
		n.SessionType = constant.SingleChatType
		n.RecvID = v
		Notification(&n)
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
	setGroupInfo(groupID, GroupInfoChangedTips.Group)
	setOpUserInfo(opUserID, groupID, GroupInfoChangedTips.OpUser)
	GroupInfoChangedTips.ChangedType = changedType
	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(GroupInfoChangedTips)
	tips.DefaultTips = "GroupInfoChangedNotification"
	n.Content, _ = json.Marshal(tips)
	n.RecvID = groupID
	Notification(&n)
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
	setGroupInfo(req.GroupID, MemberLeaveTips.Group)
	setOpUserInfo(req.OpUserID, req.GroupID, MemberLeaveTips.LeaverUser)

	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(MemberLeaveTips)
	tips.DefaultTips = "MemberLeaveNotification"
	n.Content, _ = json.Marshal(tips)
	n.RecvID = req.GroupID
	Notification(&n)

	n.SessionType = constant.SingleChatType
	n.RecvID = req.OpUserID
	Notification(&n)
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
	setGroupInfo(req.GroupID, MemberLeaveTips.Group)
	setOpUserInfo(req.OpUserID, req.GroupID, MemberLeaveTips.EntrantUser)

	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(MemberLeaveTips)
	tips.DefaultTips = "MemberEnterNotification"
	n.Content, _ = json.Marshal(tips)
	n.RecvID = req.GroupID
	Notification(&n)

}
