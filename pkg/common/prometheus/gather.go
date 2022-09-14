package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// user rpc
	UserLoginCounter    prometheus.Counter
	UserRegisterCounter prometheus.Counter

	SeqGetSuccessCounter prometheus.Counter
	SeqGetFailedCounter  prometheus.Counter
	SeqSetSuccessCounter prometheus.Counter
	SeqSetFailedCounter  prometheus.Counter
)

func NewUserLoginCounter() {
	UserLoginCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "user_login",
		Help: "The number of user login",
	})
}

func NewSeqGetSuccessCounter() {
	SeqGetSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "seq_get_success",
		Help: "The number of successful get seq",
	})
}
func NewSeqGetFailedCounter() {
	SeqGetFailedCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "seq_get_failed",
		Help: "The number of failed get seq",
	})
}

func NewSeqSetSuccessCounter() {
	SeqSetSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "seq_set_success",
		Help: "The number of successful set seq",
	})
}
func NewSeqSetFailedCounter() {
	SeqSetFailedCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "seq_set_failed",
		Help: "The number of failed set seq",
	})
}
