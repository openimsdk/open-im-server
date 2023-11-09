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

package controller

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/convert"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
	unrelationtb "github.com/openimsdk/open-im-server/v3/pkg/common/db/table/unrelation"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/unrelation"
	"github.com/openimsdk/open-im-server/v3/pkg/common/kafka"
	"github.com/openimsdk/open-im-server/v3/pkg/common/prommetrics"
	"go.mongodb.org/mongo-driver/mongo"

	pbmsg "github.com/OpenIMSDK/protocol/msg"
	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/tools/utils"
)

const (
	updateKeyMsg = iota
	updateKeyRevoke
)

type CommonMsgDatabase interface {
	// 批量插入消息
	BatchInsertChat2DB(ctx context.Context, conversationID string, msgs []*sdkws.MsgData, currentMaxSeq int64) error
	// 撤回消息
	RevokeMsg(ctx context.Context, conversationID string, seq int64, revoke *unrelationtb.RevokeModel) error
	// mark as read
	MarkSingleChatMsgsAsRead(ctx context.Context, userID string, conversationID string, seqs []int64) error
	// 刪除redis中消息缓存
	DeleteMessagesFromCache(ctx context.Context, conversationID string, seqs []int64) error
	DelUserDeleteMsgsList(ctx context.Context, conversationID string, seqs []int64)
	// incrSeq然后批量插入缓存
	BatchInsertChat2Cache(ctx context.Context, conversationID string, msgs []*sdkws.MsgData) (seq int64, isNewConversation bool, err error)

	//  通过seqList获取mongo中写扩散消息
	GetMsgBySeqsRange(ctx context.Context, userID string, conversationID string, begin, end, num, userMaxSeq int64) (minSeq int64, maxSeq int64, seqMsg []*sdkws.MsgData, err error)
	// 通过seqList获取大群在 mongo里面的消息
	GetMsgBySeqs(ctx context.Context, userID string, conversationID string, seqs []int64) (minSeq int64, maxSeq int64, seqMsg []*sdkws.MsgData, err error)
	// 删除会话消息重置最小seq， remainTime为消息保留的时间单位秒,超时消息删除， 传0删除所有消息(此方法不删除redis cache)
	DeleteConversationMsgsAndSetMinSeq(ctx context.Context, conversationID string, remainTime int64) error
	// 用户标记删除过期消息返回标记删除的seq列表
	UserMsgsDestruct(ctx context.Context, userID string, conversationID string, destructTime int64, lastMsgDestructTime time.Time) (seqs []int64, err error)

	// 用户根据seq删除消息
	DeleteUserMsgsBySeqs(ctx context.Context, userID string, conversationID string, seqs []int64) error
	// 物理删除消息置空
	DeleteMsgsPhysicalBySeqs(ctx context.Context, conversationID string, seqs []int64) error

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
	SetConversationUserMinSeqs(ctx context.Context, conversationID string, seqs map[string]int64) (err error)
	SetUserConversationsMinSeqs(ctx context.Context, userID string, seqs map[string]int64) (err error)
	SetHasReadSeq(ctx context.Context, userID string, conversationID string, hasReadSeq int64) error
	GetHasReadSeqs(ctx context.Context, userID string, conversationIDs []string) (map[string]int64, error)
	GetHasReadSeq(ctx context.Context, userID string, conversationID string) (int64, error)
	UserSetHasReadSeqs(ctx context.Context, userID string, hasReadSeqs map[string]int64) error

	GetMongoMaxAndMinSeq(ctx context.Context, conversationID string) (minSeqMongo, maxSeqMongo int64, err error)
	GetConversationMinMaxSeqInMongoAndCache(ctx context.Context, conversationID string) (minSeqMongo, maxSeqMongo, minSeqCache, maxSeqCache int64, err error)
	SetSendMsgStatus(ctx context.Context, id string, status int32) error
	GetSendMsgStatus(ctx context.Context, id string) (int32, error)
	SearchMessage(ctx context.Context, req *pbmsg.SearchMessageReq) (total int32, msgData []*sdkws.MsgData, err error)

	// to mq
	MsgToMQ(ctx context.Context, key string, msg2mq *sdkws.MsgData) error
	MsgToModifyMQ(ctx context.Context, key, conversarionID string, msgs []*sdkws.MsgData) error
	MsgToPushMQ(ctx context.Context, key, conversarionID string, msg2mq *sdkws.MsgData) (int32, int64, error)
	MsgToMongoMQ(ctx context.Context, key, conversarionID string, msgs []*sdkws.MsgData, lastSeq int64) error

	RangeUserSendCount(
		ctx context.Context,
		start time.Time,
		end time.Time,
		group bool,
		ase bool,
		pageNumber int32,
		showNumber int32,
	) (msgCount int64, userCount int64, users []*unrelationtb.UserCount, dateCount map[string]int64, err error)
	RangeGroupSendCount(
		ctx context.Context,
		start time.Time,
		end time.Time,
		ase bool,
		pageNumber int32,
		showNumber int32,
	) (msgCount int64, userCount int64, groups []*unrelationtb.GroupCount, dateCount map[string]int64, err error)
	ConvertMsgsDocLen(ctx context.Context, conversationIDs []string)
}

