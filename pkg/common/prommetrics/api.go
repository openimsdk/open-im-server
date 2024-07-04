package prommetrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strconv"
	"time"
)

const ApiPath = "/metrics"

var (
	apiRegistry = prometheus.NewRegistry()
	apiCounter  = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_count",
			Help: "Total number of API calls",
		},
		[]string{"path", "code"},
	)
	httpCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_count",
			Help: "Total number of HTTP calls",
		},
		//[]string{"path", "method", "status", "duration"},
		[]string{"path", "method", "status"},
	)
)

func init() {
	apiRegistry.MustRegister(apiCounter, httpCounter)
}

func APICall(path string, apiCode int) {
	apiCounter.With(prometheus.Labels{"path": path, "code": strconv.Itoa(apiCode)}).Inc()
}

//func HttpCall(path string, method string, status int, duration time.Duration) {
//	httpCounter.With(prometheus.Labels{"path": path, "method": method, "status": strconv.Itoa(status), "duration": durationRange(duration)}).Inc()
//}

func HttpCall(path string, method string, status int) {
	httpCounter.With(prometheus.Labels{"path": path, "method": method, "status": strconv.Itoa(status)}).Inc()
}

var (
	durations = [...]time.Duration{
		time.Millisecond * 1,
		time.Millisecond * 2,
		time.Millisecond * 3,
		time.Millisecond * 4,
		time.Millisecond * 5,
		time.Millisecond * 6,
		time.Millisecond * 7,
		time.Millisecond * 8,
		time.Millisecond * 9,
		time.Millisecond * 10,
		time.Millisecond * 20,
		time.Millisecond * 30,
		time.Millisecond * 40,
		time.Millisecond * 50,
		time.Millisecond * 60,
		time.Millisecond * 70,
		time.Millisecond * 80,
		time.Millisecond * 90,
		time.Millisecond * 100,
		time.Millisecond * 200,
		time.Millisecond * 300,
		time.Millisecond * 400,
		time.Millisecond * 500,
		time.Millisecond * 600,
		time.Millisecond * 700,
		time.Millisecond * 800,
		time.Millisecond * 900,
		time.Second * 1,
		time.Second * 2,
		time.Second * 3,
		time.Second * 4,
		time.Second * 5,
		time.Second * 6,
		time.Second * 7,
		time.Second * 8,
		time.Second * 9,
		time.Second * 10,
		time.Second * 20,
		time.Second * 30,
		time.Second * 40,
		time.Second * 50,
		time.Second * 60,
		time.Second * 70,
		time.Second * 80,
		time.Second * 90,
		time.Second * 100,
	}
	maxDuration = durations[len(durations)-1]
)

func durationRange(duration time.Duration) string {
	for _, d := range durations {
		if duration <= d {
			return d.String()
		}
	}
	return ">" + maxDuration.String()
}

func ApiHandler() http.Handler {
	return promhttp.HandlerFor(apiRegistry, promhttp.HandlerOpts{})
}
