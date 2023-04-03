package main

import (
	"gin-demo/core/redis"
	"gin-demo/models"
	"gin-demo/routers"
	"github.com/fvbock/endless"
	"github.com/yangxx0612/plugins/config"
	"github.com/yangxx0612/plugins/log"
	"syscall"
	"time"
)

func init() {
	models.Setup()
	redis.RedisSetup()
}

func main() {
	endless.DefaultReadTimeOut = time.Duration(config.AppReadTimeout) * time.Second
	endless.DefaultWriteTimeOut = time.Duration(config.AppWriteTimeout) * time.Second
	endless.DefaultMaxHeaderBytes = 1 << 20
	server := endless.NewServer(config.AppAddr, routers.InitRouter())
	server.BeforeBegin = func(add string) {
		log.Loger.Printf("Actual pid is %d", syscall.Getpid())
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Loger.Printf("Actual pid is %d", syscall.Getpid())
	}
}
