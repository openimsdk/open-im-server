package controller

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	unRelationTb "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/unrelation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/unrelation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/kafka"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/prome"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/gogo/protobuf/sortkeys"

	"context"
	"errors"

	pbMsg "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/golang/protobuf/proto"
)

type CommonMsgDatabase interface {
	// 批量插入消息
	BatchInsertChat2DB(ctx context.Context, conversationID string, msgs []*sdkws.MsgData, currentMaxSeq int64) error
	// 刪除redis中消息缓存
	DeleteMessageFromCache(ctx context.Context, conversationID string, msgs []*sdkws.MsgData) error
	// incrSeq然后批量插入缓存
	BatchInsertChat2Cache(ctx context.Context, conversationID string, msgs []*sdkws.MsgData) (seq int64, isNewConversation bool, err error)
	// 删除消息 返回不存在的seqList
	DelMsgBySeqs(ctx context.Context, conversationID string, seqs []int64) (totalUnExistSeqs []int64, err error)
	//  通过seqList获取mongo中写扩散消息
	GetMsgBySeqsRange(ctx context.Context, conversationID string, begin, end, num int64) (seqMsg []*sdkws.MsgData, err error)
	// 通过seqList获取大群在 mongo里面的消息
	GetMsgBySeqs(ctx context.Context, conversationID string, seqs []int64) (seqMsg []*sdkws.MsgData, err error)
	// 删除会话消息重置最小seq， remainTime为消息保留的时间单位秒,超时消息删除， 传0删除所有消息(此方法不删除redis cache)
	DeleteConversationMsgsAndSetMinSeq(ctx context.Context, conversationID string, remainTime int64) error
	CleanUpUserConversationsMsgs(ctx context.Context, userID string, conversationIDs []string)
	SetMaxSeq(ctx context.Context, conversationID string, maxSeq int64) error
	GetMaxSeqs(ctx context.Context, conversationIDs []string) (map[string]int64, error)
	GetMaxSeq(ctx context.Context, conversationID string) (int64, error)
	SetMinSeq(ctx context.Context, conversationID string, minSeq int64) error
	GetMinSeqs(ctx context.Context, conversationIDs []string) (map[string]int64, error)
	GetMinSeq(ctx context.Context, conversationID string) (int64, error)
	GetConversationUserMinSeq(ctx context.Context, conversationID string, userID string) (int64, error)
	GetConversationUserMinSeqs(ctx context.Context, conversationID string, userIDs []string) (map[string]int64, error)
	SetConversationUserMinSeq(ctx context.Context, conversationID string, userID string, minSeq int64) error
	SetConversationUserMinSeqs(ctx context.Context, conversationID string, seqs map[string]int64) (err error)

	GetConversationMinMaxSeqInMongoAndCache(ctx context.Context, conversationID string) (minSeqMongo, maxSeqMongo, minSeqCache, maxSeqCache int64, err error)
	SetSendMsgStatus(ctx context.Context, id string, status int32) error
	GetSendMsgStatus(ctx context.Context, id string) (int32, error)

	// to mq
	MsgToMQ(ctx context.Context, key string, msg2mq *sdkws.MsgData) error
	MsgToModifyMQ(ctx context.Context, conversationID string, msgs []*sdkws.MsgData) error
	MsgToPushMQ(ctx context.Context, conversationID string, msg2mq *sdkws.MsgData) (int32, int64, error)
	MsgToMongoMQ(ctx context.Context, conversationID string, msgs []*sdkws.MsgData, lastSeq int64) error

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
		producer:         kafka.NewKafkaProducer(config.Config.Kafka.Ws2mschat.Addr, config.Config.Kafka.Ws2mschat.Topic),
		producerToMongo:  kafka.NewKafkaProducer(config.Config.Kafka.MsgToMongo.Addr, config.Config.Kafka.MsgToMongo.Topic),
		producerToPush:   kafka.NewKafkaProducer(config.Config.Kafka.Ms2pschat.Addr, config.Config.Kafka.Ms2pschat.Topic),
		producerToModify: kafka.NewKafkaProducer(config.Config.Kafka.MsgToModify.Addr, config.Config.Kafka.MsgToModify.Topic),
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
	cache             cache.MsgModel
	producer          *kafka.Producer
	producerToMongo   *kafka.Producer
	producerToModify  *kafka.Producer
	producerToPush    *kafka.Producer
	// model
	msg unRelationTb.MsgDocModel
}

