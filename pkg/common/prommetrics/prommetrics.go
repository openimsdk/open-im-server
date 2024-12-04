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
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

const commonPath = "/metrics"

var (
	baseCollector = []prometheus.Collector{
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(),
	}
)

func Init(registry *prometheus.Registry, prometheusPort int, path string, handler http.Handler, cs ...prometheus.Collector) error {
	registry.MustRegister(cs...)
	srv := http.NewServeMux()
	srv.Handle(path, handler)
	return http.ListenAndServe(fmt.Sprintf(":%d", prometheusPort), srv)
}
