package category

import (
	"context"
	"time"

	"gin-web/internal/dao/common"
	"gin-web/internal/model"
	"gorm.io/gorm"
)

type categoryDao struct {
	*gorm.DB
}

func NewCategoryDao(db *gorm.DB) CategoryDao {
	return &categoryDao{db}
}

func (dao *categoryDao) Count(ctx context.Context, conditions common.GormConditions) (int64, error) {
	var count int64
	err := dao.WithContext(ctx).Model(&model.Category{}).
		Scopes(common.NotDeleted(conditions).BuildConditions).
		Count(&count).Error
	return count, err
}

func (dao *categoryDao) SelectMany(ctx context.Context, conditions common.GormConditions, pager *common.Pagination) ([]*model.Category, error) {
	var list []*model.Category
	tx := dao.WithContext(ctx).Model(&model.Category{}).
		Scopes(common.NotDeleted(conditions).BuildConditions)
	if pager != nil {
		tx = tx.Scopes(pager.Paginate())
	}
	err := tx.Find(&list).Error
	return list, err
}

func (dao *categoryDao) SelectOne(ctx context.Context, id int64) (*model.Category, error) {
	var info *model.Category
	err := dao.WithContext(ctx).Model(&model.Category{}).
		Where("id = ? AND deleted_on = 0", id).First(&info).Error
	return info, err
}

// SelectManyByIDs 批量查询分类（列表页文章分类聚合）
func (dao *categoryDao) SelectManyByIDs(ctx context.Context, ids []int64) ([]*model.Category, error) {
	var list []*model.Category
	if len(ids) == 0 {
		return list, nil
	}
	err := dao.WithContext(ctx).Model(&model.Category{}).
		Where("id IN ? AND deleted_on = 0", ids).
		Find(&list).Error
	return list, err
}

// ExistName 分类名是否已存在（excludeID 用于编辑时排除自身）
func (dao *categoryDao) ExistName(ctx context.Context, name string, excludeID int64) (bool, error) {
	var count int64
	err := dao.WithContext(ctx).Model(&model.Category{}).
		Where("name = ? AND deleted_on = 0 AND id != ?", name, excludeID).
		Count(&count).Error
	return count > 0, err
}

// CountArticleByIDs 批量统计各分类下的文章数
func (dao *categoryDao) CountArticleByIDs(ctx context.Context, ids []int64) (map[int64]int64, error) {
	result := make(map[int64]int64, len(ids))
	if len(ids) == 0 {
		return result, nil
	}

	var rows []struct {
		CategoryID int64 `gorm:"column:category_id"`
		Total      int64 `gorm:"column:total"`
	}
	err := dao.WithContext(ctx).Model(&model.Article{}).
		Select("category_id, COUNT(*) AS total").
		Where("category_id IN ? AND deleted_on = 0", ids).
		Group("category_id").Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		result[row.CategoryID] = row.Total
	}
	return result, nil
}

func (dao *categoryDao) InsertOne(ctx context.Context, category *model.Category) error {
	return dao.WithContext(ctx).Create(category).Error
}

func (dao *categoryDao) UpdateOne(ctx context.Context, category *model.Category) error {
	return dao.WithContext(ctx).Save(category).Error
}

func (dao *categoryDao) DeleteOne(ctx context.Context, id int64) error {
	return dao.WithContext(ctx).Model(&model.Category{}).
		Where("id = ? AND deleted_on = 0", id).
		Updates(map[string]interface{}{
			"deleted_on": time.Now().Unix(),
			"state":      model.StateDisabled,
		}).Error
}
