package prome

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
	OnlineUserGauge                         prometheus.Gauge

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

	SendMsgCounter prometheus.Counter

	// conversation
	ConversationCreateSuccessCounter prometheus.Counter
	ConversationCreateFailedCounter  prometheus.Counter
)

func NewUserLoginCounter() {
	if UserLoginCounter != nil {
		return
	}
	UserLoginCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "user_login",
		Help: "The number of user login",
	})
}
func NewUserRegisterCounter() {
	if UserRegisterCounter != nil {
		return
	}
	UserRegisterCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "user_register",
		Help: "The number of user register",
	})
}

func NewSeqGetSuccessCounter() {
	if SeqGetSuccessCounter != nil {
		return
	}
	SeqGetSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "seq_get_success",
		Help: "The number of successful get seq",
	})
}
func NewSeqGetFailedCounter() {
	if SeqGetFailedCounter != nil {
		return
	}
	SeqGetFailedCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "seq_get_failed",
		Help: "The number of failed get seq",
	})
}

func NewSeqSetSuccessCounter() {
	if SeqSetSuccessCounter != nil {
		return
	}
	SeqSetSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "seq_set_success",
		Help: "The number of successful set seq",
	})
}
func NewSeqSetFailedCounter() {
	if SeqSetFailedCounter != nil {
		return
	}
	SeqSetFailedCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "seq_set_failed",
		Help: "The number of failed set seq",
	})
}

func NewApiRequestCounter() {
	if ApiRequestCounter != nil {
		return
	}
	ApiRequestCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "api_request",
		Help: "The number of api request",
	})
}

func NewApiRequestSuccessCounter() {
	if ApiRequestSuccessCounter != nil {
		return
	}
	ApiRequestSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "api_request_success",
		Help: "The number of api request success",
	})
}

func NewApiRequestFailedCounter() {
	if ApiRequestFailedCounter != nil {
		return
	}
	ApiRequestFailedCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "api_request_failed",
		Help: "The number of api request failed",
	})
}

func NewGrpcRequestCounter() {
	if GrpcRequestCounter != nil {
		return
	}
	GrpcRequestCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "grpc_request",
		Help: "The number of api request",
	})
}

func NewGrpcRequestSuccessCounter() {
	if GrpcRequestSuccessCounter != nil {
		return
	}
	GrpcRequestSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "grpc_request_success",
		Help: "The number of grpc request success",
	})
}

func NewGrpcRequestFailedCounter() {
	if GrpcRequestFailedCounter != nil {
		return
	}
	GrpcRequestFailedCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "grpc_request_failed",
		Help: "The number of grpc request failed",
	})
}

func NewSendMsgCount() {
	if SendMsgCounter != nil {
		return
	}
	SendMsgCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "send_msg",
		Help: "The number of send msg",
	})
}

func NewMsgInsertRedisSuccessCounter() {
	if MsgInsertRedisSuccessCounter != nil {
		return
	}
	MsgInsertRedisSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "msg_insert_redis_success",
		Help: "The number of successful insert msg to redis",
	})
}

func NewMsgInsertRedisFailedCounter() {
	if MsgInsertRedisFailedCounter != nil {
		return
	}
	MsgInsertRedisFailedCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "msg_insert_redis_failed",
		Help: "The number of failed insert msg to redis",
	})
}

func NewMsgInsertMongoSuccessCounter() {
	if MsgInsertMongoSuccessCounter != nil {
		return
	}
	MsgInsertMongoSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "msg_insert_mongo_success",
		Help: "The number of successful insert msg to mongo",
	})
}

func NewMsgInsertMongoFailedCounter() {
	if MsgInsertMongoFailedCounter != nil {
		return
	}
	MsgInsertMongoFailedCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "msg_insert_mongo_failed",
		Help: "The number of failed insert msg to mongo",
	})
}

func NewMsgPullFromRedisSuccessCounter() {
	if MsgPullFromRedisSuccessCounter != nil {
		return
	}
	MsgPullFromRedisSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "msg_pull_from_redis_success",
		Help: "The number of successful pull msg from redis",
	})
}

func NewMsgPullFromRedisFailedCounter() {
	if MsgPullFromRedisFailedCounter != nil {
		return
	}
	MsgPullFromRedisFailedCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "msg_pull_from_redis_failed",
		Help: "The number of failed pull msg from redis",
	})
}

func NewMsgPullFromMongoSuccessCounter() {
	if MsgPullFromMongoSuccessCounter != nil {
		return
	}
	MsgPullFromMongoSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "msg_pull_from_mongo_success",
		Help: "The number of successful pull msg from mongo",
	})
}

func NewMsgPullFromMongoFailedCounter() {
	if MsgPullFromMongoFailedCounter != nil {
		return
	}
	MsgPullFromMongoFailedCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "msg_pull_from_mongo_failed",
		Help: "The number of failed pull msg from mongo",
	})
}

