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
	"google.golang.org/grpc"
)

type Group struct {
	conn   grpc.ClientConnInterface
	Client group.GroupClient
	discov discoveryregistry.SvcDiscoveryRegistry
}

func NewGroup(discov discoveryregistry.SvcDiscoveryRegistry) *Group {
	conn, err := discov.GetConn(context.Background(), config.Config.RpcRegisterName.OpenImGroupName)
	if err != nil {
		panic(err)
	}
	client := group.NewGroupClient(conn)
	return &Group{discov: discov, conn: conn, Client: client}
}

type GroupRpcClient Group

func NewGroupRpcClient(discov discoveryregistry.SvcDiscoveryRegistry) GroupRpcClient {
	return GroupRpcClient(*NewGroup(discov))
}

func (g *GroupRpcClient) GetGroupInfos(ctx context.Context, groupIDs []string, complete bool) ([]*sdkws.GroupInfo, error) {
	resp, err := g.Client.GetGroupsInfo(ctx, &group.GetGroupsInfoReq{
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

func (g *GroupRpcClient) GetGroupInfo(ctx context.Context, groupID string) (*sdkws.GroupInfo, error) {
	groups, err := g.GetGroupInfos(ctx, []string{groupID}, true)
	if err != nil {
		return nil, err
	}
	return groups[0], nil
}

func (g *GroupRpcClient) GetGroupInfoMap(ctx context.Context, groupIDs []string, complete bool) (map[string]*sdkws.GroupInfo, error) {
	groups, err := g.GetGroupInfos(ctx, groupIDs, complete)
	if err != nil {
		return nil, err
	}
	return utils.SliceToMap(groups, func(e *sdkws.GroupInfo) string {
		return e.GroupID
	}), nil
}

func (g *GroupRpcClient) GetGroupMemberInfos(ctx context.Context, groupID string, userIDs []string, complete bool) ([]*sdkws.GroupMemberFullInfo, error) {
	resp, err := g.Client.GetGroupMembersInfo(ctx, &group.GetGroupMembersInfoReq{
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

func (g *GroupRpcClient) GetGroupMemberInfo(ctx context.Context, groupID string, userID string) (*sdkws.GroupMemberFullInfo, error) {
	members, err := g.GetGroupMemberInfos(ctx, groupID, []string{userID}, true)
	if err != nil {
		return nil, err
	}
	return members[0], nil
}

func (g *GroupRpcClient) GetGroupMemberInfoMap(ctx context.Context, groupID string, userIDs []string, complete bool) (map[string]*sdkws.GroupMemberFullInfo, error) {
	members, err := g.GetGroupMemberInfos(ctx, groupID, userIDs, true)
	if err != nil {
		return nil, err
	}
	return utils.SliceToMap(members, func(e *sdkws.GroupMemberFullInfo) string {
		return e.UserID
	}), nil
}

func (g *GroupRpcClient) GetOwnerAndAdminInfos(ctx context.Context, groupID string) ([]*sdkws.GroupMemberFullInfo, error) {
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
