package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gin-web/core/config"
	"gin-web/core/log"
	goredis "github.com/redis/go-redis/v9"
	"sync"
	"time"
)

var (
	RedisClient *goredis.Client
	redisOnce   sync.Once
)

func InitRedis() *goredis.Client {
	redisOnce.Do(func() {
		redisConfig := config.NewRedisConfig()
		RedisClient = goredis.NewClient(&goredis.Options{
			Addr:     redisConfig.Addr,
			Password: redisConfig.Password, // no password set
			DB:       redisConfig.DB,       // use default DB
		})

		if err := RedisClient.Ping(context.Background()).Err(); err != nil {
			log.Get(nil).Errorf("redis connection failed: %s", err.Error())
		}
	})
	return RedisClient
}

// SaveStruct 以json序列化存储结构体，expiration为过期时长
func SaveStruct(key string, data interface{}, expiration time.Duration) error {
	if RedisClient == nil {
		return fmt.Errorf("redis client is not initialized")
	}
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("save key `%s` json marshal error: %w", key, err)
	}
	if err := RedisClient.Set(context.Background(), key, string(jsonStr), expiration).Err(); err != nil {
		return fmt.Errorf("save key `%s` fail: %w", key, err)
	}
	return nil
}

// SaveStructNX stores a JSON value only when key does not already exist. It is
// used for cache-aside writes so concurrent misses do not continually replace
// one another's value.
func SaveStructNX(key string, data interface{}, expiration time.Duration) (bool, error) {
	if RedisClient == nil {
		return false, fmt.Errorf("redis client is not initialized")
	}
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return false, fmt.Errorf("save key `%s` json marshal error: %w", key, err)
	}
	ok, err := RedisClient.SetNX(context.Background(), key, string(jsonStr), expiration).Result()
	if err != nil {
		return false, fmt.Errorf("save key `%s` fail: %w", key, err)
	}
	return ok, nil
}

func GetStruct(key string, data interface{}) error {
	if RedisClient == nil {
		return fmt.Errorf("redis client is not initialized")
	}
	valStr, err := RedisClient.Get(context.Background(), key).Result()
	if errors.Is(err, goredis.Nil) {
		return fmt.Errorf("key `%s` not found", key)
	}
	if err != nil {
		return fmt.Errorf("get key `%s` fail: %w", key, err)
	}
	return json.Unmarshal([]byte(valStr), data)
}
