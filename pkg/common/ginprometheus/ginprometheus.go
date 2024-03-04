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

package ginprometheus

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var defaultMetricPath = "/metrics"

// counter, counter_vec, gauge, gauge_vec,
// histogram, histogram_vec, summary, summary_vec.
var (
	reqCounter = &Metric{
		ID:          "reqCnt",
		Name:        "requests_total",
		Description: "How many HTTP requests processed, partitioned by status code and HTTP method.",
		Type:        "counter_vec",
		Args:        []string{"code", "method", "handler", "host", "url"}}

	reqDuration = &Metric{
		ID:          "reqDur",
		Name:        "request_duration_seconds",
		Description: "The HTTP request latencies in seconds.",
		Type:        "histogram_vec",
		Args:        []string{"code", "method", "url"},
	}

	resSize = &Metric{
		ID:          "resSz",
		Name:        "response_size_bytes",
		Description: "The HTTP response sizes in bytes.",
		Type:        "summary"}

	reqSize = &Metric{
		ID:          "reqSz",
		Name:        "request_size_bytes",
		Description: "The HTTP request sizes in bytes.",
		Type:        "summary"}

	standardMetrics = []*Metric{
		reqCounter,
		reqDuration,
		resSize,
		reqSize,
	}
)

/*
RequestCounterURLLabelMappingFn is a function which can be supplied to the middleware to control
the cardinality of the request counter's "url" label, which might be required in some contexts.
For instance, if for a "/customer/:name" route you don't want to generate a time series for every
possible customer name, you could use this function:

	func(c *gin.Context) string {
		url := c.Request.URL.Path
		for _, p := range c.Params {
			if p.Key == "name" {
				url = strings.Replace(url, p.Value, ":name", 1)
				break
			}
		}
		return url
	}

which would map "/customer/alice" and "/customer/bob" to their template "/customer/:name".
*/
type RequestCounterURLLabelMappingFn func(c *gin.Context) string

// Metric is a definition for the name, description, type, ID, and
// prometheus.Collector type (i.e. CounterVec, Summary, etc) of each metric.
type Metric struct {
	MetricCollector prometheus.Collector
	ID              string
	Name            string
	Description     string
	Type            string
	Args            []string
}

// Prometheus contains the metrics gathered by the instance and its path.
type Prometheus struct {
	reqCnt        *prometheus.CounterVec
	reqDur        *prometheus.HistogramVec
	reqSz, resSz  prometheus.Summary
	router        *gin.Engine
	listenAddress string
	Ppg           PrometheusPushGateway

	MetricsList []*Metric
	MetricsPath string

	ReqCntURLLabelMappingFn RequestCounterURLLabelMappingFn

	// gin.Context string to use as a prometheus URL label
	URLLabelFromContext string
}

// PrometheusPushGateway contains the configuration for pushing to a Prometheus pushgateway (optional).
type PrometheusPushGateway struct {

	// Push interval in seconds
	PushIntervalSeconds time.Duration

	// Push Gateway URL in format http://domain:port
	// where JOBNAME can be any string of your choice
	PushGatewayURL string

	// Local metrics URL where metrics are fetched from, this could be omitted in the future
	// if implemented using prometheus common/expfmt instead
	MetricsURL string

	// pushgateway job name, defaults to "gin"
	Job string
}

// NewPrometheus generates a new set of metrics with a certain subsystem name.
func NewPrometheus(subsystem string, customMetricsList ...[]*Metric) *Prometheus {
	if subsystem == "" {
		subsystem = "app"
	}

	var metricsList []*Metric

	if len(customMetricsList) > 1 {
		panic("Too many args. NewPrometheus( string, <optional []*Metric> ).")
	} else if len(customMetricsList) == 1 {
		metricsList = customMetricsList[0]
	}
	metricsList = append(metricsList, standardMetrics...)

	p := &Prometheus{
		MetricsList: metricsList,
		MetricsPath: defaultMetricPath,
		ReqCntURLLabelMappingFn: func(c *gin.Context) string {
			return c.FullPath() // e.g. /user/:id , /user/:id/info
		},
	}

	p.registerMetrics(subsystem)

	return p
}

// SetPushGateway sends metrics to a remote pushgateway exposed on pushGatewayURL
// every pushIntervalSeconds. Metrics are fetched from metricsURL.
func (p *Prometheus) SetPushGateway(pushGatewayURL, metricsURL string, pushIntervalSeconds time.Duration) {
	p.Ppg.PushGatewayURL = pushGatewayURL
	p.Ppg.MetricsURL = metricsURL
	p.Ppg.PushIntervalSeconds = pushIntervalSeconds
	p.startPushTicker()
}

// SetPushGatewayJob job name, defaults to "gin".
func (p *Prometheus) SetPushGatewayJob(j string) {
	p.Ppg.Job = j
}

// SetListenAddress for exposing metrics on address. If not set, it will be exposed at the
// same address of the gin engine that is being used.
func (p *Prometheus) SetListenAddress(address string) {
	p.listenAddress = address
	if p.listenAddress != "" {
		p.router = gin.Default()
	}
}

// SetListenAddressWithRouter for using a separate router to expose metrics. (this keeps things like GET /metrics out of
// your content's access log).
func (p *Prometheus) SetListenAddressWithRouter(listenAddress string, r *gin.Engine) {
	p.listenAddress = listenAddress
	if len(p.listenAddress) > 0 {
		p.router = r
	}
}

