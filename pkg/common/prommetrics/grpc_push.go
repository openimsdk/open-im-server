package prommetrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	MsgOfflinePushFailedCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "msg_offline_push_failed_total",
		Help: "The number of msg failed offline pushed",
	})
	MsgLoneTimePushCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "msg_long_time_push_total",
		Help: "The number of messages with a push time exceeding 10 seconds",
	})
)

func RegistryPush() {
	registry.MustRegister(
		MsgOfflinePushFailedCounter,
		MsgLoneTimePushCounter,
	)
}