func NewCommonMsgDatabase(msgDocModel unrelationtb.MsgDocModelInterface, cacheModel cache.MsgModel) CommonMsgDatabase {
	return &commonMsgDatabase{
		msgDocDatabase:  msgDocModel,
		cache:           cacheModel,
		producer:        kafka.NewKafkaProducer(config.Config.Kafka.Addr, config.Config.Kafka.LatestMsgToRedis.Topic),
		producerToMongo: kafka.NewKafkaProducer(config.Config.Kafka.Addr, config.Config.Kafka.MsgToMongo.Topic),
		producerToPush:  kafka.NewKafkaProducer(config.Config.Kafka.Addr, config.Config.Kafka.MsgToPush.Topic),
	}
}

func InitCommonMsgDatabase(rdb redis.UniversalClient, database *mongo.Database) CommonMsgDatabase {
	cacheModel := cache.NewMsgCacheModel(rdb)
	msgDocModel := unrelation.NewMsgMongoDriver(database)
	CommonMsgDatabase := NewCommonMsgDatabase(msgDocModel, cacheModel)
	return CommonMsgDatabase
}

type commonMsgDatabase struct {
	msgDocDatabase   unrelationtb.MsgDocModelInterface
	msg              unrelationtb.MsgDocModel
	cache            cache.MsgModel
	producer         *kafka.Producer
	producerToMongo  *kafka.Producer
	producerToModify *kafka.Producer
	producerToPush   *kafka.Producer
}

func (db *commonMsgDatabase) MsgToMQ(ctx context.Context, key string, msg2mq *sdkws.MsgData) error {
	_, _, err := db.producer.SendMessage(ctx, key, msg2mq)
	return err
}

func (db *commonMsgDatabase) MsgToModifyMQ(ctx context.Context, key, conversationID string, messages []*sdkws.MsgData) error {
	if len(messages) > 0 {
		_, _, err := db.producerToModify.SendMessage(ctx, key, &pbmsg.MsgDataToModifyByMQ{ConversationID: conversationID, Messages: messages})
		return err
	}
	return nil
}

func (db *commonMsgDatabase) MsgToPushMQ(ctx context.Context, key, conversationID string, msg2mq *sdkws.MsgData) (int32, int64, error) {
	partition, offset, err := db.producerToPush.SendMessage(ctx, key, &pbmsg.PushMsgDataToMQ{MsgData: msg2mq, ConversationID: conversationID})
	if err != nil {
		log.ZError(ctx, "MsgToPushMQ", err, "key", key, "msg2mq", msg2mq)
		return 0, 0, err
	}
	return partition, offset, nil
}

func (db *commonMsgDatabase) MsgToMongoMQ(ctx context.Context, key, conversationID string, messages []*sdkws.MsgData, lastSeq int64) error {
	if len(messages) > 0 {
		_, _, err := db.producerToMongo.SendMessage(ctx, key, &pbmsg.MsgDataToMongoByMQ{LastSeq: lastSeq, ConversationID: conversationID, MsgData: messages})
		return err
	}
	return nil
}

func (db *commonMsgDatabase) BatchInsertBlock(ctx context.Context, conversationID string, fields []any, key int8, firstSeq int64) error {
	if len(fields) == 0 {
		return nil
	}
	num := db.msg.GetSingleGocMsgNum()
	// num = 100
	for i, field := range fields { // 检查类型
		var ok bool
		switch key {
		case updateKeyMsg:
			var msg *unrelationtb.MsgDataModel
			msg, ok = field.(*unrelationtb.MsgDataModel)
			if msg != nil && msg.Seq != firstSeq+int64(i) {
				return errs.ErrInternalServer.Wrap("seq is invalid")
			}
		case updateKeyRevoke:
			_, ok = field.(*unrelationtb.RevokeModel)
		default:
			return errs.ErrInternalServer.Wrap("key is invalid")
		}
		if !ok {
			return errs.ErrInternalServer.Wrap("field type is invalid")
		}
	}
	// 返回值为true表示数据库存在该文档，false表示数据库不存在该文档
	updateMsgModel := func(seq int64, i int) (bool, error) {
		var (
			res *mongo.UpdateResult
			err error
		)
		docID := db.msg.GetDocID(conversationID, seq)
		index := db.msg.GetMsgIndex(seq)
		field := fields[i]
		switch key {
		case updateKeyMsg:
			res, err = db.msgDocDatabase.UpdateMsg(ctx, docID, index, "msg", field)
		case updateKeyRevoke:
			res, err = db.msgDocDatabase.UpdateMsg(ctx, docID, index, "revoke", field)
		}
		if err != nil {
			return false, err
		}
		return res.MatchedCount > 0, nil
	}
	tryUpdate := true
	for i := 0; i < len(fields); i++ {
		seq := firstSeq + int64(i) // 当前seq
		if tryUpdate {
			matched, err := updateMsgModel(seq, i)
			if err != nil {
				return err
			}
			if matched {
				continue // 匹配到了，继续下一个(不一定修改)
			}
		}
		doc := unrelationtb.MsgDocModel{
			DocID: db.msg.GetDocID(conversationID, seq),
			Msg:   make([]*unrelationtb.MsgInfoModel, num),
		}
		var insert int // 插入的数量
		for j := i; j < len(fields); j++ {
			seq = firstSeq + int64(j)
			if db.msg.GetDocID(conversationID, seq) != doc.DocID {
				break
			}
			insert++
			switch key {
			case updateKeyMsg:
				doc.Msg[db.msg.GetMsgIndex(seq)] = &unrelationtb.MsgInfoModel{
					Msg: fields[j].(*unrelationtb.MsgDataModel),
				}
			case updateKeyRevoke:
				doc.Msg[db.msg.GetMsgIndex(seq)] = &unrelationtb.MsgInfoModel{
					Revoke: fields[j].(*unrelationtb.RevokeModel),
				}
			}
		}
		for i, model := range doc.Msg {
			if model == nil {
				model = &unrelationtb.MsgInfoModel{}
				doc.Msg[i] = model
			}
			if model.DelList == nil {
				doc.Msg[i].DelList = []string{}
			}
		}
		if err := db.msgDocDatabase.Create(ctx, &doc); err != nil {
			if mongo.IsDuplicateKeyError(err) {
				i--              // 存在并发,重试当前数据
				tryUpdate = true // 以修改模式
				continue
			}
			return err
		}
		tryUpdate = false // 当前以插入成功,下一块优先插入模式
		i += insert - 1   // 跳过已插入的数据
	}
	return nil
}

