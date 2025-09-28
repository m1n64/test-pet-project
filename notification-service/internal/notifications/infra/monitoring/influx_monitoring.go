package monitoring

import (
	"context"
	"go.uber.org/zap"
	"notification-service-api/internal/notifications/domain"
	"notification-service-api/pkg/utils"
	"sync"
	"sync/atomic"
	"time"
)

type seriesKey struct {
	Channel string
	Type    string
	Env     string
}

type InfluxMonitoring struct {
	influxClient *utils.InfluxDB
	logger       *zap.Logger
	totals       sync.Map
	env          string
}

func NewInfluxMonitoring(influxClient *utils.InfluxDB, logger *zap.Logger, env string) *InfluxMonitoring {
	return &InfluxMonitoring{
		influxClient: influxClient,
		logger:       logger,
		env:          env,
	}
}

func (i *InfluxMonitoring) Send(channel domain.Channel, notificationType domain.NotificationType, delta int64) {
	key := seriesKey{
		Channel: channel.String(),
		Type:    notificationType.String(),
		Env:     i.env,
	}

	ctrAny, _ := i.totals.LoadOrStore(key, new(atomic.Int64))
	ctr := ctrAny.(*atomic.Int64)

	val := ctr.Add(delta)

	tags := map[string]string{
		"channel": key.Channel,
		"type":    key.Type,
		"env":     key.Env,
	}
	fields := map[string]interface{}{
		"total": val,
	}

	if err := i.influxClient.Send("notification_stats", tags, fields, time.Now().UnixNano()); err != nil {
		i.logger.Error("Send notifications stats to Influx error", zap.Error(err))
	}
}

func (i *InfluxMonitoring) SendSuccess(channel domain.Channel, count int64) {
	i.Send(channel, domain.NotificationTypeSuccess, count)
}

func (i *InfluxMonitoring) SendError(channel domain.Channel, count int64) {
	i.Send(channel, domain.NotificationTypeError, count)
}

func (i *InfluxMonitoring) Snapshot(ctx context.Context) map[seriesKey]int64 {
	out := make(map[seriesKey]int64)
	i.totals.Range(func(k, v any) bool {
		out[k.(seriesKey)] = v.(*atomic.Int64).Load()
		return true
	})
	return out
}
