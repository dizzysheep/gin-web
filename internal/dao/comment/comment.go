package comment

import (
	"context"
	"time"

	"gin-web/internal/dao/common"
	"gin-web/internal/model"
	"gorm.io/gorm"
)

type commentDao struct {
	*gorm.DB
}

func NewCommentDao(db *gorm.DB) CommentDao {
	return &commentDao{db}
}

func (dao *commentDao) Count(ctx context.Context, conditions common.GormConditions) (int64, error) {
	var count int64
	err := dao.WithContext(ctx).Model(&model.Comment{}).
		Scopes(common.NotDeleted(conditions).BuildConditions).
		Count(&count).Error
	return count, err
}

func (dao *commentDao) SelectMany(ctx context.Context, conditions common.GormConditions, pager *common.Pagination) ([]*model.Comment, error) {
	var list []*model.Comment
	tx := dao.WithContext(ctx).Model(&model.Comment{}).
		Scopes(common.NotDeleted(conditions).BuildConditions)
	if pager != nil {
		tx = tx.Scopes(pager.Paginate())
	}
	err := tx.Find(&list).Error
	return list, err
}

func (dao *commentDao) SelectOne(ctx context.Context, id int64) (*model.Comment, error) {
	var info *model.Comment
	err := dao.WithContext(ctx).Model(&model.Comment{}).
		Where("id = ? AND deleted_on = 0", id).First(&info).Error
	return info, err
}

func (dao *commentDao) InsertOne(ctx context.Context, comment *model.Comment) error {
	return dao.WithContext(ctx).Create(comment).Error
}

// UpdateColumns 局部更新（审核状态等）
func (dao *commentDao) UpdateColumns(ctx context.Context, id int64, cols map[string]interface{}) error {
	return dao.WithContext(ctx).Model(&model.Comment{}).
		Where("id = ? AND deleted_on = 0", id).Updates(cols).Error
}

func (dao *commentDao) DeleteOne(ctx context.Context, id int64) error {
	return dao.WithContext(ctx).Model(&model.Comment{}).
		Where("id = ? AND deleted_on = 0", id).
		Updates(map[string]interface{}{
			"deleted_on": time.Now().Unix(),
			"state":      model.CommentStateReject,
		}).Error
}
