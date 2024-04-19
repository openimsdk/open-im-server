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
	"errors"
	"strconv"
	"time"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/openimsdk/open-im-server/v3/pkg/msgprocessor"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/stringutil"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
)

const msgCacheTimeout = 86400 * time.Second

const (
	maxSeq                 = "MAX_SEQ:"
	minSeq                 = "MIN_SEQ:"
	conversationUserMinSeq = "CON_USER_MIN_SEQ:"
	hasReadSeq             = "HAS_READ_SEQ:"

	getuiToken  = "GETUI_TOKEN"
	getuiTaskID = "GETUI_TASK_ID"
	FCM_TOKEN   = "FCM_TOKEN:"

	messageCache            = "MESSAGE_CACHE:"
	messageDelUserList      = "MESSAGE_DEL_USER_LIST:"
	userDelMessagesList     = "USER_DEL_MESSAGES_LIST:"
	sendMsgFailedFlag       = "SEND_MSG_FAILED_FLAG:"
	userBadgeUnreadCountSum = "USER_BADGE_UNREAD_COUNT_SUM:"
	exTypeKeyLocker         = "EX_LOCK:"
)

var concurrentLimit = 3

//type MsgModel interface {
//	SeqCache
//	ThirdCache
//	MsgCache
//}

type MsgCache interface {
	GetMessagesBySeq(ctx context.Context, conversationID string, seqs []int64) (seqMsg []*sdkws.MsgData, failedSeqList []int64, err error)
	SetMessageToCache(ctx context.Context, conversationID string, msgs []*sdkws.MsgData) (int, error)
	UserDeleteMsgs(ctx context.Context, conversationID string, seqs []int64, userID string) error
	DelUserDeleteMsgsList(ctx context.Context, conversationID string, seqs []int64)
	DeleteMessages(ctx context.Context, conversationID string, seqs []int64) error
	GetUserDelList(ctx context.Context, userID, conversationID string) (seqs []int64, err error)
	CleanUpOneConversationAllMsg(ctx context.Context, conversationID string) error
	DelMsgFromCache(ctx context.Context, userID string, seqList []int64) error
	SetSendMsgStatus(ctx context.Context, id string, status int32) error
	GetSendMsgStatus(ctx context.Context, id string) (int32, error)
	JudgeMessageReactionExist(ctx context.Context, clientMsgID string, sessionType int32) (bool, error)
	GetOneMessageAllReactionList(ctx context.Context, clientMsgID string, sessionType int32) (map[string]string, error)
	DeleteOneMessageKey(ctx context.Context, clientMsgID string, sessionType int32, subKey string) error
	SetMessageReactionExpire(ctx context.Context, clientMsgID string, sessionType int32, expiration time.Duration) (bool, error)
	GetMessageTypeKeyValue(ctx context.Context, clientMsgID string, sessionType int32, typeKey string) (string, error)
	SetMessageTypeKeyValue(ctx context.Context, clientMsgID string, sessionType int32, typeKey, value string) error
	LockMessageTypeKey(ctx context.Context, clientMsgID string, TypeKey string) error
	UnLockMessageTypeKey(ctx context.Context, clientMsgID string, TypeKey string) error
}

//func NewMsgCacheModel(client redis.UniversalClient, msgCacheTimeout int, redisConf *config.Redis) MsgModel {
//	return &msgCache{rdb: client, msgCacheTimeout: msgCacheTimeout, redisConf: redisConf}
//}

func NewMsgCache(client redis.UniversalClient, redisEnablePipeline bool) MsgCache {
	return &msgCache{rdb: client, msgCacheTimeout: msgCacheTimeout, redisEnablePipeline: redisEnablePipeline}
}

type msgCache struct {
	metaCache
	rdb                 redis.UniversalClient
	msgCacheTimeout     time.Duration
	redisEnablePipeline bool
}

func (c *msgCache) allMessageCacheKey(conversationID string) string {
	return messageCache + conversationID + "_*"
}

func (c *msgCache) getMessageCacheKey(conversationID string, seq int64) string {
	return messageCache + conversationID + "_" + strconv.Itoa(int(seq))
}

func (c *msgCache) SetMessageToCache(ctx context.Context, conversationID string, msgs []*sdkws.MsgData) (int, error) {
	if c.redisEnablePipeline {
		return c.PipeSetMessageToCache(ctx, conversationID, msgs)
	}
	return c.ParallelSetMessageToCache(ctx, conversationID, msgs)
}

