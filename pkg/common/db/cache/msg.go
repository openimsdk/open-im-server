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

	"github.com/dtm-labs/rockscache"

	unrelationtb "github.com/openimsdk/open-im-server/v3/pkg/common/db/table/unrelation"

	"github.com/openimsdk/open-im-server/v3/pkg/msgprocessor"

	"github.com/OpenIMSDK/tools/errs"

	"github.com/gogo/protobuf/jsonpb"

	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/utils"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"

	"github.com/redis/go-redis/v9"
)

const (
	maxSeq                 = "MAX_SEQ:"
	minSeq                 = "MIN_SEQ:"
	conversationUserMinSeq = "CON_USER_MIN_SEQ:"
	hasReadSeq             = "HAS_READ_SEQ:"

	appleDeviceToken = "DEVICE_TOKEN"
	getuiToken       = "GETUI_TOKEN"
	getuiTaskID      = "GETUI_TASK_ID"
	signalCache      = "SIGNAL_CACHE:"
	signalListCache  = "SIGNAL_LIST_CACHE:"
	FCM_TOKEN        = "FCM_TOKEN:"

	messageCache            = "MESSAGE_CACHE:"
	messageDelUserList      = "MESSAGE_DEL_USER_LIST:"
	userDelMessagesList     = "USER_DEL_MESSAGES_LIST:"
	sendMsgFailedFlag       = "SEND_MSG_FAILED_FLAG:"
	userBadgeUnreadCountSum = "USER_BADGE_UNREAD_COUNT_SUM:"
	exTypeKeyLocker         = "EX_LOCK:"
	uidPidToken             = "UID_PID_TOKEN_STATUS:"
)

type SeqCache interface {
	SetMaxSeq(ctx context.Context, conversationID string, maxSeq int64) error
	GetMaxSeqs(ctx context.Context, conversationIDs []string) (map[string]int64, error)
	GetMaxSeq(ctx context.Context, conversationID string) (int64, error)
	SetMinSeq(ctx context.Context, conversationID string, minSeq int64) error
	SetMinSeqs(ctx context.Context, seqs map[string]int64) error
	GetMinSeqs(ctx context.Context, conversationIDs []string) (map[string]int64, error)
	GetMinSeq(ctx context.Context, conversationID string) (int64, error)
	GetConversationUserMinSeq(ctx context.Context, conversationID string, userID string) (int64, error)
	GetConversationUserMinSeqs(ctx context.Context, conversationID string, userIDs []string) (map[string]int64, error)
	SetConversationUserMinSeq(ctx context.Context, conversationID string, userID string, minSeq int64) error
	// seqs map: key userID value minSeq
	SetConversationUserMinSeqs(ctx context.Context, conversationID string, seqs map[string]int64) (err error)
	// seqs map: key conversationID value minSeq
	SetUserConversationsMinSeqs(ctx context.Context, userID string, seqs map[string]int64) error
	// has read seq
	SetHasReadSeq(ctx context.Context, userID string, conversationID string, hasReadSeq int64) error
	// k: user, v: seq
	SetHasReadSeqs(ctx context.Context, conversationID string, hasReadSeqs map[string]int64) error
	// k: conversation, v :seq
	UserSetHasReadSeqs(ctx context.Context, userID string, hasReadSeqs map[string]int64) error
	GetHasReadSeqs(ctx context.Context, userID string, conversationIDs []string) (map[string]int64, error)
	GetHasReadSeq(ctx context.Context, userID string, conversationID string) (int64, error)
}

type thirdCache interface {
	SetFcmToken(ctx context.Context, account string, platformID int, fcmToken string, expireTime int64) (err error)
	GetFcmToken(ctx context.Context, account string, platformID int) (string, error)
	DelFcmToken(ctx context.Context, account string, platformID int) error
	IncrUserBadgeUnreadCountSum(ctx context.Context, userID string) (int, error)
	SetUserBadgeUnreadCountSum(ctx context.Context, userID string, value int) error
	GetUserBadgeUnreadCountSum(ctx context.Context, userID string) (int, error)
	SetGetuiToken(ctx context.Context, token string, expireTime int64) error
	GetGetuiToken(ctx context.Context) (string, error)
	SetGetuiTaskID(ctx context.Context, taskID string, expireTime int64) error
	GetGetuiTaskID(ctx context.Context) (string, error)
}

