package mysql

import (
	"gin-web/core/config"
	"github.com/jinzhu/gorm"
	"sync"
	"time"
)

var dbs = make(map[string]*gorm.DB)
var lock = sync.RWMutex{}

func GetMysql(dbName string) *gorm.DB {
	lock.RLock()
	db := dbs[dbName]
	lock.RUnlock()

	if db != nil {
		return db
	}

	dbConf := config.NewDbConfig(dbName)
	if len(dbConf.Type) == 0 { // 默认为mysql
		dbConf.Type = "mysql"
	}

	db, err := gorm.Open(dbConf.Type, dbConf.Dsn)
	if err != nil {
		panic(err.Error())
	}
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return dbConf.TablePrefix + defaultTableName
	}
	db.LogMode(config.DbMode)
	db.DB().SetMaxIdleConns(dbConf.MaxIdle)
	db.DB().SetMaxOpenConns(dbConf.MaxActive)
	db.DB().SetConnMaxLifetime(time.Duration(dbConf.IdleTimeout) * time.Second)

	lock.Lock()
	dbs[dbName] = db
	lock.Unlock()

	return db
}

func CloseMysql() {
	for k, db := range dbs {
		db.Close()
		delete(dbs, k)
	}
}
