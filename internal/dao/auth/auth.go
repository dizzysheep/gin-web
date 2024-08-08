package auth

import (
	"context"
	"gin-web/internal/dao/common"
	"gin-web/internal/model"
	"gorm.io/gorm"
)

type authDao struct {
	*gorm.DB
}

func NewAuthDao(db *gorm.DB) AuthDao {
	return &authDao{db}
}

func (dao *authDao) Count(ctx context.Context, conditions common.GormConditions) (int64, error) {
	var count int64
	err := dao.WithContext(ctx).Model(&model.Auth{}).Scopes(conditions.BuildConditions).Count(&count).Error
	return count, err
}

func (dao *authDao) SelectMany(ctx context.Context, conditions common.GormConditions, pager *common.Pagination) ([]*model.Auth, error) {
	var list []*model.Auth
	tx := dao.WithContext(ctx).Model(&model.Auth{}).Scopes(conditions.BuildConditions)
	if pager != nil {
		tx.Scopes(pager.Paginate())
	}
	err := tx.Find(&list).Error
	return list, err
}

func (dao *authDao) SelectOne(ctx context.Context, id int64) (*model.Auth, error) {
	var info *model.Auth
	err := dao.WithContext(ctx).Model(&model.Auth{}).Where("id = ?", id).First(&info).Error
	return info, err
}

func (dao *authDao) SelectOneByWhere(ctx context.Context, conditions common.GormConditions) (*model.Auth, error) {
	var info *model.Auth
	err := dao.WithContext(ctx).Model(&model.Auth{}).Scopes(conditions.BuildConditions).First(&info).Error
	return info, err
}

func (dao *authDao) InsertOne(ctx context.Context, model *model.Auth) error {
	return dao.WithContext(ctx).Create(model).Error
}

func (dao *authDao) UpdateOne(ctx context.Context, model *model.Auth) error {
	return dao.WithContext(ctx).Save(model).Error
}
