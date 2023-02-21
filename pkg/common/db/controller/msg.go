package controller

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/cache"
	unRelationTb "Open_IM/pkg/common/db/table/unrelation"
	"Open_IM/pkg/common/db/unrelation"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/prome"
	"Open_IM/pkg/common/tracelog"
	"github.com/gogo/protobuf/sortkeys"
	"sync"

	pbMsg "Open_IM/pkg/proto/msg"
	"Open_IM/pkg/proto/sdkws"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/golang/protobuf/proto"
)

type MsgInterface interface {
	// 批量插入消息到db
	BatchInsertChat2DB(ctx context.Context, sourceID string, msgList []*pbMsg.MsgDataToMQ, currentMaxSeq int64) error
	// 刪除redis中消息缓存
	DeleteMessageFromCache(ctx context.Context, sourceID string, msgList []*pbMsg.MsgDataToMQ) error
	// incrSeq然后批量插入缓存
	BatchInsertChat2Cache(ctx context.Context, sourceID string, msgList []*pbMsg.MsgDataToMQ) (int64, error)
	// 删除消息 返回不存在的seqList
	DelMsgBySeqs(ctx context.Context, userID string, seqs []int64) (totalUnExistSeqs []int64, err error)
	//  通过seqList获取db中写扩散消息
	GetMsgBySeqs(ctx context.Context, userID string, seqs []int64) (seqMsg []*sdkws.MsgData, err error)
	// 通过seqList获取大群在db里面的消息 没找到返回错误
	GetSuperGroupMsgBySeqs(ctx context.Context, groupID string, seqs []int64) (seqMsg []*sdkws.MsgData, err error)
	// 删除用户所有消息/cache/db然后重置seq
	CleanUpUserMsg(ctx context.Context, userID string) error
	// 删除大群消息重置群成员最小群seq, remainTime为消息保留的时间单位秒,超时消息删除， 传0删除所有消息(此方法不删除 redis cache)
	DeleteUserSuperGroupMsgsAndSetMinSeq(ctx context.Context, groupID string, userID []string, remainTime int64) error
	// 删除用户消息重置最小seq， remainTime为消息保留的时间单位秒,超时消息删除， 传0删除所有消息(此方法不删除redis cache)
	DeleteUserMsgsAndSetMinSeq(ctx context.Context, userID string, remainTime int64) error
	// 获取用户 seq mongo和redis
	GetUserMinMaxSeqInMongoAndCache(ctx context.Context, userID string) (minSeqMongo, maxSeqMongo, minSeqCache, maxSeqCache int64, err error)
	// 获取群 seq mongo和redis
	GetSuperGroupMinMaxSeqInMongoAndCache(ctx context.Context, groupID string) (minSeqMongo, maxSeqMongo, maxSeqCache int64, err error)
	// 设置群用户最小seq 直接调用cache
	SetGroupUserMinSeq(ctx context.Context, groupID, userID string, minSeq int64) (err error)
	// 设置用户最小seq 直接调用cache
	SetUserMinSeq(ctx context.Context, userID string, minSeq int64) (err error)
}

func NewMsgController(mgo *mongo.Client, rdb redis.UniversalClient) MsgInterface {
	return &MsgController{}
}

type MsgController struct {
	database MsgDatabase
}

func (m *MsgController) BatchInsertChat2DB(ctx context.Context, ID string, msgList []*pbMsg.MsgDataToMQ, currentMaxSeq int64) error {
	return m.database.BatchInsertChat2DB(ctx, ID, msgList, currentMaxSeq)
}

func (m *MsgController) DeleteMessageFromCache(ctx context.Context, sourceID string, msgList []*pbMsg.MsgDataToMQ) error {
	return m.database.DeleteMessageFromCache(ctx, sourceID, msgList)
}

func (m *MsgController) BatchInsertChat2Cache(ctx context.Context, sourceID string, msgList []*pbMsg.MsgDataToMQ) (int64, error) {
	return m.database.BatchInsertChat2Cache(ctx, sourceID, msgList)
}

func (m *MsgController) DelMsgBySeqs(ctx context.Context, userID string, seqs []int64) (totalUnExistSeqs []int64, err error) {
	return m.database.DelMsgBySeqs(ctx, userID, seqs)
}

