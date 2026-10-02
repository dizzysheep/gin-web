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
	Hostname        = "localhost"
	AppReadTimeout  = 10
	AppWriteTimeout = 10

	JwtSecret = ""

	LogLevel     = "info"
	LogFormatter = "text"
	LogTopic     = "golang_log"
	DbMode       = false

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

// GetConfigName 获取配置文件名称（不含扩展名）
func GetConfigName() string {
	return _defaultConfigName
}

// GetConfigPath 读取配置文件路径
func GetConfigPath() string {
	if configPath := os.Getenv("CONF_PATH"); configPath != "" {
		return configPath
	}

	// 依次在当前目录及其上级目录中查找 config/app.toml，
	// 便于在子目录中运行程序或单元测试时也能定位到配置文件
	workPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	for dir := workPath; dir != filepath.Dir(dir); dir = filepath.Dir(dir) {
		candidate := filepath.Join(dir, "config")
		if _, err := os.Stat(filepath.Join(candidate, GetConfigName()+"."+_defaultConfigType)); err == nil {
			return candidate
		}
	}
	return filepath.Join(workPath, "config")
}

// LoadApp 加载app运行配置
func LoadApp() {
	Hostname, _ = os.Hostname()
	AppName = viper.GetString("app.appName")
	appAddr := viper.GetString("app.appAddr")
	JwtSecret = viper.GetString("jwt.secret")
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
