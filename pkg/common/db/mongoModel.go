package db

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	pbMsg "Open_IM/pkg/proto/msg"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gogo/protobuf/sortkeys"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
	"sync"

	//"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
	"go.mongodb.org/mongo-driver/bson"

	"strconv"
	"time"
)

const cChat = "msg"
const cGroup = "group"
const cTag = "tag"
const cSendLog = "send_log"
const cWorkMoment = "work_moment"
const cCommentMsg = "comment_msg"
const cSuperGroup = "super_group"
const cUserToSuperGroup = "user_to_super_group"
const singleGocMsgNum = 5000

func GetSingleGocMsgNum() int {
	return singleGocMsgNum
}

type MsgInfo struct {
	SendTime int64
	Msg      []byte
}

type UserChat struct {
	UID string
	Msg []MsgInfo
}

type GroupMember_x struct {
	GroupID string
	UIDList []string
}

func (d *DataBases) GetMinSeqFromMongo(uid string) (MinSeq uint32, err error) {
	return 1, nil
	//var i, NB uint32
	//var seqUid string
	//session := d.mgoSession.Clone()
	//if session == nil {
	//	return MinSeq, errors.New("session == nil")
	//}
	//defer session.Close()
	//c := session.DB(config.Config.Mongo.DBDatabase).C(cChat)
	//MaxSeq, err := d.GetUserMaxSeq(uid)
	//if err != nil && err != redis.ErrNil {
	//	return MinSeq, err
	//}
	//NB = uint32(MaxSeq / singleGocMsgNum)
	//for i = 0; i <= NB; i++ {
	//	seqUid = indexGen(uid, i)
	//	n, err := c.Find(bson.M{"uid": seqUid}).Count()
	//	if err == nil && n != 0 {
	//		if i == 0 {
	//			MinSeq = 1
	//		} else {
	//			MinSeq = uint32(i * singleGocMsgNum)
	//		}
	//		break
	//	}
	//}
	//return MinSeq, nil
}

func (d *DataBases) GetMinSeqFromMongo2(uid string) (MinSeq uint32, err error) {
	return 1, nil
}

// deleteMsgByLogic
func (d *DataBases) DelMsgBySeqList(userID string, seqList []uint32, operationID string) (totalUnexistSeqList []uint32, err error) {
	log.Debug(operationID, utils.GetSelfFuncName(), "args ", userID, seqList)
	sortkeys.Uint32s(seqList)
	suffixUserID2SubSeqList := func(uid string, seqList []uint32) map[string][]uint32 {
		t := make(map[string][]uint32)
		for i := 0; i < len(seqList); i++ {
			seqUid := getSeqUid(uid, seqList[i])
			if value, ok := t[seqUid]; !ok {
				var temp []uint32
				t[seqUid] = append(temp, seqList[i])
			} else {
				t[seqUid] = append(value, seqList[i])
			}
		}
		return t
	}(userID, seqList)

	lock := sync.Mutex{}
	var wg sync.WaitGroup
	wg.Add(len(suffixUserID2SubSeqList))
	for k, v := range suffixUserID2SubSeqList {
		go func(suffixUserID string, subSeqList []uint32, operationID string) {
			defer wg.Done()
			unexistSeqList, err := d.DelMsgBySeqListInOneDoc(suffixUserID, subSeqList, operationID)
			if err != nil {
				log.Error(operationID, "DelMsgBySeqListInOneDoc failed ", err.Error(), suffixUserID, subSeqList)
				return
			}
			lock.Lock()
			totalUnexistSeqList = append(totalUnexistSeqList, unexistSeqList...)
			lock.Unlock()
		}(k, v, operationID)
	}
	return totalUnexistSeqList, err
}

func (d *DataBases) DelMsgBySeqListInOneDoc(suffixUserID string, seqList []uint32, operationID string) ([]uint32, error) {
	log.Debug(operationID, utils.GetSelfFuncName(), "args ", suffixUserID, seqList)
	seqMsgList, indexList, unexistSeqList, err := d.GetMsgAndIndexBySeqListInOneMongo2(suffixUserID, seqList, operationID)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	for i, v := range seqMsgList {
		if err := d.ReplaceMsgByIndex(suffixUserID, v, operationID, indexList[i]); err != nil {
			return nil, utils.Wrap(err, "")
		}
	}
	return unexistSeqList, nil
}

// deleteMsgByLogic
func (d *DataBases) DelMsgLogic(uid string, seqList []uint32, operationID string) error {
	sortkeys.Uint32s(seqList)
	seqMsgs, err := d.GetMsgBySeqListMongo2(uid, seqList, operationID)
	if err != nil {
		return utils.Wrap(err, "")
	}
	for _, seqMsg := range seqMsgs {
		log.NewDebug(operationID, utils.GetSelfFuncName(), *seqMsg)
		seqMsg.Status = constant.MsgDeleted
		if err = d.ReplaceMsgBySeq(uid, seqMsg, operationID); err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), "ReplaceMsgListBySeq error", err.Error())
		}
	}
	return nil
}

func (d *DataBases) ReplaceMsgByIndex(suffixUserID string, msg *open_im_sdk.MsgData, operationID string, seqIndex int) error {
	log.NewInfo(operationID, utils.GetSelfFuncName(), suffixUserID, *msg)
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cChat)
	s := fmt.Sprintf("msg.%d.msg", seqIndex)
	log.NewDebug(operationID, utils.GetSelfFuncName(), seqIndex, s)
	msg.Status = constant.MsgDeleted
	bytes, err := proto.Marshal(msg)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "proto marshal failed ", err.Error(), msg.String())
		return utils.Wrap(err, "")
	}
	updateResult, err := c.UpdateOne(ctx, bson.M{"uid": suffixUserID}, bson.M{"$set": bson.M{s: bytes}})
	log.NewInfo(operationID, utils.GetSelfFuncName(), updateResult)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "UpdateOne", err.Error())
		return utils.Wrap(err, "")
	}
	return nil
}

