package localcache

import (
	"Open_IM/pkg/proto/group"
	"context"
	"google.golang.org/grpc"
	"sync"
)

type GroupLocalCache struct {
	lock  sync.Mutex
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
	resp, err := g.group.GetGroupAbstractInfo(ctx, &group.GetGroupAbstractInfoReq{
		GroupIDs: nil,
	})
	if err != nil {
		return nil
	}
	return []string{}
}
