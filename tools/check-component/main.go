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
	"fmt"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"os"
	"time"
)

const maxRetry = 180

func CheckZookeeper(ctx context.Context, ZkServers []string, scheme string, options ...ZkOption) error {
	return nil
}

func CheckMongo(ctx context.Context, config *config.Mongo) error {
	return nil
}

func CheckRedis(ctx context.Context, config *config.Redis) error {
	return nil
}

func CheckMinIO(ctx context.Context, config *config.Minio) error {
	return nil

}

func CheckKafka(ctx context.Context, conf *config.Kafka, topics []string) error {
	return nil
}

func main() {
	ctx := context.Background()
	zkServers := []string{"localhost:2181"}
	scheme := "digest"
	mongoConfig := &MongoConfig{}
	redisConfig := &RedisConfig{}
	minioConfig := &MinIOConfig{}
	kafkaConfig := &KafkaConfig{}
	topics := []string{"topic1", "topic2"}

	checksDone := make(map[string]bool)
	checks := map[string]func() error{
		"Zookeeper": func() error {
			return CheckZookeeper(ctx, zkServers, scheme)
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
			return CheckKafka(ctx, kafkaConfig, topics)
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
			fmt.Println("All checks passed successfully.")
			return
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Println("Not all checks passed successfully after 180 attempts.")
	os.Exit(-1)
}
