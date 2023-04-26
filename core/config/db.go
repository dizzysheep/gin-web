package config

import (
	"github.com/spf13/viper"
	"strconv"
)

type Database struct {
	Type        string
	Dsn         string
	MaxIdle     int
	MaxActive   int
	IdleTimeout int
	TablePrefix string
}

func NewDbConfig(key string) *Database {
	conf := viper.GetStringMapString("db_" + key)
	maxIdle, _ := strconv.Atoi(conf["maxIdle"])
	maxActive, _ := strconv.Atoi(conf["maxActive"])
	idleTimeout, _ := strconv.Atoi(conf["idleTimeout"])

	return &Database{
		Type:        conf["type"],
		Dsn:         conf["dsn"],
		TablePrefix: conf["tablePrefix"],
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: idleTimeout,
	}
}
