package cronTask

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	server_api_params "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/golang/protobuf/proto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"

	"testing"
	"time"
)

var (
	redisClient *redis.Client
	mongoClient *mongo.Collection
)

func GenUserChat(startSeq, stopSeq, delSeq, index uint32, userID string) *db.UserChat {
	chat := &db.UserChat{UID: userID + ":" + strconv.Itoa(int(index))}
	for i := startSeq; i <= stopSeq; i++ {
		msg := server_api_params.MsgData{
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
			Seq:              uint32(i),
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
		chat.Msg = append(chat.Msg, db.MsgInfo{SendTime: int64(sendTime), Msg: bytes})
	}
	return chat
}

func SetUserMaxSeq(userID string, seq int) error {
	return redisClient.Set(context.Background(), "REDIS_USER_INCR_SEQ"+userID, seq, 0).Err()
}

func GetUserMinSeq(userID string) (uint64, error) {
	key := "REDIS_USER_MIN_SEQ:" + userID
	seq, err := redisClient.Get(context.Background(), key).Result()
	return uint64(utils.StringToInt(seq)), err
}

func CreateChat(userChat *db.UserChat) error {
	_, err := mongoClient.InsertOne(context.Background(), userChat)
	return err
}

func DelChat(uid string, index int) error {
	_, err := mongoClient.DeleteOne(context.Background(), bson.M{"uid": uid + ":" + strconv.Itoa(index)})
	return err
}

func TestDeleteMongoMsgAndResetRedisSeq(t *testing.T) {
	operationID := getCronTaskOperationID()
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:16379",
		Password: "openIM123", // no password set
		DB:       0,           // use default DB
	})
	mongoUri := fmt.Sprintf("mongodb://%s:%s@%s/%s?maxPoolSize=%d&authSource=admin",
		"root", "openIM123", "127.0.0.1:37017",
		"openIM", 100)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUri))
	mongoClient = client.Database("openIM").Collection("msg")
	testUID1 := "test_del_id1"
	err = DelChat(testUID1, 0)
	err = SetUserMaxSeq(testUID1, 600)
	userChat := GenUserChat(1, 600, 200, 0, testUID1)
	err = CreateChat(userChat)
	if err := DeleteMongoMsgAndResetRedisSeq(operationID, testUID1); err != nil {
		t.Error("checkMaxSeqWithMongo failed", testUID1)
	}
	if err := checkMaxSeqWithMongo(operationID, testUID1, constant.WriteDiffusion); err != nil {
		t.Error("checkMaxSeqWithMongo failed", testUID1)
	}
	minSeq, err := GetUserMinSeq(testUID1)
	if err != nil {
		t.Error("err is not nil", testUID1, err.Error())
	}
	if minSeq != 201 {
		t.Error("test1 is not the same", "minSeq:", minSeq, "targetSeq", 201)
	}

	testUID2 := "test_del_id2"
	err = DelChat(testUID2, 0)
	err = DelChat(testUID2, 1)
	err = SetUserMaxSeq(testUID2, 7000)
	userChat = GenUserChat(1, 4999, 5000, 0, testUID2)
	userChat2 := GenUserChat(5000, 7000, 6000, 1, testUID2)
	err = CreateChat(userChat)
	err = CreateChat(userChat2)

	if err := DeleteMongoMsgAndResetRedisSeq(operationID, testUID2); err != nil {
		t.Error("checkMaxSeqWithMongo failed", testUID2)
	}
	if err := checkMaxSeqWithMongo(operationID, testUID2, constant.WriteDiffusion); err != nil {
		t.Error("checkMaxSeqWithMongo failed", testUID2)
	}
	minSeq, err = GetUserMinSeq(testUID2)
	if err != nil {
		t.Error("err is not nil", testUID2, err.Error())
	}
	if minSeq != 6001 {
		t.Error("test2 is not the same", "minSeq:", minSeq, "targetSeq", 6001)
	}

	testUID3 := "test_del_id3"
	err = DelChat(testUID3, 0)
	err = SetUserMaxSeq(testUID3, 4999)
	userChat = GenUserChat(1, 4999, 5000, 0, testUID3)
	err = CreateChat(userChat)
	if err := DeleteMongoMsgAndResetRedisSeq(operationID, testUID3); err != nil {
		t.Error("checkMaxSeqWithMongo failed", testUID3)
	}
	if err := checkMaxSeqWithMongo(operationID, testUID3, constant.WriteDiffusion); err != nil {
		t.Error("checkMaxSeqWithMongo failed", testUID3)
	}
	minSeq, err = GetUserMinSeq(testUID3)
	if err != nil {
		t.Error("err is not nil", testUID3, err.Error())
	}
	if minSeq != 5000 {
		t.Error("test3 is not the same", "minSeq:", minSeq, "targetSeq", 5000)
	}

	testUID4 := "test_del_id4"
	err = DelChat(testUID4, 0)
	err = DelChat(testUID4, 1)
	err = DelChat(testUID4, 2)
	err = SetUserMaxSeq(testUID4, 12000)
	userChat = GenUserChat(1, 4999, 5000, 0, testUID4)
	userChat2 = GenUserChat(5000, 9999, 10000, 1, testUID4)
	userChat3 := GenUserChat(10000, 12000, 11000, 2, testUID4)
	err = CreateChat(userChat)
	err = CreateChat(userChat2)
	err = CreateChat(userChat3)
	if err := DeleteMongoMsgAndResetRedisSeq(operationID, testUID4); err != nil {
		t.Error("checkMaxSeqWithMongo failed", testUID4)
	}
	if err := checkMaxSeqWithMongo(operationID, testUID4, constant.WriteDiffusion); err != nil {
		t.Error("checkMaxSeqWithMongo failed", testUID4)
	}
	minSeq, err = GetUserMinSeq(testUID4)
	if err != nil {
		t.Error("err is not nil", testUID4, err.Error())
	}
	if minSeq != 11001 {
		t.Error("test4 is not the same", "minSeq:", minSeq, "targetSeq", 11001)
	}

	testUID5 := "test_del_id5"
	err = DelChat(testUID5, 0)
	err = DelChat(testUID5, 1)
	err = SetUserMaxSeq(testUID5, 9999)
	userChat = GenUserChat(1, 4999, 5000, 0, testUID5)
	userChat2 = GenUserChat(5000, 9999, 10000, 1, testUID5)
	err = CreateChat(userChat)
	err = CreateChat(userChat2)
	if err := DeleteMongoMsgAndResetRedisSeq(operationID, testUID5); err != nil {
		t.Error("checkMaxSeqWithMongo failed", testUID4)
	}
	if err := checkMaxSeqWithMongo(operationID, testUID5, constant.WriteDiffusion); err != nil {
		t.Error("checkMaxSeqWithMongo failed", testUID5)
	}
	minSeq, err = GetUserMinSeq(testUID5)
	if err != nil {
		t.Error("err is not nil", testUID5, err.Error())
	}
	if minSeq != 10000 {
		t.Error("test5 is not the same", "minSeq:", minSeq, "targetSeq", 10000)
	}

	testUID6 := "test_del_id6"
	err = DelChat(testUID5, 0)
	err = DelChat(testUID5, 1)
	err = DelChat(testUID5, 2)
	err = DelChat(testUID5, 3)
	userChat = GenUserChat(1, 4999, 5000, 0, testUID6)
	userChat2 = GenUserChat(5000, 9999, 10000, 1, testUID6)
	userChat3 = GenUserChat(10000, 14999, 13000, 2, testUID6)
	userChat4 := GenUserChat(15000, 19999, 0, 3, testUID6)
	err = CreateChat(userChat)
	err = CreateChat(userChat2)
	err = CreateChat(userChat3)
	err = CreateChat(userChat4)
	if err := DeleteMongoMsgAndResetRedisSeq(operationID, testUID6); err != nil {
		t.Error("checkMaxSeqWithMongo failed", testUID6)
	}
	if err := checkMaxSeqWithMongo(operationID, testUID6, constant.WriteDiffusion); err != nil {
		t.Error("checkMaxSeqWithMongo failed", testUID6)
	}
	minSeq, err = GetUserMinSeq(testUID6)
	if err != nil {
		t.Error("err is not nil", testUID6, err.Error())
	}
	if minSeq != 13001 {
		t.Error("test3 is not the same", "minSeq:", minSeq, "targetSeq", 13001)
	}
}
