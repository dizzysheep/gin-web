package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"gin-web/core/config"
	"github.com/go-redis/redis"
	"time"
)

var RedisClient *redis.Client

func InitRedis() *redis.Client {
	if RedisClient != nil {
		return RedisClient
	}
	redisConfig := config.NewRedisConfig()
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password, // no password set
		DB:       redisConfig.DB,       // use default DB
	})
	_, err := RedisClient.Ping().Result()
	if err != nil {
		fmt.Println("redis connection failed: ", err.Error())
	}
	return RedisClient
}

func SaveStruct(key string, data interface{}, second time.Duration) error {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return errors.New(fmt.Sprintf("save key `%s` json marshal error", key))
	}
	err = RedisClient.Set(key, string(jsonStr), time.Second*second).Err()
	if err != nil {
		return errors.New(fmt.Sprintf("save key `%s` fail:%s", key, err.Error()))
	}
	return nil
}

func GetStruct(key string, data interface{}) error {
	valStr, err := RedisClient.Get(key).Result()
	if redis.Nil == err {
		return errors.New(fmt.Sprintf("key `%s` not found", key))
	}
	if err != nil {
		return errors.New(fmt.Sprintf("get key `%s` fail:%s", key, err.Error()))
	}
	return json.Unmarshal([]byte(valStr), &data)
}
