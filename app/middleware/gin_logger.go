package middleware

import (
	"bytes"
	"fmt"
	"gin-web/core/config"
	"gin-web/core/log"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"time"
)

// 自定义一个结构体，实现 gin.ResponseWriter interface
type responseWriter struct {
	gin.ResponseWriter
	b *bytes.Buffer
}

// 重写 Write([]byte) (int, error) 方法
func (w responseWriter) Write(b []byte) (int, error) {
	//向一个bytes.buffer中写一份数据来为获取body使用
	w.b.Write(b)
	//完成gin.Context.Writer.Write()原有功能
	return w.ResponseWriter.Write(b)
}

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		//屏蔽调健康监测日志
		if c.Request.RequestURI == "/api/healthz" || config.IsDevEnv {
			return
		}

		startTime := time.Now()
		response := responseWriter{c.Writer, bytes.NewBuffer([]byte{})}
		c.Writer = response

		c.Next()

		ST := fmt.Sprintf("%d ms", time.Since(startTime).Milliseconds())
		statusCode := c.Writer.Status()
		responseStr := response.b.String()
		responseStr = log.HideSensitiveInfo(response.b.String())

		// 二.从标准记录器创建一个条目，并向其中添加多个字段(隐式添加 log 本身的时间戳,信息等 fields )
		entry := log.Get(c).WithFields(logrus.Fields{
			"status":        statusCode,
			"user_agent":    c.Request.UserAgent(),
			"response":      responseStr,
			"request_body":  getRequestBody(c),
			"response_time": ST,
		})

		// Errors 保存了使用当前context的所有中间件/handler 所产生的全部错误信息。
		// 源码注释： Errors is a list of errors attached to all the handlers/middlewares who used this context.
		// 三.将系统内部的错误 log 出去
		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		}

		// 四.根据状态码决定打印 log 的等级
		if statusCode >= 500 {
			entry.Error()
		} else if statusCode >= 400 {
			entry.Warn()
		} else {
			entry.Info()
		}
	}
}

// getRequestBody 缓存获取的http请求参数，错误日志中记录请求信息
// 防并发协程调用竟态问题，在 GinLogger() 初始化一次，避免协程并发下多次同时初始化缓存
func getRequestBody(c *gin.Context) string {
	ginCacheKey := "logrusRequestParam"
	rebody := ""

	//如果已在 gin.Context 已缓存数据，直接获取使用
	if cb, ok := c.Get(ginCacheKey); ok {
		if rebody, ok := cb.(string); ok {
			return rebody
		}
	}

	switch c.Request.Method {
	case "GET":
		return c.Request.URL.RawQuery
	case "POST":
		// 读取 c.Request.Body 的数据，
		if raw, err := c.GetRawData(); err == nil && len(raw) > 0 {
			//重写  c.Request.Body 的数据，否则在 handler 下参数为空了
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(raw))
			rebody = string(raw)
		} else {
			//form格式
			_ = c.Request.ParseMultipartForm(128)
			for k, v := range c.Request.Form {
				if len(v) == 1 {
					rebody += fmt.Sprintf("%s=%v&", k, v[0])
				} else {
					rebody += fmt.Sprintf("%s=%v&", k, v)
				}
			}
		}
	}

	// 缓存用于多次调用
	c.Set(ginCacheKey, rebody)
	return rebody
}
