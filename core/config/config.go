package config

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	_defaultConfigName = "app"
	_defaultConfigType = "toml"
)

var (
	AppName         = "default"
	AppAddr         = "127.0.0.1:8080"
	AppPath         = ""
	Hostname        = "localhost"
	AppReadTimeout  = 10
	AppWriteTimeout = 10

	JwtSecret = ""

	LogLevel     = "info"
	LogFormatter = "text"
	LogTopic     = "golang_log"
	DbMode       = false

	// skywalking服务器配置
	SkywalkingHost                 = ""
	SkywalkingSwitch               = false
	SkywalkingSamplingRate float64 = 0

	RunMode = "debug"

	// IsDevEnv 开发环境标志
	IsDevEnv = false
	// IsTestEnv 测试环境标志
	IsTestEnv = false
	// IsProdEnv 生产环境标志
	IsProdEnv = false
	// Env 运行环境
	Env = "dev"
)

func init() {
	ReadConfig()
	LoadApp()
}

// ReadConfig 读取配置文件
func ReadConfig() {
	viper.SetConfigType(_defaultConfigType)
	viper.SetConfigName(GetConfigName())
	viper.AddConfigPath(GetConfigPath())
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	viper.AutomaticEnv()
}

// GetConfigName 获取配置文件名称
func GetConfigName() string {
	return _defaultConfigName + "." + _defaultConfigType
}

// GetConfigPath 读取配置文件路径
func GetConfigPath() string {
	configPath := os.Getenv("CONF_PATH")
	if configPath == "" {
		var err error
		if AppPath, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
			panic(err)
		}
		workPath, err := os.Getwd()
		configPath = filepath.Join(workPath, "config") // 先找当前目录conf下面的文件
	}
	return configPath
}

// LoadApp 加载app运行配置
func LoadApp() {
	Hostname, _ = os.Hostname()
	AppName = viper.GetString("app.appName")
	appAddr := viper.GetString("app.appAddr")
	JwtSecret = viper.GetString("app.jwtSecret")
	RunMode = viper.GetString("app.runMode")
	Env = viper.GetString("app.env")
	if appAddr != "" {
		AppAddr = appAddr
	}
	appReadTimeout := viper.GetInt("app.readTimeout")
	if appReadTimeout != 0 {
		AppReadTimeout = appReadTimeout
	}
	appWriteTimeout := viper.GetInt("app.writeTimeout")
	if appWriteTimeout != 0 {
		AppWriteTimeout = appWriteTimeout
	}
	DbMode = viper.GetBool("log.dbMode")
	LogFormatter = viper.GetString("log.logFormatter")
	LogLevel = viper.GetString("log.LogLevel")
	if Env == "dev" {
		IsDevEnv = true
	}
	// skywalking服务器配置
	SkywalkingSwitch = GetBool("skywalking.SwitchOn")
	SkywalkingHost = GetString("skywalking.Host")
	SkywalkingSamplingRate = GetFloat64("skywalking.SamplingRate")
}

// GetFloat64 获取浮点数配置
func GetFloat64(key string) float64 {
	return viper.GetFloat64(key)
}

// GetString 获取字符串配置
func GetString(key string) string {
	return viper.GetString(key)
}

// GetInt 获取整数配置
func GetInt(key string) int {
	return viper.GetInt(key)
}

// GetInt32 获取 int32 配置
func GetInt32(key string) int32 {
	return viper.GetInt32(key)
}

// GetInt64 获取 int64 配置
func GetInt64(key string) int64 {
	return viper.GetInt64(key)
}

// GetDuration 获取时间配置
func GetDuration(key string) time.Duration {
	return viper.GetDuration(key)
}

// GetBool 获取配置布尔配置
func GetBool(key string) bool {
	return viper.GetBool(key)
}

// GetStringSlice 获取配置字符串切片配置
func GetStringSlice(key string) []string {
	var value []string
	strs := GetString(key)
	// 未填写，直接返回
	if len(strs) == 0 {
		return value
	}
	// 解析json失败， apollo配置的是json切片， app.toml配置的是逗号连接的字符串
	if err := json.Unmarshal([]byte(strs), &value); err != nil {
		value = strings.Split(strs, ",")
	}
	return value
}

// GetStringMap 获取配置map配置
func GetStringMap(key string) map[string]interface{} {
	var value map[string]interface{}
	_ = json.Unmarshal([]byte(GetString(key)), &value)
	return value
}

// GetStringMapString 获取配置字符串map配置
func GetStringMapString(key string) map[string]string {
	value := GetStringMap(key)
	return cast.ToStringMapString(value)
}

// IsSet 配置是否存在
func IsSet(key string) bool {
	value := GetString(key)
	if len(value) <= 0 {
		return false
	}
	return true
}
