package db

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	pbMsg "Open_IM/pkg/proto/chat"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (d *DataBases) BatchInsertChat(userID string, msgList []*pbMsg.MsgDataToMQ, operationID string) error {
	newTime := getCurrentTimestampByMill()
	if len(msgList) > GetSingleGocMsgNum() {
		return errors.New("too large")
	}
	isInit := false
	currentMaxSeq, err := d.GetUserMaxSeq(userID)
	if err == nil {

	} else if err == redis.ErrNil {
		isInit = true
		currentMaxSeq = 0
	} else {
		return utils.Wrap(err, "")
	}
	var remain uint64
	if currentMaxSeq < uint64(GetSingleGocMsgNum()) {
		remain = uint64(GetSingleGocMsgNum()-1) - (currentMaxSeq % uint64(GetSingleGocMsgNum()))
	}
	remain = uint64(GetSingleGocMsgNum()) - (currentMaxSeq % uint64(GetSingleGocMsgNum()))
	insertCounter := uint64(0)
	msgListToMongo := make([]MsgInfo, 0)
	msgListToMongoNext := make([]MsgInfo, 0)
	seqUid := ""
	seqUidNext := ""
	log.Debug(operationID, "remain ", remain, "insertCounter ", insertCounter, "currentMaxSeq ", currentMaxSeq, userID, len(msgList))
	for _, m := range msgList {
		log.Debug(operationID, "msg node ", m.String(), m.MsgData.ClientMsgID)
		currentMaxSeq++
		sMsg := MsgInfo{}
		sMsg.SendTime = m.MsgData.SendTime
		m.MsgData.Seq = uint32(currentMaxSeq)
		if sMsg.Msg, err = proto.Marshal(m.MsgData); err != nil {
			return utils.Wrap(err, "")
		}
		if isInit {
			msgListToMongoNext = append(msgListToMongoNext, sMsg)
			seqUidNext = getSeqUid(userID, uint32(currentMaxSeq))
			log.Debug(operationID, "msgListToMongoNext ", seqUidNext, m.MsgData.Seq, m.MsgData.ClientMsgID, insertCounter, remain)
			continue
		}

		if insertCounter < remain {
			msgListToMongo = append(msgListToMongo, sMsg)
			insertCounter++
			seqUid = getSeqUid(userID, uint32(currentMaxSeq))
			log.Debug(operationID, "msgListToMongo ", seqUid, m.MsgData.Seq, m.MsgData.ClientMsgID, insertCounter, remain)
		} else {
			msgListToMongoNext = append(msgListToMongoNext, sMsg)
			seqUidNext = getSeqUid(userID, uint32(currentMaxSeq))
			log.Debug(operationID, "msgListToMongoNext ", seqUidNext, m.MsgData.Seq, m.MsgData.ClientMsgID, insertCounter, remain)
		}
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cChat)

	if seqUid != "" {
		filter := bson.M{"uid": seqUid}
		log.NewDebug(operationID, "filter ", seqUid, "list ", msgListToMongo)
		err := c.FindOneAndUpdate(ctx, filter, bson.M{"$push": bson.M{"msg": bson.M{"$each": msgListToMongo}}}).Err()
		if err != nil {
			log.Error(operationID, "FindOneAndUpdate failed ", err.Error(), filter)
			return utils.Wrap(err, "")
		}
	}
	if seqUidNext != "" {
		filter := bson.M{"uid": seqUidNext}
		sChat := UserChat{}
		sChat.UID = seqUidNext
		sChat.Msg = msgListToMongoNext
		log.NewDebug(operationID, "filter ", seqUidNext, "list ", msgListToMongoNext)
		if _, err = c.InsertOne(ctx, &sChat); err != nil {
			log.NewError(operationID, "InsertOne failed", filter, err.Error(), sChat)
			return utils.Wrap(err, "")
		}
	}
	log.NewWarn(operationID, "batch mgo  cost time ", getCurrentTimestampByMill()-newTime, userID, len(msgList))
	return utils.Wrap(d.SetUserMaxSeq(userID, uint64(currentMaxSeq)), "")
}