func (d *DataBases) ReplaceMsgBySeq(uid string, msg *open_im_sdk.MsgData, operationID string) error {
	log.NewInfo(operationID, utils.GetSelfFuncName(), uid, *msg)
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cChat)
	uid = getSeqUid(uid, msg.Seq)
	seqIndex := getMsgIndex(msg.Seq)
	s := fmt.Sprintf("msg.%d.msg", seqIndex)
	log.NewDebug(operationID, utils.GetSelfFuncName(), seqIndex, s)
	bytes, err := proto.Marshal(msg)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "proto marshal", err.Error())
		return utils.Wrap(err, "")
	}

	updateResult, err := c.UpdateOne(
		ctx, bson.M{"uid": uid},
		bson.M{"$set": bson.M{s: bytes}})
	log.NewInfo(operationID, utils.GetSelfFuncName(), updateResult)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "UpdateOne", err.Error())
		return utils.Wrap(err, "")
	}
	return nil
}

func (d *DataBases) GetMsgBySeqList(uid string, seqList []uint32, operationID string) (seqMsg []*open_im_sdk.MsgData, err error) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), uid, seqList)
	var hasSeqList []uint32
	singleCount := 0
	session := d.mgoSession.Clone()
	if session == nil {
		return nil, errors.New("session == nil")
	}
	defer session.Close()
	c := session.DB(config.Config.Mongo.DBDatabase).C(cChat)
	m := func(uid string, seqList []uint32) map[string][]uint32 {
		t := make(map[string][]uint32)
		for i := 0; i < len(seqList); i++ {
			seqUid := getSeqUid(uid, seqList[i])
			if value, ok := t[seqUid]; !ok {
				var temp []uint32
				t[seqUid] = append(temp, seqList[i])
			} else {
				t[seqUid] = append(value, seqList[i])
			}
		}
		return t
	}(uid, seqList)
	sChat := UserChat{}
	for seqUid, value := range m {
		if err = c.Find(bson.M{"uid": seqUid}).One(&sChat); err != nil {
			log.NewError(operationID, "not find seqUid", seqUid, value, uid, seqList, err.Error())
			continue
		}
		singleCount = 0
		for i := 0; i < len(sChat.Msg); i++ {
			msg := new(open_im_sdk.MsgData)
			if err = proto.Unmarshal(sChat.Msg[i].Msg, msg); err != nil {
				log.NewError(operationID, "Unmarshal err", seqUid, value, uid, seqList, err.Error())
				return nil, err
			}
			if isContainInt32(msg.Seq, value) {
				seqMsg = append(seqMsg, msg)
				hasSeqList = append(hasSeqList, msg.Seq)
				singleCount++
				if singleCount == len(value) {
					break
				}
			}
		}
	}
	if len(hasSeqList) != len(seqList) {
		var diff []uint32
		diff = utils.Difference(hasSeqList, seqList)
		exceptionMSg := genExceptionMessageBySeqList(diff)
		seqMsg = append(seqMsg, exceptionMSg...)

	}
	return seqMsg, nil
}

func (d *DataBases) GetMsgBySeqListMongo2(uid string, seqList []uint32, operationID string) (seqMsg []*open_im_sdk.MsgData, err error) {
	var hasSeqList []uint32
	singleCount := 0
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cChat)

	m := func(uid string, seqList []uint32) map[string][]uint32 {
		t := make(map[string][]uint32)
		for i := 0; i < len(seqList); i++ {
			seqUid := getSeqUid(uid, seqList[i])
			if value, ok := t[seqUid]; !ok {
				var temp []uint32
				t[seqUid] = append(temp, seqList[i])
			} else {
				t[seqUid] = append(value, seqList[i])
			}
		}
		return t
	}(uid, seqList)
	sChat := UserChat{}
	for seqUid, value := range m {
		if err = c.FindOne(ctx, bson.M{"uid": seqUid}).Decode(&sChat); err != nil {
			log.NewError(operationID, "not find seqUid", seqUid, value, uid, seqList, err.Error())
			continue
		}
		singleCount = 0
		for i := 0; i < len(sChat.Msg); i++ {
			msg := new(open_im_sdk.MsgData)
			if err = proto.Unmarshal(sChat.Msg[i].Msg, msg); err != nil {
				log.NewError(operationID, "Unmarshal err", seqUid, value, uid, seqList, err.Error())
				return nil, err
			}
			if isContainInt32(msg.Seq, value) {
				seqMsg = append(seqMsg, msg)
				hasSeqList = append(hasSeqList, msg.Seq)
				singleCount++
				if singleCount == len(value) {
					break
				}
			}
		}
	}
	if len(hasSeqList) != len(seqList) {
		var diff []uint32
		diff = utils.Difference(hasSeqList, seqList)
		exceptionMSg := genExceptionMessageBySeqList(diff)
		seqMsg = append(seqMsg, exceptionMSg...)

	}
	return seqMsg, nil
}
func (d *DataBases) GetSuperGroupMsgBySeqListMongo(groupID string, seqList []uint32, operationID string) (seqMsg []*open_im_sdk.MsgData, err error) {
	var hasSeqList []uint32
	singleCount := 0
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cChat)

	m := func(uid string, seqList []uint32) map[string][]uint32 {
		t := make(map[string][]uint32)
		for i := 0; i < len(seqList); i++ {
			seqUid := getSeqUid(uid, seqList[i])
			if value, ok := t[seqUid]; !ok {
				var temp []uint32
				t[seqUid] = append(temp, seqList[i])
			} else {
				t[seqUid] = append(value, seqList[i])
			}
		}
		return t
	}(groupID, seqList)
	sChat := UserChat{}
	for seqUid, value := range m {
		if err = c.FindOne(ctx, bson.M{"uid": seqUid}).Decode(&sChat); err != nil {
			log.NewError(operationID, "not find seqGroupID", seqUid, value, groupID, seqList, err.Error())
			continue
		}
		singleCount = 0
		for i := 0; i < len(sChat.Msg); i++ {
			msg := new(open_im_sdk.MsgData)
			if err = proto.Unmarshal(sChat.Msg[i].Msg, msg); err != nil {
				log.NewError(operationID, "Unmarshal err", seqUid, value, groupID, seqList, err.Error())
				return nil, err
			}
			if isContainInt32(msg.Seq, value) {
				seqMsg = append(seqMsg, msg)
				hasSeqList = append(hasSeqList, msg.Seq)
				singleCount++
				if singleCount == len(value) {
					break
				}
			}
		}
	}
	if len(hasSeqList) != len(seqList) {
		var diff []uint32
		diff = utils.Difference(hasSeqList, seqList)
		exceptionMSg := genExceptionSuperGroupMessageBySeqList(diff, groupID)
		seqMsg = append(seqMsg, exceptionMSg...)

	}
	return seqMsg, nil
}

