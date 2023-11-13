package prommetrics

import (
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"

	config2 "github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/ginprometheus"
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
	case config2.Config.RpcRegisterName.OpenImMessageGatewayName:
		return []prometheus.Collector{OnlineUserGauge}
	case config2.Config.RpcRegisterName.OpenImMsgName:
		return []prometheus.Collector{SingleChatMsgProcessSuccessCounter, SingleChatMsgProcessFailedCounter, GroupChatMsgProcessSuccessCounter, GroupChatMsgProcessFailedCounter}
	case "Transfer":
		return []prometheus.Collector{MsgInsertRedisSuccessCounter, MsgInsertRedisFailedCounter, MsgInsertMongoSuccessCounter, MsgInsertMongoFailedCounter, SeqSetFailedCounter}
	case config2.Config.RpcRegisterName.OpenImPushName:
		return []prometheus.Collector{MsgOfflinePushFailedCounter}
	case config2.Config.RpcRegisterName.OpenImAuthName:
		return []prometheus.Collector{UserLoginCounter}
	default:
		return nil
	}
}

func GetGinCusMetrics(name string) []*ginprometheus.Metric {
	switch name {
	case "Api":
		return []*ginprometheus.Metric{ApiCustomCnt}
	default:
		return []*ginprometheus.Metric{ApiCustomCnt}
	}
}
