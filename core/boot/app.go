package boot

import (
	"gin-web/core/redis"
	"gin-web/global"
	"github.com/yangxx0612/plugins/mysql"
)

func InitApp() {
	InitDB()
	redis.RedisSetup()
}

func InitDB() {
	global.BlogDB = mysql.GetMysql("blog")
}
