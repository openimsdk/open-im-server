package notification

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/tokenverify"
	"Open_IM/pkg/common/tracelog"
	pbGroup "Open_IM/pkg/proto/group"
	"Open_IM/pkg/proto/sdkws"
	"Open_IM/pkg/utils"
	"context"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (c *Check) setOpUserInfo(opUserID, groupID string, groupMemberInfo *sdkws.GroupMemberFullInfo) error {
	if tokenverify.IsManagerUserID(opUserID) {
		user, err := c.user.GetUsersInfos(context.Background(), []string{opUserID}, true)
		if err != nil {
			return err
		}
		groupMemberInfo.GroupID = groupID
		groupMemberInfo.UserID = user[0].UserID
		groupMemberInfo.Nickname = user[0].Nickname
		groupMemberInfo.AppMangerLevel = user[0].AppMangerLevel
		groupMemberInfo.FaceURL = user[0].FaceURL
		return nil
	}
	u, err := c.group.GetGroupMemberInfo(context.Background(), groupID, opUserID)
	if err == nil {
		*groupMemberInfo = *u
		return nil
	}
	user, err := c.user.GetUsersInfos(context.Background(), []string{opUserID}, true)
	if err != nil {
		return err
	}
	groupMemberInfo.GroupID = groupID
	groupMemberInfo.UserID = user[0].UserID
	groupMemberInfo.Nickname = user[0].Nickname
	groupMemberInfo.AppMangerLevel = user[0].AppMangerLevel
	groupMemberInfo.FaceURL = user[0].FaceURL

	return nil
}

func (c *Check) setGroupInfo(groupID string, groupInfo *sdkws.GroupInfo) error {
	group, err := c.group.GetGroupInfos(context.Background(), []string{groupID}, true)
	if err != nil {
		return err
	}
	*groupInfo = *group[0]
	return nil
}

func (c *Check) setGroupMemberInfo(groupID, userID string, groupMemberInfo *sdkws.GroupMemberFullInfo) error {
	groupMember, err := c.group.GetGroupMemberInfo(context.Background(), groupID, userID)
	if err == nil {
		*groupMemberInfo = *groupMember
		return nil
	}
	user, err := c.user.GetUsersInfos(context.Background(), []string{userID}, true)
	if err != nil {
		return err
	}
	groupMemberInfo.GroupID = groupID
	groupMemberInfo.UserID = user[0].UserID
	groupMemberInfo.Nickname = user[0].Nickname
	groupMemberInfo.AppMangerLevel = user[0].AppMangerLevel
	groupMemberInfo.FaceURL = user[0].FaceURL
	return nil
}

func (c *Check) setGroupOwnerInfo(groupID string, groupMemberInfo *sdkws.GroupMemberFullInfo) error {
	group, err := c.group.GetGroupInfo(context.Background(), groupID)
	if err != nil {
		return err
	}
	groupMember, err := c.group.GetGroupMemberInfo(context.Background(), groupID, group.OwnerUserID)
	if err != nil {
		return err
	}
	*groupMemberInfo = *groupMember
	return nil
}

func (c *Check) setPublicUserInfo(userID string, publicUserInfo *sdkws.PublicUserInfo) error {
	user, err := c.user.GetPublicUserInfos(context.Background(), []string{userID}, true)
	if err != nil {
		return err
	}
	*publicUserInfo = *user[0]
	return nil
}

