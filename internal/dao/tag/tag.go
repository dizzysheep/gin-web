package tag

import (
	"context"
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
	err := dao.WithContext(ctx).Model(&model.Tag{}).Scopes(conditions.BuildConditions).Count(&count).Error
	return count, err
}

func (dao *tagDao) SelectMany(ctx context.Context, conditions common.GormConditions, pager *common.Pagination) ([]*model.Tag, error) {
	var list []*model.Tag
	tx := dao.WithContext(ctx).Model(&model.Tag{}).Scopes(conditions.BuildConditions)
	if pager != nil {
		tx.Scopes(pager.Paginate())
	}
	err := tx.Find(&list).Error
	return list, err
}

func (dao *tagDao) SelectOne(ctx context.Context, id int64) (*model.Tag, error) {
	var info *model.Tag
	err := dao.WithContext(ctx).Model(&model.Tag{}).Where("id = ?", id).First(&info).Error
	return info, err
}

func (dao *tagDao) InsertOne(ctx context.Context, model *model.Tag) error {
	return dao.WithContext(ctx).Create(model).Error
}

func (dao *tagDao) UpdateOne(ctx context.Context, model *model.Tag) error {
	return dao.WithContext(ctx).Save(model).Error
}
