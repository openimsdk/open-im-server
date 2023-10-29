package prom_metrics

import ginProm "github.com/openimsdk/open-im-server/v3/pkg/common/ginPrometheus"

/*
labels := prometheus.Labels{"label_one": "any", "label_two": "value"}
G_api_metrics.MetricsMap["custom_total"].MetricCollector.(*prometheus.CounterVec).With(labels).Inc()
*/
var (
	G_api_metrics *ginProm.CusMetrics
)

func init() {

	CustomCnt := &ginProm.Metric{
		Name:        "custom_total",
		Description: "Custom counter events.",
		Type:        "counter_vec",
		Args:        []string{"label_one", "label_two"},
	}
	tMetrics := make(map[string]*ginProm.Metric)
	tMetrics["custom_total"] = CustomCnt
	G_api_metrics = &ginProm.CusMetrics{MetricsMap: tMetrics}
}
