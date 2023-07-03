package tools

import (
	"context"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"

	unRelationTb "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/unrelation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/unrelation"

	"testing"
	"time"
)

func GenMsgDoc(startSeq, stopSeq, delSeq, index int64, conversationID string) *unRelationTb.MsgDocModel {
	msgDoc := &unRelationTb.MsgDocModel{DocID: conversationID + strconv.Itoa(int(index))}
	for i := 0; i < 5000; i++ {
		msgDoc.Msg = append(msgDoc.Msg, &unRelationTb.MsgInfoModel{})
	}
	for i := startSeq; i <= stopSeq; i++ {
		msg := &unRelationTb.MsgDataModel{
			SendID:           "sendID1",
			RecvID:           "recvID1",
			GroupID:          "",
			ClientMsgID:      "xxx",
			ServerMsgID:      "xxx",
			SenderPlatformID: 1,
			SenderNickname:   "testNickName",
			SenderFaceURL:    "testFaceURL",
			SessionType:      1,
			MsgFrom:          100,
			ContentType:      101,
			Content:          "testContent",
			Seq:              i,
			CreateTime:       time.Now().Unix(),
			Status:           1,
		}
		if i <= delSeq {
			msg.SendTime = 10000
		} else {
			msg.SendTime = utils.GetCurrentTimestampByMill()
		}
		msgDoc.Msg[i-1] = &unRelationTb.MsgInfoModel{Msg: msg}
	}
	return msgDoc
}

