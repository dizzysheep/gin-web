package skywalking

import (
	"context"
	"gin-web/core/config"
	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter"
)

var Reporter go2sky.Reporter
var Tracer *go2sky.Tracer
var SamplingRate float64
var Ctx context.Context
var ctxKey = "CONTEXT4SW"

//上下文对象，兼容gin框架上下文gin.Context
type swContext interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
	GetHeader(headerKey string) string
}

//skywalking连接初始化
//serviceName 服务名    hostName 主机名
func init() {
	Reporter = nil
	Tracer = nil
	Ctx = context.Background()

	//skywalking开关
	switchOn := config.SkywalkingSwitch
	if !switchOn {
		return
	}

	//采样率
	SamplingRate = config.SkywalkingSamplingRate
	if SamplingRate == 0 {
		SamplingRate = 1
	}

	host := config.SkywalkingHost
	if host == "" {
		return
	}

	r, err := reporter.NewGRPCReporter(host)
	if err != nil {
		return
	}
	Reporter = r

	//hostName,_ := os.Hostname()
	appName := "go::" + config.AppName
	tracer, err := go2sky.NewTracer(appName, go2sky.WithReporter(r), go2sky.WithSampler(SamplingRate))
	if err != nil {
		return
	}
	Tracer = tracer
}

//往context中设置skywalking上下文
func SetContext(c swContext, ctx context.Context) {
	c.Set(ctxKey, ctx)
}

//从context中获取设置的skywalking上下文
func GetContext(c swContext) context.Context {
	ctx, ok := c.Get(ctxKey)
	if ok {
		return ctx.(context.Context)
	} else {
		return context.Background()
	}
}

//创建出口Span，并返回追踪标识Header头信息
//operationName 操作名，一般为正在请求的接口名
func CreateSwExitSpan(operationName string, ctx context.Context) (span go2sky.Span, headers [][]string) {
	if Reporter == nil || Tracer == nil {
		return nil, headers
	}
	span, _ = Tracer.CreateExitSpan(ctx, operationName, "Upstream_Service", func(headerKey, headerValue string) error {
		headers = append(headers, []string{headerKey, headerValue})
		return nil
	})
	return span, headers
}

//创建入口Span,从Header读取追踪标识
//operationName 操作名，一般为当前被请求接口名
//c gin框架的request上下文
func CreateEntrySpan(operationName string, c swContext) (span go2sky.Span, ctx context.Context) {
	if Reporter == nil || Tracer == nil {
		return nil, ctx
	}
	span, ctx, _ = Tracer.CreateEntrySpan(context.Background(), operationName, func(headerKey string) (string, error) {
		return c.GetHeader(headerKey), nil
	})
	return span, ctx
}