func (c *Check) groupNotification(contentType int32, m proto.Message, sendID, groupID, recvUserID, operationID string) {
	log.Info(operationID, utils.GetSelfFuncName(), "args: ", contentType, sendID, groupID, recvUserID)
	var err error
	var tips sdkws.TipsComm
	tips.Detail, err = proto.Marshal(m)
	if err != nil {
		log.Error(operationID, "Marshal failed ", err.Error(), m.String())
		return
	}
	marshaler := jsonpb.Marshaler{
		OrigName:     true,
		EnumsAsInts:  false,
		EmitDefaults: false,
	}
	tips.JsonDetail, _ = marshaler.MarshalToString(m)
	var nickname, toNickname string
	if sendID != "" {

		from, err := c.user.GetUsersInfos(context.Background(), []string{sendID}, true)
		if err != nil {
			return
		}
		nickname = from[0].Nickname
	}
	if recvUserID != "" {
		to, err := c.user.GetUsersInfos(context.Background(), []string{recvUserID}, true)
		if err != nil {
			return
		}
		toNickname = to[0].Nickname
	}

	cn := config.Config.Notification
	switch contentType {
	case constant.GroupCreatedNotification:
		tips.DefaultTips = nickname + " " + cn.GroupCreated.DefaultTips.Tips
	case constant.GroupInfoSetNotification:
		tips.DefaultTips = nickname + " " + cn.GroupInfoSet.DefaultTips.Tips
	case constant.JoinGroupApplicationNotification:
		tips.DefaultTips = nickname + " " + cn.JoinGroupApplication.DefaultTips.Tips
	case constant.MemberQuitNotification:
		tips.DefaultTips = nickname + " " + cn.MemberQuit.DefaultTips.Tips
	case constant.GroupApplicationAcceptedNotification: //
		tips.DefaultTips = toNickname + " " + cn.GroupApplicationAccepted.DefaultTips.Tips
	case constant.GroupApplicationRejectedNotification: //
		tips.DefaultTips = toNickname + " " + cn.GroupApplicationRejected.DefaultTips.Tips
	case constant.GroupOwnerTransferredNotification: //
		tips.DefaultTips = toNickname + " " + cn.GroupOwnerTransferred.DefaultTips.Tips
	case constant.MemberKickedNotification: //
		tips.DefaultTips = toNickname + " " + cn.MemberKicked.DefaultTips.Tips
	case constant.MemberInvitedNotification: //
		tips.DefaultTips = toNickname + " " + cn.MemberInvited.DefaultTips.Tips
	case constant.MemberEnterNotification:
		tips.DefaultTips = toNickname + " " + cn.MemberEnter.DefaultTips.Tips
	case constant.GroupDismissedNotification:
		tips.DefaultTips = toNickname + "" + cn.GroupDismissed.DefaultTips.Tips
	case constant.GroupMutedNotification:
		tips.DefaultTips = toNickname + "" + cn.GroupMuted.DefaultTips.Tips
	case constant.GroupCancelMutedNotification:
		tips.DefaultTips = toNickname + "" + cn.GroupCancelMuted.DefaultTips.Tips
	case constant.GroupMemberMutedNotification:
		tips.DefaultTips = toNickname + "" + cn.GroupMemberMuted.DefaultTips.Tips
	case constant.GroupMemberCancelMutedNotification:
		tips.DefaultTips = toNickname + "" + cn.GroupMemberCancelMuted.DefaultTips.Tips
	case constant.GroupMemberInfoSetNotification:
		tips.DefaultTips = toNickname + "" + cn.GroupMemberInfoSet.DefaultTips.Tips
	case constant.GroupMemberSetToAdminNotification:
		tips.DefaultTips = toNickname + "" + cn.GroupMemberSetToAdmin.DefaultTips.Tips
	case constant.GroupMemberSetToOrdinaryUserNotification:
		tips.DefaultTips = toNickname + "" + cn.GroupMemberSetToOrdinary.DefaultTips.Tips
	default:
		log.Error(operationID, "contentType failed ", contentType)
		return
	}

	var n NotificationMsg
	n.SendID = sendID
	if groupID != "" {
		n.RecvID = groupID

		group, err := c.group.GetGroupInfo(context.Background(), groupID)
		if err != nil {
			return
		}
		switch group.GroupType {
		case constant.NormalGroup:
			n.SessionType = constant.GroupChatType
		default:
			n.SessionType = constant.SuperGroupChatType
		}
	} else {
		n.RecvID = recvUserID
		n.SessionType = constant.SingleChatType
	}
	n.ContentType = contentType
	n.OperationID = operationID
	n.Content, err = proto.Marshal(&tips)
	if err != nil {
		log.Error(operationID, "Marshal failed ", err.Error(), tips.String())
		return
	}
	c.Notification(&n)
}

