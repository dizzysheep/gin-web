package router

import (
	"testing"

	"gin-web/app/handler"
	"github.com/gin-gonic/gin"
)

// 路由注册冒烟测试：静态段(archive/slug)与参数段(:id)共存不应panic，
// handler为nil仅注册不调用
func TestRouteRegistration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	g := gin.New()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("route registration panic: %v", r)
		}
	}()

	UseIn(g, &handler.Handlers{})

	if len(g.Routes()) == 0 {
		t.Fatal("no routes registered")
	}
}
