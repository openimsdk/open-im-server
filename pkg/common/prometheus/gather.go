package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	//auth rpc
	UserLoginCounter    prometheus.Counter
	UserRegisterCounter prometheus.Counter

	//seg
	SeqGetSuccessCounter prometheus.Counter
	SeqGetFailedCounter  prometheus.Counter
	SeqSetSuccessCounter prometheus.Counter
	SeqSetFailedCounter  prometheus.Counter

	//msg-db
	MsgInsertRedisSuccessCounter   prometheus.Counter
	MsgInsertRedisFailedCounter    prometheus.Counter
	MsgInsertMongoSuccessCounter   prometheus.Counter
	MsgInsertMongoFailedCounter    prometheus.Counter
	MsgPullFromRedisSuccessCounter prometheus.Counter
	MsgPullFromRedisFailedCounter  prometheus.Counter
	MsgPullFromMongoSuccessCounter prometheus.Counter
	MsgPullFromMongoFailedCounter  prometheus.Counter

	//msg-ws
	MsgRecvTotalCounter          prometheus.Counter
	GetNewestSeqTotalCounter     prometheus.Counter
	PullMsgBySeqListTotalCounter prometheus.Counter

	SingleChatMsgRecvSuccessCounter         prometheus.Counter
	GroupChatMsgRecvSuccessCounter          prometheus.Counter
	WorkSuperGroupChatMsgRecvSuccessCounter prometheus.Counter

	//msg-msg
	SingleChatMsgProcessSuccessCounter         prometheus.Counter
	SingleChatMsgProcessFailedCounter          prometheus.Counter
	GroupChatMsgProcessSuccessCounter          prometheus.Counter
	GroupChatMsgProcessFailedCounter           prometheus.Counter
	WorkSuperGroupChatMsgProcessSuccessCounter prometheus.Counter
	WorkSuperGroupChatMsgProcessFailedCounter  prometheus.Counter

	//msg-push
	MsgOnlinePushSuccessCounter  prometheus.Counter
	MsgOfflinePushSuccessCounter prometheus.Counter
	MsgOfflinePushFailedCounter  prometheus.Counter
	// api
	ApiRequestCounter        prometheus.Counter
	ApiRequestSuccessCounter prometheus.Counter
	ApiRequestFailedCounter  prometheus.Counter

	// grpc
	GrpcRequestCounter        prometheus.Counter
	GrpcRequestSuccessCounter prometheus.Counter
	GrpcRequestFailedCounter  prometheus.Counter

	SendMsgCounter               prometheus.Counter

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

func NewApiRequestCounter() {
	ApiRequestCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "api_request",
		Help: "The number of api request",
	})
}

func NewApiRequestSuccessCounter() {
	ApiRequestSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "api_request_success",
		Help: "The number of api request success",
	})
}

func NewApiRequestFailedCounter() {
	ApiRequestFailedCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "api_request_failed",
		Help: "The number of api request failed",
	})
}

func NewGrpcRequestCounter() {
	GrpcRequestCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "grpc_request",
		Help: "The number of api request",
	})
}

func NewGrpcRequestSuccessCounter() {
	GrpcRequestSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "grpc_request_success",
		Help: "The number of grpc request success",
	})
}

func NewGrpcRequestFailedCounter() {
	GrpcRequestFailedCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "grpc_request_failed",
		Help: "The number of grpc request failed",
	})
}

func NewSendMsgCount() {
	SendMsgCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "send_msg",
		Help: "The number of send msg",
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

func NewMsgPullFromRedisSuccessCounter() {
	MsgPullFromRedisSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "msg_pull_from_redis_success",
		Help: "The number of successful pull msg from redis",
	})
}

func NewMsgPullFromRedisFailedCounter() {
	MsgPullFromRedisFailedCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "msg_pull_from_redis_failed",
		Help: "The number of failed pull msg from redis",
	})
}

func NewMsgPullFromMongoSuccessCounter() {
	MsgPullFromMongoSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "msg_pull_from_redis_success",
		Help: "The number of successful pull msg from mongo",
	})
}

func NewMsgPullFromMongoFailedCounter() {
	MsgPullFromMongoFailedCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "msg_pull_from_mongo_failed",
		Help: "The number of failed pull msg from mongo",
	})
}

func NewMsgRecvTotalCounter() {
	MsgRecvTotalCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "msg_recv_total",
		Help: "The number of msg received",
	})
}

func NewGetNewestSeqTotalCounter() {
	GetNewestSeqTotalCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "get_newest_seq_total",
		Help: "the number of get newest seq",
	})
}
func NewPullMsgBySeqListTotalCounter() {
	PullMsgBySeqListTotalCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pull_msg_by_seq_list_total",
		Help: "The number of pull msg by seq list",
	})
}

func NewSingleChatMsgRecvSuccessCounter() {
	SingleChatMsgRecvSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "single_chat_msg_recv_success",
		Help: "The number of single chat msg successful received ",
	})
}

func NewGroupChatMsgRecvSuccessCounter() {
	GroupChatMsgRecvSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "group_chat_msg_recv_success",
		Help: "The number of group chat msg successful received",
	})
}

func NewWorkSuperGroupChatMsgRecvSuccessCounter() {
	WorkSuperGroupChatMsgRecvSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "work_super_group_chat_msg_recv_success",
		Help: "The number of work/super group chat msg successful received",
	})
}

func NewSingleChatMsgProcessSuccessCounter() {
	SingleChatMsgProcessSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "single_chat_msg_process_success",
		Help: "The number of single chat msg successful processed",
	})
}

func NewSingleChatMsgProcessFailedCounter() {
	SingleChatMsgProcessFailedCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "single_chat_msg_process_failed",
		Help: "The number of single chat msg failed processed",
	})
}

func NewGroupChatMsgProcessSuccessCounter() {
	GroupChatMsgProcessSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "group_chat_msg_process_success",
		Help: "The number of group chat msg successful processed",
	})
}

func NewGroupChatMsgProcessFailedCounter() {
	GroupChatMsgProcessFailedCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "group_chat_msg_process_failed",
		Help: "The number of group chat msg failed processed",
	})
}

func NewWorkSuperGroupChatMsgProcessSuccessCounter() {
	WorkSuperGroupChatMsgProcessSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "work_super_group_chat_msg_process_success",
		Help: "The number of work/super group chat msg successful processed",
	})
}
func NewWorkSuperGroupChatMsgProcessFailedCounter() {
	WorkSuperGroupChatMsgProcessFailedCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "work_super_group_chat_msg_process_failed",
		Help: "The number of work/super group chat msg failed processed",
	})
}

func NewMsgOnlinePushSuccessCounter() {
	MsgOnlinePushSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "msg_online_push_success",
		Help: "The number of msg successful online pushed",
	})
}

func NewMsgOfflinePushSuccessCounter() {
	MsgOfflinePushSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "msg_offline_push_success",
		Help: "The number of msg successful offline pushed",
	})
}
func NewMsgOfflinePushFailedCounter() {
	MsgOfflinePushFailedCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "msg_offline_push_failed",
		Help: "The number of msg failed offline pushed",
	})
}
