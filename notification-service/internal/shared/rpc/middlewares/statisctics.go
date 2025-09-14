package middlewares

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"math"
	"notification-service-api/pkg/utils"
	"os"
	"sync/atomic"
	"time"
)

var inFlight atomic.Int64

func StatisticsMiddleware(influx *utils.InfluxDB, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		_ = inFlight.Add(1)

		c.Next()

		lat := time.Since(start)
		cur := inFlight.Add(-1)

		tags := map[string]string{
			"service": "notification",
			"env":     os.Getenv("SERVICE_ENV"),
		}

		fields := map[string]interface{}{
			"counter":       1,
			"duration_ms":   float64(lat.Milliseconds()),
			"bytes_written": int64(math.Max(0.0, float64(c.Writer.Size()))),
			"inflight":      cur,
		}

		err := influx.Send("notification_rpc", tags, fields, time.Now().UnixNano())
		if err != nil {
			logger.Error("Send to Influx error", zap.Error(err))
		}
	}
}
