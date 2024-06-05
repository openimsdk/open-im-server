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
	"encoding/json"
	"errors"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"strings"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/convert"
	"github.com/openimsdk/open-im-server/v3/pkg/common/prommetrics"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/protocol/constant"
	pbmsg "github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mq/kafka"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/openimsdk/tools/utils/timeutil"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	updateKeyMsg = iota
	updateKeyRevoke
)

// CommonMsgDatabase defines the interface for message database operations.
type CommonMsgDatabase interface {
	// BatchInsertChat2DB inserts a batch of messages into the database for a specific conversation.
	BatchInsertChat2DB(ctx context.Context, conversationID string, msgs []*sdkws.MsgData, currentMaxSeq int64) error
	// RevokeMsg revokes a message in a conversation.
	RevokeMsg(ctx context.Context, conversationID string, seq int64, revoke *model.RevokeModel) error
	// MarkSingleChatMsgsAsRead marks messages as read for a single chat by sequence numbers.
	MarkSingleChatMsgsAsRead(ctx context.Context, userID string, conversationID string, seqs []int64) error
	// DeleteMessagesFromCache deletes message caches from Redis by sequence numbers.
	DeleteMessagesFromCache(ctx context.Context, conversationID string, seqs []int64) error
	// BatchInsertChat2Cache increments the sequence number and then batch inserts messages into the cache.
	BatchInsertChat2Cache(ctx context.Context, conversationID string, msgs []*sdkws.MsgData) (seq int64, isNewConversation bool, err error)
	// GetMsgBySeqsRange retrieves messages from MongoDB by a range of sequence numbers.
	GetMsgBySeqsRange(ctx context.Context, userID string, conversationID string, begin, end, num, userMaxSeq int64) (minSeq int64, maxSeq int64, seqMsg []*sdkws.MsgData, err error)
	// GetMsgBySeqs retrieves messages for large groups from MongoDB by sequence numbers.
	GetMsgBySeqs(ctx context.Context, userID string, conversationID string, seqs []int64) (minSeq int64, maxSeq int64, seqMsg []*sdkws.MsgData, err error)
	// DeleteConversationMsgsAndSetMinSeq deletes conversation messages and resets the minimum sequence number. If `remainTime` is 0, all messages are deleted (this method does not delete Redis
	// cache).
	DeleteConversationMsgsAndSetMinSeq(ctx context.Context, conversationID string, remainTime int64) error
	// UserMsgsDestruct marks messages for deletion based on destruct time and returns a list of sequence numbers for marked messages.
	UserMsgsDestruct(ctx context.Context, userID string, conversationID string, destructTime int64, lastMsgDestructTime time.Time) (seqs []int64, err error)
	// DeleteUserMsgsBySeqs allows a user to delete messages based on sequence numbers.
	DeleteUserMsgsBySeqs(ctx context.Context, userID string, conversationID string, seqs []int64) error
	// DeleteMsgsPhysicalBySeqs physically deletes messages by emptying them based on sequence numbers.
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
	FindOneByDocIDs(ctx context.Context, docIDs []string, seqs map[string]int64) (map[string]*sdkws.MsgData, error)

	// to mq
	MsgToMQ(ctx context.Context, key string, msg2mq *sdkws.MsgData) error
	MsgToPushMQ(ctx context.Context, key, conversarionID string, msg2mq *sdkws.MsgData) (int32, int64, error)
	MsgToMongoMQ(ctx context.Context, key, conversarionID string, msgs []*sdkws.MsgData, lastSeq int64) error

	RangeUserSendCount(ctx context.Context, start time.Time, end time.Time, group bool, ase bool, pageNumber int32, showNumber int32) (msgCount int64, userCount int64, users []*model.UserCount, dateCount map[string]int64, err error)
	RangeGroupSendCount(ctx context.Context, start time.Time, end time.Time, ase bool, pageNumber int32, showNumber int32) (msgCount int64, userCount int64, groups []*model.GroupCount, dateCount map[string]int64, err error)
	ConvertMsgsDocLen(ctx context.Context, conversationIDs []string)

	// clear msg
	GetBeforeMsg(ctx context.Context, ts int64, limit int) ([]*model.MsgDocModel, error)
	DeleteDocMsgBefore(ctx context.Context, ts int64, doc *model.MsgDocModel) ([]int, error)
}

