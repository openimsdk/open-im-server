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
	"context"
	"flag"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/OpenIMSDK/tools/log"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/s3/cos"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/s3/minio"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/s3/oss"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/unrelation"
	"github.com/openimsdk/open-im-server/v3/pkg/common/discoveryregister/zookeeper"
	"github.com/openimsdk/open-im-server/v3/pkg/common/kafka"
	"os"
	"strings"
	"time"

	"github.com/OpenIMSDK/tools/component"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"

	"gopkg.in/yaml.v3"
)

const (
	// defaultCfgPath is the default path of the configuration file.
	defaultCfgPath        = "../../../../../config/config.yaml"
	maxRetry              = 300
	componentStartErrCode = 6000
)

var (
	cfgPath           = flag.String("c", defaultCfgPath, "Path to the configuration file")
	ErrComponentStart = errs.NewCodeError(componentStartErrCode, "ComponentStartErr")
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
	function func() (string, error)
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

		var (
			err     error
			errInfo string
		)
		allSuccess := true
		disruptions := true
		for _, check := range checks {
			errInfo, err = check.function()
			if err != nil {
				component.ErrorPrint(fmt.Sprintf("Starting %s failed, %v, the conneted info is:%s", check.name, err, errInfo))
				log.ZError(context.Background(), errInfo, err)
				allSuccess = false
				break
			} else {
				component.SuccessPrint(fmt.Sprintf("%s connected successfully, the addr is:%s", check.name, errInfo))
				log.ZError(context.Background(), errInfo, err)
			}
			if check.name == "kafka" && errs.Unwrap(err) == ErrComponentStart {
				disruptions = false
			}
		}

		if allSuccess {
			component.SuccessPrint("All components started successfully!")
			log.ZInfo(context.Background(), errInfo, err)
			return
		}

		if disruptions {
			component.ErrorPrint(fmt.Sprintf("component check exit,err:  %v", err))
			return
		}
	}
}

// checkMongo checks the MongoDB connection without retries
func checkMongo() (string, error) {
	_, err := unrelation.NewMongo()
	if err != nil {
		if config.Config.Mongo.Uri != "" {
			return config.Config.Mongo.Uri, err
		}
		uriFormat := "mongodb://%s/%s?maxPoolSize=%s"
		if config.Config.Mongo.Username != "" && config.Config.Mongo.Password != "" {
			uriFormat = "mongodb://%s:%s@%s/%s?maxPoolSize=%s"
			return fmt.Sprintf(uriFormat, config.Config.Mongo.Username, config.Config.Mongo.Password, config.Config.Mongo.Address, config.Config.Mongo.Database, config.Config.Mongo.MaxPoolSize), err
		}
		return fmt.Sprintf(uriFormat, config.Config.Mongo.Address, config.Config.Mongo.Database, config.Config.Mongo.MaxPoolSize), err
	}
	return strings.Join(config.Config.Mongo.Address, ","), nil
}

// checkRedis checks the Redis connection
func checkRedis() (string, error) {
	_, err := cache.NewRedis()
	if err != nil {
		uriFormat := "The username is:%s, the password is:%s, the address is:%s, the clusterMode is:%t"
		return fmt.Sprintf(uriFormat, config.Config.Redis.Username, config.Config.Redis.Password, config.Config.Redis.Address, config.Config.Redis.ClusterMode), err
	}
	return strings.Join(config.Config.Redis.Address, ","), err
}

// checkMinio checks the MinIO connection
func checkMinio() (string, error) {

	rdb, err := cache.NewRedis()

	enable := config.Config.Object.Enable
	switch config.Config.Object.Enable {
	case "minio":
		_, err = minio.NewMinio(cache.NewMinioCache(rdb))
	case "cos":
		_, err = cos.NewCos()
	case "oss":
		_, err = oss.NewOSS()
	default:
		err = fmt.Errorf("invalid object enable: %s", enable)
	}
	if err != nil {
		uriFormat := "The apiURL is:%s, the endpoint is:%s, the signEndpoint is:%s."
		return fmt.Sprintf(uriFormat, config.Config.Object.ApiURL, config.Config.Object.Minio.Endpoint, config.Config.Object.Minio.SignEndpoint), err
	}
	return config.Config.Object.Minio.Endpoint, nil
}

