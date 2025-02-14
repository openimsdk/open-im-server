package mcache

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/open-im-server/v3/pkg/localcache"
	"github.com/openimsdk/open-im-server/v3/pkg/localcache/lru"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/redis/go-redis/v9"
)

var (
	memMsgCache     lru.LRU[string, *model.MsgInfoModel]
	initMemMsgCache sync.Once
)

func NewMsgCache(cache database.Cache, msgDocDatabase database.Msg) cache.MsgCache {
	initMemMsgCache.Do(func() {
		memMsgCache = lru.NewLayLRU[string, *model.MsgInfoModel](1024*8, time.Hour, time.Second*10, localcache.EmptyTarget{}, nil)
	})
	return &msgCache{
		cache:          cache,
		msgDocDatabase: msgDocDatabase,
		memMsgCache:    memMsgCache,
	}
}

type msgCache struct {
	cache          database.Cache
	msgDocDatabase database.Msg
	memMsgCache    lru.LRU[string, *model.MsgInfoModel]
}

func (x *msgCache) getSendMsgKey(id string) string {
	return cachekey.GetSendMsgKey(id)
}

func (x *msgCache) SetSendMsgStatus(ctx context.Context, id string, status int32) error {
	return x.cache.Set(ctx, x.getSendMsgKey(id), strconv.Itoa(int(status)), time.Hour*24)
}

func (x *msgCache) GetSendMsgStatus(ctx context.Context, id string) (int32, error) {
	key := x.getSendMsgKey(id)
	res, err := x.cache.Get(ctx, []string{key})
	if err != nil {
		return 0, err
	}
	val, ok := res[key]
	if !ok {
		return 0, errs.Wrap(redis.Nil)
	}
	status, err := strconv.Atoi(val)
	if err != nil {
		return 0, errs.WrapMsg(err, "GetSendMsgStatus strconv.Atoi error", "val", val)
	}
	return int32(status), nil
}

func (x *msgCache) getMsgCacheKey(conversationID string, seq int64) string {
	return cachekey.GetMsgCacheKey(conversationID, seq)

}

func (x *msgCache) GetMessageBySeqs(ctx context.Context, conversationID string, seqs []int64) ([]*model.MsgInfoModel, error) {
	if len(seqs) == 0 {
		return nil, nil
	}
	keys := make([]string, 0, len(seqs))
	keySeq := make(map[string]int64, len(seqs))
	for _, seq := range seqs {
		key := x.getMsgCacheKey(conversationID, seq)
		keys = append(keys, key)
		keySeq[key] = seq
	}
	res, err := x.memMsgCache.GetBatch(keys, func(keys []string) (map[string]*model.MsgInfoModel, error) {
		findSeqs := make([]int64, 0, len(keys))
		for _, key := range keys {
			seq, ok := keySeq[key]
			if !ok {
				continue
			}
			findSeqs = append(findSeqs, seq)
		}
		res, err := x.msgDocDatabase.FindSeqs(ctx, conversationID, seqs)
		if err != nil {
			return nil, err
		}
		kv := make(map[string]*model.MsgInfoModel)
		for i := range res {
			msg := res[i]
			if msg == nil || msg.Msg == nil || msg.Msg.Seq <= 0 {
				continue
			}
			key := x.getMsgCacheKey(conversationID, msg.Msg.Seq)
			kv[key] = msg
		}
		return kv, nil
	})
	if err != nil {
		return nil, err
	}
	return datautil.Values(res), nil
}

func (x msgCache) DelMessageBySeqs(ctx context.Context, conversationID string, seqs []int64) error {
	if len(seqs) == 0 {
		return nil
	}
	for _, seq := range seqs {
		x.memMsgCache.Del(x.getMsgCacheKey(conversationID, seq))
	}
	return nil
}

func (x *msgCache) SetMessageBySeqs(ctx context.Context, conversationID string, msgs []*model.MsgInfoModel) error {
	for i := range msgs {
		msg := msgs[i]
		if msg == nil || msg.Msg == nil || msg.Msg.Seq <= 0 {
			continue
		}
		x.memMsgCache.Set(x.getMsgCacheKey(conversationID, msg.Msg.Seq), msg)
	}
	return nil
}
