package prommetrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	OnlineUserGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "online_user_num",
		Help: "The number of online user num",
	})
)
