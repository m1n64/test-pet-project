package utils

import (
	"context"
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"log"
	"notification-service-api/internal/shared/queue/notifications"
	"sync"
	"sync/atomic"
	"time"
)

type RabbitMQConnection struct {
	conn   *amqp.Connection
	chans  []*amqp.Channel
	idx    uint64
	mu     sync.Mutex
	dsn    string
	logger *zap.Logger
}

type HandlerFunc func(ctx context.Context, d amqp.Delivery) error

type ConsumeOptions struct {
	Queue        string
	Workers      int
	Prefetch     int
	ConsumerTag  string
	Args         amqp.Table
	RetryBackoff time.Duration
	RetryMax     int
}

var (
	instance *RabbitMQConnection
	once     sync.Once
)

func ConnectRabbitMQ(rabbitURL string, logger *zap.Logger) *RabbitMQConnection {
	poolSize := 32

	once.Do(func() {
		instance = &RabbitMQConnection{dsn: rabbitURL, logger: logger}
		if err := instance.connect(poolSize); err != nil {
			log.Fatalf("RabbitMQ: initial connect failed: %v", err)
		}
		log.Printf("RabbitMQ connection established")
	})
	return instance
}

func GetRabbitMQInstance() *RabbitMQConnection {
	if instance == nil {
		log.Fatalf("RabbitMQ connection is not initialized. Call ConnectRabbitMQ first.")
	}
	return instance
}

func (r *RabbitMQConnection) Publish(ctx context.Context, exchange, routingKey string, msg amqp.Publishing) error {
	ch, ok := r.nextChan()
	if !ok {
		return amqp.ErrClosed
	}

	errCh := make(chan error, 1)
	go func(ch *amqp.Channel) {
		errCh <- ch.Publish(exchange, routingKey, false, false, msg)
	}(ch)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		if err != nil {
			r.reopenChannel(ch)
			return err
		}
		return nil
	}
}

func (r *RabbitMQConnection) PublishMsgpack(ctx context.Context, exchange, routingKey string, body []byte, headers amqp.Table, correlationID *string) error {
	corr := ""
	if correlationID != nil {
		corr = *correlationID
	}
	msg := amqp.Publishing{
		DeliveryMode:  amqp.Persistent,
		ContentType:   "application/x-msgpack",
		Body:          body,
		Headers:       headers,
		Timestamp:     time.Now(),
		CorrelationId: corr,
	}
	return r.Publish(ctx, exchange, routingKey, msg)
}

func (r *RabbitMQConnection) Consume(ctx context.Context, opts ConsumeOptions, handler HandlerFunc) error {
	if opts.Workers <= 0 {
		opts.Workers = 1
	}
	if opts.Prefetch <= 0 {
		opts.Prefetch = opts.Workers
	}
	if opts.RetryBackoff <= 0 {
		opts.RetryBackoff = time.Second
	}

	errCh := make(chan error, opts.Workers)

	for i := 0; i < opts.Workers; i++ {
		go r.runConsumerWorker(ctx, i, opts, handler, errCh)
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

func (r *RabbitMQConnection) Channel() (*amqp.Channel, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.conn == nil {
		return nil, amqp.ErrClosed
	}
	return r.conn.Channel()
}

func (r *RabbitMQConnection) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.conn != nil {
		_ = r.conn.Close()
	}
}

func (r *RabbitMQConnection) ConsumeSimple(ctx context.Context, queue string, workers, prefetch int, handler HandlerFunc) error {
	return r.Consume(ctx, ConsumeOptions{
		Queue:    queue,
		Workers:  workers,
		Prefetch: prefetch,
		Args:     nil,
	}, handler)
}

func (r *RabbitMQConnection) IsConnected() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.conn != nil && !r.conn.IsClosed()
}

func (r *RabbitMQConnection) connect(poolSize int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.conn != nil {
		for _, ch := range r.chans {
			_ = ch.Close()
		}
		_ = r.conn.Close()
		r.chans = nil
	}

	conn, err := amqp.Dial(r.dsn)
	if err != nil {
		return err
	}
	r.conn = conn

	chans := make([]*amqp.Channel, 0, poolSize)
	for i := 0; i < poolSize; i++ {
		ch, err := conn.Channel()
		if err != nil {
			for _, c := range chans {
				_ = c.Close()
			}
			_ = conn.Close()
			return err
		}
		chans = append(chans, ch)
	}
	r.chans = chans
	atomic.StoreUint64(&r.idx, 0)

	go r.watchConn(poolSize)

	return nil
}

func (r *RabbitMQConnection) watchConn(poolSize int) {
	errCh := r.conn.NotifyClose(make(chan *amqp.Error, 1))
	if err := <-errCh; err != nil {
		log.Printf("RabbitMQ connection closed, reason: %v", err)

		backoff := time.Second
		for {
			time.Sleep(backoff)
			if backoff < 30*time.Second {
				backoff *= 2
			}
			if err := r.connect(poolSize); err != nil {
				log.Printf("RabbitMQ reconnect failed: %v (retrying)", err)
				continue
			}
			log.Printf("RabbitMQ reconnected, pool rebuilt: %d channels", poolSize)
			return
		}
	}
}

func (r *RabbitMQConnection) nextChan() (*amqp.Channel, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.conn == nil || len(r.chans) == 0 {
		return nil, false
	}
	i := atomic.AddUint64(&r.idx, 1)
	return r.chans[int(i)%len(r.chans)], true
}

