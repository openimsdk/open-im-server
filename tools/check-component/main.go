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
	"github.com/openimsdk/open-im-server/v3/pkg/common/cmd"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/redisutil"
	"github.com/openimsdk/tools/discovery/zookeeper"
	"github.com/openimsdk/tools/mq/kafka"
	"github.com/openimsdk/tools/s3/minio"
	"github.com/openimsdk/tools/system/program"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

const maxRetry = 180

func CheckZookeeper(ctx context.Context, config *config.ZooKeeper) error {
	// Temporary disable logging
	originalLogger := log.Default().Writer()
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(originalLogger) // Ensure logging is restored
	return zookeeper.Check(ctx, config.Address, config.Schema, zookeeper.WithUserNameAndPassword(config.Username, config.Password))
}

func CheckMongo(ctx context.Context, config *config.Mongo) error {
	return mongoutil.Check(ctx, config.Build())
}

func CheckRedis(ctx context.Context, config *config.Redis) error {
	return redisutil.Check(ctx, config.Build())
}

func CheckMinIO(ctx context.Context, config *config.Minio) error {
	return minio.Check(ctx, config.Build())
}

func CheckKafka(ctx context.Context, conf *config.Kafka) error {
	return kafka.Check(ctx, conf.Build(), []string{conf.ToMongoTopic, conf.ToRedisTopic, conf.ToPushTopic})
}

func initConfig(configDir string) (*config.Mongo, *config.Redis, *config.Kafka, *config.Minio, *config.ZooKeeper, error) {
	var (
		mongoConfig     = &config.Mongo{}
		redisConfig     = &config.Redis{}
		kafkaConfig     = &config.Kafka{}
		minioConfig     = &config.Minio{}
		zookeeperConfig = &config.ZooKeeper{}
	)
	err := config.LoadConfig(filepath.Join(configDir, cmd.MongodbConfigFileName), cmd.ConfigEnvPrefixMap[cmd.MongodbConfigFileName], mongoConfig)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	err = config.LoadConfig(filepath.Join(configDir, cmd.RedisConfigFileName), cmd.ConfigEnvPrefixMap[cmd.RedisConfigFileName], redisConfig)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	err = config.LoadConfig(filepath.Join(configDir, cmd.KafkaConfigFileName), cmd.ConfigEnvPrefixMap[cmd.KafkaConfigFileName], kafkaConfig)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	err = config.LoadConfig(filepath.Join(configDir, cmd.MinioConfigFileName), cmd.ConfigEnvPrefixMap[cmd.MinioConfigFileName], minioConfig)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	err = config.LoadConfig(filepath.Join(configDir, cmd.ZookeeperConfigFileName), cmd.ConfigEnvPrefixMap[cmd.ZookeeperConfigFileName], zookeeperConfig)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	return mongoConfig, redisConfig, kafkaConfig, minioConfig, zookeeperConfig, nil
}

func main() {
	var index int
	var configDir string
	flag.IntVar(&index, "i", 0, "Index number")
	defaultConfigDir := filepath.Join("..", "..", "..", "..", "..", "config")
	flag.StringVar(&configDir, "c", defaultConfigDir, "Configuration dir")
	flag.Parse()

	fmt.Printf("%s Index: %d, Config Path: %s\n", filepath.Base(os.Args[0]), index, configDir)

	mongoConfig, redisConfig, kafkaConfig, minioConfig, zookeeperConfig, err := initConfig(configDir)
	if err != nil {
		program.ExitWithError(err)
	}

	ctx := context.Background()
	err = performChecks(ctx, mongoConfig, redisConfig, kafkaConfig, minioConfig, zookeeperConfig, maxRetry)
	if err != nil {
		// Assume program.ExitWithError logs the error and exits.
		// Replace with your error handling logic as necessary.
		program.ExitWithError(err)
	}
}

func performChecks(ctx context.Context, mongoConfig *config.Mongo, redisConfig *config.Redis, kafkaConfig *config.Kafka, minioConfig *config.Minio, zookeeperConfig *config.ZooKeeper, maxRetry int) error {
	checksDone := make(map[string]bool)

	checks := map[string]func() error{
		"Zookeeper": func() error {
			return CheckZookeeper(ctx, zookeeperConfig)
		},
		"Mongo": func() error {
			return CheckMongo(ctx, mongoConfig)
		},
		"Redis": func() error {
			return CheckRedis(ctx, redisConfig)
		},
		"MinIO": func() error {
			return CheckMinIO(ctx, minioConfig)
		},
		"Kafka": func() error {
			return CheckKafka(ctx, kafkaConfig)
		},
	}

	for i := 0; i < maxRetry; i++ {
		allSuccess := true
		for name, check := range checks {
			if !checksDone[name] {
				if err := check(); err != nil {
					fmt.Printf("%s check failed: %v\n", name, err)
					allSuccess = false
				} else {
					fmt.Printf("%s check succeeded.\n", name)
					checksDone[name] = true
				}
			}
		}

		if allSuccess {
			fmt.Println("All components checks passed successfully.")
			return nil
		}

		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("not all components checks passed successfully after %d attempts", maxRetry)
}
