package api

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/openimsdk/tools/apiresp"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/stability/ratelimit"
	"github.com/openimsdk/tools/stability/ratelimit/bbr"
)

type RateLimiter struct {
	Enable       bool          `yaml:"enable"`
	Window       time.Duration `yaml:"window"`       // time duration per window
	Bucket       int           `yaml:"bucket"`       // bucket number for each window
	CPUThreshold int64         `yaml:"cpuThreshold"` // CPU threshold; valid range 0–1000 (1000 = 100%)
}

func RateLimitMiddleware(config *RateLimiter) gin.HandlerFunc {
	if !config.Enable {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	limiter := bbr.NewBBRLimiter(
		bbr.WithWindow(config.Window),
		bbr.WithBucket(config.Bucket),
		bbr.WithCPUThreshold(config.CPUThreshold),
	)

	return func(c *gin.Context) {
		status := limiter.Stat()

		c.Header("X-BBR-CPU", strconv.FormatInt(status.CPU, 10))
		c.Header("X-BBR-MinRT", strconv.FormatInt(status.MinRt, 10))
		c.Header("X-BBR-MaxPass", strconv.FormatInt(status.MaxPass, 10))
		c.Header("X-BBR-MaxInFlight", strconv.FormatInt(status.MaxInFlight, 10))
		c.Header("X-BBR-InFlight", strconv.FormatInt(status.InFlight, 10))

		done, err := limiter.Allow()
		if err != nil {

			c.Header("X-RateLimit-Policy", "BBR")
			c.Header("Retry-After", calculateBBRRetryAfter(status))
			c.Header("X-RateLimit-Limit", strconv.FormatInt(status.MaxInFlight, 10))
			c.Header("X-RateLimit-Remaining", "0") // There is no concept of remaining quota in BBR.

			fmt.Println("rate limited:", err, "path:", c.Request.URL.Path)
			log.ZWarn(c, "rate limited", err, "path", c.Request.URL.Path)
			c.AbortWithStatus(http.StatusTooManyRequests)
			apiresp.GinError(c, errs.NewCodeError(http.StatusTooManyRequests, "too many requests, please try again later"))
			return
		}

		c.Next()
		done(ratelimit.DoneInfo{})
	}
}

func calculateBBRRetryAfter(status bbr.Stat) string {
	loadRatio := float64(status.CPU) / float64(status.CPU)

	if loadRatio < 0.8 {
		return "1"
	}
	if loadRatio < 0.95 {
		return "2"
	}

	backoff := 1 + int64(math.Pow(loadRatio-0.95, 2)*50)
	if backoff > 5 {
		backoff = 5
	}
	return strconv.FormatInt(backoff, 10)
}