func (db *commonMsgDatabase) BatchInsertChat2DB(ctx context.Context, conversationID string, msgList []*sdkws.MsgData, currentMaxSeq int64) error {
	if len(msgList) == 0 {
		return errs.ErrArgs.Wrap("msgList is empty")
	}
	msgs := make([]any, len(msgList))
	for i, msg := range msgList {
		if msg == nil {
			continue
		}
		var offlinePushModel *unrelationtb.OfflinePushModel
		if msg.OfflinePushInfo != nil {
			offlinePushModel = &unrelationtb.OfflinePushModel{
				Title:         msg.OfflinePushInfo.Title,
				Desc:          msg.OfflinePushInfo.Desc,
				Ex:            msg.OfflinePushInfo.Ex,
				IOSPushSound:  msg.OfflinePushInfo.IOSPushSound,
				IOSBadgeCount: msg.OfflinePushInfo.IOSBadgeCount,
			}
		}
		msgs[i] = &unrelationtb.MsgDataModel{
			SendID:           msg.SendID,
			RecvID:           msg.RecvID,
			GroupID:          msg.GroupID,
			ClientMsgID:      msg.ClientMsgID,
			ServerMsgID:      msg.ServerMsgID,
			SenderPlatformID: msg.SenderPlatformID,
			SenderNickname:   msg.SenderNickname,
			SenderFaceURL:    msg.SenderFaceURL,
			SessionType:      msg.SessionType,
			MsgFrom:          msg.MsgFrom,
			ContentType:      msg.ContentType,
			Content:          string(msg.Content),
			Seq:              msg.Seq,
			SendTime:         msg.SendTime,
			CreateTime:       msg.CreateTime,
			Status:           msg.Status,
			Options:          msg.Options,
			OfflinePush:      offlinePushModel,
			AtUserIDList:     msg.AtUserIDList,
			AttachedInfo:     msg.AttachedInfo,
			Ex:               msg.Ex,
		}
	}
	return db.BatchInsertBlock(ctx, conversationID, msgs, updateKeyMsg, msgList[0].Seq)
}

func (db *commonMsgDatabase) RevokeMsg(ctx context.Context, conversationID string, seq int64, revoke *unrelationtb.RevokeModel) error {
	return db.BatchInsertBlock(ctx, conversationID, []any{revoke}, updateKeyRevoke, seq)
}

func (db *commonMsgDatabase) MarkSingleChatMsgsAsRead(ctx context.Context, userID string, conversationID string, totalSeqs []int64) error {
	for docID, seqs := range db.msg.GetDocIDSeqsMap(conversationID, totalSeqs) {
		var indexes []int64
		for _, seq := range seqs {
			indexes = append(indexes, db.msg.GetMsgIndex(seq))
		}
		log.ZDebug(ctx, "MarkSingleChatMsgsAsRead", "userID", userID, "docID", docID, "indexes", indexes)
		if err := db.msgDocDatabase.MarkSingleChatMsgsAsRead(ctx, userID, docID, indexes); err != nil {
			log.ZError(ctx, "MarkSingleChatMsgsAsRead", err, "userID", userID, "docID", docID, "indexes", indexes)
			return err
		}
	}
	return nil
}

func (db *commonMsgDatabase) DeleteMessagesFromCache(ctx context.Context, conversationID string, seqs []int64) error {
	return db.cache.DeleteMessages(ctx, conversationID, seqs)
}

func (db *commonMsgDatabase) DelUserDeleteMsgsList(ctx context.Context, conversationID string, seqs []int64) {
	db.cache.DelUserDeleteMsgsList(ctx, conversationID, seqs)
}

func (db *commonMsgDatabase) BatchInsertChat2Cache(ctx context.Context, conversationID string, msgs []*sdkws.MsgData) (seq int64, isNew bool, err error) {
	currentMaxSeq, err := db.cache.GetMaxSeq(ctx, conversationID)
	if err != nil && errs.Unwrap(err) != redis.Nil {
		log.ZError(ctx, "db.cache.GetMaxSeq", err)
		return 0, false, err
	}
	lenList := len(msgs)
	if int64(lenList) > db.msg.GetSingleGocMsgNum() {
		return 0, false, errors.New("too large")
	}
	if lenList < 1 {
		return 0, false, errors.New("too short as 0")
	}
	if errs.Unwrap(err) == redis.Nil {
		isNew = true
	}
	lastMaxSeq := currentMaxSeq
	userSeqMap := make(map[string]int64)
	for _, m := range msgs {
		currentMaxSeq++
		m.Seq = currentMaxSeq
		userSeqMap[m.SendID] = m.Seq
	}
	failedNum, err := db.cache.SetMessageToCache(ctx, conversationID, msgs)
	if err != nil {
		prommetrics.MsgInsertRedisFailedCounter.Add(float64(failedNum))
		log.ZError(ctx, "setMessageToCache error", err, "len", len(msgs), "conversationID", conversationID)
	} else {
		prommetrics.MsgInsertRedisSuccessCounter.Inc()
	}
	err = db.cache.SetMaxSeq(ctx, conversationID, currentMaxSeq)
	if err != nil {
		log.ZError(ctx, "db.cache.SetMaxSeq error", err, "conversationID", conversationID)
		prommetrics.SeqSetFailedCounter.Inc()
	}
	err2 := db.cache.SetHasReadSeqs(ctx, conversationID, userSeqMap)
	if err != nil {
		log.ZError(ctx, "SetHasReadSeqs error", err2, "userSeqMap", userSeqMap, "conversationID", conversationID)
		prommetrics.SeqSetFailedCounter.Inc()
	}
	return lastMaxSeq, isNew, utils.Wrap(err, "")
}

