package checks

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mq/kafka"
)

type KafkaCheck struct {
	Kafka *config.Kafka
}


func CheckKafka(ctx context.Context, config *KafkaCheck) error {
	kafkaConfig := &kafka.Config{
		Addr:     config.Kafka.Addr,
		Username: config.Kafka.Username,
		Password: config.Kafka.Password,
	}

	requiredTopics := []string{
		config.Kafka.MsgToMongo.Topic,
		config.Kafka.MsgToPush.Topic,
		config.Kafka.LatestMsgToRedis.Topic,
	}

	log.CInfo(ctx, "Checking Kafka connection", "Address", kafkaConfig.Addr, "Topics", requiredTopics)

	err := kafka.CheckKafka(ctx, kafkaConfig, requiredTopics)
	if err != nil {
		log.CInfo(ctx, "Kafka connection failed", "error", err)
		return err
	}

	log.CInfo(ctx, "Kafka connection and required topics verified successfully")
	return nil
}