func (d *DataBases) GetMsgAndIndexBySeqListInOneMongo2(suffixUserID string, seqList []uint32, operationID string) (seqMsg []*open_im_sdk.MsgData, indexList []int, unexistSeqList []uint32, err error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cChat)
	sChat := UserChat{}
	if err = c.FindOne(ctx, bson.M{"uid": suffixUserID}).Decode(&sChat); err != nil {
		log.NewError(operationID, "not find seqUid", suffixUserID, err.Error())
		return nil, nil, nil, utils.Wrap(err, "")
	}
	singleCount := 0
	var hasSeqList []uint32
	for i := 0; i < len(sChat.Msg); i++ {
		msg := new(open_im_sdk.MsgData)
		if err = proto.Unmarshal(sChat.Msg[i].Msg, msg); err != nil {
			log.NewError(operationID, "Unmarshal err", msg.String(), err.Error())
			return nil, nil, nil, err
		}
		if isContainInt32(msg.Seq, seqList) {
			indexList = append(indexList, i)
			seqMsg = append(seqMsg, msg)
			hasSeqList = append(hasSeqList, msg.Seq)
			singleCount++
			if singleCount == len(seqList) {
				break
			}
		}
	}
	for _, i := range seqList {
		if isContainInt32(i, hasSeqList) {
			continue
		}
		unexistSeqList = append(unexistSeqList, i)
	}
	return seqMsg, indexList, unexistSeqList, nil
}

func genExceptionMessageBySeqList(seqList []uint32) (exceptionMsg []*open_im_sdk.MsgData) {
	for _, v := range seqList {
		msg := new(open_im_sdk.MsgData)
		msg.Seq = v
		exceptionMsg = append(exceptionMsg, msg)
	}
	return exceptionMsg
}

func genExceptionSuperGroupMessageBySeqList(seqList []uint32, groupID string) (exceptionMsg []*open_im_sdk.MsgData) {
	for _, v := range seqList {
		msg := new(open_im_sdk.MsgData)
		msg.Seq = v
		msg.GroupID = groupID
		msg.SessionType = constant.SuperGroupChatType
		exceptionMsg = append(exceptionMsg, msg)
	}
	return exceptionMsg
}

func (d *DataBases) SaveUserChatMongo2(uid string, sendTime int64, m *pbMsg.MsgDataToDB) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cChat)
	newTime := getCurrentTimestampByMill()
	operationID := ""
	seqUid := getSeqUid(uid, m.MsgData.Seq)
	filter := bson.M{"uid": seqUid}
	var err error
	sMsg := MsgInfo{}
	sMsg.SendTime = sendTime
	if sMsg.Msg, err = proto.Marshal(m.MsgData); err != nil {
		return utils.Wrap(err, "")
	}
	err = c.FindOneAndUpdate(ctx, filter, bson.M{"$push": bson.M{"msg": sMsg}}).Err()
	log.NewWarn(operationID, "get mgoSession cost time", getCurrentTimestampByMill()-newTime)
	if err != nil {
		sChat := UserChat{}
		sChat.UID = seqUid
		sChat.Msg = append(sChat.Msg, sMsg)
		if _, err = c.InsertOne(ctx, &sChat); err != nil {
			log.NewDebug(operationID, "InsertOne failed", filter)
			return utils.Wrap(err, "")
		}
	} else {
		log.NewDebug(operationID, "FindOneAndUpdate ok", filter)
	}

	log.NewDebug(operationID, "find mgo uid cost time", getCurrentTimestampByMill()-newTime)
	return nil
}

