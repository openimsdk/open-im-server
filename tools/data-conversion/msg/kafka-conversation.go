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

package data_conversion

import (
	"context"
	"encoding/json"
	"fmt"
	. "github.com/OpenIMSDK/Open-IM-Server/tools/conversion/common"
	pbmsg "github.com/OpenIMSDK/Open-IM-Server/tools/conversion/proto/msg"
	"github.com/OpenIMSDK/protocol/constant"
	msgv3 "github.com/OpenIMSDK/protocol/msg"
	"github.com/OpenIMSDK/protocol/sdkws"
	openKeeper "github.com/OpenIMSDK/tools/discoveryregistry/zookeeper"
	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/mw"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
	"time"
)

var consumer sarama.Consumer
var producerV2 sarama.SyncProducer
var wg sync.WaitGroup

var msgRpcClient msgv3.MsgClient

func init() {

	//Producer
	config := sarama.NewConfig()            // Instantiate a sarama Config
	config.Producer.Return.Successes = true // Whether to enable the successes channel to be notified after the message is sent successfully
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll        // Set producer Message Reply level 0 1 all
	config.Producer.Partitioner = sarama.NewHashPartitioner // Set the hash-key automatic hash partition. When sending a message, you must specify the key value of the message. If there is no key, the partition will be selected randomly

	client, err := sarama.NewSyncProducer([]string{KafkaAddr}, config)
	if err != nil {
		fmt.Println("producer closed, err:", err)
	}
	producerV2 = client

	//Consumer
	consumerT, err := sarama.NewConsumer([]string{KafkaAddr}, sarama.NewConfig())
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
	}
	consumer = consumerT

	msgRpcClient = NewMessage()
}

func SendMessage() {
	// construct a message
	msg := &sarama.ProducerMessage{}
	msg.Topic = Topic
	msg.Value = sarama.StringEncoder("this is a test log")

	// Send a message
	pid, offset, err := producerV2.SendMessage(msg)
	if err != nil {
		fmt.Println("send msg failed, err:", err)
	}
	fmt.Printf("pid:%v offset:%v\n", pid, offset)
}

func GetMessage() {
	partitionList, err := consumer.Partitions(Topic) // Get all partitions according to topic
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
	}
	fmt.Println(partitionList)
	//var ch chan int

	if err != nil {
		fmt.Printf("rpc err:%s", err)
	}

	for partition := range partitionList {
		pc, err := consumer.ConsumePartition(Topic, int32(partition), sarama.OffsetOldest)
		if err != nil {
			panic(err)
		}
		wg.Add(1)
		defer pc.AsyncClose()

		go func(sarama.PartitionConsumer) {
			defer wg.Done()
			for msg := range pc.Messages() {
				Transfer([]*sarama.ConsumerMessage{msg})

				//V2
				//msgFromMQV2 := pbmsg.MsgDataToMQ{}
				//err := proto.Unmarshal(msg.Value, &msgFromMQV2)
				//if err != nil {
				//	fmt.Printf("err:%s \n", err)
				//}
				//fmt.Printf("msg:%s \n", &msgFromMQV2)

				//V3
				//msgFromMQ := &sdkws.MsgData{}
				//err = proto.Unmarshal(msg.Value, msgFromMQ)
				//if err != nil {
				//	fmt.Printf("err:%s \n", err)
				//}
				//fmt.Printf("msg:%s \n", &msgFromMQ)
				//fmt.Printf("Partition:%d, Offset:%d, Key:%s, Value:%s\n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
			}

		}(pc)

	}
	wg.Wait()
	consumer.Close()
	//_ = <-ch
}

type TextElem struct {
	Content string `json:"content"`
}