func (r *RabbitMQConnection) reopenChannel(bad *amqp.Channel) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.conn == nil {
		return
	}
	_ = bad.Close()
	nch, err := r.conn.Channel()
	if err != nil {
		log.Printf("RabbitMQ: reopen channel failed: %v", err)
		return
	}
	for i, ch := range r.chans {
		if ch == bad {
			r.chans[i] = nch
			return
		}
	}
	_ = nch.Close()
}

func (r *RabbitMQConnection) runConsumerWorker(
	ctx context.Context,
	workerID int,
	opts ConsumeOptions,
	handler HandlerFunc,
	errCh chan<- error,
) {
	backoff := opts.RetryBackoff

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		r.mu.Lock()
		conn := r.conn
		r.mu.Unlock()

		if conn == nil {
			time.Sleep(backoff)
			backoff = nextBackoff(backoff)
			continue
		}

		ch, err := conn.Channel()
		if err != nil {
			r.logger.Info(fmt.Sprintf("[consumer:%d] channel open failed: %v", workerID, err))
			time.Sleep(backoff)
			backoff = nextBackoff(backoff)
			continue
		}

		if err := ch.Qos(opts.Prefetch, 0, false); err != nil {
			_ = ch.Close()
			r.logger.Error(fmt.Sprintf("[consumer:%d] qos failed: %v", workerID, err))
			time.Sleep(backoff)
			backoff = nextBackoff(backoff)
			continue
		}

		msgs, err := ch.Consume(
			opts.Queue,
			opts.ConsumerTag, // consumer tag
			false,            // auto-ack = false
			false,            // exclusive
			false,            // no-local (ignored by RabbitMQ)
			false,            // no-wait
			opts.Args,
		)
		if err != nil {
			_ = ch.Close()
			r.logger.Error(fmt.Sprintf("[consumer:%d] consume failed: %v", workerID, err))
			time.Sleep(backoff)
			backoff = nextBackoff(backoff)
			continue
		}

		notifyClose := ch.NotifyClose(make(chan *amqp.Error, 1))

		r.logger.Info(fmt.Sprintf("[consumer:%d] started on queue=%s prefetch=%d", workerID, opts.Queue, opts.Prefetch))

	consumeLoop:
		for {
			select {
			case <-ctx.Done():
				_ = ch.Close()
				return
			case amqpErr := <-notifyClose:
				if amqpErr != nil {
					r.logger.Error(fmt.Sprintf("[consumer:%d] channel closed: %v", workerID, amqpErr))
				} else {
					r.logger.Info(fmt.Sprintf("[consumer:%d] channel closed", workerID))
				}
				break consumeLoop

			case d, ok := <-msgs:
				localLogger := r.logger.With(zap.String("request_id", d.CorrelationId))
				if !ok {
					localLogger.Warn(fmt.Sprintf("[consumer:%d] msgs channel closed by broker", workerID))
					break consumeLoop
				}
				attempts := getRetryCount(d.Headers) + 1

				// TODO: finish DLQ and Retry queues
				if err := handler(ctx, d); err != nil {
					if attempts >= int64(opts.RetryMax) {
						pubErr := r.Publish(ctx, notifications.ExchangeDLX, opts.Queue+".dlq", amqp.Publishing{
							DeliveryMode:  amqp.Persistent,
							ContentType:   d.ContentType,
							Body:          d.Body,
							Headers:       withRetryCount(d.Headers, attempts),
							CorrelationId: d.CorrelationId,
							MessageId:     d.MessageId,
							Timestamp:     time.Now(),
						})
						if pubErr != nil {
							_ = d.Nack(false, true)
							localLogger.Error(fmt.Sprintf("[consumer:%d] publish to DLX failed: %v", workerID, pubErr))
							continue
						}
						_ = d.Ack(false)
						localLogger.Info(fmt.Sprintf("[consumer:%d] publish to DLX succeeded", workerID))

						continue
					}

					pubErr := r.Publish(ctx, notifications.ExchangeRetry, opts.Queue+".retry", amqp.Publishing{
						DeliveryMode:  amqp.Persistent,
						ContentType:   d.ContentType,
						Body:          d.Body,
						Headers:       withRetryCount(d.Headers, attempts),
						CorrelationId: d.CorrelationId,
						MessageId:     d.MessageId,
						Timestamp:     time.Now(),
					})
					if pubErr != nil {
						_ = d.Nack(false, true)
						localLogger.Error(fmt.Sprintf("[consumer:%d] publish to retry failed: %v", workerID, pubErr))
						continue
					}

					_ = d.Ack(false)
					localLogger.Info(fmt.Sprintf("[consumer:%d] publish to retry succeeded", workerID))
					continue
				}

				_ = d.Ack(false)
			}
		}

		_ = ch.Close()
		time.Sleep(backoff)
		backoff = nextBackoff(backoff)
	}
}

func nextBackoff(cur time.Duration) time.Duration {
	next := time.Duration(float64(cur) * 2)
	if next > 30*time.Second {
		next = 30 * time.Second
	}
	return next
}

func getRetryCount(h amqp.Table) int64 {
	if h == nil {
		return 0
	}
	switch v := h["x-retry-count"].(type) {
	case int64:
		return v
	case int32:
		return int64(v)
	case float64:
		return int64(v)
	}
	return 0
}

func withRetryCount(h amqp.Table, n int64) amqp.Table {
	nh := amqp.Table{}
	for k, v := range h {
		nh[k] = v
	}
	nh["x-retry-count"] = n
	return nh
}

var ErrNotConnected = errors.New("rabbitmq: not connected")