//
//func (d *DataBases) SaveUserChatListMongo2(uid string, sendTime int64, msgList []*pbMsg.MsgDataToDB) error {
//	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
//	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cChat)
//	newTime := getCurrentTimestampByMill()
//	operationID := ""
//	seqUid := ""
//	msgListToMongo := make([]MsgInfo, 0)
//
//	for _, m := range msgList {
//		seqUid = getSeqUid(uid, m.MsgData.Seq)
//		var err error
//		sMsg := MsgInfo{}
//		sMsg.SendTime = sendTime
//		if sMsg.Msg, err = proto.Marshal(m.MsgData); err != nil {
//			return utils.Wrap(err, "")
//		}
//		msgListToMongo = append(msgListToMongo, sMsg)
//	}
//
//	filter := bson.M{"uid": seqUid}
//	log.NewDebug(operationID, "filter ", seqUid)
//	err := c.FindOneAndUpdate(ctx, filter, bson.M{"$push": bson.M{"msg": bson.M{"$each": msgListToMongo}}}).Err()
//	log.NewWarn(operationID, "get mgoSession cost time", getCurrentTimestampByMill()-newTime)
//	if err != nil {
//		sChat := UserChat{}
//		sChat.UID = seqUid
//		sChat.Msg = msgListToMongo
//
//		if _, err = c.InsertOne(ctx, &sChat); err != nil {
//			log.NewError(operationID, "InsertOne failed", filter, err.Error(), sChat)
//			return utils.Wrap(err, "")
//		}
//	} else {
//		log.NewDebug(operationID, "FindOneAndUpdate ok", filter)
//	}
//
//	log.NewDebug(operationID, "find mgo uid cost time", getCurrentTimestampByMill()-newTime)
//	return nil
//}

func (d *DataBases) SaveUserChat(uid string, sendTime int64, m *pbMsg.MsgDataToDB) error {
	var seqUid string
	newTime := getCurrentTimestampByMill()
	session := d.mgoSession.Clone()
	if session == nil {
		return errors.New("session == nil")
	}
	defer session.Close()
	log.NewDebug("", "get mgoSession cost time", getCurrentTimestampByMill()-newTime)
	c := session.DB(config.Config.Mongo.DBDatabase).C(cChat)
	seqUid = getSeqUid(uid, m.MsgData.Seq)
	n, err := c.Find(bson.M{"uid": seqUid}).Count()
	if err != nil {
		return err
	}
	log.NewDebug("", "find mgo uid cost time", getCurrentTimestampByMill()-newTime)
	sMsg := MsgInfo{}
	sMsg.SendTime = sendTime
	if sMsg.Msg, err = proto.Marshal(m.MsgData); err != nil {
		return err
	}
	if n == 0 {
		sChat := UserChat{}
		sChat.UID = seqUid
		sChat.Msg = append(sChat.Msg, sMsg)
		err = c.Insert(&sChat)
		if err != nil {
			return err
		}
	} else {
		err = c.Update(bson.M{"uid": seqUid}, bson.M{"$push": bson.M{"msg": sMsg}})
		if err != nil {
			return err
		}
	}
	log.NewDebug("", "insert mgo data cost time", getCurrentTimestampByMill()-newTime)
	return nil
}

func (d *DataBases) DelUserChat(uid string) error {
	return nil
	//session := d.mgoSession.Clone()
	//if session == nil {
	//	return errors.New("session == nil")
	//}
	//defer session.Close()
	//
	//c := session.DB(config.Config.Mongo.DBDatabase).C(cChat)
	//
	//delTime := time.Now().Unix() - int64(config.Config.Mongo.DBRetainChatRecords)*24*3600
	//if err := c.Update(bson.M{"uid": uid}, bson.M{"$pull": bson.M{"msg": bson.M{"sendtime": bson.M{"$lte": delTime}}}}); err != nil {
	//	return err
	//}
	//
	//return nil
}

func (d *DataBases) DelUserChatMongo2(uid string) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cChat)
	filter := bson.M{"uid": uid}

	delTime := time.Now().Unix() - int64(config.Config.Mongo.DBRetainChatRecords)*24*3600
	if _, err := c.UpdateOne(ctx, filter, bson.M{"$pull": bson.M{"msg": bson.M{"sendtime": bson.M{"$lte": delTime}}}}); err != nil {
		return utils.Wrap(err, "")
	}
	return nil
}

func (d *DataBases) MgoUserCount() (int, error) {
	return 0, nil
	//session := d.mgoSession.Clone()
	//if session == nil {
	//	return 0, errors.New("session == nil")
	//}
	//defer session.Close()
	//
	//c := session.DB(config.Config.Mongo.DBDatabase).C(cChat)
	//
	//return c.Find(nil).Count()
}

func (d *DataBases) MgoSkipUID(count int) (string, error) {
	return "", nil
	//session := d.mgoSession.Clone()
	//if session == nil {
	//	return "", errors.New("session == nil")
	//}
	//defer session.Close()
	//
	//c := session.DB(config.Config.Mongo.DBDatabase).C(cChat)
	//
	//sChat := UserChat{}
	//c.Find(nil).Skip(count).Limit(1).One(&sChat)
	//return sChat.UID, nil
}

func (d *DataBases) GetGroupMember(groupID string) []string {
	return nil
	//groupInfo := GroupMember_x{}
	//groupInfo.GroupID = groupID
	//groupInfo.UIDList = make([]string, 0)
	//
	//session := d.mgoSession.Clone()
	//if session == nil {
	//	return groupInfo.UIDList
	//}
	//defer session.Close()
	//
	//c := session.DB(config.Config.Mongo.DBDatabase).C(cGroup)
	//
	//if err := c.Find(bson.M{"groupid": groupInfo.GroupID}).One(&groupInfo); err != nil {
	//	return groupInfo.UIDList
	//}
	//
	//return groupInfo.UIDList
}