// 创建群后调用
func (c *Check) GroupCreatedNotification(operationID, opUserID, groupID string, initMemberList []string) {
	GroupCreatedTips := sdkws.GroupCreatedTips{Group: &sdkws.GroupInfo{},
		OpUser: &sdkws.GroupMemberFullInfo{}, GroupOwnerUser: &sdkws.GroupMemberFullInfo{}}
	if err := c.setOpUserInfo(opUserID, groupID, GroupCreatedTips.OpUser); err != nil {
		log.NewError(operationID, "setOpUserInfo failed ", err.Error(), opUserID, groupID, GroupCreatedTips.OpUser)
		return
	}
	err := c.setGroupInfo(groupID, GroupCreatedTips.Group)
	if err != nil {
		log.Error(operationID, "setGroupInfo failed ", groupID, GroupCreatedTips.Group)
		return
	}

	if err := c.setGroupOwnerInfo(groupID, GroupCreatedTips.GroupOwnerUser); err != nil {
		log.Error(operationID, "setGroupOwnerInfo failed", err.Error(), groupID)
		return
	}
	for _, v := range initMemberList {
		var groupMemberInfo sdkws.GroupMemberFullInfo
		if err := c.setGroupMemberInfo(groupID, v, &groupMemberInfo); err != nil {
			log.Error(operationID, "setGroupMemberInfo failed ", err.Error(), groupID, v)
			continue
		}
		GroupCreatedTips.MemberList = append(GroupCreatedTips.MemberList, &groupMemberInfo)
		if len(GroupCreatedTips.MemberList) == constant.MaxNotificationNum {
			break
		}
	}
	c.groupNotification(constant.GroupCreatedNotification, &GroupCreatedTips, opUserID, groupID, "", operationID)
}

// 群信息改变后掉用
// groupName := ""
//
//	notification := ""
//	introduction := ""
//	faceURL := ""
func (c *Check) GroupInfoSetNotification(operationID, opUserID, groupID string, groupName, notification, introduction, faceURL string, needVerification *wrapperspb.Int32Value) {
	GroupInfoChangedTips := sdkws.GroupInfoSetTips{Group: &sdkws.GroupInfo{}, OpUser: &sdkws.GroupMemberFullInfo{}}
	if err := c.setGroupInfo(groupID, GroupInfoChangedTips.Group); err != nil {
		log.Error(operationID, "setGroupInfo failed ", err.Error(), groupID)
		return
	}
	GroupInfoChangedTips.Group.GroupName = groupName
	GroupInfoChangedTips.Group.Notification = notification
	GroupInfoChangedTips.Group.Introduction = introduction
	GroupInfoChangedTips.Group.FaceURL = faceURL
	if needVerification != nil {
		GroupInfoChangedTips.Group.NeedVerification = needVerification.Value
	}

	if err := c.setOpUserInfo(opUserID, groupID, GroupInfoChangedTips.OpUser); err != nil {
		log.Error(operationID, "setOpUserInfo failed ", err.Error(), opUserID, groupID)
		return
	}
	c.groupNotification(constant.GroupInfoSetNotification, &GroupInfoChangedTips, opUserID, groupID, "", operationID)
}

