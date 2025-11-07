package prommetrics

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const commonPath = "/metrics"

var registry = &prometheusRegistry{prometheus.NewRegistry()}

type prometheusRegistry struct {
	*prometheus.Registry
}

func (x *prometheusRegistry) MustRegister(cs ...prometheus.Collector) {
	for _, c := range cs {
		if err := x.Registry.Register(c); err != nil {
			if errors.As(err, &prometheus.AlreadyRegisteredError{}) {
				continue
			}
			panic(err)
		}
	}
}

func init() {
	registry.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(),
	)
}

var (
	baseCollector = []prometheus.Collector{
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(),
	}
)

func Init(registry *prometheus.Registry, listener net.Listener, path string, handler http.Handler, cs ...prometheus.Collector) error {
	registry.MustRegister(cs...)
	srv := http.NewServeMux()
	srv.Handle(path, handler)
	return http.Serve(listener, srv)
}

func RegistryAll() {
	RegistryApi()
	RegistryAuth()
	RegistryMsg()
	RegistryMsgGateway()
	RegistryPush()
	RegistryUser()
	RegistryRpc()
	RegistryTransfer()
}

func Start(listener net.Listener) error {
	srv := http.NewServeMux()
	srv.Handle(commonPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	return http.Serve(listener, srv)
}

const (
	APIKeyName             = "api"
	MessageTransferKeyName = "message-transfer"

	TTL = 300
)

type Target struct {
	Target string            `json:"target"`
	Labels map[string]string `json:"labels"`
}

type RespTarget struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

func BuildDiscoveryKeyPrefix(name string) string {
	return fmt.Sprintf("%s/%s/%s", "openim", "prometheus_discovery", name)
}

func BuildDiscoveryKey(name string, index int) string {
	return fmt.Sprintf("%s/%s/%s/%d", "openim", "prometheus_discovery", name, index)
}

func BuildDefaultTarget(host string, ip int) Target {
	return Target{
		Target: fmt.Sprintf("%s:%d", host, ip),
		Labels: map[string]string{
			"namespace": "default",
		},
	}
}
