package controller

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/convert"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/kafka"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/utils/datautil"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	pbmsg "github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"go.mongodb.org/mongo-driver/mongo"
)

type MsgTransferDatabase interface {
	// BatchInsertChat2DB inserts a batch of messages into the database for a specific conversation.
	BatchInsertChat2DB(ctx context.Context, conversationID string, msgs []*sdkws.MsgData, currentMaxSeq int64) error
	// DeleteMessagesFromCache deletes message caches from Redis by sequence numbers.
	DeleteMessagesFromCache(ctx context.Context, conversationID string, seqs []int64) error

	// BatchInsertChat2Cache increments the sequence number and then batch inserts messages into the cache.
	BatchInsertChat2Cache(ctx context.Context, conversationID string, msgs []*sdkws.MsgData) (seq int64, isNewConversation bool, userHasReadMap map[string]int64, err error)

	SetHasReadSeqs(ctx context.Context, conversationID string, userSeqMap map[string]int64) error

	SetHasReadSeqToDB(ctx context.Context, conversationID string, userSeqMap map[string]int64) error

	// to mq
	MsgToPushMQ(ctx context.Context, key, conversationID string, msg2mq *sdkws.MsgData) (int32, int64, error)
	MsgToMongoMQ(ctx context.Context, key, conversationID string, msgs []*sdkws.MsgData, lastSeq int64) error
}

func NewMsgTransferDatabase(msgDocModel database.Msg, msg cache.MsgCache, seqUser cache.SeqUser, seqConversation cache.SeqConversationCache, kafkaConf *config.Kafka) (MsgTransferDatabase, error) {
	conf, err := kafka.BuildProducerConfig(*kafkaConf.Build())
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
	return &msgTransferDatabase{
		msgDocDatabase:  msgDocModel,
		msgCache:        msg,
		seqUser:         seqUser,
		seqConversation: seqConversation,
		producerToMongo: producerToMongo,
		producerToPush:  producerToPush,
	}, nil
}

type msgTransferDatabase struct {
	msgDocDatabase  database.Msg
	msgTable        model.MsgDocModel
	msgCache        cache.MsgCache
	seqConversation cache.SeqConversationCache
	seqUser         cache.SeqUser
	producerToMongo *kafka.Producer
	producerToPush  *kafka.Producer
}

func (db *msgTransferDatabase) BatchInsertChat2DB(ctx context.Context, conversationID string, msgList []*sdkws.MsgData, currentMaxSeq int64) error {
	if len(msgList) == 0 {
		return errs.ErrArgs.WrapMsg("msgList is empty")
	}
	msgs := make([]any, len(msgList))
	seqs := make([]int64, len(msgList))
	for i, msg := range msgList {
		if msg == nil {
			continue
		}
		seqs[i] = msg.Seq
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
		if msg.Status == constant.MsgStatusSending {
			msg.Status = constant.MsgStatusSendSuccess
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
	if err := db.BatchInsertBlock(ctx, conversationID, msgs, updateKeyMsg, msgList[0].Seq); err != nil {
		return err
	}
	//return db.msgCache.DelMessageBySeqs(ctx, conversationID, seqs)
	return nil
}

func (db *msgTransferDatabase) BatchInsertBlock(ctx context.Context, conversationID string, fields []any, key int8, firstSeq int64) error {
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

func (db *msgTransferDatabase) DeleteMessagesFromCache(ctx context.Context, conversationID string, seqs []int64) error {
	return db.msgCache.DelMessageBySeqs(ctx, conversationID, seqs)
}

func (db *msgTransferDatabase) BatchInsertChat2Cache(ctx context.Context, conversationID string, msgs []*sdkws.MsgData) (seq int64, isNew bool, userHasReadMap map[string]int64, err error) {
	lenList := len(msgs)
	if int64(lenList) > db.msgTable.GetSingleGocMsgNum() {
		return 0, false, nil, errs.New("message count exceeds limit", "limit", db.msgTable.GetSingleGocMsgNum()).Wrap()
	}
	if lenList < 1 {
		return 0, false, nil, errs.New("no messages to insert", "minCount", 1).Wrap()
	}
	currentMaxSeq, err := db.seqConversation.Malloc(ctx, conversationID, int64(len(msgs)))
	if err != nil {
		log.ZError(ctx, "storage.seq.Malloc", err)
		return 0, false, nil, err
	}
	isNew = currentMaxSeq == 0
	lastMaxSeq := currentMaxSeq
	userSeqMap := make(map[string]int64)
	seqs := make([]int64, 0, lenList)
	for _, m := range msgs {
		currentMaxSeq++
		m.Seq = currentMaxSeq
		userSeqMap[m.SendID] = m.Seq
		seqs = append(seqs, m.Seq)
	}
	msgToDB := func(msg *sdkws.MsgData) *model.MsgInfoModel {
		return &model.MsgInfoModel{
			Msg: convert.MsgPb2DB(msg),
		}
	}
	if err := db.msgCache.SetMessageBySeqs(ctx, conversationID, datautil.Slice(msgs, msgToDB)); err != nil {
		return 0, false, nil, err
	}
	return lastMaxSeq, isNew, userSeqMap, nil
}

func (db *msgTransferDatabase) SetHasReadSeqs(ctx context.Context, conversationID string, userSeqMap map[string]int64) error {
	for userID, seq := range userSeqMap {
		if err := db.seqUser.SetUserReadSeq(ctx, conversationID, userID, seq); err != nil {
			return err
		}
	}
	return nil
}

func (db *msgTransferDatabase) SetHasReadSeqToDB(ctx context.Context, conversationID string, userSeqMap map[string]int64) error {
	for userID, seq := range userSeqMap {
		if err := db.seqUser.SetUserReadSeqToDB(ctx, conversationID, userID, seq); err != nil {
			return err
		}
	}
	return nil
}

func (db *msgTransferDatabase) MsgToPushMQ(ctx context.Context, key, conversationID string, msg2mq *sdkws.MsgData) (int32, int64, error) {
	partition, offset, err := db.producerToPush.SendMessage(ctx, key, &pbmsg.PushMsgDataToMQ{MsgData: msg2mq, ConversationID: conversationID})
	if err != nil {
		log.ZError(ctx, "MsgToPushMQ", err, "key", key, "msg2mq", msg2mq)
		return 0, 0, err
	}
	return partition, offset, nil
}

func (db *msgTransferDatabase) MsgToMongoMQ(ctx context.Context, key, conversationID string, messages []*sdkws.MsgData, lastSeq int64) error {
	if len(messages) > 0 {
		_, _, err := db.producerToMongo.SendMessage(ctx, key, &pbmsg.MsgDataToMongoByMQ{LastSeq: lastSeq, ConversationID: conversationID, MsgData: messages})
		if err != nil {
			log.ZError(ctx, "MsgToMongoMQ", err, "key", key, "conversationID", conversationID, "lastSeq", lastSeq)
			return err
		}
	}
	return nil
}