func (c *msgCache) PipeSetMessageToCache(ctx context.Context, conversationID string, msgs []*sdkws.MsgData) (int, error) {
	pipe := c.rdb.Pipeline()
	for _, msg := range msgs {
		s, err := msgprocessor.Pb2String(msg)
		if err != nil {
			return 0, err
		}

		key := c.getMessageCacheKey(conversationID, msg.Seq)
		_ = pipe.Set(ctx, key, s, c.msgCacheTimeout)
	}

	results, err := pipe.Exec(ctx)
	if err != nil {
		return 0, errs.Wrap(err)
	}

	for _, res := range results {
		if res.Err() != nil {
			return 0, errs.Wrap(err)
		}
	}

	return len(msgs), nil
}

func (c *msgCache) ParallelSetMessageToCache(ctx context.Context, conversationID string, msgs []*sdkws.MsgData) (int, error) {
	wg := errgroup.Group{}
	wg.SetLimit(concurrentLimit)

	for _, msg := range msgs {
		msg := msg // closure safe var
		wg.Go(func() error {
			s, err := msgprocessor.Pb2String(msg)
			if err != nil {
				return errs.Wrap(err)
			}

			key := c.getMessageCacheKey(conversationID, msg.Seq)
			if err := c.rdb.Set(ctx, key, s, c.msgCacheTimeout).Err(); err != nil {
				return errs.Wrap(err)
			}
			return nil
		})
	}

	err := wg.Wait()
	if err != nil {
		return 0, errs.WrapMsg(err, "wg.Wait failed")
	}

	return len(msgs), nil
}

func (c *msgCache) getMessageDelUserListKey(conversationID string, seq int64) string {
	return messageDelUserList + conversationID + ":" + strconv.Itoa(int(seq))
}

func (c *msgCache) getUserDelList(conversationID, userID string) string {
	return userDelMessagesList + conversationID + ":" + userID
}

func (c *msgCache) UserDeleteMsgs(ctx context.Context, conversationID string, seqs []int64, userID string) error {
	for _, seq := range seqs {
		delUserListKey := c.getMessageDelUserListKey(conversationID, seq)
		userDelListKey := c.getUserDelList(conversationID, userID)
		err := c.rdb.SAdd(ctx, delUserListKey, userID).Err()
		if err != nil {
			return errs.Wrap(err)
		}
		err = c.rdb.SAdd(ctx, userDelListKey, seq).Err()
		if err != nil {
			return errs.Wrap(err)
		}
		if err := c.rdb.Expire(ctx, delUserListKey, c.msgCacheTimeout).Err(); err != nil {
			return errs.Wrap(err)
		}
		if err := c.rdb.Expire(ctx, userDelListKey, c.msgCacheTimeout).Err(); err != nil {
			return errs.Wrap(err)
		}
	}

	return nil
	// pipe := c.rdb.Pipeline()
	// for _, seq := range seqs {
	//	delUserListKey := c.getMessageDelUserListKey(conversationID, seq)
	//	userDelListKey := c.getUserDelList(conversationID, userID)
	//	err := pipe.SAdd(ctx, delUserListKey, userID).Err()
	//	if err != nil {
	//		return errs.Wrap(err)
	//	}
	//	err = pipe.SAdd(ctx, userDelListKey, seq).Err()
	//	if err != nil {
	//		return errs.Wrap(err)
	//	}
	//	if err := pipe.Expire(ctx, delUserListKey, time.Duration(config.Config.MsgCacheTimeout)*time.Second).Err(); err != nil {
	//		return errs.Wrap(err)
	//	}
	//	if err := pipe.Expire(ctx, userDelListKey, time.Duration(config.Config.MsgCacheTimeout)*time.Second).Err(); err != nil {
	//		return errs.Wrap(err)
	//	}
	//}
	// _, err := pipe.Exec(ctx)
	// return errs.Wrap(err)
}

func (c *msgCache) GetUserDelList(ctx context.Context, userID, conversationID string) (seqs []int64, err error) {
	result, err := c.rdb.SMembers(ctx, c.getUserDelList(conversationID, userID)).Result()
	if err != nil {
		return nil, errs.Wrap(err)
	}
	seqs = make([]int64, len(result))
	for i, v := range result {
		seqs[i] = stringutil.StringToInt64(v)
	}

	return seqs, nil
}