func (db *commonMsgDatabase) getMsgBySeqs(ctx context.Context, userID, conversationID string, seqs []int64) (totalMsgs []*sdkws.MsgData, err error) {
	for docID, seqs := range db.msg.GetDocIDSeqsMap(conversationID, seqs) {
		// log.ZDebug(ctx, "getMsgBySeqs", "docID", docID, "seqs", seqs)
		msgs, err := db.findMsgInfoBySeq(ctx, userID, docID, seqs)
		if err != nil {
			return nil, err
		}
		for _, msg := range msgs {
			totalMsgs = append(totalMsgs, convert.MsgDB2Pb(msg.Msg))
		}
	}
	return totalMsgs, nil
}

func (db *commonMsgDatabase) findMsgInfoBySeq(ctx context.Context, userID, docID string, seqs []int64) (totalMsgs []*unrelationtb.MsgInfoModel, err error) {
	msgs, err := db.msgDocDatabase.GetMsgBySeqIndexIn1Doc(ctx, userID, docID, seqs)
	for _, msg := range msgs {
		if msg.IsRead {
			msg.Msg.IsRead = true
		}
	}
	return msgs, err
}

func (db *commonMsgDatabase) getMsgBySeqsRange(ctx context.Context, userID string, conversationID string, allSeqs []int64, begin, end int64) (seqMsgs []*sdkws.MsgData, err error) {
	log.ZDebug(ctx, "getMsgBySeqsRange", "conversationID", conversationID, "allSeqs", allSeqs, "begin", begin, "end", end)
	for docID, seqs := range db.msg.GetDocIDSeqsMap(conversationID, allSeqs) {
		log.ZDebug(ctx, "getMsgBySeqsRange", "docID", docID, "seqs", seqs)
		msgs, err := db.findMsgInfoBySeq(ctx, userID, docID, seqs)
		if err != nil {
			return nil, err
		}
		for _, msg := range msgs {
			if msg.IsRead {
				msg.Msg.IsRead = true
			}
			seqMsgs = append(seqMsgs, convert.MsgDB2Pb(msg.Msg))
		}
	}
	return seqMsgs, nil
}

