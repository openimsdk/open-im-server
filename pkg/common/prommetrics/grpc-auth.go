package prommetrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	UserLoginCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "user_login_total",
		Help: "The number of user login",
	})
)
