package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/lgxju/gogretago/config"
	"github.com/redis/go-redis/v9"
)

// CacheService provides Redis caching operations
type CacheService struct {
	client  *redis.Client
	enabled bool
	prefix  string
	ttl     time.Duration
}

// NewCacheService creates and connects a Redis cache service
func NewCacheService() *CacheService {
	cfg := config.Get()

	if !cfg.CacheEnabled {
		log.Println("Cache disabled")
		return &CacheService{enabled: false}
	}

	opt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		log.Printf("Failed to parse Redis URL, cache disabled: %v", err)
		return &CacheService{enabled: false}
	}

	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("Failed to connect to Redis, cache disabled: %v", err)
		return &CacheService{enabled: false}
	}

	log.Println("Redis cache connected")
	return &CacheService{
		client:  client,
		enabled: true,
		prefix:  cfg.CacheKeyPrefix,
		ttl:     15 * time.Minute,
	}
}

// Get retrieves a cached value by key
func (c *CacheService) Get(ctx context.Context, key string, dest interface{}) (bool, error) {
	if !c.enabled {
		return false, nil
	}
	val, err := c.client.Get(ctx, c.prefix+key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return false, err
	}
	return true, nil
}

// Set stores a value in the cache
func (c *CacheService) Set(ctx context.Context, key string, value interface{}) error {
	if !c.enabled {
		return nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, c.prefix+key, data, c.ttl).Err()
}

// Delete removes a cached value by key
func (c *CacheService) Delete(ctx context.Context, key string) error {
	if !c.enabled {
		return nil
	}
	return c.client.Del(ctx, c.prefix+key).Err()
}

// InvalidatePattern removes all cached keys matching the pattern
func (c *CacheService) InvalidatePattern(ctx context.Context, pattern string) error {
	if !c.enabled {
		return nil
	}
	iter := c.client.Scan(ctx, 0, c.prefix+pattern, 100).Iterator()
	for iter.Next(ctx) {
		c.client.Del(ctx, iter.Val())
	}
	return iter.Err()
}

// BuildKey creates a cache key from parts
func BuildKey(parts ...interface{}) string {
	key := ""
	for i, p := range parts {
		if i > 0 {
			key += ":"
		}
		key += fmt.Sprintf("%v", p)
	}
	return key
}
