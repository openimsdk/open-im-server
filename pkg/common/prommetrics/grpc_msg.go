package prommetrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	SingleChatMsgProcessSuccessCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "single_chat_msg_process_success_total",
		Help: "The number of single chat msg successful processed",
	})
	SingleChatMsgProcessFailedCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "single_chat_msg_process_failed_total",
		Help: "The number of single chat msg failed processed",
	})
	GroupChatMsgProcessSuccessCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "group_chat_msg_process_success_total",
		Help: "The number of group chat msg successful processed",
	})
	GroupChatMsgProcessFailedCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "group_chat_msg_process_failed_total",
		Help: "The number of group chat msg failed processed",
	})
)
