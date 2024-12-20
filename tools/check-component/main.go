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
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/cmd"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/redisutil"
	"github.com/openimsdk/tools/discovery/etcd"
	"github.com/openimsdk/tools/discovery/zookeeper"
	"github.com/openimsdk/tools/mq/kafka"
	"github.com/openimsdk/tools/s3/minio"
	"github.com/openimsdk/tools/system/program"
	"github.com/openimsdk/tools/utils/runtimeenv"
)

const maxRetry = 180

const (
	MountConfigFilePath = "CONFIG_PATH"
	DeploymentType      = "DEPLOYMENT_TYPE"
	KUBERNETES          = "kubernetes"
)

func CheckZookeeper(ctx context.Context, config *config.ZooKeeper) error {
	// Temporary disable logging
	originalLogger := log.Default().Writer()
	log.SetOutput(io.Discard)
	defer log.SetOutput(originalLogger) // Ensure logging is restored
	return zookeeper.Check(ctx, config.Address, config.Schema, zookeeper.WithUserNameAndPassword(config.Username, config.Password))
}

func CheckEtcd(ctx context.Context, config *config.Etcd) error {
	return etcd.Check(ctx, config.Address, "/check_openim_component",
		true,
		etcd.WithDialTimeout(10*time.Second),
		etcd.WithMaxCallSendMsgSize(20*1024*1024),
		etcd.WithUsernameAndPassword(config.Username, config.Password))
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
	return kafka.CheckHealth(ctx, conf.Build())
}

func initConfig(configDir string) (*config.Mongo, *config.Redis, *config.Kafka, *config.Minio, *config.Discovery, error) {
	var (
		mongoConfig = &config.Mongo{}
		redisConfig = &config.Redis{}
		kafkaConfig = &config.Kafka{}
		minioConfig = &config.Minio{}
		discovery   = &config.Discovery{}
		thirdConfig = &config.Third{}
	)
	runtimeEnv := runtimeenv.PrintRuntimeEnvironment()

	err := config.Load(configDir, config.MongodbConfigFileName, cmd.ConfigEnvPrefixMap[config.MongodbConfigFileName], runtimeEnv, mongoConfig)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	err = config.Load(configDir, config.RedisConfigFileName, cmd.ConfigEnvPrefixMap[config.RedisConfigFileName], runtimeEnv, redisConfig)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	err = config.Load(configDir, config.KafkaConfigFileName, cmd.ConfigEnvPrefixMap[config.KafkaConfigFileName], runtimeEnv, kafkaConfig)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	err = config.Load(configDir, config.OpenIMRPCThirdCfgFileName, cmd.ConfigEnvPrefixMap[config.OpenIMRPCThirdCfgFileName], runtimeEnv, thirdConfig)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	if thirdConfig.Object.Enable == "minio" {
		err = config.Load(configDir, config.MinioConfigFileName, cmd.ConfigEnvPrefixMap[config.MinioConfigFileName], runtimeEnv, minioConfig)
		if err != nil {
			return nil, nil, nil, nil, nil, err
		}
	} else {
		minioConfig = nil
	}
	err = config.Load(configDir, config.DiscoveryConfigFilename, cmd.ConfigEnvPrefixMap[config.DiscoveryConfigFilename], runtimeEnv, discovery)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	return mongoConfig, redisConfig, kafkaConfig, minioConfig, discovery, nil
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

func performChecks(ctx context.Context, mongoConfig *config.Mongo, redisConfig *config.Redis, kafkaConfig *config.Kafka, minioConfig *config.Minio, discovery *config.Discovery, maxRetry int) error {
	checksDone := make(map[string]bool)

	checks := map[string]func(ctx context.Context) error{
		"Mongo": func(ctx context.Context) error {
			return CheckMongo(ctx, mongoConfig)
		},
		"Redis": func(ctx context.Context) error {
			return CheckRedis(ctx, redisConfig)
		},
		"Kafka": func(ctx context.Context) error {
			return CheckKafka(ctx, kafkaConfig)
		},
	}
	if minioConfig != nil {
		checks["MinIO"] = func(ctx context.Context) error {
			return CheckMinIO(ctx, minioConfig)
		}
	}
	if discovery.Enable == "etcd" {
		checks["Etcd"] = func(ctx context.Context) error {
			return CheckEtcd(ctx, &discovery.Etcd)
		}
	}

	for i := 0; i < maxRetry; i++ {
		allSuccess := true
		for name, check := range checks {
			if !checksDone[name] {
				if err := check(ctx); err != nil {
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
