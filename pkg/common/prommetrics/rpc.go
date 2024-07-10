package prommetrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strconv"
)

const RpcPath = "/metrics"

var (
	rpcRegistry = prometheus.NewRegistry()
	rpcCounter  = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rpc_count",
			Help: "Total number of RPC calls",
		},
		[]string{"name", "path", "code"},
	)
)

func init() {
	rpcRegistry.MustRegister(rpcCounter)
}

func RPCCall(name string, path string, code int) {
	rpcCounter.With(prometheus.Labels{"name": name, "path": path, "code": strconv.Itoa(code)}).Inc()
}

func RPCHandler() http.Handler {
	return promhttp.HandlerFor(rpcRegistry, promhttp.HandlerOpts{})
}