func (c *Check) GroupMutedNotification(operationID, opUserID, groupID string) {
	tips := sdkws.GroupMutedTips{Group: &sdkws.GroupInfo{},
		OpUser: &sdkws.GroupMemberFullInfo{}}
	if err := c.setGroupInfo(groupID, tips.Group); err != nil {
		log.Error(operationID, "setGroupInfo failed ", err.Error(), groupID)
		return
	}
	if err := c.setOpUserInfo(opUserID, groupID, tips.OpUser); err != nil {
		log.Error(operationID, "setOpUserInfo failed ", err.Error(), opUserID, groupID)
		return
	}
	c.groupNotification(constant.GroupMutedNotification, &tips, opUserID, groupID, "", operationID)
}

func (c *Check) GroupCancelMutedNotification(operationID, opUserID, groupID string) {
	tips := sdkws.GroupCancelMutedTips{Group: &sdkws.GroupInfo{},
		OpUser: &sdkws.GroupMemberFullInfo{}}
	if err := c.setGroupInfo(groupID, tips.Group); err != nil {
		log.Error(operationID, "setGroupInfo failed ", err.Error(), groupID)
		return
	}
	if err := c.setOpUserInfo(opUserID, groupID, tips.OpUser); err != nil {
		log.Error(operationID, "setOpUserInfo failed ", err.Error(), opUserID, groupID)
		return
	}
	c.groupNotification(constant.GroupCancelMutedNotification, &tips, opUserID, groupID, "", operationID)
}

func (c *Check) GroupMemberMutedNotification(operationID, opUserID, groupID, groupMemberUserID string, mutedSeconds uint32) {
	tips := sdkws.GroupMemberMutedTips{Group: &sdkws.GroupInfo{},
		OpUser: &sdkws.GroupMemberFullInfo{}, MutedUser: &sdkws.GroupMemberFullInfo{}}
	tips.MutedSeconds = mutedSeconds
	if err := c.setGroupInfo(groupID, tips.Group); err != nil {
		log.Error(operationID, "setGroupInfo failed ", err.Error(), groupID)
		return
	}
	if err := c.setOpUserInfo(opUserID, groupID, tips.OpUser); err != nil {
		log.Error(operationID, "setOpUserInfo failed ", err.Error(), opUserID, groupID)
		return
	}
	if err := c.setGroupMemberInfo(groupID, groupMemberUserID, tips.MutedUser); err != nil {
		log.Error(operationID, "setGroupMemberInfo failed ", err.Error(), groupID, groupMemberUserID)
		return
	}
	c.groupNotification(constant.GroupMemberMutedNotification, &tips, opUserID, groupID, "", operationID)
}

func (c *Check) GroupMemberInfoSetNotification(operationID, opUserID, groupID, groupMemberUserID string) {
	tips := sdkws.GroupMemberInfoSetTips{Group: &sdkws.GroupInfo{},
		OpUser: &sdkws.GroupMemberFullInfo{}, ChangedUser: &sdkws.GroupMemberFullInfo{}}
	if err := c.setGroupInfo(groupID, tips.Group); err != nil {
		log.Error(operationID, "setGroupInfo failed ", err.Error(), groupID)
		return
	}
	if err := c.setOpUserInfo(opUserID, groupID, tips.OpUser); err != nil {
		log.Error(operationID, "setOpUserInfo failed ", err.Error(), opUserID, groupID)
		return
	}
	if err := c.setGroupMemberInfo(groupID, groupMemberUserID, tips.ChangedUser); err != nil {
		log.Error(operationID, "setGroupMemberInfo failed ", err.Error(), groupID, groupMemberUserID)
		return
	}
	c.groupNotification(constant.GroupMemberInfoSetNotification, &tips, opUserID, groupID, "", operationID)
}

