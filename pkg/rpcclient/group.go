package rpcclient

import (
	"context"
	"strings"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/group"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
)

type GroupClient struct {
	MetaClient
}

func NewGroupClient(client discoveryregistry.SvcDiscoveryRegistry) *GroupClient {
	return &GroupClient{
		MetaClient: MetaClient{
			client:          client,
			rpcRegisterName: config.Config.RpcRegisterName.OpenImGroupName,
		},
	}
}

func (g *GroupClient) GetGroupInfos(ctx context.Context, groupIDs []string, complete bool) ([]*sdkws.GroupInfo, error) {
	cc, err := g.getConn(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := group.NewGroupClient(cc).GetGroupsInfo(ctx, &group.GetGroupsInfoReq{
		GroupIDs: groupIDs,
	})
	if err != nil {
		return nil, err
	}
	if complete {
		if ids := utils.Single(groupIDs, utils.Slice(resp.GroupInfos, func(e *sdkws.GroupInfo) string {
			return e.GroupID
		})); len(ids) > 0 {
			return nil, errs.ErrGroupIDNotFound.Wrap(strings.Join(ids, ","))
		}
	}
	return resp.GroupInfos, nil
}

func (g *GroupClient) GetGroupInfo(ctx context.Context, groupID string) (*sdkws.GroupInfo, error) {
	groups, err := g.GetGroupInfos(ctx, []string{groupID}, true)
	if err != nil {
		return nil, err
	}
	return groups[0], nil
}

func (g *GroupClient) GetGroupInfoMap(ctx context.Context, groupIDs []string, complete bool) (map[string]*sdkws.GroupInfo, error) {
	groups, err := g.GetGroupInfos(ctx, groupIDs, complete)
	if err != nil {
		return nil, err
	}
	return utils.SliceToMap(groups, func(e *sdkws.GroupInfo) string {
		return e.GroupID
	}), nil
}

func (g *GroupClient) GetGroupMemberInfos(ctx context.Context, groupID string, userIDs []string, complete bool) ([]*sdkws.GroupMemberFullInfo, error) {
	cc, err := g.getConn(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := group.NewGroupClient(cc).GetGroupMembersInfo(ctx, &group.GetGroupMembersInfoReq{
		GroupID: groupID,
		UserIDs: userIDs,
	})
	if err != nil {
		return nil, err
	}
	if complete {
		if ids := utils.Single(userIDs, utils.Slice(resp.Members, func(e *sdkws.GroupMemberFullInfo) string {
			return e.UserID
		})); len(ids) > 0 {
			return nil, errs.ErrNotInGroupYet.Wrap(strings.Join(ids, ","))
		}
	}
	return resp.Members, nil
}

func (g *GroupClient) GetGroupMemberInfo(ctx context.Context, groupID string, userID string) (*sdkws.GroupMemberFullInfo, error) {
	members, err := g.GetGroupMemberInfos(ctx, groupID, []string{userID}, true)
	if err != nil {
		return nil, err
	}
	return members[0], nil
}

func (g *GroupClient) GetGroupMemberInfoMap(ctx context.Context, groupID string, userIDs []string, complete bool) (map[string]*sdkws.GroupMemberFullInfo, error) {
	members, err := g.GetGroupMemberInfos(ctx, groupID, userIDs, true)
	if err != nil {
		return nil, err
	}
	return utils.SliceToMap(members, func(e *sdkws.GroupMemberFullInfo) string {
		return e.UserID
	}), nil
}

func (g *GroupClient) GetOwnerAndAdminInfos(ctx context.Context, groupID string) ([]*sdkws.GroupMemberFullInfo, error) {
	cc, err := g.getConn(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := group.NewGroupClient(cc).GetGroupMemberRoleLevel(ctx, &group.GetGroupMemberRoleLevelReq{
		GroupID:    groupID,
		RoleLevels: []int32{constant.GroupOwner, constant.GroupAdmin},
	})
	return resp.Members, err
}

func (g *GroupClient) GetOwnerInfo(ctx context.Context, groupID string) (*sdkws.GroupMemberFullInfo, error) {
	cc, err := g.getConn(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := group.NewGroupClient(cc).GetGroupMemberRoleLevel(ctx, &group.GetGroupMemberRoleLevelReq{
		GroupID:    groupID,
		RoleLevels: []int32{constant.GroupOwner},
	})
	return resp.Members[0], err
}

func (g *GroupClient) GetGroupMemberIDs(ctx context.Context, groupID string) ([]string, error) {
	cc, err := g.getConn(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := group.NewGroupClient(cc).GetGroupMemberUserIDs(ctx, &group.GetGroupMemberUserIDsReq{
		GroupID: groupID,
	})
	return resp.UserIDs, err
}
