package controller

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/convert"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	unRelationTb "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/unrelation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/unrelation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/kafka"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/prome"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"

	"context"
	"errors"

	pbMsg "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	updateKeyMsg = iota
	updateKeyRevoke
)

type CommonMsgDatabase interface {
	// 批量插入消息
	BatchInsertChat2DB(ctx context.Context, conversationID string, msgs []*sdkws.MsgData, currentMaxSeq int64) error
	// 撤回消息
	RevokeMsg(ctx context.Context, conversationID string, seq int64, revoke *unRelationTb.RevokeModel) error
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

	GetMongoMaxAndMinSeq(ctx context.Context, conversationID string) (maxSeq, minSeq int64, err error)
	GetConversationMinMaxSeqInMongoAndCache(ctx context.Context, conversationID string) (minSeqMongo, maxSeqMongo, minSeqCache, maxSeqCache int64, err error)
	SetSendMsgStatus(ctx context.Context, id string, status int32) error
	GetSendMsgStatus(ctx context.Context, id string) (int32, error)

	// to mq
	MsgToMQ(ctx context.Context, key string, msg2mq *sdkws.MsgData) error
	MsgToModifyMQ(ctx context.Context, key, conversarionID string, msgs []*sdkws.MsgData) error
	MsgToPushMQ(ctx context.Context, key, conversarionID string, msg2mq *sdkws.MsgData) (int32, int64, error)
	MsgToMongoMQ(ctx context.Context, key, conversarionID string, msgs []*sdkws.MsgData, lastSeq int64) error

	// modify
	JudgeMessageReactionExist(ctx context.Context, clientMsgID string, sessionType int32) (bool, error)
	SetMessageTypeKeyValue(ctx context.Context, clientMsgID string, sessionType int32, typeKey, value string) error
	SetMessageReactionExpire(ctx context.Context, clientMsgID string, sessionType int32, expiration time.Duration) (bool, error)
	GetExtendMsg(ctx context.Context, conversationID string, sessionType int32, clientMsgID string, maxMsgUpdateTime int64) (*pbMsg.ExtendMsg, error)
	InsertOrUpdateReactionExtendMsgSet(ctx context.Context, conversationID string, sessionType int32, clientMsgID string, msgFirstModifyTime int64, reactionExtensionList map[string]*sdkws.KeyValue) error
	GetMessageTypeKeyValue(ctx context.Context, clientMsgID string, sessionType int32, typeKey string) (string, error)
	GetOneMessageAllReactionList(ctx context.Context, clientMsgID string, sessionType int32) (map[string]string, error)
	DeleteOneMessageKey(ctx context.Context, clientMsgID string, sessionType int32, subKey string) error
	DeleteReactionExtendMsgSet(ctx context.Context, conversationID string, sessionType int32, clientMsgID string, msgFirstModifyTime int64, reactionExtensionList map[string]*sdkws.KeyValue) error
}

func NewCommonMsgDatabase(msgDocModel unRelationTb.MsgDocModelInterface, cacheModel cache.MsgModel) CommonMsgDatabase {
	return &commonMsgDatabase{
		msgDocDatabase:   msgDocModel,
		cache:            cacheModel,
		producer:         kafka.NewKafkaProducer(config.Config.Kafka.Addr, config.Config.Kafka.LatestMsgToRedis.Topic),
		producerToMongo:  kafka.NewKafkaProducer(config.Config.Kafka.Addr, config.Config.Kafka.MsgToMongo.Topic),
		producerToPush:   kafka.NewKafkaProducer(config.Config.Kafka.Addr, config.Config.Kafka.MsgToPush.Topic),
		producerToModify: kafka.NewKafkaProducer(config.Config.Kafka.Addr, config.Config.Kafka.MsgToModify.Topic),
	}
}

func InitCommonMsgDatabase(rdb redis.UniversalClient, database *mongo.Database) CommonMsgDatabase {
	cacheModel := cache.NewMsgCacheModel(rdb)
	msgDocModel := unrelation.NewMsgMongoDriver(database)
	CommonMsgDatabase := NewCommonMsgDatabase(msgDocModel, cacheModel)
	return CommonMsgDatabase
}

