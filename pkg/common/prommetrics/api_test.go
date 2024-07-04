package prommetrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"testing"
	"time"
)

var (
	apiCallCounter1 = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_calls_total",
			Help: "Total number of API calls",
		},
		[]string{"endpoint", "status", "code", "error"},
	)
	registerer *prometheus.Registry
)

func init() {
	registerer = prometheus.NewRegistry()
	registerer.MustRegister(apiCallCounter1)
}

func recordAPICall(endpoint string, status string) {
	apiCallCounter1.With(prometheus.Labels{"endpoint": endpoint, "status": status, "code": "200", "error": "ArgsError"}).Inc()
}

func TestName(t *testing.T) {
	go func() {
		for i := 0; ; i++ {
			recordAPICall("/api/test", "success")
			time.Sleep(time.Second)
		}
	}()

	go func() {
		for i := 0; ; i++ {
			recordAPICall("/api/test", "failed")
			time.Sleep(time.Second * 3)
		}
	}()
	http.Handle("/metrics", promhttp.HandlerFor(registerer, promhttp.HandlerOpts{}))
	if err := http.ListenAndServe(":2112", nil); err != nil {
		panic(err)
	}
}

func TestName2(t *testing.T) {
	var d time.Duration
	d = time.Second / 900
	fmt.Println(durationRange(d))
}
