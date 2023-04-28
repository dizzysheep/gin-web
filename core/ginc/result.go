package ginc

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	// 通用错误
	SuccessCode = 1
	FailCode    = 0
)

type Response struct {
	Code      int         `json:"code"`
	Data      interface{} `json:"data"`
	Msg       string      `json:"msg"`
	RequestId string      `json:"request_id"`
}

type PageInfo struct {
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
	List  interface{} `json:"list"`
}

// BuildResponse 构建返回格式
func BuildResponse(code int, msg string, data interface{}) *Response {
	response := &Response{
		Code: code,
		Data: data,
		Msg:  msg,
	}
	return response
}

// SetRequestId 添加requestId
func (resp *Response) SetRequestId(ctx *gin.Context) *Response {
	if ctx == nil {
		return resp
	}
	resp.RequestId = GetTraceID(ctx)
	return resp
}

func (resp *Response) toJson(c *gin.Context) {
	c.JSON(http.StatusOK, resp)
}

func Ok(c *gin.Context, data interface{}) {
	res := BuildResponse(SuccessCode, "ok", data).SetRequestId(c)
	c.JSON(http.StatusOK, res)
}

func Fail(c *gin.Context, msg string) {
	res := BuildResponse(FailCode, msg, nil).SetRequestId(c)
	c.JSON(http.StatusOK, res)
}

func FailData(c *gin.Context, msg string, data interface{}) {
	res := BuildResponse(FailCode, msg, data).SetRequestId(c)
	c.JSON(http.StatusOK, res)
}

func InternalServerError(c *gin.Context) {
	res := BuildResponse(FailCode, "internal server error", nil).SetRequestId(c)
	c.JSON(http.StatusInternalServerError, res)
}
