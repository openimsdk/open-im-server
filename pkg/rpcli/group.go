package rpcli

import (
	"context"
	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/sdkws"
)

func NewGroupClient(cli group.GroupClient) *GroupClient {
	return &GroupClient{cli}
}

type GroupClient struct {
	group.GroupClient
}

func (x *GroupClient) cli() group.GroupClient {
	return x.GroupClient
}

func (x *GroupClient) GetGroupsInfo(ctx context.Context, groupIDs []string) ([]*sdkws.GroupInfo, error) {
	req := &group.GetGroupsInfoReq{GroupIDs: groupIDs}
	return extractField(ctx, x.cli().GetGroupsInfo, req, (*group.GetGroupsInfoResp).GetGroupInfos)
}

func (x *GroupClient) GetGroupInfo(ctx context.Context, groupID string) (*sdkws.GroupInfo, error) {
	return firstValue(x.GetGroupsInfo(ctx, []string{groupID}))
}
