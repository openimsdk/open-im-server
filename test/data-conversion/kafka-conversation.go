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
	"fmt"

	"github.com/IBM/sarama"
)

var (
	topic = "latestMsgToRedis"
	addr  = "127.0.0.1:9092"
)

var (
	consumer sarama.Consumer
	producer sarama.SyncProducer
)

func init() {
	// Producer
	config := sarama.NewConfig()            // Instantiate a sarama Config
	config.Producer.Return.Successes = true // Whether to enable the successes channel to be notified after the message is sent successfully
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll        // Set producer Message Reply level 0 1 all
	config.Producer.Partitioner = sarama.NewHashPartitioner // Set the hash-key automatic hash partition. When sending a message, you must specify the key value of the message. If there is no key, the partition will be selected randomly

	client, err := sarama.NewSyncProducer([]string{addr}, config)
	if err != nil {
		fmt.Println("producer closed, err:", err)
	}
	producer = client

	// Consumer
	consumerT, err := sarama.NewConsumer([]string{addr}, sarama.NewConfig())
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
	}
	consumer = consumerT
}

func SendMessage() {
	// construct a message
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder("this is a test log")

	// Send a message
	pid, offset, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println("send msg failed, err:", err)
	}
	fmt.Printf("pid:%v offset:%v\n", pid, offset)
}

func GetMessage() {
	partitionList, err := consumer.Partitions(topic) // Get all partitions according to topic
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
	}
	fmt.Println(partitionList)
	for partition := range partitionList { // iterate over all partitions
		// Create a corresponding partition consumer for each partition
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
		}
		// Asynchronously consume information from each partition
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v", msg.Partition, msg.Offset, msg.Key, msg.Value)
			}
		}(pc)
	}
}
