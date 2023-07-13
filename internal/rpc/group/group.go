// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package group

import (
	"context"
	"fmt"
	pbConversation "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/conversation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/wrapperspb"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient/notification"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/convert"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mw/specialerror"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/controller"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/relation"
	relationTb "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/unrelation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/tokenverify"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	pbGroup "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/group"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"google.golang.org/grpc"
)

func Start(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error {
	db, err := relation.NewGormDB()
	if err != nil {
		return err
	}
	if err := db.AutoMigrate(&relationTb.GroupModel{}, &relationTb.GroupMemberModel{}, &relationTb.GroupRequestModel{}); err != nil {
		return err
	}
	mongo, err := unrelation.NewMongo()
	if err != nil {
		return err
	}
	rdb, err := cache.NewRedis()
	if err != nil {
		return err
	}
	userRpcClient := rpcclient.NewUserRpcClient(client)
	msgRpcClient := rpcclient.NewMessageRpcClient(client)
	conversationRpcClient := rpcclient.NewConversationRpcClient(client)
	database := controller.InitGroupDatabase(db, rdb, mongo.GetDatabase())
	pbGroup.RegisterGroupServer(server, &groupServer{
		GroupDatabase: database,
		User:          userRpcClient,
		Notification: notification.NewGroupNotificationSender(database, &msgRpcClient, &userRpcClient, func(ctx context.Context, userIDs []string) ([]notification.CommonUser, error) {
			users, err := userRpcClient.GetUsersInfo(ctx, userIDs)
			if err != nil {
				return nil, err
			}
			return utils.Slice(users, func(e *sdkws.UserInfo) notification.CommonUser { return e }), nil
		}),
		conversationRpcClient: conversationRpcClient,
		msgRpcClient:          msgRpcClient,
	})
	return nil
}

type groupServer struct {
	GroupDatabase         controller.GroupDatabase
	User                  rpcclient.UserRpcClient
	Notification          *notification.GroupNotificationSender
	conversationRpcClient rpcclient.ConversationRpcClient
	msgRpcClient          rpcclient.MessageRpcClient
}

func (s *groupServer) CheckGroupAdmin(ctx context.Context, groupID string) error {
	if !tokenverify.IsAppManagerUid(ctx) {
		groupMember, err := s.GroupDatabase.TakeGroupMember(ctx, groupID, mcontext.GetOpUserID(ctx))
		if err != nil {
			return err
		}
		if !(groupMember.RoleLevel == constant.GroupOwner || groupMember.RoleLevel == constant.GroupAdmin) {
			return errs.ErrNoPermission.Wrap("no group owner or admin")
		}
	}
	return nil
}

func (s *groupServer) GetUsernameMap(ctx context.Context, userIDs []string, complete bool) (map[string]string, error) {
	if len(userIDs) == 0 {
		return map[string]string{}, nil
	}
	users, err := s.User.GetPublicUserInfos(ctx, userIDs, complete)
	if err != nil {
		return nil, err
	}
	return utils.SliceToMapAny(users, func(e *sdkws.PublicUserInfo) (string, string) {
		return e.UserID, e.Nickname
	}), nil
}

func (s *groupServer) IsNotFound(err error) bool {
	return errs.ErrRecordNotFound.Is(specialerror.ErrCode(errs.Unwrap(err)))
}

func (s *groupServer) GenGroupID(ctx context.Context, groupID *string) error {
	if *groupID != "" {
		_, err := s.GroupDatabase.TakeGroup(ctx, *groupID)
		if err == nil {
			return errs.ErrGroupIDExisted.Wrap("group id existed " + *groupID)
		} else if s.IsNotFound(err) {
			return nil
		} else {
			return err
		}
	}
	for i := 0; i < 10; i++ {
		id := utils.Md5(strings.Join([]string{mcontext.GetOperationID(ctx), strconv.FormatInt(time.Now().UnixNano(), 10), strconv.Itoa(rand.Int())}, ",;,"))
		bi := big.NewInt(0)
		bi.SetString(id[0:8], 16)
		id = bi.String()
		_, err := s.GroupDatabase.TakeGroup(ctx, id)
		if err == nil {
			continue
		} else if s.IsNotFound(err) {
			*groupID = id
			return nil
		} else {
			return err
		}
	}
	return errs.ErrData.Wrap("group id gen error")
}

func (s *groupServer) CreateGroup(ctx context.Context, req *pbGroup.CreateGroupReq) (*pbGroup.CreateGroupResp, error) {
	if req.OwnerUserID == "" {
		return nil, errs.ErrArgs.Wrap("no group owner")
	}
	if err := tokenverify.CheckAccessV3(ctx, req.OwnerUserID); err != nil {
		return nil, err
	}
	userIDs := append(append(req.MemberUserIDs, req.AdminUserIDs...), req.OwnerUserID)
	opUserID := mcontext.GetOpUserID(ctx)
	if !utils.Contain(opUserID, userIDs...) {
		userIDs = append(userIDs, opUserID)
	}
	if utils.Duplicate(userIDs) {
		return nil, errs.ErrArgs.Wrap("group member repeated")
	}
	userMap, err := s.User.GetUsersInfoMap(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	if len(userMap) != len(userIDs) {
		return nil, errs.ErrUserIDNotFound.Wrap("user not found")
	}
	if err := CallbackBeforeCreateGroup(ctx, req); err != nil {
		return nil, err
	}
	var groupMembers []*relationTb.GroupMemberModel
	group := convert.Pb2DBGroupInfo(req.GroupInfo)
	if err := s.GenGroupID(ctx, &group.GroupID); err != nil {
		return nil, err
	}
	joinGroup := func(userID string, roleLevel int32) error {
		groupMember := convert.Pb2DbGroupMember(userMap[userID])
		groupMember.Nickname = ""
		groupMember.GroupID = group.GroupID
		groupMember.RoleLevel = roleLevel
		groupMember.OperatorUserID = mcontext.GetOpUserID(ctx)
		groupMember.JoinSource = constant.JoinByInvitation
		groupMember.InviterUserID = mcontext.GetOpUserID(ctx)
		groupMember.JoinTime = time.Now()
		groupMember.MuteEndTime = time.Unix(0, 0)
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
		if err := s.GroupDatabase.CreateSuperGroup(ctx, group.GroupID, userIDs); err != nil {
			return nil, err
		}
	} else {
		for _, userID := range req.AdminUserIDs {
			if err := joinGroup(userID, constant.GroupAdmin); err != nil {
				return nil, err
			}
		}
		for _, userID := range req.MemberUserIDs {
			if err := joinGroup(userID, constant.GroupOrdinaryUsers); err != nil {
				return nil, err
			}
		}
	}
	if err := s.GroupDatabase.CreateGroup(ctx, []*relationTb.GroupModel{group}, groupMembers); err != nil {
		return nil, err
	}
	resp := &pbGroup.CreateGroupResp{GroupInfo: &sdkws.GroupInfo{}}
	resp.GroupInfo = convert.Db2PbGroupInfo(group, req.OwnerUserID, uint32(len(userIDs)))
	resp.GroupInfo.MemberCount = uint32(len(userIDs))
	if req.GroupInfo.GroupType == constant.SuperGroup {
		go func() {
			for _, userID := range userIDs {
				s.Notification.SuperGroupNotification(ctx, userID, userID)
			}
		}()
	} else {
		//s.Notification.GroupCreatedNotification(ctx, group, groupMembers, userMap)
		tips := &sdkws.GroupCreatedTips{
			Group:          resp.GroupInfo,
			OperationTime:  group.CreateTime.UnixMilli(),
			GroupOwnerUser: s.groupMemberDB2PB(groupMembers[0], userMap[groupMembers[0].UserID].AppMangerLevel),
		}
		for _, member := range groupMembers {
			member.Nickname = userMap[member.UserID].Nickname
			tips.MemberList = append(tips.MemberList, s.groupMemberDB2PB(member, userMap[member.UserID].AppMangerLevel))
			if member.UserID == opUserID {
				tips.OpUser = s.groupMemberDB2PB(member, userMap[member.UserID].AppMangerLevel)
				break
			}
		}
		s.Notification.GroupCreatedNotification(ctx, tips)
	}
	return resp, nil
}

func (s *groupServer) GetJoinedGroupList(ctx context.Context, req *pbGroup.GetJoinedGroupListReq) (*pbGroup.GetJoinedGroupListResp, error) {
	resp := &pbGroup.GetJoinedGroupListResp{}
	if err := tokenverify.CheckAccessV3(ctx, req.FromUserID); err != nil {
		return nil, err
	}
	var pageNumber, showNumber int32
	if req.Pagination != nil {
		pageNumber = req.Pagination.PageNumber
		showNumber = req.Pagination.ShowNumber
	}
	//total, members, err := s.GroupDatabase.PageGroupMember(ctx, nil, []string{req.FromUserID}, nil, pageNumber, showNumber)
	total, members, err := s.GroupDatabase.PageGetJoinGroup(ctx, req.FromUserID, pageNumber, showNumber)
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
	groups, err := s.GroupDatabase.FindGroup(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	groupMemberNum, err := s.GroupDatabase.MapGroupMemberNum(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	owners, err := s.FindGroupMember(ctx, groupIDs, nil, []int32{constant.GroupOwner})
	if err != nil {
		return nil, err
	}
	ownerMap := utils.SliceToMap(owners, func(e *relationTb.GroupMemberModel) string {
		return e.GroupID
	})
	resp.Groups = utils.Slice(utils.Order(groupIDs, groups, func(group *relationTb.GroupModel) string {
		return group.GroupID
	}), func(group *relationTb.GroupModel) *sdkws.GroupInfo {
		var userID string
		if user := ownerMap[group.GroupID]; user != nil {
			userID = user.UserID
		}
		return convert.Db2PbGroupInfo(group, userID, groupMemberNum[group.GroupID])
	})
	return resp, nil
}

func (s *groupServer) InviteUserToGroup(ctx context.Context, req *pbGroup.InviteUserToGroupReq) (*pbGroup.InviteUserToGroupResp, error) {
	resp := &pbGroup.InviteUserToGroupResp{}
	if len(req.InvitedUserIDs) == 0 {
		return nil, errs.ErrArgs.Wrap("user empty")
	}
	if utils.Duplicate(req.InvitedUserIDs) {
		return nil, errs.ErrArgs.Wrap("userID duplicate")
	}
	group, err := s.GroupDatabase.TakeGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if group.Status == constant.GroupStatusDismissed {
		return nil, errs.ErrDismissedAlready.Wrap()
	}
	userMap, err := s.User.GetUsersInfoMap(ctx, req.InvitedUserIDs)
	if err != nil {
		return nil, err
	}
	var groupMember *relationTb.GroupMemberModel
	var opUserID string
	if !tokenverify.IsAppManagerUid(ctx) {
		opUserID = mcontext.GetOpUserID(ctx)
		groupMembers, err := s.FindGroupMember(ctx, []string{req.GroupID}, []string{opUserID}, nil)
		if err != nil {
			return nil, err
		}
		if len(groupMembers) <= 0 {
			return nil, errs.ErrNoPermission.Wrap("not in group")
		}
		groupMember = groupMembers[0]
	}
	if group.NeedVerification == constant.AllNeedVerification {
		if !tokenverify.IsAppManagerUid(ctx) {
			if !(groupMember.RoleLevel == constant.GroupOwner || groupMember.RoleLevel == constant.GroupAdmin) {
				var requests []*relationTb.GroupRequestModel
				for _, userID := range req.InvitedUserIDs {
					requests = append(requests, &relationTb.GroupRequestModel{
						UserID:        userID,
						GroupID:       req.GroupID,
						JoinSource:    constant.JoinByInvitation,
						InviterUserID: opUserID,
						ReqTime:       time.Now(),
						HandledTime:   time.Unix(0, 0),
					})
				}
				if err := s.GroupDatabase.CreateGroupRequest(ctx, requests); err != nil {
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
		if err := s.GroupDatabase.CreateSuperGroupMember(ctx, req.GroupID, req.InvitedUserIDs); err != nil {
			return nil, err
		}
		if err := s.conversationRpcClient.GroupChatFirstCreateConversation(ctx, req.GroupID, req.InvitedUserIDs); err != nil {
			return nil, err
		}
		for _, userID := range req.InvitedUserIDs {
			s.Notification.SuperGroupNotification(ctx, userID, userID)
		}
	} else {
		opUserID := mcontext.GetOpUserID(ctx)
		var groupMembers []*relationTb.GroupMemberModel
		for _, userID := range req.InvitedUserIDs {
			member := convert.Pb2DbGroupMember(userMap[userID])
			member.Nickname = ""
			member.GroupID = req.GroupID
			member.RoleLevel = constant.GroupOrdinaryUsers
			member.OperatorUserID = opUserID
			member.InviterUserID = opUserID
			member.JoinSource = constant.JoinByInvitation
			member.JoinTime = time.Now()
			member.MuteEndTime = time.Unix(0, 0)
			if err := CallbackBeforeMemberJoinGroup(ctx, member, group.Ex); err != nil {
				return nil, err
			}
			groupMembers = append(groupMembers, member)
		}
		if err := s.GroupDatabase.CreateGroup(ctx, nil, groupMembers); err != nil {
			return nil, err
		}
		if err := s.conversationRpcClient.GroupChatFirstCreateConversation(ctx, req.GroupID, req.InvitedUserIDs); err != nil {
			return nil, err
		}
		s.Notification.MemberInvitedNotification(ctx, req.GroupID, req.Reason, req.InvitedUserIDs)
	}
	return resp, nil
}

func (s *groupServer) GetGroupAllMember(ctx context.Context, req *pbGroup.GetGroupAllMemberReq) (*pbGroup.GetGroupAllMemberResp, error) {
	resp := &pbGroup.GetGroupAllMemberResp{}
	group, err := s.GroupDatabase.TakeGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if group.GroupType == constant.SuperGroup {
		return nil, errs.ErrArgs.Wrap("unsupported super group")
	}
	members, err := s.FindGroupMember(ctx, []string{req.GroupID}, nil, nil)
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
		return convert.Db2PbGroupMember(e)
	})
	return resp, nil
}

func (s *groupServer) GetGroupMemberList(ctx context.Context, req *pbGroup.GetGroupMemberListReq) (*pbGroup.GetGroupMemberListResp, error) {
	resp := &pbGroup.GetGroupMemberListResp{}
	total, members, err := s.PageGetGroupMember(ctx, req.GroupID, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	log.ZDebug(ctx, "GetGroupMemberList", "total", total, "members", members, "length", len(members))
	if err != nil {
		return nil, err
	}
	resp.Total = total
	resp.Members = utils.Batch(convert.Db2PbGroupMember, members)
	log.ZDebug(ctx, "GetGroupMemberList", "resp", resp, "length", len(resp.Members))
	return resp, nil
}

func (s *groupServer) KickGroupMember(ctx context.Context, req *pbGroup.KickGroupMemberReq) (*pbGroup.KickGroupMemberResp, error) {
	resp := &pbGroup.KickGroupMemberResp{}
	group, err := s.GroupDatabase.TakeGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if len(req.KickedUserIDs) == 0 {
		return nil, errs.ErrArgs.Wrap("KickedUserIDs empty")
	}
	if utils.IsDuplicateStringSlice(req.KickedUserIDs) {
		return nil, errs.ErrArgs.Wrap("KickedUserIDs duplicate")
	}
	opUserID := mcontext.GetOpUserID(ctx)
	if utils.IsContain(opUserID, req.KickedUserIDs) {
		return nil, errs.ErrArgs.Wrap("opUserID in KickedUserIDs")
	}
	if group.GroupType == constant.SuperGroup {
		if err := s.GroupDatabase.DeleteSuperGroupMember(ctx, req.GroupID, req.KickedUserIDs); err != nil {
			return nil, err
		}
		go func() {
			for _, userID := range req.KickedUserIDs {
				s.Notification.SuperGroupNotification(ctx, userID, userID)
			}
		}()
	} else {
		members, err := s.FindGroupMember(ctx, []string{req.GroupID}, append(req.KickedUserIDs, opUserID), nil)
		if err != nil {
			return nil, err
		}
		memberMap := make(map[string]*relationTb.GroupMemberModel)
		for i, member := range members {
			memberMap[member.UserID] = members[i]
		}
		isAppManagerUid := tokenverify.IsAppManagerUid(ctx)
		opMember := memberMap[opUserID]
		for _, userID := range req.KickedUserIDs {
			member, ok := memberMap[userID]
			if !ok {
				return nil, errs.ErrUserIDNotFound.Wrap(userID)
			}
			if !isAppManagerUid {
				if opMember == nil {
					return nil, errs.ErrNoPermission.Wrap("opUserID no in group")
				}
				switch opMember.RoleLevel {
				case constant.GroupOwner:
				case constant.GroupAdmin:
					if member.RoleLevel == constant.GroupOwner || member.RoleLevel == constant.GroupAdmin {
						return nil, errs.ErrNoPermission.Wrap("group admins cannot remove the group owner and other admins")
					}
				case constant.GroupOrdinaryUsers:
					return nil, errs.ErrNoPermission.Wrap("opUserID no permission")
				default:
					return nil, errs.ErrNoPermission.Wrap("opUserID roleLevel unknown")
				}
			}
		}
		num, err := s.GroupDatabase.FindGroupMemberNum(ctx, req.GroupID)
		if err != nil {
			return nil, err
		}
		owner, err := s.FindGroupMember(ctx, []string{req.GroupID}, nil, []int32{constant.GroupOwner})
		if err != nil {
			return nil, err
		}
		if err := s.GroupDatabase.DeleteGroupMember(ctx, group.GroupID, req.KickedUserIDs); err != nil {
			return nil, err
		}
		tips := &sdkws.MemberKickedTips{
			Group: &sdkws.GroupInfo{
				GroupID:      group.GroupID,
				GroupName:    group.GroupName,
				Notification: group.Notification,
				Introduction: group.Introduction,
				FaceURL:      group.FaceURL,
				//OwnerUserID:            owner[0].UserID,
				CreateTime:             group.CreateTime.UnixMilli(),
				MemberCount:            num,
				Ex:                     group.Ex,
				Status:                 group.Status,
				CreatorUserID:          group.CreatorUserID,
				GroupType:              group.GroupType,
				NeedVerification:       group.NeedVerification,
				LookMemberInfo:         group.LookMemberInfo,
				ApplyMemberFriend:      group.ApplyMemberFriend,
				NotificationUpdateTime: group.NotificationUpdateTime.UnixMilli(),
				NotificationUserID:     group.NotificationUserID,
			},
			KickedUserList: []*sdkws.GroupMemberFullInfo{},
		}
		if len(owner) > 0 {
			tips.Group.OwnerUserID = owner[0].UserID
		}
		if opMember, ok := memberMap[opUserID]; ok {
			tips.OpUser = convert.Db2PbGroupMember(opMember)
		}
		for _, userID := range req.KickedUserIDs {
			tips.KickedUserList = append(tips.KickedUserList, convert.Db2PbGroupMember(memberMap[userID]))
		}
		s.Notification.MemberKickedNotification(ctx, tips)
	}
	if err := s.deleteMemberAndSetConversationSeq(ctx, req.GroupID, req.KickedUserIDs); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *groupServer) GetGroupMembersInfo(ctx context.Context, req *pbGroup.GetGroupMembersInfoReq) (*pbGroup.GetGroupMembersInfoResp, error) {
	resp := &pbGroup.GetGroupMembersInfoResp{}
	if len(req.UserIDs) == 0 {
		return nil, errs.ErrArgs.Wrap("userIDs empty")
	}
	if req.GroupID == "" {
		return nil, errs.ErrArgs.Wrap("groupID empty")
	}
	members, err := s.FindGroupMember(ctx, []string{req.GroupID}, req.UserIDs, nil)
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
		return convert.Db2PbGroupMember(e)
	})
	return resp, nil
}

func (s *groupServer) GetGroupApplicationList(ctx context.Context, req *pbGroup.GetGroupApplicationListReq) (*pbGroup.GetGroupApplicationListResp, error) {
	pageNumber, showNumber := utils.GetPage(req.Pagination)

	groupIDs, err := s.GroupDatabase.FindUserManagedGroupID(ctx, req.FromUserID)
	if err != nil {
		return nil, err
	}
	resp := &pbGroup.GetGroupApplicationListResp{}
	if len(groupIDs) == 0 {
		return resp, nil
	}
	total, groupRequests, err := s.GroupDatabase.PageGroupRequest(ctx, groupIDs, pageNumber, showNumber)
	if err != nil {
		return nil, err
	}
	resp.Total = total
	if len(groupRequests) == 0 {
		return resp, nil
	}
	var userIDs []string

	for _, gr := range groupRequests {
		userIDs = append(userIDs, gr.UserID)
	}
	userIDs = utils.Distinct(userIDs)
	userMap, err := s.User.GetPublicUserInfoMap(ctx, userIDs, true)
	if err != nil {
		return nil, err
	}
	groups, err := s.GroupDatabase.FindGroup(ctx, utils.Distinct(groupIDs))
	if err != nil {
		return nil, err
	}
	groupMap := utils.SliceToMap(groups, func(e *relationTb.GroupModel) string {
		return e.GroupID
	})
	if ids := utils.Single(utils.Keys(groupMap), groupIDs); len(ids) > 0 {
		return nil, errs.ErrGroupIDNotFound.Wrap(strings.Join(ids, ","))
	}
	groupMemberNumMap, err := s.GroupDatabase.MapGroupMemberNum(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	owners, err := s.FindGroupMember(ctx, groupIDs, nil, []int32{constant.GroupOwner})
	if err != nil {
		return nil, err
	}
	ownerMap := utils.SliceToMap(owners, func(e *relationTb.GroupMemberModel) string {
		return e.GroupID
	})
	resp.GroupRequests = utils.Slice(groupRequests, func(e *relationTb.GroupRequestModel) *sdkws.GroupRequest {
		return convert.Db2PbGroupRequest(e, userMap[e.UserID], convert.Db2PbGroupInfo(groupMap[e.GroupID], ownerMap[e.GroupID].UserID, groupMemberNumMap[e.GroupID]))
	})
	return resp, nil
}

func (s *groupServer) GetGroupsInfo(ctx context.Context, req *pbGroup.GetGroupsInfoReq) (*pbGroup.GetGroupsInfoResp, error) {
	resp := &pbGroup.GetGroupsInfoResp{}
	if len(req.GroupIDs) == 0 {
		return nil, errs.ErrArgs.Wrap("groupID is empty")
	}
	groups, err := s.GroupDatabase.FindGroup(ctx, req.GroupIDs)
	if err != nil {
		return nil, err
	}
	groupMemberNumMap, err := s.GroupDatabase.MapGroupMemberNum(ctx, req.GroupIDs)
	if err != nil {
		return nil, err
	}
	owners, err := s.FindGroupMember(ctx, req.GroupIDs, nil, []int32{constant.GroupOwner})
	if err != nil {
		return nil, err
	}
	ownerMap := utils.SliceToMap(owners, func(e *relationTb.GroupMemberModel) string {
		return e.GroupID
	})
	resp.GroupInfos = utils.Slice(groups, func(e *relationTb.GroupModel) *sdkws.GroupInfo {
		var ownerUserID string
		if owner, ok := ownerMap[e.GroupID]; ok {
			ownerUserID = owner.UserID
		}
		return convert.Db2PbGroupInfo(e, ownerUserID, groupMemberNumMap[e.GroupID])
	})
	return resp, nil
}

func (s *groupServer) GroupApplicationResponse(ctx context.Context, req *pbGroup.GroupApplicationResponseReq) (*pbGroup.GroupApplicationResponseResp, error) {
	defer log.ZInfo(ctx, utils.GetFuncName()+" Return")
	if !utils.Contain(req.HandleResult, constant.GroupResponseAgree, constant.GroupResponseRefuse) {
		return nil, errs.ErrArgs.Wrap("HandleResult unknown")
	}
	if !tokenverify.IsAppManagerUid(ctx) {
		groupMember, err := s.GroupDatabase.TakeGroupMember(ctx, req.GroupID, mcontext.GetOpUserID(ctx))
		if err != nil {
			return nil, err
		}
		if !(groupMember.RoleLevel == constant.GroupOwner || groupMember.RoleLevel == constant.GroupAdmin) {
			return nil, errs.ErrNoPermission.Wrap("no group owner or admin")
		}
	}
	group, err := s.GroupDatabase.TakeGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	groupRequest, err := s.GroupDatabase.TakeGroupRequest(ctx, req.GroupID, req.FromUserID)
	if err != nil {
		return nil, err
	}
	if groupRequest.HandleResult != 0 {
		return nil, errs.ErrGroupRequestHandled.Wrap("group request already processed")
	}
	var inGroup bool
	_, err = s.GroupDatabase.TakeGroupMember(ctx, req.GroupID, req.FromUserID)
	if err == nil {
		inGroup = true // 已经在群里了
	} else if !s.IsNotFound(err) {
		return nil, err
	}
	if _, err := s.User.GetPublicUserInfo(ctx, req.FromUserID); err != nil {
		return nil, err
	}
	var member *relationTb.GroupMemberModel
	if (!inGroup) && req.HandleResult == constant.GroupResponseAgree {
		member = &relationTb.GroupMemberModel{
			GroupID:        req.GroupID,
			UserID:         req.FromUserID,
			Nickname:       "",
			FaceURL:        "",
			RoleLevel:      constant.GroupOrdinaryUsers,
			JoinTime:       time.Now(),
			JoinSource:     groupRequest.JoinSource,
			MuteEndTime:    time.Unix(0, 0),
			InviterUserID:  groupRequest.InviterUserID,
			OperatorUserID: mcontext.GetOpUserID(ctx),
			Ex:             groupRequest.Ex,
		}
		if err = CallbackBeforeMemberJoinGroup(ctx, member, group.Ex); err != nil {
			return nil, err
		}
	}
	if err := s.GroupDatabase.HandlerGroupRequest(ctx, req.GroupID, req.FromUserID, req.HandledMsg, req.HandleResult, member); err != nil {
		return nil, err
	}
	if err := s.conversationRpcClient.GroupChatFirstCreateConversation(ctx, req.GroupID, []string{req.FromUserID}); err != nil {
		return nil, err
	}
	switch req.HandleResult {
	case constant.GroupResponseAgree:
		s.Notification.GroupApplicationAcceptedNotification(ctx, req)
	case constant.GroupResponseRefuse:
		s.Notification.GroupApplicationRejectedNotification(ctx, req)
	}
	if member != nil {
		s.Notification.MemberEnterNotification(ctx, req)
	}
	return &pbGroup.GroupApplicationResponseResp{}, nil
}

func (s *groupServer) JoinGroup(ctx context.Context, req *pbGroup.JoinGroupReq) (resp *pbGroup.JoinGroupResp, err error) {
	defer log.ZInfo(ctx, "JoinGroup.Return")
	user, err := s.User.GetUserInfo(ctx, req.InviterUserID)
	if err != nil {
		return nil, err
	}
	group, err := s.GroupDatabase.TakeGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if group.Status == constant.GroupStatusDismissed {
		return nil, errs.ErrDismissedAlready.Wrap()
	}
	_, err = s.GroupDatabase.TakeGroupMember(ctx, req.GroupID, req.InviterUserID)
	if err == nil {
		return nil, errs.ErrArgs.Wrap("already in group")
	} else if !s.IsNotFound(err) && utils.Unwrap(err) != errs.ErrRecordNotFound {
		return nil, err
	}
	log.ZInfo(ctx, "JoinGroup.groupInfo", "group", group, "eq", group.NeedVerification == constant.Directly)
	resp = &pbGroup.JoinGroupResp{}
	if group.NeedVerification == constant.Directly {
		if group.GroupType == constant.SuperGroup {
			return nil, errs.ErrGroupTypeNotSupport.Wrap()
		}
		groupMember := convert.Pb2DbGroupMember(user)
		groupMember.GroupID = group.GroupID
		groupMember.RoleLevel = constant.GroupOrdinaryUsers
		groupMember.OperatorUserID = mcontext.GetOpUserID(ctx)
		groupMember.JoinSource = constant.JoinByInvitation
		groupMember.InviterUserID = req.InviterUserID
		groupMember.JoinTime = time.Now()
		groupMember.MuteEndTime = time.Unix(0, 0)
		if err := CallbackBeforeMemberJoinGroup(ctx, groupMember, group.Ex); err != nil {
			return nil, err
		}
		if err := s.GroupDatabase.CreateGroup(ctx, nil, []*relationTb.GroupMemberModel{groupMember}); err != nil {
			return nil, err
		}
		if err := s.conversationRpcClient.GroupChatFirstCreateConversation(ctx, req.GroupID, []string{req.InviterUserID}); err != nil {
			return nil, err
		}
		s.Notification.MemberEnterDirectlyNotification(ctx, req.GroupID, req.InviterUserID)
		return resp, nil
	}
	groupRequest := relationTb.GroupRequestModel{
		UserID:      req.InviterUserID,
		ReqMsg:      req.ReqMessage,
		GroupID:     req.GroupID,
		JoinSource:  req.JoinSource,
		ReqTime:     time.Now(),
		HandledTime: time.Unix(0, 0),
	}
	if err := s.GroupDatabase.CreateGroupRequest(ctx, []*relationTb.GroupRequestModel{&groupRequest}); err != nil {
		return nil, err
	}
	s.Notification.JoinGroupApplicationNotification(ctx, req)
	return resp, nil
}

func (s *groupServer) QuitGroup(ctx context.Context, req *pbGroup.QuitGroupReq) (*pbGroup.QuitGroupResp, error) {
	resp := &pbGroup.QuitGroupResp{}
	group, err := s.GroupDatabase.TakeGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if group.GroupType == constant.SuperGroup {
		if err := s.GroupDatabase.DeleteSuperGroupMember(ctx, req.GroupID, []string{mcontext.GetOpUserID(ctx)}); err != nil {
			return nil, err
		}
		s.Notification.SuperGroupNotification(ctx, mcontext.GetOpUserID(ctx), mcontext.GetOpUserID(ctx))
	} else {
		info, err := s.TakeGroupMember(ctx, req.GroupID, mcontext.GetOpUserID(ctx))
		if err != nil {
			return nil, err
		}
		if info.RoleLevel == constant.GroupOwner {
			return nil, errs.ErrNoPermission.Wrap("group owner can't quit")
		}
		err = s.GroupDatabase.DeleteGroupMember(ctx, req.GroupID, []string{mcontext.GetOpUserID(ctx)})
		if err != nil {
			return nil, err
		}
		s.Notification.MemberQuitNotification(ctx, s.groupMemberDB2PB(info, 0))
	}
	if err := s.deleteMemberAndSetConversationSeq(ctx, req.GroupID, []string{mcontext.GetOpUserID(ctx)}); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *groupServer) deleteMemberAndSetConversationSeq(ctx context.Context, groupID string, userIDs []string) error {
	conevrsationID := utils.GetConversationIDBySessionType(constant.SuperGroupChatType, groupID)
	maxSeq, err := s.msgRpcClient.GetConversationMaxSeq(ctx, conevrsationID)
	if err != nil {
		return err
	}
	return s.conversationRpcClient.SetConversationMaxSeq(ctx, userIDs, conevrsationID, maxSeq)
}

func (s *groupServer) SetGroupInfo(ctx context.Context, req *pbGroup.SetGroupInfoReq) (*pbGroup.SetGroupInfoResp, error) {
	var opMember *relationTb.GroupMemberModel
	if !tokenverify.IsAppManagerUid(ctx) {
		var err error
		opMember, err = s.TakeGroupMember(ctx, req.GroupInfoForSet.GroupID, mcontext.GetOpUserID(ctx))
		if err != nil {
			return nil, err
		}
		if !(opMember.RoleLevel == constant.GroupOwner || opMember.RoleLevel == constant.GroupAdmin) {
			return nil, errs.ErrNoPermission.Wrap("no group owner or admin")
		}
	}
	group, err := s.GroupDatabase.TakeGroup(ctx, req.GroupInfoForSet.GroupID)
	if err != nil {
		return nil, err
	}
	if group.Status == constant.GroupStatusDismissed {
		return nil, utils.Wrap(errs.ErrDismissedAlready, "")
	}
	resp := &pbGroup.SetGroupInfoResp{}
	count, err := s.GroupDatabase.FindGroupMemberNum(ctx, group.GroupID)
	if err != nil {
		return nil, err
	}
	owner, err := s.TakeGroupOwner(ctx, group.GroupID)
	if err != nil {
		return nil, err
	}
	data := UpdateGroupInfoMap(ctx, req.GroupInfoForSet)
	if len(data) == 0 {
		return resp, nil
	}
	if err := s.GroupDatabase.UpdateGroup(ctx, group.GroupID, data); err != nil {
		return nil, err
	}
	group, err = s.GroupDatabase.TakeGroup(ctx, req.GroupInfoForSet.GroupID)
	if err != nil {
		return nil, err
	}
	tips := &sdkws.GroupInfoSetTips{
		Group:    s.groupDB2PB(group, owner.UserID, count),
		MuteTime: 0,
		OpUser:   &sdkws.GroupMemberFullInfo{},
	}
	if opMember != nil {
		tips.OpUser = s.groupMemberDB2PB(opMember, 0)
	}
	var num int
	if req.GroupInfoForSet.Notification != "" {
		go func() {
			nctx := mcontext.NewCtx("@@@" + mcontext.GetOperationID(ctx))
			conversation := &pbConversation.ConversationReq{
				ConversationID:   utils.GetConversationIDBySessionType(constant.SuperGroupChatType, req.GroupInfoForSet.GroupID),
				ConversationType: constant.SuperGroupChatType,
				GroupID:          req.GroupInfoForSet.GroupID,
			}
			resp, err := s.GetGroupMemberUserIDs(nctx, &pbGroup.GetGroupMemberUserIDsReq{GroupID: req.GroupInfoForSet.GroupID})
			if err != nil {
				log.ZWarn(ctx, "GetGroupMemberIDs", err)
				return
			}
			conversation.GroupAtType = &wrapperspb.Int32Value{Value: constant.GroupNotification}
			if err := s.conversationRpcClient.SetConversations(nctx, resp.UserIDs, conversation); err != nil {
				log.ZWarn(ctx, "SetConversations", err, resp.UserIDs, conversation)
			}
		}()
		num++
		s.Notification.GroupInfoSetAnnouncementNotification(ctx, &sdkws.GroupInfoSetAnnouncementTips{Group: tips.Group, OpUser: tips.OpUser})

	}
	switch len(data) - num {
	case 0:
	case 1:
		if req.GroupInfoForSet.GroupName == "" {
			s.Notification.GroupInfoSetNotification(ctx, tips)
		} else {
			s.Notification.GroupInfoSetNameNotification(ctx, &sdkws.GroupInfoSetNameTips{Group: tips.Group, OpUser: tips.OpUser})
		}
	default:
		s.Notification.GroupInfoSetNotification(ctx, tips)
	}
	return resp, nil
}

func (s *groupServer) TransferGroupOwner(ctx context.Context, req *pbGroup.TransferGroupOwnerReq) (*pbGroup.TransferGroupOwnerResp, error) {
	resp := &pbGroup.TransferGroupOwnerResp{}
	group, err := s.GroupDatabase.TakeGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if group.Status == constant.GroupStatusDismissed {
		return nil, errs.ErrDismissedAlready.Wrap("")
	}
	if req.OldOwnerUserID == req.NewOwnerUserID {
		return nil, errs.ErrArgs.Wrap("OldOwnerUserID == NewOwnerUserID")
	}
	members, err := s.FindGroupMember(ctx, []string{req.GroupID}, []string{req.OldOwnerUserID, req.NewOwnerUserID}, nil)
	if err != nil {
		return nil, err
	}
	memberMap := utils.SliceToMap(members, func(e *relationTb.GroupMemberModel) string { return e.UserID })
	if ids := utils.Single([]string{req.OldOwnerUserID, req.NewOwnerUserID}, utils.Keys(memberMap)); len(ids) > 0 {
		return nil, errs.ErrArgs.Wrap("user not in group " + strings.Join(ids, ","))
	}
	oldOwner := memberMap[req.OldOwnerUserID]
	if oldOwner == nil {
		return nil, errs.ErrArgs.Wrap("OldOwnerUserID not in group " + req.NewOwnerUserID)
	}
	newOwner := memberMap[req.NewOwnerUserID]
	if newOwner == nil {
		return nil, errs.ErrArgs.Wrap("NewOwnerUser not in group " + req.NewOwnerUserID)
	}
	if !tokenverify.IsAppManagerUid(ctx) {
		if !(mcontext.GetOpUserID(ctx) == oldOwner.UserID && oldOwner.RoleLevel == constant.GroupOwner) {
			return nil, errs.ErrNoPermission.Wrap("no permission transfer group owner")
		}
	}
	if err := s.GroupDatabase.TransferGroupOwner(ctx, req.GroupID, req.OldOwnerUserID, req.NewOwnerUserID, newOwner.RoleLevel); err != nil {
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
		groups, err = s.GroupDatabase.FindGroup(ctx, []string{req.GroupID})
		resp.Total = uint32(len(groups))
	} else {
		resp.Total, groups, err = s.GroupDatabase.SearchGroup(ctx, req.GroupName, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	}
	if err != nil {
		return nil, err
	}
	groupIDs := utils.Slice(groups, func(e *relationTb.GroupModel) string {
		return e.GroupID
	})
	ownerMembers, err := s.FindGroupMember(ctx, groupIDs, nil, []int32{constant.GroupOwner})
	if err != nil {
		return nil, err
	}
	ownerMemberMap := utils.SliceToMap(ownerMembers, func(e *relationTb.GroupMemberModel) string {
		return e.GroupID
	})
	if ids := utils.Single(groupIDs, utils.Keys(ownerMemberMap)); len(ids) > 0 {
		return nil, errs.ErrDatabase.Wrap("group not owner " + strings.Join(ids, ","))
	}
	groupMemberNumMap, err := s.GroupDatabase.MapGroupMemberNum(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	resp.Groups = utils.Slice(groups, func(group *relationTb.GroupModel) *pbGroup.CMSGroup {
		member := ownerMemberMap[group.GroupID]
		return convert.Db2PbCMSGroup(group, member.UserID, member.Nickname, uint32(groupMemberNumMap[group.GroupID]))
	})
	return resp, nil
}

func (s *groupServer) GetGroupMembersCMS(ctx context.Context, req *pbGroup.GetGroupMembersCMSReq) (*pbGroup.GetGroupMembersCMSResp, error) {
	resp := &pbGroup.GetGroupMembersCMSResp{}
	total, members, err := s.GroupDatabase.SearchGroupMember(ctx, req.UserName, []string{req.GroupID}, nil, nil, req.Pagination.PageNumber, req.Pagination.ShowNumber)
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
		return convert.Db2PbGroupMember(e)
	})
	return resp, nil
}

func (s *groupServer) GetUserReqApplicationList(ctx context.Context, req *pbGroup.GetUserReqApplicationListReq) (*pbGroup.GetUserReqApplicationListResp, error) {
	resp := &pbGroup.GetUserReqApplicationListResp{}
	user, err := s.User.GetPublicUserInfo(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	var pageNumber, showNumber int32
	if req.Pagination != nil {
		pageNumber = req.Pagination.PageNumber
		showNumber = req.Pagination.ShowNumber
	}
	total, requests, err := s.GroupDatabase.PageGroupRequestUser(ctx, req.UserID, pageNumber, showNumber)
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
	groups, err := s.GroupDatabase.FindGroup(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	groupMap := utils.SliceToMap(groups, func(e *relationTb.GroupModel) string {
		return e.GroupID
	})
	if ids := utils.Single(groupIDs, utils.Keys(groupMap)); len(ids) > 0 {
		return nil, errs.ErrGroupIDNotFound.Wrap(strings.Join(ids, ","))
	}
	owners, err := s.FindGroupMember(ctx, groupIDs, nil, []int32{constant.GroupOwner})
	if err != nil {
		return nil, err
	}
	ownerMap := utils.SliceToMap(owners, func(e *relationTb.GroupMemberModel) string {
		return e.GroupID
	})
	if ids := utils.Single(groupIDs, utils.Keys(ownerMap)); len(ids) > 0 {
		return nil, errs.ErrData.Wrap("group no owner", strings.Join(ids, ","))
	}
	groupMemberNum, err := s.GroupDatabase.MapGroupMemberNum(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	resp.GroupRequests = utils.Slice(requests, func(e *relationTb.GroupRequestModel) *sdkws.GroupRequest {
		return convert.Db2PbGroupRequest(e, user, convert.Db2PbGroupInfo(groupMap[e.GroupID], ownerMap[e.GroupID].UserID, uint32(groupMemberNum[e.GroupID])))
	})
	return resp, nil
}

func (s *groupServer) DismissGroup(ctx context.Context, req *pbGroup.DismissGroupReq) (*pbGroup.DismissGroupResp, error) {
	defer log.ZInfo(ctx, "DismissGroup.return")
	resp := &pbGroup.DismissGroupResp{}
	owner, err := s.TakeGroupOwner(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if !tokenverify.IsAppManagerUid(ctx) {
		if owner.UserID != mcontext.GetOpUserID(ctx) {
			return nil, errs.ErrNoPermission.Wrap("not group owner")
		}
	}
	group, err := s.GroupDatabase.TakeGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if req.DeleteMember == false && group.Status == constant.GroupStatusDismissed {
		return nil, errs.ErrDismissedAlready.Wrap("group status is dismissed")
	}
	//if group.Status == constant.GroupStatusDismissed {
	//	return nil, errs.ErrArgs.Wrap("group status is dismissed")
	//}
	if err := s.GroupDatabase.DismissGroup(ctx, req.GroupID, req.DeleteMember); err != nil {
		return nil, err
	}
	if group.GroupType == constant.SuperGroup {
		if err := s.GroupDatabase.DeleteSuperGroup(ctx, group.GroupID); err != nil {
			return nil, err
		}
	} else {
		if !req.DeleteMember {
			num, err := s.GroupDatabase.FindGroupMemberNum(ctx, req.GroupID)
			if err != nil {
				return nil, err
			}
			//s.Notification.GroupDismissedNotification(ctx, req)
			tips := &sdkws.GroupDismissedTips{
				Group:  s.groupDB2PB(group, owner.UserID, num),
				OpUser: &sdkws.GroupMemberFullInfo{},
			}
			if mcontext.GetOpUserID(ctx) == owner.UserID {
				tips.OpUser = s.groupMemberDB2PB(owner, 0)
			}
			s.Notification.GroupDismissedNotification(ctx, tips)
		}
	}
	return resp, nil
}

func (s *groupServer) MuteGroupMember(ctx context.Context, req *pbGroup.MuteGroupMemberReq) (*pbGroup.MuteGroupMemberResp, error) {
	resp := &pbGroup.MuteGroupMemberResp{}
	//if err := tokenverify.CheckAccessV3(ctx, req.UserID); err != nil {
	//	return nil, err
	//}
	member, err := s.TakeGroupMember(ctx, req.GroupID, req.UserID)
	if err != nil {
		return nil, err
	}
	if !tokenverify.IsAppManagerUid(ctx) {
		opMember, err := s.TakeGroupMember(ctx, req.GroupID, mcontext.GetOpUserID(ctx))
		if err != nil {
			return nil, err
		}
		switch member.RoleLevel {
		case constant.GroupOwner:
			return nil, errs.ErrNoPermission.Wrap("set group owner mute")
		case constant.GroupAdmin:
			if opMember.RoleLevel != constant.GroupOwner {
				return nil, errs.ErrNoPermission.Wrap("set group admin mute")
			}
		case constant.GroupOrdinaryUsers:
			if !(opMember.RoleLevel == constant.GroupAdmin || opMember.RoleLevel == constant.GroupOwner) {
				return nil, errs.ErrNoPermission.Wrap("set group ordinary users mute")
			}
		}
	}
	data := UpdateGroupMemberMutedTimeMap(time.Now().Add(time.Second * time.Duration(req.MutedSeconds)))
	if err := s.GroupDatabase.UpdateGroupMember(ctx, member.GroupID, member.UserID, data); err != nil {
		return nil, err
	}
	s.Notification.GroupMemberMutedNotification(ctx, req.GroupID, req.UserID, req.MutedSeconds)
	return resp, nil
}

func (s *groupServer) CancelMuteGroupMember(ctx context.Context, req *pbGroup.CancelMuteGroupMemberReq) (*pbGroup.CancelMuteGroupMemberResp, error) {
	resp := &pbGroup.CancelMuteGroupMemberResp{}
	//member, err := s.GroupDatabase.TakeGroupMember(ctx, req.GroupID, req.UserID)
	//if err != nil {
	//	return nil, err
	//}
	//if !(mcontext.GetOpUserID(ctx) == req.UserID || tokenverify.IsAppManagerUid(ctx)) {
	//	opMember, err := s.GroupDatabase.TakeGroupMember(ctx, req.GroupID, mcontext.GetOpUserID(ctx))
	//	if err != nil {
	//		return nil, err
	//	}
	//	if opMember.RoleLevel <= member.RoleLevel {
	//		return nil, errs.ErrNoPermission.Wrap(fmt.Sprintf("self RoleLevel %d target %d", opMember.RoleLevel, member.RoleLevel))
	//	}
	//}
	//if err := tokenverify.CheckAccessV3(ctx, req.UserID); err != nil {
	//	return nil, err
	//}
	member, err := s.TakeGroupMember(ctx, req.GroupID, req.UserID)
	if err != nil {
		return nil, err
	}
	if !tokenverify.IsAppManagerUid(ctx) {
		opMember, err := s.TakeGroupMember(ctx, req.GroupID, mcontext.GetOpUserID(ctx))
		if err != nil {
			return nil, err
		}
		switch member.RoleLevel {
		case constant.GroupOwner:
			return nil, errs.ErrNoPermission.Wrap("set group owner mute")
		case constant.GroupAdmin:
			if opMember.RoleLevel != constant.GroupOwner {
				return nil, errs.ErrNoPermission.Wrap("set group admin mute")
			}
		case constant.GroupOrdinaryUsers:
			if !(opMember.RoleLevel == constant.GroupAdmin || opMember.RoleLevel == constant.GroupOwner) {
				return nil, errs.ErrNoPermission.Wrap("set group ordinary users mute")
			}
		}
	}
	data := UpdateGroupMemberMutedTimeMap(time.Unix(0, 0))
	if err := s.GroupDatabase.UpdateGroupMember(ctx, member.GroupID, member.UserID, data); err != nil {
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
	if err := s.GroupDatabase.UpdateGroup(ctx, req.GroupID, UpdateGroupStatusMap(constant.GroupStatusMuted)); err != nil {
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
	if err := s.GroupDatabase.UpdateGroup(ctx, req.GroupID, UpdateGroupStatusMap(constant.GroupOk)); err != nil {
		return nil, err
	}
	s.Notification.GroupCancelMutedNotification(ctx, req.GroupID)
	return resp, nil
}

func (s *groupServer) SetGroupMemberInfo(ctx context.Context, req *pbGroup.SetGroupMemberInfoReq) (*pbGroup.SetGroupMemberInfoResp, error) {
	resp := &pbGroup.SetGroupMemberInfoResp{}
	if len(req.Members) == 0 {
		return nil, errs.ErrArgs.Wrap("members empty")
	}
	duplicateMap := make(map[[2]string]struct{})
	userIDMap := make(map[string]struct{})
	groupIDMap := make(map[string]struct{})
	for _, member := range req.Members {
		key := [...]string{member.GroupID, member.UserID}
		if _, ok := duplicateMap[key]; ok {
			return nil, errs.ErrArgs.Wrap("group user duplicate")
		}
		duplicateMap[key] = struct{}{}
		userIDMap[member.UserID] = struct{}{}
		groupIDMap[member.GroupID] = struct{}{}
	}
	groupIDs := utils.Keys(groupIDMap)
	userIDs := utils.Keys(userIDMap)
	members, err := s.FindGroupMember(ctx, groupIDs, append(userIDs, mcontext.GetOpUserID(ctx)), nil)
	if err != nil {
		return nil, err
	}
	for _, member := range members {
		delete(duplicateMap, [...]string{member.GroupID, member.UserID})
	}
	if len(duplicateMap) > 0 {
		return nil, errs.ErrArgs.Wrap("user not found" + strings.Join(utils.Slice(utils.Keys(duplicateMap), func(e [2]string) string {
			return fmt.Sprintf("[group: %s user: %s]", e[0], e[1])
		}), ","))
	}
	memberMap := utils.SliceToMap(members, func(e *relationTb.GroupMemberModel) [2]string {
		return [...]string{e.GroupID, e.UserID}
	})
	if !tokenverify.IsAppManagerUid(ctx) {
		opUserID := mcontext.GetOpUserID(ctx)
		for _, member := range req.Members {
			if member.RoleLevel != nil {
				switch member.RoleLevel.Value {
				case constant.GroupOrdinaryUsers, constant.GroupAdmin:
				default:
					return nil, errs.ErrArgs.Wrap("invalid role level")
				}
			}
			opMember, ok := memberMap[[...]string{member.GroupID, opUserID}]
			if !ok {
				return nil, errs.ErrArgs.Wrap(fmt.Sprintf("user %s not in group %s", opUserID, member.GroupID))
			}
			if member.UserID == opUserID {
				if member.RoleLevel != nil {
					return nil, errs.ErrNoPermission.Wrap("can not change self role level")
				}
				continue
			}
			if opMember.RoleLevel == constant.GroupOrdinaryUsers {
				return nil, errs.ErrNoPermission.Wrap("ordinary users can not change other role level")
			}
			dbMember, ok := memberMap[[...]string{member.GroupID, member.UserID}]
			if !ok {
				return nil, errs.ErrRecordNotFound.Wrap(fmt.Sprintf("user %s not in group %s", member.UserID, member.GroupID))
			}
			//if opMember.RoleLevel == constant.GroupOwner {
			//	continue
			//}
			//if dbMember.RoleLevel == constant.GroupOwner {
			//	return nil, errs.ErrNoPermission.Wrap("change group owner")
			//}
			//if opMember.RoleLevel == constant.GroupAdmin && dbMember.RoleLevel == constant.GroupAdmin {
			//	return nil, errs.ErrNoPermission.Wrap("admin can not change other admin role info")
			//}
			switch opMember.RoleLevel {
			case constant.GroupOrdinaryUsers:
				return nil, errs.ErrNoPermission.Wrap("ordinary users can not change other role level")
			case constant.GroupAdmin:
				if dbMember.RoleLevel != constant.GroupOrdinaryUsers {
					return nil, errs.ErrNoPermission.Wrap("admin can not change other role level")
				}
				if member.RoleLevel != nil {
					return nil, errs.ErrNoPermission.Wrap("admin can not change other role level")
				}
			case constant.GroupOwner:
				//if member.RoleLevel != nil && member.RoleLevel.Value == constant.GroupOwner {
				//	return nil, errs.ErrNoPermission.Wrap("owner only one")
				//}
			}
		}
	}
	for _, member := range req.Members {
		if member.RoleLevel == nil {
			continue
		}
		if memberMap[[...]string{member.GroupID, member.UserID}].RoleLevel == constant.GroupOwner {
			return nil, errs.ErrArgs.Wrap(fmt.Sprintf("group %s user %s is owner", member.GroupID, member.UserID))
		}
	}
	for i := 0; i < len(req.Members); i++ {
		if err := CallbackBeforeSetGroupMemberInfo(ctx, req.Members[i]); err != nil {
			return nil, err
		}
	}
	if err = s.GroupDatabase.UpdateGroupMembers(ctx, utils.Slice(req.Members, func(e *pbGroup.SetGroupMemberInfo) *relationTb.BatchUpdateGroupMember {
		return &relationTb.BatchUpdateGroupMember{
			GroupID: e.GroupID,
			UserID:  e.UserID,
			Map:     UpdateGroupMemberMap(e),
		}
	})); err != nil {
		return nil, err
	}
	for _, member := range req.Members {
		if member.RoleLevel != nil {
			switch member.RoleLevel.Value {
			case constant.GroupAdmin:
				s.Notification.GroupMemberSetToAdminNotification(ctx, member.GroupID, member.UserID)
			case constant.GroupOrdinaryUsers:
				s.Notification.GroupMemberSetToOrdinaryUserNotification(ctx, member.GroupID, member.UserID)
			}
		}
		if member.Nickname != nil || member.FaceURL != nil || member.Ex != nil {
			log.ZDebug(ctx, "setGroupMemberInfo notification", "member", member.UserID)
			if err := s.Notification.GroupMemberInfoSetNotification(ctx, member.GroupID, member.UserID); err != nil {
				log.ZError(ctx, "setGroupMemberInfo notification failed", err, "member", member.UserID, "groupID", member.GroupID)
			}
		}
	}
	return resp, nil
}

func (s *groupServer) GetGroupAbstractInfo(ctx context.Context, req *pbGroup.GetGroupAbstractInfoReq) (*pbGroup.GetGroupAbstractInfoResp, error) {
	resp := &pbGroup.GetGroupAbstractInfoResp{}
	if len(req.GroupIDs) == 0 {
		return nil, errs.ErrArgs.Wrap("groupIDs empty")
	}
	if utils.Duplicate(req.GroupIDs) {
		return nil, errs.ErrArgs.Wrap("groupIDs duplicate")
	}
	groups, err := s.GroupDatabase.FindGroup(ctx, req.GroupIDs)
	if err != nil {
		return nil, err
	}
	if ids := utils.Single(req.GroupIDs, utils.Slice(groups, func(group *relationTb.GroupModel) string {
		return group.GroupID
	})); len(ids) > 0 {
		return nil, errs.ErrGroupIDNotFound.Wrap("not found group " + strings.Join(ids, ","))
	}
	groupUserMap, err := s.GroupDatabase.MapGroupMemberUserID(ctx, req.GroupIDs)
	if err != nil {
		return nil, err
	}
	if ids := utils.Single(req.GroupIDs, utils.Keys(groupUserMap)); len(ids) > 0 {
		return nil, errs.ErrGroupIDNotFound.Wrap(fmt.Sprintf("group %s not found member", strings.Join(ids, ",")))
	}
	resp.GroupAbstractInfos = utils.Slice(groups, func(group *relationTb.GroupModel) *pbGroup.GroupAbstractInfo {
		users := groupUserMap[group.GroupID]
		return convert.Db2PbGroupAbstractInfo(group.GroupID, users.MemberNum, users.Hash)
	})
	return resp, nil
}

func (s *groupServer) GetUserInGroupMembers(ctx context.Context, req *pbGroup.GetUserInGroupMembersReq) (*pbGroup.GetUserInGroupMembersResp, error) {
	resp := &pbGroup.GetUserInGroupMembersResp{}
	if len(req.GroupIDs) == 0 {
		return nil, errs.ErrArgs.Wrap("groupIDs empty")
	}
	members, err := s.FindGroupMember(ctx, []string{req.UserID}, req.GroupIDs, nil)
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
		return convert.Db2PbGroupMember(e)
	})
	return resp, nil
}

func (s *groupServer) GetGroupMemberUserIDs(ctx context.Context, req *pbGroup.GetGroupMemberUserIDsReq) (resp *pbGroup.GetGroupMemberUserIDsResp, err error) {
	resp = &pbGroup.GetGroupMemberUserIDsResp{}
	resp.UserIDs, err = s.GroupDatabase.FindGroupMemberUserID(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *groupServer) GetGroupMemberRoleLevel(ctx context.Context, req *pbGroup.GetGroupMemberRoleLevelReq) (*pbGroup.GetGroupMemberRoleLevelResp, error) {
	resp := &pbGroup.GetGroupMemberRoleLevelResp{}
	if len(req.RoleLevels) == 0 {
		return nil, errs.ErrArgs.Wrap("RoleLevels empty")
	}
	members, err := s.FindGroupMember(ctx, []string{req.GroupID}, nil, req.RoleLevels)
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
		return convert.Db2PbGroupMember(e)
	})
	return resp, nil
}
