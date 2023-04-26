package main

import (
	"context"
	"gin-web/app/routers"
	"gin-web/core/boot"
	"gin-web/core/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	boot.InitApp()

	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:    config.AppAddr,
		Handler: routers.InitRouter(),
	}
	// 异步启动 HTTP 服务器
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
