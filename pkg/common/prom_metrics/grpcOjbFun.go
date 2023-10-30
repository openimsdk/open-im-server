package prom_metrics

import (
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

type GrpcCusMetricsMap struct {
	MetricsMap map[string]prometheus.Collector
}

func (m *GrpcCusMetricsMap) MetricList() []prometheus.Collector {
	var ret []prometheus.Collector
	for _, v := range m.MetricsMap {
		ret = append(ret, v)
	}
	return ret
}

func NewGrpcPromObj(cusMetrics []prometheus.Collector) (*prometheus.Registry, *grpc_prometheus.ServerMetrics, error) {
	////////////////////////////////////////////////////////
	reg := prometheus.NewRegistry()
	grpcMetrics := grpc_prometheus.NewServerMetrics()
	grpcMetrics.EnableHandlingTimeHistogram()
	cusMetrics = append(cusMetrics, grpcMetrics, prometheus.NewGoCollector())
	reg.MustRegister(cusMetrics...)
	return reg, grpcMetrics, nil
}
func GetGrpcCusMetrics(name string) (*GrpcCusMetricsMap, error) {
	switch name {
	case "MessageGateway":
		return G_grpc_msggateway_metrics, nil
	case "User":
		return G_grpc_msggateway_metrics, nil
	case "Msg":
		return G_grpc_msggateway_metrics, nil
	case "Conversation":
		return G_grpc_msggateway_metrics, nil
	default:
		return G_grpc_msggateway_metrics, nil
	}
}
