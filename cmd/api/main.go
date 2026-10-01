package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gin-web/app/middleware"
	"gin-web/app/router"
	"gin-web/inject/api"
	"gin-web/internal/job"
	"gin-web/pkg/boostrap"
)

// @title gin-web 个人博客 API
// @version 1.0
// @description 个人博客后端 RESTful API（设计文档见 doc/blog-backend-design.md）
// @BasePath /api
func main() {
	boostrap.InitApp()

	ginApp := boostrap.NewGin()
	appContainer := api.NewAppContainer()

	//启动服务
	middleware.UserIn(ginApp)
	router.UseIn(ginApp, appContainer.Handlers)
	service := boostrap.SetupServer(ginApp)

	// 后台任务：浏览量刷盘
	jobCtx, jobCancel := context.WithCancel(context.Background())
	defer jobCancel()
	go job.RunViewCountFlusher(jobCtx)

	// 监听信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	//退出服务
	jobCancel()
	boostrap.ShutdownServer(service, 30*time.Second)
}
