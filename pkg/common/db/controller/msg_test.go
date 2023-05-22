package controller

import (
	"context"
	"fmt"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/unrelation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
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
	ctx := context.Background()

	msgs := make([]*sdkws.MsgData, 0, 15000)

	for i := 0; i < cap(msgs); i++ {
		msgs = append(msgs, &sdkws.MsgData{
			Content:  []byte(fmt.Sprintf("test-%d", i)),
			SendTime: time.Now().UnixMilli(),
		})
	}
	err = db.BatchInsertChat2DB(ctx, "test", msgs, 4999)
	if err != nil {
		panic(err)
	}

}
