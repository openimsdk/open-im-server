package prommetrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	MsgOfflinePushFailedCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "msg_offline_push_failed_total",
		Help: "The number of msg failed offline pushed",
	})
)