func (m *MsgController) GetMsgBySeqs(ctx context.Context, userID string, seqs []int64) (seqMsg []*sdkws.MsgData, err error) {
	return m.database.GetMsgBySeqs(ctx, userID, seqs)
}

func (m *MsgController) GetSuperGroupMsgBySeqs(ctx context.Context, groupID string, seqs []int64) (seqMsg []*sdkws.MsgData, err error) {
	return m.database.GetSuperGroupMsgBySeqs(ctx, groupID, seqs)
}

func (m *MsgController) CleanUpUserMsg(ctx context.Context, userID string) error {
	return m.database.CleanUpUserMsg(ctx, userID)
}

func (m *MsgController) DeleteUserSuperGroupMsgsAndSetMinSeq(ctx context.Context, groupID string, userIDs []string, remainTime int64) error {
	return m.database.DeleteUserSuperGroupMsgsAndSetMinSeq(ctx, groupID, userIDs, remainTime)
}

func (m *MsgController) DeleteUserMsgsAndSetMinSeq(ctx context.Context, userID string, remainTime int64) error {
	return m.database.DeleteUserMsgsAndSetMinSeq(ctx, userID, remainTime)
}

func (m *MsgController) GetUserMinMaxSeqInMongoAndCache(ctx context.Context, userID string) (minSeqMongo, maxSeqMongo, minSeqCache, maxSeqCache int64, err error) {
	return m.database.GetUserMinMaxSeqInMongoAndCache(ctx, userID)
}

func (m *MsgController) GetSuperGroupMinMaxSeqInMongoAndCache(ctx context.Context, groupID string) (minSeqMongo, maxSeqMongo, maxSeqCache int64, err error) {
	return m.database.GetSuperGroupMinMaxSeqInMongoAndCache(ctx, groupID)
}

func (m *MsgController) SetGroupUserMinSeq(ctx context.Context, groupID, userID string, minSeq int64) (err error) {
	return m.database.SetGroupUserMinSeq(ctx, groupID, userID, minSeq)
}

func (m *MsgController) SetUserMinSeq(ctx context.Context, userID string, minSeq int64) (err error) {
	return m.database.SetUserMinSeq(ctx, userID, minSeq)
}

type MsgDatabaseInterface interface {
	// 批量插入消息
	BatchInsertChat2DB(ctx context.Context, sourceID string, msgList []*pbMsg.MsgDataToMQ, currentMaxSeq int64) error
	// 刪除redis中消息缓存
	DeleteMessageFromCache(ctx context.Context, sourceID string, msgList []*pbMsg.MsgDataToMQ) error
	// incrSeq然后批量插入缓存
	BatchInsertChat2Cache(ctx context.Context, sourceID string, msgList []*pbMsg.MsgDataToMQ) (int64, error)
	// 删除消息 返回不存在的seqList
	DelMsgBySeqs(ctx context.Context, userID string, seqs []int64) (totalUnExistSeqs []int64, err error)
	// 获取群ID或者UserID最新一条在mongo里面的消息
	//  通过seqList获取mongo中写扩散消息
	GetMsgBySeqs(ctx context.Context, userID string, seqs []int64) (seqMsg []*sdkws.MsgData, err error)
	// 通过seqList获取大群在 mongo里面的消息
	GetSuperGroupMsgBySeqs(ctx context.Context, groupID string, seqs []int64) (seqMsg []*sdkws.MsgData, err error)
	// 删除用户所有消息/redis/mongo然后重置seq
	CleanUpUserMsg(ctx context.Context, userID string) error
	// 删除大群消息重置群成员最小群seq, remainTime为消息保留的时间单位秒,超时消息删除， 传0删除所有消息(此方法不删除 redis cache)
	DeleteUserSuperGroupMsgsAndSetMinSeq(ctx context.Context, groupID string, userIDs []string, remainTime int64) error
	// 删除用户消息重置最小seq， remainTime为消息保留的时间单位秒,超时消息删除， 传0删除所有消息(此方法不删除redis cache)
	DeleteUserMsgsAndSetMinSeq(ctx context.Context, userID string, remainTime int64) error
	// 获取用户 seq mongo和redis
	GetUserMinMaxSeqInMongoAndCache(ctx context.Context, userID string) (minSeqMongo, maxSeqMongo, minSeqCache, maxSeqCache int64, err error)
	// 获取群 seq mongo和redis
	GetSuperGroupMinMaxSeqInMongoAndCache(ctx context.Context, groupID string) (minSeqMongo, maxSeqMongo, maxSeqCache int64, err error)
	// 设置群用户最小seq 直接调用cache
	SetGroupUserMinSeq(ctx context.Context, groupID, userID string, minSeq int64) (err error)
	// 设置用户最小seq 直接调用cache
	SetUserMinSeq(ctx context.Context, userID string, minSeq int64) (err error)
}