// GetMsgBySeqsRange In the context of group chat, we have the following parameters:
//
// "maxSeq" of a conversation: It represents the maximum value of messages in the group conversation.
// "minSeq" of a conversation (default: 1): It represents the minimum value of messages in the group conversation.
//
// For a user's perspective regarding the group conversation, we have the following parameters:
//
// "userMaxSeq": It represents the user's upper limit for message retrieval in the group. If not set (default: 0),
// it means the upper limit is the same as the conversation's "maxSeq".
// "userMinSeq": It represents the user's starting point for message retrieval in the group. If not set (default: 0),
// it means the starting point is the same as the conversation's "minSeq".
//
// The scenarios for these parameters are as follows:
//
// For users who have been kicked out of the group, "userMaxSeq" can be set as the maximum value they had before
// being kicked out. This limits their ability to retrieve messages up to a certain point.
// For new users joining the group, if they don't need to receive old messages,
// "userMinSeq" can be set as the same value as the conversation's "maxSeq" at the moment they join the group.
// This ensures that their message retrieval starts from the point they joined.
func (db *commonMsgDatabase) GetMsgBySeqsRange(ctx context.Context, userID string, conversationID string, begin, end, num, userMaxSeq int64) (int64, int64, []*sdkws.MsgData, error) {
	userMinSeq, err := db.cache.GetConversationUserMinSeq(ctx, conversationID, userID)
	if err != nil && errs.Unwrap(err) != redis.Nil {
		return 0, 0, nil, err
	}
	minSeq, err := db.cache.GetMinSeq(ctx, conversationID)
	if err != nil && errs.Unwrap(err) != redis.Nil {
		return 0, 0, nil, err
	}
	if userMinSeq > minSeq {
		minSeq = userMinSeq
	}
	//"minSeq" represents the startSeq value that the user can retrieve.
	if minSeq > end {
		log.ZInfo(ctx, "minSeq > end", "minSeq", minSeq, "end", end)
		return 0, 0, nil, nil
	}
	maxSeq, err := db.cache.GetMaxSeq(ctx, conversationID)
	if err != nil && errs.Unwrap(err) != redis.Nil {
		return 0, 0, nil, err
	}
	log.ZDebug(ctx, "GetMsgBySeqsRange", "userMinSeq", userMinSeq, "conMinSeq", minSeq, "conMaxSeq", maxSeq, "userMaxSeq", userMaxSeq)
	if userMaxSeq != 0 {
		if userMaxSeq < maxSeq {
			maxSeq = userMaxSeq
		}
	}
	//"maxSeq" represents the endSeq value that the user can retrieve.

	if begin < minSeq {
		begin = minSeq
	}
	if end > maxSeq {
		end = maxSeq
	}
	//"begin" and "end" represent the actual startSeq and endSeq values that the user can retrieve.
	if end < begin {
		return 0, 0, nil, errs.ErrArgs.Wrap("seq end < begin")
	}
	var seqs []int64
	if end-begin+1 <= num {
		for i := begin; i <= end; i++ {
			seqs = append(seqs, i)
		}
	} else {
		for i := end - num + 1; i <= end; i++ {
			seqs = append(seqs, i)
		}
	}

	//167 178 10
	//if end-num <  {
	//
	//}
	//var seqs []int64
	//for i := end; i > end-num; i-- {
	//	if i >= begin {
	//		seqs = append([]int64{i}, seqs...)
	//	} else {
	//		break
	//	}
	//}
	if len(seqs) == 0 {
		return 0, 0, nil, nil
	}
	newBegin := seqs[0]
	newEnd := seqs[len(seqs)-1]
	log.ZDebug(ctx, "GetMsgBySeqsRange", "first seqs", seqs, "newBegin", newBegin, "newEnd", newEnd)
	cachedMsgs, failedSeqs, err := db.cache.GetMessagesBySeq(ctx, conversationID, seqs)
	if err != nil {
		if err != redis.Nil {

			log.ZError(ctx, "get message from redis exception", err, "conversationID", conversationID, "seqs", seqs)
		}
	}
	var successMsgs []*sdkws.MsgData
	if len(cachedMsgs) > 0 {
		delSeqs, err := db.cache.GetUserDelList(ctx, userID, conversationID)
		if err != nil && errs.Unwrap(err) != redis.Nil {
			return 0, 0, nil, err
		}
		var cacheDelNum int
		for _, msg := range cachedMsgs {
			if !utils.Contain(msg.Seq, delSeqs...) {
				successMsgs = append(successMsgs, msg)
			} else {
				cacheDelNum += 1
			}
		}
		log.ZDebug(ctx, "get delSeqs from redis", "delSeqs", delSeqs, "userID", userID, "conversationID", conversationID, "cacheDelNum", cacheDelNum)
		var reGetSeqsCache []int64
		for i := 1; i <= cacheDelNum; {
			newSeq := newBegin - int64(i)
			if newSeq >= begin {
				if !utils.Contain(newSeq, delSeqs...) {
					log.ZDebug(ctx, "seq del in cache, a new seq in range append", "new seq", newSeq)
					reGetSeqsCache = append(reGetSeqsCache, newSeq)
					i++
				}
			} else {
				break
			}
		}
		if len(reGetSeqsCache) > 0 {
			log.ZDebug(ctx, "reGetSeqsCache", "reGetSeqsCache", reGetSeqsCache)
			cachedMsgs, failedSeqs2, err := db.cache.GetMessagesBySeq(ctx, conversationID, reGetSeqsCache)
			if err != nil {
				if err != redis.Nil {

					log.ZError(ctx, "get message from redis exception", err, "conversationID", conversationID, "seqs", reGetSeqsCache)
				}
			}
			failedSeqs = append(failedSeqs, failedSeqs2...)
			successMsgs = append(successMsgs, cachedMsgs...)
		}
	}
	log.ZDebug(ctx, "get msgs from cache", "successMsgs", successMsgs)
	if len(failedSeqs) != 0 {
		log.ZDebug(ctx, "msgs not exist in redis", "seqs", failedSeqs)
	}
	// get from cache or db

	if len(failedSeqs) > 0 {
		mongoMsgs, err := db.getMsgBySeqsRange(ctx, userID, conversationID, failedSeqs, begin, end)
		if err != nil {

			return 0, 0, nil, err
		}
		successMsgs = append(mongoMsgs, successMsgs...)
	}

	return minSeq, maxSeq, successMsgs, nil
}

func (db *commonMsgDatabase) GetMsgBySeqs(ctx context.Context, userID string, conversationID string, seqs []int64) (int64, int64, []*sdkws.MsgData, error) {
	userMinSeq, err := db.cache.GetConversationUserMinSeq(ctx, conversationID, userID)
	if err != nil && errs.Unwrap(err) != redis.Nil {
		return 0, 0, nil, err
	}
	minSeq, err := db.cache.GetMinSeq(ctx, conversationID)
	if err != nil && errs.Unwrap(err) != redis.Nil {
		return 0, 0, nil, err
	}
	maxSeq, err := db.cache.GetMaxSeq(ctx, conversationID)
	if err != nil && errs.Unwrap(err) != redis.Nil {
		return 0, 0, nil, err
	}
	if userMinSeq < minSeq {
		minSeq = userMinSeq
	}
	var newSeqs []int64
	for _, seq := range seqs {
		if seq >= minSeq && seq <= maxSeq {
			newSeqs = append(newSeqs, seq)
		}
	}
	successMsgs, failedSeqs, err := db.cache.GetMessagesBySeq(ctx, conversationID, newSeqs)
	if err != nil {
		if err != redis.Nil {
			log.ZError(ctx, "get message from redis exception", err, "failedSeqs", failedSeqs, "conversationID", conversationID)
		}
	}
	log.ZInfo(
		ctx,
		"db.cache.GetMessagesBySeq",
		"userID",
		userID,
		"conversationID",
		conversationID,
		"seqs",
		seqs,
		"successMsgs",
		len(successMsgs),
		"failedSeqs",
		failedSeqs,
		"conversationID",
		conversationID,
	)

	if len(failedSeqs) > 0 {
		mongoMsgs, err := db.getMsgBySeqs(ctx, userID, conversationID, failedSeqs)
		if err != nil {

			return 0, 0, nil, err
		}

		successMsgs = append(successMsgs, mongoMsgs...)
	}
	return minSeq, maxSeq, successMsgs, nil
}

