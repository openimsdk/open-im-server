package config

import (
	"strings"

	"github.com/openimsdk/tools/errs"
)

const (
	QueueEngineKafka  = "kafka"
	QueueEngineRedis  = "redis"
	QueueEngineMemory = "memory"
)

func NormalizeQueueEngine(engine string) string {
	switch strings.ToLower(strings.TrimSpace(engine)) {
	case "", "kafka":
		return QueueEngineKafka
	case "redis":
		return QueueEngineRedis
	case "memory":
		return QueueEngineMemory
	default:
		return strings.ToLower(strings.TrimSpace(engine))
	}
}

func ValidateQueueEngine(engine string, standalone bool) (string, error) {
	normalized := NormalizeQueueEngine(engine)
	switch normalized {
	case QueueEngineKafka, QueueEngineRedis:
		return normalized, nil
	case QueueEngineMemory:
		if standalone {
			return normalized, nil
		}
		return "", errs.ErrArgs.WrapMsg("unsupported queue engine for microservice deployment", "engine", engine)
	default:
		return "", errs.ErrArgs.WrapMsg("unsupported queue engine", "engine", engine)
	}
}
