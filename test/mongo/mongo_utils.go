// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mongo

import (
	"Open_IM/pkg/common/config"
	server_api_params "Open_IM/pkg/proto/sdk_ws"
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var (
	Client *mongo.Client
)

type MsgInfo struct {
	SendTime int64
	Msg      []byte
}

type UserChat struct {
	UID string
	Msg []MsgInfo
}

func GetUserAllChat(uid string) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
	collection := Client.Database(config.Config.Mongo.DBDatabase).Collection("msg")
	var userChatList []UserChat
	uid = uid + ":"
	filter := bson.M{"uid": bson.M{"$regex": uid}}
	//filter := bson.M{"uid": "17726378428:0"}
	result, err := collection.Find(context.Background(), filter)
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
			fmt.Println("seq: ", msgData.Seq, "status: ", msgData.Status,
				"sendID: ", msgData.SendID, "recvID: ", msgData.RecvID,
				"sendTime: ", msgData.SendTime,
				"clientMsgID: ", msgData.ClientMsgID,
				"serverMsgID: ", msgData.ServerMsgID,
				"content: ", string(msgData.Content))
		}
	}
}
