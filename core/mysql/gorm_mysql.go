package mysql

import (
	"context"
	"gin-web/core/config"
	"gin-web/core/crypto"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"sync"
	"time"
)

var dbs = make(map[string]*gorm.DB)
var lock = sync.RWMutex{}

const InstanceIdKey = "gorm:instance_id"

func GetMysql(dbName string) *gorm.DB {
	var err error
	lock.RLock()
	db := dbs[dbName]
	lock.RUnlock()
	if db != nil {
		return db
	}

	dbConf := config.NewDbConfig(dbName)
	var configOption *gorm.Config
	if config.IsDevEnv {
		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // 输出目标
			logger.Config{
				SlowThreshold:             time.Second, // 慢查询阈值
				LogLevel:                  logger.Info, // 日志级别
				IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound错误
				Colorful:                  true,        // 开启彩色输出
			},
		)
		configOption = &gorm.Config{Logger: newLogger}
	} else {
		configOption = &gorm.Config{Logger: NewGormLogger()}
	}

	db, err = gorm.Open(mysql.Open(dbConf.Dsn), configOption)
	if err != nil {
		panic(err.Error())
	}

	//开启skywalkingSwitch
	if config.SkywalkingSwitch {
		InitGormHook(db)
	}
	db = db.Set(InstanceIdKey, crypto.Md5(dbConf.Dsn))
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(dbConf.MaxIdle)
	sqlDB.SetMaxOpenConns(dbConf.MaxActive)
	sqlDB.SetConnMaxLifetime(time.Duration(dbConf.IdleTimeout) * time.Second)

	lock.Lock()
	dbs[dbName] = db
	lock.Unlock()
	return db
}

func GetGinContext(ctx context.Context) *gin.Context {
	ginContext, ok := ctx.(*gin.Context)
	if !ok {
		return nil
	}
	return ginContext
}
