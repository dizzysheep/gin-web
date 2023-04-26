// 操作请求 gin.Context 信息
package ginc

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// TraceIDKey 请求唯一标识，类型：string
	TraceIDKey = "x-trace-id"
	// StartTimeKey 请求开始时间，类型：time.Time
	StartTimeKey = "x-start-time"
	// UserIDKey 用户 ID，未登录则为 0，类型：int64
	UserIDKey = "x-user-id"
	// UserIPKey 用户 IP，类型：string
	UserIPKey = "x-user-ip"
	// PlatformKey 用户使用平台，ios, android, pc
	PlatformKey = "x-platform"
	// BuildKey 客户端构建版本号
	BuildKey = "x-build"
	// VersionKey 客户端版本号
	VersionKey = "x-version"
	// AccessKeyKey 移动端支付令牌
	AccessKeyKey = "x-access-key"
	// DeviceKey 移动 app 设备标识，android, phone, pad
	DeviceKey = "x-device"
	// MobiAppKey 移动 app 标识，android, phone, pad
	MobiAppKey = "x-mobi-app"
	// UserPortKey 用户端口
	UserPortKey = "x-user-port"
	// ManageUserKey 管理后台用户名
	ManageUserKey = "x-manage-user"
	// BuvidKey 非登录用户标识
	BuvidKey = "x-buvid"
)

// GetUserID 获取当前登录用户 ID
func GetUserID(ctx *gin.Context) int64 {
	uid, _ := ctx.Value(UserIDKey).(int64)
	return uid
}

// GetUserIP 获取用户 IP
func GetUserIP(ctx *gin.Context) string {
	ip, _ := ctx.Value(UserIPKey).(string)
	return ip
}

// GetUserPort 获取用户端口
func GetUserPort(ctx *gin.Context) string {
	port, _ := ctx.Value(UserPortKey).(string)
	return port
}

// GetPlatform 获取用户平台
func GetPlatform(ctx *gin.Context) string {
	platform, _ := ctx.Value(PlatformKey).(string)
	return platform
}

// IsIOSPlatform 判断是否为 IOS 平台
func IsIOSPlatform(ctx *gin.Context) bool {
	return GetPlatform(ctx) == "ios"
}

// GetTraceID 获取用户请求标识
func GetTraceID(ctx *gin.Context) string {
	id, _ := ctx.Value(TraceIDKey).(string)
	return id
}

// GetBuild 获取客户端构建版本号
func GetBuild(ctx *gin.Context) string {
	build, _ := ctx.Value(BuildKey).(string)
	return build
}

// GetDevice 获取用户设备，配合 GetPlatform 使用
func GetDevice(ctx *gin.Context) string {
	device, _ := ctx.Value(DeviceKey).(string)
	return device
}

// GetMobiApp 获取 APP 标识
func GetMobiApp(ctx *gin.Context) string {
	app, _ := ctx.Value(MobiAppKey).(string)
	return app
}

// GetVersion 获取客户端版本
func GetVersion(ctx *gin.Context) string {
	version, _ := ctx.Value(VersionKey).(string)
	return version
}

// GetAccessKey 获取客户端版本
func GetAccessKey(ctx *gin.Context) string {
	key, _ := ctx.Value(AccessKeyKey).(string)
	return key
}

// WithTraceID 注入 trace_id
func WithTraceID(ctx *gin.Context) string {
	traceID := NewTraceID()

	ctx.Writer.Header().Set(TraceIDKey, traceID)
	ctx.Set(TraceIDKey, traceID)
	return traceID
}

// WithStartTime 注入 trace_id
func WithStartTime(ctx *gin.Context, startTime time.Time) {
	ctx.Set(StartTimeKey, startTime)
	return
}

// GetManageUser 获取管理后台用户名
func GetManageUser(ctx *gin.Context) string {
	user, _ := ctx.Value(ManageUserKey).(string)
	return user
}

// GetBuvid 获取用户 buvid
func GetBuvid(ctx *gin.Context) string {
	buvid, _ := ctx.Value(BuvidKey).(string)
	return buvid
}

// NewTraceID 生成
func NewTraceID() string {
	return uuid.New().String()
}
