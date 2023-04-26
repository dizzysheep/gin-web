package log

import (
	"fmt"
	"gin-web/core/config"
	"gin-web/core/ginc"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net"
)

// Logger logger
type Logger = *logrus.Entry

var Loger Logger

var ServerIP = localIP()

func init() {
	Loger = logrus.WithFields(logrus.Fields{
		"app_name":    config.AppName,
		"server_host": config.AppAddr,
		"env":         config.Env,
		"server_ip":   ServerIP,
	})

	Loger.Logger.SetReportCaller(true)
}

func init() {
	setLogConfig()
	NewFileHook()
}

func setLogConfig() {
	logrus.SetFormatter(&logrus.TextFormatter{})
	if config.LogFormatter == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	//设置日志输出级别
	logLevel := getLogLevel(config.LogLevel)
	logrus.SetLevel(logLevel)
}

func getLogLevel(lvl string) logrus.Level {
	level, ok := logrus.ParseLevel(lvl)
	if ok == nil {
		return level
	}

	if config.IsDevEnv { // 开发环境设置为debug级别
		return logrus.DebugLevel
	}

	return logrus.ErrorLevel
}

func localIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// Get 获取日志实例
func Get(c *gin.Context) (log Logger) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("log panic,err:", err)
			log = Loger
		}
	}()

	//获取request信息
	if c == nil || c.Request == nil {
		return Loger
	}

	path := c.Request.URL.Path
	method := c.Request.Method
	raw := c.Request.URL.RawQuery

	//参数获取，兼容各种格式的请求
	request_param := ""
	if method == "GET" {
		request_param = raw
	} else if method == "POST" {
		if ct, ok := c.Request.Header["Content-Type"]; ok && ct[0] == "application/json" {
			raw, _ := c.GetRawData()
			request_param = string(raw)
		} else {
			//form格式
			_ = c.Request.ParseMultipartForm(128)
			data := c.Request.Form
			for k, v := range data {
				if len(v) == 1 {
					request_param += fmt.Sprintf("%s=%v&", k, v[0])
				} else {
					request_param += fmt.Sprintf("%s=%v&", k, v)
				}
			}
		}
	}

	//请求时间戳（毫秒数）
	var microTime int64 = 0
	microTimeInterface, exists := c.Get(StartMicroTime)
	if exists {
		microTime = microTimeInterface.(int64)
	}

	return logrus.WithFields(logrus.Fields{
		"is_gin":        true,
		"trace_id":      ginc.GetTraceID(c),
		"method":        method,
		"path":          fmt.Sprintf("【%s】%s", method, path),
		"request_param": request_param,
		"microtime":     microTime,
	})
}