func (db *commonMsgDatabase) DeleteConversationMsgsAndSetMinSeq(ctx context.Context, conversationID string, remainTime int64) error {
	var delStruct delMsgRecursionStruct
	var skip int64
	minSeq, err := db.deleteMsgRecursion(ctx, conversationID, skip, &delStruct, remainTime)
	if err != nil {
		return err
	}
	log.ZInfo(ctx, "DeleteConversationMsgsAndSetMinSeq", "conversationID", conversationID, "minSeq", minSeq)
	if minSeq == 0 {
		return nil
	}
	if remainTime == 0 {
		err = db.cache.CleanUpOneConversationAllMsg(ctx, conversationID)
		if err != nil {
			log.ZWarn(ctx, "CleanUpOneUserAllMsg", err, "conversationID", conversationID)
		}
	}
	return db.cache.SetMinSeq(ctx, conversationID, minSeq)
}

func (db *commonMsgDatabase) UserMsgsDestruct(ctx context.Context, userID string, conversationID string, destructTime int64, lastMsgDestructTime time.Time) (seqs []int64, err error) {
	var index int64
	for {
		// from oldest 2 newest
		msgDocModel, err := db.msgDocDatabase.GetMsgDocModelByIndex(ctx, conversationID, index, 1)
		if err != nil || msgDocModel.DocID == "" {
			if err != nil {
				if err == unrelation.ErrMsgListNotExist {
					log.ZDebug(ctx, "not doc find", "conversationID", conversationID, "userID", userID, "index", index)
				} else {
					log.ZError(ctx, "deleteMsgRecursion GetUserMsgListByIndex failed", err, "conversationID", conversationID, "index", index)
				}
			}
			// 获取报错，或者获取不到了，物理删除并且返回seq delMongoMsgsPhysical(delStruct.delDocIDList), 结束递归
			break
		}
		index++
		//&& msgDocModel.Msg[0].Msg.SendTime > lastMsgDestructTime.UnixMilli()
		if len(msgDocModel.Msg) > 0 {
			i := 0
			var over bool
			for _, msg := range msgDocModel.Msg {
				i++
				if msg != nil && msg.Msg != nil && msg.Msg.SendTime+destructTime*1000 <= time.Now().UnixMilli() {
					if msg.Msg.SendTime+destructTime*1000 > lastMsgDestructTime.UnixMilli() && !utils.Contain(userID, msg.DelList...) {
						seqs = append(seqs, msg.Msg.Seq)
					}
				} else {
					log.ZDebug(ctx, "all msg need destruct is found", "conversationID", conversationID, "userID", userID, "index", index, "stop index", i)
					over = true
					break
				}
			}
			if over {
				break
			}
		}
	}

	log.ZDebug(ctx, "UserMsgsDestruct", "conversationID", conversationID, "userID", userID, "seqs", seqs)
	if len(seqs) > 0 {
		userMinSeq := seqs[len(seqs)-1] + 1
		currentUserMinSeq, err := db.cache.GetConversationUserMinSeq(ctx, conversationID, userID)
		if err != nil && errs.Unwrap(err) != redis.Nil {
			return nil, err
		}
		if currentUserMinSeq < userMinSeq {
			if err := db.cache.SetConversationUserMinSeq(ctx, conversationID, userID, userMinSeq); err != nil {
				return nil, err
			}
		}
	}
	return seqs, nil
}

// this is struct for recursion.
type delMsgRecursionStruct struct {
	minSeq    int64
	delDocIDs []string
}

func (d *delMsgRecursionStruct) getSetMinSeq() int64 {
	return d.minSeq
}

// index 0....19(del) 20...69
// seq 70
// set minSeq 21
// recursion 删除list并且返回设置的最小seq.
func (db *commonMsgDatabase) deleteMsgRecursion(ctx context.Context, conversationID string, index int64, delStruct *delMsgRecursionStruct, remainTime int64) (int64, error) {
	// find from oldest list
	msgDocModel, err := db.msgDocDatabase.GetMsgDocModelByIndex(ctx, conversationID, index, 1)
	if err != nil || msgDocModel.DocID == "" {
		if err != nil {
			if err == unrelation.ErrMsgListNotExist {
				log.ZDebug(ctx, "deleteMsgRecursion ErrMsgListNotExist", "conversationID", conversationID, "index:", index)
			} else {
				log.ZError(ctx, "deleteMsgRecursion GetUserMsgListByIndex failed", err, "conversationID", conversationID, "index", index)
			}
		}
		// 获取报错，或者获取不到了，物理删除并且返回seq delMongoMsgsPhysical(delStruct.delDocIDList), 结束递归
		err = db.msgDocDatabase.DeleteDocs(ctx, delStruct.delDocIDs)
		if err != nil {
			return 0, err
		}
		return delStruct.getSetMinSeq() + 1, nil
	}
	log.ZDebug(ctx, "doc info", "conversationID", conversationID, "index", index, "docID", msgDocModel.DocID, "len", len(msgDocModel.Msg))
	if int64(len(msgDocModel.Msg)) > db.msg.GetSingleGocMsgNum() {
		log.ZWarn(ctx, "msgs too large", nil, "lenth", len(msgDocModel.Msg), "docID:", msgDocModel.DocID)
	}
	if msgDocModel.IsFull() && msgDocModel.Msg[len(msgDocModel.Msg)-1].Msg.SendTime+(remainTime*1000) < utils.GetCurrentTimestampByMill() {
		log.ZDebug(ctx, "doc is full and all msg is expired", "docID", msgDocModel.DocID)
		delStruct.delDocIDs = append(delStruct.delDocIDs, msgDocModel.DocID)
		delStruct.minSeq = msgDocModel.Msg[len(msgDocModel.Msg)-1].Msg.Seq
	} else {
		var delMsgIndexs []int
		for i, MsgInfoModel := range msgDocModel.Msg {
			if MsgInfoModel != nil && MsgInfoModel.Msg != nil {
				if utils.GetCurrentTimestampByMill() > MsgInfoModel.Msg.SendTime+(remainTime*1000) {
					delMsgIndexs = append(delMsgIndexs, i)
				}
			}
		}
		if len(delMsgIndexs) > 0 {
			if err := db.msgDocDatabase.DeleteMsgsInOneDocByIndex(ctx, msgDocModel.DocID, delMsgIndexs); err != nil {
				log.ZError(ctx, "deleteMsgRecursion DeleteMsgsInOneDocByIndex failed", err, "conversationID", conversationID, "index", index)
			}
			delStruct.minSeq = int64(msgDocModel.Msg[delMsgIndexs[len(delMsgIndexs)-1]].Msg.Seq)
		}
	}
	seq, err := db.deleteMsgRecursion(ctx, conversationID, index+1, delStruct, remainTime)
	return seq, err
}

