package user

import (
	"context"
	"gin-web/internal/dao/common"
	"gin-web/internal/model"
	"gorm.io/gorm"
)

type userDao struct {
	*gorm.DB
}

func NewUserDao(db *gorm.DB) UserDao {
	return &userDao{db}
}

func (dao *userDao) SelectOne(ctx context.Context, id int64) (*model.User, error) {
	var info *model.User
	err := dao.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).First(&info).Error
	return info, err
}

func (dao *userDao) SelectOneByUsername(ctx context.Context, username string) (*model.User, error) {
	var info *model.User
	conditions := common.GormConditions{
		&common.EqCond{Field: "username", Value: username},
	}
	err := dao.WithContext(ctx).Model(&model.User{}).Scopes(conditions.BuildConditions).First(&info).Error
	return info, err
}

func (dao *userDao) UpdateOne(ctx context.Context, user *model.User) error {
	return dao.WithContext(ctx).Save(user).Error
}