// SetMetricsPath set metrics paths.
func (p *Prometheus) SetMetricsPath(e *gin.Engine) error {

	if p.listenAddress != "" {
		p.router.GET(p.MetricsPath, prometheusHandler())
		return p.runServer()
	} else {
		e.GET(p.MetricsPath, prometheusHandler())
		return nil
	}
}

// SetMetricsPathWithAuth set metrics paths with authentication.
func (p *Prometheus) SetMetricsPathWithAuth(e *gin.Engine, accounts gin.Accounts) error {

	if p.listenAddress != "" {
		p.router.GET(p.MetricsPath, gin.BasicAuth(accounts), prometheusHandler())
		return p.runServer()
	} else {
		e.GET(p.MetricsPath, gin.BasicAuth(accounts), prometheusHandler())
		return nil
	}

}

func (p *Prometheus) runServer() error {
	return p.router.Run(p.listenAddress)
}

func (p *Prometheus) getMetrics() []byte {
	response, err := http.Get(p.Ppg.MetricsURL)
	if err != nil {
		return nil
	}

	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)
	return body
}

var hostname, _ = os.Hostname()

func (p *Prometheus) getPushGatewayURL() string {
	if p.Ppg.Job == "" {
		p.Ppg.Job = "gin"
	}
	return p.Ppg.PushGatewayURL + "/metrics/job/" + p.Ppg.Job + "/instance/" + hostname
}

func (p *Prometheus) sendMetricsToPushGateway(metrics []byte) {
	req, err := http.NewRequest("POST", p.getPushGatewayURL(), bytes.NewBuffer(metrics))
	if err != nil {
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending to push gateway error:", err.Error())
	}

	resp.Body.Close()
}

func (p *Prometheus) startPushTicker() {
	ticker := time.NewTicker(time.Second * p.Ppg.PushIntervalSeconds)
	go func() {
		for range ticker.C {
			p.sendMetricsToPushGateway(p.getMetrics())
		}
	}()
}

// NewMetric associates prometheus.Collector based on Metric.Type.
func NewMetric(m *Metric, subsystem string) prometheus.Collector {
	var metric prometheus.Collector
	switch m.Type {
	case "counter_vec":
		metric = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
			m.Args,
		)
	case "counter":
		metric = prometheus.NewCounter(
			prometheus.CounterOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
		)
	case "gauge_vec":
		metric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
			m.Args,
		)
	case "gauge":
		metric = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
		)
	case "histogram_vec":
		metric = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
			m.Args,
		)
	case "histogram":
		metric = prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
		)
	case "summary_vec":
		metric = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
			m.Args,
		)
	case "summary":
		metric = prometheus.NewSummary(
			prometheus.SummaryOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
		)
	}
	return metric
}

func (p *Prometheus) registerMetrics(subsystem string) {
	for _, metricDef := range p.MetricsList {
		metric := NewMetric(metricDef, subsystem)
		if err := prometheus.Register(metric); err != nil {
			fmt.Println("could not be registered in Prometheus,metricDef.Name:", metricDef.Name, "   error:", err.Error())
		}

		switch metricDef {
		case reqCounter:
			p.reqCnt = metric.(*prometheus.CounterVec)
		case reqDuration:
			p.reqDur = metric.(*prometheus.HistogramVec)
		case resSize:
			p.resSz = metric.(prometheus.Summary)
		case reqSize:
			p.reqSz = metric.(prometheus.Summary)
		}
		metricDef.MetricCollector = metric
	}
}

// Use adds the middleware to a gin engine.
func (p *Prometheus) Use(e *gin.Engine) error {
	e.Use(p.HandlerFunc())
	return p.SetMetricsPath(e)
}

// UseWithAuth adds the middleware to a gin engine with BasicAuth.
func (p *Prometheus) UseWithAuth(e *gin.Engine, accounts gin.Accounts) error {
	e.Use(p.HandlerFunc())
	return p.SetMetricsPathWithAuth(e, accounts)
}

// HandlerFunc defines handler function for middleware.
func (p *Prometheus) HandlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == p.MetricsPath {
			c.Next()
			return
		}

		start := time.Now()
		reqSz := computeApproximateRequestSize(c.Request)

		c.Next()

		status := strconv.Itoa(c.Writer.Status())
		elapsed := float64(time.Since(start)) / float64(time.Second)
		resSz := float64(c.Writer.Size())

		url := p.ReqCntURLLabelMappingFn(c)
		if len(p.URLLabelFromContext) > 0 {
			u, found := c.Get(p.URLLabelFromContext)
			if !found {
				u = "unknown"
			}
			url = u.(string)
		}
		p.reqDur.WithLabelValues(status, c.Request.Method, url).Observe(elapsed)
		p.reqCnt.WithLabelValues(status, c.Request.Method, c.HandlerName(), c.Request.Host, url).Inc()
		p.reqSz.Observe(float64(reqSz))
		p.resSz.Observe(resSz)
	}
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func computeApproximateRequestSize(r *http.Request) int {
	var s int
	if r.URL != nil {
		s = len(r.URL.Path)
	}

	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)

	// r.FormData and r.MultipartForm are assumed to be included in r.URL.

	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}
	return s
}
