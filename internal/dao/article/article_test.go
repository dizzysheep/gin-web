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

// 验证带分页参数时SQL确实携带LIMIT/OFFSET，且默认过滤软删除
func TestSelectManyAppliesPagination(t *testing.T) {
	dao, mock := newArticleDao(t)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `blog_article` WHERE deleted_on = 0 LIMIT 10 OFFSET 20")).
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
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `blog_article` WHERE deleted_on = 0")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title"}))

	_, err := dao.SelectMany(context.Background(), common.GormConditions{}, nil)
	if err != nil {
		t.Fatalf("SelectMany err: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations not met: %v", err)
	}
}

// 按标签查询走子查询过滤关联
func TestSelectManyByTagID(t *testing.T) {
	dao, mock := newArticleDao(t)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `blog_article` WHERE id IN (SELECT `article_id` FROM `blog_article_tag` WHERE tag_id = ?) AND deleted_on = 0")).
		WithArgs(int64(5)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title"}))

	_, err := dao.SelectManyByTagID(context.Background(), common.GormConditions{}, nil, 5)
	if err != nil {
		t.Fatalf("SelectManyByTagID err: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations not met: %v", err)
	}
}

func TestCount(t *testing.T) {
	dao, mock := newArticleDao(t)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `blog_article` WHERE deleted_on = 0")).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(42))

	total, err := dao.Count(context.Background(), common.GormConditions{})
	if err != nil {
		t.Fatalf("Count err: %v", err)
	}
	if total != 42 {
		t.Fatalf("total = %d, want 42", total)
	}
}

// Create 在事务内写入文章并清理旧标签关联
func TestCreate(t *testing.T) {
	dao, mock := newArticleDao(t)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `blog_article`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `blog_article_tag` WHERE article_id = ?")).
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	if err := dao.Create(context.Background(), &model.Article{Title: "hello"}, nil); err != nil {
		t.Fatalf("Create err: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations not met: %v", err)
	}
}

// Delete 软删除文章并清理标签关联
func TestDelete(t *testing.T) {
	dao, mock := newArticleDao(t)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `blog_article` SET `deleted_on`=?,`state`=? WHERE id = ? AND deleted_on = 0")).
		WithArgs(sqlmock.AnyArg(), int8(0), int64(7)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `blog_article_tag` WHERE article_id = ?")).
		WithArgs(int64(7)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	if err := dao.Delete(context.Background(), 7); err != nil {
		t.Fatalf("Delete err: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations not met: %v", err)
	}
}

// IncrViewCount 浏览量原子累加
func TestIncrViewCount(t *testing.T) {
	dao, mock := newArticleDao(t)
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `blog_article` SET view_count=view_count+? WHERE id = ? AND deleted_on = 0")).
		WithArgs(int64(3), int64(7)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := dao.IncrViewCount(context.Background(), 7, 3); err != nil {
		t.Fatalf("IncrViewCount err: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations not met: %v", err)
	}
}
