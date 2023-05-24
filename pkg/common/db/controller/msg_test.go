package controller

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/unrelation"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
	"sync"
	"testing"
	"time"
)

func Test_BatchInsertChat2DB(t *testing.T) {
	config.Config.Mongo.DBAddress = []string{"192.168.44.128:37017"}
	config.Config.Mongo.DBTimeout = 60
	config.Config.Mongo.DBDatabase = "openIM"
	config.Config.Mongo.DBSource = "admin"
	config.Config.Mongo.DBUserName = "root"
	config.Config.Mongo.DBPassword = "openIM123"
	config.Config.Mongo.DBMaxPoolSize = 100
	config.Config.Mongo.DBRetainChatRecords = 3650
	config.Config.Mongo.ChatRecordsClearTime = "0 2 * * 3"

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

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ch
			for i := 0; i < 500; i++ {
				filter := bson.M{"doc_id": "test:0"}
				update := bson.M{
					"$addToSet": bson.M{
						"msgs.7.del_list": bson.M{"$each": []string{strconv.Itoa(i + 1)}},
					},
				}
				_, err := c.UpdateOne(context.Background(), filter, update)
				if err != nil {
					t.Fatal(err)
				}
			}
		}()
	}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ch
			for i := 0; i < 500; i++ {
				filter := bson.M{"doc_id": "test:0"}
				update := bson.M{
					"$addToSet": bson.M{
						"msgs.7.read_list": bson.M{"$each": []string{strconv.Itoa(200 + i + 1)}},
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
