package redis

import (
	"context"
	"github.com/dtm-labs/rockscache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/open-im-server/v3/pkg/msgprocessor"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/redis/go-redis/v9"
	"time"
) //

// msgCacheTimeout is  expiration time of message cache, 86400 seconds
const msgCacheTimeout = time.Hour * 24

func NewMsgCache(client redis.UniversalClient) cache.MsgCache {
	return &msgCache{rdb: client}
}

type msgCache struct {
	rdb            redis.UniversalClient
	rcClient       *rockscache.Client
	msgDocDatabase database.Msg
	msgTable       model.MsgDocModel
}

func (c *msgCache) getMessageCacheKey(conversationID string, seq int64) string {
	return cachekey.GetMessageCacheKey(conversationID, seq)
}

func (c *msgCache) getSendMsgKey(id string) string {
	return cachekey.GetSendMsgKey(id)
}

func (c *msgCache) SetMessagesToCache(ctx context.Context, conversationID string, msgs []*sdkws.MsgData) (int, error) {
	msgMap := datautil.SliceToMap(msgs, func(msg *sdkws.MsgData) string {
		return c.getMessageCacheKey(conversationID, msg.Seq)
	})
	keys := datautil.Slice(msgs, func(msg *sdkws.MsgData) string {
		return c.getMessageCacheKey(conversationID, msg.Seq)
	})
	err := ProcessKeysBySlot(ctx, c.rdb, keys, func(ctx context.Context, slot int64, keys []string) error {
		var values []string
		for _, key := range keys {
			if msg, ok := msgMap[key]; ok {
				s, err := msgprocessor.Pb2String(msg)
				if err != nil {
					return err
				}
				values = append(values, s)
			}
		}
		return LuaSetBatchWithCommonExpire(ctx, c.rdb, keys, values, int(msgCacheTimeout/time.Second))
	})
	if err != nil {
		return 0, err
	}
	return len(msgs), nil
}

func (c *msgCache) DeleteMessagesFromCache(ctx context.Context, conversationID string, seqs []int64) error {
	var keys []string
	for _, seq := range seqs {
		keys = append(keys, c.getMessageCacheKey(conversationID, seq))
	}
	return ProcessKeysBySlot(ctx, c.rdb, keys, func(ctx context.Context, slot int64, keys []string) error {
		return LuaDeleteBatch(ctx, c.rdb, keys)
	})
}

func (c *msgCache) SetSendMsgStatus(ctx context.Context, id string, status int32) error {
	return errs.Wrap(c.rdb.Set(ctx, c.getSendMsgKey(id), status, time.Hour*24).Err())
}

func (c *msgCache) GetSendMsgStatus(ctx context.Context, id string) (int32, error) {
	result, err := c.rdb.Get(ctx, c.getSendMsgKey(id)).Int()
	return int32(result), errs.Wrap(err)
}

func (c *msgCache) GetMessagesBySeq(ctx context.Context, conversationID string, seqs []int64) (seqMsgs []*sdkws.MsgData, failedSeqs []int64, err error) {
	var keys []string
	keySeqMap := make(map[string]int64, 10)
	for _, seq := range seqs {
		key := c.getMessageCacheKey(conversationID, seq)
		keys = append(keys, key)
		keySeqMap[key] = seq
	}
	err = ProcessKeysBySlot(ctx, c.rdb, keys, func(ctx context.Context, slot int64, keys []string) error {
		result, err := LuaGetBatch(ctx, c.rdb, keys)
		if err != nil {
			return err
		}
		for i, value := range result {
			seq := keySeqMap[keys[i]]
			if value == nil {
				failedSeqs = append(failedSeqs, seq)
				continue
			}

			msg := &sdkws.MsgData{}
			msgString, ok := value.(string)
			if !ok || msgprocessor.String2Pb(msgString, msg) != nil {
				failedSeqs = append(failedSeqs, seq)
				continue
			}
			seqMsgs = append(seqMsgs, msg)

		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return seqMsgs, failedSeqs, nil
}

func (c *msgCache) GetMessageBySeqs(ctx context.Context, conversationID string, seqs []int64) ([]*model.MsgInfoModel, error) {
	if len(seqs) == 0 {
		return nil, nil
	}
	getKey := func(seq int64) string {
		return cachekey.GetMessageCacheKeyV2(conversationID, seq)
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
		return cachekey.GetMessageCacheKeyV2(conversationID, seq)
	})
	slotKeys, err := groupKeysBySlot(ctx, getRocksCacheRedisClient(c.rcClient), keys)
	if err != nil {
		return err
	}
	for _, keys := range slotKeys {
		if err := c.rcClient.TagAsDeletedBatch2(ctx, keys); err != nil {
			return err
		}
	}
	return nil
}
