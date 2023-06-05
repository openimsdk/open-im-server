package notification

import (
	"context"
	"fmt"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/controller"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	pbGroup "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/group"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
)

func NewGroupNotificationSender(db controller.GroupDatabase, sdr discoveryregistry.SvcDiscoveryRegistry, fn func(ctx context.Context, userIDs []string) ([]CommonUser, error)) *GroupNotificationSender {
	return &GroupNotificationSender{
		NotificationSender: rpcclient.NewNotificationSender(rpcclient.WithDiscov(sdr)),
		getUsersInfo:       fn,
		db:                 db,
	}
}

type GroupNotificationSender struct {
	*rpcclient.NotificationSender
	getUsersInfo func(ctx context.Context, userIDs []string) ([]CommonUser, error)
	db           controller.GroupDatabase
}

func (g *GroupNotificationSender) getUser(ctx context.Context, userID string) (*sdkws.PublicUserInfo, error) {
	users, err := g.getUsersInfo(ctx, []string{userID})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errs.ErrUserIDNotFound.Wrap(fmt.Sprintf("user %s not found", userID))
	}
	return &sdkws.PublicUserInfo{
		UserID:   users[0].GetUserID(),
		Nickname: users[0].GetNickname(),
		FaceURL:  users[0].GetFaceURL(),
		Ex:       users[0].GetEx(),
	}, nil
}

func (g *GroupNotificationSender) getGroupInfo(ctx context.Context, groupID string) (*sdkws.GroupInfo, error) {
	gm, err := g.db.TakeGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}
	userIDs, err := g.db.FindGroupMemberUserID(ctx, groupID)
	if err != nil {
		return nil, err
	}
	owner, err := g.db.FindGroupMember(ctx, []string{groupID}, nil, []int32{constant.GroupOwner})
	if err != nil {
		return nil, err
	}
	if len(owner) == 0 {
		return nil, errs.ErrInternalServer.Wrap(fmt.Sprintf("group %s owner not found", groupID))
	}
	return &sdkws.GroupInfo{
		GroupID:                gm.GroupID,
		GroupName:              gm.GroupName,
		Notification:           gm.Notification,
		Introduction:           gm.Introduction,
		FaceURL:                gm.FaceURL,
		OwnerUserID:            owner[0].UserID,
		CreateTime:             gm.CreateTime.UnixMilli(),
		MemberCount:            uint32(len(userIDs)),
		Ex:                     gm.Ex,
		Status:                 gm.Status,
		CreatorUserID:          gm.CreatorUserID,
		GroupType:              gm.GroupType,
		NeedVerification:       gm.NeedVerification,
		LookMemberInfo:         gm.LookMemberInfo,
		ApplyMemberFriend:      gm.ApplyMemberFriend,
		NotificationUpdateTime: gm.NotificationUpdateTime.UnixMilli(),
		NotificationUserID:     gm.NotificationUserID,
	}, nil
}

