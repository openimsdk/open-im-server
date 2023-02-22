package localcache

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/discoveryregistry"
	"Open_IM/pkg/proto/group"
	"context"
	"sync"
)

type GroupLocalCacheInterface interface {
	GetGroupMemberIDs(ctx context.Context, groupID string) ([]string, error)
}

type GroupLocalCache struct {
	lock   sync.Mutex
	cache  map[string]GroupMemberIDsHash
	client discoveryregistry.SvcDiscoveryRegistry
}

type GroupMemberIDsHash struct {
	memberListHash uint64
	userIDs        []string
}

func NewGroupMemberIDsLocalCache(client discoveryregistry.SvcDiscoveryRegistry) *GroupLocalCache {
	return &GroupLocalCache{
		cache:  make(map[string]GroupMemberIDsHash, 0),
		client: client,
	}
}

func (g *GroupLocalCache) GetGroupMemberIDs(ctx context.Context, groupID string) ([]string, error) {
	g.lock.Lock()
	defer g.lock.Unlock()
	conn, err := g.client.GetConn(config.Config.RpcRegisterName.OpenImGroupName)
	if err != nil {
		return nil, err
	}
	client := group.NewGroupClient(conn)
	resp, err := client.GetGroupAbstractInfo(ctx, &group.GetGroupAbstractInfoReq{
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
	groupMembersResp, err := client.GetGroupMemberUserID(ctx, &group.GetGroupMemberUserIDReq{
		GroupID: groupID,
	})
	if err != nil {
		return nil, err
	}
	g.cache[groupID] = GroupMemberIDsHash{
		memberListHash: resp.GroupAbstractInfos[0].GroupMemberListHash,
		userIDs:        groupMembersResp.UserIDs,
	}
	return g.cache[groupID].userIDs, nil
}
