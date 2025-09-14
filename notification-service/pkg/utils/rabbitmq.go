package utils

import (
	"context"
	"github.com/streadway/amqp"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type RabbitMQConnection struct {
	conn  *amqp.Connection
	chans []*amqp.Channel
	idx   uint64
	mu    sync.Mutex
	dsn   string
}

var (
	instance *RabbitMQConnection
	once     sync.Once
)

func ConnectRabbitMQ(rabbitURL string) *RabbitMQConnection {
	poolSize := 32

	once.Do(func() {
		instance = &RabbitMQConnection{dsn: rabbitURL}
		if err := instance.connect(poolSize); err != nil {
			log.Fatalf("RabbitMQ: initial connect failed: %v", err)
		}
		log.Println("RabbitMQ connection established")
	})
	return instance
}

func GetRabbitMQInstance() *RabbitMQConnection {
	if instance == nil {
		log.Fatalf("RabbitMQ connection is not initialized. Call ConnectRabbitMQ first.")
	}
	return instance
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