type commonMsgDatabase struct {
	msgDocDatabase    unRelationTb.MsgDocModelInterface
	extendMsgDatabase unRelationTb.ExtendMsgSetModelInterface
	extendMsgSetModel unRelationTb.ExtendMsgSetModel
	msg               unRelationTb.MsgDocModel
	cache             cache.MsgModel
	producer          *kafka.Producer
	producerToMongo   *kafka.Producer
	producerToModify  *kafka.Producer
	producerToPush    *kafka.Producer
}

func (db *commonMsgDatabase) MsgToMQ(ctx context.Context, key string, msg2mq *sdkws.MsgData) error {
	_, _, err := db.producer.SendMessage(ctx, key, msg2mq)
	return err
}

func (db *commonMsgDatabase) MsgToModifyMQ(ctx context.Context, key, conversationID string, messages []*sdkws.MsgData) error {
	if len(messages) > 0 {
		_, _, err := db.producerToModify.SendMessage(ctx, key, &pbMsg.MsgDataToModifyByMQ{ConversationID: conversationID, Messages: messages})
		return err
	}
	return nil
}

func (db *commonMsgDatabase) MsgToPushMQ(ctx context.Context, key, conversationID string, msg2mq *sdkws.MsgData) (int32, int64, error) {
	partition, offset, err := db.producerToPush.SendMessage(ctx, key, &pbMsg.PushMsgDataToMQ{MsgData: msg2mq, ConversationID: conversationID})
	if err != nil {
		log.ZError(ctx, "MsgToPushMQ", err, "key", key, "msg2mq", msg2mq)
		return 0, 0, err
	}
	return partition, offset, nil
}

func (db *commonMsgDatabase) MsgToMongoMQ(ctx context.Context, key, conversationID string, messages []*sdkws.MsgData, lastSeq int64) error {
	if len(messages) > 0 {
		_, _, err := db.producerToMongo.SendMessage(ctx, key, &pbMsg.MsgDataToMongoByMQ{LastSeq: lastSeq, ConversationID: conversationID, MsgData: messages})
		return err
	}
	return nil
}