type MsgDatabase struct {
	mgo   unRelationTb.MsgDocModelInterface
	cache cache.Cache
	msg   unRelationTb.MsgDocModel
}

func NewMsgDatabase(mgo *mongo.Client, rdb redis.UniversalClient) MsgDatabaseInterface {
	return &MsgDatabase{}
}

func (db *MsgDatabase) BatchInsertChat2DB(ctx context.Context, sourceID string, msgList []*pbMsg.MsgDataToMQ, currentMaxSeq int64) error {
	//newTime := utils.GetCurrentTimestampByMill()
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
		//log.Debug(operationID, "msg node ", m.String(), m.MsgData.ClientMsgID)
		currentMaxSeq++
		sMsg := unRelationTb.MsgInfoModel{}
		sMsg.SendTime = m.MsgData.SendTime
		m.MsgData.Seq = currentMaxSeq
		if sMsg.Msg, err = proto.Marshal(m.MsgData); err != nil {
			return utils.Wrap(err, "")
		}
		if insertCounter < remain {
			msgsToMongo = append(msgsToMongo, sMsg)
			insertCounter++
			docID = db.msg.GetDocID(sourceID, currentMaxSeq)
			//log.Debug(operationID, "msgListToMongo ", seqUid, m.MsgData.Seq, m.MsgData.ClientMsgID, insertCounter, remain, "userID: ", userID)
		} else {
			msgsToMongoNext = append(msgsToMongoNext, sMsg)
			docIDNext = db.msg.GetDocID(sourceID, currentMaxSeq)
			//log.Debug(operationID, "msgListToMongoNext ", seqUidNext, m.MsgData.Seq, m.MsgData.ClientMsgID, insertCounter, remain, "userID: ", userID)
		}
	}

	if docID != "" {
		//filter := bson.M{"uid": seqUid}
		//log.NewDebug(operationID, "filter ", seqUid, "list ", msgListToMongo, "userID: ", userID)
		//err := c.FindOneAndUpdate(ctx, filter, bson.M{"$push": bson.M{"msg": bson.M{"$each": msgsToMongo}}}).Err()
		err = db.mgo.PushMsgsToDoc(ctx, docID, msgsToMongo)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				doc := &unRelationTb.MsgDocModel{}
				doc.DocID = docID
				doc.Msg = msgsToMongo
				if err = db.mgo.Create(ctx, doc); err != nil {
					prome.PromeInc(prome.MsgInsertMongoFailedCounter)
					//log.NewError(operationID, "InsertOne failed", filter, err.Error(), sChat)
					return utils.Wrap(err, "")
				}
				prome.PromeInc(prome.MsgInsertMongoSuccessCounter)
			} else {
				prome.PromeInc(prome.MsgInsertMongoFailedCounter)
				//log.Error(operationID, "FindOneAndUpdate failed ", err.Error(), filter)
				return utils.Wrap(err, "")
			}
		} else {
			prome.PromeInc(prome.MsgInsertMongoSuccessCounter)
		}
	}
	if docIDNext != "" {
		nextDoc := &unRelationTb.MsgDocModel{}
		nextDoc.DocID = docIDNext
		nextDoc.Msg = msgsToMongoNext
		//log.NewDebug(operationID, "filter ", seqUidNext, "list ", msgListToMongoNext, "userID: ", userID)
		if err = db.mgo.Create(ctx, nextDoc); err != nil {
			prome.PromeInc(prome.MsgInsertMongoFailedCounter)
			//log.NewError(operationID, "InsertOne failed", filter, err.Error(), sChat)
			return utils.Wrap(err, "")
		}
		prome.PromeInc(prome.MsgInsertMongoSuccessCounter)
	}
	//log.Debug(operationID, "batch mgo  cost time ", mongo2.getCurrentTimestampByMill()-newTime, userID, len(msgList))
	return nil
}

