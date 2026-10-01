package article

import (
	"context"
	"time"

	"gin-web/internal/dao/common"
	"gin-web/internal/model"
	"gorm.io/gorm"
)

type articleDao struct {
	*gorm.DB
}

func NewArticleDao(db *gorm.DB) ArticleDao {
	return &articleDao{db}
}

func (dao *articleDao) Count(ctx context.Context, conditions common.GormConditions) (int64, error) {
	var count int64
	err := dao.WithContext(ctx).Model(&model.Article{}).
		Scopes(common.NotDeleted(conditions).BuildConditions).
		Count(&count).Error
	return count, err
}

func (dao *articleDao) SelectMany(ctx context.Context, conditions common.GormConditions, pager *common.Pagination) ([]*model.Article, error) {
	var list []*model.Article
	tx := dao.WithContext(ctx).Model(&model.Article{}).
		Scopes(common.NotDeleted(conditions).BuildConditions)
	if pager != nil {
		tx = tx.Scopes(pager.Paginate())
	}
	err := tx.Find(&list).Error
	return list, err
}

// CountByTagID 按标签统计文章数（子查询过滤关联）
func (dao *articleDao) CountByTagID(ctx context.Context, conditions common.GormConditions, tagID int64) (int64, error) {
	var count int64
	err := dao.WithContext(ctx).Model(&model.Article{}).
		Where("id IN (?)", dao.tagArticleIDs(tagID)).
		Scopes(common.NotDeleted(conditions).BuildConditions).
		Count(&count).Error
	return count, err
}

// SelectManyByTagID 按标签查询文章列表
func (dao *articleDao) SelectManyByTagID(ctx context.Context, conditions common.GormConditions, pager *common.Pagination, tagID int64) ([]*model.Article, error) {
	var list []*model.Article
	tx := dao.WithContext(ctx).Model(&model.Article{}).
		Where("id IN (?)", dao.tagArticleIDs(tagID)).
		Scopes(common.NotDeleted(conditions).BuildConditions)
	if pager != nil {
		tx = tx.Scopes(pager.Paginate())
	}
	err := tx.Find(&list).Error
	return list, err
}

func (dao *articleDao) SelectOne(ctx context.Context, id int64) (*model.Article, error) {
	var info *model.Article
	err := dao.WithContext(ctx).Model(&model.Article{}).
		Where("id = ? AND deleted_on = 0", id).First(&info).Error
	return info, err
}

func (dao *articleDao) SelectOneBySlug(ctx context.Context, slug string) (*model.Article, error) {
	var info *model.Article
	err := dao.WithContext(ctx).Model(&model.Article{}).
		Where("slug = ? AND deleted_on = 0", slug).First(&info).Error
	return info, err
}

func (dao *articleDao) ExistSlug(ctx context.Context, slug string, excludeID int64) (bool, error) {
	var count int64
	err := dao.WithContext(ctx).Model(&model.Article{}).
		Where("slug = ? AND deleted_on = 0 AND id != ?", slug, excludeID).
		Count(&count).Error
	return count > 0, err
}

// Archive 按月归档已发布文章
func (dao *articleDao) Archive(ctx context.Context) ([]*model.ArchiveRow, error) {
	var rows []*model.ArchiveRow
	err := dao.WithContext(ctx).Model(&model.Article{}).
		Select("FROM_UNIXTIME(published_on, '%Y-%m') AS month, COUNT(*) AS total").
		Where("deleted_on = 0 AND state = 1 AND is_draft = 0").
		Group("month").Order("month DESC").
		Scan(&rows).Error
	return rows, err
}

func (dao *articleDao) Create(ctx context.Context, article *model.Article, tagIDs []int64) error {
	return dao.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(article).Error; err != nil {
			return err
		}
		return replaceTagRelations(tx, article.ID, tagIDs)
	})
}

// Update 更新文章，withTags 为 true 时同步替换标签关联
func (dao *articleDao) Update(ctx context.Context, article *model.Article, tagIDs []int64, withTags bool) error {
	return dao.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(article).Error; err != nil {
			return err
		}
		if withTags {
			return replaceTagRelations(tx, article.ID, tagIDs)
		}
		return nil
	})
}

// UpdateColumns 局部更新（发布/置顶/启停/浏览量）
func (dao *articleDao) UpdateColumns(ctx context.Context, id int64, cols map[string]interface{}) error {
	return dao.WithContext(ctx).Model(&model.Article{}).
		Where("id = ? AND deleted_on = 0", id).Updates(cols).Error
}

// Delete 软删除文章并清理标签关联
func (dao *articleDao) Delete(ctx context.Context, id int64) error {
	return dao.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&model.Article{}).
			Where("id = ? AND deleted_on = 0", id).
			Updates(map[string]interface{}{
				"deleted_on": time.Now().Unix(),
				"state":      model.StateDisabled,
			}).Error
		if err != nil {
			return err
		}
		return tx.Where("article_id = ?", id).Delete(&model.ArticleTag{}).Error
	})
}

// IncrViewCount 浏览量落库（定时任务调用）
func (dao *articleDao) IncrViewCount(ctx context.Context, id int64, delta int64) error {
	return dao.WithContext(ctx).Exec(
		"UPDATE `blog_article` SET view_count=view_count+? WHERE id = ? AND deleted_on = 0",
		delta, id,
	).Error
}

// tagArticleIDs 构造按标签查文章的子查询
func (dao *articleDao) tagArticleIDs(tagID int64) *gorm.DB {
	return dao.Session(&gorm.Session{NewDB: true}).
		Model(&model.ArticleTag{}).Select("article_id").Where("tag_id = ?", tagID)
}

// replaceTagRelations 全量替换文章标签关联
func replaceTagRelations(tx *gorm.DB, articleID int64, tagIDs []int64) error {
	if err := tx.Where("article_id = ?", articleID).Delete(&model.ArticleTag{}).Error; err != nil {
		return err
	}
	if len(tagIDs) == 0 {
		return nil
	}
	relations := make([]*model.ArticleTag, 0, len(tagIDs))
	for _, tagID := range tagIDs {
		relations = append(relations, &model.ArticleTag{ArticleID: articleID, TagID: tagID})
	}
	return tx.Create(&relations).Error
}