func (d *DataBases) AddGroupMember(groupID, uid string) error {
	return nil
	//session := d.mgoSession.Clone()
	//if session == nil {
	//	return errors.New("session == nil")
	//}
	//defer session.Close()
	//
	//c := session.DB(config.Config.Mongo.DBDatabase).C(cGroup)
	//
	//n, err := c.Find(bson.M{"groupid": groupID}).Count()
	//if err != nil {
	//	return err
	//}
	//
	//if n == 0 {
	//	groupInfo := GroupMember_x{}
	//	groupInfo.GroupID = groupID
	//	groupInfo.UIDList = append(groupInfo.UIDList, uid)
	//	err = c.Insert(&groupInfo)
	//	if err != nil {
	//		return err
	//	}
	//} else {
	//	err = c.Update(bson.M{"groupid": groupID}, bson.M{"$addToSet": bson.M{"uidlist": uid}})
	//	if err != nil {
	//		return err
	//	}
	//}
	//
	//return nil
}

func (d *DataBases) DelGroupMember(groupID, uid string) error {
	return nil
	//session := d.mgoSession.Clone()
	//if session == nil {
	//	return errors.New("session == nil")
	//}
	//defer session.Close()
	//
	//c := session.DB(config.Config.Mongo.DBDatabase).C(cGroup)
	//
	//if err := c.Update(bson.M{"groupid": groupID}, bson.M{"$pull": bson.M{"uidlist": uid}}); err != nil {
	//	return err
	//}
	//
	//return nil
}

type Tag struct {
	UserID   string   `bson:"user_id"`
	TagID    string   `bson:"tag_id"`
	TagName  string   `bson:"tag_name"`
	UserList []string `bson:"user_list"`
}

func (d *DataBases) GetUserTags(userID string) ([]Tag, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cTag)
	var tags []Tag
	cursor, err := c.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return tags, err
	}
	if err = cursor.All(ctx, &tags); err != nil {
		return tags, err
	}
	return tags, nil
}

func (d *DataBases) CreateTag(userID, tagName string, userList []string) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cTag)
	tagID := generateTagID(tagName, userID)
	tag := Tag{
		UserID:   userID,
		TagID:    tagID,
		TagName:  tagName,
		UserList: userList,
	}
	_, err := c.InsertOne(ctx, tag)
	return err
}

func (d *DataBases) GetTagByID(userID, tagID string) (Tag, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cTag)
	var tag Tag
	err := c.FindOne(ctx, bson.M{"user_id": userID, "tag_id": tagID}).Decode(&tag)
	return tag, err
}

func (d *DataBases) DeleteTag(userID, tagID string) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cTag)
	_, err := c.DeleteOne(ctx, bson.M{"user_id": userID, "tag_id": tagID})
	return err
}

func (d *DataBases) SetTag(userID, tagID, newName string, increaseUserIDList []string, reduceUserIDList []string) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cTag)
	var tag Tag
	if err := c.FindOne(ctx, bson.M{"tag_id": tagID, "user_id": userID}).Decode(&tag); err != nil {
		return err
	}
	if newName != "" {
		_, err := c.UpdateOne(ctx, bson.M{"user_id": userID, "tag_id": tagID}, bson.M{"$set": bson.M{"tag_name": newName}})
		if err != nil {
			return err
		}
	}
	tag.UserList = append(tag.UserList, increaseUserIDList...)
	tag.UserList = utils.RemoveRepeatedStringInList(tag.UserList)
	for _, v := range reduceUserIDList {
		for i2, v2 := range tag.UserList {
			if v == v2 {
				tag.UserList[i2] = ""
			}
		}
	}
	var newUserList []string
	for _, v := range tag.UserList {
		if v != "" {
			newUserList = append(newUserList, v)
		}
	}
	_, err := c.UpdateOne(ctx, bson.M{"user_id": userID, "tag_id": tagID}, bson.M{"$set": bson.M{"user_list": newUserList}})
	if err != nil {
		return err
	}
	return nil
}

func (d *DataBases) GetUserIDListByTagID(userID, tagID string) ([]string, error) {
	var tag Tag
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cTag)
	_ = c.FindOne(ctx, bson.M{"user_id": userID, "tag_id": tagID}).Decode(&tag)
	return tag.UserList, nil
}

type TagUser struct {
	UserID   string `bson:"user_id"`
	UserName string `bson:"user_name"`
}

type TagSendLog struct {
	UserList         []TagUser `bson:"tag_list"`
	SendID           string    `bson:"send_id"`
	SenderPlatformID int32     `bson:"sender_platform_id"`
	Content          string    `bson:"content"`
	SendTime         int64     `bson:"send_time"`
}

func (d *DataBases) SaveTagSendLog(tagSendLog *TagSendLog) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cSendLog)
	_, err := c.InsertOne(ctx, tagSendLog)
	return err
}

func (d *DataBases) GetTagSendLogs(userID string, showNumber, pageNumber int32) ([]TagSendLog, error) {
	var tagSendLogs []TagSendLog
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cSendLog)
	findOpts := options.Find().SetLimit(int64(showNumber)).SetSkip(int64(showNumber) * (int64(pageNumber) - 1)).SetSort(bson.M{"send_time": -1})
	cursor, err := c.Find(ctx, bson.M{"send_id": userID}, findOpts)
	if err != nil {
		return tagSendLogs, err
	}
	err = cursor.All(ctx, &tagSendLogs)
	if err != nil {
		return tagSendLogs, err
	}
	return tagSendLogs, nil
}