func (c *Check) GroupMemberRoleLevelChangeNotification(operationID, opUserID, groupID, groupMemberUserID string, notificationType int32) {
	if notificationType != constant.GroupMemberSetToAdminNotification && notificationType != constant.GroupMemberSetToOrdinaryUserNotification {
		log.NewError(operationID, utils.GetSelfFuncName(), "invalid notificationType: ", notificationType)
		return
	}
	tips := sdkws.GroupMemberInfoSetTips{Group: &sdkws.GroupInfo{},
		OpUser: &sdkws.GroupMemberFullInfo{}, ChangedUser: &sdkws.GroupMemberFullInfo{}}
	if err := c.setGroupInfo(groupID, tips.Group); err != nil {
		log.Error(operationID, "setGroupInfo failed ", err.Error(), groupID)
		return
	}
	if err := c.setOpUserInfo(opUserID, groupID, tips.OpUser); err != nil {
		log.Error(operationID, "setOpUserInfo failed ", err.Error(), opUserID, groupID)
		return
	}
	if err := c.setGroupMemberInfo(groupID, groupMemberUserID, tips.ChangedUser); err != nil {
		log.Error(operationID, "setGroupMemberInfo failed ", err.Error(), groupID, groupMemberUserID)
		return
	}
	c.groupNotification(notificationType, &tips, opUserID, groupID, "", operationID)
}

func (c *Check) GroupMemberCancelMutedNotification(operationID, opUserID, groupID, groupMemberUserID string) {
	tips := sdkws.GroupMemberCancelMutedTips{Group: &sdkws.GroupInfo{},
		OpUser: &sdkws.GroupMemberFullInfo{}, MutedUser: &sdkws.GroupMemberFullInfo{}}
	if err := c.setGroupInfo(groupID, tips.Group); err != nil {
		log.Error(operationID, "setGroupInfo failed ", err.Error(), groupID)
		return
	}
	if err := c.setOpUserInfo(opUserID, groupID, tips.OpUser); err != nil {
		log.Error(operationID, "setOpUserInfo failed ", err.Error(), opUserID, groupID)
		return
	}
	if err := c.setGroupMemberInfo(groupID, groupMemberUserID, tips.MutedUser); err != nil {
		log.Error(operationID, "setGroupMemberInfo failed ", err.Error(), groupID, groupMemberUserID)
		return
	}
	c.groupNotification(constant.GroupMemberCancelMutedNotification, &tips, opUserID, groupID, "", operationID)
}

//	message ReceiveJoinApplicationTips{
//	 GroupInfo Group = 1;
//	 PublicUserInfo Applicant  = 2;
//	 string 	Reason = 3;
//	}  apply->all managers GroupID              string   `protobuf:"bytes,1,opt,name=GroupID" json:"GroupID,omitempty"`
//
//	ReqMessage           string   `protobuf:"bytes,2,opt,name=ReqMessage" json:"ReqMessage,omitempty"`
//	OpUserID             string   `protobuf:"bytes,3,opt,name=OpUserID" json:"OpUserID,omitempty"`
//	OperationID          string   `protobuf:"bytes,4,opt,name=OperationID" json:"OperationID,omitempty"`
//
// 申请进群后调用
func (c *Check) JoinGroupApplicationNotification(ctx context.Context, req *pbGroup.JoinGroupReq) {
	JoinGroupApplicationTips := sdkws.JoinGroupApplicationTips{Group: &sdkws.GroupInfo{}, Applicant: &sdkws.PublicUserInfo{}}
	err := c.setGroupInfo(req.GroupID, JoinGroupApplicationTips.Group)
	if err != nil {

		return
	}
	if err = c.setPublicUserInfo(tracelog.GetOpUserID(ctx), JoinGroupApplicationTips.Applicant); err != nil {

		return
	}
	JoinGroupApplicationTips.ReqMsg = req.ReqMessage

	managerList, err := imdb.GetOwnerManagerByGroupID(req.GroupID)
	if err != nil {

		return
	}
	for _, v := range managerList {
		c.groupNotification(constant.JoinGroupApplicationNotification, &JoinGroupApplicationTips, tracelog.GetOpUserID(ctx), "", v.UserID, utils.OperationID(ctx))

	}
}

