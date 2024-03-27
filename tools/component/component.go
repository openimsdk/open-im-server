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
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/OpenIMSDK/tools/component"
	kfk "github.com/openimsdk/tools/mq/kafka"

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

// checkMongo checks the MongoDB connection without retries
func checkMongo(config *config.GlobalConfig) error {
	mongoStu := &component.Mongo{
		URL:         config.Mongo.Uri,
		Address:     config.Mongo.Address,
		Database:    config.Mongo.Database,
		Username:    config.Mongo.Username,
		Password:    config.Mongo.Password,
		MaxPoolSize: config.Mongo.MaxPoolSize,
	}
	err := component.CheckMongo(mongoStu)

	return err
}

// checkRedis checks the Redis connection
func checkRedis(config *config.GlobalConfig) error {
	redisStu := &component.Redis{
		Address:  config.Redis.Address,
		Username: config.Redis.Username,
		Password: config.Redis.Password,
	}
	err := component.CheckRedis(redisStu)
	return err
}

// checkMinio checks the MinIO connection
func checkMinio(config *config.GlobalConfig) error {
	if strings.Contains(config.Object.ApiURL, "127.0.0.1") {
		return errs.Wrap(errApiURL)
	}
	if config.Object.Enable != "minio" {
		return errs.Wrap(errMinioNotEnabled)
	}
	if strings.Contains(config.Object.Minio.Endpoint, "127.0.0.1") {
		return errs.Wrap(errSignEndPoint)
	}

	minio := &component.Minio{
		ApiURL:          config.Object.ApiURL,
		Endpoint:        config.Object.Minio.Endpoint,
		AccessKeyID:     config.Object.Minio.AccessKeyID,
		SecretAccessKey: config.Object.Minio.SecretAccessKey,
		SignEndpoint:    config.Object.Minio.SignEndpoint,
		UseSSL:          getEnv("MINIO_USE_SSL", "false"),
	}
	err := component.CheckMinio(minio)
	return err
}

// checkZookeeper checks the Zookeeper connection
func checkZookeeper(config *config.GlobalConfig) error {
	zkStu := &component.Zookeeper{
		Schema:   config.Zookeeper.Schema,
		ZkAddr:   config.Zookeeper.ZkAddr,
		Username: config.Zookeeper.Username,
		Password: config.Zookeeper.Password,
	}
	err := component.CheckZookeeper(zkStu)
	return err
}

// checkKafka checks the Kafka connection
func checkKafka(config *config.GlobalConfig) error {
	topics := []string{
		config.Kafka.MsgToMongo.Topic,
		config.Kafka.MsgToPush.Topic,
		config.Kafka.LatestMsgToRedis.Topic,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	return kfk.CheckKafka(ctx, &config.Kafka.Config, topics)
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
