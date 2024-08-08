package response

type ErrCode int

func (e ErrCode) Code() int {
	return int(e)
}

type CustomError struct {
	Code    ErrCode `json:"code"`
	Message string  `json:"message"`
}

func (c *CustomError) Error() string {
	return c.Code.String()
}

//go:generate stringer -type ErrCode -linecomment

const (
	SuccessOk ErrCode = 1 //ok

	//ErrFail 200 -- 通用错误码
	ErrFail          ErrCode = 200000 //服务器错误
	ErrInvalidParams ErrCode = 200001 // 非法请求参数
	ErrNoAccess      ErrCode = 200002 // 无访问权限
	ErrNoFound       ErrCode = 200003 // 找不到资源
	ErrDb            ErrCode = 200004 // 数据库出错
	ErrCache         ErrCode = 200005 // 缓存出错
	ErrLimitExceeded ErrCode = 200006 // 请求国泰，请稍后重试
)
