package prommetrics

import (
	"net"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	MsgInsertRedisSuccessCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "msg_insert_redis_success_total",
		Help: "The number of successful insert msg to redis",
	})
	MsgInsertRedisFailedCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "msg_insert_redis_failed_total",
		Help: "The number of failed insert msg to redis",
	})
	MsgInsertMongoSuccessCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "msg_insert_mongo_success_total",
		Help: "The number of successful insert msg to mongo",
	})
	MsgInsertMongoFailedCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "msg_insert_mongo_failed_total",
		Help: "The number of failed insert msg to mongo",
	})
	SeqSetFailedCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "seq_set_failed_total",
		Help: "The number of failed set seq",
	})
)

func RegistryTransfer() {
	registry.MustRegister(
		MsgInsertRedisSuccessCounter,
		MsgInsertRedisFailedCounter,
		MsgInsertMongoSuccessCounter,
		MsgInsertMongoFailedCounter,
		SeqSetFailedCounter,
	)
}

func TransferInit(listener net.Listener) error {
	reg := prometheus.NewRegistry()
	cs := append(
		baseCollector,
		MsgInsertRedisSuccessCounter,
		MsgInsertRedisFailedCounter,
		MsgInsertMongoSuccessCounter,
		MsgInsertMongoFailedCounter,
		SeqSetFailedCounter,
	)
	return Init(reg, listener, commonPath, promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}), cs...)
}