func TestDeleteMongoMsgAndResetRedisSeq(t *testing.T) {
	operationID := "test"
	rdb, err := cache.NewRedis()
	if err != nil {
		return
	}
	mgo, err := unrelation.NewMongo()
	if err != nil {
		return
	}
	cacheModel := cache.NewMsgCacheModel(rdb)
	mongoClient := mgo.GetDatabase().Collection(unRelationTb.MsgDocModel{}.TableName())
	ctx := context.Background()
	ctx = mcontext.SetOperationID(ctx, operationID)

	testUID1 := "test_del_id1"
	var conversationID string
	conversationID = utils.GetConversationIDBySessionType(constant.SuperGroupChatType, testUID1)
	_, err = mongoClient.DeleteOne(ctx, bson.M{"doc_id": conversationID + ":" + strconv.Itoa(0)})
	if err != nil {
		t.Error("DeleteOne failed")
		return
	}

	err = cacheModel.SetMaxSeq(ctx, conversationID, 600)
	if err != nil {
		t.Error("SetUserMaxSeq failed")
	}
	msgDoc := GenMsgDoc(1, 600, 200, 0, conversationID)
	if _, err := mongoClient.InsertOne(ctx, msgDoc); err != nil {
		t.Error("InsertOne failed", conversationID)
	}

	msgTools, err := InitMsgTool()
	if err != nil {
		t.Error("init failed")
		return
	}
	msgTools.ClearConversationsMsg(ctx, []string{conversationID})
	minSeqMongo, maxSeqMongo, minSeqCache, maxSeqCache, err := msgTools.msgDatabase.GetConversationMinMaxSeqInMongoAndCache(
		ctx,
		conversationID,
	)
	if err != nil {
		t.Error("GetSuperGroupMinMaxSeqInMongoAndCache failed")
		return
	}
	if maxSeqCache != maxSeqMongo {
		t.Error("checkMaxSeqWithMongo failed", conversationID)
	}
	if minSeqMongo != minSeqCache {
		t.Error("minSeqMongo != minSeqCache", minSeqMongo, minSeqCache)
	}
	if minSeqCache != 201 {
		t.Error("test1 is not the same", "minSeq:", minSeqCache, "targetSeq", 201)
	}

	/////// uid2

	testUID2 := "test_del_id2"
	conversationID = utils.GetConversationIDBySessionType(constant.SuperGroupChatType, testUID2)

	_, err = mongoClient.DeleteOne(ctx, bson.M{"uid": conversationID + ":" + strconv.Itoa(0)})
	if err != nil {
		t.Error("delete failed")
	}
	_, err = mongoClient.DeleteOne(ctx, bson.M{"uid": conversationID + ":" + strconv.Itoa(1)})
	if err != nil {
		t.Error("delete failed")
	}

	err = cacheModel.SetMaxSeq(ctx, conversationID, 7000)
	if err != nil {
		t.Error("SetUserMaxSeq failed")
	}
	msgDoc = GenMsgDoc(1, 4999, 5000, 0, conversationID)
	msgDoc2 := GenMsgDoc(5000, 7000, 6000, 1, conversationID)
	if _, err := mongoClient.InsertOne(ctx, msgDoc); err != nil {
		t.Error("InsertOne failed", testUID1)
	}
	if _, err := mongoClient.InsertOne(ctx, msgDoc2); err != nil {
		t.Error("InsertOne failed", testUID1)
	}

	msgTools.ClearConversationsMsg(ctx, []string{conversationID})
	minSeqMongo, maxSeqMongo, minSeqCache, maxSeqCache, err = msgTools.msgDatabase.GetConversationMinMaxSeqInMongoAndCache(
		ctx,
		conversationID,
	)
	if err != nil {
		t.Error("GetSuperGroupMinMaxSeqInMongoAndCache failed")
		return
	}
	if maxSeqCache != maxSeqMongo {
		t.Error("checkMaxSeqWithMongo failed", conversationID)
	}
	if minSeqMongo != minSeqCache {
		t.Error("minSeqMongo != minSeqCache", minSeqMongo, minSeqCache)
	}
	if minSeqCache != 6001 {
		t.Error("test1 is not the same", "minSeq:", minSeqCache, "targetSeq", 201)
	}

	/////// uid3
	testUID3 := "test_del_id3"
	conversationID = utils.GetConversationIDBySessionType(constant.SuperGroupChatType, testUID3)
	_, err = mongoClient.DeleteOne(ctx, bson.M{"uid": conversationID + ":" + strconv.Itoa(0)})
	if err != nil {
		t.Error("delete failed")
	}
	err = cacheModel.SetMaxSeq(ctx, conversationID, 4999)
	if err != nil {
		t.Error("SetUserMaxSeq failed")
	}
	msgDoc = GenMsgDoc(1, 4999, 5000, 0, conversationID)
	if _, err := mongoClient.InsertOne(ctx, msgDoc); err != nil {
		t.Error("InsertOne failed", conversationID)
	}

	msgTools.ClearConversationsMsg(ctx, []string{conversationID})
	minSeqMongo, maxSeqMongo, minSeqCache, maxSeqCache, err = msgTools.msgDatabase.GetConversationMinMaxSeqInMongoAndCache(
		ctx,
		conversationID,
	)
	if err != nil {
		t.Error("GetSuperGroupMinMaxSeqInMongoAndCache failed")
		return
	}
	if maxSeqCache != maxSeqMongo {
		t.Error("checkMaxSeqWithMongo failed", conversationID)
	}
	if minSeqMongo != minSeqCache {
		t.Error("minSeqMongo != minSeqCache", minSeqMongo, minSeqCache)
	}
	if minSeqCache != 5000 {
		t.Error("test1 is not the same", "minSeq:", minSeqCache, "targetSeq", 201)
	}

	//// uid4
	testUID4 := "test_del_id4"
	conversationID = utils.GetConversationIDBySessionType(constant.SuperGroupChatType, testUID4)
	_, err = mongoClient.DeleteOne(ctx, bson.M{"uid": conversationID + ":" + strconv.Itoa(0)})
	if err != nil {
		t.Error("delete failed")
	}
	_, err = mongoClient.DeleteOne(ctx, bson.M{"uid": conversationID + ":" + strconv.Itoa(1)})
	if err != nil {
		t.Error("delete failed")
	}
	_, err = mongoClient.DeleteOne(ctx, bson.M{"uid": conversationID + ":" + strconv.Itoa(2)})
	if err != nil {
		t.Error("delete failed")
	}

	err = cacheModel.SetMaxSeq(ctx, conversationID, 12000)
	msgDoc = GenMsgDoc(1, 4999, 5000, 0, conversationID)
	msgDoc2 = GenMsgDoc(5000, 9999, 10000, 1, conversationID)
	msgDoc3 := GenMsgDoc(10000, 12000, 11000, 2, conversationID)
	if _, err := mongoClient.InsertOne(ctx, msgDoc); err != nil {
		t.Error("InsertOne failed", conversationID)
	}
	if _, err := mongoClient.InsertOne(ctx, msgDoc2); err != nil {
		t.Error("InsertOne failed", conversationID)
	}
	if _, err := mongoClient.InsertOne(ctx, msgDoc3); err != nil {
		t.Error("InsertOne failed", conversationID)
	}

	msgTools.ClearConversationsMsg(ctx, []string{conversationID})
	if err != nil {
		t.Error("GetSuperGroupMinMaxSeqInMongoAndCache failed")
		return
	}
	minSeqMongo, maxSeqMongo, minSeqCache, maxSeqCache, err = msgTools.msgDatabase.GetConversationMinMaxSeqInMongoAndCache(
		ctx,
		conversationID,
	)
	if err != nil {
		t.Error("GetSuperGroupMinMaxSeqInMongoAndCache failed")
		return
	}
	if maxSeqCache != maxSeqMongo {
		t.Error("checkMaxSeqWithMongo failed", conversationID)
	}
	if minSeqMongo != minSeqCache {
		t.Error("minSeqMongo != minSeqCache", minSeqMongo, minSeqCache)
	}
	if minSeqCache != 5000 {
		t.Error("test1 is not the same", "minSeq:", minSeqCache)
	}

	testUID5 := "test_del_id5"
	conversationID = utils.GetConversationIDBySessionType(constant.SuperGroupChatType, testUID5)

	_, err = mongoClient.DeleteOne(ctx, bson.M{"uid": conversationID + ":" + strconv.Itoa(0)})
	if err != nil {
		t.Error("delete failed")
	}
	_, err = mongoClient.DeleteOne(ctx, bson.M{"uid": conversationID + ":" + strconv.Itoa(1)})
	if err != nil {
		t.Error("delete failed")
	}
	err = cacheModel.SetMaxSeq(ctx, conversationID, 9999)
	msgDoc = GenMsgDoc(1, 4999, 5000, 0, conversationID)
	msgDoc2 = GenMsgDoc(5000, 9999, 10000, 1, conversationID)
	if _, err := mongoClient.InsertOne(ctx, msgDoc); err != nil {
		t.Error("InsertOne failed", conversationID)
	}
	if _, err := mongoClient.InsertOne(ctx, msgDoc2); err != nil {
		t.Error("InsertOne failed", conversationID)
	}

	msgTools.ClearConversationsMsg(ctx, []string{conversationID})
	if err != nil {
		t.Error("GetSuperGroupMinMaxSeqInMongoAndCache failed")
		return
	}
	minSeqMongo, maxSeqMongo, minSeqCache, maxSeqCache, err = msgTools.msgDatabase.GetConversationMinMaxSeqInMongoAndCache(
		ctx,
		conversationID,
	)
	if err != nil {
		t.Error("GetSuperGroupMinMaxSeqInMongoAndCache failed")
		return
	}
	if maxSeqCache != maxSeqMongo {
		t.Error("checkMaxSeqWithMongo failed", conversationID)
	}
	if minSeqMongo != minSeqCache {
		t.Error("minSeqMongo != minSeqCache", minSeqMongo, minSeqCache)
	}
	if minSeqCache != 10000 {
		t.Error("test1 is not the same", "minSeq:", minSeqCache)
	}

	testUID6 := "test_del_id6"
	conversationID = utils.GetConversationIDBySessionType(constant.SuperGroupChatType, testUID6)

	_, err = mongoClient.DeleteOne(ctx, bson.M{"uid": conversationID + ":" + strconv.Itoa(0)})
	if err != nil {
		t.Error("delete failed")
	}
	_, err = mongoClient.DeleteOne(ctx, bson.M{"uid": conversationID + ":" + strconv.Itoa(1)})
	if err != nil {
		t.Error("delete failed")
	}
	_, err = mongoClient.DeleteOne(ctx, bson.M{"uid": conversationID + ":" + strconv.Itoa(2)})
	if err != nil {
		t.Error("delete failed")
	}
	_, err = mongoClient.DeleteOne(ctx, bson.M{"uid": conversationID + ":" + strconv.Itoa(3)})
	if err != nil {
		t.Error("delete failed")
	}
	msgDoc = GenMsgDoc(1, 4999, 5000, 0, conversationID)
	msgDoc2 = GenMsgDoc(5000, 9999, 10000, 1, conversationID)
	msgDoc3 = GenMsgDoc(10000, 14999, 13000, 2, conversationID)
	msgDoc4 := GenMsgDoc(15000, 19999, 0, 3, conversationID)
	if _, err := mongoClient.InsertOne(ctx, msgDoc); err != nil {
		t.Error("InsertOne failed", testUID4)
	}
	if _, err := mongoClient.InsertOne(ctx, msgDoc2); err != nil {
		t.Error("InsertOne failed", testUID4)
	}
	if _, err := mongoClient.InsertOne(ctx, msgDoc3); err != nil {
		t.Error("InsertOne failed", testUID4)
	}
	if _, err := mongoClient.InsertOne(ctx, msgDoc4); err != nil {
		t.Error("InsertOne failed", testUID4)
	}
	minSeqMongo, maxSeqMongo, minSeqCache, maxSeqCache, err = msgTools.msgDatabase.GetConversationMinMaxSeqInMongoAndCache(
		ctx,
		conversationID,
	)
	if err != nil {
		t.Error("GetSuperGroupMinMaxSeqInMongoAndCache failed")
		return
	}
	if maxSeqCache != maxSeqMongo {
		t.Error("checkMaxSeqWithMongo failed", conversationID)
	}
	if minSeqMongo != minSeqCache {
		t.Error("minSeqMongo != minSeqCache", minSeqMongo, minSeqCache)
	}
	if minSeqCache != 13001 {
		t.Error("test1 is not the same", "minSeq:", minSeqCache)
	}
}
