package option

import (
	"context"

	"gin-web/internal/dao/common"
	"gin-web/internal/model"
	"gorm.io/gorm"
)

type optionDao struct {
	*gorm.DB
}

func NewOptionDao(db *gorm.DB) OptionDao {
	return &optionDao{db}
}

// SelectAll 全量查询启用中的配置
func (dao *optionDao) SelectAll(ctx context.Context) ([]*model.Option, error) {
	var list []*model.Option
	conditions := common.GormConditions{
		&common.NotDeletedCond{},
		&common.EqCond{Field: "state", Value: model.StateEnabled},
	}
	err := dao.WithContext(ctx).Model(&model.Option{}).
		Scopes(conditions.BuildConditions).Find(&list).Error
	return list, err
}

// SelectOneByKey 按配置键查询
func (dao *optionDao) SelectOneByKey(ctx context.Context, key string) (*model.Option, error) {
	var info *model.Option
	err := dao.WithContext(ctx).Model(&model.Option{}).
		Where("`key` = ? AND deleted_on = 0", key).First(&info).Error
	return info, err
}

// Upsert 新增或更新配置项
func (dao *optionDao) Upsert(ctx context.Context, option *model.Option) error {
	exist, err := dao.SelectOneByKey(ctx, option.Key)
	if err != nil {
		if !isRecordNotFound(err) {
			return err
		}
		return dao.WithContext(ctx).Create(option).Error
	}

	exist.Value = option.Value
	exist.ModifiedOn = option.ModifiedOn
	exist.ModifiedBy = option.ModifiedBy
	return dao.WithContext(ctx).Save(exist).Error
}

func isRecordNotFound(err error) bool {
	return err == gorm.ErrRecordNotFound
}