func (c *Check) MemberQuitNotification(req *pbGroup.QuitGroupReq) {
	MemberQuitTips := sdkws.MemberQuitTips{Group: &sdkws.GroupInfo{}, QuitUser: &sdkws.GroupMemberFullInfo{}}
	if err := c.setGroupInfo(req.GroupID, MemberQuitTips.Group); err != nil {

		return
	}
	if err := c.setOpUserInfo(tracelog.GetOpUserID(), req.GroupID, MemberQuitTips.QuitUser); err != nil {
		return
	}

	c.groupNotification(constant.MemberQuitNotification, &MemberQuitTips, tracelog.GetOpUserID(), req.GroupID, "", req.OperationID)
}

//	message ApplicationProcessedTips{
//	 GroupInfo Group = 1;
//	 GroupMemberFullInfo OpUser = 2;
//	 int32 Result = 3;
//	 string 	Reason = 4;
//	}
//
// 处理进群请求后调用
func (c *Check) GroupApplicationAcceptedNotification(req *pbGroup.GroupApplicationResponseReq) {
	GroupApplicationAcceptedTips := sdkws.GroupApplicationAcceptedTips{Group: &sdkws.GroupInfo{}, OpUser: &sdkws.GroupMemberFullInfo{}, HandleMsg: req.HandledMsg}
	if err := c.setGroupInfo(req.GroupID, GroupApplicationAcceptedTips.Group); err != nil {
		log.NewError(req.OperationID, "setGroupInfo failed ", err.Error(), req.GroupID, GroupApplicationAcceptedTips.Group)
		return
	}
	if err := c.setOpUserInfo(req.OpUserID, req.GroupID, GroupApplicationAcceptedTips.OpUser); err != nil {
		log.Error(req.OperationID, "setOpUserInfo failed", req.OpUserID, req.GroupID, GroupApplicationAcceptedTips.OpUser)
		return
	}

	c.groupNotification(constant.GroupApplicationAcceptedNotification, &GroupApplicationAcceptedTips, req.OpUserID, "", req.FromUserID, req.OperationID)
	adminList, err := imdb.GetOwnerManagerByGroupID(req.GroupID)
	if err != nil {
		log.Error(req.OperationID, "GetOwnerManagerByGroupID failed", req.GroupID)
		return
	}
	for _, v := range adminList {
		if v.UserID == req.OpUserID {
			continue
		}
		GroupApplicationAcceptedTips.ReceiverAs = 1
		c.groupNotification(constant.GroupApplicationAcceptedNotification, &GroupApplicationAcceptedTips, req.OpUserID, "", v.UserID, req.OperationID)
	}
}

func (c *Check) GroupApplicationRejectedNotification(req *pbGroup.GroupApplicationResponseReq) {
	GroupApplicationRejectedTips := sdkws.GroupApplicationRejectedTips{Group: &sdkws.GroupInfo{}, OpUser: &sdkws.GroupMemberFullInfo{}, HandleMsg: req.HandledMsg}
	if err := c.setGroupInfo(req.GroupID, GroupApplicationRejectedTips.Group); err != nil {
		log.NewError(req.OperationID, "setGroupInfo failed ", err.Error(), req.GroupID, GroupApplicationRejectedTips.Group)
		return
	}
	if err := c.setOpUserInfo(req.OpUserID, req.GroupID, GroupApplicationRejectedTips.OpUser); err != nil {
		log.Error(req.OperationID, "setOpUserInfo failed", req.OpUserID, req.GroupID, GroupApplicationRejectedTips.OpUser)
		return
	}
	c.groupNotification(constant.GroupApplicationRejectedNotification, &GroupApplicationRejectedTips, req.OpUserID, "", req.FromUserID, req.OperationID)
	adminList, err := imdb.GetOwnerManagerByGroupID(req.GroupID)
	if err != nil {
		log.Error(req.OperationID, "GetOwnerManagerByGroupID failed", req.GroupID)
		return
	}
	for _, v := range adminList {
		if v.UserID == req.OpUserID {
			continue
		}
		GroupApplicationRejectedTips.ReceiverAs = 1
		c.groupNotification(constant.GroupApplicationRejectedNotification, &GroupApplicationRejectedTips, req.OpUserID, "", v.UserID, req.OperationID)
	}
}

