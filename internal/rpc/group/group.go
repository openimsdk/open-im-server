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
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/common"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database/mgo"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/open-im-server/v3/pkg/common/webhook"
	"github.com/openimsdk/open-im-server/v3/pkg/localcache"

	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
	"github.com/openimsdk/open-im-server/v3/pkg/common/convert"
	"github.com/openimsdk/open-im-server/v3/pkg/common/servererrs"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/msgprocessor"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient/grouphash"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient/notification"
	"github.com/openimsdk/protocol/constant"
	pbconversation "github.com/openimsdk/protocol/conversation"
	pbgroup "github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/protocol/wrapperspb"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/redisutil"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/mw/specialerror"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/openimsdk/tools/utils/encrypt"
	"google.golang.org/grpc"
)

type groupServer struct {
	db                    controller.GroupDatabase
	user                  rpcclient.UserRpcClient
	notification          *GroupNotificationSender
	conversationRpcClient rpcclient.ConversationRpcClient
	msgRpcClient          rpcclient.MessageRpcClient
	config                *Config
	webhookClient         *webhook.Client
}

type Config struct {
	RpcConfig          config.Group
	RedisConfig        config.Redis
	MongodbConfig      config.Mongo
	NotificationConfig config.Notification
	Share              config.Share
	WebhooksConfig     config.Webhooks
	LocalCacheConfig   config.LocalCache
	Discovery          config.Discovery
}

func Start(ctx context.Context, config *Config, client discovery.SvcDiscoveryRegistry, server *grpc.Server) error {
	mgocli, err := mongoutil.NewMongoDB(ctx, config.MongodbConfig.Build())
	if err != nil {
		return err
	}
	rdb, err := redisutil.NewRedisClient(ctx, config.RedisConfig.Build())
	if err != nil {
		return err
	}
	groupDB, err := mgo.NewGroupMongo(mgocli.GetDB())
	if err != nil {
		return err
	}
	groupMemberDB, err := mgo.NewGroupMember(mgocli.GetDB())
	if err != nil {
		return err
	}
	groupRequestDB, err := mgo.NewGroupRequestMgo(mgocli.GetDB())
	if err != nil {
		return err
	}
	userRpcClient := rpcclient.NewUserRpcClient(client, config.Share.RpcRegisterName.User, config.Share.IMAdminUserID)
	msgRpcClient := rpcclient.NewMessageRpcClient(client, config.Share.RpcRegisterName.Msg)
	conversationRpcClient := rpcclient.NewConversationRpcClient(client, config.Share.RpcRegisterName.Conversation)
	var gs groupServer
	database := controller.NewGroupDatabase(rdb, &config.LocalCacheConfig, groupDB, groupMemberDB, groupRequestDB, mgocli.GetTx(), grouphash.NewGroupHashFromGroupServer(&gs))
	gs.db = database
	gs.user = userRpcClient
	gs.notification = NewGroupNotificationSender(database, &msgRpcClient, &userRpcClient, config, func(ctx context.Context, userIDs []string) ([]notification.CommonUser, error) {
		users, err := userRpcClient.GetUsersInfo(ctx, userIDs)
		if err != nil {
			return nil, err
		}
		return datautil.Slice(users, func(e *sdkws.UserInfo) notification.CommonUser { return e }), nil
	})
	localcache.InitLocalCache(&config.LocalCacheConfig)
	gs.conversationRpcClient = conversationRpcClient
	gs.msgRpcClient = msgRpcClient
	gs.config = config
	gs.webhookClient = webhook.NewWebhookClient(config.WebhooksConfig.URL)
	pbgroup.RegisterGroupServer(server, &gs)
	return nil
}

func (s *groupServer) NotificationUserInfoUpdate(ctx context.Context, req *pbgroup.NotificationUserInfoUpdateReq) (*pbgroup.NotificationUserInfoUpdateResp, error) {
	members, err := s.db.FindGroupMemberUser(ctx, nil, req.UserID)
	if err != nil {
		return nil, err
	}
	groupIDs := make([]string, 0, len(members))
	for _, member := range members {
		if member.Nickname != "" && member.FaceURL != "" {
			continue
		}
		groupIDs = append(groupIDs, member.GroupID)
	}
	for _, groupID := range groupIDs {
		if err := s.db.MemberGroupIncrVersion(ctx, groupID, []string{req.UserID}, model.VersionStateUpdate); err != nil {
			return nil, err
		}
	}
	for _, groupID := range groupIDs {
		s.notification.GroupMemberInfoSetNotification(ctx, groupID, req.UserID)
	}
	if err = s.db.DeleteGroupMemberHash(ctx, groupIDs); err != nil {
		return nil, err
	}
	return &pbgroup.NotificationUserInfoUpdateResp{}, nil
}

func (s *groupServer) CheckGroupAdmin(ctx context.Context, groupID string) error {
	if !authverify.IsAppManagerUid(ctx, s.config.Share.IMAdminUserID) {
		groupMember, err := s.db.TakeGroupMember(ctx, groupID, mcontext.GetOpUserID(ctx))
		if err != nil {
			return err
		}
		if !(groupMember.RoleLevel == constant.GroupOwner || groupMember.RoleLevel == constant.GroupAdmin) {
			return errs.ErrNoPermission.WrapMsg("no group owner or admin")
		}
	}
	return nil
}

func (s *groupServer) GetPublicUserInfoMap(ctx context.Context, userIDs []string, complete bool) (map[string]*sdkws.PublicUserInfo, error) {
	if len(userIDs) == 0 {
		return map[string]*sdkws.PublicUserInfo{}, nil
	}
	users, err := s.user.GetPublicUserInfos(ctx, userIDs, complete)
	if err != nil {
		return nil, err
	}
	return datautil.SliceToMapAny(users, func(e *sdkws.PublicUserInfo) (string, *sdkws.PublicUserInfo) {
		return e.UserID, e
	}), nil
}

func (s *groupServer) IsNotFound(err error) bool {
	return errs.ErrRecordNotFound.Is(specialerror.ErrCode(errs.Unwrap(err)))
}

func (s *groupServer) GenGroupID(ctx context.Context, groupID *string) error {
	if *groupID != "" {
		_, err := s.db.TakeGroup(ctx, *groupID)
		if err == nil {
			return servererrs.ErrGroupIDExisted.WrapMsg("group id existed " + *groupID)
		} else if s.IsNotFound(err) {
			return nil
		} else {
			return err
		}
	}
	for i := 0; i < 10; i++ {
		id := encrypt.Md5(strings.Join([]string{mcontext.GetOperationID(ctx), strconv.FormatInt(time.Now().UnixNano(), 10), strconv.Itoa(rand.Int())}, ",;,"))
		bi := big.NewInt(0)
		bi.SetString(id[0:8], 16)
		id = bi.String()
		_, err := s.db.TakeGroup(ctx, id)
		if err == nil {
			continue
		} else if s.IsNotFound(err) {
			*groupID = id
			return nil
		} else {
			return err
		}
	}
	return servererrs.ErrData.WrapMsg("group id gen error")
}

