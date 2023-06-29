package controller

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	unRelationTb "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/unrelation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/unrelation"
	"go.mongodb.org/mongo-driver/bson"
)

func Test_BatchInsertChat2DB(t *testing.T) {
	config.Config.Mongo.Address = []string{"192.168.44.128:37017"}
	// config.Config.Mongo.Timeout = 60
	config.Config.Mongo.Database = "openIM"
	// config.Config.Mongo.Source = "admin"
	config.Config.Mongo.Username = "root"
	config.Config.Mongo.Password = "openIM123"
	config.Config.Mongo.MaxPoolSize = 100
	config.Config.RetainChatRecords = 3650
	config.Config.ChatRecordsClearTime = "0 2 * * 3"

	mongo, err := unrelation.NewMongo()
	if err != nil {
		t.Fatal(err)
	}
	err = mongo.GetDatabase().Client().Ping(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	db := &commonMsgDatabase{
		msgDocDatabase: unrelation.NewMsgMongoDriver(mongo.GetDatabase()),
	}

	//ctx := context.Background()
	//msgs := make([]*sdkws.MsgData, 0, 1)
	//for i := 0; i < cap(msgs); i++ {
	//	msgs = append(msgs, &sdkws.MsgData{
	//		Content:  []byte(fmt.Sprintf("test-%d", i)),
	//		SendTime: time.Now().UnixMilli(),
	//	})
	//}
	//err = db.BatchInsertChat2DB(ctx, "test", msgs, 0)
	//if err != nil {
	//	panic(err)
	//}

	_ = db.BatchInsertChat2DB
	c := mongo.GetDatabase().Collection("msg")

	ch := make(chan int)
	rand.Seed(time.Now().UnixNano())

	index := 10

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(channelID int) {
			defer wg.Done()
			<-ch
			var arr []string
			for i := 0; i < 500; i++ {
				arr = append(arr, strconv.Itoa(i+1))
			}
			rand.Shuffle(len(arr), func(i, j int) {
				arr[i], arr[j] = arr[j], arr[i]
			})
			for j, s := range arr {
				if j == 0 {
					fmt.Printf("channnelID: %d, arr[0]: %s\n", channelID, arr[j])
				}
				filter := bson.M{"doc_id": "test:0"}
				update := bson.M{
					"$addToSet": bson.M{
						fmt.Sprintf("msgs.%d.del_list", index): bson.M{"$each": []string{s}},
					},
				}
				_, err := c.UpdateOne(context.Background(), filter, update)
				if err != nil {
					t.Fatal(err)
				}
			}
		}(i)
	}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ch
			var arr []string
			for i := 0; i < 500; i++ {
				arr = append(arr, strconv.Itoa(1001+i))
			}
			rand.Shuffle(len(arr), func(i, j int) {
				arr[i], arr[j] = arr[j], arr[i]
			})
			for _, s := range arr {
				filter := bson.M{"doc_id": "test:0"}
				update := bson.M{
					"$addToSet": bson.M{
						fmt.Sprintf("msgs.%d.read_list", index): bson.M{"$each": []string{s}},
					},
				}
				_, err := c.UpdateOne(context.Background(), filter, update)
				if err != nil {
					t.Fatal(err)
				}
			}
		}()
	}

	time.Sleep(time.Second * 2)

	close(ch)

	wg.Wait()

}

