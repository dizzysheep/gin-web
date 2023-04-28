package service

import (
	"gin-web/global"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetBlogDB(ctx *gin.Context) *gorm.DB {
	return global.BlogDB.WithContext(ctx)
}
