package notification2

import (
	"context"
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
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func NewGroupNotificationSender(db controller.GroupDatabase, sdr discoveryregistry.SvcDiscoveryRegistry, fn func(ctx context.Context, userIDs []string) ([]rpcclient.CommonUser, error)) *GroupNotificationSender {
	return &GroupNotificationSender{
		msgClient:    rpcclient.NewMsgClient(sdr),
		getUsersInfo: fn,
		db:           db,
	}
}

type GroupNotificationSender struct {
	msgClient *rpcclient.MsgClient
	// 找不到报错
	getUsersInfo func(ctx context.Context, userIDs []string) ([]rpcclient.CommonUser, error)
	db           controller.GroupDatabase
}

func (g *GroupNotificationSender) getGroupInfo(ctx context.Context, groupID string) (*sdkws.GroupInfo, error) {
	gm, err := g.db.TakeGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}
	return &sdkws.GroupInfo{
		GroupID:      gm.GroupID,
		GroupName:    gm.GroupName,
		Notification: gm.Notification,
		Introduction: gm.Introduction,
		FaceURL:      gm.FaceURL,
		//OwnerUserID:            gm.OwnerUserID,
		CreateTime: gm.CreateTime.UnixMilli(),
		//MemberCount:            gm.MemberCount,
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

func (g *GroupNotificationSender) groupNotification(ctx context.Context, contentType int32, m proto.Message, sendID, groupID, recvUserID string) (err error) {
	var tips sdkws.TipsComm
	tips.Detail, err = proto.Marshal(m)
	if err != nil {
		return err
	}
	marshaler := jsonpb.Marshaler{
		OrigName:     true,
		EnumsAsInts:  false,
		EmitDefaults: false,
	}
	tips.JsonDetail, err = marshaler.MarshalToString(m)
	if err != nil {
		return err
	}
	fromUserNickname, toUserNickname, err := g.getFromToUserNickname(ctx, sendID, recvUserID)
	if err != nil {
		return err
	}
	//cn := config.Config.Notification
	switch contentType {
	case constant.GroupCreatedNotification:
		tips.DefaultTips = fromUserNickname
	case constant.GroupInfoSetNotification:
		tips.DefaultTips = fromUserNickname
	case constant.JoinGroupApplicationNotification:
		tips.DefaultTips = fromUserNickname
	case constant.MemberQuitNotification:
		tips.DefaultTips = fromUserNickname
	case constant.GroupApplicationAcceptedNotification:
		tips.DefaultTips = toUserNickname
	case constant.GroupApplicationRejectedNotification:
		tips.DefaultTips = toUserNickname
	case constant.GroupOwnerTransferredNotification:
		tips.DefaultTips = toUserNickname
	case constant.MemberKickedNotification:
		tips.DefaultTips = toUserNickname
	case constant.MemberInvitedNotification:
		tips.DefaultTips = toUserNickname
	case constant.MemberEnterNotification:
		tips.DefaultTips = toUserNickname
	case constant.GroupDismissedNotification:
		tips.DefaultTips = toUserNickname
	case constant.GroupMutedNotification:
		tips.DefaultTips = toUserNickname
	case constant.GroupCancelMutedNotification:
		tips.DefaultTips = toUserNickname
	case constant.GroupMemberMutedNotification:
		tips.DefaultTips = toUserNickname
	case constant.GroupMemberCancelMutedNotification:
		tips.DefaultTips = toUserNickname
	case constant.GroupMemberInfoSetNotification:
		tips.DefaultTips = toUserNickname
	case constant.GroupMemberSetToAdminNotification:
		tips.DefaultTips = toUserNickname
	case constant.GroupMemberSetToOrdinaryUserNotification:
		tips.DefaultTips = toUserNickname
	default:
		return errs.ErrInternalServer.Wrap("unknown group notification type")
	}
	var n rpcclient.NotificationMsg
	n.SendID = sendID
	n.RecvID = recvUserID
	n.ContentType = contentType
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.Content, err = proto.Marshal(&tips)
	if err != nil {
		return
	}
	return g.msgClient.Notification(ctx, &n)
}

func (g *GroupNotificationSender) GroupCreatedNotification(ctx context.Context, group *relation.GroupModel, members []*relation.GroupMemberModel, userMap map[string]*sdkws.UserInfo) (err error) {
	defer log.ZDebug(ctx, "return")
	defer func() {
		if err != nil {
			log.ZError(ctx, "GroupCreatedNotification failed", err)
		}
	}()
	groupInfo, err := g.mergeGroupFull(ctx, group.GroupID, group, &members, &userMap)
	if err != nil {
		return err
	}
	return g.groupNotification(ctx, constant.GroupCreatedNotification, groupInfo, mcontext.GetOpUserID(ctx), group.GroupID, "")
}

func (g *GroupNotificationSender) mergeGroupFull(ctx context.Context, groupID string, group *relation.GroupModel, ms *[]*relation.GroupMemberModel, users *map[string]*sdkws.UserInfo) (groupInfo *sdkws.GroupCreatedTips, err error) {
	defer log.ZDebug(ctx, "return")
	defer func() {
		if err != nil {
			log.ZError(ctx, "mergeGroupFull failed", err)
		}
	}()
	if group == nil {
		group, err = g.db.TakeGroup(ctx, groupID)
		if err != nil {
			return nil, err
		}
	}
	var members []*relation.GroupMemberModel
	if ms == nil || len(*ms) == 0 {
		members, err = g.db.FindGroupMember(ctx, []string{groupID}, nil, nil)
		if err != nil {
			return nil, err
		}
		if ms != nil {
			*ms = members
		}
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
	groupInfo = &sdkws.GroupCreatedTips{Group: g.groupDB2PB(group, opUserID, uint32(len(members))),
		OpUser: opUserMember, GroupOwnerUser: groupOwnerMember}
	return groupInfo, nil
}

func (g *GroupNotificationSender) GroupInfoSetNotification(ctx context.Context, group *relation.GroupModel, members []*relation.GroupMemberModel, needVerification *int32) error {
	groupInfo, err := g.mergeGroupFull(ctx, group.GroupID, group, &members, nil)
	if err != nil {
		return err
	}
	groupInfoChangedTips := &sdkws.GroupInfoSetTips{Group: groupInfo.Group, OpUser: groupInfo.GroupOwnerUser}
	if needVerification != nil {
		groupInfoChangedTips.Group.NeedVerification = *needVerification
	}
	return g.groupNotification(ctx, constant.GroupInfoSetNotification, groupInfoChangedTips, mcontext.GetOpUserID(ctx), group.GroupID, "")
}

//func (g *GroupNotificationSender) mergeGroupAndUser(ctx context.Context, groupID string, userIDs []string) (*sdkws.GroupInfo, map[string]*sdkws.UserInfo, error) {
//	//g.groupDB2PB(group, opUserID, uint32(len(members))
//	groupInfo, err := g.db.TakeGroup(ctx, groupID)
//	if err != nil {
//		return nil, nil, err
//	}
//	owner, err := g.db.TakeGroupOwner(ctx, groupID)
//	if err != nil {
//		return nil, nil, err
//	}
//	memberUserIDs, err := g.db.FindGroupMemberUserID(ctx, groupID)
//	if err != nil {
//		return nil, nil, err
//	}
//	g.getUsersInfoMap(ctx, memberUserIDs)
//	return nil, nil, nil
//}

func (g *GroupNotificationSender) JoinGroupApplicationNotification(ctx context.Context, req *pbGroup.JoinGroupReq) error {
	var members []*relation.GroupMemberModel
	var userMap map[string]*sdkws.UserInfo
	groupInfo, err := g.mergeGroupFull(ctx, req.GroupID, nil, &members, &userMap)
	if err != nil {
		return err
	}
	joinGroupApplicationTips := &sdkws.JoinGroupApplicationTips{Group: groupInfo.Group}
	for _, member := range members {
		if member.UserID == req.InviterUserID {
			if user := userMap[member.UserID]; user != nil {
				joinGroupApplicationTips.Applicant = &sdkws.PublicUserInfo{
					UserID:   user.UserID,
					Nickname: user.Nickname,
					FaceURL:  user.FaceURL,
					Ex:       user.Ex,
				}
			}
			break
		}
	}
	joinGroupApplicationTips.ReqMsg = req.ReqMessage
	for _, member := range members {
		if member.RoleLevel == constant.GroupOwner || member.RoleLevel == constant.GroupAdmin {
			err := g.groupNotification(ctx, constant.JoinGroupApplicationNotification, joinGroupApplicationTips, mcontext.GetOpUserID(ctx), "", member.UserID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *GroupNotificationSender) MemberQuitNotification(ctx context.Context, req *pbGroup.QuitGroupReq) error {
	var members []*relation.GroupMemberModel
	var userMap map[string]*sdkws.UserInfo
	groupInfo, err := g.mergeGroupFull(ctx, req.GroupID, nil, &members, &userMap)
	if err != nil {
		return err
	}
	opUserID := mcontext.GetOpUserID(ctx)
	memberQuitTips := &sdkws.MemberQuitTips{Group: groupInfo.Group, QuitUser: &sdkws.GroupMemberFullInfo{}}
	for _, member := range members {
		if member.UserID == opUserID {
			if user := userMap[member.UserID]; user != nil {
				memberQuitTips.QuitUser = g.groupMemberDB2PB(member, user.AppMangerLevel)
			}
			break
		}
	}
	for _, member := range members {
		if member.RoleLevel == constant.GroupOwner || member.RoleLevel == constant.GroupAdmin {
			err := g.groupNotification(ctx, constant.JoinGroupApplicationNotification, memberQuitTips, mcontext.GetOpUserID(ctx), "", member.UserID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *GroupNotificationSender) GroupApplicationAcceptedNotification(ctx context.Context, req *pbGroup.GroupApplicationResponseReq) error {
	var members []*relation.GroupMemberModel
	groupInfo, err := g.mergeGroupFull(ctx, req.GroupID, nil, &members, nil)
	if err != nil {
		return err
	}
	groupApplicationAcceptedTips := &sdkws.GroupApplicationAcceptedTips{Group: groupInfo.Group, OpUser: groupInfo.OpUser, HandleMsg: req.HandledMsg}
	err = g.groupNotification(ctx, constant.GroupApplicationAcceptedNotification, groupApplicationAcceptedTips, mcontext.GetOpUserID(ctx), "", req.FromUserID)
	if err != nil {
		log.ZError(ctx, "groupNotification failed", err)
	}
	groupApplicationAcceptedTips.ReceiverAs = 1
	for _, member := range members {
		if member.RoleLevel == constant.GroupOwner || member.RoleLevel == constant.GroupAdmin {
			err = g.groupNotification(ctx, constant.GroupApplicationAcceptedNotification, groupApplicationAcceptedTips, mcontext.GetOpUserID(ctx), "", req.FromUserID)
			if err != nil {
				log.ZError(ctx, "groupNotification failed", err)
			}
		}
	}
	return nil
}

func (g *GroupNotificationSender) GroupApplicationRejectedNotification(ctx context.Context, req *pbGroup.GroupApplicationResponseReq) error {
	var members []*relation.GroupMemberModel
	groupInfo, err := g.mergeGroupFull(ctx, req.GroupID, nil, &members, nil)
	if err != nil {
		return err
	}
	groupApplicationRejectedTips := sdkws.GroupApplicationRejectedTips{Group: groupInfo.Group, OpUser: groupInfo.OpUser, HandleMsg: req.HandledMsg}
	if err := g.groupNotification(ctx, constant.GroupApplicationRejectedNotification, &groupApplicationRejectedTips, mcontext.GetOpUserID(ctx), "", req.FromUserID); err != nil {
		log.ZError(ctx, "groupNotification failed", err)
	}
	for _, member := range members {
		if member.RoleLevel == constant.GroupOwner || member.RoleLevel == constant.GroupAdmin {
			if err := g.groupNotification(ctx, constant.GroupApplicationRejectedNotification, &groupApplicationRejectedTips, mcontext.GetOpUserID(ctx), "", req.FromUserID); err != nil {
				log.ZError(ctx, "groupNotification failed", err)
			}
		}
	}
	return nil
}

func (g *GroupNotificationSender) GroupOwnerTransferredNotification(ctx context.Context, req *pbGroup.TransferGroupOwnerReq) error {
	var members []*relation.GroupMemberModel
	groupInfo, err := g.mergeGroupFull(ctx, req.GroupID, nil, &members, nil)
	if err != nil {
		return err
	}
	groupOwnerTransferredTips := &sdkws.GroupOwnerTransferredTips{Group: groupInfo.Group, OpUser: groupInfo.OpUser, NewGroupOwner: groupInfo.GroupOwnerUser}
	return g.groupNotification(ctx, constant.GroupOwnerTransferredNotification, groupOwnerTransferredTips, mcontext.GetOpUserID(ctx), req.GroupID, "")
}

func (g *GroupNotificationSender) MemberKickedNotification(ctx context.Context, req *pbGroup.KickGroupMemberReq, kickedUserIDList []string) error {
	var members []*relation.GroupMemberModel
	groupInfo, err := g.mergeGroupFull(ctx, req.GroupID, nil, &members, nil)
	if err != nil {
		return err
	}
	memberKickedTips := &sdkws.MemberKickedTips{Group: groupInfo.Group, OpUser: groupInfo.OpUser}
	//for _, v := range kickedUserIDList {
	//	var groupMemberInfo sdkws.GroupMemberFullInfo
	//	if err := c.setGroupMemberInfo(ctx, req.GroupID, v, &groupMemberInfo); err != nil {
	//		continue
	//	}
	//	MemberKickedTips.KickedUserList = append(MemberKickedTips.KickedUserList, &groupMemberInfo)
	//}
	return g.groupNotification(ctx, constant.MemberKickedNotification, memberKickedTips, mcontext.GetOpUserID(ctx), req.GroupID, "")
}

func (g *GroupNotificationSender) MemberInvitedNotification(ctx context.Context, groupID, reason string, invitedUserIDList []string) error {
	var members []*relation.GroupMemberModel
	var userMap map[string]*sdkws.UserInfo
	groupInfo, err := g.mergeGroupFull(ctx, groupID, nil, &members, &userMap)
	if err != nil {
		return err
	}
	memberInvitedTips := &sdkws.MemberInvitedTips{Group: groupInfo.Group, OpUser: groupInfo.OpUser}
	groupMembers, err := g.db.FindGroupMember(ctx, []string{groupID}, invitedUserIDList, nil)
	if err != nil {
		return err
	}
	for _, member := range groupMembers {
		user, ok := userMap[member.UserID]
		if !ok {
			continue
		}
		memberInvitedTips.InvitedUserList = append(memberInvitedTips.InvitedUserList, g.groupMemberDB2PB(member, user.AppMangerLevel))
	}
	return g.groupNotification(ctx, constant.MemberInvitedNotification, memberInvitedTips, mcontext.GetOpUserID(ctx), groupID, "")
}

func (g *GroupNotificationSender) MemberEnterNotification(ctx context.Context, req *pbGroup.GroupApplicationResponseReq) error {
	var members []*relation.GroupMemberModel
	var userMap map[string]*sdkws.UserInfo
	groupInfo, err := g.mergeGroupFull(ctx, req.GroupID, nil, &members, &userMap)
	if err != nil {
		return err
	}
	MemberEnterTips := sdkws.MemberEnterTips{Group: groupInfo.Group}
	for _, member := range members {
		if member.UserID == req.FromUserID {
			if user := userMap[member.UserID]; user != nil {
				MemberEnterTips.EntrantUser = g.groupMemberDB2PB(member, user.AppMangerLevel)
			}
			break
		}
	}
	return g.groupNotification(ctx, constant.MemberEnterNotification, &MemberEnterTips, mcontext.GetOpUserID(ctx), req.GroupID, "")
}

func (g *GroupNotificationSender) groupMemberFullInfo(members []*relation.GroupMemberModel, userMap map[string]*sdkws.UserInfo, userID string) *sdkws.GroupMemberFullInfo {
	for _, member := range members {
		if member.UserID == userID {
			if user := userMap[member.UserID]; user != nil {
				return g.groupMemberDB2PB(member, user.AppMangerLevel)
			}
			return g.groupMemberDB2PB(member, 0)
		}
	}
	return nil
}

func (g *GroupNotificationSender) GroupDismissedNotification(ctx context.Context, req *pbGroup.DismissGroupReq) error {
	var members []*relation.GroupMemberModel
	var userMap map[string]*sdkws.UserInfo
	groupInfo, err := g.mergeGroupFull(ctx, req.GroupID, nil, &members, &userMap)
	if err != nil {
		return err
	}
	tips := &sdkws.GroupDismissedTips{Group: groupInfo.Group, OpUser: groupInfo.OpUser}
	return g.groupNotification(ctx, constant.GroupDismissedNotification, tips, mcontext.GetOpUserID(ctx), req.GroupID, "")
}

func (g *GroupNotificationSender) GroupMemberMutedNotification(ctx context.Context, groupID, groupMemberUserID string, mutedSeconds uint32) error {
	var members []*relation.GroupMemberModel
	var userMap map[string]*sdkws.UserInfo
	groupInfo, err := g.mergeGroupFull(ctx, groupID, nil, &members, &userMap)
	if err != nil {
		return err
	}
	tips := sdkws.GroupMemberMutedTips{Group: groupInfo.Group, MutedSeconds: mutedSeconds,
		OpUser: groupInfo.OpUser, MutedUser: g.groupMemberFullInfo(members, userMap, groupMemberUserID)}
	tips.MutedSeconds = mutedSeconds
	return g.groupNotification(ctx, constant.GroupMemberMutedNotification, &tips, mcontext.GetOpUserID(ctx), groupID, "")
}

func (g *GroupNotificationSender) GroupMemberCancelMutedNotification(ctx context.Context, groupID, groupMemberUserID string) error {
	var members []*relation.GroupMemberModel
	var userMap map[string]*sdkws.UserInfo
	groupInfo, err := g.mergeGroupFull(ctx, groupID, nil, &members, &userMap)
	if err != nil {
		return err
	}
	tips := sdkws.GroupMemberCancelMutedTips{Group: groupInfo.Group,
		OpUser: groupInfo.OpUser, MutedUser: g.groupMemberFullInfo(members, userMap, groupMemberUserID)}
	return g.groupNotification(ctx, constant.GroupMemberCancelMutedNotification, &tips, mcontext.GetOpUserID(ctx), groupID, "")
}

func (g *GroupNotificationSender) GroupMutedNotification(ctx context.Context, groupID string) error {
	var members []*relation.GroupMemberModel
	var userMap map[string]*sdkws.UserInfo
	groupInfo, err := g.mergeGroupFull(ctx, groupID, nil, &members, &userMap)
	if err != nil {
		return err
	}
	tips := sdkws.GroupMutedTips{Group: groupInfo.Group, OpUser: groupInfo.OpUser}
	return g.groupNotification(ctx, constant.GroupMutedNotification, &tips, mcontext.GetOpUserID(ctx), groupID, "")
}

func (g *GroupNotificationSender) GroupCancelMutedNotification(ctx context.Context, groupID string) error {
	var members []*relation.GroupMemberModel
	var userMap map[string]*sdkws.UserInfo
	groupInfo, err := g.mergeGroupFull(ctx, groupID, nil, &members, &userMap)
	if err != nil {
		return err
	}
	tips := sdkws.GroupCancelMutedTips{Group: groupInfo.Group, OpUser: groupInfo.OpUser}
	return g.groupNotification(ctx, constant.GroupMutedNotification, &tips, mcontext.GetOpUserID(ctx), groupID, "")
}

func (g *GroupNotificationSender) GroupMemberInfoSetNotification(ctx context.Context, groupID, groupMemberUserID string) error {
	var members []*relation.GroupMemberModel
	var userMap map[string]*sdkws.UserInfo
	groupInfo, err := g.mergeGroupFull(ctx, groupID, nil, &members, &userMap)
	if err != nil {
		return err
	}
	tips := sdkws.GroupMemberInfoSetTips{Group: groupInfo.Group,
		OpUser: groupInfo.OpUser, ChangedUser: g.groupMemberFullInfo(members, userMap, groupMemberUserID)}
	return g.groupNotification(ctx, constant.GroupMemberCancelMutedNotification, &tips, mcontext.GetOpUserID(ctx), groupID, "")
}

func (g *GroupNotificationSender) GroupMemberSetToAdminNotification(ctx context.Context, groupID, groupMemberUserID string, notificationType int32) error {
	var members []*relation.GroupMemberModel
	var userMap map[string]*sdkws.UserInfo
	groupInfo, err := g.mergeGroupFull(ctx, groupID, nil, &members, &userMap)
	if err != nil {
		return err
	}
	tips := sdkws.GroupMemberInfoSetTips{Group: groupInfo.Group,
		OpUser: groupInfo.OpUser, ChangedUser: g.groupMemberFullInfo(members, userMap, groupMemberUserID)}
	return g.groupNotification(ctx, constant.GroupMemberCancelMutedNotification, &tips, mcontext.GetOpUserID(ctx), groupID, "")
}

func (g *GroupNotificationSender) MemberEnterDirectlyNotification(ctx context.Context, groupID string, entrantUserID string) error {
	var members []*relation.GroupMemberModel
	var userMap map[string]*sdkws.UserInfo
	groupInfo, err := g.mergeGroupFull(ctx, groupID, nil, &members, &userMap)
	if err != nil {
		return err
	}
	tips := sdkws.MemberEnterTips{Group: groupInfo.Group, EntrantUser: g.groupMemberFullInfo(members, userMap, entrantUserID)}
	return g.groupNotification(ctx, constant.GroupMemberCancelMutedNotification, &tips, mcontext.GetOpUserID(ctx), groupID, "")
}

type NotificationMsg struct {
	SendID         string
	RecvID         string
	Content        []byte //  sdkws.TipsComm
	MsgFrom        int32
	ContentType    int32
	SessionType    int32
	SenderNickname string
	SenderFaceURL  string
}

func (g *GroupNotificationSender) SuperGroupNotification(ctx context.Context, sendID, recvID string) error {
	n := &NotificationMsg{
		SendID:      sendID,
		RecvID:      recvID,
		MsgFrom:     constant.SysMsgType,
		ContentType: constant.SuperGroupUpdateNotification,
		SessionType: constant.SingleChatType,
	}
	_ = n // todo
	//g.Notification(ctx, n)
	return nil
}
