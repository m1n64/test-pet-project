package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgraph-io/ristretto"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"sync"
	"time"
)

type MultiCache struct {
	logger      *zap.Logger
	l2          *redis.Client
	l1          *ristretto.Cache
	defaultTTL  time.Duration
	costPerItem int64
	ctx         context.Context
	mu          sync.RWMutex
	keys        map[string]struct{}
	until       map[string]time.Time
	keySet      string
}

var ErrMultiCacheValueBinary = errors.New("error to convert value to binary")

func NewMultiCache(redis *redis.Client, logger *zap.Logger, maxItems int64) (*MultiCache, error) {
	if maxItems <= 0 {
		maxItems = 1000000
	}

	l1, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: maxItems * 10,
		MaxCost:     maxItems,
		BufferItems: 64,
	})
	if err != nil {
		return nil, err
	}

	return &MultiCache{
		logger:      logger,
		l2:          redis,
		l1:          l1,
		costPerItem: 1,
		ctx:         context.Background(),
		keys:        make(map[string]struct{}),
		until:       make(map[string]time.Time),
		keySet:      "app:cache:keys",
	}, nil
}

func (c *MultiCache) Set(key string, value interface{}, ttl time.Duration) bool {
	var isSet bool
	if ttl > 0 {
		isSet = c.l1.SetWithTTL(key, value, c.costPerItem, ttl)
	} else {
		isSet = c.l1.Set(key, value, c.costPerItem)
	}

	c.l1.Wait()

	c.mu.Lock()
	c.keys[key] = struct{}{}
	if ttl > 0 {
		c.until[key] = time.Now().Add(ttl)
	} else {
		delete(c.until, key)
	}
	c.mu.Unlock()

	return isSet
}

func (c *MultiCache) Get(key string) (MultiCacheValue, bool, error) {
	if val, ok := c.l1.Get(key); ok {
		return MultiCacheValue{
			v: val,
		}, true, nil
	}

	obj, ttl, ok, err := c.getFromRedis(key)
	if err != nil {
		return MultiCacheValue{}, false, err
	}
	if ok {
		c.Set(key, obj, ttl)

		return MultiCacheValue{v: obj}, true, nil
	}

	return MultiCacheValue{}, false, nil
}

func (c *MultiCache) Del(key string) error {
	c.l1.Del(key)

	err := c.l2.Del(c.ctx, key).Err()

	c.mu.Lock()
	delete(c.keys, key)
	delete(c.until, key)
	c.mu.Unlock()

	return err
}

func (c *MultiCache) WarmFromRedis() error {
	const page = int64(100)
	var cursor uint64

	c.mu.Lock()
	c.keys = make(map[string]struct{})
	c.until = make(map[string]time.Time)
	c.mu.Unlock()

	for {
		keys, next, err := c.l2.Scan(c.ctx, cursor, c.keySet, page).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			return err
		}

		for _, key := range keys {
			obj, ttl, ok, err := c.getFromRedis(key)
			if err != nil || !ok {
				continue
			}

			c.Set(key, obj, ttl)

			if err := c.l2.Del(c.ctx, key).Err(); err != nil {
				c.logger.Error(fmt.Sprintf("[MultiCache] Error deleting %s key from Redis", key), zap.Error(err))
			}
		}

		cursor = next
		if cursor == 0 {
			break
		}
	}

	return c.l2.Del(c.ctx, c.keySet).Err()
}

func (c *MultiCache) FlushToRedisOnce() error {
	c.mu.RLock()
	keys := make([]string, 0, len(c.keys))
	for key := range c.keys {
		keys = append(keys, key)
	}
	c.mu.RUnlock()

	for _, key := range keys {
		if val, ok := c.l1.Get(key); ok {
			ttl := c.ttlLeft(key)

			err := c.setOnRedis(key, val, ttl)
			if err != nil {
				return err
			}

			_ = c.l2.SAdd(c.ctx, c.keySet, key).Err()
		}
	}

	return nil
}

func (c *MultiCache) StartAutoFlush(interval time.Duration) context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := c.FlushToRedisOnce(); err != nil {
					c.logger.Error("[MultiCache] Error flushing to Redis", zap.Error(err))
				}
			}
		}
	}()

	return cancel
}

func (c *MultiCache) setOnRedis(key string, value interface{}, ttl time.Duration) error {
	blob, err := Pack(value)
	if err != nil {
		return err
	}

	if err := c.l2.Set(c.ctx, key, blob, ttl).Err(); err != nil {
		return err
	}

	return nil
}

func (c *MultiCache) getFromRedis(key string) (any, time.Duration, bool, error) {
	pipe := c.l2.Pipeline()
	pttlCmd := pipe.PTTL(c.ctx, key)
	getCmd := pipe.Get(c.ctx, key)
	_, _ = pipe.Exec(c.ctx)

	if err := getCmd.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, 0, false, nil
		}
		return nil, 0, false, err
	}

	b, err := getCmd.Bytes()
	if err != nil {
		return nil, 0, false, err
	}
	ttl := pttlCmd.Val()

	val, _, err := Unpack(b)
	if err != nil {
		return nil, ttl, false, err
	}

	return val, ttl, true, nil
}

func (c *MultiCache) ttlLeft(key string) time.Duration {
	c.mu.RLock()
	until, ok := c.until[key]
	c.mu.RUnlock()
	if !ok || until.IsZero() {
		return 0
	}
	d := time.Until(until)
	if d < 0 {
		return 0
	}
	return d
}
