package prommetrics

import ginProm "github.com/openimsdk/open-im-server/v3/pkg/common/ginprometheus"

/*
labels := prometheus.Labels{"label_one": "any", "label_two": "value"}
ApiCustomCnt.MetricCollector.(*prometheus.CounterVec).With(labels).Inc()
*/
var (
	ApiCustomCnt = &ginProm.Metric{
		Name:        "custom_total",
		Description: "Custom counter events.",
		Type:        "counter_vec",
		Args:        []string{"label_one", "label_two"},
	}
)