func (db *commonMsgDatabase) DeleteMsgsPhysicalBySeqs(ctx context.Context, conversationID string, allSeqs []int64) error {
	if err := db.cache.DeleteMessages(ctx, conversationID, allSeqs); err != nil {
		return err
	}
	for docID, seqs := range db.msg.GetDocIDSeqsMap(conversationID, allSeqs) {
		var indexes []int
		for _, seq := range seqs {
			indexes = append(indexes, int(db.msg.GetMsgIndex(seq)))
		}
		if err := db.msgDocDatabase.DeleteMsgsInOneDocByIndex(ctx, docID, indexes); err != nil {
			return err
		}
	}
	return nil
}

func (db *commonMsgDatabase) DeleteUserMsgsBySeqs(ctx context.Context, userID string, conversationID string, seqs []int64) error {
	cachedMsgs, _, err := db.cache.GetMessagesBySeq(ctx, conversationID, seqs)
	if err != nil && errs.Unwrap(err) != redis.Nil {
		log.ZWarn(ctx, "DeleteUserMsgsBySeqs", err, "conversationID", conversationID, "seqs", seqs)
		return err
	}
	if len(cachedMsgs) > 0 {
		var cacheSeqs []int64
		for _, msg := range cachedMsgs {
			cacheSeqs = append(cacheSeqs, msg.Seq)
		}
		if err := db.cache.UserDeleteMsgs(ctx, conversationID, cacheSeqs, userID); err != nil {
			return err
		}
	}

	for docID, seqs := range db.msg.GetDocIDSeqsMap(conversationID, seqs) {
		for _, seq := range seqs {
			if _, err := db.msgDocDatabase.PushUnique(ctx, docID, db.msg.GetMsgIndex(seq), "del_list", []string{userID}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (db *commonMsgDatabase) DeleteMsgsBySeqs(ctx context.Context, conversationID string, seqs []int64) error {
	return nil
}

func (db *commonMsgDatabase) CleanUpUserConversationsMsgs(ctx context.Context, user string, conversationIDs []string) {
	for _, conversationID := range conversationIDs {
		maxSeq, err := db.cache.GetMaxSeq(ctx, conversationID)
		if err != nil {
			if err == redis.Nil {
				log.ZInfo(ctx, "max seq is nil", "conversationID", conversationID)
			} else {
				log.ZError(ctx, "get max seq failed", err, "conversationID", conversationID)
			}
			continue
		}
		if err := db.cache.SetMinSeq(ctx, conversationID, maxSeq+1); err != nil {
			log.ZError(ctx, "set min seq failed", err, "conversationID", conversationID, "minSeq", maxSeq+1)
		}
	}
}

func (db *commonMsgDatabase) SetMaxSeq(ctx context.Context, conversationID string, maxSeq int64) error {
	return db.cache.SetMaxSeq(ctx, conversationID, maxSeq)
}

func (db *commonMsgDatabase) GetMaxSeqs(ctx context.Context, conversationIDs []string) (map[string]int64, error) {
	return db.cache.GetMaxSeqs(ctx, conversationIDs)
}

func (db *commonMsgDatabase) GetMaxSeq(ctx context.Context, conversationID string) (int64, error) {
	return db.cache.GetMaxSeq(ctx, conversationID)
}

func (db *commonMsgDatabase) SetMinSeq(ctx context.Context, conversationID string, minSeq int64) error {
	return db.cache.SetMinSeq(ctx, conversationID, minSeq)
}

func (db *commonMsgDatabase) SetMinSeqs(ctx context.Context, seqs map[string]int64) error {
	return db.cache.SetMinSeqs(ctx, seqs)
}

func (db *commonMsgDatabase) GetMinSeqs(ctx context.Context, conversationIDs []string) (map[string]int64, error) {
	return db.cache.GetMinSeqs(ctx, conversationIDs)
}

func (db *commonMsgDatabase) GetMinSeq(ctx context.Context, conversationID string) (int64, error) {
	return db.cache.GetMinSeq(ctx, conversationID)
}

func (db *commonMsgDatabase) GetConversationUserMinSeq(ctx context.Context, conversationID string, userID string) (int64, error) {
	return db.cache.GetConversationUserMinSeq(ctx, conversationID, userID)
}

func (db *commonMsgDatabase) GetConversationUserMinSeqs(ctx context.Context, conversationID string, userIDs []string) (map[string]int64, error) {
	return db.cache.GetConversationUserMinSeqs(ctx, conversationID, userIDs)
}

func (db *commonMsgDatabase) SetConversationUserMinSeq(ctx context.Context, conversationID string, userID string, minSeq int64) error {
	return db.cache.SetConversationUserMinSeq(ctx, conversationID, userID, minSeq)
}

func (db *commonMsgDatabase) SetConversationUserMinSeqs(ctx context.Context, conversationID string, seqs map[string]int64) (err error) {
	return db.cache.SetConversationUserMinSeqs(ctx, conversationID, seqs)
}

func (db *commonMsgDatabase) SetUserConversationsMinSeqs(ctx context.Context, userID string, seqs map[string]int64) error {
	return db.cache.SetUserConversationsMinSeqs(ctx, userID, seqs)
}

func (db *commonMsgDatabase) UserSetHasReadSeqs(ctx context.Context, userID string, hasReadSeqs map[string]int64) error {
	return db.cache.UserSetHasReadSeqs(ctx, userID, hasReadSeqs)
}

func (db *commonMsgDatabase) SetHasReadSeq(ctx context.Context, userID string, conversationID string, hasReadSeq int64) error {
	return db.cache.SetHasReadSeq(ctx, userID, conversationID, hasReadSeq)
}

func (db *commonMsgDatabase) GetHasReadSeqs(ctx context.Context, userID string, conversationIDs []string) (map[string]int64, error) {
	return db.cache.GetHasReadSeqs(ctx, userID, conversationIDs)
}

func (db *commonMsgDatabase) GetHasReadSeq(ctx context.Context, userID string, conversationID string) (int64, error) {
	return db.cache.GetHasReadSeq(ctx, userID, conversationID)
}

func (db *commonMsgDatabase) SetSendMsgStatus(ctx context.Context, id string, status int32) error {
	return db.cache.SetSendMsgStatus(ctx, id, status)
}

func (db *commonMsgDatabase) GetSendMsgStatus(ctx context.Context, id string) (int32, error) {
	return db.cache.GetSendMsgStatus(ctx, id)
}

func (db *commonMsgDatabase) GetConversationMinMaxSeqInMongoAndCache(ctx context.Context, conversationID string) (minSeqMongo, maxSeqMongo, minSeqCache, maxSeqCache int64, err error) {
	minSeqMongo, maxSeqMongo, err = db.GetMinMaxSeqMongo(ctx, conversationID)
	if err != nil {
		return
	}
	minSeqCache, err = db.cache.GetMinSeq(ctx, conversationID)
	if err != nil {
		return
	}
	maxSeqCache, err = db.cache.GetMaxSeq(ctx, conversationID)
	if err != nil {
		return
	}
	return
}

func (db *commonMsgDatabase) GetMongoMaxAndMinSeq(ctx context.Context, conversationID string) (minSeqMongo, maxSeqMongo int64, err error) {
	return db.GetMinMaxSeqMongo(ctx, conversationID)
}

func (db *commonMsgDatabase) GetMinMaxSeqMongo(ctx context.Context, conversationID string) (minSeqMongo, maxSeqMongo int64, err error) {
	oldestMsgMongo, err := db.msgDocDatabase.GetOldestMsg(ctx, conversationID)
	if err != nil {
		return
	}
	minSeqMongo = oldestMsgMongo.Msg.Seq
	newestMsgMongo, err := db.msgDocDatabase.GetNewestMsg(ctx, conversationID)
	if err != nil {
		return
	}
	maxSeqMongo = newestMsgMongo.Msg.Seq
	return
}

func (db *commonMsgDatabase) RangeUserSendCount(
	ctx context.Context,
	start time.Time,
	end time.Time,
	group bool,
	ase bool,
	pageNumber int32,
	showNumber int32,
) (msgCount int64, userCount int64, users []*unrelationtb.UserCount, dateCount map[string]int64, err error) {
	return db.msgDocDatabase.RangeUserSendCount(ctx, start, end, group, ase, pageNumber, showNumber)
}

func (db *commonMsgDatabase) RangeGroupSendCount(
	ctx context.Context,
	start time.Time,
	end time.Time,
	ase bool,
	pageNumber int32,
	showNumber int32,
) (msgCount int64, userCount int64, groups []*unrelationtb.GroupCount, dateCount map[string]int64, err error) {
	return db.msgDocDatabase.RangeGroupSendCount(ctx, start, end, ase, pageNumber, showNumber)
}

func (db *commonMsgDatabase) SearchMessage(ctx context.Context, req *pbmsg.SearchMessageReq) (total int32, msgData []*sdkws.MsgData, err error) {
	var totalMsgs []*sdkws.MsgData
	total, msgs, err := db.msgDocDatabase.SearchMessage(ctx, req)
	if err != nil {
		return 0, nil, err
	}
	for _, msg := range msgs {
		if msg.IsRead {
			msg.Msg.IsRead = true
		}
		totalMsgs = append(totalMsgs, convert.MsgDB2Pb(msg.Msg))
	}
	return total, totalMsgs, nil
}

func (db *commonMsgDatabase) ConvertMsgsDocLen(ctx context.Context, conversationIDs []string) {
	db.msgDocDatabase.ConvertMsgsDocLen(ctx, conversationIDs)
}
