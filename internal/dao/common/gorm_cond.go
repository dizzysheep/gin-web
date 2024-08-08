package common

import (
	"fmt"
	"gorm.io/gorm"
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
		return query.Where(fmt.Sprintf("%s LIKE ?", c.Field), fmt.Sprintf("%%%s%%", c.Value))
	}

	return query
}

type GormConditions []GormCond

func (c GormConditions) BuildConditions(query *gorm.DB) *gorm.DB {
	for _, cond := range c {
		query = cond.BuildCond(query)
	}
	return query
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