func (g *GroupNotificationSender) getGroupMembers(ctx context.Context, groupID string, userIDs []string) ([]*sdkws.GroupMemberFullInfo, error) {
	//members, err := g.db.FindGroupMember(ctx, []string{groupID}, userIDs, []int32{constant.GroupOwner})
	members, err := g.db.FindGroupMember(ctx, []string{groupID}, userIDs, nil)
	if err != nil {
		return nil, err
	}
	users, err := g.getUsersInfoMap(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	res := make([]*sdkws.GroupMemberFullInfo, 0, len(members))
	for _, member := range members {
		var appMangerLevel int32
		if user := users[member.UserID]; user != nil {
			appMangerLevel = user.AppMangerLevel
			if member.Nickname == "" {
				member.Nickname = user.Nickname
			}
			if member.FaceURL == "" {
				member.FaceURL = user.FaceURL
			}
		}
		res = append(res, g.groupMemberDB2PB(member, appMangerLevel))
		delete(users, member.UserID)
	}
	for userID, info := range users {
		if info.AppMangerLevel == constant.AppAdmin {
			res = append(res, &sdkws.GroupMemberFullInfo{
				GroupID:        groupID,
				UserID:         userID,
				Nickname:       info.Nickname,
				FaceURL:        info.FaceURL,
				AppMangerLevel: info.AppMangerLevel,
			})
		}
	}
	return res, nil
}

func (g *GroupNotificationSender) getGroupMemberMap(ctx context.Context, groupID string, userIDs []string) (map[string]*sdkws.GroupMemberFullInfo, error) {
	members, err := g.getGroupMembers(ctx, groupID, userIDs)
	if err != nil {
		return nil, err
	}
	m := make(map[string]*sdkws.GroupMemberFullInfo)
	for i, member := range members {
		m[member.UserID] = members[i]
	}
	return m, nil
}

func (g *GroupNotificationSender) getGroupMember(ctx context.Context, groupID string, userID string) (*sdkws.GroupMemberFullInfo, error) {
	members, err := g.getGroupMembers(ctx, groupID, []string{userID})
	if err != nil {
		return nil, err
	}
	if len(members) == 0 {
		return nil, errs.ErrInternalServer.Wrap(fmt.Sprintf("group %s member %s not found", groupID, userID))
	}
	return members[0], nil
}

func (g *GroupNotificationSender) getGroupOwnerAndAdminUserID(ctx context.Context, groupID string) ([]string, error) {
	members, err := g.db.FindGroupMember(ctx, []string{groupID}, nil, []int32{constant.GroupOwner, constant.GroupAdmin})
	if err != nil {
		return nil, err
	}
	fn := func(e *relation.GroupMemberModel) string { return e.UserID }
	return utils.Slice(members, fn), nil
}

func (g *GroupNotificationSender) groupDB2PB(group *relation.GroupModel, ownerUserID string, memberCount uint32) *sdkws.GroupInfo {
	return &sdkws.GroupInfo{
		GroupID:                group.GroupID,
		GroupName:              group.GroupName,
		Notification:           group.Notification,
		Introduction:           group.Introduction,
		FaceURL:                group.FaceURL,
		OwnerUserID:            ownerUserID,
		CreateTime:             group.CreateTime.UnixMilli(),
		MemberCount:            memberCount,
		Ex:                     group.Ex,
		Status:                 group.Status,
		CreatorUserID:          group.CreatorUserID,
		GroupType:              group.GroupType,
		NeedVerification:       group.NeedVerification,
		LookMemberInfo:         group.LookMemberInfo,
		ApplyMemberFriend:      group.ApplyMemberFriend,
		NotificationUpdateTime: group.NotificationUpdateTime.UnixMilli(),
		NotificationUserID:     group.NotificationUserID,
	}
}

func (g *GroupNotificationSender) groupMemberDB2PB(member *relation.GroupMemberModel, appMangerLevel int32) *sdkws.GroupMemberFullInfo {
	return &sdkws.GroupMemberFullInfo{
		GroupID:        member.GroupID,
		UserID:         member.UserID,
		RoleLevel:      member.RoleLevel,
		JoinTime:       member.JoinTime.UnixMilli(),
		Nickname:       member.Nickname,
		FaceURL:        member.FaceURL,
		AppMangerLevel: appMangerLevel,
		JoinSource:     member.JoinSource,
		OperatorUserID: member.OperatorUserID,
		Ex:             member.Ex,
		MuteEndTime:    member.MuteEndTime.UnixMilli(),
		InviterUserID:  member.InviterUserID,
	}
}

func (g *GroupNotificationSender) getUsersInfoMap(ctx context.Context, userIDs []string) (map[string]*sdkws.UserInfo, error) {
	users, err := g.getUsersInfo(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	result := make(map[string]*sdkws.UserInfo)
	for _, user := range users {
		result[user.GetUserID()] = user.(*sdkws.UserInfo)
	}
	return result, nil
}

func (g *GroupNotificationSender) getFromToUserNickname(ctx context.Context, fromUserID, toUserID string) (string, string, error) {
	users, err := g.getUsersInfoMap(ctx, []string{fromUserID, toUserID})
	if err != nil {
		return "", "", nil
	}
	return users[fromUserID].Nickname, users[toUserID].Nickname, nil
}

func (g *GroupNotificationSender) mergeGroupFull(ctx context.Context, groupID string, group *relation.GroupModel, ms *[]*relation.GroupMemberModel, users *map[string]*sdkws.UserInfo) (groupInfo *sdkws.GroupCreatedTips, err error) {
	defer log.ZDebug(ctx, "return")
	defer func() {
		if err != nil {
			log.ZError(ctx, utils.GetFuncName(1)+" failed", err)
		}
	}()
	if group == nil {
		group, err = g.db.TakeGroup(ctx, groupID)
		if err != nil {
			return nil, err
		}
	}
	if ms == nil {
		var temp []*relation.GroupMemberModel
		ms = &temp
	}
	if users == nil {
		temp := make(map[string]*sdkws.UserInfo)
		users = &temp
	}
	var members []*relation.GroupMemberModel
	if len(*ms) == 0 {
		members, err = g.db.FindGroupMember(ctx, []string{groupID}, nil, nil)
		if err != nil {
			return nil, err
		}
		*ms = members
	} else {
		members = *ms
	}
	opUserID := mcontext.GetOpUserID(ctx)
	var userMap map[string]*sdkws.UserInfo
	if users == nil || len(*users) == 0 {
		userIDs := utils.Slice(members, func(e *relation.GroupMemberModel) string { return e.UserID })
		userIDs = append(userIDs, opUserID)
		userMap, err = g.getUsersInfoMap(ctx, userIDs)
		if err != nil {
			return nil, err
		}
		if users != nil {
			*users = userMap
		}
	} else {
		userMap = *users
	}
	var (
		opUserMember     *sdkws.GroupMemberFullInfo
		groupOwnerMember *sdkws.GroupMemberFullInfo
	)
	for _, member := range members {
		if member.UserID == opUserID {
			opUserMember = g.groupMemberDB2PB(member, userMap[member.UserID].AppMangerLevel)
		}
		if member.RoleLevel == constant.GroupOwner {
			groupOwnerMember = g.groupMemberDB2PB(member, userMap[member.UserID].AppMangerLevel)
		}
		if opUserMember != nil && groupOwnerMember != nil {
			break
		}
	}
	if opUser := userMap[opUserID]; opUser != nil && opUserMember == nil {
		opUserMember = &sdkws.GroupMemberFullInfo{
			GroupID:        group.GroupID,
			UserID:         opUser.UserID,
			Nickname:       opUser.Nickname,
			FaceURL:        opUser.FaceURL,
			AppMangerLevel: opUser.AppMangerLevel,
		}
	}
	groupInfo = &sdkws.GroupCreatedTips{Group: g.groupDB2PB(group, groupOwnerMember.UserID, uint32(len(members))),
		OpUser: opUserMember, GroupOwnerUser: groupOwnerMember}
	return groupInfo, nil
}

//func (g *GroupNotificationSender) GroupCreatedNotification(ctx context.Context, group *relation.GroupModel, members []*relation.GroupMemberModel, userMap map[string]*sdkws.UserInfo) (err error) {
//	defer log.ZDebug(ctx, "GroupCreatedNotification.return")
//	defer func() {
//		if err != nil {
//			log.ZError(ctx, utils.GetFuncName(1)+" failed", err)
//		}
//	}()
//	groupInfo, err := g.mergeGroupFull(ctx, group.GroupID, group, &members, &userMap)
//	if err != nil {
//		return err
//	}
//	return g.msgClient.Notification(ctx, mcontext.GetOpUserID(ctx), group.GroupID, constant.GroupCreatedNotification, groupInfo)
//}

func (g *GroupNotificationSender) GroupCreatedNotification(ctx context.Context, tips *sdkws.GroupCreatedTips) (err error) {
	return g.Notification(ctx, mcontext.GetOpUserID(ctx), tips.Group.GroupID, constant.GroupCreatedNotification, tips)
}

//func (g *GroupNotificationSender) GroupInfoSetNotification(ctx context.Context, group *relation.GroupModel, members []*relation.GroupMemberModel, needVerification *int32) (err error) {
//	groupInfo, err := g.mergeGroupFull(ctx, group.GroupID, group, &members, nil)
//	if err != nil {
//		return err
//	}
//	tips := &sdkws.GroupInfoSetTips{Group: groupInfo.Group, OpUser: groupInfo.GroupOwnerUser}
//	if needVerification != nil {
//		tips.Group.NeedVerification = *needVerification
//	}
//	return g.msgClient.Notification(ctx, mcontext.GetOpUserID(ctx), group.GroupID, constant.GroupInfoSetNotification, tips)
//}

func (g *GroupNotificationSender) GroupInfoSetNotification(ctx context.Context, tips *sdkws.GroupInfoSetTips) (err error) {
	//groupInfo, err := g.mergeGroupFull(ctx, group.GroupID, group, &members, nil)
	//if err != nil {
	//	return err
	//}
	//tips := &sdkws.GroupInfoSetTips{Group: groupInfo.Group, OpUser: groupInfo.GroupOwnerUser}
	//if needVerification != nil {
	//	tips.Group.NeedVerification = *needVerification
	//}
	return g.Notification(ctx, mcontext.GetOpUserID(ctx), tips.Group.GroupID, constant.GroupInfoSetNotification, tips)
}

func (g *GroupNotificationSender) JoinGroupApplicationNotification(ctx context.Context, req *pbGroup.JoinGroupReq) (err error) {
	defer log.ZDebug(ctx, "return")
	defer func() {
		if err != nil {
			log.ZError(ctx, utils.GetFuncName(1)+" failed", err)
		}
	}()
	group, err := g.getGroupInfo(ctx, req.GroupID)
	if err != nil {
		return err
	}
	user, err := g.getUser(ctx, req.InviterUserID)
	if err != nil {
		return err
	}
	userIDs, err := g.getGroupOwnerAndAdminUserID(ctx, req.GroupID)
	if err != nil {
		return err
	}
	userIDs = append(userIDs, req.InviterUserID, mcontext.GetOpUserID(ctx))
	tips := &sdkws.JoinGroupApplicationTips{Group: group, Applicant: user, ReqMsg: req.ReqMessage}
	for _, userID := range utils.Distinct(userIDs) {
		err = g.Notification(ctx, mcontext.GetOpUserID(ctx), userID, constant.JoinGroupApplicationNotification, tips)
		if err != nil {
			log.ZError(ctx, "JoinGroupApplicationNotification failed", err, "group", req.GroupID, "userID", userID)
		}
	}
	return nil
}

func (g *GroupNotificationSender) MemberQuitNotification(ctx context.Context, req *pbGroup.QuitGroupReq) (err error) {
	defer log.ZDebug(ctx, "return")
	defer func() {
		if err != nil {
			log.ZError(ctx, utils.GetFuncName(1)+" failed", err)
		}
	}()
	group, err := g.getGroupInfo(ctx, req.GroupID)
	if err != nil {
		return err
	}
	opUserID := mcontext.GetOpUserID(ctx)
	user, err := g.getUser(ctx, opUserID)
	if err != nil {
		return err
	}
	//userIDs, err := g.getGroupOwnerAndAdminUserID(ctx, req.GroupID)
	//if err != nil {
	//	return err
	//}
	tips := &sdkws.MemberQuitTips{Group: group, QuitUser: &sdkws.GroupMemberFullInfo{
		GroupID:  group.GroupID,
		UserID:   user.UserID,
		Nickname: user.Nickname,
		FaceURL:  user.FaceURL,
	}}
	//for _, userID := range append(userIDs, opUserID) {
	//	err = g.msgClient.Notification(ctx, mcontext.GetOpUserID(ctx), userID, constant.MemberQuitNotification, tips)
	//	if err != nil {
	//		log.ZError(ctx, "MemberQuitNotification failed", err, "group", req.GroupID, "userID", userID)
	//	}
	//}
	return g.Notification(ctx, mcontext.GetOpUserID(ctx), req.GroupID, constant.MemberQuitNotification, tips)
}

func (g *GroupNotificationSender) GroupApplicationAcceptedNotification(ctx context.Context, req *pbGroup.GroupApplicationResponseReq) (err error) {
	defer log.ZDebug(ctx, "return")
	defer func() {
		if err != nil {
			log.ZError(ctx, utils.GetFuncName(1)+" failed", err)
		}
	}()
	group, err := g.getGroupInfo(ctx, req.GroupID)
	if err != nil {
		return err
	}
	user, err := g.getGroupMember(ctx, req.GroupID, mcontext.GetOpUserID(ctx))
	if err != nil {
		return err
	}
	userIDs, err := g.getGroupOwnerAndAdminUserID(ctx, req.GroupID)
	if err != nil {
		return err
	}
	tips := &sdkws.GroupApplicationAcceptedTips{Group: group, OpUser: user, HandleMsg: req.HandledMsg, ReceiverAs: 1}
	for _, userID := range append(userIDs, mcontext.GetOpUserID(ctx)) {
		err = g.Notification(ctx, mcontext.GetOpUserID(ctx), userID, constant.GroupApplicationAcceptedNotification, tips)
		if err != nil {
			log.ZError(ctx, "failed", err)
		}
	}
	return nil
}

func (g *GroupNotificationSender) GroupApplicationRejectedNotification(ctx context.Context, req *pbGroup.GroupApplicationResponseReq) (err error) {
	defer log.ZDebug(ctx, "return")
	defer func() {
		if err != nil {
			log.ZError(ctx, utils.GetFuncName(1)+" failed", err)
		}
	}()
	group, err := g.getGroupInfo(ctx, req.GroupID)
	if err != nil {
		return err
	}
	user, err := g.getGroupMember(ctx, req.GroupID, mcontext.GetOpUserID(ctx))
	if err != nil {
		return err
	}
	userIDs, err := g.getGroupOwnerAndAdminUserID(ctx, req.GroupID)
	if err != nil {
		return err
	}
	tips := &sdkws.GroupApplicationRejectedTips{Group: group, OpUser: user, HandleMsg: req.HandledMsg}
	for _, userID := range append(userIDs, mcontext.GetOpUserID(ctx)) {
		err = g.Notification(ctx, mcontext.GetOpUserID(ctx), userID, constant.GroupApplicationRejectedNotification, tips)
		if err != nil {
			log.ZError(ctx, "failed", err)
		}
	}
	return nil
}

func (g *GroupNotificationSender) GroupOwnerTransferredNotification(ctx context.Context, req *pbGroup.TransferGroupOwnerReq) (err error) {
	defer log.ZDebug(ctx, "return")
	defer func() {
		if err != nil {
			log.ZError(ctx, utils.GetFuncName(1)+" failed", err)
		}
	}()
	group, err := g.getGroupInfo(ctx, req.GroupID)
	if err != nil {
		return err
	}
	opUserID := mcontext.GetOpUserID(ctx)
	member, err := g.getGroupMemberMap(ctx, req.GroupID, []string{opUserID, req.NewOwnerUserID})
	if err != nil {
		return err
	}
	tips := &sdkws.GroupOwnerTransferredTips{Group: group, OpUser: member[opUserID], NewGroupOwner: member[req.NewOwnerUserID]}
	return g.Notification(ctx, mcontext.GetOpUserID(ctx), group.GroupID, constant.GroupOwnerTransferredNotification, tips)
}

func (g *GroupNotificationSender) MemberKickedNotification(ctx context.Context, tips *sdkws.MemberKickedTips) (err error) {
	//defer log.ZDebug(ctx, "return")
	//defer func() {
	//	if err != nil {
	//		log.ZError(ctx, utils.GetFuncName(1)+" failed", err)
	//	}
	//}()
	//group, err := g.getGroupInfo(ctx, tips.GroupID)
	//if err != nil {
	//	return err
	//}
	//user, err := g.getGroupMember(ctx, tips.GroupID, mcontext.GetOpUserID(ctx))
	//if err != nil {
	//	return err
	//}
	//tips := &sdkws.MemberKickedTips{Group: group, OpUser: user}
	//for _, v := range kickedUserIDList {
	//	var groupMemberInfo sdkws.GroupMemberFullInfo
	//	if err := c.setGroupMemberInfo(ctx, tips.GroupID, v, &groupMemberInfo); err != nil {
	//		continue
	//	}
	//	MemberKickedTips.KickedUserList = append(MemberKickedTips.KickedUserList, &groupMemberInfo)
	//}
	return g.Notification(ctx, mcontext.GetOpUserID(ctx), tips.Group.GroupID, constant.MemberKickedNotification, tips)
}

func (g *GroupNotificationSender) MemberInvitedNotification(ctx context.Context, groupID, reason string, invitedUserIDList []string) (err error) {
	defer log.ZDebug(ctx, "return")
	defer func() {
		if err != nil {
			log.ZError(ctx, utils.GetFuncName(1)+" failed", err)
		}
	}()
	group, err := g.getGroupInfo(ctx, groupID)
	if err != nil {
		return err
	}
	opUser, err := g.getGroupMember(ctx, groupID, mcontext.GetOpUserID(ctx))
	if err != nil {
		return err
	}
	users, err := g.getGroupMembers(ctx, groupID, invitedUserIDList)
	if err != nil {
		return err
	}
	tips := &sdkws.MemberInvitedTips{Group: group, OpUser: opUser, InvitedUserList: users}
	return g.Notification(ctx, mcontext.GetOpUserID(ctx), group.GroupID, constant.MemberInvitedNotification, tips)
}

func (g *GroupNotificationSender) MemberEnterNotification(ctx context.Context, req *pbGroup.GroupApplicationResponseReq) (err error) {
	defer log.ZDebug(ctx, "return")
	defer func() {
		if err != nil {
			log.ZError(ctx, utils.GetFuncName(1)+" failed", err)
		}
	}()
	group, err := g.getGroupInfo(ctx, req.GroupID)
	if err != nil {
		return err
	}
	user, err := g.getGroupMember(ctx, req.GroupID, req.FromUserID)
	if err != nil {
		return err
	}
	tips := &sdkws.MemberEnterTips{Group: group, EntrantUser: user}
	return g.Notification(ctx, mcontext.GetOpUserID(ctx), group.GroupID, constant.MemberEnterNotification, tips)
}

func (g *GroupNotificationSender) GroupDismissedNotification(ctx context.Context, tips *sdkws.GroupDismissedTips) (err error) {
	//defer log.ZDebug(ctx, "return")
	//defer func() {
	//	if err != nil {
	//		log.ZError(ctx, utils.GetFuncName(1)+" failed", err)
	//	}
	//}()
	//group, err := g.getGroupInfo(ctx, req.GroupID)
	//if err != nil {
	//	return err
	//}
	//users, err := g.getGroupMembers(ctx, req.GroupID, []string{mcontext.GetOpUserID(ctx)})
	//if err != nil {
	//	return err
	//}
	//tips := &sdkws.GroupDismissedTips{Group: group}
	//if len(users) > 0 {
	//	tips.OpUser = users[0]
	//} else {
	//	tips.OpUser = &sdkws.GroupMemberFullInfo{
	//		GroupID: group.GroupID,
	//		UserID:  mcontext.GetOpUserID(ctx),
	//	}
	//}
	//return g.Notification(ctx, mcontext.GetOpUserID(ctx), group.GroupID, constant.GroupDismissedNotification, tips)
	return g.Notification(ctx, mcontext.GetOpUserID(ctx), tips.Group.GroupID, constant.GroupDismissedNotification, tips)
}

func (g *GroupNotificationSender) GroupMemberMutedNotification(ctx context.Context, groupID, groupMemberUserID string, mutedSeconds uint32) (err error) {
	defer log.ZDebug(ctx, "return")
	defer func() {
		if err != nil {
			log.ZError(ctx, utils.GetFuncName(1)+" failed", err)
		}
	}()
	group, err := g.getGroupInfo(ctx, groupID)
	if err != nil {
		return err
	}
	user, err := g.getGroupMemberMap(ctx, groupID, []string{mcontext.GetOpUserID(ctx), groupMemberUserID})
	if err != nil {
		return err
	}
	tips := &sdkws.GroupMemberMutedTips{Group: group, MutedSeconds: mutedSeconds,
		OpUser: user[mcontext.GetOpUserID(ctx)], MutedUser: user[groupMemberUserID]}
	return g.Notification(ctx, mcontext.GetOpUserID(ctx), group.GroupID, constant.GroupMemberMutedNotification, tips)
}

func (g *GroupNotificationSender) GroupMemberCancelMutedNotification(ctx context.Context, groupID, groupMemberUserID string) (err error) {
	defer log.ZDebug(ctx, "return")
	defer func() {
		if err != nil {
			log.ZError(ctx, utils.GetFuncName(1)+" failed", err)
		}
	}()
	group, err := g.getGroupInfo(ctx, groupID)
	if err != nil {
		return err
	}
	user, err := g.getGroupMemberMap(ctx, groupID, []string{mcontext.GetOpUserID(ctx), groupMemberUserID})
	if err != nil {
		return err
	}
	tips := &sdkws.GroupMemberCancelMutedTips{Group: group, OpUser: user[mcontext.GetOpUserID(ctx)], MutedUser: user[groupMemberUserID]}
	return g.Notification(ctx, mcontext.GetOpUserID(ctx), group.GroupID, constant.GroupMemberCancelMutedNotification, tips)
}

func (g *GroupNotificationSender) GroupMutedNotification(ctx context.Context, groupID string) (err error) {
	defer log.ZDebug(ctx, "return")
	defer func() {
		if err != nil {
			log.ZError(ctx, utils.GetFuncName(1)+" failed", err)
		}
	}()
	group, err := g.getGroupInfo(ctx, groupID)
	if err != nil {
		return err
	}
	users, err := g.getGroupMembers(ctx, groupID, []string{mcontext.GetOpUserID(ctx)})
	if err != nil {
		return err
	}
	tips := &sdkws.GroupMutedTips{Group: group}
	if len(users) > 0 {
		tips.OpUser = users[0]
	} else {
		tips.OpUser = &sdkws.GroupMemberFullInfo{UserID: mcontext.GetOpUserID(ctx), GroupID: groupID}
	}
	return g.Notification(ctx, mcontext.GetOpUserID(ctx), group.GroupID, constant.GroupMutedNotification, tips)
}

func (g *GroupNotificationSender) GroupCancelMutedNotification(ctx context.Context, groupID string) (err error) {
	defer log.ZDebug(ctx, "return")
	defer func() {
		if err != nil {
			log.ZError(ctx, utils.GetFuncName(1)+" failed", err)
		}
	}()
	group, err := g.getGroupInfo(ctx, groupID)
	if err != nil {
		return err
	}
	users, err := g.getGroupMembers(ctx, groupID, []string{mcontext.GetOpUserID(ctx)})
	if err != nil {
		return err
	}
	tips := &sdkws.GroupCancelMutedTips{Group: group}
	if len(users) > 0 {
		tips.OpUser = users[0]
	} else {
		tips.OpUser = &sdkws.GroupMemberFullInfo{UserID: mcontext.GetOpUserID(ctx), GroupID: groupID}
	}
	return g.Notification(ctx, mcontext.GetOpUserID(ctx), group.GroupID, constant.GroupCancelMutedNotification, tips)
}

func (g *GroupNotificationSender) GroupMemberInfoSetNotification(ctx context.Context, groupID, groupMemberUserID string) (err error) {
	defer log.ZDebug(ctx, "return")
	defer func() {
		if err != nil {
			log.ZError(ctx, utils.GetFuncName(1)+" failed", err)
		}
	}()
	group, err := g.getGroupInfo(ctx, groupID)
	if err != nil {
		return err
	}
	user, err := g.getGroupMemberMap(ctx, groupID, []string{mcontext.GetOpUserID(ctx), groupMemberUserID})
	if err != nil {
		return err
	}
	tips := &sdkws.GroupMemberInfoSetTips{Group: group, OpUser: user[mcontext.GetOpUserID(ctx)], ChangedUser: user[groupMemberUserID]}
	return g.Notification(ctx, mcontext.GetOpUserID(ctx), group.GroupID, constant.GroupMemberInfoSetNotification, tips)
}

func (g *GroupNotificationSender) GroupMemberSetToAdminNotification(ctx context.Context, groupID, groupMemberUserID string) (err error) {
	defer log.ZDebug(ctx, "return")
	defer func() {
		if err != nil {
			log.ZError(ctx, utils.GetFuncName(1)+" failed", err)
		}
	}()
	group, err := g.getGroupInfo(ctx, groupID)
	if err != nil {
		return err
	}
	user, err := g.getGroupMemberMap(ctx, groupID, []string{mcontext.GetOpUserID(ctx), groupMemberUserID})
	if err != nil {
		return err
	}
	tips := &sdkws.GroupMemberInfoSetTips{Group: group, OpUser: user[mcontext.GetOpUserID(ctx)], ChangedUser: user[groupMemberUserID]}
	return g.Notification(ctx, mcontext.GetOpUserID(ctx), group.GroupID, constant.GroupMemberSetToAdminNotification, tips)
}

func (g *GroupNotificationSender) GroupMemberSetToOrdinaryUserNotification(ctx context.Context, groupID, groupMemberUserID string) (err error) {
	defer log.ZDebug(ctx, "return")
	defer func() {
		if err != nil {
			log.ZError(ctx, utils.GetFuncName(1)+" failed", err)
		}
	}()
	group, err := g.getGroupInfo(ctx, groupID)
	if err != nil {
		return err
	}
	user, err := g.getGroupMemberMap(ctx, groupID, []string{mcontext.GetOpUserID(ctx), groupMemberUserID})
	if err != nil {
		return err
	}
	tips := &sdkws.GroupMemberInfoSetTips{Group: group, OpUser: user[mcontext.GetOpUserID(ctx)], ChangedUser: user[groupMemberUserID]}
	return g.Notification(ctx, mcontext.GetOpUserID(ctx), group.GroupID, constant.GroupMemberSetToOrdinaryUserNotification, tips)
}

func (g *GroupNotificationSender) MemberEnterDirectlyNotification(ctx context.Context, groupID string, entrantUserID string) (err error) {
	defer log.ZDebug(ctx, "return")
	defer func() {
		if err != nil {
			log.ZError(ctx, utils.GetFuncName(1)+" failed", err)
		}
	}()
	group, err := g.getGroupInfo(ctx, groupID)
	if err != nil {
		return err
	}
	user, err := g.getGroupMember(ctx, groupID, entrantUserID)
	if err != nil {
		return err
	}
	tips := &sdkws.MemberEnterTips{Group: group, EntrantUser: user}
	return g.Notification(ctx, mcontext.GetOpUserID(ctx), group.GroupID, constant.MemberEnterNotification, tips)
}

func (g *GroupNotificationSender) SuperGroupNotification(ctx context.Context, sendID, recvID string) (err error) {
	defer log.ZDebug(ctx, "return")
	defer func() {
		if err != nil {
			log.ZError(ctx, utils.GetFuncName(1)+" failed", err)
		}
	}()
	err = g.Notification(ctx, sendID, recvID, constant.SuperGroupUpdateNotification, nil)
	return err
}
