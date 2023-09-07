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
	"database/sql"
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"

	"github.com/IBM/sarama"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/utils"
	"github.com/go-zookeeper/zk"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/kafka"

	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	// defaultCfgPath is the default path of the configuration file
	defaultCfgPath           = "../../../../../config/config.yaml"
	minioHealthCheckDuration = 1
	maxRetry                 = 100
	componentStartErrCode    = 6000
	configErrCode            = 6001
)

var (
	cfgPath = flag.String("c", defaultCfgPath, "Path to the configuration file")

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
	function func() error
}

func main() {
	flag.Parse()

	if err := initCfg(); err != nil {
		fmt.Printf("Read config failed: %v\n", err)
		return
	}

	checks := []checkFunc{
		{name: "Mysql", function: checkMysql},
		{name: "Mongo", function: checkMongo},
		{name: "Minio", function: checkMinio},
		{name: "Redis", function: checkRedis},
		{name: "Zookeeper", function: checkZookeeper},
		{name: "Kafka", function: checkKafka},
	}

	for i := 0; i < maxRetry; i++ {
		if i != 0 {
			time.Sleep(3 * time.Second)
		}
		fmt.Printf("Checking components Round %v...\n", i+1)

		allSuccess := true
		for _, check := range checks {
			err := check.function()
			if err != nil {
				errorPrint(fmt.Sprintf("Starting %s failed: %v", check.name, err))
				allSuccess = false
				break
			} else {
				successPrint(fmt.Sprintf("%s starts successfully", check.name))
			}
		}

		if allSuccess {
			successPrint("All components started successfully!")
			return
		}
	}
	os.Exit(1)
}

func exactIP(urll string) string {
	u, _ := url.Parse(urll)
	host, _, err := net.SplitHostPort(u.Host)
	if err != nil {
		host = u.Host
	}
	if strings.HasSuffix(host, ":") {
		host = host[0 : len(host)-1]
	}
	return host
}

func checkMysql() error {
	var sqlDB *sql.DB
	defer func() {
		if sqlDB != nil {
			sqlDB.Close()
		}
	}()
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		config.Config.Mysql.Username, config.Config.Mysql.Password, config.Config.Mysql.Address[0], "mysql")
	db, err := gorm.Open(mysql.Open(dsn), nil)
	if err != nil {
		return errs.Wrap(err)
	} else {
		sqlDB, err = db.DB()
		err = sqlDB.Ping()
		if err != nil {
			return errs.Wrap(err)
		}
	}
	return nil
}

func checkMongo() error {
	var client *mongo.Client
	defer func() {
		if client != nil {
			client.Disconnect(context.TODO())
		}
	}()
	mongodbHosts := ""
	for i, v := range config.Config.Mongo.Address {
		if i == len(config.Config.Mongo.Address)-1 {
			mongodbHosts += v
		} else {
			mongodbHosts += v + ","
		}
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(
		fmt.Sprintf("mongodb://%v:%v@%v/?authSource=admin",
			config.Config.Mongo.Username, config.Config.Mongo.Password, mongodbHosts)))
	if err != nil {
		return errs.Wrap(err)
	} else {
		err = client.Ping(context.TODO(), &readpref.ReadPref{})
		if err != nil {
			return errs.Wrap(err)
		}
	}
	return nil
}

func checkMinio() error {
	if config.Config.Object.Enable == "minio" {
		conf := config.Config.Object.Minio
		u, _ := url.Parse(conf.Endpoint)
		minioClient, err := minio.New(u.Host, &minio.Options{
			Creds:  credentials.NewStaticV4(conf.AccessKeyID, conf.SecretAccessKey, ""),
			Secure: u.Scheme == "https",
		})
		if err != nil {
			return errs.Wrap(err)
		}

		cancel, err := minioClient.HealthCheck(time.Duration(minioHealthCheckDuration) * time.Second)
		defer func() {
			if cancel != nil {
				cancel()
			}
		}()
		if err != nil {
			return errs.Wrap(err)
		} else {
			if minioClient.IsOffline() {
				return ErrComponentStart.Wrap("Minio server is offline")
			}
		}
		if exactIP(config.Config.Object.ApiURL) == "127.0.0.1" || exactIP(config.Config.Object.Minio.SignEndpoint) == "127.0.0.1" {
			return ErrConfig.Wrap("apiURL or Minio SignEndpoint endpoint contain 127.0.0.1.")
		}
	}
	return nil
}

func checkRedis() error {
	var redisClient redis.UniversalClient
	defer func() {
		if redisClient != nil {
			redisClient.Close()
		}
	}()
	if len(config.Config.Redis.Address) > 1 {
		redisClient = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    config.Config.Redis.Address,
			Username: config.Config.Redis.Username,
			Password: config.Config.Redis.Password,
		})
	} else {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     config.Config.Redis.Address[0],
			Username: config.Config.Redis.Username,
			Password: config.Config.Redis.Password,
		})
	}
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		return errs.Wrap(err)
	}
	return nil
}

func checkZookeeper() error {
	var c *zk.Conn
	defer func() {
		if c != nil {
			c.Close()
		}
	}()
	c, _, err := zk.Connect(config.Config.Zookeeper.ZkAddr, time.Second)
	if err != nil {
		return errs.Wrap(err)
	} else {
		if config.Config.Zookeeper.Username != "" && config.Config.Zookeeper.Password != "" {
			if err := c.AddAuth("digest", []byte(config.Config.Zookeeper.Username+":"+config.Config.Zookeeper.Password)); err != nil {
				return errs.Wrap(err)
			}
		}
		_, _, err = c.Get("/")
		if err != nil {
			return errs.Wrap(err)
		}
	}
	return nil
}

func checkKafka() error {
	var kafkaClient sarama.Client
	defer func() {
		if kafkaClient != nil {
			kafkaClient.Close()
		}
	}()
	cfg := sarama.NewConfig()
	if config.Config.Kafka.Username != "" && config.Config.Kafka.Password != "" {
		cfg.Net.SASL.Enable = true
		cfg.Net.SASL.User = config.Config.Kafka.Username
		cfg.Net.SASL.Password = config.Config.Kafka.Password
	}
	kafka.SetupTLSConfig(cfg)
	kafkaClient, err := sarama.NewClient(config.Config.Kafka.Addr, cfg)
	if err != nil {
		return errs.Wrap(err)
	} else {
		topics, err := kafkaClient.Topics()
		if err != nil {
			return err
		}
		if !utils.IsContain(config.Config.Kafka.MsgToMongo.Topic, topics) {
			return ErrComponentStart.Wrap(fmt.Sprintf("kafka doesn't contain topic:%v", config.Config.Kafka.MsgToMongo.Topic))
		}
		if !utils.IsContain(config.Config.Kafka.MsgToPush.Topic, topics) {
			return ErrComponentStart.Wrap(fmt.Sprintf("kafka doesn't contain topic:%v", config.Config.Kafka.MsgToPush.Topic))
		}
		if !utils.IsContain(config.Config.Kafka.LatestMsgToRedis.Topic, topics) {
			return ErrComponentStart.Wrap(fmt.Sprintf("kafka doesn't contain topic:%v", config.Config.Kafka.LatestMsgToRedis.Topic))
		}
	}
	return nil
}

func errorPrint(s string) {
	fmt.Printf("\x1b[%dm%v\x1b[0m\n", 31, s)
}

func successPrint(s string) {
	fmt.Printf("\x1b[%dm%v\x1b[0m\n", 32, s)
}

func warningPrint(s string) {
	fmt.Printf("\x1b[%dmWarning: But %v\x1b[0m\n", 33, s)
}
