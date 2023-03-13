package notification

import (
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/log"
	"OpenIM/pkg/common/tokenverify"
	"OpenIM/pkg/common/tracelog"
	pbGroup "OpenIM/pkg/proto/group"
	"OpenIM/pkg/proto/sdkws"
	"OpenIM/pkg/proto/wrapperspb"
	"OpenIM/pkg/utils"
	"context"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func (c *Check) setOpUserInfo(ctx context.Context, groupID string, groupMemberInfo *sdkws.GroupMemberFullInfo) error {
	opUserID := tracelog.GetOpUserID(ctx)
	if tokenverify.IsManagerUserID(opUserID) {
		user, err := c.user.GetUsersInfos(ctx, []string{opUserID}, true)
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
	u, err := c.group.GetGroupMemberInfo(ctx, groupID, opUserID)
	if err == nil {
		*groupMemberInfo = *u
		return nil
	}
	user, err := c.user.GetUsersInfos(ctx, []string{opUserID}, true)
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

func (c *Check) setGroupInfo(ctx context.Context, groupID string, groupInfo *sdkws.GroupInfo) error {
	group, err := c.group.GetGroupInfos(ctx, []string{groupID}, true)
	if err != nil {
		return err
	}
	*groupInfo = *group[0]
	return nil
}

func (c *Check) setGroupMemberInfo(ctx context.Context, groupID, userID string, groupMemberInfo *sdkws.GroupMemberFullInfo) error {
	groupMember, err := c.group.GetGroupMemberInfo(ctx, groupID, userID)
	if err == nil {
		*groupMemberInfo = *groupMember
		return nil
	}
	user, err := c.user.GetUsersInfos(ctx, []string{userID}, true)
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

func (c *Check) setGroupOwnerInfo(ctx context.Context, groupID string, groupMemberInfo *sdkws.GroupMemberFullInfo) error {
	group, err := c.group.GetGroupInfo(ctx, groupID)
	if err != nil {
		return err
	}
	groupMember, err := c.group.GetGroupMemberInfo(ctx, groupID, group.OwnerUserID)
	if err != nil {
		return err
	}
	*groupMemberInfo = *groupMember
	return nil
}

func (c *Check) setPublicUserInfo(ctx context.Context, userID string, publicUserInfo *sdkws.PublicUserInfo) error {
	user, err := c.user.GetPublicUserInfos(ctx, []string{userID}, true)
	if err != nil {
		return err
	}
	*publicUserInfo = *user[0]
	return nil
}

func (c *Check) groupNotification(ctx context.Context, contentType int32, m proto.Message, sendID, groupID, recvUserID string) {
	var err error
	var tips sdkws.TipsComm
	tips.Detail, err = proto.Marshal(m)
	if err != nil {
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

		from, err := c.user.GetUsersInfos(ctx, []string{sendID}, true)
		if err != nil {
			return
		}
		nickname = from[0].Nickname
	}
	if recvUserID != "" {
		to, err := c.user.GetUsersInfos(ctx, []string{recvUserID}, true)
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
		return
	}

	var n NotificationMsg
	n.SendID = sendID
	if groupID != "" {
		n.RecvID = groupID

		group, err := c.group.GetGroupInfo(ctx, groupID)
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
	n.Content, err = proto.Marshal(&tips)
	if err != nil {
		return
	}
	c.Notification(ctx, &n)
}

// 创建群后调用
func (c *Check) GroupCreatedNotification(ctx context.Context, groupID string, initMemberList []string) {
	GroupCreatedTips := sdkws.GroupCreatedTips{Group: &sdkws.GroupInfo{},
		OpUser: &sdkws.GroupMemberFullInfo{}, GroupOwnerUser: &sdkws.GroupMemberFullInfo{}}
	if err := c.setOpUserInfo(ctx, groupID, GroupCreatedTips.OpUser); err != nil {
		return
	}
	err := c.setGroupInfo(ctx, groupID, GroupCreatedTips.Group)
	if err != nil {
		return
	}

	if err := c.setGroupOwnerInfo(ctx, groupID, GroupCreatedTips.GroupOwnerUser); err != nil {
		return
	}
	for _, v := range initMemberList {
		var groupMemberInfo sdkws.GroupMemberFullInfo
		if err := c.setGroupMemberInfo(ctx, groupID, v, &groupMemberInfo); err != nil {
			continue
		}
		GroupCreatedTips.MemberList = append(GroupCreatedTips.MemberList, &groupMemberInfo)
		if len(GroupCreatedTips.MemberList) == constant.MaxNotificationNum {
			break
		}
	}

	c.groupNotification(ctx, constant.GroupCreatedNotification, &GroupCreatedTips, tracelog.GetOpUserID(ctx), groupID, "")
}

// 群信息改变后掉用
// groupName := ""
//
//	notification := ""
//	introduction := ""
//	faceURL := ""
func (c *Check) GroupInfoSetNotification(ctx context.Context, groupID string, groupName, notification, introduction, faceURL string, needVerification *wrapperspb.Int32Value) {
	GroupInfoChangedTips := sdkws.GroupInfoSetTips{Group: &sdkws.GroupInfo{}, OpUser: &sdkws.GroupMemberFullInfo{}}
	if err := c.setGroupInfo(ctx, groupID, GroupInfoChangedTips.Group); err != nil {
		return
	}
	GroupInfoChangedTips.Group.GroupName = groupName
	GroupInfoChangedTips.Group.Notification = notification
	GroupInfoChangedTips.Group.Introduction = introduction
	GroupInfoChangedTips.Group.FaceURL = faceURL
	if needVerification != nil {
		GroupInfoChangedTips.Group.NeedVerification = needVerification.Value
	}

	if err := c.setOpUserInfo(ctx, groupID, GroupInfoChangedTips.OpUser); err != nil {
		return
	}
	c.groupNotification(ctx, constant.GroupInfoSetNotification, &GroupInfoChangedTips, tracelog.GetOpUserID(ctx), groupID, "")
}

func (c *Check) GroupMutedNotification(ctx context.Context, groupID string) {
	tips := sdkws.GroupMutedTips{Group: &sdkws.GroupInfo{},
		OpUser: &sdkws.GroupMemberFullInfo{}}
	if err := c.setGroupInfo(ctx, groupID, tips.Group); err != nil {
		return
	}
	if err := c.setOpUserInfo(ctx, groupID, tips.OpUser); err != nil {
		return
	}
	c.groupNotification(ctx, constant.GroupMutedNotification, &tips, tracelog.GetOpUserID(ctx), groupID, "")
}

func (c *Check) GroupCancelMutedNotification(ctx context.Context, groupID string) {
	tips := sdkws.GroupCancelMutedTips{Group: &sdkws.GroupInfo{},
		OpUser: &sdkws.GroupMemberFullInfo{}}
	if err := c.setGroupInfo(ctx, groupID, tips.Group); err != nil {
		return
	}
	if err := c.setOpUserInfo(ctx, groupID, tips.OpUser); err != nil {
		return
	}
	c.groupNotification(ctx, constant.GroupCancelMutedNotification, &tips, tracelog.GetOpUserID(ctx), groupID, "")
}

func (c *Check) GroupMemberMutedNotification(ctx context.Context, groupID, groupMemberUserID string, mutedSeconds uint32) {
	tips := sdkws.GroupMemberMutedTips{Group: &sdkws.GroupInfo{},
		OpUser: &sdkws.GroupMemberFullInfo{}, MutedUser: &sdkws.GroupMemberFullInfo{}}
	tips.MutedSeconds = mutedSeconds
	if err := c.setGroupInfo(ctx, groupID, tips.Group); err != nil {
		return
	}
	if err := c.setOpUserInfo(ctx, groupID, tips.OpUser); err != nil {
		return
	}
	if err := c.setGroupMemberInfo(ctx, groupID, groupMemberUserID, tips.MutedUser); err != nil {
		return
	}
	c.groupNotification(ctx, constant.GroupMemberMutedNotification, &tips, tracelog.GetOpUserID(ctx), groupID, "")
}

func (c *Check) GroupMemberInfoSetNotification(ctx context.Context, groupID, groupMemberUserID string) {
	tips := sdkws.GroupMemberInfoSetTips{Group: &sdkws.GroupInfo{},
		OpUser: &sdkws.GroupMemberFullInfo{}, ChangedUser: &sdkws.GroupMemberFullInfo{}}
	if err := c.setGroupInfo(ctx, groupID, tips.Group); err != nil {
		return
	}
	if err := c.setOpUserInfo(ctx, groupID, tips.OpUser); err != nil {
		return
	}
	if err := c.setGroupMemberInfo(ctx, groupID, groupMemberUserID, tips.ChangedUser); err != nil {
		return
	}
	c.groupNotification(ctx, constant.GroupMemberInfoSetNotification, &tips, tracelog.GetOpUserID(ctx), groupID, "")
}

func (c *Check) GroupMemberRoleLevelChangeNotification(ctx context.Context, operationID, opUserID, groupID, groupMemberUserID string, notificationType int32) {
	if notificationType != constant.GroupMemberSetToAdminNotification && notificationType != constant.GroupMemberSetToOrdinaryUserNotification {
		log.NewError(operationID, utils.GetSelfFuncName(), "invalid notificationType: ", notificationType)
		return
	}
	tips := sdkws.GroupMemberInfoSetTips{Group: &sdkws.GroupInfo{},
		OpUser: &sdkws.GroupMemberFullInfo{}, ChangedUser: &sdkws.GroupMemberFullInfo{}}
	if err := c.setGroupInfo(ctx, groupID, tips.Group); err != nil {
		log.Error(operationID, "setGroupInfo failed ", err.Error(), groupID)
		return
	}
	if err := c.setOpUserInfo(ctx, groupID, tips.OpUser); err != nil {
		log.Error(operationID, "setOpUserInfo failed ", err.Error(), opUserID, groupID)
		return
	}
	if err := c.setGroupMemberInfo(ctx, groupID, groupMemberUserID, tips.ChangedUser); err != nil {
		log.Error(operationID, "setGroupMemberInfo failed ", err.Error(), groupID, groupMemberUserID)
		return
	}
	c.groupNotification(ctx, notificationType, &tips, tracelog.GetOpUserID(ctx), groupID, "")
}

func (c *Check) GroupMemberCancelMutedNotification(ctx context.Context, groupID, groupMemberUserID string) {
	tips := sdkws.GroupMemberCancelMutedTips{Group: &sdkws.GroupInfo{},
		OpUser: &sdkws.GroupMemberFullInfo{}, MutedUser: &sdkws.GroupMemberFullInfo{}}
	if err := c.setGroupInfo(ctx, groupID, tips.Group); err != nil {
		return
	}
	if err := c.setOpUserInfo(ctx, groupID, tips.OpUser); err != nil {
		return
	}
	if err := c.setGroupMemberInfo(ctx, groupID, groupMemberUserID, tips.MutedUser); err != nil {
		return
	}
	c.groupNotification(ctx, constant.GroupMemberCancelMutedNotification, &tips, tracelog.GetOpUserID(ctx), groupID, "")
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
	err := c.setGroupInfo(ctx, req.GroupID, JoinGroupApplicationTips.Group)
	if err != nil {

		return
	}
	if err = c.setPublicUserInfo(ctx, tracelog.GetOpUserID(ctx), JoinGroupApplicationTips.Applicant); err != nil {

		return
	}
	JoinGroupApplicationTips.ReqMsg = req.ReqMessage
	managerList, err := c.group.GetOwnerAndAdminInfos(ctx, req.GroupID)
	if err != nil {
		return
	}
	for _, v := range managerList {
		c.groupNotification(ctx, constant.JoinGroupApplicationNotification, &JoinGroupApplicationTips, tracelog.GetOpUserID(ctx), "", v.UserID)
	}
}

func (c *Check) MemberQuitNotification(ctx context.Context, req *pbGroup.QuitGroupReq) {
	MemberQuitTips := sdkws.MemberQuitTips{Group: &sdkws.GroupInfo{}, QuitUser: &sdkws.GroupMemberFullInfo{}}
	if err := c.setGroupInfo(ctx, req.GroupID, MemberQuitTips.Group); err != nil {
		return
	}
	if err := c.setOpUserInfo(ctx, req.GroupID, MemberQuitTips.QuitUser); err != nil {
		return
	}

	c.groupNotification(ctx, constant.MemberQuitNotification, &MemberQuitTips, tracelog.GetOpUserID(ctx), req.GroupID, "")
}

//	message ApplicationProcessedTips{
//	 GroupInfo Group = 1;
//	 GroupMemberFullInfo OpUser = 2;
//	 int32 Result = 3;
//	 string 	Reason = 4;
//	}
//
// 处理进群请求后调用
func (c *Check) GroupApplicationAcceptedNotification(ctx context.Context, req *pbGroup.GroupApplicationResponseReq) {
	GroupApplicationAcceptedTips := sdkws.GroupApplicationAcceptedTips{Group: &sdkws.GroupInfo{}, OpUser: &sdkws.GroupMemberFullInfo{}, HandleMsg: req.HandledMsg}
	if err := c.setGroupInfo(ctx, req.GroupID, GroupApplicationAcceptedTips.Group); err != nil {
		return
	}
	if err := c.setOpUserInfo(ctx, req.GroupID, GroupApplicationAcceptedTips.OpUser); err != nil {
		return
	}

	c.groupNotification(ctx, constant.GroupApplicationAcceptedNotification, &GroupApplicationAcceptedTips, tracelog.GetOpUserID(ctx), "", req.FromUserID)
	adminList, err := c.group.GetOwnerAndAdminInfos(ctx, req.GroupID)
	if err != nil {
		return
	}
	for _, v := range adminList {
		if v.UserID == tracelog.GetOpUserID(ctx) {
			continue
		}
		GroupApplicationAcceptedTips.ReceiverAs = 1
		c.groupNotification(ctx, constant.GroupApplicationAcceptedNotification, &GroupApplicationAcceptedTips, tracelog.GetOpUserID(ctx), "", v.UserID)
	}
}

func (c *Check) GroupApplicationRejectedNotification(ctx context.Context, req *pbGroup.GroupApplicationResponseReq) {
	GroupApplicationRejectedTips := sdkws.GroupApplicationRejectedTips{Group: &sdkws.GroupInfo{}, OpUser: &sdkws.GroupMemberFullInfo{}, HandleMsg: req.HandledMsg}
	if err := c.setGroupInfo(ctx, req.GroupID, GroupApplicationRejectedTips.Group); err != nil {
		return
	}
	if err := c.setOpUserInfo(ctx, req.GroupID, GroupApplicationRejectedTips.OpUser); err != nil {
		return
	}
	c.groupNotification(ctx, constant.GroupApplicationRejectedNotification, &GroupApplicationRejectedTips, tracelog.GetOpUserID(ctx), "", req.FromUserID)
	adminList, err := c.group.GetOwnerAndAdminInfos(ctx, req.GroupID)
	if err != nil {
		return
	}
	for _, v := range adminList {
		if v.UserID == tracelog.GetOpUserID(ctx) {
			continue
		}
		GroupApplicationRejectedTips.ReceiverAs = 1
		c.groupNotification(ctx, constant.GroupApplicationRejectedNotification, &GroupApplicationRejectedTips, tracelog.GetOpUserID(ctx), "", v.UserID)
	}
}

func (c *Check) GroupOwnerTransferredNotification(ctx context.Context, req *pbGroup.TransferGroupOwnerReq) {
	GroupOwnerTransferredTips := sdkws.GroupOwnerTransferredTips{Group: &sdkws.GroupInfo{}, OpUser: &sdkws.GroupMemberFullInfo{}, NewGroupOwner: &sdkws.GroupMemberFullInfo{}}
	if err := c.setGroupInfo(ctx, req.GroupID, GroupOwnerTransferredTips.Group); err != nil {
		return
	}
	if err := c.setOpUserInfo(ctx, req.GroupID, GroupOwnerTransferredTips.OpUser); err != nil {
		return
	}
	if err := c.setGroupMemberInfo(ctx, req.GroupID, req.NewOwnerUserID, GroupOwnerTransferredTips.NewGroupOwner); err != nil {
		return
	}
	c.groupNotification(ctx, constant.GroupOwnerTransferredNotification, &GroupOwnerTransferredTips, tracelog.GetOpUserID(ctx), req.GroupID, "")
}

func (c *Check) GroupDismissedNotification(ctx context.Context, req *pbGroup.DismissGroupReq) {
	tips := sdkws.GroupDismissedTips{Group: &sdkws.GroupInfo{}, OpUser: &sdkws.GroupMemberFullInfo{}}
	if err := c.setGroupInfo(ctx, req.GroupID, tips.Group); err != nil {

		return
	}
	if err := c.setOpUserInfo(ctx, req.GroupID, tips.OpUser); err != nil {

		return
	}
	c.groupNotification(ctx, constant.GroupDismissedNotification, &tips, tracelog.GetOpUserID(ctx), req.GroupID, "")
}

//	message MemberKickedTips{
//	 GroupInfo Group = 1;
//	 GroupMemberFullInfo OpUser = 2;
//	 GroupMemberFullInfo KickedUser = 3;
//	 uint64 OperationTime = 4;
//	}
//
// 被踢后调用
func (c *Check) MemberKickedNotification(ctx context.Context, req *pbGroup.KickGroupMemberReq, kickedUserIDList []string) {
	MemberKickedTips := sdkws.MemberKickedTips{Group: &sdkws.GroupInfo{}, OpUser: &sdkws.GroupMemberFullInfo{}}
	if err := c.setGroupInfo(ctx, req.GroupID, MemberKickedTips.Group); err != nil {
		return
	}
	if err := c.setOpUserInfo(ctx, req.GroupID, MemberKickedTips.OpUser); err != nil {
		return
	}
	for _, v := range kickedUserIDList {
		var groupMemberInfo sdkws.GroupMemberFullInfo
		if err := c.setGroupMemberInfo(ctx, req.GroupID, v, &groupMemberInfo); err != nil {
			continue
		}
		MemberKickedTips.KickedUserList = append(MemberKickedTips.KickedUserList, &groupMemberInfo)
	}
	c.groupNotification(ctx, constant.MemberKickedNotification, &MemberKickedTips, tracelog.GetOpUserID(ctx), req.GroupID, "")
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
func (c *Check) MemberInvitedNotification(ctx context.Context, groupID, reason string, invitedUserIDList []string) {
	MemberInvitedTips := sdkws.MemberInvitedTips{Group: &sdkws.GroupInfo{}, OpUser: &sdkws.GroupMemberFullInfo{}}
	if err := c.setGroupInfo(ctx, groupID, MemberInvitedTips.Group); err != nil {
		return
	}
	if err := c.setOpUserInfo(ctx, groupID, MemberInvitedTips.OpUser); err != nil {

		return
	}
	for _, v := range invitedUserIDList {
		var groupMemberInfo sdkws.GroupMemberFullInfo
		if err := c.setGroupMemberInfo(ctx, groupID, v, &groupMemberInfo); err != nil {
			continue
		}
		MemberInvitedTips.InvitedUserList = append(MemberInvitedTips.InvitedUserList, &groupMemberInfo)
	}
	c.groupNotification(ctx, constant.MemberInvitedNotification, &MemberInvitedTips, tracelog.GetOpUserID(ctx), groupID, "")
}

// 群成员主动申请进群，管理员同意后调用，
func (c *Check) MemberEnterNotification(ctx context.Context, req *pbGroup.GroupApplicationResponseReq) {
	MemberEnterTips := sdkws.MemberEnterTips{Group: &sdkws.GroupInfo{}, EntrantUser: &sdkws.GroupMemberFullInfo{}}
	if err := c.setGroupInfo(ctx, req.GroupID, MemberEnterTips.Group); err != nil {
		return
	}
	if err := c.setGroupMemberInfo(ctx, req.GroupID, req.FromUserID, MemberEnterTips.EntrantUser); err != nil {
		return
	}
	c.groupNotification(ctx, constant.MemberEnterNotification, &MemberEnterTips, tracelog.GetOpUserID(ctx), req.GroupID, "")
}

func (c *Check) MemberEnterDirectlyNotification(ctx context.Context, groupID string, entrantUserID string, operationID string) {
	MemberEnterTips := sdkws.MemberEnterTips{Group: &sdkws.GroupInfo{}, EntrantUser: &sdkws.GroupMemberFullInfo{}}
	if err := c.setGroupInfo(ctx, groupID, MemberEnterTips.Group); err != nil {
		log.Error(operationID, "setGroupInfo failed ", err.Error(), groupID, MemberEnterTips.Group)
		return
	}
	if err := c.setGroupMemberInfo(ctx, groupID, entrantUserID, MemberEnterTips.EntrantUser); err != nil {
		log.Error(operationID, "setGroupMemberInfo failed ", err.Error(), groupID, entrantUserID, MemberEnterTips.EntrantUser)
		return
	}
	c.groupNotification(ctx, constant.MemberEnterNotification, &MemberEnterTips, entrantUserID, groupID, "")
}
