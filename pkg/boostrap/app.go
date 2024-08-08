package boostrap

import (
	"context"
	"errors"
	"gin-web/core/config"
	"gin-web/core/log"
	"gin-web/core/mysql"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func InitApp() {
}

func InitDBEngine() *gorm.DB {
	return mysql.GetMysql("blog")
}

func NewGin() *gin.Engine {
	gin.SetMode(config.RunMode)
	g := gin.New()
	g.ForwardedByClientIP = true
	return g
}

func SetupServer(handler *gin.Engine) *http.Server {
	// 创建 HTTP 服务器
	server := &http.Server{
		Addr:           config.AppAddr,
		Handler:        handler,
		ReadTimeout:    time.Duration(config.AppReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(config.AppWriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	// 异步启动 HTTP 服务器
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Get(nil).Fatalf("ListenAndServe err: %s\n", err)
		}
	}()

	return server
}

func ShutdownServer(server *http.Server, ctxTimeout time.Duration) {
	log.Get(nil).Println("Shutdown Server")

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Get(nil).Fatalf("Server Shutdown err: %s\n", err)
	}

	log.Get(nil).Println("Server exiting")
}
