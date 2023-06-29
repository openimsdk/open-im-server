package localcache

import (
	"context"
	"sync"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/group"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
)

type GroupLocalCache struct {
	lock   sync.Mutex
	cache  map[string]GroupMemberIDsHash
	client *rpcclient.Group
}

type GroupMemberIDsHash struct {
	memberListHash uint64
	userIDs        []string
}

func NewGroupLocalCache(discov discoveryregistry.SvcDiscoveryRegistry) *GroupLocalCache {
	client := rpcclient.NewGroup(discov)
	return &GroupLocalCache{
		cache:  make(map[string]GroupMemberIDsHash, 0),
		client: client,
	}
}

func (g *GroupLocalCache) GetGroupMemberIDs(ctx context.Context, groupID string) ([]string, error) {
	resp, err := g.client.Client.GetGroupAbstractInfo(ctx, &group.GetGroupAbstractInfoReq{
		GroupIDs: []string{groupID},
	})
	if err != nil {
		return nil, err
	}
	if len(resp.GroupAbstractInfos) < 1 {
		return nil, errs.ErrGroupIDNotFound
	}
	g.lock.Lock()
	defer g.lock.Unlock()
	localHashInfo, ok := g.cache[groupID]
	if ok && localHashInfo.memberListHash == resp.GroupAbstractInfos[0].GroupMemberListHash {
		return localHashInfo.userIDs, nil
	}
	groupMembersResp, err := g.client.Client.GetGroupMemberUserIDs(ctx, &group.GetGroupMemberUserIDsReq{
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
