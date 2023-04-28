package mysql

import (
	"context"
	"gin-web/core/config"
	"gin-web/core/crypto"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	db, err = gorm.Open(mysql.Open(dbConf.Dsn), &gorm.Config{
		Logger: NewGormLogger(),
	})
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
