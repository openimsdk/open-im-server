package tools

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"github.com/golang/protobuf/proto"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"

	unRelationTb "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/unrelation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/unrelation"

	"testing"
	"time"
)

func GenMsgDoc(startSeq, stopSeq, delSeq, index int64, userID string) *unRelationTb.MsgDocModel {
	msgDoc := &unRelationTb.MsgDocModel{DocID: userID + strconv.Itoa(int(index))}
	for i := startSeq; i <= stopSeq; i++ {
		msg := sdkws.MsgData{
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
			Content:          []byte("testFaceURL"),
			Seq:              i,
			SendTime:         time.Now().Unix(),
			CreateTime:       time.Now().Unix(),
			Status:           1,
		}
		bytes, _ := proto.Marshal(&msg)
		var sendTime int64
		if i <= delSeq {
			sendTime = 10000
		} else {
			sendTime = utils.GetCurrentTimestampByMill()
		}
		msgDoc.Msg = append(msgDoc.Msg, unRelationTb.MsgInfoModel{SendTime: int64(sendTime), Msg: bytes})
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
	cacheModel := cache.NewCacheModel(rdb)
	mongoClient := mgo.GetDatabase().Collection(unRelationTb.MsgDocModel{}.TableName())

	ctx := context.Background()
	mcontext.SetOperationID(ctx, operationID)
	testUID1 := "test_del_id1"
	_, err = mongoClient.DeleteOne(ctx, bson.M{"uid": testUID1 + ":" + strconv.Itoa(0)})
	if err != nil {
		t.Error("DeleteOne failed")
		return
	}
	err = cacheModel.SetUserMaxSeq(ctx, testUID1, 600)
	if err != nil {
		t.Error("SetUserMaxSeq failed")
	}
	msgDoc := GenMsgDoc(1, 600, 200, 0, testUID1)
	if _, err := mongoClient.InsertOne(ctx, msgDoc); err != nil {
		t.Error("InsertOne failed", testUID1)
	}

	msgTools, err := InitMsgTool()
	if err != nil {
		t.Error("init failed")
		return
	}
	msgTools.ClearUsersMsg(ctx, []string{testUID1})
	minSeqMongo, maxSeqMongo, minSeqCache, maxSeqCache, err := msgTools.msgDatabase.GetUserMinMaxSeqInMongoAndCache(ctx, testUID1)
	if err != nil {
		t.Error("GetSuperGroupMinMaxSeqInMongoAndCache failed")
		return
	}
	if err := msgTools.CheckMaxSeqWithMongo(ctx, testUID1, maxSeqCache, maxSeqMongo, constant.WriteDiffusion); err != nil {
		t.Error("checkMaxSeqWithMongo failed", testUID1)
	}
	if minSeqMongo != minSeqCache {
		t.Error("minSeqMongo != minSeqCache", minSeqMongo, minSeqCache)
	}
	if minSeqCache != 201 {
		t.Error("test1 is not the same", "minSeq:", minSeqCache, "targetSeq", 201)
	}

	/////// uid2

	testUID2 := "test_del_id2"
	_, err = mongoClient.DeleteOne(ctx, bson.M{"uid": testUID2 + ":" + strconv.Itoa(0)})
	if err != nil {
		t.Error("delete failed")
	}
	_, err = mongoClient.DeleteOne(ctx, bson.M{"uid": testUID2 + ":" + strconv.Itoa(1)})
	if err != nil {
		t.Error("delete failed")
	}

	err = cacheModel.SetUserMaxSeq(ctx, testUID2, 7000)
	if err != nil {
		t.Error("SetUserMaxSeq failed")
	}
	msgDoc = GenMsgDoc(1, 4999, 5000, 0, testUID2)
	msgDoc2 := GenMsgDoc(5000, 7000, 6000, 1, testUID2)
	if _, err := mongoClient.InsertOne(ctx, msgDoc); err != nil {
		t.Error("InsertOne failed", testUID1)
	}
	if _, err := mongoClient.InsertOne(ctx, msgDoc2); err != nil {
		t.Error("InsertOne failed", testUID1)
	}

	msgTools.ClearUsersMsg(ctx, []string{testUID2})
	minSeqMongo, maxSeqMongo, minSeqCache, maxSeqCache, err = msgTools.msgDatabase.GetUserMinMaxSeqInMongoAndCache(ctx, testUID2)
	if err != nil {
		t.Error("GetSuperGroupMinMaxSeqInMongoAndCache failed")
		return
	}
	if err := msgTools.CheckMaxSeqWithMongo(ctx, testUID2, maxSeqCache, maxSeqMongo, constant.WriteDiffusion); err != nil {
		t.Error("checkMaxSeqWithMongo failed", testUID2)
	}
	if minSeqMongo != minSeqCache {
		t.Error("minSeqMongo != minSeqCache", minSeqMongo, minSeqCache)
	}
	if minSeqCache != 6001 {
		t.Error("test1 is not the same", "minSeq:", minSeqCache, "targetSeq", 201)
	}

	/////// uid3
	testUID3 := "test_del_id3"
	_, err = mongoClient.DeleteOne(ctx, bson.M{"uid": testUID3 + ":" + strconv.Itoa(0)})
	if err != nil {
		t.Error("delete failed")
	}
	err = cacheModel.SetUserMaxSeq(ctx, testUID3, 4999)
	if err != nil {
		t.Error("SetUserMaxSeq failed")
	}
	msgDoc = GenMsgDoc(1, 4999, 5000, 0, testUID3)
	if _, err := mongoClient.InsertOne(ctx, msgDoc); err != nil {
		t.Error("InsertOne failed", testUID3)
	}

	msgTools.ClearUsersMsg(ctx, []string{testUID3})
	minSeqMongo, maxSeqMongo, minSeqCache, maxSeqCache, err = msgTools.msgDatabase.GetUserMinMaxSeqInMongoAndCache(ctx, testUID3)
	if err != nil {
		t.Error("GetSuperGroupMinMaxSeqInMongoAndCache failed")
		return
	}
	if err := msgTools.CheckMaxSeqWithMongo(ctx, testUID3, maxSeqCache, maxSeqMongo, constant.WriteDiffusion); err != nil {
		t.Error("checkMaxSeqWithMongo failed", testUID3)
	}
	if minSeqMongo != minSeqCache {
		t.Error("minSeqMongo != minSeqCache", minSeqMongo, minSeqCache)
	}
	if minSeqCache != 5000 {
		t.Error("test1 is not the same", "minSeq:", minSeqCache, "targetSeq", 201)
	}

	//// uid4
	testUID4 := "test_del_id4"
	_, err = mongoClient.DeleteOne(ctx, bson.M{"uid": testUID4 + ":" + strconv.Itoa(0)})
	if err != nil {
		t.Error("delete failed")
	}
	_, err = mongoClient.DeleteOne(ctx, bson.M{"uid": testUID4 + ":" + strconv.Itoa(1)})
	if err != nil {
		t.Error("delete failed")
	}
	_, err = mongoClient.DeleteOne(ctx, bson.M{"uid": testUID4 + ":" + strconv.Itoa(2)})
	if err != nil {
		t.Error("delete failed")
	}

	err = cacheModel.SetUserMaxSeq(ctx, testUID4, 12000)
	msgDoc = GenMsgDoc(1, 4999, 5000, 0, testUID4)
	msgDoc2 = GenMsgDoc(5000, 9999, 10000, 1, testUID4)
	msgDoc3 := GenMsgDoc(10000, 12000, 11000, 2, testUID4)
	if _, err := mongoClient.InsertOne(ctx, msgDoc); err != nil {
		t.Error("InsertOne failed", testUID4)
	}
	if _, err := mongoClient.InsertOne(ctx, msgDoc2); err != nil {
		t.Error("InsertOne failed", testUID4)
	}
	if _, err := mongoClient.InsertOne(ctx, msgDoc3); err != nil {
		t.Error("InsertOne failed", testUID4)
	}

	msgTools.ClearUsersMsg(ctx, []string{testUID4})
	if err != nil {
		t.Error("GetSuperGroupMinMaxSeqInMongoAndCache failed")
		return
	}
	minSeqMongo, maxSeqMongo, minSeqCache, maxSeqCache, err = msgTools.msgDatabase.GetUserMinMaxSeqInMongoAndCache(ctx, testUID4)
	if err != nil {
		t.Error("GetSuperGroupMinMaxSeqInMongoAndCache failed")
		return
	}
	if err := msgTools.CheckMaxSeqWithMongo(ctx, testUID4, maxSeqCache, maxSeqMongo, constant.WriteDiffusion); err != nil {
		t.Error("checkMaxSeqWithMongo failed", testUID4)
	}
	if minSeqMongo != minSeqCache {
		t.Error("minSeqMongo != minSeqCache", minSeqMongo, minSeqCache)
	}
	if minSeqCache != 5000 {
		t.Error("test1 is not the same", "minSeq:", minSeqCache)
	}

	testUID5 := "test_del_id5"
	_, err = mongoClient.DeleteOne(ctx, bson.M{"uid": testUID5 + ":" + strconv.Itoa(0)})
	if err != nil {
		t.Error("delete failed")
	}
	_, err = mongoClient.DeleteOne(ctx, bson.M{"uid": testUID5 + ":" + strconv.Itoa(1)})
	if err != nil {
		t.Error("delete failed")
	}
	err = cacheModel.SetUserMaxSeq(ctx, testUID5, 9999)
	msgDoc = GenMsgDoc(1, 4999, 5000, 0, testUID5)
	msgDoc2 = GenMsgDoc(5000, 9999, 10000, 1, testUID5)
	if _, err := mongoClient.InsertOne(ctx, msgDoc); err != nil {
		t.Error("InsertOne failed", testUID5)
	}
	if _, err := mongoClient.InsertOne(ctx, msgDoc2); err != nil {
		t.Error("InsertOne failed", testUID5)
	}

	msgTools.ClearUsersMsg(ctx, []string{testUID5})
	if err != nil {
		t.Error("GetSuperGroupMinMaxSeqInMongoAndCache failed")
		return
	}
	minSeqMongo, maxSeqMongo, minSeqCache, maxSeqCache, err = msgTools.msgDatabase.GetUserMinMaxSeqInMongoAndCache(ctx, testUID5)
	if err != nil {
		t.Error("GetSuperGroupMinMaxSeqInMongoAndCache failed")
		return
	}
	if err := msgTools.CheckMaxSeqWithMongo(ctx, testUID5, maxSeqCache, maxSeqMongo, constant.WriteDiffusion); err != nil {
		t.Error("checkMaxSeqWithMongo failed", testUID5)
	}
	if minSeqMongo != minSeqCache {
		t.Error("minSeqMongo != minSeqCache", minSeqMongo, minSeqCache)
	}
	if minSeqCache != 10000 {
		t.Error("test1 is not the same", "minSeq:", minSeqCache)
	}

	testUID6 := "test_del_id6"
	_, err = mongoClient.DeleteOne(ctx, bson.M{"uid": testUID6 + ":" + strconv.Itoa(0)})
	if err != nil {
		t.Error("delete failed")
	}
	_, err = mongoClient.DeleteOne(ctx, bson.M{"uid": testUID6 + ":" + strconv.Itoa(1)})
	if err != nil {
		t.Error("delete failed")
	}
	_, err = mongoClient.DeleteOne(ctx, bson.M{"uid": testUID6 + ":" + strconv.Itoa(2)})
	if err != nil {
		t.Error("delete failed")
	}
	_, err = mongoClient.DeleteOne(ctx, bson.M{"uid": testUID6 + ":" + strconv.Itoa(3)})
	if err != nil {
		t.Error("delete failed")
	}
	msgDoc = GenMsgDoc(1, 4999, 5000, 0, testUID6)
	msgDoc2 = GenMsgDoc(5000, 9999, 10000, 1, testUID6)
	msgDoc3 = GenMsgDoc(10000, 14999, 13000, 2, testUID6)
	msgDoc4 := GenMsgDoc(15000, 19999, 0, 3, testUID6)
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
	minSeqMongo, maxSeqMongo, minSeqCache, maxSeqCache, err = msgTools.msgDatabase.GetUserMinMaxSeqInMongoAndCache(ctx, testUID6)
	if err != nil {
		t.Error("GetSuperGroupMinMaxSeqInMongoAndCache failed")
		return
	}
	if err := msgTools.CheckMaxSeqWithMongo(ctx, testUID6, maxSeqCache, maxSeqMongo, constant.WriteDiffusion); err != nil {
		t.Error("checkMaxSeqWithMongo failed", testUID6)
	}
	if minSeqMongo != minSeqCache {
		t.Error("minSeqMongo != minSeqCache", minSeqMongo, minSeqCache)
	}
	if minSeqCache != 13001 {
		t.Error("test1 is not the same", "minSeq:", minSeqCache)
	}
}