func (db *MsgDatabase) DeleteMessageFromCache(ctx context.Context, userID string, msgs []*pbMsg.MsgDataToMQ) error {
	return db.cache.DeleteMessageFromCache(ctx, userID, msgs)
}

func (db *MsgDatabase) BatchInsertChat2Cache(ctx context.Context, sourceID string, msgList []*pbMsg.MsgDataToMQ) (int64, error) {
	//newTime := utils.GetCurrentTimestampByMill()
	lenList := len(msgList)
	if int64(lenList) > db.msg.GetSingleGocMsgNum() {
		return 0, errors.New("too large")
	}
	if lenList < 1 {
		return 0, errors.New("too short as 0")
	}
	// judge sessionType to get seq
	var currentMaxSeq int64
	var err error
	if msgList[0].MsgData.SessionType == constant.SuperGroupChatType {
		currentMaxSeq, err = db.cache.GetGroupMaxSeq(ctx, sourceID)
		//log.Debug(operationID, "constant.SuperGroupChatType  lastMaxSeq before add ", currentMaxSeq, "userID ", sourceID, err)
	} else {
		currentMaxSeq, err = db.cache.GetUserMaxSeq(ctx, sourceID)
		//log.Debug(operationID, "constant.SingleChatType  lastMaxSeq before add ", currentMaxSeq, "userID ", sourceID, err)
	}
	if err != nil && err != redis.Nil {
		prome.PromeInc(prome.SeqGetFailedCounter)
		return 0, utils.Wrap(err, "")
	}
	prome.PromeInc(prome.SeqGetSuccessCounter)
	lastMaxSeq := currentMaxSeq
	for _, m := range msgList {
		currentMaxSeq++
		m.MsgData.Seq = currentMaxSeq
		//log.Debug(operationID, "cache msg node ", m.String(), m.MsgData.ClientMsgID, "userID: ", sourceID, "seq: ", currentMaxSeq)
	}
	//log.Debug(operationID, "SetMessageToCache ", sourceID, len(msgList))
	failedNum, err := db.cache.SetMessageToCache(ctx, sourceID, msgList)
	if err != nil {
		prome.PromeAdd(prome.MsgInsertRedisFailedCounter, failedNum)
		//log.Error(operationID, "setMessageToCache failed, continue ", err.Error(), len(msgList), sourceID)
	} else {
		prome.PromeInc(prome.MsgInsertRedisSuccessCounter)
	}
	//log.Debug(operationID, "batch to redis  cost time ", mongo2.getCurrentTimestampByMill()-newTime, sourceID, len(msgList))
	if msgList[0].MsgData.SessionType == constant.SuperGroupChatType {
		err = db.cache.SetGroupMaxSeq(ctx, sourceID, currentMaxSeq)
	} else {
		err = db.cache.SetUserMaxSeq(ctx, sourceID, currentMaxSeq)
	}
	if err != nil {
		prome.PromeInc(prome.SeqSetFailedCounter)
	} else {
		prome.PromeInc(prome.SeqSetSuccessCounter)
	}
	return lastMaxSeq, utils.Wrap(err, "")
}

