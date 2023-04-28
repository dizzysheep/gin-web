package mysql

import (
	"context"
	"fmt"
	"gin-web/core/log"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

type GormLogger struct{}

func NewGormLogger() *GormLogger {
	return &GormLogger{}
}

func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	log.Get(GetGinContext(ctx)).WithFields(logrus.Fields{
		"context": ctx,
		"data":    utils.ToString(data),
	}).Info(msg)
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	log.Get(GetGinContext(ctx)).WithFields(logrus.Fields{
		"context": ctx,
		"data":    utils.ToString(data),
	}).Warn(msg)
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	log.Get(GetGinContext(ctx)).WithFields(logrus.Fields{
		"data": utils.ToString(data),
	}).Error(msg)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if err != nil {
		log.Get(GetGinContext(ctx)).Error(ctx, "SQL Trace Error", err)
		return
	}

	sql, rows := fc()
	log.Get(GetGinContext(ctx)).WithFields(logrus.Fields{
		"sql":     sql,
		"rows":    rows,
		"elapsed": fmt.Sprintf("%.3f", time.Since(begin).Seconds()),
	}).Debug("SQL Trace")
}
