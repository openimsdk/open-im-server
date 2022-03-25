package db

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	pbMsg "Open_IM/pkg/proto/chat"
	officePb "Open_IM/pkg/proto/office"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	"fmt"
	"github.com/gogo/protobuf/sortkeys"
	"math/rand"

	//"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
	"gopkg.in/mgo.v2/bson"

	"strconv"
	"time"
)

const cChat = "msg"
const cGroup = "group"
const cTag = "tag"
const cSendLog = "sendLog"
const singleGocMsgNum = 5000

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

func genExceptionMessageBySeqList(seqList []uint32) (exceptionMsg []*open_im_sdk.MsgData) {
	for _, v := range seqList {
		msg := new(open_im_sdk.MsgData)
		msg.Seq = v
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
	log.NewDebug(operationID, "get mgoSession cost time", getCurrentTimestampByMill()-newTime)
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
	TagID    string   `bson:"tagID"`
	TagName  string   `bson:"tagName"`
	UserList []string `bson:"userList"`
}

type TagsStruct struct {
	Uid  string         `bson:"uid"`
	Tags map[string]Tag `bson:"tags"`
}

type TagSendLogStruct struct {
}

func (d *DataBases) GetUserTags(userID string) ([]Tag, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cTag)
	var tagStruct TagsStruct
	var tags []Tag
	_ = c.FindOne(ctx, bson.M{"uid": userID}).Decode(&tagStruct)
	for _, v := range tagStruct.Tags {
		tags = append(tags, v)
	}
	return tags, nil
}

func (d *DataBases) CreateTag(userID, tagName string, userList []string) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cTag)
	tagID := generateTagID(tagName, userID)
	tag := Tag{
		TagID:    tagID,
		TagName:  tagName,
		UserList: userList,
	}
	_, err := c.InsertOne(ctx, TagsStruct{
		Uid:  userID,
		Tags: map[string]Tag{tagID: tag},
	})
	return err
}

func (d *DataBases) DeleteTag(userID, tagID string) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cTag)
	_, err := c.DeleteOne(ctx, bson.M{"uid": userID, "tags": bson.M{"$unset": tagID}})
	return err
}

func (d *DataBases) SetTag(userID, tagID, newName string, increaseUserList []string, reduceUserIDList []string) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cTag)
	_, err := c.UpdateOne(ctx, bson.M{"uid": userID, "tags": tagID}, bson.M{"tagName": newName})
	if err != nil {
		return err
	}
	_, err = c.InsertOne(ctx, bson.M{"uid": userID, "tags": bson.M{tagID: ""}})
	if err != nil {
		return err
	}
	//_, err = c.InsertOne(ctx)
	//if err != nil {
	//	return err
	//}
	return nil
}

func (d *DataBases) GetUserIDListByTagID(userID, tagID string) ([]string, error) {
	var tagIDList []string
	var tagStruct TagsStruct
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cTag)
	_ = c.FindOne(ctx, bson.M{"uid": userID}).Decode(&tagStruct)
	for k, tag := range tagStruct.Tags {
		if k == tagID {
			tagIDList = tag.UserList
		}
	}
	return tagIDList, nil
}

type TagSendLog struct {
	TagID            string   `bson:"tagID"`
	SendID           string   `bson:"sendID"`
	SenderPlatformID int32    `bson:"senderPlatformID"`
	Content          string   `bson:"content"`
	ContentType      int32    `bson:"contentType"`
	SendTime         int64    `bson:"sendTime"`
	UserList         []string `bson:"userList"`
}

func (d *DataBases) SaveTagSendLog(sendReq *officePb.SendMsg2TagReq) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cSendLog)
	tagSendLog := TagSendLog{
		TagID:            sendReq.TagID,
		SendID:           sendReq.SendID,
		SenderPlatformID: sendReq.SenderPlatformID,
		Content:          sendReq.Content,
		ContentType:      sendReq.ContentType,
		SendTime:         time.Now().Unix(),
	}
	_, err := c.InsertOne(ctx, tagSendLog)
	return err
}

func (d *DataBases) GetTagSendLogs(userID string, showNumber, pageNumber int32) ([]*TagSendLog, error) {
	var tagSendLogs []*TagSendLog
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cSendLog)
	cursor, err := c.Find(ctx, bson.M{"sendID": userID})
	if err != nil {
		return tagSendLogs, err
	}
	err = cursor.Decode(&tagSendLogs)
	if err != nil {
		return tagSendLogs, err
	}
	return tagSendLogs, nil
}

func generateTagID(tagName, userID string) string {
	return utils.Md5(tagName + userID + strconv.Itoa(rand.Int()))
}

func getCurrentTimestampByMill() int64 {
	return time.Now().UnixNano() / 1e6
}

func getSeqUid(uid string, seq uint32) string {
	seqSuffix := seq / singleGocMsgNum
	return indexGen(uid, seqSuffix)
}

func getMsgIndex(seq uint32) int {
	seqSuffix := seq / singleGocMsgNum
	var index uint32
	if seqSuffix == 0 {
		index = (seq - seqSuffix*5000) - 1
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
func indexGen(uid string, seqSuffix uint32) string {
	return uid + ":" + strconv.FormatInt(int64(seqSuffix), 10)
}
