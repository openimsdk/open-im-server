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
	"strconv"
	"strings"
	"time"

	"github.com/openimsdk/tools/utils/jsonutil"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/convert"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/protocol/constant"
	pbmsg "github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mq/kafka"
	"github.com/openimsdk/tools/utils/datautil"
)

const (
	updateKeyMsg = iota
	updateKeyRevoke
)

// CommonMsgDatabase defines the interface for message database operations.
type CommonMsgDatabase interface {
	// RevokeMsg revokes a message in a conversation.
	RevokeMsg(ctx context.Context, conversationID string, seq int64, revoke *model.RevokeModel) error
	// MarkSingleChatMsgsAsRead marks messages as read for a single chat by sequence numbers.
	MarkSingleChatMsgsAsRead(ctx context.Context, userID string, conversationID string, seqs []int64) error
	// GetMsgBySeqsRange retrieves messages from MongoDB by a range of sequence numbers.
	GetMsgBySeqsRange(ctx context.Context, userID string, conversationID string, begin, end, num, userMaxSeq int64) (minSeq int64, maxSeq int64, seqMsg []*sdkws.MsgData, err error)
	// GetMsgBySeqs retrieves messages for large groups from MongoDB by sequence numbers.
	GetMsgBySeqs(ctx context.Context, userID string, conversationID string, seqs []int64) (minSeq int64, maxSeq int64, seqMsg []*sdkws.MsgData, err error)

	GetMessagesBySeqWithBounds(ctx context.Context, userID string, conversationID string, seqs []int64, pullOrder sdkws.PullOrder) (bool, int64, []*sdkws.MsgData, error)
	// DeleteUserMsgsBySeqs allows a user to delete messages based on sequence numbers.
	DeleteUserMsgsBySeqs(ctx context.Context, userID string, conversationID string, seqs []int64) error
	// DeleteMsgsPhysicalBySeqs physically deletes messages by emptying them based on sequence numbers.
	DeleteMsgsPhysicalBySeqs(ctx context.Context, conversationID string, seqs []int64) error
	//SetMaxSeq(ctx context.Context, conversationID string, maxSeq int64) error
	GetMaxSeqs(ctx context.Context, conversationIDs []string) (map[string]int64, error)
	GetMaxSeq(ctx context.Context, conversationID string) (int64, error)
	SetMinSeqs(ctx context.Context, seqs map[string]int64) error
	SetMinSeq(ctx context.Context, conversationID string, seq int64) error

	SetUserConversationsMinSeqs(ctx context.Context, userID string, seqs map[string]int64) (err error)
	SetHasReadSeq(ctx context.Context, userID string, conversationID string, hasReadSeq int64) error
	GetHasReadSeqs(ctx context.Context, userID string, conversationIDs []string) (map[string]int64, error)
	GetHasReadSeq(ctx context.Context, userID string, conversationID string) (int64, error)
	UserSetHasReadSeqs(ctx context.Context, userID string, hasReadSeqs map[string]int64) error

	GetMaxSeqsWithTime(ctx context.Context, conversationIDs []string) (map[string]database.SeqTime, error)
	GetMaxSeqWithTime(ctx context.Context, conversationID string) (database.SeqTime, error)
	GetCacheMaxSeqWithTime(ctx context.Context, conversationIDs []string) (map[string]database.SeqTime, error)

	SetSendMsgStatus(ctx context.Context, id string, status int32) error
	GetSendMsgStatus(ctx context.Context, id string) (int32, error)
	SearchMessage(ctx context.Context, req *pbmsg.SearchMessageReq) (total int64, msgData []*pbmsg.SearchedMsgData, err error)
	FindOneByDocIDs(ctx context.Context, docIDs []string, seqs map[string]int64) (map[string]*sdkws.MsgData, error)

	// to mq
	MsgToMQ(ctx context.Context, key string, msg2mq *sdkws.MsgData) error

	RangeUserSendCount(ctx context.Context, start time.Time, end time.Time, group bool, ase bool, pageNumber int32, showNumber int32) (msgCount int64, userCount int64, users []*model.UserCount, dateCount map[string]int64, err error)
	RangeGroupSendCount(ctx context.Context, start time.Time, end time.Time, ase bool, pageNumber int32, showNumber int32) (msgCount int64, userCount int64, groups []*model.GroupCount, dateCount map[string]int64, err error)

	GetRandBeforeMsg(ctx context.Context, ts int64, limit int) ([]*model.MsgDocModel, error)

	SetUserConversationsMaxSeq(ctx context.Context, conversationID string, userID string, seq int64) error
	SetUserConversationsMinSeq(ctx context.Context, conversationID string, userID string, seq int64) error

	DeleteDoc(ctx context.Context, docID string) error

	GetLastMessageSeqByTime(ctx context.Context, conversationID string, time int64) (int64, error)

	GetLastMessage(ctx context.Context, conversationIDS []string, userID string) (map[string]*sdkws.MsgData, error)
}

