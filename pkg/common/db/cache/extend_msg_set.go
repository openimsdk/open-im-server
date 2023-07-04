// Copyright © 2023 OpenIM. All rights reserved.
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

package cache

import (
	"context"
	"time"

	"github.com/dtm-labs/rockscache"
	"github.com/redis/go-redis/v9"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/unrelation"
)

const (
	extendMsgSetCache = "EXTEND_MSG_SET_CACHE:"
	extendMsgCache    = "EXTEND_MSG_CACHE:"
)

type ExtendMsgSetCache interface {
	metaCache
	NewCache() ExtendMsgSetCache
	GetExtendMsg(
		ctx context.Context,
		conversationID string,
		sessionType int32,
		clientMsgID string,
		firstModifyTime int64,
	) (extendMsg *unrelation.ExtendMsgModel, err error)
	DelExtendMsg(clientMsgID string) ExtendMsgSetCache
}

type ExtendMsgSetCacheRedis struct {
	metaCache
	expireTime     time.Duration
	rcClient       *rockscache.Client
	extendMsgSetDB unrelation.ExtendMsgSetModelInterface
}

func NewExtendMsgSetCacheRedis(
	rdb redis.UniversalClient,
	extendMsgSetDB unrelation.ExtendMsgSetModelInterface,
	options rockscache.Options,
) ExtendMsgSetCache {
	rcClient := rockscache.NewClient(rdb, options)
	return &ExtendMsgSetCacheRedis{
		metaCache:      NewMetaCacheRedis(rcClient),
		expireTime:     time.Second * 30 * 60,
		extendMsgSetDB: extendMsgSetDB,
		rcClient:       rcClient,
	}
}

func (e *ExtendMsgSetCacheRedis) NewCache() ExtendMsgSetCache {
	return &ExtendMsgSetCacheRedis{
		metaCache:      NewMetaCacheRedis(e.rcClient, e.metaCache.GetPreDelKeys()...),
		expireTime:     e.expireTime,
		extendMsgSetDB: e.extendMsgSetDB,
		rcClient:       e.rcClient,
	}
}

func (e *ExtendMsgSetCacheRedis) getKey(clientMsgID string) string {
	return extendMsgCache + clientMsgID
}

func (e *ExtendMsgSetCacheRedis) GetExtendMsg(
	ctx context.Context,
	conversationID string,
	sessionType int32,
	clientMsgID string,
	firstModifyTime int64,
) (extendMsg *unrelation.ExtendMsgModel, err error) {
	return getCache(
		ctx,
		e.rcClient,
		e.getKey(clientMsgID),
		e.expireTime,
		func(ctx context.Context) (*unrelation.ExtendMsgModel, error) {
			return e.extendMsgSetDB.TakeExtendMsg(ctx, conversationID, sessionType, clientMsgID, firstModifyTime)
		},
	)
}

func (e *ExtendMsgSetCacheRedis) DelExtendMsg(clientMsgID string) ExtendMsgSetCache {
	new := e.NewCache()
	new.AddKeys(e.getKey(clientMsgID))
	return new
}