func (db *commonMsgDatabase) MsgToMQ(ctx context.Context, key string, msg2mq *sdkws.MsgData) error {
	_, _, err := db.producer.SendMessage(ctx, key, msg2mq)
	return err
}

func (db *commonMsgDatabase) MsgToModifyMQ(ctx context.Context, conversationID string, messages []*sdkws.MsgData) error {
	if len(messages) > 0 {
		_, _, err := db.producerToModify.SendMessage(ctx, conversationID, &pbMsg.MsgDataToModifyByMQ{ConversationID: conversationID, Messages: messages})
		return err
	}
	return nil
}

func (db *commonMsgDatabase) MsgToPushMQ(ctx context.Context, conversationID string, msg2mq *sdkws.MsgData) (int32, int64, error) {
	partition, offset, err := db.producerToPush.SendMessage(ctx, conversationID, &pbMsg.PushMsgDataToMQ{MsgData: msg2mq, ConversationID: conversationID})
	if err != nil {
		log.ZError(ctx, "MsgToPushMQ", err, "key", conversationID, "msg2mq", msg2mq)
	}
	return partition, offset, err
}

func (db *commonMsgDatabase) MsgToMongoMQ(ctx context.Context, conversationID string, messages []*sdkws.MsgData, lastSeq int64) error {
	if len(messages) > 0 {
		_, _, err := db.producerToModify.SendMessage(ctx, conversationID, &pbMsg.MsgDataToMongoByMQ{LastSeq: lastSeq, ConversationID: conversationID, MsgData: messages})
		return err
	}
	return nil
}

func (db *commonMsgDatabase) BatchInsertChat2DB(ctx context.Context, conversationID string, msgList []*sdkws.MsgData, currentMaxSeq int64) error {
	if int64(len(msgList)) > db.msg.GetSingleGocMsgNum() {
		return errors.New("too large")
	}
	var remain int64
	blk0 := db.msg.GetSingleGocMsgNum() - 1
	//currentMaxSeq 4998
	if currentMaxSeq < db.msg.GetSingleGocMsgNum() {
		remain = blk0 - currentMaxSeq //1
	} else {
		excludeBlk0 := currentMaxSeq - blk0 //=1
		//(5000-1)%5000 == 4999
		remain = (db.msg.GetSingleGocMsgNum() - (excludeBlk0 % db.msg.GetSingleGocMsgNum())) % db.msg.GetSingleGocMsgNum()
	}
	//remain=1
	var insertCounter int64
	msgsToMongo := make([]unRelationTb.MsgInfoModel, 0)
	msgsToMongoNext := make([]unRelationTb.MsgInfoModel, 0)
	docID := ""
	docIDNext := ""
	var err error
	for _, m := range msgList {
		currentMaxSeq++
		sMsg := unRelationTb.MsgInfoModel{}
		sMsg.SendTime = m.SendTime
		m.Seq = currentMaxSeq
		if sMsg.Msg, err = proto.Marshal(m); err != nil {
			return utils.Wrap(err, "")
		}
		if insertCounter < remain {
			msgsToMongo = append(msgsToMongo, sMsg)
			insertCounter++
			docID = db.msg.GetDocID(conversationID, currentMaxSeq)
		} else {
			msgsToMongoNext = append(msgsToMongoNext, sMsg)
			docIDNext = db.msg.GetDocID(conversationID, currentMaxSeq)
		}
	}

	if docID != "" {
		err = db.msgDocDatabase.PushMsgsToDoc(ctx, docID, msgsToMongo)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				doc := &unRelationTb.MsgDocModel{}
				doc.DocID = docID
				doc.Msg = msgsToMongo
				if err = db.msgDocDatabase.Create(ctx, doc); err != nil {
					prome.Inc(prome.MsgInsertMongoFailedCounter)
					return err
				}
				prome.Inc(prome.MsgInsertMongoSuccessCounter)
			} else {
				prome.Inc(prome.MsgInsertMongoFailedCounter)
				return err
			}
		} else {
			log.ZDebug(ctx, "PushMsgsToDoc success", "docID", docID, "len", len(msgsToMongo))
			prome.Inc(prome.MsgInsertMongoSuccessCounter)
		}
	}
	if docIDNext != "" {
		nextDoc := &unRelationTb.MsgDocModel{}
		nextDoc.DocID = docIDNext
		nextDoc.Msg = msgsToMongoNext
		log.ZDebug(ctx, "create next doc", "docIDNext", docIDNext, "len", len(nextDoc.Msg))
		if err = db.msgDocDatabase.Create(ctx, nextDoc); err != nil {
			prome.Inc(prome.MsgInsertMongoFailedCounter)
			return utils.Wrap(err, "")
		}
		prome.Inc(prome.MsgInsertMongoSuccessCounter)
	}
	return nil
}

