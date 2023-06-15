package cache

import (
	"context"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/unrelation"
	"github.com/dtm-labs/rockscache"
	"github.com/redis/go-redis/v9"
)

const (
	extendMsgSetCache = "EXTEND_MSG_SET_CACHE:"
	extendMsgCache    = "EXTEND_MSG_CACHE:"
)

type ExtendMsgSetCache interface {
	metaCache
	NewCache() ExtendMsgSetCache
	GetExtendMsg(ctx context.Context, conversationID string, sessionType int32, clientMsgID string, firstModifyTime int64) (extendMsg *unrelation.ExtendMsgModel, err error)
	DelExtendMsg(clientMsgID string) ExtendMsgSetCache
}

type ExtendMsgSetCacheRedis struct {
	metaCache
	expireTime     time.Duration
	rcClient       *rockscache.Client
	extendMsgSetDB unrelation.ExtendMsgSetModelInterface
}

func NewExtendMsgSetCacheRedis(rdb redis.UniversalClient, extendMsgSetDB unrelation.ExtendMsgSetModelInterface, options rockscache.Options) ExtendMsgSetCache {
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

func (e *ExtendMsgSetCacheRedis) GetExtendMsg(ctx context.Context, conversationID string, sessionType int32, clientMsgID string, firstModifyTime int64) (extendMsg *unrelation.ExtendMsgModel, err error) {
	return getCache(ctx, e.rcClient, e.getKey(clientMsgID), e.expireTime, func(ctx context.Context) (*unrelation.ExtendMsgModel, error) {
		return e.extendMsgSetDB.TakeExtendMsg(ctx, conversationID, sessionType, clientMsgID, firstModifyTime)
	})
}

func (e *ExtendMsgSetCacheRedis) DelExtendMsg(clientMsgID string) ExtendMsgSetCache {
	new := e.NewCache()
	new.AddKeys(e.getKey(clientMsgID))
	return new
}
