package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/tools/utils"
	"github.com/Shopify/sarama"
	"github.com/go-zookeeper/zk"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	sqlDriver                = "mysql"
	minioHealthCheckDuration = 1
	maxRetry                 = 3
)

func main() {

	for i := 0; i < maxRetry; i++ {
		success := 1
		// Check MySQL
		db, err := sql.Open(sqlDriver, fmt.Sprintf("%s:%s@tcp(%s)/",
			config.Config.Mysql.Username, config.Config.Mysql.Password, config.Config.Mysql.Address))
		if err != nil {
			fmt.Printf("Cannot connect to MySQL: %v", err)
			success = 0
		}
		err = db.Ping()
		if err != nil {
			fmt.Printf("ping mysql failed: %v. please make sure your mysql service has started", err)
			success = 0
		}
		db.Close()

		// Check MongoDB
		client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(
			fmt.Sprintf("mongodb://%s:%s@%s",
				config.Config.Mongo.Username, config.Config.Mongo.Password, config.Config.Mongo.Address)))
		if err != nil {
			fmt.Printf("Cannot connect to MongoDB: %v", err)
			success = 0
		}
		err = client.Ping(context.TODO(), &readpref.ReadPref{})
		if err != nil {
			fmt.Printf("ping mysql failed: %v. please make sure your mysql service has started", err)
			success = 0
		}
		client.Disconnect(context.TODO())

		// Check Minio
		if config.Config.Object.Enable == "minio" {
			if exactIP(config.Config.Object.ApiURL) == "127.0.0.1" && config.Config.Object.ApiURL == config.Config.Object.Minio.Endpoint {
				fmt.Printf("ApiURL contain the same address with Endpoint: %v. please modify the config file", config.Config.Object.ApiURL)
			}
			minioClient, err := minio.New(config.Config.Object.Minio.Endpoint, &minio.Options{
				Creds:  credentials.NewStaticV4(config.Config.Object.Minio.AccessKeyID, config.Config.Object.Minio.SecretAccessKey, ""),
				Secure: false,
			})
			if err != nil {
				fmt.Printf("Cannot connect to Minio: %v", err)
				success = 0
			}
			cancel, err := minioClient.HealthCheck(time.Duration(minioHealthCheckDuration))
			if err != nil {
				fmt.Printf("starting minio health check failed:%v", err)
				success = 0
			}
			if minioClient.IsOffline() {
				fmt.Printf("Error: minio server is offline.")
				success = 0
			}
			cancel()
		}

		// Check Redis
		var redisClient redis.UniversalClient
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
		_, err = redisClient.Ping(context.Background()).Result()
		if err != nil {
			fmt.Printf("Cannot connect to Redis: %v", err)
			success = 0
		}

		// Check Zookeeper
		c, _, err := zk.Connect(config.Config.Zookeeper.ZkAddr, time.Second)
		if err != nil {
			fmt.Printf("Cannot connect to Zookeeper: %v", err)
			success = 0
		}
		c.Close()

		// Check Kafka
		kafkaClient, err := sarama.NewClient(config.Config.Kafka.Addr, &sarama.Config{})
		if err != nil {
			fmt.Printf("Cannot connect to Kafka: %v", err)
			success = 0
		} else {
			topics, err := kafkaClient.Topics()
			if err != nil {
				fmt.Println("get kafka topic error")
				success = 0
			}
			if !utils.IsContain(config.Config.Kafka.MsgToMongo.Topic, topics) {
				fmt.Printf("kafka doesn't contain topic:%v", config.Config.Kafka.MsgToMongo.Topic)
				success = 0
			}
			if !utils.IsContain(config.Config.Kafka.MsgToPush.Topic, topics) {
				fmt.Printf("kafka doesn't contain topic:%v", config.Config.Kafka.MsgToPush.Topic)
				success = 0
			}
			if !utils.IsContain(config.Config.Kafka.LatestMsgToRedis.Topic, topics) {
				fmt.Printf("kafka doesn't contain topic:%v", config.Config.Kafka.LatestMsgToRedis.Topic)
				success = 0
			}
		}
		kafkaClient.Close()
		if success == 1 {
			fmt.Println("all compose check pass")
			return
		}
		time.Sleep(3 * time.Second)
	}
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
