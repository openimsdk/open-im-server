package redis

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/msgprocessor"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/redis/go-redis/v9"
	"time"
) //

// msgCacheTimeout is  expiration time of message cache, 86400 seconds
const msgCacheTimeout = 86400

func NewMsgCache(client redis.UniversalClient) cache.MsgCache {
	return &msgCache{rdb: client}
}

type msgCache struct {
	rdb redis.UniversalClient
}

func (c *msgCache) getMessageCacheKey(conversationID string, seq int64) string {
	return cachekey.GetMessageCacheKey(conversationID, seq)
}
func (c *msgCache) getMessageDelUserListKey(conversationID string, seq int64) string {
	return cachekey.GetMessageDelUserListKey(conversationID, seq)
}

func (c *msgCache) getUserDelList(conversationID, userID string) string {
	return cachekey.GetUserDelListKey(conversationID, userID)
}

func (c *msgCache) getSendMsgKey(id string) string {
	return cachekey.GetSendMsgKey(id)
}

func (c *msgCache) getLockMessageTypeKey(clientMsgID string, TypeKey string) string {
	return cachekey.GetLockMessageTypeKey(clientMsgID, TypeKey)
}

func (c *msgCache) getMessageReactionExPrefix(clientMsgID string, sessionType int32) string {
	return cachekey.GetMessageReactionExKey(clientMsgID, sessionType)
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
		return LuaSetBatchWithCommonExpire(ctx, c.rdb, keys, values, msgCacheTimeout)
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

func (c *msgCache) LockMessageTypeKey(ctx context.Context, clientMsgID string, TypeKey string) error {
	key := c.getLockMessageTypeKey(clientMsgID, TypeKey)
	return errs.Wrap(c.rdb.SetNX(ctx, key, 1, time.Minute).Err())
}

func (c *msgCache) UnLockMessageTypeKey(ctx context.Context, clientMsgID string, TypeKey string) error {
	key := c.getLockMessageTypeKey(clientMsgID, TypeKey)
	return errs.Wrap(c.rdb.Del(ctx, key).Err())
}

func (c *msgCache) JudgeMessageReactionExist(ctx context.Context, clientMsgID string, sessionType int32) (bool, error) {
	n, err := c.rdb.Exists(ctx, c.getMessageReactionExPrefix(clientMsgID, sessionType)).Result()
	if err != nil {
		return false, errs.Wrap(err)
	}

	return n > 0, nil
}

func (c *msgCache) SetMessageTypeKeyValue(ctx context.Context, clientMsgID string, sessionType int32, typeKey, value string) error {
	return errs.Wrap(c.rdb.HSet(ctx, c.getMessageReactionExPrefix(clientMsgID, sessionType), typeKey, value).Err())
}

func (c *msgCache) SetMessageReactionExpire(ctx context.Context, clientMsgID string, sessionType int32, expiration time.Duration) (bool, error) {
	val, err := c.rdb.Expire(ctx, c.getMessageReactionExPrefix(clientMsgID, sessionType), expiration).Result()
	return val, errs.Wrap(err)
}

func (c *msgCache) GetMessageTypeKeyValue(ctx context.Context, clientMsgID string, sessionType int32, typeKey string) (string, error) {
	val, err := c.rdb.HGet(ctx, c.getMessageReactionExPrefix(clientMsgID, sessionType), typeKey).Result()
	return val, errs.Wrap(err)
}

func (c *msgCache) GetOneMessageAllReactionList(ctx context.Context, clientMsgID string, sessionType int32) (map[string]string, error) {
	val, err := c.rdb.HGetAll(ctx, c.getMessageReactionExPrefix(clientMsgID, sessionType)).Result()
	return val, errs.Wrap(err)
}

func (c *msgCache) DeleteOneMessageKey(ctx context.Context, clientMsgID string, sessionType int32, subKey string) error {
	return errs.Wrap(c.rdb.HDel(ctx, c.getMessageReactionExPrefix(clientMsgID, sessionType), subKey).Err())
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
