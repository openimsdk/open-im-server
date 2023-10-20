package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/IBM/sarama"
	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/protocol/msg"
	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/tools/mw"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/openimsdk/open-im-server/v3/pkg/apistruct"
	pbmsg "github.com/openimsdk/open-im-server/v3/tools/data-conversion/openim/proto/msg"
)

func main() {

	var (
		topic       = "ws2ms_chat"      // v2版本配置文件kafka.topic.ws2ms_chat
		kafkaAddr   = "127.0.0.1:9092"  // v2版本配置文件kafka.topic.addr
		rpcAddr     = "127.0.0.1:10130" // v3版本配置文件rpcPort.openImMessagePort
		adminUserID = "openIM123456"    // v3版本管理员userID
		concurrency = 1                 // 并发数量
	)

	getRpcConn := func() (*grpc.ClientConn, error) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		return grpc.DialContext(ctx, rpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), mw.GrpcClient())
	}
	conn, err := getRpcConn()
	if err != nil {
		log.Println("get rpc conn", err)
		return
	}
	defer conn.Close()

	msgClient := msg.NewMsgClient(conn)

	conf := sarama.NewConfig()
	conf.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumer([]string{kafkaAddr}, conf)
	if err != nil {
		log.Println("kafka consumer conn", err)
		return
	}
	partitions, err := consumer.Partitions(topic) // Get all partitions according to topic
	if err != nil {
		log.Println("kafka partitions", err)
		return
	}

	if len(partitions) == 0 {
		log.Println("kafka partitions is empty")
		return
	}
	log.Println("kafka partitions", partitions)

	msgCh := make(chan *pbmsg.MsgDataToMQ, concurrency*2)

	var kfkWg sync.WaitGroup

	distinct := make(map[string]struct{})
	var lock sync.Mutex

	for _, partition := range partitions {
		kfkWg.Add(1)
		go func(partition int32) {
			defer kfkWg.Done()
			pc, err := consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
			if err != nil {
				log.Printf("kafka Consume Partition %d failed %s\n", partition, err)
				return
			}
			defer pc.Close()
			ch := pc.Messages()
			for {
				select {
				case <-time.After(time.Second * 10): // 10s读取不到就关闭
					return
				case message, ok := <-ch:
					if !ok {
						return
					}
					msgFromMQV2 := pbmsg.MsgDataToMQ{}
					err := proto.Unmarshal(message.Value, &msgFromMQV2)
					if err != nil {
						log.Printf("kafka msg partition %d offset %d unmarshal failed %s\n", message.Partition, message.Offset, message.Value)
						continue
					}
					if msgFromMQV2.MsgData == nil || msgFromMQV2.OperationID == "" {
						continue
					}
					if msgFromMQV2.MsgData.ContentType < constant.ContentTypeBegin || msgFromMQV2.MsgData.ContentType > constant.AdvancedText {
						continue
					}
					lock.Lock()
					_, exist := distinct[msgFromMQV2.MsgData.ClientMsgID]
					if !exist {
						distinct[msgFromMQV2.MsgData.ClientMsgID] = struct{}{}
					}
					lock.Unlock()
					if exist {
						continue
					}
					msgCh <- &msgFromMQV2
				}
			}
		}(partition)
	}

	go func() {
		kfkWg.Wait()
		close(msgCh)
	}()

	var msgWg sync.WaitGroup

	var (
		success int64
		failed  int64
	)
	for i := 0; i < concurrency; i++ {
		msgWg.Add(1)
		go func() {
			defer msgWg.Done()
			for message := range msgCh {
				HandlerV2Msg(msgClient, adminUserID, message, &success, &failed)
			}
		}()
	}

	msgWg.Wait()
	log.Printf("total %d success %d failed %d\n", success+failed, success, failed)
}

func HandlerV2Msg(msgClient msg.MsgClient, adminUserID string, msgFromMQV2 *pbmsg.MsgDataToMQ, success *int64, failed *int64) {
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
		SendTime:         msgFromMQV2.MsgData.SendTime,
		CreateTime:       msgFromMQV2.MsgData.CreateTime,
		Status:           msgFromMQV2.MsgData.Status,
		IsRead:           false,
		Options:          msgFromMQV2.MsgData.Options,
		AtUserIDList:     msgFromMQV2.MsgData.AtUserIDList,
		AttachedInfo:     msgFromMQV2.MsgData.AttachedInfo,
		Ex:               msgFromMQV2.MsgData.Ex,
	}

	if msgFromMQV2.MsgData.OfflinePushInfo != nil {
		msgData.OfflinePushInfo = &sdkws.OfflinePushInfo{
			Title:         msgFromMQV2.MsgData.OfflinePushInfo.Title,
			Desc:          msgFromMQV2.MsgData.OfflinePushInfo.Desc,
			Ex:            msgFromMQV2.MsgData.OfflinePushInfo.Ex,
			IOSPushSound:  msgFromMQV2.MsgData.OfflinePushInfo.IOSPushSound,
			IOSBadgeCount: msgFromMQV2.MsgData.OfflinePushInfo.IOSBadgeCount,
			SignalInfo:    "",
		}
	}
	switch msgData.ContentType {
	case constant.Text:
		data, err := json.Marshal(apistruct.TextElem{
			Content: string(msgFromMQV2.MsgData.Content),
		})
		if err != nil {
			return
		}
		msgData.Content = data
	default:
		msgData.Content = msgFromMQV2.MsgData.Content
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	ctx = context.WithValue(context.Background(), constant.OperationID, msgFromMQV2.OperationID)
	ctx = context.WithValue(ctx, constant.OpUserID, adminUserID)

	resp, err := msgClient.SendMsg(ctx, &msg.SendMsgReq{MsgData: msgData})
	if err != nil {
		atomic.AddInt64(failed, 1)
		log.Printf("send msg %+v failed %s\n", msgData, err)
		return
	}
	atomic.AddInt64(success, 1)
	log.Printf("send msg success %+v resp %+v\n", msgData, resp)
}
