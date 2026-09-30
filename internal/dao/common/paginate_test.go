package common

import (
	"strings"
	"testing"

	"gin-web/internal/model"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("init sqlmock: %v", err)
	}
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("open gorm: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	return db, mock
}

// buildSQL 以DryRun模式构建查询，返回生成的SQL与绑定参数，不实际执行
func buildSQL(t *testing.T, conds GormConditions, pager *Pagination) (string, []interface{}) {
	t.Helper()
	db, _ := newMockDB(t)
	tx := db.Session(&gorm.Session{DryRun: true}).Model(&model.Tag{}).Scopes(conds.BuildConditions)
	if pager != nil {
		tx = tx.Scopes(pager.Paginate())
	}
	if err := tx.Find(&[]model.Tag{}).Error; err != nil {
		t.Fatalf("build query: %v", err)
	}
	return tx.Statement.SQL.String(), tx.Statement.Vars
}

func containsVar(vars []interface{}, want interface{}) bool {
	for _, v := range vars {
		if v == want {
			return true
		}
	}
	return false
}

// 验证分页Scope会生成LIMIT/OFFSET子句（曾因丢弃Scopes返回值导致分页失效）
func TestPaginateScope(t *testing.T) {
	sql, _ := buildSQL(t, GormConditions{}, &Pagination{Offset: 20, PageSize: 10})

	// GORM构建SQL时会将LIMIT/OFFSET的值直接内联
	if !strings.Contains(sql, "LIMIT 10") {
		t.Fatalf("expect LIMIT 10 in sql, got: %s", sql)
	}
	if !strings.Contains(sql, "OFFSET 20") {
		t.Fatalf("expect OFFSET 20 in sql, got: %s", sql)
	}
}

func TestNoPaginateWithoutPager(t *testing.T) {
	sql, _ := buildSQL(t, GormConditions{}, nil)
	if strings.Contains(sql, "LIMIT") || strings.Contains(sql, "OFFSET") {
		t.Fatalf("unexpected LIMIT/OFFSET in sql: %s", sql)
	}
}