// checkZookeeper checks the Zookeeper connection
func checkZookeeper() (string, error) {
	_, err := zookeeper.NewZookeeperDiscoveryRegister()
	if err != nil {
		if config.Config.Zookeeper.Username != "" && config.Config.Zookeeper.Password != "" {
			return fmt.Sprintf("The addr is:%s,the schema is:%s, the username is:%s, the password is:%s.", config.Config.Zookeeper.ZkAddr, config.Config.Zookeeper.Schema, config.Config.Zookeeper.Username, config.Config.Zookeeper.Password), err
		}
		return fmt.Sprintf("The addr is:%s,the schema is:%s", config.Config.Zookeeper.ZkAddr, config.Config.Zookeeper.Schema), err
	}
	return strings.Join(config.Config.Zookeeper.ZkAddr, ","), nil
}

// checkKafka checks the Kafka connection
func checkKafka() (string, error) {

	// Prioritize environment variables
	kafkaStu := &component.Kafka{
		Username: config.Config.Kafka.Username,
		Password: config.Config.Kafka.Password,
		Addr:     config.Config.Kafka.Addr,
	}

	str, kafkaClient, err := component.CheckKafka(kafkaStu)
	if err != nil {
		return "", err
	}
	defer kafkaClient.Close()

	// Verify if necessary topics exist
	topics, err := kafkaClient.Topics()
	if err != nil {
		return "", errs.Wrap(err)
	}

	requiredTopics := []string{
		config.Config.Kafka.MsgToMongo.Topic,
		config.Config.Kafka.MsgToPush.Topic,
		config.Config.Kafka.LatestMsgToRedis.Topic,
	}

	for _, requiredTopic := range requiredTopics {
		if !isTopicPresent(requiredTopic, topics) {
			return "", ErrComponentStart.Wrap(fmt.Sprintf("Kafka doesn't contain topic: %v", requiredTopic))
		}
	}

	kafka.NewMConsumerGroup(&kafka.MConsumerGroupConfig{
		KafkaVersion:   sarama.V2_0_0_0,
		OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false,
	}, []string{config.Config.Kafka.LatestMsgToRedis.Topic},
		config.Config.Kafka.Addr, config.Config.Kafka.ConsumerGroupID.MsgToRedis)

	kafka.NewMConsumerGroup(&kafka.MConsumerGroupConfig{
		KafkaVersion:   sarama.V2_0_0_0,
		OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false,
	}, []string{config.Config.Kafka.MsgToMongo.Topic},
		config.Config.Kafka.Addr, config.Config.Kafka.ConsumerGroupID.MsgToMongo)

	kafka.NewMConsumerGroup(&kafka.MConsumerGroupConfig{
		KafkaVersion:   sarama.V2_0_0_0,
		OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false,
	}, []string{config.Config.Kafka.MsgToPush.Topic}, config.Config.Kafka.Addr,
		config.Config.Kafka.ConsumerGroupID.MsgToPush)

	return str, nil
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
	config.Config.Mongo.Uri = getEnv("MONGO_URI", config.Config.Mongo.Uri)
	config.Config.Mongo.Username = getEnv("MONGO_OPENIM_USERNAME", config.Config.Mongo.Username)
	config.Config.Mongo.Password = getEnv("MONGO_OPENIM_PASSWORD", config.Config.Mongo.Password)
	config.Config.Kafka.Username = getEnv("KAFKA_USERNAME", config.Config.Kafka.Username)
	config.Config.Kafka.Password = getEnv("KAFKA_PASSWORD", config.Config.Kafka.Password)
	config.Config.Kafka.Addr = strings.Split(getEnv("KAFKA_ADDRESS", strings.Join(config.Config.Kafka.Addr, ",")), ",")

}

// Helper function to get environment variable or default value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