type WorkMoment struct {
	WorkMomentID         string            `bson:"work_moment_id"`
	UserID               string            `bson:"user_id"`
	UserName             string            `bson:"user_name"`
	FaceURL              string            `bson:"face_url"`
	Content              string            `bson:"content"`
	LikeUserList         []*WorkMomentUser `bson:"like_user_list"`
	AtUserList           []*WorkMomentUser `bson:"at_user_list"`
	PermissionUserList   []*WorkMomentUser `bson:"permission_user_list"`
	Comments             []*Comment        `bson:"comments"`
	PermissionUserIDList []string          `bson:"permission_user_id_list"`
	Permission           int32             `bson:"permission"`
	CreateTime           int32             `bson:"create_time"`
}

type WorkMomentUser struct {
	UserID   string `bson:"user_id"`
	UserName string `bson:"user_name"`
}

type Comment struct {
	UserID        string `bson:"user_id" json:"user_id"`
	UserName      string `bson:"user_name" json:"user_name"`
	ReplyUserID   string `bson:"reply_user_id" json:"reply_user_id"`
	ReplyUserName string `bson:"reply_user_name" json:"reply_user_name"`
	ContentID     string `bson:"content_id" json:"content_id"`
	Content       string `bson:"content" json:"content"`
	CreateTime    int32  `bson:"create_time" json:"create_time"`
}

func (d *DataBases) CreateOneWorkMoment(workMoment *WorkMoment) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cWorkMoment)
	workMomentID := generateWorkMomentID(workMoment.UserID)
	workMoment.WorkMomentID = workMomentID
	workMoment.CreateTime = int32(time.Now().Unix())
	_, err := c.InsertOne(ctx, workMoment)
	return err
}

func (d *DataBases) DeleteOneWorkMoment(workMomentID string) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cWorkMoment)
	_, err := c.DeleteOne(ctx, bson.M{"work_moment_id": workMomentID})
	return err
}

func (d *DataBases) DeleteComment(workMomentID, contentID, opUserID string) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cWorkMoment)
	_, err := c.UpdateOne(ctx, bson.D{{"work_moment_id", workMomentID},
		{"$or", bson.A{
			bson.D{{"user_id", opUserID}},
			bson.D{{"comments", bson.M{"$elemMatch": bson.M{"user_id": opUserID}}}},
		},
		}}, bson.M{"$pull": bson.M{"comments": bson.M{"content_id": contentID}}})
	return err
}

func (d *DataBases) GetWorkMomentByID(workMomentID string) (*WorkMoment, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cWorkMoment)
	workMoment := &WorkMoment{}
	err := c.FindOne(ctx, bson.M{"work_moment_id": workMomentID}).Decode(workMoment)
	return workMoment, err
}

func (d *DataBases) LikeOneWorkMoment(likeUserID, userName, workMomentID string) (*WorkMoment, bool, error) {
	workMoment, err := d.GetWorkMomentByID(workMomentID)
	if err != nil {
		return nil, false, err
	}
	var isAlreadyLike bool
	for i, user := range workMoment.LikeUserList {
		if likeUserID == user.UserID {
			isAlreadyLike = true
			workMoment.LikeUserList = append(workMoment.LikeUserList[0:i], workMoment.LikeUserList[i+1:]...)
		}
	}
	if !isAlreadyLike {
		workMoment.LikeUserList = append(workMoment.LikeUserList, &WorkMomentUser{UserID: likeUserID, UserName: userName})
	}
	log.NewDebug("", utils.GetSelfFuncName(), workMoment)
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cWorkMoment)
	_, err = c.UpdateOne(ctx, bson.M{"work_moment_id": workMomentID}, bson.M{"$set": bson.M{"like_user_list": workMoment.LikeUserList}})
	return workMoment, !isAlreadyLike, err
}

func (d *DataBases) SetUserWorkMomentsLevel(userID string, level int32) error {
	return nil
}

func (d *DataBases) CommentOneWorkMoment(comment *Comment, workMomentID string) (WorkMoment, error) {
	comment.ContentID = generateWorkMomentCommentID(workMomentID)
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cWorkMoment)
	var workMoment WorkMoment
	err := c.FindOneAndUpdate(ctx, bson.M{"work_moment_id": workMomentID}, bson.M{"$push": bson.M{"comments": comment}}).Decode(&workMoment)
	return workMoment, err
}

func (d *DataBases) GetUserSelfWorkMoments(userID string, showNumber, pageNumber int32) ([]WorkMoment, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cWorkMoment)
	var workMomentList []WorkMoment
	findOpts := options.Find().SetLimit(int64(showNumber)).SetSkip(int64(showNumber) * (int64(pageNumber) - 1)).SetSort(bson.M{"create_time": -1})
	result, err := c.Find(ctx, bson.M{"user_id": userID}, findOpts)
	if err != nil {
		return workMomentList, nil
	}
	err = result.All(ctx, &workMomentList)
	return workMomentList, err
}