func (s *groupServer) CreateGroup(ctx context.Context, req *pbgroup.CreateGroupReq) (*pbgroup.CreateGroupResp, error) {
	if req.GroupInfo.GroupType != constant.WorkingGroup {
		return nil, errs.ErrArgs.WrapMsg(fmt.Sprintf("group type only supports %d", constant.WorkingGroup))
	}
	if req.OwnerUserID == "" {
		return nil, errs.ErrArgs.WrapMsg("no group owner")
	}
	if err := authverify.CheckAccessV3(ctx, req.OwnerUserID, s.config.Share.IMAdminUserID); err != nil {

		return nil, err
	}
	userIDs := append(append(req.MemberUserIDs, req.AdminUserIDs...), req.OwnerUserID)
	opUserID := mcontext.GetOpUserID(ctx)
	if !datautil.Contain(opUserID, userIDs...) {
		userIDs = append(userIDs, opUserID)
	}

	if datautil.Duplicate(userIDs) {
		return nil, errs.ErrArgs.WrapMsg("group member repeated")
	}

	userMap, err := s.user.GetUsersInfoMap(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	if len(userMap) != len(userIDs) {
		return nil, servererrs.ErrUserIDNotFound.WrapMsg("user not found")
	}

	if err := s.webhookBeforeCreateGroup(ctx, &s.config.WebhooksConfig.BeforeCreateGroup, req); err != nil && err != servererrs.ErrCallbackContinue {
		return nil, err
	}

	var groupMembers []*model.GroupMember
	group := convert.Pb2DBGroupInfo(req.GroupInfo)
	if err := s.GenGroupID(ctx, &group.GroupID); err != nil {
		return nil, err
	}

	joinGroup := func(userID string, roleLevel int32) error {
		groupMember := &model.GroupMember{
			GroupID:        group.GroupID,
			UserID:         userID,
			RoleLevel:      roleLevel,
			OperatorUserID: opUserID,
			JoinSource:     constant.JoinByInvitation,
			InviterUserID:  opUserID,
			JoinTime:       time.Now(),
			MuteEndTime:    time.UnixMilli(0),
		}

		if err := s.webhookBeforeMemberJoinGroup(ctx, &s.config.WebhooksConfig.BeforeMemberJoinGroup, groupMember, group.Ex); err != nil && err != servererrs.ErrCallbackContinue {
			return err
		}
		groupMembers = append(groupMembers, groupMember)
		return nil
	}
	if err := joinGroup(req.OwnerUserID, constant.GroupOwner); err != nil {
		return nil, err
	}
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
	if err := s.db.CreateGroup(ctx, []*model.Group{group}, groupMembers); err != nil {
		return nil, err
	}
	resp := &pbgroup.CreateGroupResp{GroupInfo: &sdkws.GroupInfo{}}

	resp.GroupInfo = convert.Db2PbGroupInfo(group, req.OwnerUserID, uint32(len(userIDs)))
	resp.GroupInfo.MemberCount = uint32(len(userIDs))
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
	s.notification.GroupCreatedNotification(ctx, tips)

	reqCallBackAfter := &pbgroup.CreateGroupReq{
		MemberUserIDs: userIDs,
		GroupInfo:     resp.GroupInfo,
		OwnerUserID:   req.OwnerUserID,
		AdminUserIDs:  req.AdminUserIDs,
	}

	s.webhookAfterCreateGroup(ctx, &s.config.WebhooksConfig.AfterCreateGroup, reqCallBackAfter)

	return resp, nil
}

func (s *groupServer) GetJoinedGroupList(ctx context.Context, req *pbgroup.GetJoinedGroupListReq) (*pbgroup.GetJoinedGroupListResp, error) {
	if err := authverify.CheckAccessV3(ctx, req.FromUserID, s.config.Share.IMAdminUserID); err != nil {
		return nil, err
	}
	total, members, err := s.db.PageGetJoinGroup(ctx, req.FromUserID, req.Pagination)
	if err != nil {
		return nil, err
	}
	var resp pbgroup.GetJoinedGroupListResp
	resp.Total = uint32(total)
	if len(members) == 0 {
		return &resp, nil
	}
	groupIDs := datautil.Slice(members, func(e *model.GroupMember) string {
		return e.GroupID
	})
	groups, err := s.db.FindGroup(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	groupMemberNum, err := s.db.MapGroupMemberNum(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	owners, err := s.db.FindGroupsOwner(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	if err := s.PopulateGroupMember(ctx, members...); err != nil {
		return nil, err
	}
	ownerMap := datautil.SliceToMap(owners, func(e *model.GroupMember) string {
		return e.GroupID
	})
	resp.Groups = datautil.Slice(datautil.Order(groupIDs, groups, func(group *model.Group) string {
		return group.GroupID
	}), func(group *model.Group) *sdkws.GroupInfo {
		var userID string
		if user := ownerMap[group.GroupID]; user != nil {
			userID = user.UserID
		}
		return convert.Db2PbGroupInfo(group, userID, groupMemberNum[group.GroupID])
	})
	return &resp, nil
}

func (s *groupServer) InviteUserToGroup(ctx context.Context, req *pbgroup.InviteUserToGroupReq) (*pbgroup.InviteUserToGroupResp, error) {
	if len(req.InvitedUserIDs) == 0 {
		return nil, errs.ErrArgs.WrapMsg("user empty")
	}
	if datautil.Duplicate(req.InvitedUserIDs) {
		return nil, errs.ErrArgs.WrapMsg("userID duplicate")
	}
	group, err := s.db.TakeGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	if group.Status == constant.GroupStatusDismissed {
		return nil, servererrs.ErrDismissedAlready.WrapMsg("group dismissed checking group status found it dismissed")
	}

	userMap, err := s.user.GetUsersInfoMap(ctx, req.InvitedUserIDs)
	if err != nil {
		return nil, err
	}

	if len(userMap) != len(req.InvitedUserIDs) {
		return nil, errs.ErrRecordNotFound.WrapMsg("user not found")
	}

	var groupMember *model.GroupMember
	var opUserID string
	if !authverify.IsAppManagerUid(ctx, s.config.Share.IMAdminUserID) {
		opUserID = mcontext.GetOpUserID(ctx)
		var err error
		groupMember, err = s.db.TakeGroupMember(ctx, req.GroupID, opUserID)
		if err != nil {
			return nil, err
		}
		if err := s.PopulateGroupMember(ctx, groupMember); err != nil {
			return nil, err
		}
	}

	if err := s.webhookBeforeInviteUserToGroup(ctx, &s.config.WebhooksConfig.BeforeInviteUserToGroup, req); err != nil && err != servererrs.ErrCallbackContinue {
		return nil, err
	}

	if group.NeedVerification == constant.AllNeedVerification {
		if !authverify.IsAppManagerUid(ctx, s.config.Share.IMAdminUserID) {
			if !(groupMember.RoleLevel == constant.GroupOwner || groupMember.RoleLevel == constant.GroupAdmin) {
				var requests []*model.GroupRequest
				for _, userID := range req.InvitedUserIDs {
					requests = append(requests, &model.GroupRequest{
						UserID:        userID,
						GroupID:       req.GroupID,
						JoinSource:    constant.JoinByInvitation,
						InviterUserID: opUserID,
						ReqTime:       time.Now(),
						HandledTime:   time.Unix(0, 0),
					})
				}
				if err := s.db.CreateGroupRequest(ctx, requests); err != nil {
					return nil, err
				}
				for _, request := range requests {
					s.notification.JoinGroupApplicationNotification(ctx, &pbgroup.JoinGroupReq{
						GroupID:       request.GroupID,
						ReqMessage:    request.ReqMsg,
						JoinSource:    request.JoinSource,
						InviterUserID: request.InviterUserID,
					})
				}
				return &pbgroup.InviteUserToGroupResp{}, nil
			}
		}
	}
	var groupMembers []*model.GroupMember
	for _, userID := range req.InvitedUserIDs {
		member := &model.GroupMember{
			GroupID:        req.GroupID,
			UserID:         userID,
			RoleLevel:      constant.GroupOrdinaryUsers,
			OperatorUserID: opUserID,
			InviterUserID:  opUserID,
			JoinSource:     constant.JoinByInvitation,
			JoinTime:       time.Now(),
			MuteEndTime:    time.UnixMilli(0),
		}

		if err := s.webhookBeforeMemberJoinGroup(ctx, &s.config.WebhooksConfig.BeforeMemberJoinGroup, groupMember, group.Ex); err != nil && err != servererrs.ErrCallbackContinue {
			return nil, err
		}
		groupMembers = append(groupMembers, member)

	}
	if err := s.db.CreateGroup(ctx, nil, groupMembers); err != nil {
		return nil, err
	}
	if err := s.conversationRpcClient.GroupChatFirstCreateConversation(ctx, req.GroupID, req.InvitedUserIDs); err != nil {
		return nil, err
	}
	s.notification.MemberInvitedNotification(ctx, req.GroupID, req.Reason, req.InvitedUserIDs)
	return &pbgroup.InviteUserToGroupResp{}, nil
}

func (s *groupServer) GetGroupAllMember(ctx context.Context, req *pbgroup.GetGroupAllMemberReq) (*pbgroup.GetGroupAllMemberResp, error) {
	members, err := s.db.FindGroupMemberAll(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if err := s.PopulateGroupMember(ctx, members...); err != nil {
		return nil, err
	}
	var resp pbgroup.GetGroupAllMemberResp
	resp.Members = datautil.Slice(members, func(e *model.GroupMember) *sdkws.GroupMemberFullInfo {
		return convert.Db2PbGroupMember(e)
	})
	return &resp, nil
}

func (s *groupServer) GetGroupMemberList(ctx context.Context, req *pbgroup.GetGroupMemberListReq) (*pbgroup.GetGroupMemberListResp, error) {
	var (
		total   int64
		members []*model.GroupMember
		err     error
	)
	if req.Keyword == "" {
		total, members, err = s.db.PageGetGroupMember(ctx, req.GroupID, req.Pagination)
	} else {
		members, err = s.db.FindGroupMemberAll(ctx, req.GroupID)
	}
	if err != nil {
		return nil, err
	}
	if err := s.PopulateGroupMember(ctx, members...); err != nil {
		return nil, err
	}
	if req.Keyword != "" {
		groupMembers := make([]*model.GroupMember, 0)
		for _, member := range members {
			if member.UserID == req.Keyword {
				groupMembers = append(groupMembers, member)
				total++
				continue
			}
			if member.Nickname == req.Keyword {
				groupMembers = append(groupMembers, member)
				total++
				continue
			}
		}

		members := datautil.Paginate(groupMembers, int(req.Pagination.GetPageNumber()), int(req.Pagination.GetShowNumber()))
		return &pbgroup.GetGroupMemberListResp{
			Total:   uint32(total),
			Members: datautil.Batch(convert.Db2PbGroupMember, members),
		}, nil
	}
	return &pbgroup.GetGroupMemberListResp{
		Total:   uint32(total),
		Members: datautil.Batch(convert.Db2PbGroupMember, members),
	}, nil
}

func (s *groupServer) KickGroupMember(ctx context.Context, req *pbgroup.KickGroupMemberReq) (*pbgroup.KickGroupMemberResp, error) {
	group, err := s.db.TakeGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if len(req.KickedUserIDs) == 0 {
		return nil, errs.ErrArgs.WrapMsg("KickedUserIDs empty")
	}
	if datautil.Duplicate(req.KickedUserIDs) {
		return nil, errs.ErrArgs.WrapMsg("KickedUserIDs duplicate")
	}
	opUserID := mcontext.GetOpUserID(ctx)
	if datautil.Contain(opUserID, req.KickedUserIDs...) {
		return nil, errs.ErrArgs.WrapMsg("opUserID in KickedUserIDs")
	}
	owner, err := s.db.TakeGroupOwner(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if datautil.Contain(owner.UserID, req.KickedUserIDs...) {
		return nil, errs.ErrArgs.WrapMsg("ownerUID can not Kick")
	}

	members, err := s.db.FindGroupMembers(ctx, req.GroupID, append(req.KickedUserIDs, opUserID))
	if err != nil {
		return nil, err
	}
	if err := s.PopulateGroupMember(ctx, members...); err != nil {
		return nil, err
	}
	memberMap := make(map[string]*model.GroupMember)
	for i, member := range members {
		memberMap[member.UserID] = members[i]
	}
	isAppManagerUid := authverify.IsAppManagerUid(ctx, s.config.Share.IMAdminUserID)
	opMember := memberMap[opUserID]
	for _, userID := range req.KickedUserIDs {
		member, ok := memberMap[userID]
		if !ok {
			return nil, servererrs.ErrUserIDNotFound.WrapMsg(userID)
		}
		if !isAppManagerUid {
			if opMember == nil {
				return nil, errs.ErrNoPermission.WrapMsg("opUserID no in group")
			}
			switch opMember.RoleLevel {
			case constant.GroupOwner:
			case constant.GroupAdmin:
				if member.RoleLevel == constant.GroupOwner || member.RoleLevel == constant.GroupAdmin {
					return nil, errs.ErrNoPermission.WrapMsg("group admins cannot remove the group owner and other admins")
				}
			case constant.GroupOrdinaryUsers:
				return nil, errs.ErrNoPermission.WrapMsg("opUserID no permission")
			default:
				return nil, errs.ErrNoPermission.WrapMsg("opUserID roleLevel unknown")
			}
		}
	}
	num, err := s.db.FindGroupMemberNum(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	ownerUserIDs, err := s.db.GetGroupRoleLevelMemberIDs(ctx, req.GroupID, constant.GroupOwner)
	if err != nil {
		return nil, err
	}
	var ownerUserID string
	if len(ownerUserIDs) > 0 {
		ownerUserID = ownerUserIDs[0]
	}
	if err := s.db.DeleteGroupMember(ctx, group.GroupID, req.KickedUserIDs); err != nil {
		return nil, err
	}
	tips := &sdkws.MemberKickedTips{
		Group: &sdkws.GroupInfo{
			GroupID:                group.GroupID,
			GroupName:              group.GroupName,
			Notification:           group.Notification,
			Introduction:           group.Introduction,
			FaceURL:                group.FaceURL,
			OwnerUserID:            ownerUserID,
			CreateTime:             group.CreateTime.UnixMilli(),
			MemberCount:            num - uint32(len(req.KickedUserIDs)),
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
	if opMember, ok := memberMap[opUserID]; ok {
		tips.OpUser = convert.Db2PbGroupMember(opMember)
	}
	for _, userID := range req.KickedUserIDs {
		tips.KickedUserList = append(tips.KickedUserList, convert.Db2PbGroupMember(memberMap[userID]))
	}
	s.notification.MemberKickedNotification(ctx, tips)
	if err := s.deleteMemberAndSetConversationSeq(ctx, req.GroupID, req.KickedUserIDs); err != nil {
		return nil, err
	}
	s.webhookAfterKickGroupMember(ctx, &s.config.WebhooksConfig.AfterKickGroupMember, req)

	return &pbgroup.KickGroupMemberResp{}, nil
}

func (s *groupServer) GetGroupMembersInfo(ctx context.Context, req *pbgroup.GetGroupMembersInfoReq) (*pbgroup.GetGroupMembersInfoResp, error) {
	if len(req.UserIDs) == 0 {
		return nil, errs.ErrArgs.WrapMsg("userIDs empty")
	}
	if req.GroupID == "" {
		return nil, errs.ErrArgs.WrapMsg("groupID empty")
	}
	members, err := s.getGroupMembersInfo(ctx, req.GroupID, req.UserIDs)
	if err != nil {
		return nil, err
	}
	return &pbgroup.GetGroupMembersInfoResp{
		Members: members,
	}, nil
}

func (s *groupServer) getGroupMembersInfo(ctx context.Context, groupID string, userIDs []string) ([]*sdkws.GroupMemberFullInfo, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	members, err := s.db.FindGroupMembers(ctx, groupID, userIDs)
	if err != nil {
		return nil, err
	}
	if err := s.PopulateGroupMember(ctx, members...); err != nil {
		return nil, err
	}
	return datautil.Slice(members, func(e *model.GroupMember) *sdkws.GroupMemberFullInfo {
		return convert.Db2PbGroupMember(e)
	}), nil
}

// GetGroupApplicationList handles functions that get a list of group requests.
func (s *groupServer) GetGroupApplicationList(ctx context.Context, req *pbgroup.GetGroupApplicationListReq) (*pbgroup.GetGroupApplicationListResp, error) {
	groupIDs, err := s.db.FindUserManagedGroupID(ctx, req.FromUserID)
	if err != nil {
		return nil, err
	}
	resp := &pbgroup.GetGroupApplicationListResp{}
	if len(groupIDs) == 0 {
		return resp, nil
	}
	total, groupRequests, err := s.db.PageGroupRequest(ctx, groupIDs, req.Pagination)
	if err != nil {
		return nil, err
	}
	resp.Total = uint32(total)
	if len(groupRequests) == 0 {
		return resp, nil
	}
	var userIDs []string

	for _, gr := range groupRequests {
		userIDs = append(userIDs, gr.UserID)
	}
	userIDs = datautil.Distinct(userIDs)
	userMap, err := s.user.GetPublicUserInfoMap(ctx, userIDs, true)
	if err != nil {
		return nil, err
	}
	groups, err := s.db.FindGroup(ctx, datautil.Distinct(groupIDs))
	if err != nil {
		return nil, err
	}
	groupMap := datautil.SliceToMap(groups, func(e *model.Group) string {
		return e.GroupID
	})
	if ids := datautil.Single(datautil.Keys(groupMap), groupIDs); len(ids) > 0 {
		return nil, servererrs.ErrGroupIDNotFound.WrapMsg(strings.Join(ids, ","))
	}
	groupMemberNumMap, err := s.db.MapGroupMemberNum(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	owners, err := s.db.FindGroupsOwner(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	if err := s.PopulateGroupMember(ctx, owners...); err != nil {
		return nil, err
	}
	ownerMap := datautil.SliceToMap(owners, func(e *model.GroupMember) string {
		return e.GroupID
	})
	resp.GroupRequests = datautil.Slice(groupRequests, func(e *model.GroupRequest) *sdkws.GroupRequest {
		var ownerUserID string
		if owner, ok := ownerMap[e.GroupID]; ok {
			ownerUserID = owner.UserID
		}
		return convert.Db2PbGroupRequest(e, userMap[e.UserID], convert.Db2PbGroupInfo(groupMap[e.GroupID], ownerUserID, groupMemberNumMap[e.GroupID]))
	})
	return resp, nil
}

func (s *groupServer) GetGroupsInfo(ctx context.Context, req *pbgroup.GetGroupsInfoReq) (*pbgroup.GetGroupsInfoResp, error) {
	if len(req.GroupIDs) == 0 {
		return nil, errs.ErrArgs.WrapMsg("groupID is empty")
	}
	groups, err := s.getGroupsInfo(ctx, req.GroupIDs)
	if err != nil {
		return nil, err
	}
	return &pbgroup.GetGroupsInfoResp{
		GroupInfos: groups,
	}, nil
}

func (s *groupServer) getGroupsInfo(ctx context.Context, groupIDs []string) ([]*sdkws.GroupInfo, error) {
	if len(groupIDs) == 0 {
		return nil, nil
	}
	groups, err := s.db.FindGroup(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	groupMemberNumMap, err := s.db.MapGroupMemberNum(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	owners, err := s.db.FindGroupsOwner(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	if err := s.PopulateGroupMember(ctx, owners...); err != nil {
		return nil, err
	}
	ownerMap := datautil.SliceToMap(owners, func(e *model.GroupMember) string {
		return e.GroupID
	})
	return datautil.Slice(groups, func(e *model.Group) *sdkws.GroupInfo {
		var ownerUserID string
		if owner, ok := ownerMap[e.GroupID]; ok {
			ownerUserID = owner.UserID
		}
		return convert.Db2PbGroupInfo(e, ownerUserID, groupMemberNumMap[e.GroupID])
	}), nil
}

func (s *groupServer) GroupApplicationResponse(ctx context.Context, req *pbgroup.GroupApplicationResponseReq) (*pbgroup.GroupApplicationResponseResp, error) {
	if !datautil.Contain(req.HandleResult, constant.GroupResponseAgree, constant.GroupResponseRefuse) {
		return nil, errs.ErrArgs.WrapMsg("HandleResult unknown")
	}
	if !authverify.IsAppManagerUid(ctx, s.config.Share.IMAdminUserID) {
		groupMember, err := s.db.TakeGroupMember(ctx, req.GroupID, mcontext.GetOpUserID(ctx))
		if err != nil {
			return nil, err
		}
		if !(groupMember.RoleLevel == constant.GroupOwner || groupMember.RoleLevel == constant.GroupAdmin) {
			return nil, errs.ErrNoPermission.WrapMsg("no group owner or admin")
		}
	}
	group, err := s.db.TakeGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	groupRequest, err := s.db.TakeGroupRequest(ctx, req.GroupID, req.FromUserID)
	if err != nil {
		return nil, err
	}
	if groupRequest.HandleResult != 0 {
		return nil, servererrs.ErrGroupRequestHandled.WrapMsg("group request already processed")
	}
	var inGroup bool
	if _, err := s.db.TakeGroupMember(ctx, req.GroupID, req.FromUserID); err == nil {
		inGroup = true // Already in group
	} else if !s.IsNotFound(err) {
		return nil, err
	}
	if _, err := s.user.GetPublicUserInfo(ctx, req.FromUserID); err != nil {
		return nil, err
	}
	var member *model.GroupMember
	if (!inGroup) && req.HandleResult == constant.GroupResponseAgree {
		member = &model.GroupMember{
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
		if err := s.webhookBeforeMemberJoinGroup(ctx, &s.config.WebhooksConfig.BeforeMemberJoinGroup, member, group.Ex); err != nil && err != servererrs.ErrCallbackContinue {
			return nil, err
		}
	}
	log.ZDebug(ctx, "GroupApplicationResponse", "inGroup", inGroup, "HandleResult", req.HandleResult, "member", member)
	if err := s.db.HandlerGroupRequest(ctx, req.GroupID, req.FromUserID, req.HandledMsg, req.HandleResult, member); err != nil {
		return nil, err
	}
	switch req.HandleResult {
	case constant.GroupResponseAgree:
		if err := s.conversationRpcClient.GroupChatFirstCreateConversation(ctx, req.GroupID, []string{req.FromUserID}); err != nil {
			return nil, err
		}
		s.notification.GroupApplicationAcceptedNotification(ctx, req)
		if member == nil {
			log.ZDebug(ctx, "GroupApplicationResponse", "member is nil")
		} else {
			s.notification.MemberEnterNotification(ctx, req.GroupID, req.FromUserID)
		}
	case constant.GroupResponseRefuse:
		s.notification.GroupApplicationRejectedNotification(ctx, req)
	}

	return &pbgroup.GroupApplicationResponseResp{}, nil
}

func (s *groupServer) JoinGroup(ctx context.Context, req *pbgroup.JoinGroupReq) (*pbgroup.JoinGroupResp, error) {
	user, err := s.user.GetUserInfo(ctx, req.InviterUserID)
	if err != nil {
		return nil, err
	}
	group, err := s.db.TakeGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if group.Status == constant.GroupStatusDismissed {
		return nil, servererrs.ErrDismissedAlready.Wrap()
	}

	reqCall := &callbackstruct.CallbackJoinGroupReq{
		GroupID:    req.GroupID,
		GroupType:  string(group.GroupType),
		ApplyID:    req.InviterUserID,
		ReqMessage: req.ReqMessage,
		Ex:         req.Ex,
	}

	if err := s.webhookBeforeApplyJoinGroup(ctx, &s.config.WebhooksConfig.BeforeApplyJoinGroup, reqCall); err != nil && err != servererrs.ErrCallbackContinue {
		return nil, err
	}

	_, err = s.db.TakeGroupMember(ctx, req.GroupID, req.InviterUserID)
	if err == nil {
		return nil, errs.ErrArgs.Wrap()
	} else if !s.IsNotFound(err) && errs.Unwrap(err) != errs.ErrRecordNotFound {
		return nil, err
	}
	log.ZDebug(ctx, "JoinGroup.groupInfo", "group", group, "eq", group.NeedVerification == constant.Directly)
	if group.NeedVerification == constant.Directly {
		groupMember := &model.GroupMember{
			GroupID:        group.GroupID,
			UserID:         user.UserID,
			RoleLevel:      constant.GroupOrdinaryUsers,
			OperatorUserID: mcontext.GetOpUserID(ctx),
			InviterUserID:  req.InviterUserID,
			JoinTime:       time.Now(),
			MuteEndTime:    time.UnixMilli(0),
		}

		if err := s.webhookBeforeMemberJoinGroup(ctx, &s.config.WebhooksConfig.BeforeMemberJoinGroup, groupMember, group.Ex); err != nil && err != servererrs.ErrCallbackContinue {
			return nil, err
		}

		if err := s.db.CreateGroup(ctx, nil, []*model.GroupMember{groupMember}); err != nil {
			return nil, err
		}

		if err := s.conversationRpcClient.GroupChatFirstCreateConversation(ctx, req.GroupID, []string{req.InviterUserID}); err != nil {
			return nil, err
		}
		s.notification.MemberEnterNotification(ctx, req.GroupID, req.InviterUserID)
		s.webhookAfterJoinGroup(ctx, &s.config.WebhooksConfig.AfterJoinGroup, req)

		return &pbgroup.JoinGroupResp{}, nil
	}
	groupRequest := model.GroupRequest{
		UserID:      req.InviterUserID,
		ReqMsg:      req.ReqMessage,
		GroupID:     req.GroupID,
		JoinSource:  req.JoinSource,
		ReqTime:     time.Now(),
		HandledTime: time.Unix(0, 0),
		Ex:          req.Ex,
	}
	if err = s.db.CreateGroupRequest(ctx, []*model.GroupRequest{&groupRequest}); err != nil {
		return nil, err
	}
	s.notification.JoinGroupApplicationNotification(ctx, req)
	return &pbgroup.JoinGroupResp{}, nil
}

func (s *groupServer) QuitGroup(ctx context.Context, req *pbgroup.QuitGroupReq) (*pbgroup.QuitGroupResp, error) {
	if req.UserID == "" {
		req.UserID = mcontext.GetOpUserID(ctx)
	} else {
		if err := authverify.CheckAccessV3(ctx, req.UserID, s.config.Share.IMAdminUserID); err != nil {
			return nil, err
		}
	}
	member, err := s.db.TakeGroupMember(ctx, req.GroupID, req.UserID)
	if err != nil {
		return nil, err
	}
	if member.RoleLevel == constant.GroupOwner {
		return nil, errs.ErrNoPermission.WrapMsg("group owner can't quit")
	}
	if err := s.PopulateGroupMember(ctx, member); err != nil {
		return nil, err
	}
	err = s.db.DeleteGroupMember(ctx, req.GroupID, []string{req.UserID})
	if err != nil {
		return nil, err
	}
	s.notification.MemberQuitNotification(ctx, s.groupMemberDB2PB(member, 0))
	if err := s.deleteMemberAndSetConversationSeq(ctx, req.GroupID, []string{req.UserID}); err != nil {
		return nil, err
	}
	s.webhookAfterQuitGroup(ctx, &s.config.WebhooksConfig.AfterQuitGroup, req)

	return &pbgroup.QuitGroupResp{}, nil
}

func (s *groupServer) deleteMemberAndSetConversationSeq(ctx context.Context, groupID string, userIDs []string) error {
	conevrsationID := msgprocessor.GetConversationIDBySessionType(constant.ReadGroupChatType, groupID)
	maxSeq, err := s.msgRpcClient.GetConversationMaxSeq(ctx, conevrsationID)
	if err != nil {
		return err
	}
	return s.conversationRpcClient.SetConversationMaxSeq(ctx, userIDs, conevrsationID, maxSeq)
}

func (s *groupServer) SetGroupInfo(ctx context.Context, req *pbgroup.SetGroupInfoReq) (*pbgroup.SetGroupInfoResp, error) {
	var opMember *model.GroupMember
	if !authverify.IsAppManagerUid(ctx, s.config.Share.IMAdminUserID) {
		var err error
		opMember, err = s.db.TakeGroupMember(ctx, req.GroupInfoForSet.GroupID, mcontext.GetOpUserID(ctx))
		if err != nil {
			return nil, err
		}
		if !(opMember.RoleLevel == constant.GroupOwner || opMember.RoleLevel == constant.GroupAdmin) {
			return nil, errs.ErrNoPermission.WrapMsg("no group owner or admin")
		}
		if err := s.PopulateGroupMember(ctx, opMember); err != nil {
			return nil, err
		}
	}

	if err := s.webhookBeforeSetGroupInfo(ctx, &s.config.WebhooksConfig.BeforeSetGroupInfo, req); err != nil && err != servererrs.ErrCallbackContinue {
		return nil, err
	}

	group, err := s.db.TakeGroup(ctx, req.GroupInfoForSet.GroupID)
	if err != nil {
		return nil, err
	}
	if group.Status == constant.GroupStatusDismissed {
		return nil, servererrs.ErrDismissedAlready.Wrap()
	}

	count, err := s.db.FindGroupMemberNum(ctx, group.GroupID)
	if err != nil {
		return nil, err
	}
	owner, err := s.db.TakeGroupOwner(ctx, group.GroupID)
	if err != nil {
		return nil, err
	}
	if err := s.PopulateGroupMember(ctx, owner); err != nil {
		return nil, err
	}
	update := UpdateGroupInfoMap(ctx, req.GroupInfoForSet)
	if len(update) == 0 {
		return &pbgroup.SetGroupInfoResp{}, nil
	}
	if err := s.db.UpdateGroup(ctx, group.GroupID, update); err != nil {
		return nil, err
	}
	group, err = s.db.TakeGroup(ctx, req.GroupInfoForSet.GroupID)
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
	num := len(update)
	if req.GroupInfoForSet.Notification != "" {
		num--
		func() {
			conversation := &pbconversation.ConversationReq{
				ConversationID:   msgprocessor.GetConversationIDBySessionType(constant.ReadGroupChatType, req.GroupInfoForSet.GroupID),
				ConversationType: constant.ReadGroupChatType,
				GroupID:          req.GroupInfoForSet.GroupID,
			}
			resp, err := s.GetGroupMemberUserIDs(ctx, &pbgroup.GetGroupMemberUserIDsReq{GroupID: req.GroupInfoForSet.GroupID})
			if err != nil {
				log.ZWarn(ctx, "GetGroupMemberIDs", err)
				return
			}
			conversation.GroupAtType = &wrapperspb.Int32Value{Value: constant.GroupNotification}
			if err := s.conversationRpcClient.SetConversations(ctx, resp.UserIDs, conversation); err != nil {
				log.ZWarn(ctx, "SetConversations", err, resp.UserIDs, conversation)
			}
		}()
		s.notification.GroupInfoSetAnnouncementNotification(ctx, &sdkws.GroupInfoSetAnnouncementTips{Group: tips.Group, OpUser: tips.OpUser})
	}
	if req.GroupInfoForSet.GroupName != "" {
		num--
		s.notification.GroupInfoSetNameNotification(ctx, &sdkws.GroupInfoSetNameTips{Group: tips.Group, OpUser: tips.OpUser})
	}
	if num > 0 {
		s.notification.GroupInfoSetNotification(ctx, tips)
	}

	s.webhookAfterSetGroupInfo(ctx, &s.config.WebhooksConfig.AfterSetGroupInfo, req)

	return &pbgroup.SetGroupInfoResp{}, nil
}

func (s *groupServer) TransferGroupOwner(ctx context.Context, req *pbgroup.TransferGroupOwnerReq) (*pbgroup.TransferGroupOwnerResp, error) {
	group, err := s.db.TakeGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if group.Status == constant.GroupStatusDismissed {
		return nil, servererrs.ErrDismissedAlready.Wrap()
	}
	if req.OldOwnerUserID == req.NewOwnerUserID {
		return nil, errs.ErrArgs.WrapMsg("OldOwnerUserID == NewOwnerUserID")
	}
	members, err := s.db.FindGroupMembers(ctx, req.GroupID, []string{req.OldOwnerUserID, req.NewOwnerUserID})
	if err != nil {
		return nil, err
	}
	if err := s.PopulateGroupMember(ctx, members...); err != nil {
		return nil, err
	}
	memberMap := datautil.SliceToMap(members, func(e *model.GroupMember) string { return e.UserID })
	if ids := datautil.Single([]string{req.OldOwnerUserID, req.NewOwnerUserID}, datautil.Keys(memberMap)); len(ids) > 0 {
		return nil, errs.ErrArgs.WrapMsg("user not in group " + strings.Join(ids, ","))
	}
	oldOwner := memberMap[req.OldOwnerUserID]
	if oldOwner == nil {
		return nil, errs.ErrArgs.WrapMsg("OldOwnerUserID not in group " + req.NewOwnerUserID)
	}
	newOwner := memberMap[req.NewOwnerUserID]
	if newOwner == nil {
		return nil, errs.ErrArgs.WrapMsg("NewOwnerUser not in group " + req.NewOwnerUserID)
	}
	if !authverify.IsAppManagerUid(ctx, s.config.Share.IMAdminUserID) {
		if !(mcontext.GetOpUserID(ctx) == oldOwner.UserID && oldOwner.RoleLevel == constant.GroupOwner) {
			return nil, errs.ErrNoPermission.WrapMsg("no permission transfer group owner")
		}
	}
	if err := s.db.TransferGroupOwner(ctx, req.GroupID, req.OldOwnerUserID, req.NewOwnerUserID, newOwner.RoleLevel); err != nil {
		return nil, err
	}

	s.webhookAfterTransferGroupOwner(ctx, &s.config.WebhooksConfig.AfterTransferGroupOwner, req)

	s.notification.GroupOwnerTransferredNotification(ctx, req)
	return &pbgroup.TransferGroupOwnerResp{}, nil
}

func (s *groupServer) GetGroups(ctx context.Context, req *pbgroup.GetGroupsReq) (*pbgroup.GetGroupsResp, error) {
	var (
		group []*model.Group
		err   error
	)
	var resp pbgroup.GetGroupsResp
	if req.GroupID != "" {
		group, err = s.db.FindGroup(ctx, []string{req.GroupID})
		resp.Total = uint32(len(group))
	} else {
		var total int64
		total, group, err = s.db.SearchGroup(ctx, req.GroupName, req.Pagination)
		resp.Total = uint32(total)
	}

	if err != nil {
		return nil, err
	}

	groupIDs := datautil.Slice(group, func(e *model.Group) string {
		return e.GroupID
	})

	ownerMembers, err := s.db.FindGroupsOwner(ctx, groupIDs)
	if err != nil {
		return nil, err
	}

	ownerMemberMap := datautil.SliceToMap(ownerMembers, func(e *model.GroupMember) string {
		return e.GroupID
	})
	groupMemberNumMap, err := s.db.MapGroupMemberNum(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	resp.Groups = datautil.Slice(group, func(group *model.Group) *pbgroup.CMSGroup {
		var (
			userID   string
			username string
		)
		if member, ok := ownerMemberMap[group.GroupID]; ok {
			userID = member.UserID
			username = member.Nickname
		}
		return convert.Db2PbCMSGroup(group, userID, username, groupMemberNumMap[group.GroupID])
	})
	return &resp, nil
}

func (s *groupServer) GetGroupMembersCMS(ctx context.Context, req *pbgroup.GetGroupMembersCMSReq) (*pbgroup.GetGroupMembersCMSResp, error) {
	total, members, err := s.db.SearchGroupMember(ctx, req.UserName, req.GroupID, req.Pagination)
	if err != nil {
		return nil, err
	}
	var resp pbgroup.GetGroupMembersCMSResp
	resp.Total = uint32(total)
	if err := s.PopulateGroupMember(ctx, members...); err != nil {
		return nil, err
	}
	resp.Members = datautil.Slice(members, func(e *model.GroupMember) *sdkws.GroupMemberFullInfo {
		return convert.Db2PbGroupMember(e)
	})
	return &resp, nil
}

func (s *groupServer) GetUserReqApplicationList(ctx context.Context, req *pbgroup.GetUserReqApplicationListReq) (*pbgroup.GetUserReqApplicationListResp, error) {
	user, err := s.user.GetPublicUserInfo(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	total, requests, err := s.db.PageGroupRequestUser(ctx, req.UserID, req.Pagination)
	if err != nil {
		return nil, err
	}
	if len(requests) == 0 {
		return &pbgroup.GetUserReqApplicationListResp{Total: uint32(total)}, nil
	}
	groupIDs := datautil.Distinct(datautil.Slice(requests, func(e *model.GroupRequest) string {
		return e.GroupID
	}))
	groups, err := s.db.FindGroup(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	groupMap := datautil.SliceToMap(groups, func(e *model.Group) string {
		return e.GroupID
	})
	owners, err := s.db.FindGroupsOwner(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	if err := s.PopulateGroupMember(ctx, owners...); err != nil {
		return nil, err
	}
	ownerMap := datautil.SliceToMap(owners, func(e *model.GroupMember) string {
		return e.GroupID
	})
	groupMemberNum, err := s.db.MapGroupMemberNum(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	return &pbgroup.GetUserReqApplicationListResp{
		Total: uint32(total),
		GroupRequests: datautil.Slice(requests, func(e *model.GroupRequest) *sdkws.GroupRequest {
			var ownerUserID string
			if owner, ok := ownerMap[e.GroupID]; ok {
				ownerUserID = owner.UserID
			}
			return convert.Db2PbGroupRequest(e, user, convert.Db2PbGroupInfo(groupMap[e.GroupID], ownerUserID, groupMemberNum[e.GroupID]))
		}),
	}, nil
}

func (s *groupServer) DismissGroup(ctx context.Context, req *pbgroup.DismissGroupReq) (*pbgroup.DismissGroupResp, error) {
	owner, err := s.db.TakeGroupOwner(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if !authverify.IsAppManagerUid(ctx, s.config.Share.IMAdminUserID) {
		if owner.UserID != mcontext.GetOpUserID(ctx) {
			return nil, errs.ErrNoPermission.WrapMsg("not group owner")
		}
	}
	if err := s.PopulateGroupMember(ctx, owner); err != nil {
		return nil, err
	}
	group, err := s.db.TakeGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if !req.DeleteMember && group.Status == constant.GroupStatusDismissed {
		return nil, servererrs.ErrDismissedAlready.WrapMsg("group status is dismissed")
	}
	if err := s.db.DismissGroup(ctx, req.GroupID, req.DeleteMember); err != nil {
		return nil, err
	}
	if !req.DeleteMember {
		num, err := s.db.FindGroupMemberNum(ctx, req.GroupID)
		if err != nil {
			return nil, err
		}
		tips := &sdkws.GroupDismissedTips{
			Group:  s.groupDB2PB(group, owner.UserID, num),
			OpUser: &sdkws.GroupMemberFullInfo{},
		}
		if mcontext.GetOpUserID(ctx) == owner.UserID {
			tips.OpUser = s.groupMemberDB2PB(owner, 0)
		}
		s.notification.GroupDismissedNotification(ctx, tips)
	}
	membersID, err := s.db.FindGroupMemberUserID(ctx, group.GroupID)
	if err != nil {
		return nil, err
	}
	cbReq := &callbackstruct.CallbackDisMissGroupReq{
		GroupID:   req.GroupID,
		OwnerID:   owner.UserID,
		MembersID: membersID,
		GroupType: string(group.GroupType),
	}

	s.webhookAfterDismissGroup(ctx, &s.config.WebhooksConfig.AfterDismissGroup, cbReq)

	return &pbgroup.DismissGroupResp{}, nil
}

func (s *groupServer) MuteGroupMember(ctx context.Context, req *pbgroup.MuteGroupMemberReq) (*pbgroup.MuteGroupMemberResp, error) {
	member, err := s.db.TakeGroupMember(ctx, req.GroupID, req.UserID)
	if err != nil {
		return nil, err
	}
	if err := s.PopulateGroupMember(ctx, member); err != nil {
		return nil, err
	}
	if !authverify.IsAppManagerUid(ctx, s.config.Share.IMAdminUserID) {
		opMember, err := s.db.TakeGroupMember(ctx, req.GroupID, mcontext.GetOpUserID(ctx))
		if err != nil {
			return nil, err
		}
		switch member.RoleLevel {
		case constant.GroupOwner:
			return nil, errs.ErrNoPermission.WrapMsg("set group owner mute")
		case constant.GroupAdmin:
			if opMember.RoleLevel != constant.GroupOwner {
				return nil, errs.ErrNoPermission.WrapMsg("set group admin mute")
			}
		case constant.GroupOrdinaryUsers:
			if !(opMember.RoleLevel == constant.GroupAdmin || opMember.RoleLevel == constant.GroupOwner) {
				return nil, errs.ErrNoPermission.WrapMsg("set group ordinary users mute")
			}
		}
	}
	data := UpdateGroupMemberMutedTimeMap(time.Now().Add(time.Second * time.Duration(req.MutedSeconds)))
	if err := s.db.UpdateGroupMember(ctx, member.GroupID, member.UserID, data); err != nil {
		return nil, err
	}
	s.notification.GroupMemberMutedNotification(ctx, req.GroupID, req.UserID, req.MutedSeconds)
	return &pbgroup.MuteGroupMemberResp{}, nil
}

func (s *groupServer) CancelMuteGroupMember(ctx context.Context, req *pbgroup.CancelMuteGroupMemberReq) (*pbgroup.CancelMuteGroupMemberResp, error) {
	member, err := s.db.TakeGroupMember(ctx, req.GroupID, req.UserID)
	if err != nil {
		return nil, err
	}
	if err := s.PopulateGroupMember(ctx, member); err != nil {
		return nil, err
	}
	if !authverify.IsAppManagerUid(ctx, s.config.Share.IMAdminUserID) {
		opMember, err := s.db.TakeGroupMember(ctx, req.GroupID, mcontext.GetOpUserID(ctx))
		if err != nil {
			return nil, err
		}
		switch member.RoleLevel {
		case constant.GroupOwner:
			return nil, errs.ErrNoPermission.WrapMsg("set group owner mute")
		case constant.GroupAdmin:
			if opMember.RoleLevel != constant.GroupOwner {
				return nil, errs.ErrNoPermission.WrapMsg("set group admin mute")
			}
		case constant.GroupOrdinaryUsers:
			if !(opMember.RoleLevel == constant.GroupAdmin || opMember.RoleLevel == constant.GroupOwner) {
				return nil, errs.ErrNoPermission.WrapMsg("set group ordinary users mute")
			}
		}
	}
	data := UpdateGroupMemberMutedTimeMap(time.Unix(0, 0))
	if err := s.db.UpdateGroupMember(ctx, member.GroupID, member.UserID, data); err != nil {
		return nil, err
	}
	s.notification.GroupMemberCancelMutedNotification(ctx, req.GroupID, req.UserID)
	return &pbgroup.CancelMuteGroupMemberResp{}, nil
}

func (s *groupServer) MuteGroup(ctx context.Context, req *pbgroup.MuteGroupReq) (*pbgroup.MuteGroupResp, error) {
	if err := s.CheckGroupAdmin(ctx, req.GroupID); err != nil {
		return nil, err
	}
	if err := s.db.UpdateGroup(ctx, req.GroupID, UpdateGroupStatusMap(constant.GroupStatusMuted)); err != nil {
		return nil, err
	}
	s.notification.GroupMutedNotification(ctx, req.GroupID)
	return &pbgroup.MuteGroupResp{}, nil
}

func (s *groupServer) CancelMuteGroup(ctx context.Context, req *pbgroup.CancelMuteGroupReq) (*pbgroup.CancelMuteGroupResp, error) {
	if err := s.CheckGroupAdmin(ctx, req.GroupID); err != nil {
		return nil, err
	}
	if err := s.db.UpdateGroup(ctx, req.GroupID, UpdateGroupStatusMap(constant.GroupOk)); err != nil {
		return nil, err
	}
	s.notification.GroupCancelMutedNotification(ctx, req.GroupID)
	return &pbgroup.CancelMuteGroupResp{}, nil
}

func (s *groupServer) SetGroupMemberInfo(ctx context.Context, req *pbgroup.SetGroupMemberInfoReq) (*pbgroup.SetGroupMemberInfoResp, error) {
	if len(req.Members) == 0 {
		return nil, errs.ErrArgs.WrapMsg("members empty")
	}
	opUserID := mcontext.GetOpUserID(ctx)
	if opUserID == "" {
		return nil, errs.ErrNoPermission.WrapMsg("no op user id")
	}
	isAppManagerUid := authverify.IsAppManagerUid(ctx, s.config.Share.IMAdminUserID)
	for i := range req.Members {
		req.Members[i].FaceURL = nil
	}
	groupMembers := make(map[string][]*pbgroup.SetGroupMemberInfo)
	for i, member := range req.Members {
		if member.RoleLevel != nil {
			switch member.RoleLevel.Value {
			case constant.GroupOwner:
				return nil, errs.ErrNoPermission.WrapMsg("cannot set ungroup owner")
			case constant.GroupAdmin, constant.GroupOrdinaryUsers:
			default:
				return nil, errs.ErrArgs.WrapMsg("invalid role level")
			}
		}
		groupMembers[member.GroupID] = append(groupMembers[member.GroupID], req.Members[i])
	}
	for groupID, members := range groupMembers {
		temp := make(map[string]struct{})
		userIDs := make([]string, 0, len(members)+1)
		for _, member := range members {
			if _, ok := temp[member.UserID]; ok {
				return nil, errs.ErrArgs.WrapMsg(fmt.Sprintf("repeat group %s user %s", member.GroupID, member.UserID))
			}
			temp[member.UserID] = struct{}{}
			userIDs = append(userIDs, member.UserID)
		}
		if _, ok := temp[opUserID]; !ok {
			userIDs = append(userIDs, opUserID)
		}
		dbMembers, err := s.db.FindGroupMembers(ctx, groupID, userIDs)
		if err != nil {
			return nil, err
		}
		opUserIndex := -1
		for i, member := range dbMembers {
			if member.UserID == opUserID {
				opUserIndex = i
				break
			}
		}
		switch len(userIDs) - len(dbMembers) {
		case 0:
			if !isAppManagerUid {
				roleLevel := dbMembers[opUserIndex].RoleLevel
				if roleLevel != constant.GroupOwner {
					switch roleLevel {
					case constant.GroupAdmin:
						for _, member := range dbMembers {
							if member.RoleLevel == constant.GroupOwner {
								return nil, errs.ErrNoPermission.WrapMsg("admin can not change group owner")
							}
							if member.RoleLevel == constant.GroupAdmin && member.UserID != opUserID {
								return nil, errs.ErrNoPermission.WrapMsg("admin can not change other group admin")
							}
						}
					case constant.GroupOrdinaryUsers:
						for _, member := range dbMembers {
							if !(member.RoleLevel == constant.GroupOrdinaryUsers && member.UserID == opUserID) {
								return nil, errs.ErrNoPermission.WrapMsg("ordinary users can not change other role level")
							}
						}
					default:
						for _, member := range dbMembers {
							if member.RoleLevel >= roleLevel {
								return nil, errs.ErrNoPermission.WrapMsg("can not change higher role level")
							}
						}
					}
				}
			}
		case 1:
			if opUserIndex >= 0 {
				return nil, errs.ErrArgs.WrapMsg("user not in group")
			}
			if !isAppManagerUid {
				return nil, errs.ErrNoPermission.WrapMsg("user not in group")
			}
		default:
			return nil, errs.ErrArgs.WrapMsg("user not in group")
		}
	}

	for i := 0; i < len(req.Members); i++ {

		if err := s.webhookBeforeSetGroupMemberInfo(ctx, &s.config.WebhooksConfig.BeforeSetGroupMemberInfo, req.Members[i]); err != nil && err != servererrs.ErrCallbackContinue {
			return nil, err
		}

	}
	if err := s.db.UpdateGroupMembers(ctx, datautil.Slice(req.Members, func(e *pbgroup.SetGroupMemberInfo) *common.BatchUpdateGroupMember {
		return &common.BatchUpdateGroupMember{
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
				s.notification.GroupMemberSetToAdminNotification(ctx, member.GroupID, member.UserID)
			case constant.GroupOrdinaryUsers:
				s.notification.GroupMemberSetToOrdinaryUserNotification(ctx, member.GroupID, member.UserID)
			}
		}
		if member.Nickname != nil || member.FaceURL != nil || member.Ex != nil {
			s.notification.GroupMemberInfoSetNotification(ctx, member.GroupID, member.UserID)
		}
	}
	for i := 0; i < len(req.Members); i++ {
		s.webhookAfterSetGroupMemberInfo(ctx, &s.config.WebhooksConfig.AfterSetGroupMemberInfo, req.Members[i])
	}

	return &pbgroup.SetGroupMemberInfoResp{}, nil
}

func (s *groupServer) GetGroupAbstractInfo(ctx context.Context, req *pbgroup.GetGroupAbstractInfoReq) (*pbgroup.GetGroupAbstractInfoResp, error) {
	if len(req.GroupIDs) == 0 {
		return nil, errs.ErrArgs.WrapMsg("groupIDs empty")
	}
	if datautil.Duplicate(req.GroupIDs) {
		return nil, errs.ErrArgs.WrapMsg("groupIDs duplicate")
	}
	groups, err := s.db.FindGroup(ctx, req.GroupIDs)
	if err != nil {
		return nil, err
	}
	if ids := datautil.Single(req.GroupIDs, datautil.Slice(groups, func(group *model.Group) string {
		return group.GroupID
	})); len(ids) > 0 {
		return nil, servererrs.ErrGroupIDNotFound.WrapMsg("not found group " + strings.Join(ids, ","))
	}
	groupUserMap, err := s.db.MapGroupMemberUserID(ctx, req.GroupIDs)
	if err != nil {
		return nil, err
	}
	if ids := datautil.Single(req.GroupIDs, datautil.Keys(groupUserMap)); len(ids) > 0 {
		return nil, servererrs.ErrGroupIDNotFound.WrapMsg(fmt.Sprintf("group %s not found member", strings.Join(ids, ",")))
	}
	return &pbgroup.GetGroupAbstractInfoResp{
		GroupAbstractInfos: datautil.Slice(groups, func(group *model.Group) *pbgroup.GroupAbstractInfo {
			users := groupUserMap[group.GroupID]
			return convert.Db2PbGroupAbstractInfo(group.GroupID, users.MemberNum, users.Hash)
		}),
	}, nil
}

func (s *groupServer) GetUserInGroupMembers(ctx context.Context, req *pbgroup.GetUserInGroupMembersReq) (*pbgroup.GetUserInGroupMembersResp, error) {
	if len(req.GroupIDs) == 0 {
		return nil, errs.ErrArgs.WrapMsg("groupIDs empty")
	}
	members, err := s.db.FindGroupMemberUser(ctx, req.GroupIDs, req.UserID)
	if err != nil {
		return nil, err
	}
	if err := s.PopulateGroupMember(ctx, members...); err != nil {
		return nil, err
	}
	return &pbgroup.GetUserInGroupMembersResp{
		Members: datautil.Slice(members, func(e *model.GroupMember) *sdkws.GroupMemberFullInfo {
			return convert.Db2PbGroupMember(e)
		}),
	}, nil
}

func (s *groupServer) GetGroupMemberUserIDs(ctx context.Context, req *pbgroup.GetGroupMemberUserIDsReq) (*pbgroup.GetGroupMemberUserIDsResp, error) {
	userIDs, err := s.db.FindGroupMemberUserID(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	return &pbgroup.GetGroupMemberUserIDsResp{
		UserIDs: userIDs,
	}, nil
}

func (s *groupServer) GetGroupMemberRoleLevel(ctx context.Context, req *pbgroup.GetGroupMemberRoleLevelReq) (*pbgroup.GetGroupMemberRoleLevelResp, error) {
	if len(req.RoleLevels) == 0 {
		return nil, errs.ErrArgs.WrapMsg("RoleLevels empty")
	}
	members, err := s.db.FindGroupMemberRoleLevels(ctx, req.GroupID, req.RoleLevels)
	if err != nil {
		return nil, err
	}
	if err := s.PopulateGroupMember(ctx, members...); err != nil {
		return nil, err
	}
	return &pbgroup.GetGroupMemberRoleLevelResp{
		Members: datautil.Slice(members, func(e *model.GroupMember) *sdkws.GroupMemberFullInfo {
			return convert.Db2PbGroupMember(e)
		}),
	}, nil
}

func (s *groupServer) GetGroupUsersReqApplicationList(ctx context.Context, req *pbgroup.GetGroupUsersReqApplicationListReq) (*pbgroup.GetGroupUsersReqApplicationListResp, error) {
	requests, err := s.db.FindGroupRequests(ctx, req.GroupID, req.UserIDs)
	if err != nil {
		return nil, err
	}
	if len(requests) == 0 {
		return &pbgroup.GetGroupUsersReqApplicationListResp{}, nil
	}
	groupIDs := datautil.Distinct(datautil.Slice(requests, func(e *model.GroupRequest) string {
		return e.GroupID
	}))
	groups, err := s.db.FindGroup(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	groupMap := datautil.SliceToMap(groups, func(e *model.Group) string {
		return e.GroupID
	})
	if ids := datautil.Single(groupIDs, datautil.Keys(groupMap)); len(ids) > 0 {
		return nil, servererrs.ErrGroupIDNotFound.WrapMsg(strings.Join(ids, ","))
	}
	owners, err := s.db.FindGroupsOwner(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	if err := s.PopulateGroupMember(ctx, owners...); err != nil {
		return nil, err
	}
	ownerMap := datautil.SliceToMap(owners, func(e *model.GroupMember) string {
		return e.GroupID
	})
	groupMemberNum, err := s.db.MapGroupMemberNum(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	return &pbgroup.GetGroupUsersReqApplicationListResp{
		Total: int64(len(requests)),
		GroupRequests: datautil.Slice(requests, func(e *model.GroupRequest) *sdkws.GroupRequest {
			var ownerUserID string
			if owner, ok := ownerMap[e.GroupID]; ok {
				ownerUserID = owner.UserID
			}
			return convert.Db2PbGroupRequest(e, nil, convert.Db2PbGroupInfo(groupMap[e.GroupID], ownerUserID, groupMemberNum[e.GroupID]))
		}),
	}, nil
}
