package prom_metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

/*
labels := prometheus.Labels{"label_one": "any", "label_two": "value"}
G_grpc_msggateway_metrics.MetricsMap["demo_server_say_hello_method_handle_count"].(*prometheus.CounterVec).With(labels).Inc()
*/
var (
	G_grpc_msggateway_metrics *GrpcCusMetricsMap
)

func init() {
	customizedCounterMetric := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "demo_server_say_hello_method_handle_count",
		Help: "Total number of RPCs handled on the server.",
	}, []string{"name"})
	tMetrics := make(map[string]prometheus.Collector)
	tMetrics["demo_server_say_hello_method_handle_count"] = customizedCounterMetric
	G_grpc_msggateway_metrics = &GrpcCusMetricsMap{MetricsMap: tMetrics}
}