type MsgModel interface {
	SeqCache
	thirdCache
	AddTokenFlag(ctx context.Context, userID string, platformID int, token string, flag int) error
	GetTokensWithoutError(ctx context.Context, userID string, platformID int) (map[string]int, error)
	SetTokenMapByUidPid(ctx context.Context, userID string, platformID int, m map[string]int) error
	DeleteTokenByUidPid(ctx context.Context, userID string, platformID int, fields []string) error
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

func NewMsgCacheModel(client redis.UniversalClient) MsgModel {
	return &msgCache{rdb: client}
}

type msgCache struct {
	metaCache
	rdb            redis.UniversalClient
	expireTime     time.Duration
	rcClient       *rockscache.Client
	msgDocDatabase unrelationtb.MsgDocModelInterface
}

func (c *msgCache) getMaxSeqKey(conversationID string) string {
	return maxSeq + conversationID
}

func (c *msgCache) getMinSeqKey(conversationID string) string {
	return minSeq + conversationID
}

func (c *msgCache) getHasReadSeqKey(conversationID string, userID string) string {
	return hasReadSeq + userID + ":" + conversationID
}

func (c *msgCache) setSeq(ctx context.Context, conversationID string, seq int64, getkey func(conversationID string) string) error {
	return utils.Wrap1(c.rdb.Set(ctx, getkey(conversationID), seq, 0).Err())
}

func (c *msgCache) getSeq(ctx context.Context, conversationID string, getkey func(conversationID string) string) (int64, error) {
	return utils.Wrap2(c.rdb.Get(ctx, getkey(conversationID)).Int64())
}

func (c *msgCache) getSeqs(ctx context.Context, items []string, getkey func(s string) string) (m map[string]int64, err error) {
	m = make(map[string]int64, len(items))
	for i, v := range items {
		res, err := c.rdb.Get(ctx, getkey(v)).Result()
		if err != nil && err != redis.Nil {
			return nil, errs.Wrap(err)
		}
		val := utils.StringToInt64(res)
		if val != 0 {
			m[items[i]] = val
		}
	}

	return m, nil

	//pipe := c.rdb.Pipeline()
	//for _, v := range items {
	//	if err := pipe.Get(ctx, getkey(v)).Err(); err != nil && err != redis.Nil {
	//		return nil, errs.Wrap(err)
	//	}
	//}
	//result, err := pipe.Exec(ctx)
	//if err != nil && err != redis.Nil {
	//	return nil, errs.Wrap(err)
	//}
	//m = make(map[string]int64, len(items))
	//for i, v := range result {
	//	seq := v.(*redis.StringCmd)
	//	if seq.Err() != nil && seq.Err() != redis.Nil {
	//		return nil, errs.Wrap(v.Err())
	//	}
	//	val := utils.StringToInt64(seq.Val())
	//	if val != 0 {
	//		m[items[i]] = val
	//	}
	//}
	//return m, nil
}

func (c *msgCache) SetMaxSeq(ctx context.Context, conversationID string, maxSeq int64) error {
	return c.setSeq(ctx, conversationID, maxSeq, c.getMaxSeqKey)
}

func (c *msgCache) GetMaxSeqs(ctx context.Context, conversationIDs []string) (m map[string]int64, err error) {
	return c.getSeqs(ctx, conversationIDs, c.getMaxSeqKey)
}

func (c *msgCache) GetMaxSeq(ctx context.Context, conversationID string) (int64, error) {
	return c.getSeq(ctx, conversationID, c.getMaxSeqKey)
}

func (c *msgCache) SetMinSeq(ctx context.Context, conversationID string, minSeq int64) error {
	return c.setSeq(ctx, conversationID, minSeq, c.getMinSeqKey)
}

func (c *msgCache) setSeqs(ctx context.Context, seqs map[string]int64, getkey func(key string) string) error {
	for conversationID, seq := range seqs {
		if err := c.rdb.Set(ctx, getkey(conversationID), seq, 0).Err(); err != nil {
			return errs.Wrap(err)
		}
	}
	return nil
	//pipe := c.rdb.Pipeline()
	//for k, seq := range seqs {
	//	err := pipe.Set(ctx, getkey(k), seq, 0).Err()
	//	if err != nil {
	//		return errs.Wrap(err)
	//	}
	//}
	//_, err := pipe.Exec(ctx)
	//return err
}

func (c *msgCache) SetMinSeqs(ctx context.Context, seqs map[string]int64) error {
	return c.setSeqs(ctx, seqs, c.getMinSeqKey)
}

func (c *msgCache) GetMinSeqs(ctx context.Context, conversationIDs []string) (map[string]int64, error) {
	return c.getSeqs(ctx, conversationIDs, c.getMinSeqKey)
}

func (c *msgCache) GetMinSeq(ctx context.Context, conversationID string) (int64, error) {
	return c.getSeq(ctx, conversationID, c.getMinSeqKey)
}

func (c *msgCache) getConversationUserMinSeqKey(conversationID, userID string) string {
	return conversationUserMinSeq + conversationID + "u:" + userID
}

func (c *msgCache) GetConversationUserMinSeq(ctx context.Context, conversationID string, userID string) (int64, error) {
	return utils.Wrap2(c.rdb.Get(ctx, c.getConversationUserMinSeqKey(conversationID, userID)).Int64())
}

func (c *msgCache) GetConversationUserMinSeqs(ctx context.Context, conversationID string, userIDs []string) (m map[string]int64, err error) {
	return c.getSeqs(ctx, userIDs, func(userID string) string {
		return c.getConversationUserMinSeqKey(conversationID, userID)
	})
}

func (c *msgCache) SetConversationUserMinSeq(ctx context.Context, conversationID string, userID string, minSeq int64) error {
	return utils.Wrap1(c.rdb.Set(ctx, c.getConversationUserMinSeqKey(conversationID, userID), minSeq, 0).Err())
}

func (c *msgCache) SetConversationUserMinSeqs(ctx context.Context, conversationID string, seqs map[string]int64) (err error) {
	return c.setSeqs(ctx, seqs, func(userID string) string {
		return c.getConversationUserMinSeqKey(conversationID, userID)
	})
}

func (c *msgCache) SetUserConversationsMinSeqs(ctx context.Context, userID string, seqs map[string]int64) (err error) {
	return c.setSeqs(ctx, seqs, func(conversationID string) string {
		return c.getConversationUserMinSeqKey(conversationID, userID)
	})
}

func (c *msgCache) SetHasReadSeq(ctx context.Context, userID string, conversationID string, hasReadSeq int64) error {
	return utils.Wrap1(c.rdb.Set(ctx, c.getHasReadSeqKey(conversationID, userID), hasReadSeq, 0).Err())
}

func (c *msgCache) SetHasReadSeqs(ctx context.Context, conversationID string, hasReadSeqs map[string]int64) error {
	return c.setSeqs(ctx, hasReadSeqs, func(userID string) string {
		return c.getHasReadSeqKey(conversationID, userID)
	})
}

func (c *msgCache) UserSetHasReadSeqs(ctx context.Context, userID string, hasReadSeqs map[string]int64) error {
	return c.setSeqs(ctx, hasReadSeqs, func(conversationID string) string {
		return c.getHasReadSeqKey(conversationID, userID)
	})
}

func (c *msgCache) GetHasReadSeqs(ctx context.Context, userID string, conversationIDs []string) (map[string]int64, error) {
	return c.getSeqs(ctx, conversationIDs, func(conversationID string) string {
		return c.getHasReadSeqKey(conversationID, userID)
	})
}

func (c *msgCache) GetHasReadSeq(ctx context.Context, userID string, conversationID string) (int64, error) {
	return utils.Wrap2(c.rdb.Get(ctx, c.getHasReadSeqKey(conversationID, userID)).Int64())
}

func (c *msgCache) AddTokenFlag(ctx context.Context, userID string, platformID int, token string, flag int) error {
	key := uidPidToken + userID + ":" + constant.PlatformIDToName(platformID)

	return errs.Wrap(c.rdb.HSet(ctx, key, token, flag).Err())
}

func (c *msgCache) GetTokensWithoutError(ctx context.Context, userID string, platformID int) (map[string]int, error) {
	key := uidPidToken + userID + ":" + constant.PlatformIDToName(platformID)
	m, err := c.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, errs.Wrap(err)
	}
	mm := make(map[string]int)
	for k, v := range m {
		mm[k] = utils.StringToInt(v)
	}

	return mm, nil
}

