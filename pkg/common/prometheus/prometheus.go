package prometheus

import (
	"Open_IM/pkg/common/config"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartPromeSrv(promethuesPort int) error {
	if config.Config.Prometheus.Enable {
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(":"+strconv.Itoa(promethuesPort), nil)
		return err
	}
	return nil
}

func PrometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func PromeInc(counter prometheus.Counter) {
	if config.Config.Prometheus.Enable {
		counter.Inc()
	}
}

func PromeAdd(counter prometheus.Counter, add int) {
	if config.Config.Prometheus.Enable {
		counter.Add(float64(add))
	}
}
