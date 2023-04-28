package service

import (
	"gin-web/internal/dao"
	"gin-web/internal/dto"
	"gin-web/internal/models"
	"github.com/gin-gonic/gin"
)

type ArticleService struct {
	dao *dao.ArticleDao
}

func NewArticleService(ctx *gin.Context) *ArticleService {
	return &ArticleService{dao.NewArticleDao(GetBlogDB(ctx))}
}

func (svc *ArticleService) GetListByPage(req dto.SearchArticleReqDTO, pageNo, pageSize int) (total int64, list []models.Article, err error) {
	return svc.dao.GetArticleByPage(req, pageNo, pageSize)
}

func (svc *ArticleService) GetArticle(id int) (*models.Article, error) {
	return svc.dao.GetArticle(id)
}

func (svc *ArticleService) CreateArticle(article dto.CreateArticleReqDTO) (int, error) {
	articleModel := &models.Article{
		TagID:     article.TagId,
		Title:     article.Title,
		Desc:      article.Desc,
		Content:   article.Content,
		CreatedBy: article.CreatedBy,
		State:     article.State,
	}
	return svc.dao.CreateArticle(articleModel)
}

func (svc *ArticleService) EditArticle() error {

	return nil
}

func (svc *ArticleService) DeleteArticle(id int) error {
	return svc.dao.DeleteArticle(id)
}
