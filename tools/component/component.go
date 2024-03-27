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
	"strconv"
	"strings"
	"time"

	"github.com/OpenIMSDK/tools/component"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/s3/util"
	"gopkg.in/yaml.v2"

	"github.com/openimsdk/tools/errs"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
)

const (
	// defaultCfgPath is the default path of the configuration file.
	defaultCfgPath = "../../../../../config/config.yaml"
	maxRetry       = 100
)

var (
	cfgPath = flag.String("c", defaultCfgPath, "Path to the configuration file")
)

func initCfg() (*config.GlobalConfig, error) {
	data, err := os.ReadFile(*cfgPath)
	if err != nil {
		return nil, errs.WrapMsg(err, "ReadFile unmarshal failed")
	}

	conf := config.NewGlobalConfig()
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		return nil, errs.WrapMsg(err, "InitConfig unmarshal failed")
	}
	return conf, nil
}

type checkFunc struct {
	name     string
	function func(*config.GlobalConfig) error
	flag     bool
	config   *config.GlobalConfig
}

// colorErrPrint prints formatted string in red to stderr
func colorErrPrint(msg string) {
	// ANSI escape code for red text
	const redColor = "\033[31m"
	// ANSI escape code to reset color
	const resetColor = "\033[0m"
	msg = redColor + msg + resetColor
	// Print to stderr in red
	fmt.Fprintf(os.Stderr, "%s\n", msg)
}

func colorSuccessPrint(format string, a ...interface{}) {
	// ANSI escape code for green text is \033[32m
	// \033[0m resets the color
	fmt.Printf("\033[32m"+format+"\033[0m", a...)
}

func main() {
	flag.Parse()

	conf, err := initCfg()
	if err != nil {
		fmt.Printf("Read config failed: %v\n", err)
		return
	}

	err = configGetEnv(conf)
	if err != nil {
		fmt.Printf("configGetEnv failed, err:%v", err)
		return
	}

	checks := []checkFunc{
		{name: "Mongo", function: checkMongo, config: conf},
		{name: "Redis", function: checkRedis, config: conf},
		{name: "Zookeeper", function: checkZookeeper, config: conf},
		{name: "Kafka", function: checkKafka, config: conf},
	}
	if conf.Object.Enable == "minio" {
		checks = append(checks, checkFunc{name: "Minio", function: checkMinio, config: conf})
	}

	for i := 0; i < maxRetry; i++ {
		if i != 0 {
			time.Sleep(1 * time.Second)
		}
		fmt.Printf("Checking components round %v...\n", i+1)

		var err error
		allSuccess := true
		for index, check := range checks {
			if !check.flag {
				err = check.function(check.config)
				if err != nil {
					allSuccess = false
					colorErrPrint(fmt.Sprintf("Check component: %s, failed: %v", check.name, err.Error()))

					if check.name == "Minio" {
						if errors.Is(err, errMinioNotEnabled) ||
							errors.Is(err, errSignEndPoint) ||
							errors.Is(err, errApiURL) {
							checks[index].flag = true
							continue
						}
						break
					}
				} else {
					checks[index].flag = true
					util.SuccessPrint(fmt.Sprintf("%s connected successfully", check.name))
				}
			}
		}
		if allSuccess {
			component.SuccessPrint("All components started successfully!")
			return
		}
	}
	component.ErrorPrint("Some components checked failed!")
	os.Exit(-1)
}

var errMinioNotEnabled = errors.New("minio.Enable is not configured to use MinIO")

var errSignEndPoint = errors.New("minio.signEndPoint contains 127.0.0.1, causing issues with image sending")
var errApiURL = errors.New("object.apiURL contains 127.0.0.1, causing issues with image sending")

