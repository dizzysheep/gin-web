package result

import (
	"fmt"
)

const (
	ResultFail = 0
	ResultOk   = 1
)

const (
	ErrorMsgSuccess         = "OK"
	ErrorMsgFail            = "ERROR"
	ErrorMsgInvalidArgument = "InvalidArgument"
	ErrorMsgNotLogin        = "Unauthenticated"
)

type M map[string]interface{}

type Result struct {
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	Data      interface{} `json:"data"`
	RequestId string      `json:"request_id"`
}

type PageInfo struct {
	Total int         `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
	List  interface{} `json:"list"`
}

// 1:data 2:error_msg 3:code
func Ok(data interface{}, errs ...interface{}) *Result {
	ret := &Result{
		Code: ResultOk,
		Msg:  ErrorMsgSuccess,
	}
	ret.Data = data
	ret.dealError(errs)
	return ret
}

// 1:PageInfo 2:error_msg 3:code
func Page(page *PageInfo, errs ...interface{}) *Result {
	ret := &Result{
		Code: ResultOk,
		Msg:  ErrorMsgSuccess,
	}
	ret.Data = page
	ret.dealError(errs)
	return ret
}

// 1:error_msg 2:code
func FailData(data interface{}, errs ...interface{}) *Result {
	ret := &Result{
		Code: ResultFail,
		Msg:  ErrorMsgFail,
	}
	ret.Data = data
	ret.dealError(errs)
	return ret
}

// 1:error_msg 2:code
func Fail(errs ...interface{}) *Result {
	ret := &Result{
		Code: ResultFail,
		Msg:  ErrorMsgFail,
	}
	ret.dealError(errs)
	return ret
}

func InvalidArgumentError() *Result {
	ret := &Result{
		Code: ResultFail,
		Msg:  ErrorMsgInvalidArgument,
	}

	return ret
}

func NotLoginError() *Result {
	ret := &Result{
		Code: ResultFail,
		Msg:  ErrorMsgNotLogin,
	}

	return ret
}

func (s *Result) String() string {
	return fmt.Sprintf("Code:[%d] Msg:[%s]", s.Code, s.Msg)
}

func (s *Result) SetRequestId(requestId string) *Result {
	s.RequestId = requestId
	return s
}

func (s *Result) dealError(errs []interface{}) {
	for i, err := range errs {
		if val, ok := err.(TransferMsg); ok { // error_msg
			s.Msg = val.Transform()
			continue
		}

		if i == 0 { // error_msg
			if val, ok := err.(string); ok { // error_msg
				s.Msg = val
			}
			continue
		}
	}
}

type TransferMsg interface {
	Transform() string
}