func (c *msgCache) SetTokenMapByUidPid(ctx context.Context, userID string, platform int, m map[string]int) error {
	key := uidPidToken + userID + ":" + constant.PlatformIDToName(platform)
	mm := make(map[string]interface{})
	for k, v := range m {
		mm[k] = v
	}

	return errs.Wrap(c.rdb.HSet(ctx, key, mm).Err())
}

func (c *msgCache) DeleteTokenByUidPid(ctx context.Context, userID string, platform int, fields []string) error {
	key := uidPidToken + userID + ":" + constant.PlatformIDToName(platform)

	return errs.Wrap(c.rdb.HDel(ctx, key, fields...).Err())
}

func (c *msgCache) getMessageCacheKey(conversationID string, seq int64) string {
	return messageCache + conversationID + "_" + strconv.Itoa(int(seq))
}

func (c *msgCache) allMessageCacheKey(conversationID string) string {
	return messageCache + conversationID + "_*"
}

func (c *msgCache) GetMessagesBySeq(ctx context.Context, conversationID string, seqs []int64) (seqMsgs []*sdkws.MsgData, failedSeqs []int64, err error) {
	for _, seq := range seqs {
		res, err := c.rdb.Get(ctx, c.getMessageCacheKey(conversationID, seq)).Result()
		if err != nil {
			log.ZError(ctx, "GetMessagesBySeq failed", err, "conversationID", conversationID, "seq", seq)
			failedSeqs = append(failedSeqs, seq)
			continue
		}
		msg := sdkws.MsgData{}
		if err = msgprocessor.String2Pb(res, &msg); err != nil {
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
	//pipe := c.rdb.Pipeline()
	//for _, v := range seqs {
	//	// MESSAGE_CACHE:169.254.225.224_reliability1653387820_0_1
	//	key := c.getMessageCacheKey(conversationID, v)
	//	if err := pipe.Get(ctx, key).Err(); err != nil && err != redis.Nil {
	//		return nil, nil, err
	//	}
	//}
	//result, err := pipe.Exec(ctx)
	//for i, v := range result {
	//	cmd := v.(*redis.StringCmd)
	//	if cmd.Err() != nil {
	//		failedSeqs = append(failedSeqs, seqs[i])
	//	} else {
	//		msg := sdkws.MsgData{}
	//		err = msgprocessor.String2Pb(cmd.Val(), &msg)
	//		if err == nil {
	//			if msg.Status != constant.MsgDeleted {
	//				seqMsgs = append(seqMsgs, &msg)
	//				continue
	//			}
	//		} else {
	//			log.ZWarn(ctx, "UnmarshalString failed", err, "conversationID", conversationID, "seq", seqs[i], "msg", cmd.Val())
	//		}
	//		failedSeqs = append(failedSeqs, seqs[i])
	//	}
	//}
	//return seqMsgs, failedSeqs, err
}

func (c *msgCache) SetMessageToCache(ctx context.Context, conversationID string, msgs []*sdkws.MsgData) (int, error) {
	for _, msg := range msgs {
		s, err := msgprocessor.Pb2String(msg)
		if err != nil {
			return 0, errs.Wrap(err)
		}
		key := c.getMessageCacheKey(conversationID, msg.Seq)
		if err := c.rdb.Set(ctx, key, s, time.Duration(config.Config.MsgCacheTimeout)*time.Second).Err(); err != nil {
			return 0, errs.Wrap(err)
		}
	}
	return len(msgs), nil
	//pipe := c.rdb.Pipeline()
	//var failedMsgs []*sdkws.MsgData
	//for _, msg := range msgs {
	//	key := c.getMessageCacheKey(conversationID, msg.Seq)
	//	s, err := msgprocessor.Pb2String(msg)
	//	if err != nil {
	//		return 0, errs.Wrap(err)
	//	}
	//	err = pipe.Set(ctx, key, s, time.Duration(config.Config.MsgCacheTimeout)*time.Second).Err()
	//	if err != nil {
	//		failedMsgs = append(failedMsgs, msg)
	//		log.ZWarn(ctx, "set msg 2 cache failed", err, "msg", failedMsgs)
	//	}
	//}
	//_, err := pipe.Exec(ctx)
	//return len(failedMsgs), err
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
		if err := c.rdb.Expire(ctx, delUserListKey, time.Duration(config.Config.MsgCacheTimeout)*time.Second).Err(); err != nil {
			return errs.Wrap(err)
		}
		if err := c.rdb.Expire(ctx, userDelListKey, time.Duration(config.Config.MsgCacheTimeout)*time.Second).Err(); err != nil {
			return errs.Wrap(err)
		}
	}

	return nil
	//pipe := c.rdb.Pipeline()
	//for _, seq := range seqs {
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
	//_, err := pipe.Exec(ctx)
	//return errs.Wrap(err)
}

func (c *msgCache) GetUserDelList(ctx context.Context, userID, conversationID string) (seqs []int64, err error) {
	result, err := c.rdb.SMembers(ctx, c.getUserDelList(conversationID, userID)).Result()
	if err != nil {
		return nil, errs.Wrap(err)
	}
	seqs = make([]int64, len(result))
	for i, v := range result {
		seqs[i] = utils.StringToInt64(v)
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
	//for _, seq := range seqs {
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
	for _, seq := range seqs {
		if err := c.rdb.Del(ctx, c.getMessageCacheKey(conversationID, seq)).Err(); err != nil {
			return errs.Wrap(err)
		}
	}
	return nil
	//pipe := c.rdb.Pipeline()
	//for _, seq := range seqs {
	//	if err := pipe.Del(ctx, c.getMessageCacheKey(conversationID, seq)).Err(); err != nil {
	//		return errs.Wrap(err)
	//	}
	//}
	//_, err := pipe.Exec(ctx)
	//return errs.Wrap(err)
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
	//pipe := c.rdb.Pipeline()
	//for _, v := range vals {
	//	if err := pipe.Del(ctx, v).Err(); err != nil {
	//		return errs.Wrap(err)
	//	}
	//}
	//_, err = pipe.Exec(ctx)
	//return errs.Wrap(err)
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
		if err := c.rdb.Set(ctx, key, s, time.Duration(config.Config.MsgCacheTimeout)*time.Second).Err(); err != nil {
			return errs.Wrap(err)
		}
	}

	return nil
}

func (c *msgCache) SetGetuiToken(ctx context.Context, token string, expireTime int64) error {
	return errs.Wrap(c.rdb.Set(ctx, getuiToken, token, time.Duration(expireTime)*time.Second).Err())
}

func (c *msgCache) GetGetuiToken(ctx context.Context) (string, error) {
	return utils.Wrap2(c.rdb.Get(ctx, getuiToken).Result())
}

func (c *msgCache) SetGetuiTaskID(ctx context.Context, taskID string, expireTime int64) error {
	return errs.Wrap(c.rdb.Set(ctx, getuiTaskID, taskID, time.Duration(expireTime)*time.Second).Err())
}

func (c *msgCache) GetGetuiTaskID(ctx context.Context) (string, error) {
	return utils.Wrap2(c.rdb.Get(ctx, getuiTaskID).Result())
}

func (c *msgCache) SetSendMsgStatus(ctx context.Context, id string, status int32) error {
	return errs.Wrap(c.rdb.Set(ctx, sendMsgFailedFlag+id, status, time.Hour*24).Err())
}

func (c *msgCache) GetSendMsgStatus(ctx context.Context, id string) (int32, error) {
	result, err := c.rdb.Get(ctx, sendMsgFailedFlag+id).Int()

	return int32(result), errs.Wrap(err)
}

func (c *msgCache) SetFcmToken(ctx context.Context, account string, platformID int, fcmToken string, expireTime int64) (err error) {
	return errs.Wrap(c.rdb.Set(ctx, FCM_TOKEN+account+":"+strconv.Itoa(platformID), fcmToken, time.Duration(expireTime)*time.Second).Err())
}

func (c *msgCache) GetFcmToken(ctx context.Context, account string, platformID int) (string, error) {
	return utils.Wrap2(c.rdb.Get(ctx, FCM_TOKEN+account+":"+strconv.Itoa(platformID)).Result())
}

func (c *msgCache) DelFcmToken(ctx context.Context, account string, platformID int) error {
	return errs.Wrap(c.rdb.Del(ctx, FCM_TOKEN+account+":"+strconv.Itoa(platformID)).Err())
}

func (c *msgCache) IncrUserBadgeUnreadCountSum(ctx context.Context, userID string) (int, error) {
	seq, err := c.rdb.Incr(ctx, userBadgeUnreadCountSum+userID).Result()

	return int(seq), errs.Wrap(err)
}

func (c *msgCache) SetUserBadgeUnreadCountSum(ctx context.Context, userID string, value int) error {
	return errs.Wrap(c.rdb.Set(ctx, userBadgeUnreadCountSum+userID, value, 0).Err())
}

func (c *msgCache) GetUserBadgeUnreadCountSum(ctx context.Context, userID string) (int, error) {
	return utils.Wrap2(c.rdb.Get(ctx, userBadgeUnreadCountSum+userID).Int())
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
	case constant.GroupChatType:
		return "EX_GROUP_" + clientMsgID
	case constant.SuperGroupChatType:
		return "EX_SUPER_GROUP_" + clientMsgID
	case constant.NotificationChatType:
		return "EX_NOTIFICATION" + clientMsgID
	}

	return ""
}