func NewMsgRecvTotalCounter() {
	if MsgRecvTotalCounter != nil {
		return
	}
	MsgRecvTotalCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "msg_recv_total",
		Help: "The number of msg received",
	})
}

func NewGetNewestSeqTotalCounter() {
	if GetNewestSeqTotalCounter != nil {
		return
	}
	GetNewestSeqTotalCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "get_newest_seq_total",
		Help: "the number of get newest seq",
	})
}
func NewPullMsgBySeqListTotalCounter() {
	if PullMsgBySeqListTotalCounter != nil {
		return
	}
	PullMsgBySeqListTotalCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pull_msg_by_seq_list_total",
		Help: "The number of pull msg by seq list",
	})
}

func NewSingleChatMsgRecvSuccessCounter() {
	if SingleChatMsgRecvSuccessCounter != nil {
		return
	}
	SingleChatMsgRecvSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "single_chat_msg_recv_success",
		Help: "The number of single chat msg successful received ",
	})
}

func NewGroupChatMsgRecvSuccessCounter() {
	if GroupChatMsgRecvSuccessCounter != nil {
		return
	}
	GroupChatMsgRecvSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "group_chat_msg_recv_success",
		Help: "The number of group chat msg successful received",
	})
}

func NewWorkSuperGroupChatMsgRecvSuccessCounter() {
	if WorkSuperGroupChatMsgRecvSuccessCounter != nil {
		return
	}
	WorkSuperGroupChatMsgRecvSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "work_super_group_chat_msg_recv_success",
		Help: "The number of work/super group chat msg successful received",
	})
}

func NewOnlineUserGauges() {
	if OnlineUserGauge != nil {
		return
	}
	OnlineUserGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "online_user_num",
		Help: "The number of online user num",
	})
}

func NewSingleChatMsgProcessSuccessCounter() {
	if SingleChatMsgProcessSuccessCounter != nil {
		return
	}
	SingleChatMsgProcessSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "single_chat_msg_process_success",
		Help: "The number of single chat msg successful processed",
	})
}

func NewSingleChatMsgProcessFailedCounter() {
	if SingleChatMsgProcessFailedCounter != nil {
		return
	}
	SingleChatMsgProcessFailedCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "single_chat_msg_process_failed",
		Help: "The number of single chat msg failed processed",
	})
}

func NewGroupChatMsgProcessSuccessCounter() {
	if GroupChatMsgProcessSuccessCounter != nil {
		return
	}
	GroupChatMsgProcessSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "group_chat_msg_process_success",
		Help: "The number of group chat msg successful processed",
	})
}

func NewGroupChatMsgProcessFailedCounter() {
	if GroupChatMsgProcessFailedCounter != nil {
		return
	}
	GroupChatMsgProcessFailedCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "group_chat_msg_process_failed",
		Help: "The number of group chat msg failed processed",
	})
}

func NewWorkSuperGroupChatMsgProcessSuccessCounter() {
	if WorkSuperGroupChatMsgProcessSuccessCounter != nil {
		return
	}
	WorkSuperGroupChatMsgProcessSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "work_super_group_chat_msg_process_success",
		Help: "The number of work/super group chat msg successful processed",
	})
}
func NewWorkSuperGroupChatMsgProcessFailedCounter() {
	if WorkSuperGroupChatMsgProcessFailedCounter != nil {
		return
	}
	WorkSuperGroupChatMsgProcessFailedCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "work_super_group_chat_msg_process_failed",
		Help: "The number of work/super group chat msg failed processed",
	})
}

func NewMsgOnlinePushSuccessCounter() {
	if MsgOnlinePushSuccessCounter != nil {
		return
	}
	MsgOnlinePushSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "msg_online_push_success",
		Help: "The number of msg successful online pushed",
	})
}

func NewMsgOfflinePushSuccessCounter() {
	if MsgOfflinePushSuccessCounter != nil {
		return
	}
	MsgOfflinePushSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "msg_offline_push_success",
		Help: "The number of msg successful offline pushed",
	})
}
func NewMsgOfflinePushFailedCounter() {
	if MsgOfflinePushFailedCounter != nil {
		return
	}
	MsgOfflinePushFailedCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "msg_offline_push_failed",
		Help: "The number of msg failed offline pushed",
	})
}

func NewConversationCreateSuccessCounter() {
	if ConversationCreateSuccessCounter != nil {
		return
	}
	ConversationCreateSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "conversation_push_success",
		Help: "The number of conversation successful pushed",
	})
}

func NewConversationCreateFailedCounter() {
	if ConversationCreateFailedCounter != nil {
		return
	}
	ConversationCreateFailedCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "conversation_push_failed",
		Help: "The number of conversation failed pushed",
	})
}