func NewCommonMsgDatabase(msgDocModel database.Msg, msg cache.MsgCache, seq cache.SeqCache, kafkaConf *config.Kafka) (CommonMsgDatabase, error) {
	conf, err := kafka.BuildProducerConfig(*kafkaConf.Build())
	if err != nil {
		return nil, err
	}
	producerToRedis, err := kafka.NewKafkaProducer(conf, kafkaConf.Address, kafkaConf.ToRedisTopic)
	if err != nil {
		return nil, err
	}
	producerToMongo, err := kafka.NewKafkaProducer(conf, kafkaConf.Address, kafkaConf.ToMongoTopic)
	if err != nil {
		return nil, err
	}
	producerToPush, err := kafka.NewKafkaProducer(conf, kafkaConf.Address, kafkaConf.ToPushTopic)
	if err != nil {
		return nil, err
	}
	return &commonMsgDatabase{
		msgDocDatabase:  msgDocModel,
		msg:             msg,
		seq:             seq,
		producer:        producerToRedis,
		producerToMongo: producerToMongo,
		producerToPush:  producerToPush,
	}, nil
}

//func InitCommonMsgDatabase(rdb redis.UniversalClient, database *mongo.Database, config *tools.CronTaskConfig) (CommonMsgDatabase, error) {
//	msgDocModel, err := database.NewMsgMongo(database)
//	if err != nil {
//		return nil, err
//	}
//	//todo MsgCacheTimeout
//	msg := cache.NewMsgCache(rdb, 86400, config.RedisConfig.EnablePipeline)
//	seq := cache.NewSeqCache(rdb)
//	return NewCommonMsgDatabase(msgDocModel, msg, seq, &config.KafkaConfig)
//}

type commonMsgDatabase struct {
	msgDocDatabase  database.Msg
	msgTable        model.MsgDocModel
	msg             cache.MsgCache
	seq             cache.SeqCache
	producer        *kafka.Producer
	producerToMongo *kafka.Producer
	producerToPush  *kafka.Producer
}

func (db *commonMsgDatabase) MsgToMQ(ctx context.Context, key string, msg2mq *sdkws.MsgData) error {
	_, _, err := db.producer.SendMessage(ctx, key, msg2mq)
	return err
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
	num := db.msgTable.GetSingleGocMsgNum()
	// num = 100
	for i, field := range fields { // Check the type of the field
		var ok bool
		switch key {
		case updateKeyMsg:
			var msg *model.MsgDataModel
			msg, ok = field.(*model.MsgDataModel)
			if msg != nil && msg.Seq != firstSeq+int64(i) {
				return errs.ErrInternalServer.WrapMsg("seq is invalid")
			}
		case updateKeyRevoke:
			_, ok = field.(*model.RevokeModel)
		default:
			return errs.ErrInternalServer.WrapMsg("key is invalid")
		}
		if !ok {
			return errs.ErrInternalServer.WrapMsg("field type is invalid")
		}
	}
	// Returns true if the document exists in the database, false if the document does not exist in the database
	updateMsgModel := func(seq int64, i int) (bool, error) {
		var (
			res *mongo.UpdateResult
			err error
		)
		docID := db.msgTable.GetDocID(conversationID, seq)
		index := db.msgTable.GetMsgIndex(seq)
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
		seq := firstSeq + int64(i) // Current sequence number
		if tryUpdate {
			matched, err := updateMsgModel(seq, i)
			if err != nil {
				return err
			}
			if matched {
				continue // The current data has been updated, skip the current data
			}
		}
		doc := model.MsgDocModel{
			DocID: db.msgTable.GetDocID(conversationID, seq),
			Msg:   make([]*model.MsgInfoModel, num),
		}
		var insert int // Inserted data number
		for j := i; j < len(fields); j++ {
			seq = firstSeq + int64(j)
			if db.msgTable.GetDocID(conversationID, seq) != doc.DocID {
				break
			}
			insert++
			switch key {
			case updateKeyMsg:
				doc.Msg[db.msgTable.GetMsgIndex(seq)] = &model.MsgInfoModel{
					Msg: fields[j].(*model.MsgDataModel),
				}
			case updateKeyRevoke:
				doc.Msg[db.msgTable.GetMsgIndex(seq)] = &model.MsgInfoModel{
					Revoke: fields[j].(*model.RevokeModel),
				}
			}
		}
		for i, msgInfo := range doc.Msg {
			if msgInfo == nil {
				msgInfo = &model.MsgInfoModel{}
				doc.Msg[i] = msgInfo
			}
			if msgInfo.DelList == nil {
				doc.Msg[i].DelList = []string{}
			}
		}
		if err := db.msgDocDatabase.Create(ctx, &doc); err != nil {
			if mongo.IsDuplicateKeyError(err) {
				i--              // already inserted
				tryUpdate = true // next block use update mode
				continue
			}
			return err
		}
		tryUpdate = false // The current block is inserted successfully, and the next block is inserted preferentially
		i += insert - 1   // Skip the inserted data
	}
	return nil
}

