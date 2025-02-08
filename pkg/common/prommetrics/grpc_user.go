package prommetrics

import "github.com/prometheus/client_golang/prometheus"

var (
	UserRegisterCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "user_register_total",
		Help: "The number of user login",
	})
)

func RegistryUser() {
	registry.MustRegister(UserRegisterCounter)
}
