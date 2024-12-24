// Copyright © 2024 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rpccache

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcli"
	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/tools/utils/datautil"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/localcache"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/redis/go-redis/v9"
)

func NewGroupLocalCache(client *rpcli.GroupClient, localCache *config.LocalCache, cli redis.UniversalClient) *GroupLocalCache {
	lc := localCache.Group
	log.ZDebug(context.Background(), "GroupLocalCache", "topic", lc.Topic, "slotNum", lc.SlotNum, "slotSize", lc.SlotSize, "enable", lc.Enable())
	x := &GroupLocalCache{
		client: client,
		local: localcache.New[[]byte](
			localcache.WithLocalSlotNum(lc.SlotNum),
			localcache.WithLocalSlotSize(lc.SlotSize),
			localcache.WithLinkSlotNum(lc.SlotNum),
			localcache.WithLocalSuccessTTL(lc.Success()),
			localcache.WithLocalFailedTTL(lc.Failed()),
		),
	}
	if lc.Enable() {
		go subscriberRedisDeleteCache(context.Background(), cli, lc.Topic, x.local.DelLocal)
	}
	return x
}

type GroupLocalCache struct {
	client *rpcli.GroupClient
	local  localcache.Cache[[]byte]
}

func (g *GroupLocalCache) getGroupMemberIDs(ctx context.Context, groupID string) (val *group.GetGroupMemberUserIDsResp, err error) {
	log.ZDebug(ctx, "GroupLocalCache getGroupMemberIDs req", "groupID", groupID)
	defer func() {
		if err == nil {
			log.ZDebug(ctx, "GroupLocalCache getGroupMemberIDs return", "groupID", groupID, "value", val)
		} else {
			log.ZError(ctx, "GroupLocalCache getGroupMemberIDs return", err, "groupID", groupID)
		}
	}()
	var cache cacheProto[group.GetGroupMemberUserIDsResp]
	return cache.Unmarshal(g.local.Get(ctx, cachekey.GetGroupMemberIDsKey(groupID), func(ctx context.Context) ([]byte, error) {
		log.ZDebug(ctx, "GroupLocalCache getGroupMemberIDs rpc", "groupID", groupID)
		return cache.Marshal(g.client.GroupClient.GetGroupMemberUserIDs(ctx, &group.GetGroupMemberUserIDsReq{GroupID: groupID}))
	}))
}

func (g *GroupLocalCache) GetGroupMember(ctx context.Context, groupID, userID string) (val *sdkws.GroupMemberFullInfo, err error) {
	log.ZDebug(ctx, "GroupLocalCache GetGroupInfo req", "groupID", groupID, "userID", userID)
	defer func() {
		if err == nil {
			log.ZDebug(ctx, "GroupLocalCache GetGroupInfo return", "groupID", groupID, "userID", userID, "value", val)
		} else {
			log.ZError(ctx, "GroupLocalCache GetGroupInfo return", err, "groupID", groupID, "userID", userID)
		}
	}()
	var cache cacheProto[sdkws.GroupMemberFullInfo]
	return cache.Unmarshal(g.local.Get(ctx, cachekey.GetGroupMemberInfoKey(groupID, userID), func(ctx context.Context) ([]byte, error) {
		log.ZDebug(ctx, "GroupLocalCache GetGroupInfo rpc", "groupID", groupID, "userID", userID)
		return cache.Marshal(g.client.GetGroupMemberCache(ctx, groupID, userID))
	}))
}

func (g *GroupLocalCache) GetGroupInfo(ctx context.Context, groupID string) (val *sdkws.GroupInfo, err error) {
	log.ZDebug(ctx, "GroupLocalCache GetGroupInfo req", "groupID", groupID)
	defer func() {
		if err == nil {
			log.ZDebug(ctx, "GroupLocalCache GetGroupInfo return", "groupID", groupID, "value", val)
		} else {
			log.ZError(ctx, "GroupLocalCache GetGroupInfo return", err, "groupID", groupID)
		}
	}()
	var cache cacheProto[sdkws.GroupInfo]
	return cache.Unmarshal(g.local.Get(ctx, cachekey.GetGroupInfoKey(groupID), func(ctx context.Context) ([]byte, error) {
		log.ZDebug(ctx, "GroupLocalCache GetGroupInfo rpc", "groupID", groupID)
		return cache.Marshal(g.client.GetGroupInfoCache(ctx, groupID))
	}))
}

func (g *GroupLocalCache) GetGroupMemberIDs(ctx context.Context, groupID string) ([]string, error) {
	res, err := g.getGroupMemberIDs(ctx, groupID)
	if err != nil {
		return nil, err
	}
	return res.UserIDs, nil
}

func (g *GroupLocalCache) GetGroupMemberIDMap(ctx context.Context, groupID string) (map[string]struct{}, error) {
	res, err := g.getGroupMemberIDs(ctx, groupID)
	if err != nil {
		return nil, err
	}
	return datautil.SliceSet(res.UserIDs), nil
}

func (g *GroupLocalCache) GetGroupInfos(ctx context.Context, groupIDs []string) ([]*sdkws.GroupInfo, error) {
	groupInfos := make([]*sdkws.GroupInfo, 0, len(groupIDs))
	for _, groupID := range groupIDs {
		groupInfo, err := g.GetGroupInfo(ctx, groupID)
		if err != nil {
			if errs.ErrRecordNotFound.Is(err) {
				continue
			}
			return nil, err
		}
		groupInfos = append(groupInfos, groupInfo)
	}
	return groupInfos, nil
}

func (g *GroupLocalCache) GetGroupMembers(ctx context.Context, groupID string, userIDs []string) ([]*sdkws.GroupMemberFullInfo, error) {
	members := make([]*sdkws.GroupMemberFullInfo, 0, len(userIDs))
	for _, userID := range userIDs {
		member, err := g.GetGroupMember(ctx, groupID, userID)
		if err != nil {
			if errs.ErrRecordNotFound.Is(err) {
				continue
			}
			return nil, err
		}
		members = append(members, member)
	}
	return members, nil
}

func (g *GroupLocalCache) GetGroupMemberInfoMap(ctx context.Context, groupID string, userIDs []string) (map[string]*sdkws.GroupMemberFullInfo, error) {
	members := make(map[string]*sdkws.GroupMemberFullInfo)
	for _, userID := range userIDs {
		member, err := g.GetGroupMember(ctx, groupID, userID)
		if err != nil {
			if errs.ErrRecordNotFound.Is(err) {
				continue
			}
			return nil, err
		}
		members[userID] = member
	}
	return members, nil
}
