package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/AsaHero/dastyor-bot/pkg/config"
	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	client          *redis.Client
	storageDeadline time.Duration
	prefix          string
}

func NewRedisStorage(cfg *config.Config) (*RedisStorage, error) {
	dbNum, err := strconv.ParseInt(cfg.Redis.DB, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DB number for redis storage: %v", err)
	}

	storageDeadline, err := time.ParseDuration(cfg.Redis.StorageDeadline)
	if err != nil {
		return nil, fmt.Errorf("failed to parse storage deadline for redis: %v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       int(dbNum),
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return &RedisStorage{
		client:          client,
		prefix:          cfg.APP,
		storageDeadline: storageDeadline,
	}, nil
}

func (s *RedisStorage) prefixKey(key string) string {
	return fmt.Sprintf("%s:%s", s.prefix, key)
}

// Save stores any value in Redis, with special handling for structs
func (r *RedisStorage) Save(ctx context.Context, key string, value any) error {
	// Check if the value is nil
	if value == nil {
		return fmt.Errorf("cannot store nil value")
	}

	// Get the reflection value and type
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Special handling for structs
	if val.Kind() == reflect.Struct {
		data, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal struct: %v", err)
		}
		return r.client.Set(ctx, r.prefixKey(key), data, r.storageDeadline).Err()
	}

	// For non-struct types, store directly
	return r.client.Set(ctx, r.prefixKey(key), value, r.storageDeadline).Err()
}

// Get retrieves a value from Redis and automatically unmarshals if necessary
func (r *RedisStorage) Get(ctx context.Context, key string, dest any) error {
	// Check if destination is nil
	if dest == nil {
		return fmt.Errorf("destination cannot be nil")
	}

	// Get reflection value of destination
	val := reflect.ValueOf(dest)
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("destination must be a pointer")
	}
	val = val.Elem()

	// Get the value from Redis
	result, err := r.client.Get(ctx, r.prefixKey(key)).Bytes()
	if err != nil {
		return fmt.Errorf("failed to get value from Redis: %v", err)
	}

	// If destination is a struct, unmarshal the JSON
	if val.Kind() == reflect.Struct {
		return json.Unmarshal(result, dest)
	}

	// For non-struct types, decode directly
	return r.client.Get(ctx, r.prefixKey(key)).Scan(dest)
}

func (r *RedisStorage) Stop() {
	r.client.Close()
}