func (db *commonMsgDatabase) BatchInsertChat2DB(ctx context.Context, conversationID string, msgList []*sdkws.MsgData, currentMaxSeq int64) error {
	if len(msgList) == 0 {
		return errs.ErrArgs.WrapMsg("msgList is empty")
	}
	msgs := make([]any, len(msgList))
	for i, msg := range msgList {
		if msg == nil {
			continue
		}
		var offlinePushModel *model.OfflinePushModel
		if msg.OfflinePushInfo != nil {
			offlinePushModel = &model.OfflinePushModel{
				Title:         msg.OfflinePushInfo.Title,
				Desc:          msg.OfflinePushInfo.Desc,
				Ex:            msg.OfflinePushInfo.Ex,
				IOSPushSound:  msg.OfflinePushInfo.IOSPushSound,
				IOSBadgeCount: msg.OfflinePushInfo.IOSBadgeCount,
			}
		}
		msgs[i] = &model.MsgDataModel{
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

func (db *commonMsgDatabase) RevokeMsg(ctx context.Context, conversationID string, seq int64, revoke *model.RevokeModel) error {
	return db.BatchInsertBlock(ctx, conversationID, []any{revoke}, updateKeyRevoke, seq)
}

func (db *commonMsgDatabase) MarkSingleChatMsgsAsRead(ctx context.Context, userID string, conversationID string, totalSeqs []int64) error {
	for docID, seqs := range db.msgTable.GetDocIDSeqsMap(conversationID, totalSeqs) {
		var indexes []int64
		for _, seq := range seqs {
			indexes = append(indexes, db.msgTable.GetMsgIndex(seq))
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
	return db.msg.DeleteMessagesFromCache(ctx, conversationID, seqs)
}

func (db *commonMsgDatabase) BatchInsertChat2Cache(ctx context.Context, conversationID string, msgs []*sdkws.MsgData) (seq int64, isNew bool, err error) {
	currentMaxSeq, err := db.seq.GetMaxSeq(ctx, conversationID)
	if err != nil && errs.Unwrap(err) != redis.Nil {
		log.ZError(ctx, "storage.seq.GetMaxSeq", err)
		return 0, false, err
	}
	lenList := len(msgs)
	if int64(lenList) > db.msgTable.GetSingleGocMsgNum() {
		return 0, false, errs.New("message count exceeds limit", "limit", db.msgTable.GetSingleGocMsgNum()).Wrap()
	}
	if lenList < 1 {
		return 0, false, errs.New("no messages to insert", "minCount", 1).Wrap()
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

	failedNum, err := db.msg.SetMessagesToCache(ctx, conversationID, msgs)
	if err != nil {
		prommetrics.MsgInsertRedisFailedCounter.Add(float64(failedNum))
		log.ZError(ctx, "setMessageToCache error", err, "len", len(msgs), "conversationID", conversationID)
	} else {
		prommetrics.MsgInsertRedisSuccessCounter.Inc()
	}

	err = db.seq.SetMaxSeq(ctx, conversationID, currentMaxSeq)
	if err != nil {
		log.ZError(ctx, "storage.seq.SetMaxSeq error", err, "conversationID", conversationID)
		prommetrics.SeqSetFailedCounter.Inc()
	}

	err = db.seq.SetHasReadSeqs(ctx, conversationID, userSeqMap)
	if err != nil {
		log.ZError(ctx, "SetHasReadSeqs error", err, "userSeqMap", userSeqMap, "conversationID", conversationID)
		prommetrics.SeqSetFailedCounter.Inc()
	}
	return lastMaxSeq, isNew, errs.Wrap(err)
}

func (db *commonMsgDatabase) getMsgBySeqs(ctx context.Context, userID, conversationID string, seqs []int64) (totalMsgs []*sdkws.MsgData, err error) {
	for docID, seqs := range db.msgTable.GetDocIDSeqsMap(conversationID, seqs) {
		// log.ZDebug(ctx, "getMsgBySeqs", "docID", docID, "seqs", seqs)
		msgs, err := db.findMsgInfoBySeq(ctx, userID, docID, conversationID, seqs)
		if err != nil {
			return nil, err
		}
		for _, msg := range msgs {
			totalMsgs = append(totalMsgs, convert.MsgDB2Pb(msg.Msg))
		}
	}
	return totalMsgs, nil
}

func (db *commonMsgDatabase) handlerDBMsg(ctx context.Context, cache map[int64][]*model.MsgInfoModel, userID, conversationID string, msg *model.MsgInfoModel) {
	if msg.IsRead {
		msg.Msg.IsRead = true
	}
	if msg.Msg.ContentType != constant.Quote {
		return
	}
	if msg.Msg.Content == "" {
		return
	}
	var quoteMsg struct {
		Text              string          `json:"text,omitempty"`
		QuoteMessage      *sdkws.MsgData  `json:"quoteMessage,omitempty"`
		MessageEntityList json.RawMessage `json:"messageEntityList,omitempty"`
	}
	if err := json.Unmarshal([]byte(msg.Msg.Content), &quoteMsg); err != nil {
		log.ZError(ctx, "json.Unmarshal", err)
		return
	}
	if quoteMsg.QuoteMessage == nil || quoteMsg.QuoteMessage.ContentType == constant.MsgRevokeNotification {
		return
	}
	var msgs []*model.MsgInfoModel
	if v, ok := cache[quoteMsg.QuoteMessage.Seq]; ok {
		msgs = v
	} else {
		if quoteMsg.QuoteMessage.Seq > 0 {
			ms, err := db.msgDocDatabase.GetMsgBySeqIndexIn1Doc(ctx, userID, db.msgTable.GetDocID(conversationID, quoteMsg.QuoteMessage.Seq), []int64{quoteMsg.QuoteMessage.Seq})
			if err != nil {
				log.ZError(ctx, "GetMsgBySeqIndexIn1Doc", err, "conversationID", conversationID, "seq", quoteMsg.QuoteMessage.Seq)
				return
			}
			msgs = ms
			cache[quoteMsg.QuoteMessage.Seq] = ms
		}
	}
	if len(msgs) != 0 && msgs[0].Msg.ContentType != constant.MsgRevokeNotification {
		return
	}
	quoteMsg.QuoteMessage.ContentType = constant.MsgRevokeNotification
	if len(msgs) > 0 {
		quoteMsg.QuoteMessage.Content = []byte(msgs[0].Msg.Content)
	} else {
		quoteMsg.QuoteMessage.Content = []byte("{}")
	}
	data, err := json.Marshal(&quoteMsg)
	if err != nil {
		log.ZError(ctx, "json.Marshal", err)
		return
	}
	msg.Msg.Content = string(data)
	if _, err := db.msgDocDatabase.UpdateMsg(ctx, db.msgTable.GetDocID(conversationID, msg.Msg.Seq), db.msgTable.GetMsgIndex(msg.Msg.Seq), "msg", msg.Msg); err != nil {
		log.ZError(ctx, "UpdateMsgContent", err)
	}
}

func (db *commonMsgDatabase) findMsgInfoBySeq(ctx context.Context, userID, docID string, conversationID string, seqs []int64) (totalMsgs []*model.MsgInfoModel, err error) {
	msgs, err := db.msgDocDatabase.GetMsgBySeqIndexIn1Doc(ctx, userID, docID, seqs)
	if err != nil {
		return nil, err
	}
	tempCache := make(map[int64][]*model.MsgInfoModel)
	for _, msg := range msgs {
		db.handlerDBMsg(ctx, tempCache, userID, conversationID, msg)
	}
	return msgs, err
}

func (db *commonMsgDatabase) getMsgBySeqsRange(ctx context.Context, userID string, conversationID string, allSeqs []int64, begin, end int64) (seqMsgs []*sdkws.MsgData, err error) {
	log.ZDebug(ctx, "getMsgBySeqsRange", "conversationID", conversationID, "allSeqs", allSeqs, "begin", begin, "end", end)
	for docID, seqs := range db.msgTable.GetDocIDSeqsMap(conversationID, allSeqs) {
		log.ZDebug(ctx, "getMsgBySeqsRange", "docID", docID, "seqs", seqs)
		msgs, err := db.findMsgInfoBySeq(ctx, userID, docID, conversationID, seqs)
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
	userMinSeq, err := db.seq.GetConversationUserMinSeq(ctx, conversationID, userID)
	if err != nil && errs.Unwrap(err) != redis.Nil {
		return 0, 0, nil, err
	}
	minSeq, err := db.seq.GetMinSeq(ctx, conversationID)
	if err != nil && errs.Unwrap(err) != redis.Nil {
		return 0, 0, nil, err
	}
	if userMinSeq > minSeq {
		minSeq = userMinSeq
	}
	// "minSeq" represents the startSeq value that the user can retrieve.
	if minSeq > end {
		log.ZWarn(ctx, "minSeq > end", errs.New("minSeq>end"), "minSeq", minSeq, "end", end)
		return 0, 0, nil, nil
	}
	maxSeq, err := db.seq.GetMaxSeq(ctx, conversationID)
	if err != nil && errs.Unwrap(err) != redis.Nil {
		return 0, 0, nil, err
	}
	log.ZDebug(ctx, "GetMsgBySeqsRange", "userMinSeq", userMinSeq, "conMinSeq", minSeq, "conMaxSeq", maxSeq, "userMaxSeq", userMaxSeq)
	if userMaxSeq != 0 {
		if userMaxSeq < maxSeq {
			maxSeq = userMaxSeq
		}
	}
	// "maxSeq" represents the endSeq value that the user can retrieve.

	if begin < minSeq {
		begin = minSeq
	}
	if end > maxSeq {
		end = maxSeq
	}
	// "begin" and "end" represent the actual startSeq and endSeq values that the user can retrieve.
	if end < begin {
		return 0, 0, nil, errs.ErrArgs.WrapMsg("seq end < begin")
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

	if len(seqs) == 0 {
		return 0, 0, nil, nil
	}
	newBegin := seqs[0]
	newEnd := seqs[len(seqs)-1]
	var successMsgs []*sdkws.MsgData
	log.ZDebug(ctx, "GetMsgBySeqsRange", "first seqs", seqs, "newBegin", newBegin, "newEnd", newEnd)
	cachedMsgs, failedSeqs, err := db.msg.GetMessagesBySeq(ctx, conversationID, seqs)
	if err != nil {
		if err != redis.Nil {

			log.ZError(ctx, "get message from redis exception", err, "conversationID", conversationID, "seqs", seqs)
		}
	}
	successMsgs = append(successMsgs, cachedMsgs...)
	log.ZDebug(ctx, "get msgs from cache", "cachedMsgs", cachedMsgs)
	// get from cache or db

	if len(failedSeqs) > 0 {
		log.ZDebug(ctx, "msgs not exist in redis", "seqs", failedSeqs)
		mongoMsgs, err := db.getMsgBySeqsRange(ctx, userID, conversationID, failedSeqs, begin, end)
		if err != nil {

			return 0, 0, nil, err
		}
		successMsgs = append(mongoMsgs, successMsgs...)
	}

	return minSeq, maxSeq, successMsgs, nil
}

func (db *commonMsgDatabase) GetMsgBySeqs(ctx context.Context, userID string, conversationID string, seqs []int64) (int64, int64, []*sdkws.MsgData, error) {
	userMinSeq, err := db.seq.GetConversationUserMinSeq(ctx, conversationID, userID)
	if err != nil && errs.Unwrap(err) != redis.Nil {
		return 0, 0, nil, err
	}
	minSeq, err := db.seq.GetMinSeq(ctx, conversationID)
	if err != nil && errs.Unwrap(err) != redis.Nil {
		return 0, 0, nil, err
	}
	maxSeq, err := db.seq.GetMaxSeq(ctx, conversationID)
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
	successMsgs, failedSeqs, err := db.msg.GetMessagesBySeq(ctx, conversationID, newSeqs)
	if err != nil {
		if err != redis.Nil {
			log.ZError(ctx, "get message from redis exception", err, "failedSeqs", failedSeqs, "conversationID", conversationID)
		}
	}
	log.ZDebug(ctx, "db.seq.GetMessagesBySeq", "userID", userID, "conversationID", conversationID, "seqs",
		seqs, "len(successMsgs)", len(successMsgs), "failedSeqs", failedSeqs)

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
	log.ZDebug(ctx, "DeleteConversationMsgsAndSetMinSeq", "conversationID", conversationID, "minSeq", minSeq)
	if minSeq == 0 {
		return nil
	}
	return db.seq.SetMinSeq(ctx, conversationID, minSeq)
}

func (db *commonMsgDatabase) UserMsgsDestruct(ctx context.Context, userID string, conversationID string, destructTime int64, lastMsgDestructTime time.Time) (seqs []int64, err error) {
	var index int64
	for {
		// from oldest 2 newest
		msgDocModel, err := db.msgDocDatabase.GetMsgDocModelByIndex(ctx, conversationID, index, 1)
		if err != nil || msgDocModel.DocID == "" {
			if err != nil {
				if err == model.ErrMsgListNotExist {
					log.ZDebug(ctx, "not doc find", "conversationID", conversationID, "userID", userID, "index", index)
				} else {
					log.ZError(ctx, "deleteMsgRecursion GetUserMsgListByIndex failed", err, "conversationID", conversationID, "index", index)
				}
			}
			// If an error is reported, or the error cannot be obtained, it is physically deleted and seq delMongoMsgsPhysical(delStruct.delDocIDList) is returned to end the recursion
			break
		}
		index++
		// && msgDocModel.Msg[0].Msg.SendTime > lastMsgDestructTime.UnixMilli()
		if len(msgDocModel.Msg) > 0 {
			i := 0
			var over bool
			for _, msg := range msgDocModel.Msg {
				i++
				if msg != nil && msg.Msg != nil && msg.Msg.SendTime+destructTime*1000 <= time.Now().UnixMilli() {
					if msg.Msg.SendTime+destructTime*1000 > lastMsgDestructTime.UnixMilli() && !datautil.Contain(userID, msg.DelList...) {
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
		currentUserMinSeq, err := db.seq.GetConversationUserMinSeq(ctx, conversationID, userID)
		if err != nil && errs.Unwrap(err) != redis.Nil {
			return nil, err
		}
		if currentUserMinSeq < userMinSeq {
			if err := db.seq.SetConversationUserMinSeq(ctx, conversationID, userID, userMinSeq); err != nil {
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
// recursion deletes the list and returns the set minimum seq.
func (db *commonMsgDatabase) deleteMsgRecursion(ctx context.Context, conversationID string, index int64, delStruct *delMsgRecursionStruct, remainTime int64) (int64, error) {
	// find from oldest list
	msgDocModel, err := db.msgDocDatabase.GetMsgDocModelByIndex(ctx, conversationID, index, 1)
	if err != nil || msgDocModel.DocID == "" {
		if err != nil {
			if err == model.ErrMsgListNotExist {
				log.ZDebug(ctx, "deleteMsgRecursion ErrMsgListNotExist", "conversationID", conversationID, "index:", index)
			} else {
				log.ZError(ctx, "deleteMsgRecursion GetUserMsgListByIndex failed", err, "conversationID", conversationID, "index", index)
			}
		}
		// If an error is reported, or the error cannot be obtained, it is physically deleted and seq delMongoMsgsPhysical(delStruct.delDocIDList) is returned to end the recursion
		err = db.msgDocDatabase.DeleteDocs(ctx, delStruct.delDocIDs)
		if err != nil {
			return 0, err
		}
		return delStruct.getSetMinSeq() + 1, nil
	}
	log.ZDebug(ctx, "doc info", "conversationID", conversationID, "index", index, "docID", msgDocModel.DocID, "len", len(msgDocModel.Msg))
	if int64(len(msgDocModel.Msg)) > db.msgTable.GetSingleGocMsgNum() {
		log.ZWarn(ctx, "msgs too large", nil, "lenth", len(msgDocModel.Msg), "docID:", msgDocModel.DocID)
	}
	if msgDocModel.IsFull() && msgDocModel.Msg[len(msgDocModel.Msg)-1].Msg.SendTime+(remainTime*1000) < timeutil.GetCurrentTimestampByMill() {
		log.ZDebug(ctx, "doc is full and all msg is expired", "docID", msgDocModel.DocID)
		delStruct.delDocIDs = append(delStruct.delDocIDs, msgDocModel.DocID)
		delStruct.minSeq = msgDocModel.Msg[len(msgDocModel.Msg)-1].Msg.Seq
	} else {
		var delMsgIndexs []int
		for i, MsgInfoModel := range msgDocModel.Msg {
			if MsgInfoModel != nil && MsgInfoModel.Msg != nil {
				if timeutil.GetCurrentTimestampByMill() > MsgInfoModel.Msg.SendTime+(remainTime*1000) {
					delMsgIndexs = append(delMsgIndexs, i)
				}
			}
		}
		if len(delMsgIndexs) > 0 {
			if err = db.msgDocDatabase.DeleteMsgsInOneDocByIndex(ctx, msgDocModel.DocID, delMsgIndexs); err != nil {
				log.ZError(ctx, "deleteMsgRecursion DeleteMsgsInOneDocByIndex failed", err, "conversationID", conversationID, "index", index)
			}
			delStruct.minSeq = int64(msgDocModel.Msg[delMsgIndexs[len(delMsgIndexs)-1]].Msg.Seq)
		}
	}
	seq, err := db.deleteMsgRecursion(ctx, conversationID, index+1, delStruct, remainTime)
	return seq, err
}

func (db *commonMsgDatabase) DeleteMsgsPhysicalBySeqs(ctx context.Context, conversationID string, allSeqs []int64) error {
	if err := db.msg.DeleteMessagesFromCache(ctx, conversationID, allSeqs); err != nil {
		return err
	}
	for docID, seqs := range db.msgTable.GetDocIDSeqsMap(conversationID, allSeqs) {
		var indexes []int
		for _, seq := range seqs {
			indexes = append(indexes, int(db.msgTable.GetMsgIndex(seq)))
		}
		if err := db.msgDocDatabase.DeleteMsgsInOneDocByIndex(ctx, docID, indexes); err != nil {
			return err
		}
	}
	return nil
}

func (db *commonMsgDatabase) DeleteUserMsgsBySeqs(ctx context.Context, userID string, conversationID string, seqs []int64) error {
	if err := db.msg.DeleteMessagesFromCache(ctx, conversationID, seqs); err != nil {
		return err
	}
	for docID, seqs := range db.msgTable.GetDocIDSeqsMap(conversationID, seqs) {
		for _, seq := range seqs {
			if _, err := db.msgDocDatabase.PushUnique(ctx, docID, db.msgTable.GetMsgIndex(seq), "del_list", []string{userID}); err != nil {
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
		maxSeq, err := db.seq.GetMaxSeq(ctx, conversationID)
		if err != nil {
			if err == redis.Nil {
				log.ZDebug(ctx, "max seq is nil", "conversationID", conversationID)
			} else {
				log.ZError(ctx, "get max seq failed", err, "conversationID", conversationID)
			}
			continue
		}
		if err := db.seq.SetMinSeq(ctx, conversationID, maxSeq+1); err != nil {
			log.ZError(ctx, "set min seq failed", err, "conversationID", conversationID, "minSeq", maxSeq+1)
		}
	}
}

func (db *commonMsgDatabase) SetMaxSeq(ctx context.Context, conversationID string, maxSeq int64) error {
	return db.seq.SetMaxSeq(ctx, conversationID, maxSeq)
}

func (db *commonMsgDatabase) GetMaxSeqs(ctx context.Context, conversationIDs []string) (map[string]int64, error) {
	return db.seq.GetMaxSeqs(ctx, conversationIDs)
}

func (db *commonMsgDatabase) GetMaxSeq(ctx context.Context, conversationID string) (int64, error) {
	return db.seq.GetMaxSeq(ctx, conversationID)
}

func (db *commonMsgDatabase) SetMinSeq(ctx context.Context, conversationID string, minSeq int64) error {
	return db.seq.SetMinSeq(ctx, conversationID, minSeq)
}

func (db *commonMsgDatabase) SetMinSeqs(ctx context.Context, seqs map[string]int64) error {
	return db.seq.SetMinSeqs(ctx, seqs)
}

func (db *commonMsgDatabase) GetMinSeqs(ctx context.Context, conversationIDs []string) (map[string]int64, error) {
	return db.seq.GetMinSeqs(ctx, conversationIDs)
}

func (db *commonMsgDatabase) GetMinSeq(ctx context.Context, conversationID string) (int64, error) {
	return db.seq.GetMinSeq(ctx, conversationID)
}

func (db *commonMsgDatabase) GetConversationUserMinSeq(ctx context.Context, conversationID string, userID string) (int64, error) {
	return db.seq.GetConversationUserMinSeq(ctx, conversationID, userID)
}

func (db *commonMsgDatabase) GetConversationUserMinSeqs(ctx context.Context, conversationID string, userIDs []string) (map[string]int64, error) {
	return db.seq.GetConversationUserMinSeqs(ctx, conversationID, userIDs)
}

func (db *commonMsgDatabase) SetConversationUserMinSeq(ctx context.Context, conversationID string, userID string, minSeq int64) error {
	return db.seq.SetConversationUserMinSeq(ctx, conversationID, userID, minSeq)
}

func (db *commonMsgDatabase) SetConversationUserMinSeqs(ctx context.Context, conversationID string, seqs map[string]int64) (err error) {
	return db.seq.SetConversationUserMinSeqs(ctx, conversationID, seqs)
}

func (db *commonMsgDatabase) SetUserConversationsMinSeqs(ctx context.Context, userID string, seqs map[string]int64) error {
	return db.seq.SetUserConversationsMinSeqs(ctx, userID, seqs)
}

func (db *commonMsgDatabase) UserSetHasReadSeqs(ctx context.Context, userID string, hasReadSeqs map[string]int64) error {
	return db.seq.UserSetHasReadSeqs(ctx, userID, hasReadSeqs)
}

func (db *commonMsgDatabase) SetHasReadSeq(ctx context.Context, userID string, conversationID string, hasReadSeq int64) error {
	return db.seq.SetHasReadSeq(ctx, userID, conversationID, hasReadSeq)
}

func (db *commonMsgDatabase) GetHasReadSeqs(ctx context.Context, userID string, conversationIDs []string) (map[string]int64, error) {
	return db.seq.GetHasReadSeqs(ctx, userID, conversationIDs)
}

func (db *commonMsgDatabase) GetHasReadSeq(ctx context.Context, userID string, conversationID string) (int64, error) {
	return db.seq.GetHasReadSeq(ctx, userID, conversationID)
}

func (db *commonMsgDatabase) SetSendMsgStatus(ctx context.Context, id string, status int32) error {
	return db.msg.SetSendMsgStatus(ctx, id, status)
}

func (db *commonMsgDatabase) GetSendMsgStatus(ctx context.Context, id string) (int32, error) {
	return db.msg.GetSendMsgStatus(ctx, id)
}

func (db *commonMsgDatabase) GetConversationMinMaxSeqInMongoAndCache(ctx context.Context, conversationID string) (minSeqMongo, maxSeqMongo, minSeqCache, maxSeqCache int64, err error) {
	minSeqMongo, maxSeqMongo, err = db.GetMinMaxSeqMongo(ctx, conversationID)
	if err != nil {
		return
	}
	minSeqCache, err = db.seq.GetMinSeq(ctx, conversationID)
	if err != nil {
		return
	}
	maxSeqCache, err = db.seq.GetMaxSeq(ctx, conversationID)
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
) (msgCount int64, userCount int64, users []*model.UserCount, dateCount map[string]int64, err error) {
	return db.msgDocDatabase.RangeUserSendCount(ctx, start, end, group, ase, pageNumber, showNumber)
}

func (db *commonMsgDatabase) RangeGroupSendCount(
	ctx context.Context,
	start time.Time,
	end time.Time,
	ase bool,
	pageNumber int32,
	showNumber int32,
) (msgCount int64, userCount int64, groups []*model.GroupCount, dateCount map[string]int64, err error) {
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

func (db *commonMsgDatabase) FindOneByDocIDs(ctx context.Context, conversationIDs []string, seqs map[string]int64) (map[string]*sdkws.MsgData, error) {
	totalMsgs := make(map[string]*sdkws.MsgData)
	for _, conversationID := range conversationIDs {
		seq := seqs[conversationID]
		docID := db.msgTable.GetDocID(conversationID, seq)
		msgs, err := db.msgDocDatabase.FindOneByDocID(ctx, docID)
		if err != nil {
			return nil, err
		}
		index := db.msgTable.GetMsgIndex(seq)
		totalMsgs[conversationID] = convert.MsgDB2Pb(msgs.Msg[index].Msg)
	}
	return totalMsgs, nil
}

func (db *commonMsgDatabase) ConvertMsgsDocLen(ctx context.Context, conversationIDs []string) {
	db.msgDocDatabase.ConvertMsgsDocLen(ctx, conversationIDs)
}

func (db *commonMsgDatabase) GetBeforeMsg(ctx context.Context, ts int64, limit int) ([]*model.MsgDocModel, error) {
	return db.msgDocDatabase.GetBeforeMsg(ctx, ts, limit)
}

func (db *commonMsgDatabase) DeleteDocMsgBefore(ctx context.Context, ts int64, doc *model.MsgDocModel) ([]int, error) {
	var notNull int
	index := make([]int, 0, len(doc.Msg))
	for i, message := range doc.Msg {
		if message.Msg != nil {
			notNull++
			if message.Msg.SendTime < ts {
				index = append(index, i)
			}
		}
	}
	if len(index) == 0 {
		return index, nil
	}
	maxSeq := doc.Msg[index[len(index)-1]].Msg.Seq
	conversationID := doc.DocID[:strings.LastIndex(doc.DocID, ":")]
	if err := db.setMinSeq(ctx, conversationID, maxSeq+1); err != nil {
		return index, err
	}
	if len(index) == notNull {
		return index, db.msgDocDatabase.DeleteDoc(ctx, doc.DocID)
	} else {
		return index, db.msgDocDatabase.DeleteMsgByIndex(ctx, doc.DocID, index)
	}
}

//func (db *commonMsgDatabase) ClearMsg(ctx context.Context, ts int64) (err error) {
//	var (
//		docNum int
//		msgNum int
//		start  = time.Now()
//	)
//	for {
//		msgs, err := db.msgDocDatabase.GetBeforeMsg(ctx, ts, 100)
//		if err != nil {
//			return err
//		}
//		if len(msgs) == 0 {
//			return nil
//		}
//		for _, msg := range msgs {
//			num, err := db.deleteOneMsg(ctx, ts, msg)
//			if err != nil {
//				return err
//			}
//			docNum++
//			msgNum += num
//		}
//	}
//}

func (db *commonMsgDatabase) setMinSeq(ctx context.Context, conversationID string, seq int64) error {
	dbSeq, err := db.seq.GetMinSeq(ctx, conversationID)
	if err != nil {
		if errors.Is(errs.Unwrap(err), redis.Nil) {
			return nil
		}
		return err
	}
	if dbSeq >= seq {
		return nil
	}
	return db.seq.SetMinSeq(ctx, conversationID, seq)
}
