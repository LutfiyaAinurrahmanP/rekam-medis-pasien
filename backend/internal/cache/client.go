package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/redis/go-redis/v9"
)

// RedisClient wraps redis.Client with helper methods
type RedisClient struct {
	client *redis.Client
	cfg    *config.RedisConfig
}

// NewRedisClient creates and returns a connected RedisClient
func NewRedisClient(cfg *config.RedisConfig) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     10,
		MinIdleConns: 2,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis at %s:%s — %w", cfg.Host, cfg.Port, err)
	}

	log.Printf("✅ Redis connected at %s:%s (DB: %d)", cfg.Host, cfg.Port, cfg.DB)

	return &RedisClient{client: rdb, cfg: cfg}, nil
}

// Set encodes value to JSON and stores it with the given TTL.
// If ttl == 0, falls back to the default TTL from config.
func (r *RedisClient) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	if ttl == 0 {
		ttl = r.cfg.DefaultTTL
	}

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache set: marshal error for key %q: %w", key, err)
	}

	return r.client.Set(ctx, key, data, ttl).Err()
}

// Get retrieves the value for key and unmarshals it into dest.
// Returns ErrCacheMiss if the key does not exist.
func (r *RedisClient) Get(ctx context.Context, key string, dest any) error {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return ErrCacheMiss
		}
		return fmt.Errorf("cache get error for key %q: %w", key, err)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("cache get: unmarshal error for key %q: %w", key, err)
	}

	return nil
}

// Delete removes one or more keys from the cache.
func (r *RedisClient) Delete(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

// DeleteByPattern removes all keys matching the given glob pattern.
func (r *RedisClient) DeleteByPattern(ctx context.Context, pattern string) error {
	var cursor uint64
	for {
		keys, nextCursor, err := r.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return fmt.Errorf("cache scan error for pattern %q: %w", pattern, err)
		}
		if len(keys) > 0 {
			if err := r.client.Del(ctx, keys...).Err(); err != nil {
				return fmt.Errorf("cache delete error: %w", err)
			}
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return nil
}

// Exists returns true if the given key exists in the cache.
func (r *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	n, err := r.client.Exists(ctx, key).Result()
	return n > 0, err
}

// TTL returns the remaining TTL of the given key.
func (r *RedisClient) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, key).Result()
}

// FlushDB removes all keys in the current DB (use only in tests/dev).
func (r *RedisClient) FlushDB(ctx context.Context) error {
	return r.client.FlushDB(ctx).Err()
}

// Close closes the underlying Redis connection.
func (r *RedisClient) Close() error {
	return r.client.Close()
}

// Client returns the underlying *redis.Client (for advanced use).
func (r *RedisClient) Client() *redis.Client {
	return r.client
}

// ErrCacheMiss is returned by Get when the key does not exist.
var ErrCacheMiss = fmt.Errorf("cache miss")
