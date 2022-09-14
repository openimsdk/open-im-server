package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	UserLoginCounter    prometheus.Counter
	UserRegisterCounter prometheus.Counter

	SeqGetSuccessCounter prometheus.Counter
	SeqGetFailedCounter  prometheus.Counter
	SeqSetSuccessCounter prometheus.Counter
	SeqSetFailedCounter  prometheus.Counter

	MsgInsertRedisSuccessCounter prometheus.Counter
	MsgInsertRedisFailedCounter  prometheus.Counter

	MsgInsertMongoSuccessCounter prometheus.Counter
	MsgInsertMongoFailedCounter  prometheus.Counter
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

func NewMsgInsertRedisSuccessCounter() {
	MsgInsertRedisSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "msg_insert_redis_success",
		Help: "The number of successful insert msg to redis",
	})
}

func NewMsgInsertRedisFailedCounter() {
	MsgInsertRedisFailedCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "msg_insert_redis_failed",
		Help: "The number of failed insert msg to redis",
	})
}

func NewMsgInsertMongoSuccessCounter() {
	MsgInsertMongoSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "msg_insert_mongo_success",
		Help: "The number of successful insert msg to mongo",
	})
}

func NewMsgInsertMongoFailedCounter() {
	MsgInsertMongoFailedCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "msg_insert_mongo_failed",
		Help: "The number of failed insert msg to mongo",
	})
}