func (d *DataBases) GetUserWorkMoments(opUserID, userID string, showNumber, pageNumber int32) ([]WorkMoment, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cWorkMoment)
	var workMomentList []WorkMoment
	findOpts := options.Find().SetLimit(int64(showNumber)).SetSkip(int64(showNumber) * (int64(pageNumber) - 1)).SetSort(bson.M{"create_time": -1})
	result, err := c.Find(ctx, bson.D{ // 等价条件: select * from
		{"user_id", userID},
		{"$or", bson.A{
			bson.D{{"permission", constant.WorkMomentPermissionCantSee}, {"permission_user_id_list", bson.D{{"$nin", bson.A{opUserID}}}}},
			bson.D{{"permission", constant.WorkMomentPermissionCanSee}, {"permission_user_id_list", bson.D{{"$in", bson.A{opUserID}}}}},
			bson.D{{"permission", constant.WorkMomentPublic}},
		}},
	}, findOpts)
	if err != nil {
		return workMomentList, nil
	}
	err = result.All(ctx, &workMomentList)
	return workMomentList, err
}

func (d *DataBases) GetUserFriendWorkMoments(showNumber, pageNumber int32, userID string) ([]WorkMoment, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cWorkMoment)
	var workMomentList []WorkMoment
	findOpts := options.Find().SetLimit(int64(showNumber)).SetSkip(int64(showNumber) * (int64(pageNumber) - 1)).SetSort(bson.M{"create_time": -1})
	result, err := c.Find(ctx, bson.D{
		{"$or", bson.A{
			bson.D{{"user_id", userID}}, //self
			bson.D{
				{"$or", bson.A{
					bson.D{{"permission", constant.WorkMomentPermissionCantSee}, {"permission_user_id_list", bson.D{{"$nin", bson.A{userID}}}}},
					bson.D{{"permission", constant.WorkMomentPermissionCanSee}, {"permission_user_id_list", bson.D{{"$in", bson.A{userID}}}}},
					bson.D{{"permission", constant.WorkMomentPublic}},
				}}},
		},
		},
	}, findOpts)
	if err != nil {
		return workMomentList, err
	}
	err = result.All(ctx, &workMomentList)
	return workMomentList, err
}

type SuperGroup struct {
	GroupID string `bson:"group_id"`
	//MemberNumCount int      `bson:"member_num_count"`
	MemberIDList []string `bson:"member_id_list"`
}

type UserToSuperGroup struct {
	UserID      string   `bson:"user_id"`
	GroupIDList []string `bson:"group_id_list"`
}

func (d *DataBases) CreateSuperGroup(groupID string, initMemberIDList []string, memberNumCount int) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cSuperGroup)
	session, err := d.mongoClient.StartSession()
	if err != nil {
		return utils.Wrap(err, "start session failed")
	}
	defer session.EndSession(ctx)
	sCtx := mongo.NewSessionContext(ctx, session)
	superGroup := SuperGroup{
		GroupID:      groupID,
		MemberIDList: initMemberIDList,
	}
	_, err = c.InsertOne(sCtx, superGroup)
	if err != nil {
		_ = session.AbortTransaction(ctx)
		return utils.Wrap(err, "transaction failed")
	}
	var users []UserToSuperGroup
	for _, v := range initMemberIDList {
		users = append(users, UserToSuperGroup{
			UserID: v,
		})
	}
	upsert := true
	opts := &options.UpdateOptions{
		Upsert: &upsert,
	}
	c = d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cUserToSuperGroup)
	//_, err = c.UpdateMany(sCtx, bson.M{"user_id": bson.M{"$in": initMemberIDList}}, bson.M{"$addToSet": bson.M{"group_id_list": groupID}}, opts)
	//if err != nil {
	//	session.AbortTransaction(ctx)
	//	return utils.Wrap(err, "transaction failed")
	//}
	for _, userID := range initMemberIDList {
		_, err = c.UpdateOne(sCtx, bson.M{"user_id": userID}, bson.M{"$addToSet": bson.M{"group_id_list": groupID}}, opts)
		if err != nil {
			_ = session.AbortTransaction(ctx)
			return utils.Wrap(err, "transaction failed")
		}

	}
	return err
}

func (d *DataBases) GetSuperGroup(groupID string) (SuperGroup, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cSuperGroup)
	superGroup := SuperGroup{}
	err := c.FindOne(ctx, bson.M{"group_id": groupID}).Decode(&superGroup)
	return superGroup, err
}

func (d *DataBases) AddUserToSuperGroup(groupID string, userIDList []string) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cSuperGroup)
	session, err := d.mongoClient.StartSession()
	if err != nil {
		return utils.Wrap(err, "start session failed")
	}
	defer session.EndSession(ctx)
	sCtx := mongo.NewSessionContext(ctx, session)
	if err != nil {
		return utils.Wrap(err, "start transaction failed")
	}
	_, err = c.UpdateOne(sCtx, bson.M{"group_id": groupID}, bson.M{"$addToSet": bson.M{"member_id_list": bson.M{"$each": userIDList}}})
	if err != nil {
		_ = session.AbortTransaction(ctx)
		return utils.Wrap(err, "transaction failed")
	}
	c = d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cUserToSuperGroup)
	var users []UserToSuperGroup
	for _, v := range userIDList {
		users = append(users, UserToSuperGroup{
			UserID: v,
		})
	}
	upsert := true
	opts := &options.UpdateOptions{
		Upsert: &upsert,
	}
	for _, userID := range userIDList {
		_, err = c.UpdateOne(sCtx, bson.M{"user_id": userID}, bson.M{"$addToSet": bson.M{"group_id_list": groupID}}, opts)
		if err != nil {
			_ = session.AbortTransaction(ctx)
			return utils.Wrap(err, "transaction failed")
		}
	}
	_ = session.CommitTransaction(ctx)
	return err
}

