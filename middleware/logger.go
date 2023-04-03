package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/yangxx0612/plugins/config"
	"github.com/yangxx0612/plugins/log"
	"os"
	"time"
)

//自定义一个结构体，实现 gin.ResponseWriter interface
type responseWriter struct {
	gin.ResponseWriter
	b *bytes.Buffer
}

//重写 Write([]byte) (int, error) 方法
func (w responseWriter) Write(b []byte) (int, error) {
	//向一个bytes.buffer中写一份数据来为获取body使用
	w.b.Write(b)
	//完成gin.Context.Writer.Write()原有功能
	return w.ResponseWriter.Write(b)
}

func Logger() gin.HandlerFunc {

	return func(c *gin.Context) {
		// 一.配置所需的 Fields
		startTime := time.Now()

		//相应结果
		response := responseWriter{
			c.Writer,
			bytes.NewBuffer([]byte{}),
		}
		c.Writer = response

		c.Next()

		// 1.API 调用耗时
		spendTime := time.Since(startTime).Milliseconds()
		ST := fmt.Sprintf("%d ms", spendTime)

		// 2.主机名
		hostName, err := os.Hostname()
		if err != nil {
			hostName = "unknown"
		}

		statusCode := c.Writer.Status() // 3.状态码

		// 二.从标准记录器创建一个条目，并向其中添加多个字段(隐式添加 log 本身的时间戳,信息等 fields )
		entry := log.Loger.WithFields(logrus.Fields{
			"request_id": c.GetString("X-Request-ID"),
			"app_ame":    config.AppName,
			"host_name":  hostName,
			"domain":     c.Request.Host,
			"url":        c.Request.RequestURI,
			"method":     c.Request.Method,
			"referer":    c.Request.Referer(),
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
			"response":   response.b.String(),
			"status":     statusCode,
			"spend_time": ST,
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