func (c *msgCache) DelUserDeleteMsgsList(ctx context.Context, conversationID string, seqs []int64) {
	for _, seq := range seqs {
		delUsers, err := c.rdb.SMembers(ctx, c.getMessageDelUserListKey(conversationID, seq)).Result()
		if err != nil {
			log.ZWarn(ctx, "DelUserDeleteMsgsList failed", err, "conversationID", conversationID, "seq", seq)

			continue
		}
		if len(delUsers) > 0 {
			var failedFlag bool
			for _, userID := range delUsers {
				err = c.rdb.SRem(ctx, c.getUserDelList(conversationID, userID), seq).Err()
				if err != nil {
					failedFlag = true
					log.ZWarn(ctx, "DelUserDeleteMsgsList failed", err, "conversationID", conversationID, "seq", seq, "userID", userID)
				}
			}
			if !failedFlag {
				if err := c.rdb.Del(ctx, c.getMessageDelUserListKey(conversationID, seq)).Err(); err != nil {
					log.ZWarn(ctx, "DelUserDeleteMsgsList failed", err, "conversationID", conversationID, "seq", seq)
				}
			}
		}
	}
	// for _, seq := range seqs {
	//	delUsers, err := c.rdb.SMembers(ctx, c.getMessageDelUserListKey(conversationID, seq)).Result()
	//	if err != nil {
	//		log.ZWarn(ctx, "DelUserDeleteMsgsList failed", err, "conversationID", conversationID, "seq", seq)
	//		continue
	//	}
	//	if len(delUsers) > 0 {
	//		pipe := c.rdb.Pipeline()
	//		var failedFlag bool
	//		for _, userID := range delUsers {
	//			err = pipe.SRem(ctx, c.getUserDelList(conversationID, userID), seq).Err()
	//			if err != nil {
	//				failedFlag = true
	//				log.ZWarn(
	//					ctx,
	//					"DelUserDeleteMsgsList failed",
	//					err,
	//					"conversationID",
	//					conversationID,
	//					"seq",
	//					seq,
	//					"userID",
	//					userID,
	//				)
	//			}
	//		}
	//		if !failedFlag {
	//			if err := pipe.Del(ctx, c.getMessageDelUserListKey(conversationID, seq)).Err(); err != nil {
	//				log.ZWarn(ctx, "DelUserDeleteMsgsList failed", err, "conversationID", conversationID, "seq", seq)
	//			}
	//		}
	//		if _, err := pipe.Exec(ctx); err != nil {
	//			log.ZError(ctx, "pipe exec failed", err, "conversationID", conversationID, "seq", seq)
	//		}
	//	}
	//}
}

func (c *msgCache) DeleteMessages(ctx context.Context, conversationID string, seqs []int64) error {
	if c.redisEnablePipeline {
		return c.PipeDeleteMessages(ctx, conversationID, seqs)
	}

	return c.ParallelDeleteMessages(ctx, conversationID, seqs)
}

func (c *msgCache) ParallelDeleteMessages(ctx context.Context, conversationID string, seqs []int64) error {
	wg := errgroup.Group{}
	wg.SetLimit(concurrentLimit)

	for _, seq := range seqs {
		seq := seq
		wg.Go(func() error {
			err := c.rdb.Del(ctx, c.getMessageCacheKey(conversationID, seq)).Err()
			if err != nil {
				return errs.Wrap(err)
			}
			return nil
		})
	}

	return wg.Wait()
}

func (c *msgCache) PipeDeleteMessages(ctx context.Context, conversationID string, seqs []int64) error {
	pipe := c.rdb.Pipeline()
	for _, seq := range seqs {
		_ = pipe.Del(ctx, c.getMessageCacheKey(conversationID, seq))
	}

	results, err := pipe.Exec(ctx)
	if err != nil {
		return errs.WrapMsg(err, "pipe.del")
	}

	for _, res := range results {
		if res.Err() != nil {
			return errs.Wrap(err)
		}
	}

	return nil
}

func (c *msgCache) CleanUpOneConversationAllMsg(ctx context.Context, conversationID string) error {
	vals, err := c.rdb.Keys(ctx, c.allMessageCacheKey(conversationID)).Result()
	if errors.Is(err, redis.Nil) {
		return nil
	}
	if err != nil {
		return errs.Wrap(err)
	}
	for _, v := range vals {
		if err := c.rdb.Del(ctx, v).Err(); err != nil {
			return errs.Wrap(err)
		}
	}
	return nil
}

func (c *msgCache) DelMsgFromCache(ctx context.Context, userID string, seqs []int64) error {
	for _, seq := range seqs {
		key := c.getMessageCacheKey(userID, seq)
		result, err := c.rdb.Get(ctx, key).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				continue
			}

			return errs.Wrap(err)
		}
		var msg sdkws.MsgData
		err = jsonpb.UnmarshalString(result, &msg)
		if err != nil {
			return err
		}
		msg.Status = constant.MsgDeleted
		s, err := msgprocessor.Pb2String(&msg)
		if err != nil {
			return errs.Wrap(err)
		}
		if err := c.rdb.Set(ctx, key, s, c.msgCacheTimeout).Err(); err != nil {
			return errs.Wrap(err)
		}
	}

	return nil
}