func Transfer(consumerMessages []*sarama.ConsumerMessage) {
	for i := 0; i < len(consumerMessages); i++ {
		fmt.Printf("Partition:%d, Offset:%d, Key:%s \n", consumerMessages[i].Partition, consumerMessages[i].Offset, string(consumerMessages[i].Key))
		msgFromMQV2 := pbmsg.MsgDataToMQ{}
		err := proto.Unmarshal(consumerMessages[i].Value, &msgFromMQV2)
		if err != nil {
			fmt.Printf("err:%s \n", err)
		}
		fmt.Printf("msg:%s \n", &msgFromMQV2)
		//fmt.Printf("rpcClient:%s \n", msgRpcClient)
		if msgFromMQV2.MsgData.ContentType == constant.Text {
			text := string(msgFromMQV2.MsgData.Content)
			textElem := TextElem{
				Content: text,
			}
			msgFromMQV2.MsgData.Content, err = json.Marshal(textElem)
			if err != nil {
				fmt.Printf("test err: %s \n", err)
			}
		}
		if msgFromMQV2.MsgData.SessionType == constant.SingleChatType || msgFromMQV2.MsgData.SessionType == constant.NotificationChatType {
			if string(consumerMessages[i].Key) != msgFromMQV2.MsgData.SendID {
				continue
			}
			offlinePushInfo := &sdkws.OfflinePushInfo{
				Title:         msgFromMQV2.MsgData.OfflinePushInfo.Title,
				Desc:          msgFromMQV2.MsgData.OfflinePushInfo.Desc,
				Ex:            msgFromMQV2.MsgData.OfflinePushInfo.Ex,
				IOSPushSound:  msgFromMQV2.MsgData.OfflinePushInfo.IOSPushSound,
				IOSBadgeCount: msgFromMQV2.MsgData.OfflinePushInfo.IOSBadgeCount,
				SignalInfo:    "",
			}
			msgData := &sdkws.MsgData{
				SendID:           msgFromMQV2.MsgData.SendID,
				RecvID:           msgFromMQV2.MsgData.RecvID,
				GroupID:          msgFromMQV2.MsgData.GroupID,
				ClientMsgID:      msgFromMQV2.MsgData.ClientMsgID,
				ServerMsgID:      msgFromMQV2.MsgData.ServerMsgID,
				SenderPlatformID: msgFromMQV2.MsgData.SenderPlatformID,
				SenderNickname:   msgFromMQV2.MsgData.SenderNickname,
				SenderFaceURL:    msgFromMQV2.MsgData.SenderFaceURL,
				SessionType:      msgFromMQV2.MsgData.SessionType,
				MsgFrom:          msgFromMQV2.MsgData.MsgFrom,
				ContentType:      msgFromMQV2.MsgData.ContentType,
				Content:          msgFromMQV2.MsgData.Content,
				Seq:              int64(msgFromMQV2.MsgData.Seq),
				SendTime:         msgFromMQV2.MsgData.SendTime,
				CreateTime:       msgFromMQV2.MsgData.CreateTime,
				Status:           msgFromMQV2.MsgData.Status,
				IsRead:           false,
				Options:          msgFromMQV2.MsgData.Options,
				OfflinePushInfo:  offlinePushInfo,
				AtUserIDList:     msgFromMQV2.MsgData.AtUserIDList,
				AttachedInfo:     msgFromMQV2.MsgData.AttachedInfo,
				Ex:               msgFromMQV2.MsgData.Ex,
			}
			ctx := context.WithValue(context.Background(), "operationID", msgFromMQV2.OperationID)
			resp, err := msgRpcClient.SendMsg(ctx, &msgv3.SendMsgReq{MsgData: msgData})
			if err != nil {
				fmt.Printf("resp err: %s \n", err)
			}
			fmt.Printf("resp: %s \n", resp)
		} else if msgFromMQV2.MsgData.SessionType == constant.GroupChatType {
			if string(consumerMessages[i].Key) != msgFromMQV2.MsgData.SendID {
				continue
			}
			if msgFromMQV2.MsgData.ContentType < constant.ContentTypeBegin || msgFromMQV2.MsgData.ContentType > constant.AdvancedText {
				continue
			}
			offlinePushInfo := &sdkws.OfflinePushInfo{
				Title:         msgFromMQV2.MsgData.OfflinePushInfo.Title,
				Desc:          msgFromMQV2.MsgData.OfflinePushInfo.Desc,
				Ex:            msgFromMQV2.MsgData.OfflinePushInfo.Ex,
				IOSPushSound:  msgFromMQV2.MsgData.OfflinePushInfo.IOSPushSound,
				IOSBadgeCount: msgFromMQV2.MsgData.OfflinePushInfo.IOSBadgeCount,
				SignalInfo:    "",
			}
			msgData := &sdkws.MsgData{
				SendID:           msgFromMQV2.MsgData.SendID,
				RecvID:           msgFromMQV2.MsgData.RecvID,
				GroupID:          msgFromMQV2.MsgData.GroupID,
				ClientMsgID:      msgFromMQV2.MsgData.ClientMsgID,
				ServerMsgID:      msgFromMQV2.MsgData.ServerMsgID,
				SenderPlatformID: msgFromMQV2.MsgData.SenderPlatformID,
				SenderNickname:   msgFromMQV2.MsgData.SenderNickname,
				SenderFaceURL:    msgFromMQV2.MsgData.SenderFaceURL,
				SessionType:      constant.SuperGroupChatType,
				MsgFrom:          msgFromMQV2.MsgData.MsgFrom,
				ContentType:      msgFromMQV2.MsgData.ContentType,
				Content:          msgFromMQV2.MsgData.Content,
				Seq:              int64(msgFromMQV2.MsgData.Seq),
				SendTime:         msgFromMQV2.MsgData.SendTime,
				CreateTime:       msgFromMQV2.MsgData.CreateTime,
				Status:           msgFromMQV2.MsgData.Status,
				IsRead:           false,
				Options:          msgFromMQV2.MsgData.Options,
				OfflinePushInfo:  offlinePushInfo,
				AtUserIDList:     msgFromMQV2.MsgData.AtUserIDList,
				AttachedInfo:     msgFromMQV2.MsgData.AttachedInfo,
				Ex:               msgFromMQV2.MsgData.Ex,
			}
			ctx := context.WithValue(context.Background(), "operationID", msgFromMQV2.OperationID)
			resp, err := msgRpcClient.SendMsg(ctx, &msgv3.SendMsgReq{MsgData: msgData})
			if err != nil {
				fmt.Printf("resp err: %s \n", err)
			}
			fmt.Printf("resp: %s \n", resp)
		} else if msgFromMQV2.MsgData.SessionType == constant.SuperGroupChatType {
			if msgFromMQV2.MsgData.ContentType < constant.ContentTypeBegin || msgFromMQV2.MsgData.ContentType > constant.AdvancedText {
				continue
			}
			offlinePushInfo := &sdkws.OfflinePushInfo{
				Title:         msgFromMQV2.MsgData.OfflinePushInfo.Title,
				Desc:          msgFromMQV2.MsgData.OfflinePushInfo.Desc,
				Ex:            msgFromMQV2.MsgData.OfflinePushInfo.Ex,
				IOSPushSound:  msgFromMQV2.MsgData.OfflinePushInfo.IOSPushSound,
				IOSBadgeCount: msgFromMQV2.MsgData.OfflinePushInfo.IOSBadgeCount,
				SignalInfo:    "",
			}
			msgData := &sdkws.MsgData{
				SendID:           msgFromMQV2.MsgData.SendID,
				RecvID:           msgFromMQV2.MsgData.RecvID,
				GroupID:          msgFromMQV2.MsgData.GroupID,
				ClientMsgID:      msgFromMQV2.MsgData.ClientMsgID,
				ServerMsgID:      msgFromMQV2.MsgData.ServerMsgID,
				SenderPlatformID: msgFromMQV2.MsgData.SenderPlatformID,
				SenderNickname:   msgFromMQV2.MsgData.SenderNickname,
				SenderFaceURL:    msgFromMQV2.MsgData.SenderFaceURL,
				SessionType:      msgFromMQV2.MsgData.SessionType,
				MsgFrom:          msgFromMQV2.MsgData.MsgFrom,
				ContentType:      msgFromMQV2.MsgData.ContentType,
				Content:          msgFromMQV2.MsgData.Content,
				Seq:              int64(msgFromMQV2.MsgData.Seq),
				SendTime:         msgFromMQV2.MsgData.SendTime,
				CreateTime:       msgFromMQV2.MsgData.CreateTime,
				Status:           msgFromMQV2.MsgData.Status,
				IsRead:           false,
				Options:          msgFromMQV2.MsgData.Options,
				OfflinePushInfo:  offlinePushInfo,
				AtUserIDList:     msgFromMQV2.MsgData.AtUserIDList,
				AttachedInfo:     msgFromMQV2.MsgData.AttachedInfo,
				Ex:               msgFromMQV2.MsgData.Ex,
			}
			ctx := context.WithValue(context.Background(), "operationID", msgFromMQV2.OperationID)
			resp, err := msgRpcClient.SendMsg(ctx, &msgv3.SendMsgReq{MsgData: msgData})
			if err != nil {
				fmt.Printf("resp err: %s \n", err)
			}
			fmt.Printf("resp: %s \n", resp)
		}
		fmt.Printf("\n\n\n")
	}
}

