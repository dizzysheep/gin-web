package main

import (
	"gin-web/app/middleware"
	"gin-web/app/router"
	"gin-web/inject/api"
	"gin-web/pkg/boostrap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	boostrap.InitApp()

	ginApp := boostrap.NewGin()
	appContainer := api.NewAppContainer()

	//启动服务
	middleware.UserIn(ginApp)
	router.UseIn(ginApp, appContainer.Handlers)
	service := boostrap.SetupServer(ginApp)

	// 监听信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	//退出服务
	boostrap.ShutdownServer(service, 30*time.Second)
}
