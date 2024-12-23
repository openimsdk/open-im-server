package rpcli

import (
	"context"
	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/sdkws"
	"google.golang.org/grpc"
)

func NewGroupClient(cc grpc.ClientConnInterface) *GroupClient {
	return &GroupClient{group.NewGroupClient(cc)}
}

type GroupClient struct {
	group.GroupClient
}

func (x *GroupClient) GetGroupsInfo(ctx context.Context, groupIDs []string) ([]*sdkws.GroupInfo, error) {
	req := &group.GetGroupsInfoReq{GroupIDs: groupIDs}
	return extractField(ctx, x.GroupClient.GetGroupsInfo, req, (*group.GetGroupsInfoResp).GetGroupInfos)
}

func (x *GroupClient) GetGroupInfo(ctx context.Context, groupID string) (*sdkws.GroupInfo, error) {
	return firstValue(x.GetGroupsInfo(ctx, []string{groupID}))
}

func (x *GroupClient) GetGroupInfoCache(ctx context.Context, groupID string) (*sdkws.GroupInfo, error) {
	req := &group.GetGroupInfoCacheReq{GroupID: groupID}
	return extractField(ctx, x.GroupClient.GetGroupInfoCache, req, (*group.GetGroupInfoCacheResp).GetGroupInfo)
}

func (x *GroupClient) GetGroupMemberCache(ctx context.Context, groupID string, userID string) (*sdkws.GroupMemberFullInfo, error) {
	req := &group.GetGroupMemberCacheReq{GroupID: groupID}
	return extractField(ctx, x.GroupClient.GetGroupMemberCache, req, (*group.GetGroupMemberCacheResp).GetMember)
}

func (x *GroupClient) DismissGroup(ctx context.Context, groupID string, deleteMember bool) error {
	req := &group.DismissGroupReq{GroupID: groupID, DeleteMember: deleteMember}
	return ignoreResp(x.GroupClient.DismissGroup(ctx, req))
}

func (x *GroupClient) GetGroupMemberUserIDs(ctx context.Context, groupID string) ([]string, error) {
	req := &group.GetGroupMemberUserIDsReq{GroupID: groupID}
	return extractField(ctx, x.GroupClient.GetGroupMemberUserIDs, req, (*group.GetGroupMemberUserIDsResp).GetUserIDs)
}
