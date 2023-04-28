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

	return logrus.WithFields(logrus.Fields{
		"is_gin":   true,
		"trace_id": ginc.GetTraceID(c),
		"method":   c.Request.Method,
		"path":     c.Request.URL.Path,
	})
}
