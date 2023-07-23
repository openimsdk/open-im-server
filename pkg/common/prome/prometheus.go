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

package prome

import (
	"bytes"
	"net/http"
	"strconv"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// StartPrometheusSrv
func StartPrometheusSrv(prometheusPort int) error {
	if config.Config.Prometheus.Enable {
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(":"+strconv.Itoa(prometheusPort), nil)
		return err
	}
	return nil
}

// PrometheusHandler
func PrometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// responseBodyWriter
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write
func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// PrometheusMiddleware
func PrometheusMiddleware(c *gin.Context) {
	Inc(ApiRequestCounter)
	w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
	c.Writer = w
	c.Next()
	if c.Writer.Status() == http.StatusOK {
		Inc(ApiRequestSuccessCounter)
	} else {
		Inc(ApiRequestFailedCounter)
	}
}

// Inc
func Inc(counter prometheus.Counter) {
	if config.Config.Prometheus.Enable {
		if counter != nil {
			counter.Inc()
		}
	}
}

// Add
func Add(counter prometheus.Counter, add int) {
	if config.Config.Prometheus.Enable {
		if counter != nil {
			counter.Add(float64(add))
		}
	}
}

// GaugeInc
func GaugeInc(gauges prometheus.Gauge) {
	if config.Config.Prometheus.Enable {
		if gauges != nil {
			gauges.Inc()
		}
	}
}

// GaugeDec
func GaugeDec(gauges prometheus.Gauge) {
	if config.Config.Prometheus.Enable {
		if gauges != nil {
			gauges.Dec()
		}
	}
}
