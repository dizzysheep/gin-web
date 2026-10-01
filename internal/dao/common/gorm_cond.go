package common

import (
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"reflect"
)

type GormCond interface {
	BuildCond(query *gorm.DB) *gorm.DB
}

type EqCond struct {
	Field string
	Value interface{}
}

func (c *EqCond) BuildCond(q *gorm.DB) *gorm.DB {
	return q.Where(c.Field, c.Value)
}

type LikeCond struct {
	Field string
	Value string
}

func (c *LikeCond) BuildCond(query *gorm.DB) *gorm.DB {
	if c.Value != "" {
		return query.Where(clause.Like{
			Column: clause.Column{Name: c.Field},
			Value:  "%" + c.Value + "%",
		})
	}

	return query
}

type OrConditions struct {
	GormCond []GormCond
}

func (c *OrConditions) BuildCond(query *gorm.DB) *gorm.DB {
	if len(c.GormCond) == 0 {
		return query
	}

	// 使用 OR 组合子条件
	queryClauses := make([]clause.Expression, 0, len(c.GormCond))
	for _, cond := range c.GormCond {
		subQuery := query.Session(&gorm.Session{DryRun: true}) // 使用 DryRun 模式避免实际查询
		subQuery = cond.BuildCond(subQuery)
		queryClauses = append(queryClauses, subQuery.Statement.Clauses["WHERE"].Expression)
	}

	query = query.Where(clause.Or(queryClauses...))
	return query
}

type GormConditions []GormCond

func (c GormConditions) BuildConditions(query *gorm.DB) *gorm.DB {
	for _, cond := range c {
		query = cond.BuildCond(query)
	}
	return query
}

// NotDeletedCond 软删除过滤：deleted_on = 0
type NotDeletedCond struct{}

func (c *NotDeletedCond) BuildCond(query *gorm.DB) *gorm.DB {
	return query.Where("deleted_on = 0")
}

// NotDeleted 在业务条件前追加软删除过滤
func NotDeleted(conditions GormConditions) GormConditions {
	return append(GormConditions{&NotDeletedCond{}}, conditions...)
}

// OrderCond 排序条件，Columns 必须为代码内白名单字段，不可透传用户输入
type OrderCond struct {
	Columns []string
}

func (c *OrderCond) BuildCond(query *gorm.DB) *gorm.DB {
	if len(c.Columns) == 0 {
		return query
	}
	return query.Order(strings.Join(c.Columns, ", "))
}

func IsZeroValue(x interface{}) bool {
	switch reflect.TypeOf(x).Kind() {
	case reflect.Ptr:
		return x == nil || reflect.ValueOf(x).IsNil()
	case reflect.Struct:
		return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
	default:
		return reflect.ValueOf(x).IsZero()
	}
}