func (c *msgCache) SetSendMsgStatus(ctx context.Context, id string, status int32) error {
	return errs.Wrap(c.rdb.Set(ctx, sendMsgFailedFlag+id, status, time.Hour*24).Err())
}

func (c *msgCache) GetSendMsgStatus(ctx context.Context, id string) (int32, error) {
	result, err := c.rdb.Get(ctx, sendMsgFailedFlag+id).Int()

	return int32(result), errs.Wrap(err)
}

func (c *msgCache) LockMessageTypeKey(ctx context.Context, clientMsgID string, TypeKey string) error {
	key := exTypeKeyLocker + clientMsgID + "_" + TypeKey

	return errs.Wrap(c.rdb.SetNX(ctx, key, 1, time.Minute).Err())
}

func (c *msgCache) UnLockMessageTypeKey(ctx context.Context, clientMsgID string, TypeKey string) error {
	key := exTypeKeyLocker + clientMsgID + "_" + TypeKey

	return errs.Wrap(c.rdb.Del(ctx, key).Err())
}

func (c *msgCache) getMessageReactionExPrefix(clientMsgID string, sessionType int32) string {
	switch sessionType {
	case constant.SingleChatType:
		return "EX_SINGLE_" + clientMsgID
	case constant.WriteGroupChatType:
		return "EX_GROUP_" + clientMsgID
	case constant.ReadGroupChatType:
		return "EX_SUPER_GROUP_" + clientMsgID
	case constant.NotificationChatType:
		return "EX_NOTIFICATION" + clientMsgID
	}

	return ""
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
	if c.redisEnablePipeline {
		return c.PipeGetMessagesBySeq(ctx, conversationID, seqs)
	}

	return c.ParallelGetMessagesBySeq(ctx, conversationID, seqs)
}

func (c *msgCache) PipeGetMessagesBySeq(ctx context.Context, conversationID string, seqs []int64) (seqMsgs []*sdkws.MsgData, failedSeqs []int64, err error) {
	pipe := c.rdb.Pipeline()

	results := []*redis.StringCmd{}
	for _, seq := range seqs {
		results = append(results, pipe.Get(ctx, c.getMessageCacheKey(conversationID, seq)))
	}

	_, err = pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return seqMsgs, failedSeqs, errs.WrapMsg(err, "pipe.get")
	}

	for idx, res := range results {
		seq := seqs[idx]
		if res.Err() != nil {
			log.ZError(ctx, "GetMessagesBySeq failed", err, "conversationID", conversationID, "seq", seq, "err", res.Err())
			failedSeqs = append(failedSeqs, seq)
			continue
		}

		msg := sdkws.MsgData{}
		if err = msgprocessor.String2Pb(res.Val(), &msg); err != nil {
			log.ZError(ctx, "GetMessagesBySeq Unmarshal failed", err, "res", res, "conversationID", conversationID, "seq", seq)
			failedSeqs = append(failedSeqs, seq)
			continue
		}

		if msg.Status == constant.MsgDeleted {
			failedSeqs = append(failedSeqs, seq)
			continue
		}

		seqMsgs = append(seqMsgs, &msg)
	}

	return
}

func (c *msgCache) ParallelGetMessagesBySeq(ctx context.Context, conversationID string, seqs []int64) (seqMsgs []*sdkws.MsgData, failedSeqs []int64, err error) {
	type entry struct {
		err error
		msg *sdkws.MsgData
	}

	wg := errgroup.Group{}
	wg.SetLimit(concurrentLimit)

	results := make([]entry, len(seqs)) // set slice len/cap to length of seqs.
	for idx, seq := range seqs {
		// closure safe var
		idx := idx
		seq := seq

		wg.Go(func() error {
			res, err := c.rdb.Get(ctx, c.getMessageCacheKey(conversationID, seq)).Result()
			if err != nil {
				log.ZError(ctx, "GetMessagesBySeq failed", err, "conversationID", conversationID, "seq", seq)
				results[idx] = entry{err: err}
				return nil
			}

			msg := sdkws.MsgData{}
			if err = msgprocessor.String2Pb(res, &msg); err != nil {
				log.ZError(ctx, "GetMessagesBySeq Unmarshal failed", err, "res", res, "conversationID", conversationID, "seq", seq)
				results[idx] = entry{err: err}
				return nil
			}

			if msg.Status == constant.MsgDeleted {
				results[idx] = entry{err: err}
				return nil
			}

			results[idx] = entry{msg: &msg}
			return nil
		})
	}

	_ = wg.Wait()

	for idx, res := range results {
		if res.err != nil {
			failedSeqs = append(failedSeqs, seqs[idx])
			continue
		}

		seqMsgs = append(seqMsgs, res.msg)
	}

	return
}