func configGetEnv(config *config.GlobalConfig) error {
	config.Mongo.Uri = getEnv("MONGO_URI", config.Mongo.Uri)
	config.Mongo.Username = getEnv("MONGO_OPENIM_USERNAME", config.Mongo.Username)
	config.Mongo.Password = getEnv("MONGO_OPENIM_PASSWORD", config.Mongo.Password)
	config.Mongo.Address = getArrEnv("MONGO_ADDRESS", "MONGO_PORT", config.Mongo.Address)
	config.Mongo.Database = getEnv("MONGO_DATABASE", config.Mongo.Database)
	maxPoolSize, err := getEnvInt("MONGO_MAX_POOL_SIZE", config.Mongo.MaxPoolSize)
	if err != nil {
		return errs.WrapMsg(err, "MONGO_MAX_POOL_SIZE")
	}
	config.Mongo.MaxPoolSize = maxPoolSize

	config.Redis.Username = getEnv("REDIS_USERNAME", config.Redis.Username)
	config.Redis.Password = getEnv("REDIS_PASSWORD", config.Redis.Password)
	config.Redis.Address = getArrEnv("REDIS_ADDRESS", "REDIS_PORT", config.Redis.Address)

	config.Object.ApiURL = getEnv("OBJECT_APIURL", config.Object.ApiURL)
	config.Object.Minio.Endpoint = getEnv("MINIO_ENDPOINT", config.Object.Minio.Endpoint)
	config.Object.Minio.AccessKeyID = getEnv("MINIO_ACCESS_KEY_ID", config.Object.Minio.AccessKeyID)
	config.Object.Minio.SecretAccessKey = getEnv("MINIO_SECRET_ACCESS_KEY", config.Object.Minio.SecretAccessKey)
	config.Object.Minio.SignEndpoint = getEnv("MINIO_SIGN_ENDPOINT", config.Object.Minio.SignEndpoint)

	config.Zookeeper.Schema = getEnv("ZOOKEEPER_SCHEMA", config.Zookeeper.Schema)
	config.Zookeeper.ZkAddr = getArrEnv("ZOOKEEPER_ADDRESS", "ZOOKEEPER_PORT", config.Zookeeper.ZkAddr)
	config.Zookeeper.Username = getEnv("ZOOKEEPER_USERNAME", config.Zookeeper.Username)
	config.Zookeeper.Password = getEnv("ZOOKEEPER_PASSWORD", config.Zookeeper.Password)

	config.Kafka.Username = getEnv("KAFKA_USERNAME", config.Kafka.Username)
	config.Kafka.Password = getEnv("KAFKA_PASSWORD", config.Kafka.Password)
	config.Kafka.Addr = getArrEnv("KAFKA_ADDRESS", "KAFKA_PORT", config.Kafka.Addr)
	config.Object.Minio.Endpoint = getMinioAddr("MINIO_ENDPOINT", "MINIO_ADDRESS", "MINIO_PORT", config.Object.Minio.Endpoint)
	return nil
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

// Helper function to get environment variable or default value
func getEnvInt(key string, fallback int) (int, error) {
	if value, exists := os.LookupEnv(key); exists {
		val, err := strconv.Atoi(value)
		if err != nil {
			return 0, errs.WrapMsg(err, "string to int failed")
		}
		return val, nil
	}
	return fallback, nil
}

func getArrEnv(key1, key2 string, fallback []string) []string {
	address, addrExists := os.LookupEnv(key1)
	port, portExists := os.LookupEnv(key2)

	if addrExists && portExists {
		addresses := strings.Split(address, ",")
		for i, addr := range addresses {
			addresses[i] = addr + ":" + port
		}
		return addresses
	}

	if addrExists && !portExists {
		addresses := strings.Split(address, ",")
		for i, addr := range addresses {
			addresses[i] = addr + ":" + "0"
		}
		return addresses
	}

	if !addrExists && portExists {
		result := make([]string, len(fallback))
		for i, addr := range fallback {
			add := strings.Split(addr, ":")
			result[i] = add[0] + ":" + port
		}
		return result
	}
	return fallback
}