func NewCommonMsgDatabase(msgDocModel database.Msg, msg cache.MsgCache, seqUser cache.SeqUser, seqConversation cache.SeqConversationCache, kafkaConf *config.Kafka) (CommonMsgDatabase, error) {
	conf, err := kafka.BuildProducerConfig(*kafkaConf.Build())
	if err != nil {
		return nil, err
	}
	producerToRedis, err := kafka.NewKafkaProducer(conf, kafkaConf.Address, kafkaConf.ToRedisTopic)
	if err != nil {
		return nil, err
	}
	return &commonMsgDatabase{
		msgDocDatabase:  msgDocModel,
		msgCache:        msg,
		seqUser:         seqUser,
		seqConversation: seqConversation,
		producer:        producerToRedis,
	}, nil
}

type commonMsgDatabase struct {
	msgDocDatabase  database.Msg
	msgTable        model.MsgDocModel
	msgCache        cache.MsgCache
	seqConversation cache.SeqConversationCache
	seqUser         cache.SeqUser
	producer        *kafka.Producer
}

func (db *commonMsgDatabase) MsgToMQ(ctx context.Context, key string, msg2mq *sdkws.MsgData) error {
	_, _, err := db.producer.SendMessage(ctx, key, msg2mq)
	return err
}

func (db *commonMsgDatabase) batchInsertBlock(ctx context.Context, conversationID string, fields []any, key int8, firstSeq int64) error {
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

func (db *commonMsgDatabase) RevokeMsg(ctx context.Context, conversationID string, seq int64, revoke *model.RevokeModel) error {
	if err := db.batchInsertBlock(ctx, conversationID, []any{revoke}, updateKeyRevoke, seq); err != nil {
		return err
	}
	return db.msgCache.DelMessageBySeqs(ctx, conversationID, []int64{seq})
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
	return db.msgCache.DelMessageBySeqs(ctx, conversationID, totalSeqs)
}

func (db *commonMsgDatabase) getMsgBySeqs(ctx context.Context, userID, conversationID string, seqs []int64) (totalMsgs []*sdkws.MsgData, err error) {
	return db.GetMessageBySeqs(ctx, conversationID, userID, seqs)
}

func (db *commonMsgDatabase) handlerDBMsg(ctx context.Context, cache map[int64][]*model.MsgInfoModel, userID, conversationID string, msg *model.MsgInfoModel) {
	if msg == nil || msg.Msg == nil {
		return
	}
	if msg.IsRead {
		msg.Msg.IsRead = true
	}
	if msg.Msg.ContentType != constant.Quote {
		return
	}
	if msg.Msg.Content == "" {
		return
	}
	type MsgData struct {
		SendID           string                 `json:"sendID"`
		RecvID           string                 `json:"recvID"`
		GroupID          string                 `json:"groupID"`
		ClientMsgID      string                 `json:"clientMsgID"`
		ServerMsgID      string                 `json:"serverMsgID"`
		SenderPlatformID int32                  `json:"senderPlatformID"`
		SenderNickname   string                 `json:"senderNickname"`
		SenderFaceURL    string                 `json:"senderFaceURL"`
		SessionType      int32                  `json:"sessionType"`
		MsgFrom          int32                  `json:"msgFrom"`
		ContentType      int32                  `json:"contentType"`
		Content          string                 `json:"content"`
		Seq              int64                  `json:"seq"`
		SendTime         int64                  `json:"sendTime"`
		CreateTime       int64                  `json:"createTime"`
		Status           int32                  `json:"status"`
		IsRead           bool                   `json:"isRead"`
		Options          map[string]bool        `json:"options,omitempty"`
		OfflinePushInfo  *sdkws.OfflinePushInfo `json:"offlinePushInfo"`
		AtUserIDList     []string               `json:"atUserIDList"`
		AttachedInfo     string                 `json:"attachedInfo"`
		Ex               string                 `json:"ex"`
		KeyVersion       int32                  `json:"keyVersion"`
		DstUserIDs       []string               `json:"dstUserIDs"`
	}
	var quoteMsg struct {
		Text              string          `json:"text,omitempty"`
		QuoteMessage      *MsgData        `json:"quoteMessage,omitempty"`
		MessageEntityList json.RawMessage `json:"messageEntityList,omitempty"`
	}
	if err := json.Unmarshal([]byte(msg.Msg.Content), &quoteMsg); err != nil {
		log.ZError(ctx, "json.Unmarshal", err)
		return
	}
	if quoteMsg.QuoteMessage == nil || quoteMsg.QuoteMessage.Content == "" {
		return
	}
	if quoteMsg.QuoteMessage.Content == "e30=" {
		quoteMsg.QuoteMessage.Content = "{}"
		data, err := json.Marshal(&quoteMsg)
		if err != nil {
			return
		}
		msg.Msg.Content = string(data)
	}
	if quoteMsg.QuoteMessage.Seq <= 0 && quoteMsg.QuoteMessage.ContentType == constant.MsgRevokeNotification {
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
		quoteMsg.QuoteMessage.Content = msgs[0].Msg.Content
	} else {
		quoteMsg.QuoteMessage.Content = "{}"
	}
	data, err := json.Marshal(&quoteMsg)
	if err != nil {
		log.ZError(ctx, "json.Marshal", err)
		return
	}
	msg.Msg.Content = string(data)
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
	userMinSeq, err := db.seqUser.GetUserMinSeq(ctx, conversationID, userID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, 0, nil, err
	}
	minSeq, err := db.seqConversation.GetMinSeq(ctx, conversationID)
	if err != nil {
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
	maxSeq, err := db.seqConversation.GetMaxSeq(ctx, conversationID)
	if err != nil {
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
	successMsgs, err := db.GetMessageBySeqs(ctx, conversationID, userID, seqs)
	if err != nil {
		return 0, 0, nil, err
	}
	return minSeq, maxSeq, successMsgs, nil
}

func (db *commonMsgDatabase) GetMsgBySeqs(ctx context.Context, userID string, conversationID string, seqs []int64) (int64, int64, []*sdkws.MsgData, error) {
	userMinSeq, err := db.seqUser.GetUserMinSeq(ctx, conversationID, userID)
	if err != nil {
		return 0, 0, nil, err
	}
	minSeq, err := db.seqConversation.GetMinSeq(ctx, conversationID)
	if err != nil {
		return 0, 0, nil, err
	}
	maxSeq, err := db.seqConversation.GetMaxSeq(ctx, conversationID)
	if err != nil {
		return 0, 0, nil, err
	}
	userMaxSeq, err := db.seqUser.GetUserMaxSeq(ctx, conversationID, userID)
	if err != nil {
		return 0, 0, nil, err
	}
	if userMinSeq > minSeq {
		minSeq = userMinSeq
	}
	if userMaxSeq > 0 && userMaxSeq < maxSeq {
		maxSeq = userMaxSeq
	}
	newSeqs := make([]int64, 0, len(seqs))
	for _, seq := range seqs {
		if seq <= 0 {
			continue
		}
		if seq >= minSeq && seq <= maxSeq {
			newSeqs = append(newSeqs, seq)
		}
	}
	successMsgs, err := db.GetMessageBySeqs(ctx, conversationID, userID, newSeqs)
	if err != nil {
		return 0, 0, nil, err
	}
	return minSeq, maxSeq, successMsgs, nil
}

func (db *commonMsgDatabase) GetMessagesBySeqWithBounds(ctx context.Context, userID string, conversationID string, seqs []int64, pullOrder sdkws.PullOrder) (bool, int64, []*sdkws.MsgData, error) {
	var endSeq int64
	var isEnd bool
	userMinSeq, err := db.seqUser.GetUserMinSeq(ctx, conversationID, userID)
	if err != nil {
		return false, 0, nil, err
	}
	minSeq, err := db.seqConversation.GetMinSeq(ctx, conversationID)
	if err != nil {
		return false, 0, nil, err
	}
	maxSeq, err := db.seqConversation.GetMaxSeq(ctx, conversationID)
	if err != nil {
		return false, 0, nil, err
	}
	userMaxSeq, err := db.seqUser.GetUserMaxSeq(ctx, conversationID, userID)
	if err != nil {
		return false, 0, nil, err
	}
	if userMinSeq > minSeq {
		minSeq = userMinSeq
	}
	if userMaxSeq > 0 && userMaxSeq < maxSeq {
		maxSeq = userMaxSeq
	}
	newSeqs := make([]int64, 0, len(seqs))
	for _, seq := range seqs {
		if seq <= 0 {
			continue
		}
		// The normal range and can fetch messages
		if seq >= minSeq && seq <= maxSeq {
			newSeqs = append(newSeqs, seq)
			continue
		}
		// If the requested seq is smaller than the minimum seq and the pull order is descending (pulling older messages)
		if seq < minSeq && pullOrder == sdkws.PullOrder_PullOrderDesc {
			isEnd = true
			endSeq = minSeq
		}
		// If the requested seq is larger than the maximum seq and the pull order is ascending (pulling newer messages)
		if seq > maxSeq && pullOrder == sdkws.PullOrder_PullOrderAsc {
			isEnd = true
			endSeq = maxSeq
		}
	}
	if len(newSeqs) == 0 {
		return isEnd, endSeq, nil, nil
	}
	successMsgs, err := db.GetMessageBySeqs(ctx, conversationID, userID, newSeqs)
	if err != nil {
		return false, 0, nil, err
	}
	return isEnd, endSeq, successMsgs, nil
}

func (db *commonMsgDatabase) DeleteMsgsPhysicalBySeqs(ctx context.Context, conversationID string, allSeqs []int64) error {
	for docID, seqs := range db.msgTable.GetDocIDSeqsMap(conversationID, allSeqs) {
		var indexes []int
		for _, seq := range seqs {
			indexes = append(indexes, int(db.msgTable.GetMsgIndex(seq)))
		}
		if err := db.msgDocDatabase.DeleteMsgsInOneDocByIndex(ctx, docID, indexes); err != nil {
			return err
		}
	}
	return db.msgCache.DelMessageBySeqs(ctx, conversationID, allSeqs)
}

func (db *commonMsgDatabase) DeleteUserMsgsBySeqs(ctx context.Context, userID string, conversationID string, seqs []int64) error {
	for docID, seqs := range db.msgTable.GetDocIDSeqsMap(conversationID, seqs) {
		for _, seq := range seqs {
			if _, err := db.msgDocDatabase.PushUnique(ctx, docID, db.msgTable.GetMsgIndex(seq), "del_list", []string{userID}); err != nil {
				return err
			}
		}
	}
	return db.msgCache.DelMessageBySeqs(ctx, conversationID, seqs)
}

func (db *commonMsgDatabase) GetMaxSeqs(ctx context.Context, conversationIDs []string) (map[string]int64, error) {
	return db.seqConversation.GetMaxSeqs(ctx, conversationIDs)
}

func (db *commonMsgDatabase) GetMaxSeq(ctx context.Context, conversationID string) (int64, error) {
	return db.seqConversation.GetMaxSeq(ctx, conversationID)
}

func (db *commonMsgDatabase) SetMinSeqs(ctx context.Context, seqs map[string]int64) error {
	return db.seqConversation.SetMinSeqs(ctx, seqs)
}

func (db *commonMsgDatabase) SetUserConversationsMinSeqs(ctx context.Context, userID string, seqs map[string]int64) error {
	return db.seqUser.SetUserMinSeqs(ctx, userID, seqs)
}

func (db *commonMsgDatabase) SetUserConversationsMaxSeq(ctx context.Context, conversationID string, userID string, seq int64) error {
	return db.seqUser.SetUserMaxSeq(ctx, conversationID, userID, seq)
}

func (db *commonMsgDatabase) SetUserConversationsMinSeq(ctx context.Context, conversationID string, userID string, seq int64) error {
	return db.seqUser.SetUserMinSeq(ctx, conversationID, userID, seq)
}

func (db *commonMsgDatabase) UserSetHasReadSeqs(ctx context.Context, userID string, hasReadSeqs map[string]int64) error {
	return db.seqUser.SetUserReadSeqs(ctx, userID, hasReadSeqs)
}

func (db *commonMsgDatabase) SetHasReadSeq(ctx context.Context, userID string, conversationID string, hasReadSeq int64) error {
	return db.seqUser.SetUserReadSeq(ctx, conversationID, userID, hasReadSeq)
}

func (db *commonMsgDatabase) GetHasReadSeqs(ctx context.Context, userID string, conversationIDs []string) (map[string]int64, error) {
	return db.seqUser.GetUserReadSeqs(ctx, userID, conversationIDs)
}

func (db *commonMsgDatabase) GetHasReadSeq(ctx context.Context, userID string, conversationID string) (int64, error) {
	return db.seqUser.GetUserReadSeq(ctx, conversationID, userID)
}

func (db *commonMsgDatabase) SetSendMsgStatus(ctx context.Context, id string, status int32) error {
	return db.msgCache.SetSendMsgStatus(ctx, id, status)
}

func (db *commonMsgDatabase) GetSendMsgStatus(ctx context.Context, id string) (int32, error) {
	return db.msgCache.GetSendMsgStatus(ctx, id)
}

func (db *commonMsgDatabase) GetConversationMinMaxSeqInMongoAndCache(ctx context.Context, conversationID string) (minSeqMongo, maxSeqMongo, minSeqCache, maxSeqCache int64, err error) {
	minSeqMongo, maxSeqMongo, err = db.GetMinMaxSeqMongo(ctx, conversationID)
	if err != nil {
		return
	}
	minSeqCache, err = db.seqConversation.GetMinSeq(ctx, conversationID)
	if err != nil {
		return
	}
	maxSeqCache, err = db.seqConversation.GetMaxSeq(ctx, conversationID)
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

func (db *commonMsgDatabase) RangeUserSendCount(ctx context.Context, start time.Time, end time.Time, group bool, ase bool, pageNumber int32, showNumber int32) (msgCount int64, userCount int64, users []*model.UserCount, dateCount map[string]int64, err error) {
	return db.msgDocDatabase.RangeUserSendCount(ctx, start, end, group, ase, pageNumber, showNumber)
}

func (db *commonMsgDatabase) RangeGroupSendCount(ctx context.Context, start time.Time, end time.Time, ase bool, pageNumber int32, showNumber int32) (msgCount int64, userCount int64, groups []*model.GroupCount, dateCount map[string]int64, err error) {
	return db.msgDocDatabase.RangeGroupSendCount(ctx, start, end, ase, pageNumber, showNumber)
}

func (db *commonMsgDatabase) SearchMessage(ctx context.Context, req *pbmsg.SearchMessageReq) (total int64, msgData []*pbmsg.SearchedMsgData, err error) {
	var totalMsgs []*pbmsg.SearchedMsgData
	total, msgs, err := db.msgDocDatabase.SearchMessage(ctx, req)
	if err != nil {
		return 0, nil, err
	}
	for _, msg := range msgs {
		if msg.IsRead {
			msg.Msg.IsRead = true
		}
		searchedMsgData := &pbmsg.SearchedMsgData{MsgData: convert.MsgDB2Pb(msg.Msg)}

		if msg.Revoke != nil {
			searchedMsgData.IsRevoked = true
		}

		totalMsgs = append(totalMsgs, searchedMsgData)
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

func (db *commonMsgDatabase) GetRandBeforeMsg(ctx context.Context, ts int64, limit int) ([]*model.MsgDocModel, error) {
	return db.msgDocDatabase.GetRandBeforeMsg(ctx, ts, limit)
}

func (db *commonMsgDatabase) SetMinSeq(ctx context.Context, conversationID string, seq int64) error {
	dbSeq, err := db.seqConversation.GetMinSeq(ctx, conversationID)
	if err != nil {
		if errors.Is(errs.Unwrap(err), redis.Nil) {
			return nil
		}
		return err
	}
	if dbSeq >= seq {
		return nil
	}
	return db.seqConversation.SetMinSeq(ctx, conversationID, seq)
}

func (db *commonMsgDatabase) GetCacheMaxSeqWithTime(ctx context.Context, conversationIDs []string) (map[string]database.SeqTime, error) {
	return db.seqConversation.GetCacheMaxSeqWithTime(ctx, conversationIDs)
}

func (db *commonMsgDatabase) GetMaxSeqWithTime(ctx context.Context, conversationID string) (database.SeqTime, error) {
	return db.seqConversation.GetMaxSeqWithTime(ctx, conversationID)
}

func (db *commonMsgDatabase) GetMaxSeqsWithTime(ctx context.Context, conversationIDs []string) (map[string]database.SeqTime, error) {
	// todo: only the time in the redis cache will be taken, not the message time
	return db.seqConversation.GetMaxSeqsWithTime(ctx, conversationIDs)
}

func (db *commonMsgDatabase) DeleteDoc(ctx context.Context, docID string) error {
	index := strings.LastIndex(docID, ":")
	if index <= 0 {
		return errs.ErrInternalServer.WrapMsg("docID is invalid", "docID", docID)
	}
	docIndex, err := strconv.Atoi(docID[index+1:])
	if err != nil {
		return errs.WrapMsg(err, "strconv.Atoi", "docID", docID)
	}
	conversationID := docID[:index]
	seqs := make([]int64, db.msgTable.GetSingleGocMsgNum())
	minSeq := db.msgTable.GetMinSeq(docIndex)
	for i := range seqs {
		seqs[i] = minSeq + int64(i)
	}
	if err := db.msgDocDatabase.DeleteDoc(ctx, docID); err != nil {
		return err
	}
	return db.msgCache.DelMessageBySeqs(ctx, conversationID, seqs)
}

func (db *commonMsgDatabase) GetLastMessageSeqByTime(ctx context.Context, conversationID string, time int64) (int64, error) {
	return db.msgDocDatabase.GetLastMessageSeqByTime(ctx, conversationID, time)
}

func (db *commonMsgDatabase) handlerDeleteAndRevoked(ctx context.Context, userID string, msgs []*model.MsgInfoModel) {
	for i := range msgs {
		msg := msgs[i]
		if msg == nil || msg.Msg == nil {
			continue
		}
		msg.Msg.IsRead = msg.IsRead
		if datautil.Contain(userID, msg.DelList...) {
			msg.Msg.Content = ""
			msg.Msg.Status = constant.MsgDeleted
		}
		if msg.Revoke == nil {
			continue
		}
		msg.Msg.ContentType = constant.MsgRevokeNotification
		revokeContent := sdkws.MessageRevokedContent{
			RevokerID:                   msg.Revoke.UserID,
			RevokerRole:                 msg.Revoke.Role,
			ClientMsgID:                 msg.Msg.ClientMsgID,
			RevokerNickname:             msg.Revoke.Nickname,
			RevokeTime:                  msg.Revoke.Time,
			SourceMessageSendTime:       msg.Msg.SendTime,
			SourceMessageSendID:         msg.Msg.SendID,
			SourceMessageSenderNickname: msg.Msg.SenderNickname,
			SessionType:                 msg.Msg.SessionType,
			Seq:                         msg.Msg.Seq,
			Ex:                          msg.Msg.Ex,
		}
		data, err := jsonutil.JsonMarshal(&revokeContent)
		if err != nil {
			log.ZWarn(ctx, "handlerDeleteAndRevoked JsonMarshal MessageRevokedContent", err, "msg", msg)
			continue
		}
		elem := sdkws.NotificationElem{
			Detail: string(data),
		}
		content, err := jsonutil.JsonMarshal(&elem)
		if err != nil {
			log.ZWarn(ctx, "handlerDeleteAndRevoked JsonMarshal NotificationElem", err, "msg", msg)
			continue
		}
		msg.Msg.Content = string(content)
	}
}

func (db *commonMsgDatabase) handlerQuote(ctx context.Context, userID, conversationID string, msgs []*model.MsgInfoModel) {
	temp := make(map[int64][]*model.MsgInfoModel)
	for i := range msgs {
		db.handlerDBMsg(ctx, temp, userID, conversationID, msgs[i])
	}
}

func (db *commonMsgDatabase) GetMessageBySeqs(ctx context.Context, conversationID string, userID string, seqs []int64) ([]*sdkws.MsgData, error) {
	msgs, err := db.msgCache.GetMessageBySeqs(ctx, conversationID, seqs)
	if err != nil {
		return nil, err
	}
	db.handlerDeleteAndRevoked(ctx, userID, msgs)
	db.handlerQuote(ctx, userID, conversationID, msgs)
	seqMsgs := make(map[int64]*model.MsgInfoModel)
	for i, msg := range msgs {
		if msg.Msg == nil {
			continue
		}
		seqMsgs[msg.Msg.Seq] = msgs[i]
	}
	res := make([]*sdkws.MsgData, 0, len(seqs))
	for _, seq := range seqs {
		if v, ok := seqMsgs[seq]; ok {
			res = append(res, convert.MsgDB2Pb(v.Msg))
		} else {
			res = append(res, &sdkws.MsgData{Seq: seq, Status: constant.MsgStatusHasDeleted})
		}
	}
	return res, nil
}

func (db *commonMsgDatabase) GetLastMessage(ctx context.Context, conversationIDs []string, userID string) (map[string]*sdkws.MsgData, error) {
	res := make(map[string]*sdkws.MsgData)
	for _, conversationID := range conversationIDs {
		if _, ok := res[conversationID]; ok {
			continue
		}
		msg, err := db.msgDocDatabase.GetLastMessage(ctx, conversationID)
		if err != nil {
			if errs.Unwrap(err) == mongo.ErrNoDocuments {
				continue
			}
			return nil, err
		}
		tmp := []*model.MsgInfoModel{msg}
		db.handlerDeleteAndRevoked(ctx, userID, tmp)
		db.handlerQuote(ctx, userID, conversationID, tmp)
		res[conversationID] = convert.MsgDB2Pb(msg.Msg)
	}
	return res, nil
}
