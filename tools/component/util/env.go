package env

import (
	"os"
	"strconv"
	"strings"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/tools/errs"
)

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
