package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/redis/go-redis/v9"
) //

// msgCacheTimeout is  expiration time of message cache, 86400 seconds
const msgCacheTimeout = time.Hour * 24

func NewMsgCache(client redis.UniversalClient, db database.Msg) cache.MsgCache {
	return &msgCache{
		rcClient:       newRocksCacheClient(client),
		msgDocDatabase: db,
	}
}

type msgCache struct {
	rcClient       *rocksCacheClient
	msgDocDatabase database.Msg
}

func (c *msgCache) getSendMsgKey(id string) string {
	return cachekey.GetSendMsgKey(id)
}

func (c *msgCache) SetSendMsgStatus(ctx context.Context, id string, status int32) error {
	return errs.Wrap(c.rcClient.GetRedis().Set(ctx, c.getSendMsgKey(id), status, time.Hour*24).Err())
}

func (c *msgCache) GetSendMsgStatus(ctx context.Context, id string) (int32, error) {
	result, err := c.rcClient.GetRedis().Get(ctx, c.getSendMsgKey(id)).Int()
	return int32(result), errs.Wrap(err)
}

func (c *msgCache) GetMessageBySeqs(ctx context.Context, conversationID string, seqs []int64) ([]*model.MsgInfoModel, error) {
	if len(seqs) == 0 {
		return nil, nil
	}
	getKey := func(seq int64) string {
		return cachekey.GetMsgCacheKey(conversationID, seq)
	}
	getMsgID := func(msg *model.MsgInfoModel) int64 {
		return msg.Msg.Seq
	}
	find := func(ctx context.Context, seqs []int64) ([]*model.MsgInfoModel, error) {
		return c.msgDocDatabase.FindSeqs(ctx, conversationID, seqs)
	}
	return batchGetCache2(ctx, c.rcClient, msgCacheTimeout, seqs, getKey, getMsgID, find)
}

func (c *msgCache) DelMessageBySeqs(ctx context.Context, conversationID string, seqs []int64) error {
	if len(seqs) == 0 {
		return nil
	}
	keys := datautil.Slice(seqs, func(seq int64) string {
		return cachekey.GetMsgCacheKey(conversationID, seq)
	})
	slotKeys, err := groupKeysBySlot(ctx, c.rcClient.GetRedis(), keys)
	if err != nil {
		return err
	}
	for _, keys := range slotKeys {
		if err := c.rcClient.GetClient().TagAsDeletedBatch2(ctx, keys); err != nil {
			return err
		}
	}
	return nil
}

func (c *msgCache) SetMessageBySeqs(ctx context.Context, conversationID string, msgs []*model.MsgInfoModel) error {
	for _, msg := range msgs {
		if msg == nil || msg.Msg == nil || msg.Msg.Seq <= 0 {
			continue
		}
		data, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		if err := c.rcClient.GetClient().RawSet(ctx, cachekey.GetMsgCacheKey(conversationID, msg.Msg.Seq), string(data), msgCacheTimeout); err != nil {
			return err
		}
	}
	return nil
}
