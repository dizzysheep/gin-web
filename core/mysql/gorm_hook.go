package mysql

import (
	"fmt"
	"gin-web/core/log"
	"gin-web/core/skywalking"
	"github.com/SkyAPM/go2sky"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	v3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
)

const (
	callBackBeforeName = "gorm:before"
	callBackAfterName  = "gorm:after"
)

func before(db *gorm.DB) {
	ginCtx := GetGinContext(db.Statement.Context)
	if ginCtx == nil {
		return
	}
	span, _, err := skywalking.Tracer.CreateEntrySpan(ginCtx.Request.Context(), "gorm/mysql", func(key string) (string, error) {
		val, _ := db.Get(key)
		return cast.ToString(val), nil
	})
	if err != nil {
		log.Get(ginCtx).Printf("gorm create entry span error %v \n", err)
		return
	}
	val, _ := db.Get(InstanceIdKey)
	span.SetComponent(skywalking.ComponentIDMySQL)
	span.SetSpanLayer(v3.SpanLayer_Database)
	span.Tag(go2sky.TagDBType, "mysql")
	span.Tag(go2sky.TagDBInstance, cast.ToString(val))
	db.Set("SkyWalkingSpan", span)
}

func after(db *gorm.DB) {
	span, ok := db.Get("SkyWalkingSpan")
	ctx := GetGinContext(db.Statement.Context)
	if !ok {
		log.Get(ctx).Printf("gorm get span from db error")
		return
	}
	curSpan, ok := span.(go2sky.Span)
	if !ok {
		log.Get(ctx).Printf("gorm transfer go2sky.Span error")
		return
	}
	curSpan.Tag(go2sky.TagDBStatement, db.Statement.SQL.String())
	curSpan.Tag(go2sky.TagDBSqlParameters, fmt.Sprintf("%+v", db.Statement.Vars))
	curSpan.End()
}

func InitGormHook(db *gorm.DB) {
	// 开始前 - 并不是都用相同的方法，可以自己自定义
	_ = db.Callback().Create().Before("gorm:before_create").Register(callBackBeforeName, before)
	_ = db.Callback().Query().Before("gorm:query").Register(callBackBeforeName, before)
	_ = db.Callback().Delete().Before("gorm:before_delete").Register(callBackBeforeName, before)
	_ = db.Callback().Update().Before("gorm:setup_reflect_value").Register(callBackBeforeName, before)
	_ = db.Callback().Row().Before("gorm:row").Register(callBackBeforeName, before)
	_ = db.Callback().Raw().Before("gorm:raw").Register(callBackBeforeName, before)

	// 结束后 - 并不是都用相同的方法，可以自己自定义
	_ = db.Callback().Create().After("gorm:after_create").Register(callBackAfterName, after)
	_ = db.Callback().Query().After("gorm:after_query").Register(callBackAfterName, after)
	_ = db.Callback().Delete().After("gorm:after_delete").Register(callBackAfterName, after)
	_ = db.Callback().Update().After("gorm:after_update").Register(callBackAfterName, after)
	_ = db.Callback().Row().After("gorm:row").Register(callBackAfterName, after)
	_ = db.Callback().Raw().After("gorm:raw").Register(callBackAfterName, after)
}
