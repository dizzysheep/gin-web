package config

import (
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func NewRedisConfig() *RedisConfig {
	conf := viper.GetStringMapString("redis")
	return &RedisConfig{
		Addr:     conf["addr"],
		Password: conf["password"],
		DB:       cast.ToInt(conf["db"]),
	}
}