//GetMsgRpcService Convenient for detachment
//func GetMsgRpcService() (rpcclient.MessageRpcClient, error) {
//	client, err := openKeeper.NewClient([]string{ZkAddr}, ZKSchema,
//		openKeeper.WithFreq(time.Hour), openKeeper.WithRoundRobin(), openKeeper.WithUserNameAndPassword(ZKUsername,
//			ZKPassword), openKeeper.WithTimeout(10))
//	msgClient := rpcclient.NewMessageRpcClient(client)
//	if err != nil {
//		return msgClient, err
//	}
//	return msgClient, nil
//}

func NewMessage() msgv3.MsgClient {
	discov, err := openKeeper.NewClient(
		[]string{ZkAddr},
		ZKSchema,
		openKeeper.WithFreq(time.Hour),
		openKeeper.WithRoundRobin(),
		openKeeper.WithUserNameAndPassword(
			ZKUsername,
			ZKPassword),
		openKeeper.WithTimeout(10),
		openKeeper.WithLogger(log.NewZkLogger()),
	)
	if err != nil {
		fmt.Printf("discov, err:%s", err)
	}
	discov.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := discov.GetConn(context.Background(), MsgRpcName)
	if err != nil {
		fmt.Printf("conn, err:%s", err)
		panic(err)
	}
	client := msgv3.NewMsgClient(conn)
	return client
}