func (c *Check) GroupOwnerTransferredNotification(req *pbGroup.TransferGroupOwnerReq) {
	GroupOwnerTransferredTips := sdkws.GroupOwnerTransferredTips{Group: &sdkws.GroupInfo{}, OpUser: &sdkws.GroupMemberFullInfo{}, NewGroupOwner: &sdkws.GroupMemberFullInfo{}}
	if err := c.setGroupInfo(req.GroupID, GroupOwnerTransferredTips.Group); err != nil {
		log.NewError(req.OperationID, "setGroupInfo failed ", err.Error(), req.GroupID)
		return
	}
	if err := c.setOpUserInfo(req.OpUserID, req.GroupID, GroupOwnerTransferredTips.OpUser); err != nil {
		log.Error(req.OperationID, "setOpUserInfo failed", req.OpUserID, req.GroupID)
		return
	}
	if err := c.setGroupMemberInfo(req.GroupID, req.NewOwnerUserID, GroupOwnerTransferredTips.NewGroupOwner); err != nil {
		log.Error(req.OperationID, "setGroupMemberInfo failed", req.GroupID, req.NewOwnerUserID)
		return
	}
	c.groupNotification(constant.GroupOwnerTransferredNotification, &GroupOwnerTransferredTips, req.OpUserID, req.GroupID, "", req.OperationID)
}

func (c *Check) GroupDismissedNotification(req *pbGroup.DismissGroupReq) {
	tips := sdkws.GroupDismissedTips{Group: &sdkws.GroupInfo{}, OpUser: &sdkws.GroupMemberFullInfo{}}
	if err := c.setGroupInfo(req.GroupID, tips.Group); err != nil {
		log.NewError(req.OperationID, "setGroupInfo failed ", err.Error(), req.GroupID)
		return
	}
	if err := c.setOpUserInfo(req.OpUserID, req.GroupID, tips.OpUser); err != nil {
		log.Error(req.OperationID, "setOpUserInfo failed", req.OpUserID, req.GroupID)
		return
	}
	c.groupNotification(constant.GroupDismissedNotification, &tips, req.OpUserID, req.GroupID, "", req.OperationID)
}

//	message MemberKickedTips{
//	 GroupInfo Group = 1;
//	 GroupMemberFullInfo OpUser = 2;
//	 GroupMemberFullInfo KickedUser = 3;
//	 uint64 OperationTime = 4;
//	}
//
// 被踢后调用
func (c *Check) MemberKickedNotification(req *pbGroup.KickGroupMemberReq, kickedUserIDList []string) {
	MemberKickedTips := sdkws.MemberKickedTips{Group: &sdkws.GroupInfo{}, OpUser: &sdkws.GroupMemberFullInfo{}}
	if err := c.setGroupInfo(req.GroupID, MemberKickedTips.Group); err != nil {
		log.Error(req.OperationID, "setGroupInfo failed ", err.Error(), req.GroupID)
		return
	}
	if err := c.setOpUserInfo(req.OpUserID, req.GroupID, MemberKickedTips.OpUser); err != nil {
		log.Error(req.OperationID, "setOpUserInfo failed ", err.Error(), req.OpUserID)
		return
	}
	for _, v := range kickedUserIDList {
		var groupMemberInfo sdkws.GroupMemberFullInfo
		if err := c.setGroupMemberInfo(req.GroupID, v, &groupMemberInfo); err != nil {
			log.Error(req.OperationID, "setGroupMemberInfo failed ", err.Error(), req.GroupID, v)
			continue
		}
		MemberKickedTips.KickedUserList = append(MemberKickedTips.KickedUserList, &groupMemberInfo)
	}
	c.groupNotification(constant.MemberKickedNotification, &MemberKickedTips, req.OpUserID, req.GroupID, "", req.OperationID)
	//
	//for _, v := range kickedUserIDList {
	//	groupNotification(constant.MemberKickedNotification, &MemberKickedTips, req.OpUserID, "", v, req.OperationID)
	//}
}

