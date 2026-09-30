package tag

import (
	"context"
	"regexp"
	"testing"

	"gin-web/internal/dao/common"
	"gin-web/internal/model"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func newTagDao(t *testing.T) (TagDao, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("init sqlmock: %v", err)
	}
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
		// 测试中直接断言单条SQL，跳过GORM默认的写事务包裹
	}), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		t.Fatalf("open gorm: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	return NewTagDao(db), mock
}

// 验证删除接口真实执行DELETE（回归测试：Del曾是空实现）
func TestDeleteOne(t *testing.T) {
	dao, mock := newTagDao(t)
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `blog_tag` WHERE id = ?")).
		WithArgs(int64(7)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := dao.DeleteOne(context.Background(), 7); err != nil {
		t.Fatalf("DeleteOne err: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations not met: %v", err)
	}
}

func TestSelectOne(t *testing.T) {
	dao, mock := newTagDao(t)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `blog_tag` WHERE id = ?")).
		WithArgs(int64(3)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(3, "go"))

	tag, err := dao.SelectOne(context.Background(), 3)
	if err != nil {
		t.Fatalf("SelectOne err: %v", err)
	}
	if tag.Name != "go" {
		t.Fatalf("name = %q, want go", tag.Name)
	}
}

func TestSelectManyAppliesPagination(t *testing.T) {
	dao, mock := newTagDao(t)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `blog_tag` LIMIT 10 OFFSET 20")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

	_, err := dao.SelectMany(context.Background(), common.GormConditions{}, &common.Pagination{Offset: 20, PageSize: 10})
	if err != nil {
		t.Fatalf("SelectMany err: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations not met: %v", err)
	}
}

func TestUpdateOne(t *testing.T) {
	dao, mock := newTagDao(t)
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `blog_tag` SET")).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// 携带主键才会走UPDATE，零主键时Save会退化为INSERT
	po := &model.Tag{Model: model.Model{ID: 3}, Name: "go"}
	if err := dao.UpdateOne(context.Background(), po); err != nil {
		t.Fatalf("UpdateOne err: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations not met: %v", err)
	}
}
