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
	configErrCode         = 6001
)

var (
	cfgPath           = flag.String("c", defaultCfgPath, "Path to the configuration file")
	ErrComponentStart = errs.NewCodeError(componentStartErrCode, "ComponentStartErr")
	ErrConfig         = errs.NewCodeError(configErrCode, "Config file is incorrect")
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

	checks := []checkFunc{
		//{name: "Mysql", function: checkMysql},
		{name: "Mongo", function: checkMongo},
		{name: "Minio", function: checkMinio},
		{name: "Redis", function: checkRedis},
		{name: "Zookeeper", function: checkZookeeper},
		{name: "Kafka", function: checkKafka},
	}

	for i := 0; i < maxRetry; i++ {
		if i != 0 {
			time.Sleep(1 * time.Second)
		}
		fmt.Printf("Checking components Round %v...\n", i+1)

		allSuccess := true
		for _, check := range checks {
			str, err := check.function()
			if err != nil {
				component.ErrorPrint(fmt.Sprintf("Starting %s failed, %v", check.name, err))
				allSuccess = false
				break
			} else {
				component.SuccessPrint(fmt.Sprintf("%s connected successfully, %s", check.name, str))
			}
		}

		if allSuccess {
			component.SuccessPrint("All components started successfully!")

			return
		}
	}
	os.Exit(1)
}

// Helper function to get environment variable or default value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// checkMongo checks the MongoDB connection without retries
func checkMongo() (string, error) {
	mongo := &component.Mongo{
		Address:     config.Config.Mongo.Address,
		Database:    config.Config.Mongo.Database,
		Username:    config.Config.Mongo.Username,
		Password:    config.Config.Mongo.Password,
		MaxPoolSize: config.Config.Mongo.MaxPoolSize,
	}
	uri, uriExist := os.LookupEnv("MONGO_URI")
	if uriExist {
		mongo.URL = uri
	}

	str, err := component.CheckMongo(mongo)
	if err != nil {
		return "", err
	}
	return str, nil
}

// checkMinio checks the MinIO connection
func checkMinio() (string, error) {
	// Check if MinIO is enabled
	if config.Config.Object.Enable != "minio" {
		return "", nil
	}

	// Prioritize environment variables
	endpoint := getEnv("MINIO_ENDPOINT", config.Config.Object.Minio.Endpoint)
	address, addressExist := os.LookupEnv("MINIO_ADDRESS")
	port, portExist := os.LookupEnv("MINIO_PORT")
	if portExist && addressExist {
		endpoint = "http://" + address + ":" + port
	} else if !portExist && addressExist {
		return "", errs.Wrap(errors.New("the MINIO_ADDRESS of minio is empty"))
	} else if portExist && !addressExist {
		return "", errs.Wrap(errors.New("the MINIO_PORT of minio is empty"))
	}

	minio := &component.Minio{
		ApiURL:          config.Config.Object.ApiURL,
		Endpoint:        endpoint,
		AccessKeyID:     getEnv("MINIO_ACCESS_KEY_ID", config.Config.Object.Minio.AccessKeyID),
		SecretAccessKey: getEnv("MINIO_SECRET_ACCESS_KEY", config.Config.Object.Minio.SecretAccessKey),
		SignEndpoint:    config.Config.Object.Minio.SignEndpoint,
		UseSSL:          getEnv("MINIO_USE_SSL", "false"),
	}

	str, err := component.CheckMinio(minio)
	if err != nil {
		return "", err
	}
	return str, nil
}

// checkRedis checks the Redis connection
func checkRedis() (string, error) {
	// Prioritize environment variables
	address := getEnv("REDIS_ADDRESS", strings.Join(config.Config.Redis.Address, ","))
	username := getEnv("REDIS_USERNAME", config.Config.Redis.Username)
	password := getEnv("REDIS_PASSWORD", config.Config.Redis.Password)

	redis := &component.Redis{
		Address:  strings.Split(address, ","),
		Username: username,
		Password: password,
	}

	addr, addrExist := os.LookupEnv("REDIS_ADDRESS")
	port, portExist := os.LookupEnv("REDIS_PORT")

	if addrExist && portExist {
		addresses := strings.Split(addr, ",")
		for i, address := range addresses {
			addresses[i] = address + ":" + port
		}
		redis.Address = addresses
	} else if !addrExist && portExist {
		return "", errs.Wrap(errors.New("the REDIS_ADDRESS of minio is empty"))
	} else if addrExist && !portExist {
		return "", errs.Wrap(errors.New("the REDIS_PORT of minio is empty"))
	}

	str, err := component.CheckRedis(redis)
	if err != nil {
		return "", err
	}
	return str, nil
}

// checkZookeeper checks the Zookeeper connection
func checkZookeeper() (string, error) {
	// Prioritize environment variables

	address := getEnv("ZOOKEEPER_ADDRESS", strings.Join(config.Config.Zookeeper.ZkAddr, ","))

	zk := &component.Zookeeper{
		Schema:   getEnv("ZOOKEEPER_SCHEMA", "digest"),
		ZkAddr:   strings.Split(address, ","),
		Username: getEnv("ZOOKEEPER_USERNAME", config.Config.Zookeeper.Username),
		Password: getEnv("ZOOKEEPER_PASSWORD", config.Config.Zookeeper.Password),
	}

	address, addrExist := os.LookupEnv("ZOOKEEPER_ADDRESS")
	port, portExist := os.LookupEnv("ZOOKEEPER_PORT")
	if addrExist && portExist {
		addresses := strings.Split(address, ",")
		for i, addr := range addresses {
			addresses[i] = addr + ":" + port
		}
		zk.ZkAddr = addresses
	} else if !addrExist && portExist {
		return "", errs.Wrap(errors.New("the ZOOKEEPER_ADDRESS of minio is empty"))
	} else if addrExist && !portExist {
		return "", errs.Wrap(errors.New("the ZOOKEEPER_PORT of minio is empty"))
	}

	str, err := component.CheckZookeeper(zk)
	if err != nil {
		return "", err
	}
	return str, nil
}

// checkKafka checks the Kafka connection
func checkKafka() (string, error) {
	// Prioritize environment variables
	username := getEnv("KAFKA_USERNAME", config.Config.Kafka.Username)
	password := getEnv("KAFKA_PASSWORD", config.Config.Kafka.Password)
	address := getEnv("KAFKA_ADDRESS", strings.Join(config.Config.Kafka.Addr, ","))

	kafka := &component.Kafka{
		Username: username,
		Password: password,
		Addr:     strings.Split(address, ","),
	}

	address, addrExist := os.LookupEnv("KAFKA_ADDRESS")
	port, portExist := os.LookupEnv("KAFKA_PORT")
	if addrExist && portExist {
		addresses := strings.Split(address, ",")
		for i, addr := range addresses {
			addresses[i] = addr + ":" + port
		}
		kafka.Addr = addresses
	} else if !addrExist && portExist {
		return "", errs.Wrap(errors.New("the KAFKA_ADDRESS of minio is empty"))
	} else if addrExist && !portExist {
		return "", errs.Wrap(errors.New("the KAFKA_PORT of minio is empty"))
	}

	str, kafkaClient, err := component.CheckKafka(kafka)
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
