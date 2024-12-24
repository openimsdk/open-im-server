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

package rpcclient

import (
	"context"
	"strings"

	"github.com/openimsdk/open-im-server/v3/pkg/common/servererrs"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/system/program"
	"github.com/openimsdk/tools/utils/datautil"
)

type Group struct {
	Client group.GroupClient
	discov discovery.SvcDiscoveryRegistry
}

func NewGroup(discov discovery.SvcDiscoveryRegistry, rpcRegisterName string) *Group {
	conn, err := discov.GetConn(context.Background(), rpcRegisterName)
	if err != nil {
		program.ExitWithError(err)
	}
	client := group.NewGroupClient(conn)
	return &Group{discov: discov, Client: client}
}

type GroupRpcClient Group

func NewGroupRpcClient(discov discovery.SvcDiscoveryRegistry, rpcRegisterName string) GroupRpcClient {
	return GroupRpcClient(*NewGroup(discov, rpcRegisterName))
}

func (g *GroupRpcClient) GetGroupInfos(ctx context.Context, groupIDs []string, complete bool) ([]*sdkws.GroupInfo, error) {
	resp, err := g.Client.GetGroupsInfo(ctx, &group.GetGroupsInfoReq{
		GroupIDs: groupIDs,
	})
	if err != nil {
		return nil, err
	}
	if complete {
		if ids := datautil.Single(groupIDs, datautil.Slice(resp.GroupInfos, func(e *sdkws.GroupInfo) string {
			return e.GroupID
		})); len(ids) > 0 {
			return nil, servererrs.ErrGroupIDNotFound.WrapMsg(strings.Join(ids, ","))
		}
	}
	return resp.GroupInfos, nil
}

func (g *GroupRpcClient) GetGroupInfo(ctx context.Context, groupID string) (*sdkws.GroupInfo, error) {
	groups, err := g.GetGroupInfos(ctx, []string{groupID}, true)
	if err != nil {
		return nil, err
	}
	return groups[0], nil
}

func (g *GroupRpcClient) GetGroupInfoMap(
	ctx context.Context,
	groupIDs []string,
	complete bool,
) (map[string]*sdkws.GroupInfo, error) {
	groups, err := g.GetGroupInfos(ctx, groupIDs, complete)
	if err != nil {
		return nil, err
	}
	return datautil.SliceToMap(groups, func(e *sdkws.GroupInfo) string {
		return e.GroupID
	}), nil
}

func (g *GroupRpcClient) GetGroupMemberInfos(
	ctx context.Context,
	groupID string,
	userIDs []string,
	complete bool,
) ([]*sdkws.GroupMemberFullInfo, error) {
	resp, err := g.Client.GetGroupMembersInfo(ctx, &group.GetGroupMembersInfoReq{
		GroupID: groupID,
		UserIDs: userIDs,
	})
	if err != nil {
		return nil, err
	}
	if complete {
		if ids := datautil.Single(userIDs, datautil.Slice(resp.Members, func(e *sdkws.GroupMemberFullInfo) string {
			return e.UserID
		})); len(ids) > 0 {
			return nil, servererrs.ErrNotInGroupYet.WrapMsg(strings.Join(ids, ","))
		}
	}
	return resp.Members, nil
}

func (g *GroupRpcClient) GetGroupMemberInfo(
	ctx context.Context,
	groupID string,
	userID string,
) (*sdkws.GroupMemberFullInfo, error) {
	members, err := g.GetGroupMemberInfos(ctx, groupID, []string{userID}, true)
	if err != nil {
		return nil, err
	}
	return members[0], nil
}

func (g *GroupRpcClient) GetGroupMemberInfoMap(
	ctx context.Context,
	groupID string,
	userIDs []string,
	complete bool,
) (map[string]*sdkws.GroupMemberFullInfo, error) {
	members, err := g.GetGroupMemberInfos(ctx, groupID, userIDs, true)
	if err != nil {
		return nil, err
	}
	return datautil.SliceToMap(members, func(e *sdkws.GroupMemberFullInfo) string {
		return e.UserID
	}), nil
}

func (g *GroupRpcClient) GetOwnerAndAdminInfos(
	ctx context.Context,
	groupID string,
) ([]*sdkws.GroupMemberFullInfo, error) {
	resp, err := g.Client.GetGroupMemberRoleLevel(ctx, &group.GetGroupMemberRoleLevelReq{
		GroupID:    groupID,
		RoleLevels: []int32{constant.GroupOwner, constant.GroupAdmin},
	})
	if err != nil {
		return nil, err
	}
	return resp.Members, nil
}

func (g *GroupRpcClient) GetOwnerInfo(ctx context.Context, groupID string) (*sdkws.GroupMemberFullInfo, error) {
	resp, err := g.Client.GetGroupMemberRoleLevel(ctx, &group.GetGroupMemberRoleLevelReq{
		GroupID:    groupID,
		RoleLevels: []int32{constant.GroupOwner},
	})
	return resp.Members[0], err
}

func (g *GroupRpcClient) GetGroupMemberIDs(ctx context.Context, groupID string) ([]string, error) {
	resp, err := g.Client.GetGroupMemberUserIDs(ctx, &group.GetGroupMemberUserIDsReq{
		GroupID: groupID,
	})
	if err != nil {
		return nil, err
	}
	return resp.UserIDs, nil
}

func (g *GroupRpcClient) GetGroupInfoCache(ctx context.Context, groupID string) (*sdkws.GroupInfo, error) {
	resp, err := g.Client.GetGroupInfoCache(ctx, &group.GetGroupInfoCacheReq{
		GroupID: groupID,
	})
	if err != nil {
		return nil, err
	}
	return resp.GroupInfo, nil
}

func (g *GroupRpcClient) GetGroupMemberCache(ctx context.Context, groupID string, groupMemberID string) (*sdkws.GroupMemberFullInfo, error) {
	resp, err := g.Client.GetGroupMemberCache(ctx, &group.GetGroupMemberCacheReq{
		GroupID:       groupID,
		GroupMemberID: groupMemberID,
	})
	if err != nil {
		return nil, err
	}
	return resp.Member, nil
}

func (g *GroupRpcClient) DismissGroup(ctx context.Context, groupID string) error {
	_, err := g.Client.DismissGroup(ctx, &group.DismissGroupReq{
		GroupID:      groupID,
		DeleteMember: true,
	})
	return err
}

func (g *GroupRpcClient) NotificationUserInfoUpdate(ctx context.Context, userID string) error {
	_, err := g.Client.NotificationUserInfoUpdate(ctx, &group.NotificationUserInfoUpdateReq{
		UserID: userID,
	})
	return err
}
