package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/aegis/ratelimit"
	"github.com/go-kratos/aegis/ratelimit/bbr"
	"github.com/openimsdk/tools/apiresp"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
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

	limiter := bbr.NewLimiter(
		bbr.WithWindow(config.Window),
		bbr.WithBucket(config.Bucket),
		bbr.WithCPUThreshold(config.CPUThreshold),
	)

	return func(c *gin.Context) {
		done, err := limiter.Allow()
		if err != nil {
			log.ZWarn(c, "rate limited", err, "path", c.Request.URL.Path)
			c.AbortWithStatus(http.StatusTooManyRequests)
			apiresp.GinError(c, errs.NewCodeError(http.StatusTooManyRequests, "too many requests, please try again later"))
			return
		}

		c.Next()

		done(ratelimit.DoneInfo{})
	}
}
