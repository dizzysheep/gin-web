package skywalking

import (
	"context"
	"gin-web/core/config"
	"gin-web/core/log"
	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter"
	v3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
)

var Reporter go2sky.Reporter
var Tracer *go2sky.Tracer

const (
	ComponentIDHttp  = 2
	ComponentIDMySQL = 5
)

func InitSkyWalking() {
	var err error
	Reporter = nil
	Tracer = nil

	//skywalking开关
	if !config.SkywalkingSwitch {
		return
	}

	//采样率
	samplingRate := config.SkywalkingSamplingRate
	if samplingRate == 0 {
		samplingRate = 1
	}

	if config.SkywalkingHost == "" {
		return
	}

	Reporter, err = reporter.NewGRPCReporter(config.SkywalkingHost)
	if err != nil {
		log.Get(nil).Panic(err)
	}

	appName := "go::" + config.AppName
	Tracer, err = go2sky.NewTracer(appName, go2sky.WithReporter(Reporter), go2sky.WithSampler(samplingRate))
	if err != nil {
		log.Get(nil).Panic(err)
	}
}

func CreateHttpSwExitSpan(url, method string, ctx context.Context) (span go2sky.Span, headers [][]string) {
	if Reporter == nil || Tracer == nil {
		return nil, headers
	}
	span, _ = Tracer.CreateExitSpan(ctx, url, "Upstream_Service", func(headerKey, headerValue string) error {
		headers = append(headers, []string{headerKey, headerValue})
		return nil
	})
	span.SetComponent(ComponentIDHttp)
	span.SetSpanLayer(v3.SpanLayer_Http)
	span.Tag(go2sky.TagURL, url)
	span.Tag(go2sky.TagHTTPMethod, method)

	return span, headers
}
