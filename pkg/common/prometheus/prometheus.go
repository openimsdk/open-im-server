package prometheus

import (
	"Open_IM/pkg/common/config"
	"bytes"
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

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func PromeTheusMiddleware(c *gin.Context) {
	PromeInc(ApiRequestCounter)
	w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
	c.Writer = w
	c.Next()
	if c.Writer.Status() == http.StatusOK {
		PromeInc(ApiRequestSuccessCounter)
	} else {
		PromeInc(ApiRequestFailedCounter)
	}
}

func PromeInc(counter prometheus.Counter) {
	if config.Config.Prometheus.Enable {
		if counter != nil {
			counter.Inc()
		}
	}
}

func PromeAdd(counter prometheus.Counter, add int) {
	if config.Config.Prometheus.Enable {
		if counter != nil {
			counter.Add(float64(add))
		}
	}
}

func PromeGaugeInc(gauges prometheus.Gauge) {
	if config.Config.Prometheus.Enable {
		if gauges != nil {
			gauges.Inc()
		}
	}
}

func PromeGaugeDec(gauges prometheus.Gauge) {
	if config.Config.Prometheus.Enable {
		if gauges != nil {
			gauges.Dec()
		}
	}
}
