package link

import (
	"context"
	"time"

	"gin-web/internal/dao/common"
	"gin-web/internal/model"
	"gorm.io/gorm"
)

type linkDao struct {
	*gorm.DB
}

func NewLinkDao(db *gorm.DB) LinkDao {
	return &linkDao{db}
}

func (dao *linkDao) Count(ctx context.Context, conditions common.GormConditions) (int64, error) {
	var count int64
	err := dao.WithContext(ctx).Model(&model.Link{}).
		Scopes(common.NotDeleted(conditions).BuildConditions).
		Count(&count).Error
	return count, err
}

func (dao *linkDao) SelectMany(ctx context.Context, conditions common.GormConditions, pager *common.Pagination) ([]*model.Link, error) {
	var list []*model.Link
	tx := dao.WithContext(ctx).Model(&model.Link{}).
		Scopes(common.NotDeleted(conditions).BuildConditions)
	if pager != nil {
		tx = tx.Scopes(pager.Paginate())
	}
	err := tx.Find(&list).Error
	return list, err
}

func (dao *linkDao) SelectOne(ctx context.Context, id int64) (*model.Link, error) {
	var info *model.Link
	err := dao.WithContext(ctx).Model(&model.Link{}).
		Where("id = ? AND deleted_on = 0", id).First(&info).Error
	return info, err
}

func (dao *linkDao) InsertOne(ctx context.Context, link *model.Link) error {
	return dao.WithContext(ctx).Create(link).Error
}

func (dao *linkDao) UpdateOne(ctx context.Context, link *model.Link) error {
	return dao.WithContext(ctx).Save(link).Error
}

func (dao *linkDao) DeleteOne(ctx context.Context, id int64) error {
	return dao.WithContext(ctx).Model(&model.Link{}).
		Where("id = ? AND deleted_on = 0", id).
		Updates(map[string]interface{}{
			"deleted_on": time.Now().Unix(),
			"state":      model.StateDisabled,
		}).Error
}
