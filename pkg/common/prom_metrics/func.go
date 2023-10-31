package prom_metrics

import (
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/openimsdk/open-im-server/v3/pkg/common/ginPrometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

func NewGrpcPromObj(cusMetrics []prometheus.Collector) (*prometheus.Registry, *grpc_prometheus.ServerMetrics, error) {
	////////////////////////////////////////////////////////
	reg := prometheus.NewRegistry()
	grpcMetrics := grpc_prometheus.NewServerMetrics()
	grpcMetrics.EnableHandlingTimeHistogram()
	cusMetrics = append(cusMetrics, grpcMetrics, collectors.NewGoCollector())
	reg.MustRegister(cusMetrics...)
	return reg, grpcMetrics, nil
}

func GetGrpcCusMetrics(registerName string) []prometheus.Collector {
	switch registerName {
	case "MessageGateway":
		return []prometheus.Collector{OnlineUserGauge}
	case "Msg":
		return []prometheus.Collector{SingleChatMsgProcessSuccessCounter, SingleChatMsgProcessFailedCounter, GroupChatMsgProcessSuccessCounter, GroupChatMsgProcessFailedCounter}
	case "Transfer":
		return []prometheus.Collector{MsgInsertRedisSuccessCounter, MsgInsertRedisFailedCounter, MsgInsertMongoSuccessCounter, MsgInsertMongoFailedCounter, SeqSetFailedCounter}
	case "Push":
		return []prometheus.Collector{MsgOfflinePushFailedCounter}
	case "Auth":
		return []prometheus.Collector{UserLoginCounter}
	default:
		return nil
	}
}

func GetGinCusMetrics(name string) []*ginPrometheus.Metric {
	switch name {
	case "Api":
		return []*ginPrometheus.Metric{ApiCustomCnt}
	default:
		return []*ginPrometheus.Metric{ApiCustomCnt}
	}
}
