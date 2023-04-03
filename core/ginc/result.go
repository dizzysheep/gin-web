package ginc

import (
	"gin-demo/core/result"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Ok(c *gin.Context, data interface{}, errs ...interface{}) {
	res := result.Ok(data, errs...)
	res.RequestId = c.GetString("X-Request-ID")
	c.JSON(http.StatusOK, res)
}

func Fail(c *gin.Context, errs ...interface{}) {
	res := result.Fail(errs...)
	res.RequestId = c.GetString("X-Request-ID")
	c.JSON(http.StatusOK, res)
}
