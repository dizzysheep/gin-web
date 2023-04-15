package service

import (
	"gin-web/internal/dao"
	"gin-web/internal/dto"
	"gin-web/internal/models"
)

type ArticleService struct {
}

func NewArticleService() *ArticleService {
	return &ArticleService{}
}

// GetListByPage ---查询文章列表---
func (svc *ArticleService) GetListByPage(req dto.SearchArticleReqDTO, pageNo, pageSize int) (total int, list []models.Article, err error) {
	return dao.NewArticleDao().GetArticleByPage(req, pageNo, pageSize)
}

// GetArticle ---获取文章详情---
func (svc *ArticleService) GetArticle(id int) (*models.Article, error) {
	return dao.NewArticleDao().GetArticle(id)
}

// CreateArticle ---创建文章---
func (svc *ArticleService) CreateArticle(article dto.CreateArticleReqDTO) (int, error) {
	articleModel := &models.Article{
		TagID:     article.TagId,
		Title:     article.Title,
		Desc:      article.Desc,
		Content:   article.Content,
		CreatedBy: article.CreatedBy,
		State:     article.State,
	}
	return dao.NewArticleDao().CreateArticle(articleModel)
}

// EditArticle ---获取文章详情---
func (svc *ArticleService) EditArticle() error {

	return nil
}

// DeleteArticle ---删除文章详情---
func (svc *ArticleService) DeleteArticle(id int) error {
	return dao.NewArticleDao().DeleteArticle(id)
}
