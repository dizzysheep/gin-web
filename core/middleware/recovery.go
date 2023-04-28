package middleware

import (
	"gin-web/core/ginc"
	"gin-web/core/log"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"runtime"
	"runtime/debug"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				_, file, line, _ := runtime.Caller(2)
				stackStr := string(debug.Stack())
				log.Get(c).WithFields(logrus.Fields{
					"error_msg":   err,
					"error_file":  file,
					"error_line":  line,
					"error_stack": stackStr,
				}).Errorf("[Recovery] msg:%+v", err)
				ginc.InternalServerError(c)
			}
		}()
		c.Next()
	}
}
