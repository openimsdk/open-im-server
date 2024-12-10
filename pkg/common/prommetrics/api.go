package prommetrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net"
	"strconv"
)

var (
	apiCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_count",
			Help: "Total number of API calls",
		},
		[]string{"path", "method", "code"},
	)
	httpCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_count",
			Help: "Total number of HTTP calls",
		},
		[]string{"path", "method", "status"},
	)
)

func ApiInit(listener net.Listener) error {
	apiRegistry := prometheus.NewRegistry()
	cs := append(
		baseCollector,
		apiCounter,
		httpCounter,
	)
	return Init(apiRegistry, listener, commonPath, promhttp.HandlerFor(apiRegistry, promhttp.HandlerOpts{}), cs...)
}

func APICall(path string, method string, apiCode int) {
	apiCounter.With(prometheus.Labels{"path": path, "method": method, "code": strconv.Itoa(apiCode)}).Inc()
}

func HttpCall(path string, method string, status int) {
	httpCounter.With(prometheus.Labels{"path": path, "method": method, "status": strconv.Itoa(status)}).Inc()
}

//func ApiHandler() http.Handler {
//	return promhttp.InstrumentMetricHandler(
//		apiRegistry, promhttp.HandlerFor(apiRegistry, promhttp.HandlerOpts{}),
//	)
//}
