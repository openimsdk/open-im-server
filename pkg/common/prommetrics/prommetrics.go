// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
)

type Target struct {
	Target string            `json:"target"`
	Labels map[string]string `json:"labels"`
}

type RespTarget struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

func BuildDiscoveryKey(name string) string {
	return fmt.Sprintf("%s/%s/%s", "openim", "prometheus_discovery", name)
}

func BuildDefaultTarget(host string, ip int) Target {
	return Target{
		Target: fmt.Sprintf("%s:%d", host, ip),
		Labels: map[string]string{
			"namespace": "default",
		},
	}
}
