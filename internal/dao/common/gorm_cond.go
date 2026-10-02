package common

import (
	"reflect"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormCond interface {
	BuildCond(query *gorm.DB) *gorm.DB
}

type EqCond struct {
	Field string
	Value interface{}
}

// GteCond applies a greater-than-or-equal comparison to an internal field.
type GteCond struct {
	Field string
	Value interface{}
}

func (c *GteCond) BuildCond(query *gorm.DB) *gorm.DB {
	if c.Value == nil {
		return query
	}
	value := reflect.ValueOf(c.Value)
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return query
		}
		return query.Where(c.Field+" >= ?", value.Elem().Interface())
	}
	return query.Where(c.Field+" >= ?", c.Value)
}

// LtCond applies a less-than comparison to an internal field.
type LtCond struct {
	Field string
	Value interface{}
}

func (c *LtCond) BuildCond(query *gorm.DB) *gorm.DB {
	if c.Value == nil {
		return query
	}
	value := reflect.ValueOf(c.Value)
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return query
		}
		return query.Where(c.Field+" < ?", value.Elem().Interface())
	}
	return query.Where(c.Field+" < ?", c.Value)
}

func (c *EqCond) BuildCond(query *gorm.DB) *gorm.DB {
	// 可空指针字段（如 *int8 的 state）为 nil 时跳过；零值标量（如 is_draft=0）仍正常过滤
	if c.Value == nil {
		return query
	}
	if rv := reflect.ValueOf(c.Value); rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return query
		}
		return query.Where(c.Field, rv.Elem().Interface())
	}
	return query.Where(c.Field, c.Value)
}

type LikeCond struct {
	Field string
	Value string
}

func (c *LikeCond) BuildCond(query *gorm.DB) *gorm.DB {
	if c.Value == "" {
		return query
	}
	return query.Where(clause.Like{
		Column: clause.Column{Name: c.Field},
		Value:  "%" + escapeLike(c.Value) + "%",
	})
}

// escapeLike 转义 LIKE 通配符（\ % _），避免用户输入被当作通配符匹配
func escapeLike(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
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
		// Start each child from a clean statement so parent WHERE clauses do not
		// become part of every OR branch.
		subQuery := query.Session(&gorm.Session{NewDB: true, DryRun: true, Initialized: true})
		subQuery = cond.BuildCond(subQuery)
		// no-op 子条件（如空 LikeCond）不产生 WHERE 表达式。
		if expr := subQuery.Statement.Clauses["WHERE"].Expression; expr != nil {
			queryClauses = append(queryClauses, expr)
		}
	}
	if len(queryClauses) == 0 {
		return query
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
