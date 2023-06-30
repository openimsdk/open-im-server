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
	num, err := g.db.FindGroupMemberNum(ctx, groupID)
	if err != nil {
		return nil, err
	}
	owner, err := g.db.TakeGroupOwner(ctx, groupID)
	if err != nil {
		return nil, err
	}
	return &sdkws.GroupInfo{
		GroupID:                gm.GroupID,
		GroupName:              gm.GroupName,
		Notification:           gm.Notification,
		Introduction:           gm.Introduction,
		FaceURL:                gm.FaceURL,
		OwnerUserID:            owner.UserID,
		CreateTime:             gm.CreateTime.UnixMilli(),
		MemberCount:            num,
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
	members, err := g.db.FindGroupMember(ctx, []string{groupID}, userIDs, nil)
	if err != nil {
		return nil, err
	}
	log.ZDebug(ctx, "getGroupMembers", "members", members)
	users, err := g.getUsersInfoMap(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	log.ZDebug(ctx, "getUsersInfoMap", "users", users)
	res := make([]*sdkws.GroupMemberFullInfo, 0, len(members))
	for _, member := range members {
		user, ok := users[member.UserID]
		if !ok {
			return nil, errs.ErrUserIDNotFound.Wrap(fmt.Sprintf("group %s member %s not in user", member.GroupID, member.UserID))
		}
		if member.Nickname == "" {
			member.Nickname = user.Nickname
		}
		res = append(res, g.groupMemberDB2PB(member, user.AppMangerLevel))
		delete(users, member.UserID)
	}
	//for userID, info := range users {
	//	if info.AppMangerLevel == constant.AppAdmin {
	//		res = append(res, &sdkws.GroupMemberFullInfo{
	//			GroupID:        groupID,
	//			UserID:         userID,
	//			Nickname:       info.Nickname,
	//			FaceURL:        info.FaceURL,
	//			AppMangerLevel: info.AppMangerLevel,
	//		})
	//	}
	//}
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

func (g *GroupNotificationSender) fillOpUser(ctx context.Context, opUser **sdkws.GroupMemberFullInfo, groupID string) error {
	if opUser == nil {
		return errs.ErrInternalServer.Wrap("**sdkws.GroupMemberFullInfo is nil")
	}
	if *opUser != nil {
		return nil
	}
	userID := mcontext.GetOpUserID(ctx)
	if groupID != "" {
		member, err := g.db.TakeGroupMember(ctx, groupID, userID)
		if err == nil {
			*opUser = g.groupMemberDB2PB(member, 0)
		} else if !errs.ErrRecordNotFound.Is(err) {
			return err
		}
	}
	user, err := g.getUser(ctx, userID)
	if err != nil {
		return err
	}
	if *opUser == nil {
		*opUser = &sdkws.GroupMemberFullInfo{
			GroupID:        groupID,
			UserID:         userID,
			Nickname:       user.Nickname,
			FaceURL:        user.FaceURL,
			OperatorUserID: userID,
		}
	} else {
		if (*opUser).Nickname == "" {
			(*opUser).Nickname = user.Nickname
		}
		if (*opUser).FaceURL == "" {
			(*opUser).FaceURL = user.FaceURL
		}
	}
	return nil
}

func (g *GroupNotificationSender) GroupCreatedNotification(ctx context.Context, tips *sdkws.GroupCreatedTips) (err error) {
	if err := g.fillOpUser(ctx, &tips.OpUser, tips.Group.GroupID); err != nil {
		return err
	}
	return g.Notification(ctx, mcontext.GetOpUserID(ctx), tips.Group.GroupID, constant.GroupCreatedNotification, tips)
}

func (g *GroupNotificationSender) GroupInfoSetNotification(ctx context.Context, tips *sdkws.GroupInfoSetTips) (err error) {
	if err := g.fillOpUser(ctx, &tips.OpUser, tips.Group.GroupID); err != nil {
		return err
	}
	return g.Notification(ctx, mcontext.GetOpUserID(ctx), tips.Group.GroupID, constant.GroupInfoSetNotification, tips)
}

func (g *GroupNotificationSender) GroupInfoSetNameNotification(ctx context.Context, tips *sdkws.GroupInfoSetNameTips) (err error) {
	if err := g.fillOpUser(ctx, &tips.OpUser, tips.Group.GroupID); err != nil {
		return err
	}
	return g.Notification(ctx, mcontext.GetOpUserID(ctx), tips.Group.GroupID, constant.GroupInfoSetNameNotification, tips)
}

func (g *GroupNotificationSender) GroupInfoSetAnnouncementNotification(ctx context.Context, tips *sdkws.GroupInfoSetAnnouncementTips) (err error) {
	if err := g.fillOpUser(ctx, &tips.OpUser, tips.Group.GroupID); err != nil {
		return err
	}
	return g.Notification(ctx, mcontext.GetOpUserID(ctx), tips.Group.GroupID, constant.GroupInfoSetAnnouncementNotification, tips)
}

func (g *GroupNotificationSender) JoinGroupApplicationNotification(ctx context.Context, req *pbGroup.JoinGroupReq) (err error) {
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

func (g *GroupNotificationSender) MemberQuitNotification(ctx context.Context, member *sdkws.GroupMemberFullInfo) (err error) {
	defer log.ZDebug(ctx, "return")
	defer func() {
		if err != nil {
			log.ZError(ctx, utils.GetFuncName(1)+" failed", err)
		}
	}()
	group, err := g.getGroupInfo(ctx, member.GroupID)
	if err != nil {
		return err
	}
	tips := &sdkws.MemberQuitTips{Group: group, QuitUser: member}
	return g.Notification(ctx, mcontext.GetOpUserID(ctx), member.GroupID, constant.MemberQuitNotification, tips)
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
	userIDs, err := g.getGroupOwnerAndAdminUserID(ctx, req.GroupID)
	if err != nil {
		return err
	}
	tips := &sdkws.GroupApplicationAcceptedTips{Group: group, HandleMsg: req.HandledMsg, ReceiverAs: 1}
	if err := g.fillOpUser(ctx, &tips.OpUser, tips.Group.GroupID); err != nil {
		return err
	}
	for _, userID := range append(userIDs, mcontext.GetOpUserID(ctx)) {
		err = g.Notification(ctx, mcontext.GetOpUserID(ctx), userID, constant.GroupApplicationAcceptedNotification, tips)
		if err != nil {
			log.ZError(ctx, "failed", err)
		}
	}
	return nil
}

func (g *GroupNotificationSender) GroupApplicationRejectedNotification(ctx context.Context, req *pbGroup.GroupApplicationResponseReq) (err error) {
	group, err := g.getGroupInfo(ctx, req.GroupID)
	if err != nil {
		return err
	}
	userIDs, err := g.getGroupOwnerAndAdminUserID(ctx, req.GroupID)
	if err != nil {
		return err
	}
	tips := &sdkws.GroupApplicationRejectedTips{Group: group, HandleMsg: req.HandledMsg}
	if err := g.fillOpUser(ctx, &tips.OpUser, tips.Group.GroupID); err != nil {
		return err
	}
	for _, userID := range append(userIDs, mcontext.GetOpUserID(ctx)) {
		err = g.Notification(ctx, mcontext.GetOpUserID(ctx), userID, constant.GroupApplicationRejectedNotification, tips)
		if err != nil {
			log.ZError(ctx, "failed", err)
		}
	}
	return nil
}

func (g *GroupNotificationSender) GroupOwnerTransferredNotification(ctx context.Context, req *pbGroup.TransferGroupOwnerReq) (err error) {
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
	if err := g.fillOpUser(ctx, &tips.OpUser, tips.Group.GroupID); err != nil {
		return err
	}
	return g.Notification(ctx, mcontext.GetOpUserID(ctx), group.GroupID, constant.GroupOwnerTransferredNotification, tips)
}

func (g *GroupNotificationSender) MemberKickedNotification(ctx context.Context, tips *sdkws.MemberKickedTips) (err error) {
	if err := g.fillOpUser(ctx, &tips.OpUser, tips.Group.GroupID); err != nil {
		return err
	}
	return g.Notification(ctx, mcontext.GetOpUserID(ctx), tips.Group.GroupID, constant.MemberKickedNotification, tips)
}

func (g *GroupNotificationSender) MemberInvitedNotification(ctx context.Context, groupID, reason string, invitedUserIDList []string) (err error) {
	group, err := g.getGroupInfo(ctx, groupID)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	users, err := g.getGroupMembers(ctx, groupID, invitedUserIDList)
	if err != nil {
		return err
	}
	tips := &sdkws.MemberInvitedTips{Group: group, InvitedUserList: users}
	if err := g.fillOpUser(ctx, &tips.OpUser, tips.Group.GroupID); err != nil {
		return err
	}
	return g.Notification(ctx, mcontext.GetOpUserID(ctx), group.GroupID, constant.MemberInvitedNotification, tips)
}

func (g *GroupNotificationSender) MemberEnterNotification(ctx context.Context, req *pbGroup.GroupApplicationResponseReq) (err error) {
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
	if err := g.fillOpUser(ctx, &tips.OpUser, tips.Group.GroupID); err != nil {
		return err
	}
	return g.Notification(ctx, mcontext.GetOpUserID(ctx), tips.Group.GroupID, constant.GroupDismissedNotification, tips)
}

func (g *GroupNotificationSender) GroupMemberMutedNotification(ctx context.Context, groupID, groupMemberUserID string, mutedSeconds uint32) (err error) {
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
	if err := g.fillOpUser(ctx, &tips.OpUser, tips.Group.GroupID); err != nil {
		return err
	}
	return g.Notification(ctx, mcontext.GetOpUserID(ctx), group.GroupID, constant.GroupMemberMutedNotification, tips)
}

func (g *GroupNotificationSender) GroupMemberCancelMutedNotification(ctx context.Context, groupID, groupMemberUserID string) (err error) {
	group, err := g.getGroupInfo(ctx, groupID)
	if err != nil {
		return err
	}
	user, err := g.getGroupMemberMap(ctx, groupID, []string{mcontext.GetOpUserID(ctx), groupMemberUserID})
	if err != nil {
		return err
	}
	tips := &sdkws.GroupMemberCancelMutedTips{Group: group, OpUser: user[mcontext.GetOpUserID(ctx)], MutedUser: user[groupMemberUserID]}
	if err := g.fillOpUser(ctx, &tips.OpUser, tips.Group.GroupID); err != nil {
		return err
	}
	return g.Notification(ctx, mcontext.GetOpUserID(ctx), group.GroupID, constant.GroupMemberCancelMutedNotification, tips)
}

func (g *GroupNotificationSender) GroupMutedNotification(ctx context.Context, groupID string) (err error) {
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
	}
	if err := g.fillOpUser(ctx, &tips.OpUser, tips.Group.GroupID); err != nil {
		return err
	}
	return g.Notification(ctx, mcontext.GetOpUserID(ctx), group.GroupID, constant.GroupMutedNotification, tips)
}

func (g *GroupNotificationSender) GroupCancelMutedNotification(ctx context.Context, groupID string) (err error) {
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
	}
	if err := g.fillOpUser(ctx, &tips.OpUser, tips.Group.GroupID); err != nil {
		return err
	}
	return g.Notification(ctx, mcontext.GetOpUserID(ctx), group.GroupID, constant.GroupCancelMutedNotification, tips)
}

func (g *GroupNotificationSender) GroupMemberInfoSetNotification(ctx context.Context, groupID, groupMemberUserID string) (err error) {
	group, err := g.getGroupInfo(ctx, groupID)
	if err != nil {
		return err
	}
	user, err := g.getGroupMemberMap(ctx, groupID, []string{groupMemberUserID})
	if err != nil {
		return err
	}
	tips := &sdkws.GroupMemberInfoSetTips{Group: group, OpUser: user[mcontext.GetOpUserID(ctx)], ChangedUser: user[groupMemberUserID]}
	if err := g.fillOpUser(ctx, &tips.OpUser, tips.Group.GroupID); err != nil {
		return err
	}
	return g.Notification(ctx, mcontext.GetOpUserID(ctx), group.GroupID, constant.GroupMemberInfoSetNotification, tips)
}

func (g *GroupNotificationSender) GroupMemberSetToAdminNotification(ctx context.Context, groupID, groupMemberUserID string) (err error) {
	group, err := g.getGroupInfo(ctx, groupID)
	if err != nil {
		return err
	}
	user, err := g.getGroupMemberMap(ctx, groupID, []string{mcontext.GetOpUserID(ctx), groupMemberUserID})
	if err != nil {
		return err
	}
	tips := &sdkws.GroupMemberInfoSetTips{Group: group, OpUser: user[mcontext.GetOpUserID(ctx)], ChangedUser: user[groupMemberUserID]}
	if err := g.fillOpUser(ctx, &tips.OpUser, tips.Group.GroupID); err != nil {
		return err
	}
	return g.Notification(ctx, mcontext.GetOpUserID(ctx), group.GroupID, constant.GroupMemberSetToAdminNotification, tips)
}

func (g *GroupNotificationSender) GroupMemberSetToOrdinaryUserNotification(ctx context.Context, groupID, groupMemberUserID string) (err error) {
	group, err := g.getGroupInfo(ctx, groupID)
	if err != nil {
		return err
	}
	user, err := g.getGroupMemberMap(ctx, groupID, []string{mcontext.GetOpUserID(ctx), groupMemberUserID})
	if err != nil {
		return err
	}
	tips := &sdkws.GroupMemberInfoSetTips{Group: group, OpUser: user[mcontext.GetOpUserID(ctx)], ChangedUser: user[groupMemberUserID]}
	if err := g.fillOpUser(ctx, &tips.OpUser, tips.Group.GroupID); err != nil {
		return err
	}
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
