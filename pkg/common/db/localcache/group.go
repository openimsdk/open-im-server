package localcache

import (
	"Open_IM/pkg/proto/group"
	"context"
	"google.golang.org/grpc"
)

type GroupLocalCache struct {
	cache map[string]GroupMemberIDsHash
	rpc   *grpc.ClientConn
	group group.GroupClient
}

type GroupMemberIDsHash struct {
	memberListHash uint64
	userIDs        []string
}

func NewGroupMemberIDsLocalCache(rpc *grpc.ClientConn) GroupLocalCache {
	return GroupLocalCache{
		cache: make(map[string]GroupMemberIDsHash, 0),
		rpc:   rpc,
		group: group.NewGroupClient(rpc),
	}
}

func (g *GroupLocalCache) GetGroupMemberIDs(ctx context.Context, groupID string) []string {
	_, err := g.group.GetGroupAbstractInfo(ctx, &group.GetGroupAbstractInfoReq{
		GroupIDs: nil,
	})
	if err != nil {
		return nil
	}
	return []string{}
}
