package article

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

func newArticleDao(t *testing.T) (ArticleDao, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("init sqlmock: %v", err)
	}
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		t.Fatalf("open gorm: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	return NewArticleDao(db), mock
}

// 验证带分页参数时SQL确实携带LIMIT/OFFSET（回归测试：曾因丢弃Scopes返回值导致全表查询）
func TestSelectManyAppliesPagination(t *testing.T) {
	dao, mock := newArticleDao(t)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `blog_article` LIMIT 10 OFFSET 20")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title"}))

	_, err := dao.SelectMany(context.Background(), common.GormConditions{}, &common.Pagination{Offset: 20, PageSize: 10})
	if err != nil {
		t.Fatalf("SelectMany err: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations not met: %v", err)
	}
}

func TestSelectManyWithoutPager(t *testing.T) {
	dao, mock := newArticleDao(t)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `blog_article`")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title"}))

	_, err := dao.SelectMany(context.Background(), common.GormConditions{}, nil)
	if err != nil {
		t.Fatalf("SelectMany err: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations not met: %v", err)
	}
}

func TestCount(t *testing.T) {
	dao, mock := newArticleDao(t)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `blog_article`")).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(42))

	total, err := dao.Count(context.Background(), common.GormConditions{})
	if err != nil {
		t.Fatalf("Count err: %v", err)
	}
	if total != 42 {
		t.Fatalf("total = %d, want 42", total)
	}
}

func TestInsertOne(t *testing.T) {
	dao, mock := newArticleDao(t)
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `blog_article`")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := dao.InsertOne(context.Background(), &model.Article{Title: "hello"}); err != nil {
		t.Fatalf("InsertOne err: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations not met: %v", err)
	}
}
