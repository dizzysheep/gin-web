package dao

import (
	"gin-web/internal/dto"
	"gin-web/internal/models"
	"gorm.io/gorm"
)

type ArticleDao struct {
	*gorm.DB
}

func NewArticleDao(db *gorm.DB) *ArticleDao {
	return &ArticleDao{db}
}

func (dao *ArticleDao) CreateArticle(article *models.Article) (int, error) {
	err := dao.Create(article).Error
	if err != nil {
		return 0, err
	}

	return article.ID, nil
}

func (dao *ArticleDao) GetArticle(id int) (*models.Article, error) {
	article := &models.Article{}
	err := dao.Where("id = ?", id).First(article).Error
	if err != nil {
		return article, err
	}

	return article, nil
}

func (dao *ArticleDao) DeleteArticle(id int) error {
	return dao.Where("id = ?", id).Delete(&models.Article{}).Error
}

func (dao *ArticleDao) GetArticleByPage(req dto.SearchArticleReqDTO, pageNo, pageSize int) (total int64, list []models.Article, err error) {
	dao.SetCondition(req)
	offset := (pageNo - 1) * pageSize
	err = dao.Find(&list).Limit(pageSize).Offset(offset).Count(&total).Error
	if err != nil {
		return
	}

	return
}

// SetCondition ---设置查询条件
func (dao *ArticleDao) SetCondition(req dto.SearchArticleReqDTO) {
	if req.TagId != nil {
		dao.DB = dao.Where("tag_id = ?", req.TagId)
	}

	if req.State != nil {
		dao.DB = dao.Where("state = ?", req.State)
	}
}