func GetDB() *commonMsgDatabase {
	config.Config.Mongo.Address = []string{"192.168.44.128:37017"}
	// config.Config.Mongo.Timeout = 60
	config.Config.Mongo.Database = "openIM"
	// config.Config.Mongo.Source = "admin"
	config.Config.Mongo.Username = "root"
	config.Config.Mongo.Password = "openIM123"
	config.Config.Mongo.MaxPoolSize = 100
	config.Config.RetainChatRecords = 3650
	config.Config.ChatRecordsClearTime = "0 2 * * 3"

	mongo, err := unrelation.NewMongo()
	if err != nil {
		panic(err)
	}
	err = mongo.GetDatabase().Client().Ping(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	return &commonMsgDatabase{
		msgDocDatabase: unrelation.NewMsgMongoDriver(mongo.GetDatabase()),
	}
}

func Test_Insert(t *testing.T) {
	db := GetDB()
	ctx := context.Background()
	var arr []any
	for i := 0; i < 345; i++ {
		if i%2 == 0 {
			arr = append(arr, (*unRelationTb.MsgDataModel)(nil))
			continue
		}
		arr = append(arr, &unRelationTb.MsgDataModel{
			Seq:     int64(i),
			Content: fmt.Sprintf("test-%d", i),
		})
	}
	if err := db.BatchInsertBlock(ctx, "test", arr, updateKeyMsg, 1); err != nil {
		t.Fatal(err)
	}
}

func Test_Revoke(t *testing.T) {
	db := GetDB()
	ctx := context.Background()
	var arr []any
	for i := 0; i < 456; i++ {
		arr = append(arr, &unRelationTb.RevokeModel{
			UserID:   "uid_" + strconv.Itoa(i),
			Nickname: "uname_" + strconv.Itoa(i),
			Time:     time.Now().UnixMilli(),
		})
	}
	if err := db.BatchInsertBlock(ctx, "test", arr, updateKeyRevoke, 123); err != nil {
		t.Fatal(err)
	}
}

func Test_FindBySeq(t *testing.T) {
	if err := log.InitFromConfig("", "", 6, true, false, "", 2); err != nil {
		t.Fatal(err)
	}
	db := GetDB()
	ctx := context.Background()
	fmt.Println(db.msgDocDatabase.(*unrelation.MsgMongoDriver).GetMsgBySeqIndexIn1Doc(ctx, "100", "si_100_101:0", []int64{1}))
	//res, err := db.msgDocDatabase.GetMsgBySeqIndexIn1Doc(ctx, "123456", "test:0", []int64{1, 2, 3})
	//if err != nil {
	//	t.Fatal(err)
	//}
	//db.GetMsgBySeqs(ctx, "100", "si_100_101:0", []int64{6})
	//data, _ := json.Marshal(res)
	//fmt.Println(string(data))
}

//func Test_Delete(t *testing.T) {
//	db := GetDB()
//	ctx := context.Background()
//	var arr []any
//	for i := 0; i < 123; i++ {
//		arr = append(arr, []string{"uid_1", "uid_2"})
//	}
//	if err := db.BatchInsertBlock(ctx, "test", arr, updateKeyDel, 210); err != nil {
//		t.Fatal(err)
//	}
//}

//func Test_Delete1(t *testing.T) {
//	config.Config.Mongo.DBAddress = []string{"192.168.44.128:37017"}
//	config.Config.Mongo.DBTimeout = 60
//	config.Config.Mongo.DBDatabase = "openIM"
//	config.Config.Mongo.DBSource = "admin"
//	config.Config.Mongo.DBUserName = "root"
//	config.Config.Mongo.DBPassword = "openIM123"
//	config.Config.Mongo.DBMaxPoolSize = 100
//	config.Config.Mongo.DBRetainChatRecords = 3650
//	config.Config.Mongo.ChatRecordsClearTime = "0 2 * * 3"
//
//	mongo, err := unrelation.NewMongo()
//	if err != nil {
//		panic(err)
//	}
//	err = mongo.GetDatabase().Client().Ping(context.Background(), nil)
//	if err != nil {
//		panic(err)
//	}
//
//	c := mongo.GetClient().Database("openIM").Collection("msg")
//
//	var o unRelationTb.MsgDocModel
//
//	err = c.FindOne(context.Background(), bson.M{"doc_id": "test:0"}).Decode(&o)
//	if err != nil {
//		panic(err)
//	}
//
//	for i, model := range o.Msg {
//		fmt.Println(i, model == nil)
//	}
//
//}
