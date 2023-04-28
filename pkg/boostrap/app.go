package boostrap

import (
	"gin-web/core/mysql"
	"gin-web/core/redis"
	"gin-web/global"
)

func InitApp() {
	//数据库连接
	initDBEngine()
	redis.InitRedis()
	//skywalking.InitSkyWalking()
}

func initDBEngine() {
	global.BlogDB = mysql.GetMysql("blog")
}