func (db *commonMsgDatabase) BatchInsertBlock(ctx context.Context, conversationID string, fields []any, key int8, firstSeq int64) error {
	if len(fields) == 0 {
		return nil
	}
	num := db.msg.GetSingleGocMsgNum()
	//num = 100
	for i, field := range fields { // 检查类型
		var ok bool
		switch key {
		case updateKeyMsg:
			var msg *unRelationTb.MsgDataModel
			msg, ok = field.(*unRelationTb.MsgDataModel)
			if msg != nil && msg.Seq != firstSeq+int64(i) {
				return errs.ErrInternalServer.Wrap("seq is invalid")
			}
		case updateKeyRevoke:
			_, ok = field.(*unRelationTb.RevokeModel)
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
		doc := unRelationTb.MsgDocModel{
			DocID: db.msg.GetDocID(conversationID, seq),
			Msg:   make([]*unRelationTb.MsgInfoModel, num),
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
				doc.Msg[db.msg.GetMsgIndex(seq)] = &unRelationTb.MsgInfoModel{
					Msg: fields[j].(*unRelationTb.MsgDataModel),
				}
			case updateKeyRevoke:
				doc.Msg[db.msg.GetMsgIndex(seq)] = &unRelationTb.MsgInfoModel{
					Revoke: fields[j].(*unRelationTb.RevokeModel),
				}
			}
		}
		for i, model := range doc.Msg {
			if model == nil {
				model = &unRelationTb.MsgInfoModel{}
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
		var offlinePushModel *unRelationTb.OfflinePushModel
		if msg.OfflinePushInfo != nil {
			offlinePushModel = &unRelationTb.OfflinePushModel{
				Title:         msg.OfflinePushInfo.Title,
				Desc:          msg.OfflinePushInfo.Desc,
				Ex:            msg.OfflinePushInfo.Ex,
				IOSPushSound:  msg.OfflinePushInfo.IOSPushSound,
				IOSBadgeCount: msg.OfflinePushInfo.IOSBadgeCount,
			}
		}
		msgs[i] = &unRelationTb.MsgDataModel{
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

func (db *commonMsgDatabase) RevokeMsg(ctx context.Context, conversationID string, seq int64, revoke *unRelationTb.RevokeModel) error {
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
		prome.Inc(prome.SeqGetFailedCounter)
		return 0, false, err
	}
	prome.Inc(prome.SeqGetSuccessCounter)
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
		prome.Add(prome.MsgInsertRedisFailedCounter, failedNum)
		log.ZError(ctx, "setMessageToCache error", err, "len", len(msgs), "conversationID", conversationID)
	} else {
		prome.Inc(prome.MsgInsertRedisSuccessCounter)
	}
	err = db.cache.SetMaxSeq(ctx, conversationID, currentMaxSeq)
	if err != nil {
		prome.Inc(prome.SeqSetFailedCounter)
	} else {
		prome.Inc(prome.SeqSetSuccessCounter)
	}
	err2 := db.cache.SetHasReadSeqs(ctx, conversationID, userSeqMap)
	if err != nil {
		log.ZError(ctx, "SetHasReadSeqs error", err2, "userSeqMap", userSeqMap, "conversationID", conversationID)
		prome.Inc(prome.SeqSetFailedCounter)
	} else {
		prome.Inc(prome.SeqSetSuccessCounter)
	}
	return lastMaxSeq, isNew, utils.Wrap(err, "")
}

func (db *commonMsgDatabase) getMsgBySeqs(ctx context.Context, userID, conversationID string, seqs []int64) (totalMsgs []*sdkws.MsgData, err error) {
	for docID, seqs := range db.msg.GetDocIDSeqsMap(conversationID, seqs) {
		//log.ZDebug(ctx, "getMsgBySeqs", "docID", docID, "seqs", seqs)
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

// func (db *commonMsgDatabase) refetchDelSeqsMsgs(ctx context.Context, conversationID string, delNums, rangeBegin, begin int64) (seqMsgs []*unRelationTb.MsgDataModel, err error) {
// 	var reFetchSeqs []int64
// 	if delNums > 0 {
// 		newBeginSeq := rangeBegin - delNums
// 		if newBeginSeq >= begin {
// 			newEndSeq := rangeBegin - 1
// 			for i := newBeginSeq; i <= newEndSeq; i++ {
// 				reFetchSeqs = append(reFetchSeqs, i)
// 			}
// 		}
// 	}
// 	if len(reFetchSeqs) == 0 {
// 		return
// 	}
// 	if len(reFetchSeqs) > 0 {
// 		m := db.msg.GetDocIDSeqsMap(conversationID, reFetchSeqs)
// 		for docID, seqs := range m {
// 			msgs, _, err := db.findMsgInfoBySeq(ctx, docID, seqs)
// 			if err != nil {
// 				return nil, err
// 			}
// 			for _, msg := range msgs {
// 				if msg.Status != constant.MsgDeleted {
// 					seqMsgs = append(seqMsgs, msg)
// 				}
// 			}
// 		}
// 	}
// 	if len(seqMsgs) < int(delNums) {
// 		seqMsgs2, err := db.refetchDelSeqsMsgs(ctx, conversationID, delNums-int64(len(seqMsgs)), rangeBegin-1, begin)
// 		if err != nil {
// 			return seqMsgs, err
// 		}
// 		seqMsgs = append(seqMsgs, seqMsgs2...)
// 	}
// 	return seqMsgs, nil
// }

func (db *commonMsgDatabase) findMsgInfoBySeq(ctx context.Context, userID, docID string, seqs []int64) (totalMsgs []*unRelationTb.MsgInfoModel, err error) {
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
	if begin < minSeq {
		begin = minSeq
	}
	if end > maxSeq {
		end = maxSeq
	}
	if end < begin {
		return 0, 0, nil, errs.ErrArgs.Wrap("seq end < begin")
	}
	var seqs []int64
	for i := end; i > end-num; i-- {
		if i >= begin {
			seqs = append([]int64{i}, seqs...)
		} else {
			break
		}
	}
	if len(seqs) == 0 {
		return 0, 0, nil, nil
	}
	newBegin := seqs[0]
	newEnd := seqs[len(seqs)-1]
	log.ZDebug(ctx, "GetMsgBySeqsRange", "first seqs", seqs, "newBegin", newBegin, "newEnd", newEnd)
	cachedMsgs, failedSeqs, err := db.cache.GetMessagesBySeq(ctx, conversationID, seqs)
	if err != nil {
		if err != redis.Nil {
			prome.Add(prome.MsgPullFromRedisFailedCounter, len(failedSeqs))
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
					prome.Add(prome.MsgPullFromRedisFailedCounter, len(failedSeqs2))
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
	prome.Add(prome.MsgPullFromRedisSuccessCounter, len(successMsgs))
	if len(failedSeqs) > 0 {
		mongoMsgs, err := db.getMsgBySeqsRange(ctx, userID, conversationID, failedSeqs, begin, end)
		if err != nil {
			prome.Add(prome.MsgPullFromMongoFailedCounter, len(failedSeqs))
			return 0, 0, nil, err
		}
		prome.Add(prome.MsgPullFromMongoSuccessCounter, len(mongoMsgs))
		successMsgs = append(successMsgs, mongoMsgs...)
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
			prome.Add(prome.MsgPullFromRedisFailedCounter, len(failedSeqs))
			log.ZError(ctx, "get message from redis exception", err, "failedSeqs", failedSeqs, "conversationID", conversationID)
		}
	}
	log.ZInfo(ctx, "db.cache.GetMessagesBySeq", "userID", userID, "conversationID", conversationID, "seqs", seqs, "successMsgs", len(successMsgs), "failedSeqs", failedSeqs, "conversationID", conversationID)
	prome.Add(prome.MsgPullFromRedisSuccessCounter, len(successMsgs))
	if len(failedSeqs) > 0 {
		mongoMsgs, err := db.getMsgBySeqs(ctx, userID, conversationID, failedSeqs)
		if err != nil {
			prome.Add(prome.MsgPullFromMongoFailedCounter, len(failedSeqs))
			return 0, 0, nil, err
		}
		prome.Add(prome.MsgPullFromMongoSuccessCounter, len(mongoMsgs))
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

// this is struct for recursion
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
// recursion 删除list并且返回设置的最小seq
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
		var hasMarkDelFlag bool
		var delMsgIndexs []int
		for i, MsgInfoModel := range msgDocModel.Msg {
			if MsgInfoModel != nil && MsgInfoModel.Msg != nil {
				if utils.GetCurrentTimestampByMill() > MsgInfoModel.Msg.SendTime+(remainTime*1000) {
					delMsgIndexs = append(delMsgIndexs, i)
					hasMarkDelFlag = true
				} else {
					// 到本条消息不需要删除, minSeq置为这条消息的seq
					if len(delStruct.delDocIDs) > 0 {
						log.ZDebug(ctx, "delete docs", "delDocIDs", delStruct.delDocIDs)
					}
					if err := db.msgDocDatabase.DeleteDocs(ctx, delStruct.delDocIDs); err != nil {
						return 0, err
					}
					if hasMarkDelFlag {
						log.ZDebug(ctx, "delete msg by index", "delMsgIndexs", delMsgIndexs, "docID", msgDocModel.DocID)
						// mark del all delMsgIndexs
						if err := db.msgDocDatabase.DeleteMsgsInOneDocByIndex(ctx, msgDocModel.DocID, delMsgIndexs); err != nil {
							return delStruct.getSetMinSeq(), err
						}
					}
					return MsgInfoModel.Msg.Seq, nil
				}
			}
		}
	}
	//  继续递归 index+1
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

func (db *commonMsgDatabase) GetMongoMaxAndMinSeq(ctx context.Context, conversationID string) (maxSeq, minSeq int64, err error) {
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

func (db *commonMsgDatabase) JudgeMessageReactionExist(ctx context.Context, clientMsgID string, sessionType int32) (bool, error) {
	return db.cache.JudgeMessageReactionExist(ctx, clientMsgID, sessionType)
}

func (db *commonMsgDatabase) SetMessageTypeKeyValue(ctx context.Context, clientMsgID string, sessionType int32, typeKey, value string) error {
	return db.cache.SetMessageTypeKeyValue(ctx, clientMsgID, sessionType, typeKey, value)
}

func (db *commonMsgDatabase) SetMessageReactionExpire(ctx context.Context, clientMsgID string, sessionType int32, expiration time.Duration) (bool, error) {
	return db.cache.SetMessageReactionExpire(ctx, clientMsgID, sessionType, expiration)
}

func (db *commonMsgDatabase) GetMessageTypeKeyValue(ctx context.Context, clientMsgID string, sessionType int32, typeKey string) (string, error) {
	return db.cache.GetMessageTypeKeyValue(ctx, clientMsgID, sessionType, typeKey)
}

func (db *commonMsgDatabase) GetOneMessageAllReactionList(ctx context.Context, clientMsgID string, sessionType int32) (map[string]string, error) {
	return db.cache.GetOneMessageAllReactionList(ctx, clientMsgID, sessionType)
}

func (db *commonMsgDatabase) DeleteOneMessageKey(ctx context.Context, clientMsgID string, sessionType int32, subKey string) error {
	return db.cache.DeleteOneMessageKey(ctx, clientMsgID, sessionType, subKey)
}

func (db *commonMsgDatabase) InsertOrUpdateReactionExtendMsgSet(ctx context.Context, conversationID string, sessionType int32, clientMsgID string, msgFirstModifyTime int64, reactionExtensions map[string]*sdkws.KeyValue) error {
	return db.extendMsgDatabase.InsertOrUpdateReactionExtendMsgSet(ctx, conversationID, sessionType, clientMsgID, msgFirstModifyTime, db.extendMsgSetModel.Pb2Model(reactionExtensions))
}

func (db *commonMsgDatabase) GetExtendMsg(ctx context.Context, conversationID string, sessionType int32, clientMsgID string, maxMsgUpdateTime int64) (*pbMsg.ExtendMsg, error) {
	extendMsgSet, err := db.extendMsgDatabase.GetExtendMsgSet(ctx, conversationID, sessionType, maxMsgUpdateTime)
	if err != nil {
		return nil, err
	}
	extendMsg, ok := extendMsgSet.ExtendMsgs[clientMsgID]
	if !ok {
		return nil, errs.ErrRecordNotFound.Wrap(fmt.Sprintf("cant find client msg id: %s", clientMsgID))
	}
	reactionExtensionList := make(map[string]*pbMsg.KeyValueResp)
	for key, model := range extendMsg.ReactionExtensionList {
		reactionExtensionList[key] = &pbMsg.KeyValueResp{
			KeyValue: &sdkws.KeyValue{
				TypeKey:          model.TypeKey,
				Value:            model.Value,
				LatestUpdateTime: model.LatestUpdateTime,
			},
		}
	}
	return &pbMsg.ExtendMsg{
		ReactionExtensions: reactionExtensionList,
		ClientMsgID:        extendMsg.ClientMsgID,
		MsgFirstModifyTime: extendMsg.MsgFirstModifyTime,
		AttachedInfo:       extendMsg.AttachedInfo,
		Ex:                 extendMsg.Ex,
	}, nil
}

func (db *commonMsgDatabase) DeleteReactionExtendMsgSet(ctx context.Context, conversationID string, sessionType int32, clientMsgID string, msgFirstModifyTime int64, reactionExtensions map[string]*sdkws.KeyValue) error {
	return db.extendMsgDatabase.DeleteReactionExtendMsgSet(ctx, conversationID, sessionType, clientMsgID, msgFirstModifyTime, db.extendMsgSetModel.Pb2Model(reactionExtensions))
}
