package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// user rpc
	UserLoginCounter    prometheus.Counter
	UserRegisterCounter prometheus.Counter
)

func NewUserLoginCounter() {
	UserLoginCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "user_login",
		Help: "The number of user login",
	})
}

func NewUserRegisterCounter() {
	UserRegisterCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "user_register",
		Help: "The number of user register",
	})
}
