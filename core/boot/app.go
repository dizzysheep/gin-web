package boot

import (
	"gin-web/core/mysql"
	"gin-web/core/redis"
	"gin-web/global"
)

func InitApp() {
	InitDB()
	redis.RedisSetup()
}

func InitDB() {
	global.BlogDB = mysql.GetMysql("blog")
}
