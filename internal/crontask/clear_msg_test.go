package cronTask

import (
	"Open_IM/pkg/common/constant"
	mongo2 "Open_IM/pkg/common/db/mongo"
	server_api_params "Open_IM/pkg/proto/sdk_ws"
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/golang/protobuf/proto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"testing"
	"time"
)

var (
	redisClient *redis.Client
	mongoClient *mongo.Collection
)

func GenUserChat(startSeq, stopSeq, delSeq, index uint32, userID string) *mongo2.UserChat {
	chat := &mongo2.UserChat{UID: userID + strconv.Itoa(int(index))}
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
		sendTime := 0
		chat.Msg = append(chat.Msg, mongo2.MsgInfo{SendTime: int64(sendTime), Msg: bytes})
	}
	return chat
}

func SetUserMaxSeq(userID string, seq int) error {
	return redisClient.Set(context.Background(), "REDIS_USER_INCR_SEQ"+userID, seq, 0).Err()
}

func CreateChat(userChat *mongo2.UserChat) error {
	_, err := mongoClient.InsertOne(context.Background(), userChat)
	return err
}

func TestDeleteMongoMsgAndResetRedisSeq(t *testing.T) {
	operationID := getCronTaskOperationID()
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:16379",
		Password: "openIM123", // no password set
		DB:       13,          // use default DB
	})
	mongoUri := fmt.Sprintf("mongodb://%s:%s@%s/%s?maxPoolSize=%d&authSource=admin",
		"root", "openIM123", "127.0.0.1:37017",
		"openIM", 100)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUri))
	mongoClient = client.Database("openIM").Collection("msg")
	testUID1 := "test_del_id1"
	//testUID2 := "test_del_id2"
	//testUID3 := "test_del_id3"
	//testUID4 := "test_del_id4"
	//testUID5 := "test_del_id5"
	//testUID6 := "test_del_id6"
	err = SetUserMaxSeq(testUID1, 600)
	userChat := GenUserChat(1, 500, 200, 0, testUID1)
	err = CreateChat(userChat)

	if err := DeleteMongoMsgAndResetRedisSeq(operationID, testUID1); err != nil {
		t.Error("checkMaxSeqWithMongo failed", testUID1)
	}
	if err := checkMaxSeqWithMongo(operationID, testUID1, constant.WriteDiffusion); err != nil {
		t.Error("checkMaxSeqWithMongo failed", testUID1)
	}
	if err != nil {
		t.Error("err is not nil", testUID1, err.Error())
	}
	// testWorkingGroupIDList := []string{"test_del_id1", "test_del_id2", "test_del_id3", "test_del_id4", "test_del_id5"}
	// for _, groupID := range testWorkingGroupIDList {
	// 	operationID = groupID + "-" + operationID
	// 	log.NewDebug(operationID, utils.GetSelfFuncName(), "groupID:", groupID, "userIDList:", testUserIDList)
	// 	if err := ResetUserGroupMinSeq(operationID, groupID, testUserIDList); err != nil {
	// 		t.Error("checkMaxSeqWithMongo failed", groupID)
	// 	}
	// 	if err := checkMaxSeqWithMongo(operationID, groupID, constant.ReadDiffusion); err != nil {
	// 		t.Error("checkMaxSeqWithMongo failed", groupID)
	// 	}
	// }
}
