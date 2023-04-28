package log

import (
	"fmt"
	"gin-web/core/config"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	kl "github.com/tracer0tong/kafkalogrus"
	"os"
	"path/filepath"
	"time"
)

func NewKafkaHook(id string, formatter logrus.Formatter, brokers []string) error {
	kafkaHook, err := kl.NewKafkaLogrusHook(id, logrus.AllLevels, formatter, brokers, "KafkaInput", true, nil)
	if err != nil {
		logrus.Errorf("config es logger error. %+v", errors.WithStack(err))
		return err
	}

	logrus.AddHook(kafkaHook)
	return nil
}

func NewFileHook() {
	// 设置输出文件
	workPath, err := os.Getwd()
	if err != nil {
		fmt.Println("err:", err)
	}
	filePath := filepath.Join(workPath, config.GetString("log.filePath"))

	// 设置日志切割 rotatelogs
	writer, _ := rotatelogs.New(
		filePath+"/%Y-%m-%d.log",
		//日志最大保存时间
		rotatelogs.WithMaxAge(7*24*time.Hour),
		//设置日志切割时间间隔(1天)(隔多久分割一次)
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	// lfshook 决定哪些日志级别可用日志分割
	writeMap := lfshook.WriterMap{
		logrus.PanicLevel: writer,
		logrus.FatalLevel: writer,
		logrus.ErrorLevel: writer,
		logrus.WarnLevel:  writer,
		logrus.InfoLevel:  writer,
		logrus.DebugLevel: writer,
	}

	// 配置 lfshook
	hook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "time",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
		TimestampFormat: "2006.01.02 - 15:04:05",
	})

	//为 logrus 实例添加自定义 hook
	logrus.AddHook(hook)
}
