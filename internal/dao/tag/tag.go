package tag

import (
	"context"
	"time"

	"gin-web/internal/dao/common"
	"gin-web/internal/model"
	"gorm.io/gorm"
)

type tagDao struct {
	*gorm.DB
}

func NewTagDao(db *gorm.DB) TagDao {
	return &tagDao{db}
}

func (dao *tagDao) Count(ctx context.Context, conditions common.GormConditions) (int64, error) {
	var count int64
	err := dao.WithContext(ctx).Model(&model.Tag{}).
		Scopes(common.NotDeleted(conditions).BuildConditions).
		Count(&count).Error
	return count, err
}

func (dao *tagDao) SelectMany(ctx context.Context, conditions common.GormConditions, pager *common.Pagination) ([]*model.Tag, error) {
	var list []*model.Tag
	tx := dao.WithContext(ctx).Model(&model.Tag{}).
		Scopes(common.NotDeleted(conditions).BuildConditions)
	if pager != nil {
		tx = tx.Scopes(pager.Paginate())
	}
	err := tx.Find(&list).Error
	return list, err
}

func (dao *tagDao) SelectOne(ctx context.Context, id int64) (*model.Tag, error) {
	var info *model.Tag
	err := dao.WithContext(ctx).Model(&model.Tag{}).
		Where("id = ? AND deleted_on = 0", id).First(&info).Error
	return info, err
}

// SelectManyByIDs 批量查询标签（列表页文章标签聚合）
func (dao *tagDao) SelectManyByIDs(ctx context.Context, ids []int64) ([]*model.Tag, error) {
	var list []*model.Tag
	if len(ids) == 0 {
		return list, nil
	}
	err := dao.WithContext(ctx).Model(&model.Tag{}).
		Where("id IN ? AND deleted_on = 0", ids).
		Find(&list).Error
	return list, err
}

// CountArticleByIDs 批量统计各标签关联的文章数
func (dao *tagDao) CountArticleByIDs(ctx context.Context, ids []int64) (map[int64]int64, error) {
	result := make(map[int64]int64, len(ids))
	if len(ids) == 0 {
		return result, nil
	}

	var rows []struct {
		TagID int64 `gorm:"column:tag_id"`
		Total int64 `gorm:"column:total"`
	}
	err := dao.WithContext(ctx).Model(&model.ArticleTag{}).
		Select("tag_id, COUNT(*) AS total").
		Where("tag_id IN ?", ids).
		Group("tag_id").Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		result[row.TagID] = row.Total
	}
	return result, nil
}

// SelectRelationsByArticleIDs 批量查询文章-标签关联（列表页聚合）
func (dao *tagDao) SelectRelationsByArticleIDs(ctx context.Context, articleIDs []int64) ([]*model.ArticleTag, error) {
	var list []*model.ArticleTag
	if len(articleIDs) == 0 {
		return list, nil
	}
	err := dao.WithContext(ctx).Model(&model.ArticleTag{}).
		Where("article_id IN ?", articleIDs).
		Find(&list).Error
	return list, err
}

// CountRelationsByTagID 标签关联的文章数（删除前置校验）
func (dao *tagDao) CountRelationsByTagID(ctx context.Context, tagID int64) (int64, error) {
	var count int64
	err := dao.WithContext(ctx).Model(&model.ArticleTag{}).
		Where("tag_id = ?", tagID).Count(&count).Error
	return count, err
}

func (dao *tagDao) ExistName(ctx context.Context, name string, excludeID int64) (bool, error) {
	var count int64
	err := dao.WithContext(ctx).Model(&model.Tag{}).
		Where("name = ? AND deleted_on = 0 AND id != ?", name, excludeID).
		Count(&count).Error
	return count > 0, err
}

func (dao *tagDao) InsertOne(ctx context.Context, tag *model.Tag) error {
	return dao.WithContext(ctx).Create(tag).Error
}

func (dao *tagDao) UpdateOne(ctx context.Context, tag *model.Tag) error {
	return dao.WithContext(ctx).Save(tag).Error
}

func (dao *tagDao) DeleteOne(ctx context.Context, id int64) error {
	return dao.WithContext(ctx).Model(&model.Tag{}).
		Where("id = ? AND deleted_on = 0", id).
		Updates(map[string]interface{}{
			"deleted_on": time.Now().Unix(),
			"state":      model.StateDisabled,
		}).Error
}