func (d *DataBases) RemoverUserFromSuperGroup(groupID string, userIDList []string) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cSuperGroup)
	session, err := d.mongoClient.StartSession()
	if err != nil {
		return utils.Wrap(err, "start session failed")
	}
	defer session.EndSession(ctx)
	sCtx := mongo.NewSessionContext(ctx, session)
	_, err = c.UpdateOne(ctx, bson.M{"group_id": groupID}, bson.M{"$pull": bson.M{"member_id_list": bson.M{"$in": userIDList}}})
	if err != nil {
		_ = session.AbortTransaction(ctx)
		return utils.Wrap(err, "transaction failed")
	}
	err = d.RemoveGroupFromUser(ctx, sCtx, groupID, userIDList)
	if err != nil {
		_ = session.AbortTransaction(ctx)
		return utils.Wrap(err, "transaction failed")
	}
	_ = session.CommitTransaction(ctx)
	return err
}

func (d *DataBases) GetSuperGroupByUserID(userID string) (UserToSuperGroup, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cUserToSuperGroup)
	var user UserToSuperGroup
	return user, c.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
}

func (d *DataBases) DeleteSuperGroup(groupID string) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cSuperGroup)
	session, err := d.mongoClient.StartSession()
	if err != nil {
		return utils.Wrap(err, "start session failed")
	}
	defer session.EndSession(ctx)
	sCtx := mongo.NewSessionContext(ctx, session)
	superGroup := &SuperGroup{}
	result := c.FindOneAndDelete(sCtx, bson.M{"group_id": groupID})
	err = result.Decode(superGroup)
	if err != nil {
		session.AbortTransaction(ctx)
		return utils.Wrap(err, "transaction failed")
	}
	if err = d.RemoveGroupFromUser(ctx, sCtx, groupID, superGroup.MemberIDList); err != nil {
		session.AbortTransaction(ctx)
		return utils.Wrap(err, "transaction failed")
	}
	session.CommitTransaction(ctx)
	return nil
}

func (d *DataBases) RemoveGroupFromUser(ctx, sCtx context.Context, groupID string, userIDList []string) error {
	var users []UserToSuperGroup
	for _, v := range userIDList {
		users = append(users, UserToSuperGroup{
			UserID: v,
		})
	}
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cUserToSuperGroup)
	_, err := c.UpdateOne(sCtx, bson.M{"user_id": bson.M{"$in": userIDList}}, bson.M{"$pull": bson.M{"group_id_list": groupID}})
	if err != nil {
		return utils.Wrap(err, "UpdateOne transaction failed")
	}
	return err
}

func generateTagID(tagName, userID string) string {
	return utils.Md5(tagName + userID + strconv.Itoa(rand.Int()) + time.Now().String())
}

func generateWorkMomentID(userID string) string {
	return utils.Md5(userID + strconv.Itoa(rand.Int()) + time.Now().String())
}

func generateWorkMomentCommentID(workMomentID string) string {
	return utils.Md5(workMomentID + strconv.Itoa(rand.Int()) + time.Now().String())
}

func getCurrentTimestampByMill() int64 {
	return time.Now().UnixNano() / 1e6
}
func GetCurrentTimestampByMill() int64 {
	return time.Now().UnixNano() / 1e6
}

func getSeqUid(uid string, seq uint32) string {
	seqSuffix := seq / singleGocMsgNum
	return indexGen(uid, seqSuffix)
}

func getSeqUserIDList(userID string, maxSeq uint32) []string {
	seqMaxSuffix := maxSeq / singleGocMsgNum
	var seqUserIDList []string
	for i := 0; i <= int(seqMaxSuffix); i++ {
		seqUserID := indexGen(userID, uint32(i))
		seqUserIDList = append(seqUserIDList, seqUserID)
	}
	return seqUserIDList
}

func getSeqSuperGroupID(groupID string, seq uint32) string {
	seqSuffix := seq / singleGocMsgNum
	return superGroupIndexGen(groupID, seqSuffix)
}

func GetSeqUid(uid string, seq uint32) string {
	return getSeqUid(uid, seq)
}

func getMsgIndex(seq uint32) int {
	seqSuffix := seq / singleGocMsgNum
	var index uint32
	if seqSuffix == 0 {
		index = (seq - seqSuffix*singleGocMsgNum) - 1
	} else {
		index = seq - seqSuffix*singleGocMsgNum
	}
	return int(index)
}

func isContainInt32(target uint32, List []uint32) bool {
	for _, element := range List {
		if target == element {
			return true
		}
	}
	return false
}

func isNotContainInt32(target uint32, List []uint32) bool {
	for _, i := range List {
		if i == target {
			return false
		}
	}
	return true
}

func indexGen(uid string, seqSuffix uint32) string {
	return uid + ":" + strconv.FormatInt(int64(seqSuffix), 10)
}
func superGroupIndexGen(groupID string, seqSuffix uint32) string {
	return "super_group_" + groupID + ":" + strconv.FormatInt(int64(seqSuffix), 10)
}

func (d *DataBases) CleanUpUserMsgFromMongo(userID string, operationID string) error {
	ctx := context.Background()
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cChat)
	maxSeq, err := d.GetUserMaxSeq(userID)
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return utils.Wrap(err, "")
	}

	seqUsers := getSeqUserIDList(userID, uint32(maxSeq))
	log.Error(operationID, "getSeqUserIDList", seqUsers)
	_, err = c.DeleteMany(ctx, bson.M{"uid": bson.M{"$in": seqUsers}})
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return utils.Wrap(err, "")
}
