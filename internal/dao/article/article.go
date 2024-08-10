package article

import (
	"context"
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
	err := dao.WithContext(ctx).Model(&model.Article{}).Scopes(conditions.BuildConditions).Count(&count).Error
	return count, err
}

func (dao *articleDao) SelectMany(ctx context.Context, conditions common.GormConditions, pager *common.Pagination) ([]*model.Article, error) {
	var list []*model.Article
	tx := dao.WithContext(ctx).Model(&model.Article{}).Preload("Tag").Scopes(conditions.BuildConditions)
	if pager != nil {
		tx.Scopes(pager.Paginate())
	}
	err := tx.Find(&list).Error
	return list, err
}

func (dao *articleDao) SelectOne(ctx context.Context, id int64) (*model.Article, error) {
	var info *model.Article
	err := dao.WithContext(ctx).Model(&model.Article{}).Where("id = ?", id).First(&info).Error
	return info, err
}

func (dao *articleDao) InsertOne(ctx context.Context, model *model.Article) error {
	return dao.WithContext(ctx).Create(model).Error
}

func (dao *articleDao) UpdateOne(ctx context.Context, model *model.Article) error {
	return dao.WithContext(ctx).Save(model).Error
}
