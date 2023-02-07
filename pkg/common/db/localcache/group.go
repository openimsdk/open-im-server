package localcache

import (
	"Open_IM/pkg/common/constant"
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

func (g *GroupLocalCache) GetGroupMemberIDs(ctx context.Context, groupID string) ([]string, error) {
	g.lock.Lock()
	defer g.lock.Unlock()
	resp, err := g.group.GetGroupAbstractInfo(ctx, &group.GetGroupAbstractInfoReq{
		GroupIDs: []string{groupID},
	})
	if err != nil {
		return nil, err
	}
	if len(resp.GroupAbstractInfos) < 0 {
		return nil, constant.ErrGroupIDNotFound
	}
	localHashInfo, ok := g.cache[groupID]
	if ok && localHashInfo.memberListHash == resp.GroupAbstractInfos[0].GroupMemberListHash {
		return localHashInfo.userIDs, nil
	}
	groupMembersResp, err := g.group.GetGroupMemberList(ctx, &group.GetGroupMemberListReq{
		GroupID: groupID,
	})
	if err != nil {
		return nil, err
	}
	g.cache[groupID] = GroupMemberIDsHash{
		memberListHash: resp.GroupAbstractInfos[0].GroupMemberListHash,
		userIDs:        groupMembersResp.Members,
	}
	return g.cache[groupID].userIDs, nil
}