func (db *commonMsgDatabase) DeleteMessageFromCache(ctx context.Context, conversationID string, msgs []*sdkws.MsgData) error {
	return db.cache.DeleteMessageFromCache(ctx, conversationID, msgs)
}

func (db *commonMsgDatabase) BatchInsertChat2Cache(ctx context.Context, conversationID string, msgs []*sdkws.MsgData) (seq int64, isNew bool, err error) {
	currentMaxSeq, err := db.cache.GetMaxSeq(ctx, conversationID)
	if err != nil && err != redis.Nil {
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
	if err == redis.Nil {
		isNew = true
	}
	lastMaxSeq := currentMaxSeq
	for _, m := range msgs {
		currentMaxSeq++
		m.Seq = currentMaxSeq
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
	return lastMaxSeq, isNew, utils.Wrap(err, "")
}

func (db *commonMsgDatabase) DelMsgBySeqs(ctx context.Context, conversationID string, seqs []int64) (totalUnExistSeqs []int64, err error) {
	sortkeys.Int64s(seqs)
	docIDSeqsMap := db.msg.GetDocIDSeqsMap(conversationID, seqs)
	lock := sync.Mutex{}
	var wg sync.WaitGroup
	wg.Add(len(docIDSeqsMap))
	for k, v := range docIDSeqsMap {
		go func(docID string, seqs []int64) {
			defer wg.Done()
			unExistSeqList, err := db.DelMsgBySeqsInOneDoc(ctx, docID, seqs)
			if err != nil {
				return
			}
			lock.Lock()
			totalUnExistSeqs = append(totalUnExistSeqs, unExistSeqList...)
			lock.Unlock()
		}(k, v)
	}
	return totalUnExistSeqs, nil
}

func (db *commonMsgDatabase) DelMsgBySeqsInOneDoc(ctx context.Context, docID string, seqs []int64) (unExistSeqs []int64, err error) {
	seqMsgs, indexes, unExistSeqs, err := db.GetMsgAndIndexBySeqsInOneDoc(ctx, docID, seqs)
	if err != nil {
		return nil, err
	}
	for i, v := range seqMsgs {
		if err = db.msgDocDatabase.UpdateMsgStatusByIndexInOneDoc(ctx, docID, v, indexes[i], constant.MsgDeleted); err != nil {
			return nil, err
		}
	}
	return unExistSeqs, nil
}

func (db *commonMsgDatabase) GetMsgAndIndexBySeqsInOneDoc(ctx context.Context, docID string, seqs []int64) (seqMsgs []*sdkws.MsgData, indexes []int, unExistSeqs []int64, err error) {
	doc, err := db.msgDocDatabase.FindOneByDocID(ctx, docID)
	if err != nil {
		return nil, nil, nil, err
	}
	singleCount := 0
	var hasSeqList []int64
	for i := 0; i < len(doc.Msg); i++ {
		msgPb, err := db.unmarshalMsg(&doc.Msg[i])
		if err != nil {
			return nil, nil, nil, err
		}
		if utils.Contain(msgPb.Seq, seqs...) {
			indexes = append(indexes, i)
			seqMsgs = append(seqMsgs, msgPb)
			hasSeqList = append(hasSeqList, msgPb.Seq)
			singleCount++
			if singleCount == len(seqs) {
				break
			}
		}
	}
	for _, i := range seqs {
		if utils.Contain(i, hasSeqList...) {
			continue
		}
		unExistSeqs = append(unExistSeqs, i)
	}
	return seqMsgs, indexes, unExistSeqs, nil
}

func (db *commonMsgDatabase) GetNewestMsg(ctx context.Context, conversationID string) (msgPb *sdkws.MsgData, err error) {
	msgInfo, err := db.msgDocDatabase.GetNewestMsg(ctx, conversationID)
	if err != nil {
		return nil, err
	}
	return db.unmarshalMsg(msgInfo)
}

func (db *commonMsgDatabase) GetOldestMsg(ctx context.Context, conversationID string) (msgPb *sdkws.MsgData, err error) {
	msgInfo, err := db.msgDocDatabase.GetOldestMsg(ctx, conversationID)
	if err != nil {
		return nil, err
	}
	return db.unmarshalMsg(msgInfo)
}

func (db *commonMsgDatabase) unmarshalMsg(msgInfo *unRelationTb.MsgInfoModel) (msgPb *sdkws.MsgData, err error) {
	msgPb = &sdkws.MsgData{}
	err = proto.Unmarshal(msgInfo.Msg, msgPb)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	return msgPb, nil
}

func (db *commonMsgDatabase) getMsgBySeqs(ctx context.Context, conversationID string, seqs []int64) (seqMsgs []*sdkws.MsgData, err error) {
	var hasSeqs []int64
	singleCount := 0
	m := db.msg.GetDocIDSeqsMap(conversationID, seqs)
	for docID, value := range m {
		doc, err := db.msgDocDatabase.FindOneByDocID(ctx, docID)
		if err != nil {
			log.ZError(ctx, "get message from mongo exception", err, "docID", docID)
			continue
		}
		singleCount = 0
		for i := 0; i < len(doc.Msg); i++ {
			msgPb, err := db.unmarshalMsg(&doc.Msg[i])
			if err != nil {
				log.ZError(ctx, "unmarshal message exception", err, "docID", docID, "msg", &doc.Msg[i])
				return nil, err
			}
			if utils.Contain(msgPb.Seq, value...) {
				seqMsgs = append(seqMsgs, msgPb)
				hasSeqs = append(hasSeqs, msgPb.Seq)
				singleCount++
				if singleCount == len(value) {
					break
				}
			}
		}
	}
	if len(hasSeqs) != len(seqs) {
		var diff []int64
		var exceptionMsg []*sdkws.MsgData
		diff = utils.Difference(hasSeqs, seqs)
		exceptionMsg = db.msg.GenExceptionSuperGroupMessageBySeqs(diff, conversationID)
		seqMsgs = append(seqMsgs, exceptionMsg...)
	}
	return seqMsgs, nil
}

func (db *commonMsgDatabase) refetchDelSeqsMsgs(ctx context.Context, conversationID string, delNums, rangeBegin, begin int64) (seqMsgs []*sdkws.MsgData, err error) {
	var reFetchSeqs []int64
	if delNums > 0 {
		newBeginSeq := rangeBegin - delNums
		if newBeginSeq >= begin {
			newEndSeq := rangeBegin - 1
			for i := newBeginSeq; i <= newEndSeq; i++ {
				reFetchSeqs = append(reFetchSeqs, i)
			}
		}
	}
	if len(reFetchSeqs) == 0 {
		return
	}
	if len(reFetchSeqs) > 0 {
		m := db.msg.GetDocIDSeqsMap(conversationID, reFetchSeqs)
		for docID, seq := range m {
			msgs, _, err := db.findMsgBySeq(ctx, docID, seq)
			if err != nil {
				return nil, err
			}
			for _, msg := range msgs {
				if msg.Status != constant.MsgDeleted {
					seqMsgs = append(seqMsgs, msg)
				}
			}
		}
	}
	if len(seqMsgs) < int(delNums) {
		seqMsgs2, err := db.refetchDelSeqsMsgs(ctx, conversationID, delNums-int64(len(seqMsgs)), rangeBegin-1, begin)
		if err != nil {
			return seqMsgs, err
		}
		seqMsgs = append(seqMsgs, seqMsgs2...)
	}
	return seqMsgs, nil
}

func (db *commonMsgDatabase) findMsgBySeq(ctx context.Context, docID string, seqs []int64) (seqMsgs []*sdkws.MsgData, unExistSeqs []int64, err error) {
	beginSeq, endSeq := db.msg.GetSeqsBeginEnd(seqs)
	msgs, _, err := db.msgDocDatabase.GetMsgBySeqIndexIn1Doc(ctx, docID, beginSeq, endSeq)
	if err != nil {
		return nil, nil, err
	}
	for _, seq := range seqs {
		for i, msg := range msgs {
			if seq == msg.Seq {
				seqMsgs = append(seqMsgs, msg)
				continue
			}
			if i == len(msgs)-1 {
				unExistSeqs = append(unExistSeqs, seq)
			}
		}
	}
	msgs, _, unExistSeqs, err = db.GetMsgAndIndexBySeqsInOneDoc(ctx, docID, seqs)
	if err != nil {
		return nil, nil, err
	}
	seqMsgs = append(seqMsgs, msgs...)
	return seqMsgs, unExistSeqs, nil
}

func (db *commonMsgDatabase) getMsgBySeqsRange(ctx context.Context, conversationID string, allSeqs []int64, begin, end, num int64) (seqMsgs []*sdkws.MsgData, err error) {
	log.ZDebug(ctx, "getMsgBySeqsRange", "conversationID", conversationID, "allSeqs", allSeqs, "begin", begin, "end", end, "num", num)
	m := db.msg.GetDocIDSeqsMap(conversationID, allSeqs)
	var totalNotExistSeqs []int64
	// mongo index
	for docID, seqs := range m {
		msgs, notExistSeqs, err := db.findMsgBySeq(ctx, docID, seqs)
		if err != nil {
			return nil, err
		}
		log.ZDebug(ctx, "getMsgBySeqsRange", "docID", docID, "seqs", seqs, "unExistSeqs", notExistSeqs, "msgs", msgs)
		seqMsgs = append(seqMsgs, msgs...)
		totalNotExistSeqs = append(totalNotExistSeqs, notExistSeqs...)
	}
	log.ZDebug(ctx, "getMsgBySeqsRange", "totalNotExistSeqs", totalNotExistSeqs)
	// find by next doc
	if len(totalNotExistSeqs) > 0 {
		m = db.msg.GetDocIDSeqsMap(conversationID, totalNotExistSeqs)
		for docID, seqs := range m {
			docID = db.msg.ToNextDoc(docID)
			msgs, _, unExistSeqs, err := db.GetMsgAndIndexBySeqsInOneDoc(ctx, docID, seqs)
			if err != nil {
				log.ZError(ctx, "get message from mongo exception", err, "docID", docID, "seqs", seqs)
				continue
			}
			seqMsgs = append(seqMsgs, msgs...)
			if len(unExistSeqs) > 0 {
				log.ZWarn(ctx, "some seqs lost in mongo", err, "docID", docID, "seqs", seqs, "unExistSeqs", unExistSeqs)
			}
		}
	}
	var delSeqs []int64
	for _, msg := range seqMsgs {
		if msg.Status == constant.MsgDeleted {
			delSeqs = append(delSeqs, msg.Seq)
		}
	}
	if len(delSeqs) > 0 {
		msgs, err := db.refetchDelSeqsMsgs(ctx, conversationID, int64(len(delSeqs)), allSeqs[0], begin)
		if err != nil {
			log.ZWarn(ctx, "refetchDelSeqsMsgs", err, "delSeqs", delSeqs, "begin", begin)
		}
		seqMsgs = append(seqMsgs, msgs...)
	}
	// sort by seq
	if len(totalNotExistSeqs) > 0 || len(delSeqs) > 0 {
		sort.Sort(utils.MsgBySeq(seqMsgs))
	}
	// missSeqs为依然缺失的
	return seqMsgs, nil
}

func (db *commonMsgDatabase) GetMsgBySeqsRange(ctx context.Context, conversationID string, begin, end, num int64) (seqMsg []*sdkws.MsgData, err error) {
	var seqs []int64
	for i := end; i > end-num; i-- {
		if i >= begin {
			seqs = append(seqs, i)
		} else {
			break
		}
	}
	successMsgs, failedSeqs, err := db.cache.GetMessagesBySeq(ctx, conversationID, seqs)
	if err != nil {
		if err != redis.Nil {
			prome.Add(prome.MsgPullFromRedisFailedCounter, len(failedSeqs))
			log.ZError(ctx, "get message from redis exception", err, conversationID, seqs)
		}
	}
	// get from cache or db
	prome.Add(prome.MsgPullFromRedisSuccessCounter, len(successMsgs))
	if len(failedSeqs) > 0 {
		mongoMsgs, err := db.getMsgBySeqsRange(ctx, conversationID, failedSeqs, begin, end, num-int64(len(successMsgs)))
		if err != nil {
			prome.Add(prome.MsgPullFromMongoFailedCounter, len(failedSeqs))
			return nil, err
		}
		prome.Add(prome.MsgPullFromMongoSuccessCounter, len(mongoMsgs))
		successMsgs = append(successMsgs, mongoMsgs...)
	}
	return successMsgs, nil
}

func (db *commonMsgDatabase) GetMsgBySeqs(ctx context.Context, conversationID string, seqs []int64) (successMsgs []*sdkws.MsgData, err error) {
	successMsgs, failedSeqs, err := db.cache.GetMessagesBySeq(ctx, conversationID, seqs)
	if err != nil {
		if err != redis.Nil {
			prome.Add(prome.MsgPullFromRedisFailedCounter, len(failedSeqs))
			log.ZError(ctx, "get message from redis exception", err, "failedSeqs", failedSeqs, "conversationID", conversationID)
		}
	}
	prome.Add(prome.MsgPullFromRedisSuccessCounter, len(successMsgs))
	if len(failedSeqs) > 0 {
		mongoMsgs, err := db.getMsgBySeqs(ctx, conversationID, seqs)
		if err != nil {
			prome.Add(prome.MsgPullFromMongoFailedCounter, len(failedSeqs))
			return nil, err
		}
		prome.Add(prome.MsgPullFromMongoSuccessCounter, len(mongoMsgs))
		successMsgs = append(successMsgs, mongoMsgs...)
	}
	return successMsgs, nil
}

func (db *commonMsgDatabase) DeleteConversationMsgsAndSetMinSeq(ctx context.Context, conversationID string, remainTime int64) error {
	var delStruct delMsgRecursionStruct
	minSeq, err := db.deleteMsgRecursion(ctx, conversationID, unRelationTb.OldestList, &delStruct, remainTime)
	if err != nil {
		return err
	}
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
	msgs, err := db.msgDocDatabase.GetMsgsByIndex(ctx, conversationID, index)
	if err != nil || msgs.DocID == "" {
		if err != nil {
			if err == unrelation.ErrMsgListNotExist {
				log.ZDebug(ctx, "deleteMsgRecursion ErrMsgListNotExist", "conversationID", conversationID, "index:", index)
			} else {
				log.ZError(ctx, "deleteMsgRecursion GetUserMsgListByIndex failed", err, "conversationID", conversationID, "index", index)
			}
		}
		// 获取报错，或者获取不到了，物理删除并且返回seq delMongoMsgsPhysical(delStruct.delDocIDList), 结束递归
		err = db.msgDocDatabase.Delete(ctx, delStruct.delDocIDs)
		if err != nil {
			return 0, err
		}
		return delStruct.getSetMinSeq() + 1, nil
	}
	log.ZDebug(ctx, "conversationID", conversationID, "index:", index, "docID", msgs.DocID, "len", len(msgs.Msg))
	if int64(len(msgs.Msg)) > db.msg.GetSingleGocMsgNum() {
		log.ZWarn(ctx, "msgs too large", nil, "lenth", len(msgs.Msg), "docID:", msgs.DocID)
	}
	if msgs.Msg[len(msgs.Msg)-1].SendTime+(remainTime*1000) < utils.GetCurrentTimestampByMill() && msgs.IsFull() {
		delStruct.delDocIDs = append(delStruct.delDocIDs, msgs.DocID)
		lastMsgPb := &sdkws.MsgData{}
		err = proto.Unmarshal(msgs.Msg[len(msgs.Msg)-1].Msg, lastMsgPb)
		if err != nil {
			log.ZError(ctx, "proto.Unmarshal failed", err, "index", len(msgs.Msg)-1, "docID", msgs.DocID)
			return 0, utils.Wrap(err, "proto.Unmarshal failed")
		}
		delStruct.minSeq = lastMsgPb.Seq
	} else {
		var hasMarkDelFlag bool
		for i, msg := range msgs.Msg {
			msgPb := &sdkws.MsgData{}
			err = proto.Unmarshal(msg.Msg, msgPb)
			if err != nil {
				log.ZError(ctx, "proto.Unmarshal failed", err, "index", i, "docID", msgs.DocID)
				return 0, utils.Wrap(err, "proto.Unmarshal failed")
			}
			if utils.GetCurrentTimestampByMill() > msg.SendTime+(remainTime*1000) {
				msgPb.Status = constant.MsgDeleted
				bytes, _ := proto.Marshal(msgPb)
				msg.Msg = bytes
				msg.SendTime = 0
				hasMarkDelFlag = true
			} else {
				// 到本条消息不需要删除, minSeq置为这条消息的seq
				if err := db.msgDocDatabase.Delete(ctx, delStruct.delDocIDs); err != nil {
					return 0, err
				}
				if hasMarkDelFlag {
					if err := db.msgDocDatabase.UpdateOneDoc(ctx, msgs); err != nil {
						return delStruct.getSetMinSeq(), utils.Wrap(err, "")
					}
				}
				return msgPb.Seq, nil
			}
		}
	}
	//  继续递归 index+1
	seq, err := db.deleteMsgRecursion(ctx, conversationID, index+1, delStruct, remainTime)
	return seq, err
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
	// from cache
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

func (db *commonMsgDatabase) GetMinMaxSeqMongo(ctx context.Context, conversationID string) (minSeqMongo, maxSeqMongo int64, err error) {
	oldestMsgMongo, err := db.msgDocDatabase.GetOldestMsg(ctx, conversationID)
	if err != nil {
		return
	}
	msgPb, err := db.unmarshalMsg(oldestMsgMongo)
	if err != nil {
		return
	}
	minSeqMongo = msgPb.Seq
	newestMsgMongo, err := db.msgDocDatabase.GetNewestMsg(ctx, conversationID)
	if err != nil {
		return
	}
	msgPb, err = db.unmarshalMsg(newestMsgMongo)
	if err != nil {
		return
	}
	maxSeqMongo = msgPb.Seq
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
