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
	minSeq, err := GetUserMinSeq(testUID1)
	if err != nil {
		t.Error("err is not nil", testUID1, err.Error())
	}
	if minSeq != 201 {
		t.Error("is not the same", "minSeq:", minSeq, "targetSeq", 201)
	}
}
