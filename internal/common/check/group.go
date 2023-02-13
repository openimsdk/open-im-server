package check

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	discoveryRegistry "Open_IM/pkg/discoveryregistry"
	"Open_IM/pkg/proto/group"
	sdkws "Open_IM/pkg/proto/sdkws"
	"Open_IM/pkg/utils"
	"context"
	"google.golang.org/grpc"
	"strings"
)

type GroupChecker struct {
	zk discoveryRegistry.SvcDiscoveryRegistry
}

func NewGroupChecker(zk discoveryRegistry.SvcDiscoveryRegistry) *GroupChecker {
	return &GroupChecker{
		zk: zk,
	}
}

func (g *GroupChecker) getConn() (*grpc.ClientConn, error) {
	return g.zk.GetConn(config.Config.RpcRegisterName.OpenImGroupName)
}

func (g *GroupChecker) GetGroupInfos(ctx context.Context, groupIDs []string, complete bool) ([]*sdkws.GroupInfo, error) {
	cc, err := g.getConn()
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
			return nil, constant.ErrGroupIDNotFound.Wrap(strings.Join(ids, ","))
		}
	}
	return resp.GroupInfos, nil
}

func (g *GroupChecker) GetGroupInfo(ctx context.Context, groupID string) (*sdkws.GroupInfo, error) {
	groups, err := g.GetGroupInfos(ctx, []string{groupID}, true)
	if err != nil {
		return nil, err
	}
	return groups[0], nil
}

func (g *GroupChecker) GetGroupInfoMap(ctx context.Context, groupIDs []string, complete bool) (map[string]*sdkws.GroupInfo, error) {
	groups, err := g.GetGroupInfos(ctx, groupIDs, complete)
	if err != nil {
		return nil, err
	}
	return utils.SliceToMap(groups, func(e *sdkws.GroupInfo) string {
		return e.GroupID
	}), nil
}

func (g *GroupChecker) GetGroupMemberInfos(ctx context.Context, groupID string, userIDs []string, complete bool) ([]*sdkws.GroupMemberFullInfo, error) {
	cc, err := g.getConn()
	if err != nil {
		return nil, err
	}
	resp, err := group.NewGroupClient(cc).GetGroupMembersInfo(ctx, &group.GetGroupMembersInfoReq{
		GroupID: groupID,
		Members: userIDs,
	})
	if err != nil {
		return nil, err
	}
	if complete {
		if ids := utils.Single(userIDs, utils.Slice(resp.Members, func(e *sdkws.GroupMemberFullInfo) string {
			return e.UserID
		})); len(ids) > 0 {
			return nil, constant.ErrNotInGroupYet.Wrap(strings.Join(ids, ","))
		}
	}
	return resp.Members, nil
}

func (g *GroupChecker) GetGroupMemberInfo(ctx context.Context, groupID string, userID string) (*sdkws.GroupMemberFullInfo, error) {
	members, err := g.GetGroupMemberInfos(ctx, groupID, []string{userID}, true)
	if err != nil {
		return nil, err
	}
	return members[0], nil
}

func (g *GroupChecker) GetGroupMemberInfoMap(ctx context.Context, groupID string, userIDs []string, complete bool) (map[string]*sdkws.GroupMemberFullInfo, error) {
	members, err := g.GetGroupMemberInfos(ctx, groupID, userIDs, true)
	if err != nil {
		return nil, err
	}
	return utils.SliceToMap(members, func(e *sdkws.GroupMemberFullInfo) string {
		return e.UserID
	}), nil
}
