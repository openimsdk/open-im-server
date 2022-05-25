package mongo

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/db"
	server_api_params "Open_IM/pkg/proto/sdk_ws"
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var (
	Client *mongo.Client
)

func GetUserAllChat(uid string) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	collection := Client.Database(config.Config.Mongo.DBDatabase).Collection("msg")
	var userChatList []db.UserChat
	result, err := collection.Find(context.Background(), bson.M{"uid": primitive.Regex{Pattern: uid}})
	if err != nil {
		fmt.Println("find error", err.Error())
		return
	}
	if err := result.All(ctx, &userChatList); err != nil {
		fmt.Println(err.Error())
	}
	for _, userChat := range userChatList {
		for _, msg := range userChat.Msg {
			msgData := &server_api_params.MsgData{}
			err := proto.Unmarshal(msg.Msg, msgData)
			if err != nil {
				fmt.Println(err.Error(), msg)
				continue
			}
			fmt.Println(*msgData)
		}
	}
}