func (db *MsgDatabase) DelMsgBySeqs(ctx context.Context, userID string, seqs []int64) (totalUnExistSeqs []int64, err error) {
	sortkeys.Int64s(seqs)
	docIDSeqsMap := db.msg.GetDocIDSeqsMap(userID, seqs)
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

func (db *MsgDatabase) DelMsgBySeqsInOneDoc(ctx context.Context, docID string, seqs []int64) (unExistSeqs []int64, err error) {
	seqMsgs, indexes, unExistSeqs, err := db.GetMsgAndIndexBySeqsInOneDoc(ctx, docID, seqs)
	if err != nil {
		return nil, err
	}
	for i, v := range seqMsgs {
		if err = db.mgo.UpdateMsgStatusByIndexInOneDoc(ctx, docID, v, indexes[i], constant.MsgDeleted); err != nil {
			return nil, err
		}
	}
	return unExistSeqs, nil
}

func (db *MsgDatabase) GetMsgAndIndexBySeqsInOneDoc(ctx context.Context, docID string, seqs []int64) (seqMsgs []*sdkws.MsgData, indexes []int, unExistSeqs []int64, err error) {
	doc, err := db.mgo.FindOneByDocID(ctx, docID)
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

func (db *MsgDatabase) GetNewestMsg(ctx context.Context, sourceID string) (msgPb *sdkws.MsgData, err error) {
	msgInfo, err := db.mgo.GetNewestMsg(ctx, sourceID)
	if err != nil {
		return nil, err
	}
	return db.unmarshalMsg(msgInfo)
}

func (db *MsgDatabase) GetOldestMsg(ctx context.Context, sourceID string) (msgPb *sdkws.MsgData, err error) {
	msgInfo, err := db.mgo.GetOldestMsg(ctx, sourceID)
	if err != nil {
		return nil, err
	}
	return db.unmarshalMsg(msgInfo)
}

func (db *MsgDatabase) unmarshalMsg(msgInfo *unRelationTb.MsgInfoModel) (msgPb *sdkws.MsgData, err error) {
	msgPb = &sdkws.MsgData{}
	err = proto.Unmarshal(msgInfo.Msg, msgPb)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	return msgPb, nil
}

func (db *MsgDatabase) getMsgBySeqs(ctx context.Context, sourceID string, seqs []int64, diffusionType int) (seqMsg []*sdkws.MsgData, err error) {
	var hasSeqs []int64
	singleCount := 0
	m := db.msg.GetDocIDSeqsMap(sourceID, seqs)
	for docID, value := range m {
		doc, err := db.mgo.FindOneByDocID(ctx, docID)
		if err != nil {
			//log.NewError(operationID, "not find seqUid", seqUid, value, uid, seqList, err.Error())
			continue
		}
		singleCount = 0
		for i := 0; i < len(doc.Msg); i++ {
			msgPb, err := db.unmarshalMsg(&doc.Msg[i])
			if err != nil {
				//log.NewError(operationID, "Unmarshal err", seqUid, value, uid, seqList, err.Error())
				return nil, err
			}
			if utils.Contain(msgPb.Seq, value...) {
				seqMsg = append(seqMsg, msgPb)
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
		if diffusionType == constant.WriteDiffusion {
			exceptionMsg = db.msg.GenExceptionMessageBySeqs(diff)
		} else if diffusionType == constant.ReadDiffusion {
			exceptionMsg = db.msg.GenExceptionSuperGroupMessageBySeqs(diff, sourceID)
		}
		seqMsg = append(seqMsg, exceptionMsg...)
	}
	return seqMsg, nil
}

func (db *MsgDatabase) GetMsgBySeqs(ctx context.Context, userID string, seqs []int64) (seqMsg []*sdkws.MsgData, err error) {
	successMsgs, failedSeqs, err := db.cache.GetMessageListBySeq(ctx, userID, seqs)
	if err != nil {
		if err != redis.Nil {
			prome.PromeAdd(prome.MsgPullFromRedisFailedCounter, len(failedSeqs))
			log.Error(tracelog.GetOperationID(ctx), "get message from redis exception", err.Error(), failedSeqs)
		}
	}
	prome.PromeAdd(prome.MsgPullFromRedisSuccessCounter, len(successMsgs))
	if len(failedSeqs) > 0 {
		mongoMsgs, err := db.getMsgBySeqs(ctx, userID, seqs, constant.WriteDiffusion)
		if err != nil {
			prome.PromeAdd(prome.MsgPullFromMongoFailedCounter, len(failedSeqs))
			return nil, err
		}
		prome.PromeAdd(prome.MsgPullFromMongoSuccessCounter, len(mongoMsgs))
		successMsgs = append(successMsgs, mongoMsgs...)
	}
	return successMsgs, nil
}

func (db *MsgDatabase) GetSuperGroupMsgBySeqs(ctx context.Context, groupID string, seqs []int64) (seqMsg []*sdkws.MsgData, err error) {
	successMsgs, failedSeqs, err := db.cache.GetMessageListBySeq(ctx, groupID, seqs)
	if err != nil {
		if err != redis.Nil {
			prome.PromeAdd(prome.MsgPullFromRedisFailedCounter, len(failedSeqs))
			log.Error(tracelog.GetOperationID(ctx), "get message from redis exception", err.Error(), failedSeqs)
		}
	}
	prome.PromeAdd(prome.MsgPullFromRedisSuccessCounter, len(successMsgs))
	if len(failedSeqs) > 0 {
		mongoMsgs, err := db.getMsgBySeqs(ctx, groupID, seqs, constant.ReadDiffusion)
		if err != nil {
			prome.PromeAdd(prome.MsgPullFromMongoFailedCounter, len(failedSeqs))
			return nil, err
		}
		prome.PromeAdd(prome.MsgPullFromMongoSuccessCounter, len(mongoMsgs))
		successMsgs = append(successMsgs, mongoMsgs...)
	}
	return successMsgs, nil
}

func (db *MsgDatabase) CleanUpUserMsg(ctx context.Context, userID string) error {
	err := db.DeleteUserMsgsAndSetMinSeq(ctx, userID, 0)
	if err != nil {
		return err
	}
	err = db.cache.CleanUpOneUserAllMsg(ctx, userID)
	return utils.Wrap(err, "")
}

func (db *MsgDatabase) DeleteUserSuperGroupMsgsAndSetMinSeq(ctx context.Context, groupID string, userIDs []string, remainTime int64) error {
	var delStruct delMsgRecursionStruct
	minSeq, err := db.deleteMsgRecursion(ctx, groupID, unRelationTb.OldestList, &delStruct, remainTime)
	if err != nil {
		//log.NewError(operationID, utils.GetSelfFuncName(), groupID, "deleteMsg failed")
	}
	if minSeq == 0 {
		return nil
	}
	//log.NewDebug(operationID, utils.GetSelfFuncName(), "delMsgIDList:", delStruct, "minSeq", minSeq)
	for _, userID := range userIDs {
		userMinSeq, err := db.cache.GetGroupUserMinSeq(ctx, groupID, userID)
		if err != nil && err != redis.Nil {
			//log.NewError(operationID, utils.GetSelfFuncName(), "GetGroupUserMinSeq failed", groupID, userID, err.Error())
			continue
		}
		if userMinSeq > minSeq {
			err = db.cache.SetGroupUserMinSeq(ctx, groupID, userID, userMinSeq)
		} else {
			err = db.cache.SetGroupUserMinSeq(ctx, groupID, userID, minSeq)
		}
		if err != nil {
			//log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), groupID, userID, userMinSeq, minSeq)
		}
	}
	return nil
}

func (db *MsgDatabase) DeleteUserMsgsAndSetMinSeq(ctx context.Context, userID string, remainTime int64) error {
	var delStruct delMsgRecursionStruct
	minSeq, err := db.deleteMsgRecursion(ctx, userID, unRelationTb.OldestList, &delStruct, remainTime)
	if err != nil {
		return utils.Wrap(err, "")
	}
	if minSeq == 0 {
		return nil
	}
	return db.cache.SetUserMinSeq(ctx, userID, minSeq)
}

// this is struct for recursion
type delMsgRecursionStruct struct {
	minSeq       int64
	delDocIDList []string
}

func (d *delMsgRecursionStruct) getSetMinSeq() int64 {
	return d.minSeq
}

// index 0....19(del) 20...69
// seq 70
// set minSeq 21
// recursion 删除list并且返回设置的最小seq
func (db *MsgDatabase) deleteMsgRecursion(ctx context.Context, sourceID string, index int64, delStruct *delMsgRecursionStruct, remainTime int64) (int64, error) {
	// find from oldest list
	msgs, err := db.mgo.GetMsgsByIndex(ctx, sourceID, index)
	if err != nil || msgs.DocID == "" {
		if err != nil {
			if err == unrelation.ErrMsgListNotExist {
				//log.NewInfo(operationID, utils.GetSelfFuncName(), "ID:", sourceID, "index:", index, err.Error())
			} else {
				//log.NewError(operationID, utils.GetSelfFuncName(), "GetUserMsgListByIndex failed", err.Error(), index, ID)
			}
		}
		// 获取报错，或者获取不到了，物理删除并且返回seq delMongoMsgsPhysical(delStruct.delDocIDList)
		err = db.mgo.Delete(ctx, delStruct.delDocIDList)
		if err != nil {
			return 0, err
		}
		return delStruct.getSetMinSeq() + 1, nil
	}
	//log.NewDebug(operationID, "ID:", sourceID, "index:", index, "uid:", msgs.UID, "len:", len(msgs.Msg))
	if int64(len(msgs.Msg)) > db.msg.GetSingleGocMsgNum() {
		log.NewWarn(tracelog.GetOperationID(ctx), utils.GetSelfFuncName(), "msgs too large:", len(msgs.Msg), "docID:", msgs.DocID)
	}
	if msgs.Msg[len(msgs.Msg)-1].SendTime+(remainTime*1000) < utils.GetCurrentTimestampByMill() && msgs.IsFull() {
		delStruct.delDocIDList = append(delStruct.delDocIDList, msgs.DocID)
		lastMsgPb := &sdkws.MsgData{}
		err = proto.Unmarshal(msgs.Msg[len(msgs.Msg)-1].Msg, lastMsgPb)
		if err != nil {
			//log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), len(msgs.Msg)-1, msgs.UID)
			return 0, utils.Wrap(err, "proto.Unmarshal failed")
		}
		delStruct.minSeq = lastMsgPb.Seq
	} else {
		var hasMarkDelFlag bool
		for _, msg := range msgs.Msg {
			msgPb := &sdkws.MsgData{}
			err = proto.Unmarshal(msg.Msg, msgPb)
			if err != nil {
				//log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), len(msgs.Msg)-1, msgs.UID)
				return 0, utils.Wrap(err, "proto.Unmarshal failed")
			}
			if utils.GetCurrentTimestampByMill() > msg.SendTime+(remainTime*1000) {
				msgPb.Status = constant.MsgDeleted
				bytes, _ := proto.Marshal(msgPb)
				msg.Msg = bytes
				msg.SendTime = 0
				hasMarkDelFlag = true
			} else {
				if err := db.mgo.Delete(ctx, delStruct.delDocIDList); err != nil {
					return 0, err
				}
				if hasMarkDelFlag {
					if err := db.mgo.UpdateOneDoc(ctx, msgs); err != nil {
						return delStruct.getSetMinSeq(), utils.Wrap(err, "")
					}
				}
				return msgPb.Seq + 1, nil
			}
		}
	}
	//log.NewDebug(operationID, sourceID, "continue to", delStruct)
	//  继续递归 index+1
	seq, err := db.deleteMsgRecursion(ctx, sourceID, index+1, delStruct, remainTime)
	return seq, utils.Wrap(err, "deleteMsg failed")
}

func (db *MsgDatabase) GetUserMinMaxSeqInMongoAndCache(ctx context.Context, userID string) (minSeqMongo, maxSeqMongo, minSeqCache, maxSeqCache int64, err error) {
	minSeqMongo, maxSeqMongo, err = db.GetMinMaxSeqMongo(ctx, userID)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	// from cache
	minSeqCache, err = db.cache.GetUserMinSeq(ctx, userID)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	maxSeqCache, err = db.cache.GetUserMaxSeq(ctx, userID)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	return
}

func (db *MsgDatabase) GetSuperGroupMinMaxSeqInMongoAndCache(ctx context.Context, groupID string) (minSeqMongo, maxSeqMongo, maxSeqCache int64, err error) {
	minSeqMongo, maxSeqMongo, err = db.GetMinMaxSeqMongo(ctx, groupID)
	if err != nil {
		return 0, 0, 0, err
	}
	maxSeqCache, err = db.cache.GetGroupMaxSeq(ctx, groupID)
	if err != nil {
		return 0, 0, 0, err
	}
	return
}

func (db *MsgDatabase) GetMinMaxSeqMongo(ctx context.Context, sourceID string) (minSeqMongo, maxSeqMongo int64, err error) {
	oldestMsgMongo, err := db.mgo.GetOldestMsg(ctx, sourceID)
	if err != nil {
		return 0, 0, err
	}
	msgPb, err := db.unmarshalMsg(oldestMsgMongo)
	if err != nil {
		return 0, 0, err
	}
	minSeqMongo = msgPb.Seq
	newestMsgMongo, err := db.mgo.GetNewestMsg(ctx, sourceID)
	if err != nil {
		return 0, 0, err
	}
	msgPb, err = db.unmarshalMsg(newestMsgMongo)
	if err != nil {
		return 0, 0, err
	}
	maxSeqMongo = msgPb.Seq
	return
}

func (db *MsgDatabase) SetGroupUserMinSeq(ctx context.Context, groupID, userID string, minSeq int64) (err error) {
	return db.cache.SetGroupUserMinSeq(ctx, groupID, userID, minSeq)
}

func (db *MsgDatabase) SetUserMinSeq(ctx context.Context, userID string, minSeq int64) (err error) {
	return db.cache.SetUserMinSeq(ctx, userID, minSeq)
}
