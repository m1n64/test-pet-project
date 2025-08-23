package utils

import (
	"github.com/streadway/amqp"
	"log"
	"sync"
	"time"
)

type RabbitMQConnection struct {
	conn *amqp.Connection
	mu   sync.Mutex
	dsn  string
}

var (
	instance *RabbitMQConnection
	once     sync.Once
)

func ConnectRabbitMQ(rabbitURL string) *RabbitMQConnection {
	once.Do(func() {
		instance = &RabbitMQConnection{dsn: rabbitURL}
		if err := instance.connect(); err != nil {
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

func (r *RabbitMQConnection) connect() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	conn, err := amqp.Dial(r.dsn)
	if err != nil {
		return err
	}
	r.conn = conn

	go func() {
		errCh := conn.NotifyClose(make(chan *amqp.Error))
		if err := <-errCh; err != nil {
			log.Printf("RabbitMQ connection closed, reason: %v", err)
			for {
				time.Sleep(5 * time.Second)
				if err := r.connect(); err == nil {
					log.Println("RabbitMQ reconnected")
					return
				}
			}
		}
	}()

	return nil
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
