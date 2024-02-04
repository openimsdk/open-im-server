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

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/IBM/sarama"

	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/unrelation"
	"github.com/openimsdk/open-im-server/v3/pkg/common/discoveryregister/zookeeper"
	"github.com/openimsdk/open-im-server/v3/pkg/common/kafka"

	"github.com/OpenIMSDK/tools/component"
	"github.com/OpenIMSDK/tools/errs"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"

	"gopkg.in/yaml.v3"
)

const (
	// defaultCfgPath is the default path of the configuration file.
	defaultCfgPath = "../../../../../config/config.yaml"
	maxRetry       = 300
)

var (
	cfgPath = flag.String("c", defaultCfgPath, "Path to the configuration file")
)

func initCfg() error {
	data, err := os.ReadFile(*cfgPath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, &config.Config)
}

type checkFunc struct {
	name     string
	function func() error
	flag     bool
}

func main() {
	flag.Parse()

	if err := initCfg(); err != nil {
		fmt.Printf("Read config failed: %v\n", err)

		return
	}

	configGetEnv()

	checks := []checkFunc{
		//{name: "Mysql", function: checkMysql},
		{name: "Mongo", function: checkMongo},
		{name: "Redis", function: checkRedis},
		{name: "Minio", function: checkMinio},
		{name: "Zookeeper", function: checkZookeeper},
		{name: "Kafka", function: checkKafka},
	}

	for i := 0; i < maxRetry; i++ {
		if i != 0 {
			time.Sleep(1 * time.Second)
		}
		fmt.Printf("Checking components Round %v...\n", i+1)

		var err error
		allSuccess := true
		for index, check := range checks {
			if !check.flag {
				err = check.function()
				if err != nil {
					component.ErrorPrint(fmt.Sprintf("Starting %s failed:%v.", check.name, err))
					allSuccess = false

				} else {
					checks[index].flag = true
					component.SuccessPrint(fmt.Sprintf("%s connected successfully", check.name))
				}
			}
		}

		if allSuccess {
			component.SuccessPrint("All components started successfully!")
			return
		}
	}
}

// checkMongo checks the MongoDB connection without retries
func checkMongo() error {
	_, err := unrelation.NewMongo()
	return err
}

// checkRedis checks the Redis connection
func checkRedis() error {
	_, err := cache.NewRedis()
	return err
}

// checkMinio checks the MinIO connection
func checkMinio() error {

	// Check if MinIO is enabled
	if config.Config.Object.Enable != "minio" {
		return errs.Wrap(errors.New("minio.Enable is empty"))
	}
	minio := &component.Minio{
		ApiURL:          config.Config.Object.ApiURL,
		Endpoint:        config.Config.Object.Minio.Endpoint,
		AccessKeyID:     config.Config.Object.Minio.AccessKeyID,
		SecretAccessKey: config.Config.Object.Minio.SecretAccessKey,
		SignEndpoint:    config.Config.Object.Minio.SignEndpoint,
		UseSSL:          getEnv("MINIO_USE_SSL", "false"),
	}
	err := component.CheckMinio(minio)
	return err
}

// checkZookeeper checks the Zookeeper connection
func checkZookeeper() error {
	_, err := zookeeper.NewZookeeperDiscoveryRegister()
	return err
}

// checkKafka checks the Kafka connection
func checkKafka() error {
	// Prioritize environment variables
	kafkaStu := &component.Kafka{
		Username: config.Config.Kafka.Username,
		Password: config.Config.Kafka.Password,
		Addr:     config.Config.Kafka.Addr,
	}

	kafkaClient, err := component.CheckKafka(kafkaStu)
	if err != nil {
		return err
	}
	defer kafkaClient.Close()

	// Verify if necessary topics exist
	topics, err := kafkaClient.Topics()
	if err != nil {
		return errs.Wrap(err)
	}

	requiredTopics := []string{
		config.Config.Kafka.MsgToMongo.Topic,
		config.Config.Kafka.MsgToPush.Topic,
		config.Config.Kafka.LatestMsgToRedis.Topic,
	}

	for _, requiredTopic := range requiredTopics {
		if !isTopicPresent(requiredTopic, topics) {
			return errs.Wrap(err, fmt.Sprintf("Kafka doesn't contain topic: %v", requiredTopic))
		}
	}

	_, err = kafka.NewMConsumerGroup(&kafka.MConsumerGroupConfig{
		KafkaVersion:   sarama.V2_0_0_0,
		OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false,
	}, []string{config.Config.Kafka.LatestMsgToRedis.Topic},
		config.Config.Kafka.Addr, config.Config.Kafka.ConsumerGroupID.MsgToRedis)
	if err != nil {
		return err
	}

	_, err = kafka.NewMConsumerGroup(&kafka.MConsumerGroupConfig{
		KafkaVersion:   sarama.V2_0_0_0,
		OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false,
	}, []string{config.Config.Kafka.MsgToPush.Topic},
		config.Config.Kafka.Addr, config.Config.Kafka.ConsumerGroupID.MsgToMongo)
	if err != nil {
		return err
	}

	kafka.NewMConsumerGroup(&kafka.MConsumerGroupConfig{
		KafkaVersion:   sarama.V2_0_0_0,
		OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false,
	}, []string{config.Config.Kafka.MsgToPush.Topic}, config.Config.Kafka.Addr,
		config.Config.Kafka.ConsumerGroupID.MsgToPush)
	if err != nil {
		return err
	}

	return nil
}

// isTopicPresent checks if a topic is present in the list of topics
func isTopicPresent(topic string, topics []string) bool {
	for _, t := range topics {
		if t == topic {
			return true
		}
	}
	return false
}

func configGetEnv() {
	config.Config.Object.Minio.AccessKeyID = getEnv("MINIO_ACCESS_KEY_ID", config.Config.Object.Minio.AccessKeyID)
	config.Config.Object.Minio.SecretAccessKey = getEnv("MINIO_SECRET_ACCESS_KEY", config.Config.Object.Minio.SecretAccessKey)
	config.Config.Mongo.Uri = getEnv("MONGO_URI", config.Config.Mongo.Uri)
	config.Config.Mongo.Username = getEnv("MONGO_OPENIM_USERNAME", config.Config.Mongo.Username)
	config.Config.Mongo.Password = getEnv("MONGO_OPENIM_PASSWORD", config.Config.Mongo.Password)
	config.Config.Kafka.Username = getEnv("KAFKA_USERNAME", config.Config.Kafka.Username)
	config.Config.Kafka.Password = getEnv("KAFKA_PASSWORD", config.Config.Kafka.Password)
	config.Config.Kafka.Addr = strings.Split(getEnv("KAFKA_ADDRESS", strings.Join(config.Config.Kafka.Addr, ",")), ",")
	config.Config.Object.Minio.Endpoint = getMinioAddr("MINIO_ENDPOINT", "MINIO_ADDRESS", "MINIO_PORT", config.Config.Object.Minio.Endpoint)
}

func getMinioAddr(key1, key2, key3, fallback string) string {
	// Prioritize environment variables
	endpoint := getEnv(key1, fallback)
	address, addressExist := os.LookupEnv(key2)
	port, portExist := os.LookupEnv(key3)
	if portExist && addressExist {
		endpoint = "http://" + address + ":" + port
		return endpoint
	}
	return endpoint
}

// Helper function to get environment variable or default value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