func (c *msgCache) JudgeMessageReactionExist(ctx context.Context, clientMsgID string, sessionType int32) (bool, error) {
	n, err := c.rdb.Exists(ctx, c.getMessageReactionExPrefix(clientMsgID, sessionType)).Result()
	if err != nil {
		return false, utils.Wrap(err, "")
	}

	return n > 0, nil
}

func (c *msgCache) SetMessageTypeKeyValue(
	ctx context.Context,
	clientMsgID string,
	sessionType int32,
	typeKey, value string,
) error {
	return errs.Wrap(c.rdb.HSet(ctx, c.getMessageReactionExPrefix(clientMsgID, sessionType), typeKey, value).Err())
}

func (c *msgCache) SetMessageReactionExpire(ctx context.Context, clientMsgID string, sessionType int32, expiration time.Duration) (bool, error) {
	return utils.Wrap2(c.rdb.Expire(ctx, c.getMessageReactionExPrefix(clientMsgID, sessionType), expiration).Result())
}

func (c *msgCache) GetMessageTypeKeyValue(ctx context.Context, clientMsgID string, sessionType int32, typeKey string) (string, error) {
	return utils.Wrap2(c.rdb.HGet(ctx, c.getMessageReactionExPrefix(clientMsgID, sessionType), typeKey).Result())
}

func (c *msgCache) GetOneMessageAllReactionList(
	ctx context.Context,
	clientMsgID string,
	sessionType int32,
) (map[string]string, error) {
	return utils.Wrap2(c.rdb.HGetAll(ctx, c.getMessageReactionExPrefix(clientMsgID, sessionType)).Result())
}

func (c *msgCache) DeleteOneMessageKey(
	ctx context.Context,
	clientMsgID string,
	sessionType int32,
	subKey string,
) error {
	return errs.Wrap(c.rdb.HDel(ctx, c.getMessageReactionExPrefix(clientMsgID, sessionType), subKey).Err())
}