//	message MemberInvitedTips{
//	 GroupInfo Group = 1;
//	 GroupMemberFullInfo OpUser = 2;
//	 GroupMemberFullInfo InvitedUser = 3;
//	 uint64 OperationTime = 4;
//	}
//
// 被邀请进群后调用
func (c *Check) MemberInvitedNotification(operationID, groupID, opUserID, reason string, invitedUserIDList []string) {
	MemberInvitedTips := sdkws.MemberInvitedTips{Group: &sdkws.GroupInfo{}, OpUser: &sdkws.GroupMemberFullInfo{}}
	if err := c.setGroupInfo(groupID, MemberInvitedTips.Group); err != nil {
		log.Error(operationID, "setGroupInfo failed ", err.Error(), groupID)
		return
	}
	if err := c.setOpUserInfo(opUserID, groupID, MemberInvitedTips.OpUser); err != nil {
		log.Error(operationID, "setOpUserInfo failed ", err.Error(), opUserID, groupID)
		return
	}
	for _, v := range invitedUserIDList {
		var groupMemberInfo sdkws.GroupMemberFullInfo
		if err := c.setGroupMemberInfo(groupID, v, &groupMemberInfo); err != nil {
			log.Error(operationID, "setGroupMemberInfo failed ", err.Error(), groupID)
			continue
		}
		MemberInvitedTips.InvitedUserList = append(MemberInvitedTips.InvitedUserList, &groupMemberInfo)
	}
	c.groupNotification(constant.MemberInvitedNotification, &MemberInvitedTips, opUserID, groupID, "", operationID)
}

// 群成员主动申请进群，管理员同意后调用，
func (c *Check) MemberEnterNotification(ctx context.Context, req *pbGroup.GroupApplicationResponseReq) {
	MemberEnterTips := sdkws.MemberEnterTips{Group: &sdkws.GroupInfo{}, EntrantUser: &sdkws.GroupMemberFullInfo{}}
	if err := c.setGroupInfo(req.GroupID, MemberEnterTips.Group); err != nil {
		log.Error(req.OperationID, "setGroupInfo failed ", err.Error(), req.GroupID, MemberEnterTips.Group)
		return
	}
	if err := c.setGroupMemberInfo(req.GroupID, req.FromUserID, MemberEnterTips.EntrantUser); err != nil {
		log.Error(req.OperationID, "setGroupMemberInfo failed ", err.Error(), req.OpUserID, req.GroupID, MemberEnterTips.EntrantUser)
		return
	}
	c.groupNotification(constant.MemberEnterNotification, &MemberEnterTips, req.OpUserID, req.GroupID, "", req.OperationID)
}

func (c *Check) MemberEnterDirectlyNotification(groupID string, entrantUserID string, operationID string) {
	MemberEnterTips := sdkws.MemberEnterTips{Group: &sdkws.GroupInfo{}, EntrantUser: &sdkws.GroupMemberFullInfo{}}
	if err := c.setGroupInfo(groupID, MemberEnterTips.Group); err != nil {
		log.Error(operationID, "setGroupInfo failed ", err.Error(), groupID, MemberEnterTips.Group)
		return
	}
	if err := c.setGroupMemberInfo(groupID, entrantUserID, MemberEnterTips.EntrantUser); err != nil {
		log.Error(operationID, "setGroupMemberInfo failed ", err.Error(), groupID, entrantUserID, MemberEnterTips.EntrantUser)
		return
	}
	c.groupNotification(constant.MemberEnterNotification, &MemberEnterTips, entrantUserID, groupID, "", operationID)
}
