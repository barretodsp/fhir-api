package utils

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Logger = logrus.Logger

func GinLogger(logger *Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		fields := logrus.Fields{
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       path,
			"ip":         c.ClientIP(),
			"user-agent": c.Request.UserAgent(),
			"latency":    latency,
			"time":       end.Format(time.RFC3339),
		}

		entry := logger.WithFields(fields)

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.String())
		} else if path == "/api/v1/health" || path == "/api/v1/metrics" {
			return
		} else {
			entry.Info("request")
		}
	}
}
