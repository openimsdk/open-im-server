package group

import (
	"Open_IM/internal/common/check"
	"Open_IM/internal/common/notification"
	"Open_IM/internal/tx"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/cache"
	"Open_IM/pkg/common/db/controller"
	"Open_IM/pkg/common/db/relation"
	relationTb "Open_IM/pkg/common/db/table/relation"
	"Open_IM/pkg/common/db/unrelation"
	"Open_IM/pkg/common/tokenverify"
	"Open_IM/pkg/common/tracelog"
	pbConversation "Open_IM/pkg/proto/conversation"
	pbGroup "Open_IM/pkg/proto/group"
	"Open_IM/pkg/proto/sdkws"
	"Open_IM/pkg/utils"
	"context"
	"fmt"
	"github.com/OpenIMSDK/openKeeper"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func Start(client *openKeeper.ZkClient, server *grpc.Server) error {
	db, err := relation.NewGormDB()
	if err != nil {
		return err
	}
	if err := db.AutoMigrate(&relationTb.GroupModel{}, &relationTb.GroupMemberModel{}, &relationTb.GroupRequestModel{}); err != nil {
		return err
	}
	redis, err := cache.NewRedis()
	if err != nil {
		return err
	}
	mongo, err := unrelation.NewMongo()
	if err != nil {
		return err
	}
	pbGroup.RegisterGroupServer(server, &groupServer{
		GroupInterface: controller.NewGroupController(
			relation.NewGroupDB(db),
			relation.NewGroupMemberDB(db),
			relation.NewGroupRequest(db),
			tx.NewGorm(db),
			tx.NewMongo(mongo.GetClient()),
			unrelation.NewSuperGroupMongoDriver(mongo.GetClient()),
			redis.GetClient(),
		),
		UserCheck:           check.NewUserCheck(client),
		ConversationChecker: check.NewConversationChecker(client),
	})
	return nil
}

type groupServer struct {
	GroupInterface      controller.GroupController
	UserCheck           *check.UserCheck
	Notification        *notification.Check
	ConversationChecker *check.ConversationChecker
}

func (s *groupServer) CheckGroupAdmin(ctx context.Context, groupID string) error {
	if !tokenverify.IsAppManagerUid(ctx) {
		groupMember, err := s.GroupInterface.TakeGroupMember(ctx, groupID, tracelog.GetOpUserID(ctx))
		if err != nil {
			return err
		}
		if !(groupMember.RoleLevel == constant.GroupOwner || groupMember.RoleLevel == constant.GroupAdmin) {
			return constant.ErrNoPermission.Wrap("no group owner or admin")
		}
	}
	return nil
}

func (s *groupServer) GetUsernameMap(ctx context.Context, userIDs []string, complete bool) (map[string]string, error) {
	if len(userIDs) == 0 {
		return map[string]string{}, nil
	}
	users, err := s.UserCheck.GetPublicUserInfos(ctx, userIDs, complete)
	if err != nil {
		return nil, err
	}
	return utils.SliceToMapAny(users, func(e *sdkws.PublicUserInfo) (string, string) {
		return e.UserID, e.Nickname
	}), nil
}

func (s *groupServer) IsNotFound(err error) bool {
	return utils.Unwrap(err) == gorm.ErrRecordNotFound
}

func (s *groupServer) GenGroupID(ctx context.Context, groupID *string) error {
	if *groupID != "" {
		_, err := s.GroupInterface.TakeGroup(ctx, *groupID)
		if err == nil {
			return constant.ErrGroupIDExisted.Wrap("group id existed " + *groupID)
		} else if s.IsNotFound(err) {
			return nil
		} else {
			return err
		}
	}
	for i := 0; i < 10; i++ {
		id := utils.Md5(strings.Join([]string{tracelog.GetOperationID(ctx), strconv.FormatInt(time.Now().UnixNano(), 10), strconv.Itoa(rand.Int())}, ",;,"))
		bi := big.NewInt(0)
		bi.SetString(id[0:8], 16)
		id = bi.String()
		_, err := s.GroupInterface.TakeGroup(ctx, id)
		if err == nil {
			continue
		} else if s.IsNotFound(err) {
			*groupID = id
			return nil
		} else {
			return err
		}
	}
	return constant.ErrData.Wrap("group id gen error")
}

func (s *groupServer) CreateGroup(ctx context.Context, req *pbGroup.CreateGroupReq) (*pbGroup.CreateGroupResp, error) {
	resp := &pbGroup.CreateGroupResp{GroupInfo: &sdkws.GroupInfo{}}
	if err := tokenverify.CheckAccessV3(ctx, req.OwnerUserID); err != nil {
		return nil, err
	}
	if req.OwnerUserID == "" {
		return nil, constant.ErrArgs.Wrap("no group owner")
	}
	userIDs := append(append(req.InitMembers, req.AdminUserIDs...), req.OwnerUserID)
	if utils.Duplicate(userIDs) {
		return nil, constant.ErrArgs.Wrap("group member repeated")
	}
	userMap, err := s.UserCheck.GetUsersInfoMap(ctx, userIDs, true)
	if err != nil {
		return nil, err
	}
	if err := CallbackBeforeCreateGroup(ctx, req); err != nil {
		return nil, err
	}
	var groupMembers []*relationTb.GroupMemberModel
	group := PbToDBGroupInfo(req.GroupInfo)
	if err := s.GenGroupID(ctx, &group.GroupID); err != nil {
		return nil, err
	}
	joinGroup := func(userID string, roleLevel int32) error {
		groupMember := PbToDbGroupMember(userMap[userID])
		groupMember.Nickname = ""
		groupMember.GroupID = group.GroupID
		groupMember.RoleLevel = roleLevel
		groupMember.OperatorUserID = tracelog.GetOpUserID(ctx)
		groupMember.JoinSource = constant.JoinByInvitation
		groupMember.InviterUserID = tracelog.GetOpUserID(ctx)
		if err := CallbackBeforeMemberJoinGroup(ctx, groupMember, group.Ex); err != nil {
			return err
		}
		groupMembers = append(groupMembers, groupMember)
		return nil
	}
	if err := joinGroup(req.OwnerUserID, constant.GroupOwner); err != nil {
		return nil, err
	}
	if req.GroupInfo.GroupType == constant.SuperGroup {
		if err := s.GroupInterface.CreateSuperGroup(ctx, group.GroupID, userIDs); err != nil {
			return nil, err
		}
	} else {
		for _, userID := range req.AdminUserIDs {
			if err := joinGroup(userID, constant.GroupAdmin); err != nil {
				return nil, err
			}
		}
		for _, userID := range req.InitMembers {
			if err := joinGroup(userID, constant.GroupOrdinaryUsers); err != nil {
				return nil, err
			}
		}
	}
	if err := s.GroupInterface.CreateGroup(ctx, []*relationTb.GroupModel{group}, groupMembers); err != nil {
		return nil, err
	}
	resp.GroupInfo = DbToPbGroupInfo(group, req.OwnerUserID, uint32(len(userIDs)))
	resp.GroupInfo.MemberCount = uint32(len(userIDs))
	if req.GroupInfo.GroupType == constant.SuperGroup {
		go func() {
			for _, userID := range userIDs {
				s.Notification.SuperGroupNotification(ctx, userID, userID)
			}
		}()
	} else {
		s.Notification.GroupCreatedNotification(ctx, group.GroupID, userIDs)
	}
	return resp, nil
}

func (s *groupServer) GetJoinedGroupList(ctx context.Context, req *pbGroup.GetJoinedGroupListReq) (*pbGroup.GetJoinedGroupListResp, error) {
	resp := &pbGroup.GetJoinedGroupListResp{}
	if err := tokenverify.CheckAccessV3(ctx, req.FromUserID); err != nil {
		return nil, err
	}
	total, members, err := s.GroupInterface.PageGroupMember(ctx, nil, []string{req.FromUserID}, nil, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	if err != nil {
		return nil, err
	}
	resp.Total = total
	if len(members) == 0 {
		return resp, nil
	}
	groupIDs := utils.Slice(members, func(e *relationTb.GroupMemberModel) string {
		return e.GroupID
	})
	groups, err := s.GroupInterface.FindGroup(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	groupMemberNum, err := s.GroupInterface.MapGroupMemberNum(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	owners, err := s.GroupInterface.FindGroupMember(ctx, groupIDs, nil, []int32{constant.GroupOwner})
	if err != nil {
		return nil, err
	}
	ownerMap := utils.SliceToMap(owners, func(e *relationTb.GroupMemberModel) string {
		return e.GroupID
	})
	resp.Groups = utils.Slice(utils.Order(groupIDs, groups, func(group *relationTb.GroupModel) string {
		return group.GroupID
	}), func(group *relationTb.GroupModel) *sdkws.GroupInfo {
		return DbToPbGroupInfo(group, ownerMap[group.GroupID].UserID, groupMemberNum[group.GroupID])
	})
	return resp, nil
}

func (s *groupServer) InviteUserToGroup(ctx context.Context, req *pbGroup.InviteUserToGroupReq) (*pbGroup.InviteUserToGroupResp, error) {
	resp := &pbGroup.InviteUserToGroupResp{}
	if len(req.InvitedUserIDs) == 0 {
		return nil, constant.ErrArgs.Wrap("user empty")
	}
	if utils.Duplicate(req.InvitedUserIDs) {
		return nil, constant.ErrArgs.Wrap("userID duplicate")
	}
	group, err := s.GroupInterface.TakeGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if group.Status == constant.GroupStatusDismissed {
		return nil, constant.ErrDismissedAlready.Wrap()
	}
	members, err := s.GroupInterface.FindGroupMember(ctx, []string{group.GroupID}, nil, nil)
	if err != nil {
		return nil, err
	}
	memberMap := utils.SliceToMap(members, func(e *relationTb.GroupMemberModel) string {
		return e.UserID
	})
	if ids := utils.Single(req.InvitedUserIDs, utils.Keys(memberMap)); len(ids) > 0 {
		return nil, constant.ErrArgs.Wrap("user in group " + strings.Join(ids, ","))
	}
	userMap, err := s.UserCheck.GetUsersInfoMap(ctx, req.InvitedUserIDs, true)
	if err != nil {
		return nil, err
	}
	if group.NeedVerification == constant.AllNeedVerification {
		if !tokenverify.IsAppManagerUid(ctx) {
			opUserID := tracelog.GetOpUserID(ctx)
			member, ok := memberMap[opUserID]
			if !ok {
				return nil, constant.ErrNoPermission.Wrap("not in group")
			}
			if !(member.RoleLevel == constant.GroupOwner || member.RoleLevel == constant.GroupAdmin) {
				var requests []*relationTb.GroupRequestModel
				for _, userID := range req.InvitedUserIDs {
					requests = append(requests, &relationTb.GroupRequestModel{
						UserID:        userID,
						GroupID:       req.GroupID,
						JoinSource:    constant.JoinByInvitation,
						InviterUserID: opUserID,
					})
				}
				if err := s.GroupInterface.CreateGroupRequest(ctx, requests); err != nil {
					return nil, err
				}
				for _, request := range requests {
					s.Notification.JoinGroupApplicationNotification(ctx, &pbGroup.JoinGroupReq{
						GroupID:       request.GroupID,
						ReqMessage:    request.ReqMsg,
						JoinSource:    request.JoinSource,
						InviterUserID: request.InviterUserID,
					})
				}
				return resp, nil
			}
		}
	}
	if group.GroupType == constant.SuperGroup {
		if err := s.GroupInterface.CreateSuperGroupMember(ctx, req.GroupID, req.InvitedUserIDs); err != nil {
			return nil, err
		}
		for _, userID := range req.InvitedUserIDs {
			s.Notification.SuperGroupNotification(ctx, userID, userID)
		}
	} else {
		opUserID := tracelog.GetOpUserID(ctx)
		var groupMembers []*relationTb.GroupMemberModel
		for _, userID := range req.InvitedUserIDs {
			member := PbToDbGroupMember(userMap[userID])
			member.Nickname = ""
			member.GroupID = req.GroupID
			member.RoleLevel = constant.GroupOrdinaryUsers
			member.OperatorUserID = opUserID
			member.InviterUserID = opUserID
			member.JoinSource = constant.JoinByInvitation
			if err := CallbackBeforeMemberJoinGroup(ctx, member, group.Ex); err != nil {
				return nil, err
			}
			groupMembers = append(groupMembers, member)
		}
		if err := s.GroupInterface.CreateGroup(ctx, nil, groupMembers); err != nil {
			return nil, err
		}
		s.Notification.MemberInvitedNotification(ctx, req.GroupID, req.Reason, req.InvitedUserIDs)
	}
	return resp, nil
}

func (s *groupServer) GetGroupAllMember(ctx context.Context, req *pbGroup.GetGroupAllMemberReq) (*pbGroup.GetGroupAllMemberResp, error) {
	resp := &pbGroup.GetGroupAllMemberResp{}
	group, err := s.GroupInterface.TakeGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if group.GroupType == constant.SuperGroup {
		return nil, constant.ErrArgs.Wrap("unsupported super group")
	}
	members, err := s.GroupInterface.FindGroupMember(ctx, []string{req.GroupID}, nil, nil)
	if err != nil {
		return nil, err
	}
	nameMap, err := s.GetUsernameMap(ctx, utils.Filter(members, func(e *relationTb.GroupMemberModel) (string, bool) {
		return e.UserID, e.Nickname == ""
	}), true)
	if err != nil {
		return nil, err
	}
	resp.Members = utils.Slice(members, func(e *relationTb.GroupMemberModel) *sdkws.GroupMemberFullInfo {
		if e.Nickname == "" {
			e.Nickname = nameMap[e.UserID]
		}
		return DbToPbGroupMembersCMSResp(e)
	})
	return resp, nil
}

func (s *groupServer) GetGroupMemberList(ctx context.Context, req *pbGroup.GetGroupMemberListReq) (*pbGroup.GetGroupMemberListResp, error) {
	resp := &pbGroup.GetGroupMemberListResp{}
	total, members, err := s.GroupInterface.PageGroupMember(ctx, []string{req.GroupID}, nil, utils.If(req.Filter >= 0, []int32{req.Filter}, nil), req.Pagination.PageNumber, req.Pagination.ShowNumber)
	if err != nil {
		return nil, err
	}
	resp.Total = total
	nameMap, err := s.GetUsernameMap(ctx, utils.Filter(members, func(e *relationTb.GroupMemberModel) (string, bool) {
		return e.UserID, e.Nickname == ""
	}), true)
	if err != nil {
		return nil, err
	}
	resp.Members = utils.Slice(members, func(e *relationTb.GroupMemberModel) *sdkws.GroupMemberFullInfo {
		if e.Nickname == "" {
			e.Nickname = nameMap[e.UserID]
		}
		return DbToPbGroupMembersCMSResp(e)
	})
	return resp, nil
}

func (s *groupServer) KickGroupMember(ctx context.Context, req *pbGroup.KickGroupMemberReq) (*pbGroup.KickGroupMemberResp, error) {
	resp := &pbGroup.KickGroupMemberResp{}
	group, err := s.GroupInterface.TakeGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if len(req.KickedUserIDs) == 0 {
		return nil, constant.ErrArgs.Wrap("KickedUserIDs empty")
	}
	if utils.IsDuplicateStringSlice(req.KickedUserIDs) {
		return nil, constant.ErrArgs.Wrap("KickedUserIDs duplicate")
	}
	opUserID := tracelog.GetOpUserID(ctx)
	if utils.IsContain(opUserID, req.KickedUserIDs) {
		return nil, constant.ErrArgs.Wrap("opUserID in KickedUserIDs")
	}
	if group.GroupType == constant.SuperGroup {
		if err := s.GroupInterface.DeleteSuperGroupMember(ctx, req.GroupID, req.KickedUserIDs); err != nil {
			return nil, err
		}
		go func() {
			for _, userID := range req.KickedUserIDs {
				s.Notification.SuperGroupNotification(ctx, userID, userID)
			}
		}()
	} else {
		members, err := s.GroupInterface.FindGroupMember(ctx, []string{req.GroupID}, append(req.KickedUserIDs, opUserID), nil)
		if err != nil {
			return nil, err
		}
		memberMap := make(map[string]*relationTb.GroupMemberModel)
		for i, member := range members {
			memberMap[member.UserID] = members[i]
		}
		for _, userID := range req.KickedUserIDs {
			if _, ok := memberMap[userID]; !ok {
				return nil, constant.ErrUserIDNotFound.Wrap(userID)
			}
		}
		if !tokenverify.IsAppManagerUid(ctx) {
			member := memberMap[opUserID]
			if member == nil {
				return nil, constant.ErrNoPermission.Wrap(fmt.Sprintf("opUserID %s no in group", opUserID))
			}
			switch member.RoleLevel {
			case constant.GroupOwner:
			case constant.GroupAdmin:
				for _, member := range members {
					if member.UserID == opUserID {
						continue
					}
					if member.RoleLevel == constant.GroupOwner || member.RoleLevel == constant.GroupAdmin {
						return nil, constant.ErrNoPermission.Wrap("userID:" + member.UserID)
					}
				}
			default:
				return nil, constant.ErrNoPermission.Wrap("opUserID is OrdinaryUser")
			}
		}
		if err := s.GroupInterface.DeleteGroupMember(ctx, group.GroupID, req.KickedUserIDs); err != nil {
			return nil, err
		}
		s.Notification.MemberKickedNotification(ctx, req, req.KickedUserIDs)
	}
	return resp, nil
}

func (s *groupServer) GetGroupMembersInfo(ctx context.Context, req *pbGroup.GetGroupMembersInfoReq) (*pbGroup.GetGroupMembersInfoResp, error) {
	resp := &pbGroup.GetGroupMembersInfoResp{}
	if len(req.Members) == 0 {
		return nil, constant.ErrArgs.Wrap("members empty")
	}
	if req.GroupID == "" {
		return nil, constant.ErrArgs.Wrap("groupID empty")
	}
	members, err := s.GroupInterface.FindGroupMember(ctx, []string{req.GroupID}, req.Members, nil)
	if err != nil {
		return nil, err
	}
	nameMap, err := s.GetUsernameMap(ctx, utils.Filter(members, func(e *relationTb.GroupMemberModel) (string, bool) {
		return e.UserID, e.Nickname == ""
	}), true)
	if err != nil {
		return nil, err
	}
	resp.Members = utils.Slice(members, func(e *relationTb.GroupMemberModel) *sdkws.GroupMemberFullInfo {
		if e.Nickname == "" {
			e.Nickname = nameMap[e.UserID]
		}
		return DbToPbGroupMembersCMSResp(e)
	})
	return resp, nil
}

func (s *groupServer) GetGroupApplicationList(ctx context.Context, req *pbGroup.GetGroupApplicationListReq) (*pbGroup.GetGroupApplicationListResp, error) {
	resp := &pbGroup.GetGroupApplicationListResp{}
	total, groupRequests, err := s.GroupInterface.PageGroupRequestUser(ctx, req.FromUserID, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	if err != nil {
		return nil, err
	}
	resp.Total = total
	if len(groupRequests) == 0 {
		return resp, nil
	}
	var (
		userIDs  []string
		groupIDs []string
	)
	for _, gr := range groupRequests {
		userIDs = append(userIDs, gr.UserID)
		groupIDs = append(groupIDs, gr.GroupID)
	}
	userIDs = utils.Distinct(userIDs)
	groupIDs = utils.Distinct(groupIDs)
	userMap, err := s.UserCheck.GetPublicUserInfoMap(ctx, userIDs, true)
	if err != nil {
		return nil, err
	}
	groups, err := s.GroupInterface.FindGroup(ctx, utils.Distinct(groupIDs))
	if err != nil {
		return nil, err
	}
	groupMap := utils.SliceToMap(groups, func(e *relationTb.GroupModel) string {
		return e.GroupID
	})
	if ids := utils.Single(utils.Keys(groupMap), groupIDs); len(ids) > 0 {
		return nil, constant.ErrGroupIDNotFound.Wrap(strings.Join(ids, ","))
	}
	groupMemberNumMap, err := s.GroupInterface.MapGroupMemberNum(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	owners, err := s.GroupInterface.FindGroupMember(ctx, groupIDs, nil, []int32{constant.GroupOwner})
	if err != nil {
		return nil, err
	}
	ownerMap := utils.SliceToMap(owners, func(e *relationTb.GroupMemberModel) string {
		return e.GroupID
	})
	resp.GroupRequests = utils.Slice(groupRequests, func(e *relationTb.GroupRequestModel) *sdkws.GroupRequest {
		return DbToPbGroupRequest(e, userMap[e.UserID], DbToPbGroupInfo(groupMap[e.GroupID], ownerMap[e.GroupID].UserID, uint32(groupMemberNumMap[e.GroupID])))
	})
	return resp, nil
}

func (s *groupServer) GetGroupsInfo(ctx context.Context, req *pbGroup.GetGroupsInfoReq) (*pbGroup.GetGroupsInfoResp, error) {
	resp := &pbGroup.GetGroupsInfoResp{}
	if len(req.GroupIDs) == 0 {
		return nil, constant.ErrArgs.Wrap("groupID is empty")
	}
	groups, err := s.GroupInterface.FindGroup(ctx, req.GroupIDs)
	if err != nil {
		return nil, err
	}
	groupMemberNumMap, err := s.GroupInterface.MapGroupMemberNum(ctx, req.GroupIDs)
	if err != nil {
		return nil, err
	}
	owners, err := s.GroupInterface.FindGroupMember(ctx, req.GroupIDs, nil, []int32{constant.GroupOwner})
	if err != nil {
		return nil, err
	}
	ownerMap := utils.SliceToMap(owners, func(e *relationTb.GroupMemberModel) string {
		return e.GroupID
	})
	resp.GroupInfos = utils.Slice(groups, func(e *relationTb.GroupModel) *sdkws.GroupInfo {
		return DbToPbGroupInfo(e, ownerMap[e.GroupID].UserID, uint32(groupMemberNumMap[e.GroupID]))
	})
	return resp, nil
}

func (s *groupServer) GroupApplicationResponse(ctx context.Context, req *pbGroup.GroupApplicationResponseReq) (*pbGroup.GroupApplicationResponseResp, error) {
	resp := &pbGroup.GroupApplicationResponseResp{}
	if !utils.Contain(req.HandleResult, constant.GroupResponseAgree, constant.GroupResponseRefuse) {
		return nil, constant.ErrArgs.Wrap("HandleResult unknown")
	}
	if !tokenverify.IsAppManagerUid(ctx) {
		groupMember, err := s.GroupInterface.TakeGroupMember(ctx, req.GroupID, req.FromUserID)
		if err != nil {
			return nil, err
		}
		if !(groupMember.RoleLevel == constant.GroupOwner || groupMember.RoleLevel == constant.GroupAdmin) {
			return nil, constant.ErrNoPermission.Wrap("no group owner or admin")
		}
	}
	group, err := s.GroupInterface.TakeGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	groupRequest, err := s.GroupInterface.TakeGroupRequest(ctx, req.GroupID, req.FromUserID)
	if err != nil {
		return nil, err
	}
	if groupRequest.HandleResult != 0 {
		return nil, constant.ErrArgs.Wrap("group request already processed")
	}
	var join bool
	if _, err = s.GroupInterface.TakeGroupMember(ctx, req.GroupID, req.FromUserID); err == nil {
		join = true // 已经在群里了
	} else if !s.IsNotFound(err) {
		return nil, err
	}
	user, err := s.UserCheck.GetPublicUserInfo(ctx, req.FromUserID)
	if err != nil {
		return nil, err
	}
	var member *relationTb.GroupMemberModel
	if (!join) && req.HandleResult == constant.GroupResponseAgree {
		member = &relationTb.GroupMemberModel{
			GroupID:        req.GroupID,
			UserID:         user.UserID,
			Nickname:       user.Nickname,
			FaceURL:        user.FaceURL,
			RoleLevel:      constant.GroupOrdinaryUsers,
			JoinTime:       time.Now(),
			JoinSource:     groupRequest.JoinSource,
			InviterUserID:  groupRequest.InviterUserID,
			OperatorUserID: tracelog.GetOpUserID(ctx),
			Ex:             groupRequest.Ex,
		}
		if err = CallbackBeforeMemberJoinGroup(ctx, member, group.Ex); err != nil {
			return nil, err
		}
	}
	if err := s.GroupInterface.HandlerGroupRequest(ctx, req.GroupID, req.FromUserID, req.HandledMsg, req.HandleResult, member); err != nil {
		return nil, err
	}
	if !join {
		if req.HandleResult == constant.GroupResponseAgree {
			s.Notification.GroupApplicationAcceptedNotification(ctx, req)
			s.Notification.MemberEnterNotification(ctx, req)
		} else if req.HandleResult == constant.GroupResponseRefuse {
			s.Notification.GroupApplicationRejectedNotification(ctx, req)
		}
	}
	return resp, nil
}

func (s *groupServer) JoinGroup(ctx context.Context, req *pbGroup.JoinGroupReq) (*pbGroup.JoinGroupResp, error) {
	resp := &pbGroup.JoinGroupResp{}
	if _, err := s.UserCheck.GetPublicUserInfo(ctx, tracelog.GetOpUserID(ctx)); err != nil {
		return nil, err
	}
	group, err := s.GroupInterface.TakeGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if group.Status == constant.GroupStatusDismissed {
		return nil, constant.ErrDismissedAlready.Wrap()
	}
	if group.NeedVerification == constant.Directly {
		if group.GroupType == constant.SuperGroup {
			return nil, constant.ErrGroupTypeNotSupport.Wrap()
		}
		user, err := s.UserCheck.GetUsersInfo(ctx, tracelog.GetOpUserID(ctx))
		if err != nil {
			return nil, err
		}
		groupMember := PbToDbGroupMember(user)
		groupMember.GroupID = group.GroupID
		groupMember.RoleLevel = constant.GroupOrdinaryUsers
		groupMember.OperatorUserID = tracelog.GetOpUserID(ctx)
		groupMember.JoinSource = constant.JoinByInvitation
		groupMember.InviterUserID = tracelog.GetOpUserID(ctx)
		if err := CallbackBeforeMemberJoinGroup(ctx, groupMember, group.Ex); err != nil {
			return nil, err
		}
		if err := s.GroupInterface.CreateGroup(ctx, nil, []*relationTb.GroupMemberModel{groupMember}); err != nil {
			return nil, err
		}
		s.Notification.MemberEnterDirectlyNotification(ctx, req.GroupID, tracelog.GetOpUserID(ctx), tracelog.GetOperationID(ctx))
		return resp, nil
	}
	groupRequest := relationTb.GroupRequestModel{
		UserID:     tracelog.GetOpUserID(ctx),
		ReqMsg:     req.ReqMessage,
		GroupID:    req.GroupID,
		JoinSource: req.JoinSource,
		ReqTime:    time.Now(),
	}
	if err := s.GroupInterface.CreateGroupRequest(ctx, []*relationTb.GroupRequestModel{&groupRequest}); err != nil {
		return nil, err
	}
	s.Notification.JoinGroupApplicationNotification(ctx, req)
	return resp, nil
}

func (s *groupServer) QuitGroup(ctx context.Context, req *pbGroup.QuitGroupReq) (*pbGroup.QuitGroupResp, error) {
	resp := &pbGroup.QuitGroupResp{}
	group, err := s.GroupInterface.TakeGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if group.GroupType == constant.SuperGroup {
		if err := s.GroupInterface.DeleteSuperGroupMember(ctx, req.GroupID, []string{tracelog.GetOpUserID(ctx)}); err != nil {
			return nil, err
		}
		s.Notification.SuperGroupNotification(ctx, tracelog.GetOpUserID(ctx), tracelog.GetOpUserID(ctx))
	} else {
		_, err := s.GroupInterface.TakeGroupMember(ctx, req.GroupID, tracelog.GetOpUserID(ctx))
		if err != nil {
			return nil, err
		}
		s.Notification.MemberQuitNotification(ctx, req)
	}
	return resp, nil
}

func (s *groupServer) SetGroupInfo(ctx context.Context, req *pbGroup.SetGroupInfoReq) (*pbGroup.SetGroupInfoResp, error) {
	resp := &pbGroup.SetGroupInfoResp{}
	if !tokenverify.IsAppManagerUid(ctx) {
		groupMember, err := s.GroupInterface.TakeGroupMember(ctx, req.GroupInfoForSet.GroupID, tracelog.GetOpUserID(ctx))
		if err != nil {
			return nil, err
		}
		if !(groupMember.RoleLevel == constant.GroupOwner || groupMember.RoleLevel == constant.GroupAdmin) {
			return nil, constant.ErrNoPermission.Wrap("no group owner or admin")
		}
	}
	group, err := s.GroupInterface.TakeGroup(ctx, req.GroupInfoForSet.GroupID)
	if err != nil {
		return nil, err
	}
	if group.Status == constant.GroupStatusDismissed {
		return nil, utils.Wrap(constant.ErrDismissedAlready, "")
	}
	userIDs, err := s.GroupInterface.FindGroupMemberUserID(ctx, group.GroupID)
	if err != nil {
		return nil, err
	}
	data := UpdateGroupInfoMap(req.GroupInfoForSet)
	if len(data) > 0 {
		return resp, nil
	}
	if err := s.GroupInterface.UpdateGroup(ctx, group.GroupID, data); err != nil {
		return nil, err
	}
	group, err = s.GroupInterface.TakeGroup(ctx, req.GroupInfoForSet.GroupID)
	if err != nil {
		return nil, err
	}
	s.Notification.GroupInfoSetNotification(ctx, req.GroupInfoForSet.GroupID, group.GroupName, group.Notification, group.Introduction, group.FaceURL, req.GroupInfoForSet.NeedVerification)
	if req.GroupInfoForSet.Notification != "" {
		args := pbConversation.ModifyConversationFieldReq{
			Conversation: &pbConversation.Conversation{
				OwnerUserID:      tracelog.GetOpUserID(ctx),
				ConversationID:   utils.GetConversationIDBySessionType(group.GroupID, constant.GroupChatType),
				ConversationType: constant.GroupChatType,
				GroupID:          group.GroupID,
			},
			FieldType:  constant.FieldGroupAtType,
			UserIDList: userIDs,
		}
		if err := s.ConversationChecker.ModifyConversationField(ctx, &args); err != nil {
			tracelog.SetCtxWarn(ctx, "ModifyConversationField", err, args)
		}
	}
	return resp, nil
}

func (s *groupServer) TransferGroupOwner(ctx context.Context, req *pbGroup.TransferGroupOwnerReq) (*pbGroup.TransferGroupOwnerResp, error) {
	resp := &pbGroup.TransferGroupOwnerResp{}
	group, err := s.GroupInterface.TakeGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if group.Status == constant.GroupStatusDismissed {
		return nil, utils.Wrap(constant.ErrDismissedAlready, "")
	}
	if req.OldOwnerUserID == req.NewOwnerUserID {
		return nil, constant.ErrArgs.Wrap("OldOwnerUserID == NewOwnerUserID")
	}
	members, err := s.GroupInterface.FindGroupMember(ctx, []string{req.GroupID}, []string{req.OldOwnerUserID, req.NewOwnerUserID}, nil)
	if err != nil {
		return nil, err
	}
	memberMap := utils.SliceToMap(members, func(e *relationTb.GroupMemberModel) string { return e.UserID })
	if ids := utils.Single([]string{req.OldOwnerUserID, req.NewOwnerUserID}, utils.Keys(memberMap)); len(ids) > 0 {
		return nil, constant.ErrArgs.Wrap("user not in group " + strings.Join(ids, ","))
	}
	newOwner := memberMap[req.NewOwnerUserID]
	if newOwner == nil {
		return nil, constant.ErrArgs.Wrap("NewOwnerUser not in group " + req.NewOwnerUserID)
	}
	oldOwner := memberMap[req.OldOwnerUserID]
	if tokenverify.IsAppManagerUid(ctx) {
		if oldOwner == nil {
			oldOwner, err = s.GroupInterface.TakeGroupOwner(ctx, req.OldOwnerUserID)
			if err != nil {
				return nil, err
			}
		}
	} else {
		if oldOwner == nil {
			return nil, constant.ErrArgs.Wrap("OldOwnerUser not in group " + req.NewOwnerUserID)
		}
		if oldOwner.GroupID != tracelog.GetOpUserID(ctx) {
			return nil, constant.ErrNoPermission.Wrap(fmt.Sprintf("user %s no permission transfer group owner", tracelog.GetOpUserID(ctx)))
		}
	}
	if err := s.GroupInterface.TransferGroupOwner(ctx, req.GroupID, req.OldOwnerUserID, req.NewOwnerUserID, newOwner.RoleLevel); err != nil {
		return nil, err
	}
	s.Notification.GroupOwnerTransferredNotification(ctx, req)
	return resp, nil
}

func (s *groupServer) GetGroups(ctx context.Context, req *pbGroup.GetGroupsReq) (*pbGroup.GetGroupsResp, error) {
	resp := &pbGroup.GetGroupsResp{}
	var (
		groups []*relationTb.GroupModel
		err    error
	)
	if req.GroupID != "" {
		groups, err = s.GroupInterface.FindGroup(ctx, []string{req.GroupID})
		resp.Total = uint32(len(groups))
	} else {
		resp.Total, groups, err = s.GroupInterface.SearchGroup(ctx, req.GroupName, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	}
	if err != nil {
		return nil, err
	}
	groupIDs := utils.Slice(groups, func(e *relationTb.GroupModel) string {
		return e.GroupID
	})
	ownerMembers, err := s.GroupInterface.FindGroupMember(ctx, groupIDs, nil, []int32{constant.GroupOwner})
	if err != nil {
		return nil, err
	}
	ownerMemberMap := utils.SliceToMap(ownerMembers, func(e *relationTb.GroupMemberModel) string {
		return e.GroupID
	})
	if ids := utils.Single(groupIDs, utils.Keys(ownerMemberMap)); len(ids) > 0 {
		return nil, constant.ErrDB.Wrap("group not owner " + strings.Join(ids, ","))
	}
	groupMemberNumMap, err := s.GroupInterface.MapGroupMemberNum(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	resp.Groups = utils.Slice(groups, func(group *relationTb.GroupModel) *pbGroup.CMSGroup {
		member := ownerMemberMap[group.GroupID]
		return DbToPbCMSGroup(group, member.UserID, member.Nickname, uint32(groupMemberNumMap[group.GroupID]))
	})
	return resp, nil
}

func (s *groupServer) GetGroupMembersCMS(ctx context.Context, req *pbGroup.GetGroupMembersCMSReq) (*pbGroup.GetGroupMembersCMSResp, error) {
	resp := &pbGroup.GetGroupMembersCMSResp{}
	total, members, err := s.GroupInterface.SearchGroupMember(ctx, req.UserName, []string{req.GroupID}, nil, nil, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	if err != nil {
		return nil, err
	}
	resp.Total = total
	nameMap, err := s.GetUsernameMap(ctx, utils.Filter(members, func(e *relationTb.GroupMemberModel) (string, bool) {
		return e.UserID, e.Nickname == ""
	}), true)
	if err != nil {
		return nil, err
	}
	resp.Members = utils.Slice(members, func(e *relationTb.GroupMemberModel) *sdkws.GroupMemberFullInfo {
		if e.Nickname == "" {
			e.Nickname = nameMap[e.UserID]
		}
		return DbToPbGroupMembersCMSResp(e)
	})
	return resp, nil
}

func (s *groupServer) GetUserReqApplicationList(ctx context.Context, req *pbGroup.GetUserReqApplicationListReq) (*pbGroup.GetUserReqApplicationListResp, error) {
	resp := &pbGroup.GetUserReqApplicationListResp{}
	user, err := s.UserCheck.GetPublicUserInfo(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	total, requests, err := s.GroupInterface.PageGroupRequestUser(ctx, req.UserID, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	if err != nil {
		return nil, err
	}
	resp.Total = total
	if len(requests) == 0 {
		return resp, nil
	}
	groupIDs := utils.Distinct(utils.Slice(requests, func(e *relationTb.GroupRequestModel) string {
		return e.GroupID
	}))
	groups, err := s.GroupInterface.FindGroup(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	groupMap := utils.SliceToMap(groups, func(e *relationTb.GroupModel) string {
		return e.GroupID
	})
	if ids := utils.Single(groupIDs, utils.Keys(groupMap)); len(ids) > 0 {
		return nil, constant.ErrGroupIDNotFound.Wrap(strings.Join(ids, ","))
	}
	owners, err := s.GroupInterface.FindGroupMember(ctx, groupIDs, nil, []int32{constant.GroupOwner})
	if err != nil {
		return nil, err
	}
	ownerMap := utils.SliceToMap(owners, func(e *relationTb.GroupMemberModel) string {
		return e.GroupID
	})
	if ids := utils.Single(groupIDs, utils.Keys(ownerMap)); len(ids) > 0 {
		return nil, constant.ErrData.Wrap("group no owner", strings.Join(ids, ","))
	}
	groupMemberNum, err := s.GroupInterface.MapGroupMemberNum(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	resp.GroupRequests = utils.Slice(requests, func(e *relationTb.GroupRequestModel) *sdkws.GroupRequest {
		return DbToPbGroupRequest(e, user, DbToPbGroupInfo(groupMap[e.GroupID], ownerMap[e.GroupID].UserID, uint32(groupMemberNum[e.GroupID])))
	})
	return resp, nil
}

func (s *groupServer) DismissGroup(ctx context.Context, req *pbGroup.DismissGroupReq) (*pbGroup.DismissGroupResp, error) {
	resp := &pbGroup.DismissGroupResp{}
	if err := s.CheckGroupAdmin(ctx, req.GroupID); err != nil {
		return nil, err
	}
	group, err := s.GroupInterface.TakeGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if group.Status == constant.GroupStatusDismissed {
		return nil, constant.ErrArgs.Wrap("group status is dismissed")
	}
	if err := s.GroupInterface.DismissGroup(ctx, req.GroupID); err != nil {
		return nil, err
	}
	if group.GroupType == constant.SuperGroup {
		if err := s.GroupInterface.DeleteSuperGroup(ctx, group.GroupID); err != nil {
			return nil, err
		}
	} else {
		s.Notification.GroupDismissedNotification(ctx, req)
	}
	return resp, nil
}

func (s *groupServer) MuteGroupMember(ctx context.Context, req *pbGroup.MuteGroupMemberReq) (*pbGroup.MuteGroupMemberResp, error) {
	resp := &pbGroup.MuteGroupMemberResp{}
	member, err := s.GroupInterface.TakeGroupMember(ctx, req.GroupID, req.UserID)
	if err != nil {
		return nil, err
	}
	if !(tracelog.GetOpUserID(ctx) == req.UserID || tokenverify.IsAppManagerUid(ctx)) {
		opMember, err := s.GroupInterface.TakeGroupMember(ctx, req.GroupID, req.UserID)
		if err != nil {
			return nil, err
		}
		if opMember.RoleLevel <= member.RoleLevel {
			return nil, constant.ErrNoPermission.Wrap(fmt.Sprintf("self RoleLevel %d target %d", opMember.RoleLevel, member.RoleLevel))
		}
	}
	data := UpdateGroupMemberMutedTimeMap(time.Now().Add(time.Second * time.Duration(req.MutedSeconds)))
	if err := s.GroupInterface.UpdateGroupMember(ctx, member.GroupID, member.UserID, data); err != nil {
		return nil, err
	}
	s.Notification.GroupMemberMutedNotification(ctx, req.GroupID, req.UserID, req.MutedSeconds)
	return resp, nil
}

func (s *groupServer) CancelMuteGroupMember(ctx context.Context, req *pbGroup.CancelMuteGroupMemberReq) (*pbGroup.CancelMuteGroupMemberResp, error) {
	resp := &pbGroup.CancelMuteGroupMemberResp{}
	member, err := s.GroupInterface.TakeGroupMember(ctx, req.GroupID, req.UserID)
	if err != nil {
		return nil, err
	}
	if !(tracelog.GetOpUserID(ctx) == req.UserID || tokenverify.IsAppManagerUid(ctx)) {
		opMember, err := s.GroupInterface.TakeGroupMember(ctx, req.GroupID, tracelog.GetOpUserID(ctx))
		if err != nil {
			return nil, err
		}
		if opMember.RoleLevel <= member.RoleLevel {
			return nil, constant.ErrNoPermission.Wrap(fmt.Sprintf("self RoleLevel %d target %d", opMember.RoleLevel, member.RoleLevel))
		}
	}
	data := UpdateGroupMemberMutedTimeMap(time.Unix(0, 0))
	if err := s.GroupInterface.UpdateGroupMember(ctx, member.GroupID, member.UserID, data); err != nil {
		return nil, err
	}
	s.Notification.GroupMemberCancelMutedNotification(ctx, req.GroupID, req.UserID)
	return resp, nil
}

func (s *groupServer) MuteGroup(ctx context.Context, req *pbGroup.MuteGroupReq) (*pbGroup.MuteGroupResp, error) {
	resp := &pbGroup.MuteGroupResp{}
	if err := s.CheckGroupAdmin(ctx, req.GroupID); err != nil {
		return nil, err
	}
	if err := s.GroupInterface.UpdateGroup(ctx, req.GroupID, UpdateGroupStatusMap(constant.GroupStatusMuted)); err != nil {
		return nil, err
	}
	s.Notification.GroupMutedNotification(ctx, req.GroupID)
	return resp, nil
}

func (s *groupServer) CancelMuteGroup(ctx context.Context, req *pbGroup.CancelMuteGroupReq) (*pbGroup.CancelMuteGroupResp, error) {
	resp := &pbGroup.CancelMuteGroupResp{}
	if err := s.CheckGroupAdmin(ctx, req.GroupID); err != nil {
		return nil, err
	}
	if err := s.GroupInterface.UpdateGroup(ctx, req.GroupID, UpdateGroupStatusMap(constant.GroupOk)); err != nil {
		return nil, err
	}
	s.Notification.GroupCancelMutedNotification(ctx, req.GroupID)
	return resp, nil
}

func (s *groupServer) SetGroupMemberInfo(ctx context.Context, req *pbGroup.SetGroupMemberInfoReq) (*pbGroup.SetGroupMemberInfoResp, error) {
	resp := &pbGroup.SetGroupMemberInfoResp{}
	if len(req.Members) == 0 {
		return nil, constant.ErrArgs.Wrap("members empty")
	}
	duplicateMap := make(map[[2]string]struct{})
	userIDMap := make(map[string]struct{})
	groupIDMap := make(map[string]struct{})
	for _, member := range req.Members {
		key := [...]string{member.GroupID, member.UserID}
		if _, ok := duplicateMap[key]; ok {
			return nil, constant.ErrArgs.Wrap("group user duplicate")
		}
		duplicateMap[key] = struct{}{}
		userIDMap[member.UserID] = struct{}{}
		groupIDMap[member.GroupID] = struct{}{}
	}
	groupIDs := utils.Keys(groupIDMap)
	userIDs := utils.Keys(userIDMap)
	members, err := s.GroupInterface.FindGroupMember(ctx, groupIDs, append(userIDs, tracelog.GetOpUserID(ctx)), nil)
	if err != nil {
		return nil, err
	}
	for _, member := range members {
		delete(duplicateMap, [...]string{member.GroupID, member.UserID})
	}
	if len(duplicateMap) > 0 {
		return nil, constant.ErrArgs.Wrap("group not found" + strings.Join(utils.Slice(utils.Keys(duplicateMap), func(e [2]string) string {
			return fmt.Sprintf("[group: %s user: %s]", e[0], e[1])
		}), ","))
	}
	memberMap := utils.SliceToMap(members, func(e *relationTb.GroupMemberModel) [2]string {
		return [...]string{e.GroupID, e.UserID}
	})
	if !tokenverify.IsAppManagerUid(ctx) {
		opUserID := tracelog.GetOpUserID(ctx)
		for _, member := range members {
			if member.UserID == opUserID {
				continue
			}
			opMember, ok := memberMap[[...]string{member.GroupID, member.UserID}]
			if !ok {
				return nil, constant.ErrArgs.Wrap(fmt.Sprintf("user %s not in group %s", opUserID, member.GroupID))
			}
			if member.RoleLevel >= opMember.RoleLevel {
				return nil, constant.ErrNoPermission.Wrap(fmt.Sprintf("group %s : %s RoleLevel %d >= %s RoleLevel %d", member.GroupID, member.UserID, member.RoleLevel, opMember.UserID, opMember.RoleLevel))
			}
		}
	}
	for _, member := range req.Members {
		if member.RoleLevel == nil {
			continue
		}
		if memberMap[[...]string{member.GroupID, member.UserID}].RoleLevel == constant.GroupOwner {
			return nil, constant.ErrArgs.Wrap(fmt.Sprintf("group %s user %s is owner", member.GroupID, member.UserID))
		}
	}
	for i := 0; i < len(req.Members); i++ {
		if err := CallbackBeforeSetGroupMemberInfo(ctx, req.Members[i]); err != nil {
			return nil, err
		}
	}
	err = s.GroupInterface.UpdateGroupMembers(ctx, utils.Slice(req.Members, func(e *pbGroup.SetGroupMemberInfo) *relationTb.BatchUpdateGroupMember {
		return &relationTb.BatchUpdateGroupMember{
			GroupID: e.GroupID,
			UserID:  e.UserID,
			Map:     UpdateGroupMemberMap(e),
		}
	}))
	if err != nil {
		return nil, err
	}
	for _, member := range req.Members {
		s.Notification.GroupMemberInfoSetNotification(ctx, member.GroupID, member.UserID)
	}
	return resp, nil
}

func (s *groupServer) GetGroupAbstractInfo(ctx context.Context, req *pbGroup.GetGroupAbstractInfoReq) (*pbGroup.GetGroupAbstractInfoResp, error) {
	resp := &pbGroup.GetGroupAbstractInfoResp{}
	if len(req.GroupIDs) == 0 {
		return nil, constant.ErrArgs.Wrap("groupIDs empty")
	}
	if utils.Duplicate(req.GroupIDs) {
		return nil, constant.ErrArgs.Wrap("groupIDs duplicate")
	}
	groups, err := s.GroupInterface.FindGroup(ctx, req.GroupIDs)
	if err != nil {
		return nil, err
	}
	if ids := utils.Single(req.GroupIDs, utils.Slice(groups, func(group *relationTb.GroupModel) string {
		return group.GroupID
	})); len(ids) > 0 {
		return nil, constant.ErrGroupIDNotFound.Wrap("not found group " + strings.Join(ids, ","))
	}
	groupUserMap, err := s.GroupInterface.MapGroupMemberUserID(ctx, req.GroupIDs)
	if err != nil {
		return nil, err
	}
	if ids := utils.Single(req.GroupIDs, utils.Keys(groupUserMap)); len(ids) > 0 {
		return nil, constant.ErrGroupIDNotFound.Wrap(fmt.Sprintf("group %s not found member", strings.Join(ids, ",")))
	}
	resp.GroupAbstractInfos = utils.Slice(groups, func(group *relationTb.GroupModel) *pbGroup.GroupAbstractInfo {
		users := groupUserMap[group.GroupID]
		return DbToPbGroupAbstractInfo(group.GroupID, users.MemberNum, users.Hash)
	})
	return resp, nil
}

func (s *groupServer) GetUserInGroupMembers(ctx context.Context, req *pbGroup.GetUserInGroupMembersReq) (*pbGroup.GetUserInGroupMembersResp, error) {
	resp := &pbGroup.GetUserInGroupMembersResp{}
	if len(req.GroupIDs) == 0 {
		return nil, constant.ErrArgs.Wrap("groupIDs empty")
	}
	members, err := s.GroupInterface.FindGroupMember(ctx, []string{req.UserID}, req.GroupIDs, nil)
	if err != nil {
		return nil, err
	}
	nameMap, err := s.GetUsernameMap(ctx, utils.Filter(members, func(e *relationTb.GroupMemberModel) (string, bool) {
		return e.UserID, e.Nickname == ""
	}), true)
	if err != nil {
		return nil, err
	}
	resp.Members = utils.Slice(members, func(e *relationTb.GroupMemberModel) *sdkws.GroupMemberFullInfo {
		if e.Nickname == "" {
			e.Nickname = nameMap[e.UserID]
		}
		return DbToPbGroupMembersCMSResp(e)
	})
	return resp, nil
}

func (s *groupServer) GetGroupMemberUserID(ctx context.Context, req *pbGroup.GetGroupMemberUserIDReq) (*pbGroup.GetGroupMemberUserIDResp, error) {
	resp := &pbGroup.GetGroupMemberUserIDResp{}
	var err error
	resp.UserIDs, err = s.GroupInterface.FindGroupMemberUserID(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *groupServer) GetGroupMemberRoleLevel(ctx context.Context, req *pbGroup.GetGroupMemberRoleLevelReq) (*pbGroup.GetGroupMemberRoleLevelResp, error) {
	resp := &pbGroup.GetGroupMemberRoleLevelResp{}
	if len(req.RoleLevels) == 0 {
		return nil, constant.ErrArgs.Wrap("RoleLevels empty")
	}
	members, err := s.GroupInterface.FindGroupMember(ctx, []string{req.GroupID}, nil, req.RoleLevels)
	if err != nil {
		return nil, err
	}
	nameMap, err := s.GetUsernameMap(ctx, utils.Filter(members, func(e *relationTb.GroupMemberModel) (string, bool) {
		return e.UserID, e.Nickname == ""
	}), true)
	if err != nil {
		return nil, err
	}
	resp.Members = utils.Slice(members, func(e *relationTb.GroupMemberModel) *sdkws.GroupMemberFullInfo {
		if e.Nickname == "" {
			e.Nickname = nameMap[e.UserID]
		}
		return DbToPbGroupMembersCMSResp(e)
	})
	return resp, nil
}
