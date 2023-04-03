package service

import (
	"errors"
	"gin-demo/models"
	"gin-demo/validators/article"
	"gin-demo/validators/tag"
)

type ArticleService struct {
	Model *models.Article
}

func NewArticleService() *ArticleService {
	return &ArticleService{
		Model: models.NewArticle(),
	}
}

func (s *ArticleService) GetListByPage(params map[string]interface{}, pageNo, pageSize int) (total int, list []models.Article, err error) {
	s.HandleWhere(params)
	total, list, err = s.Model.GetListByPage(params, pageNo, pageSize)
	if err != nil {
		return total, list, err
	}
	return total, list, nil
}

func (s *ArticleService) HandleWhere(params map[string]interface{}) {
	if params["state"] == -1 {
		delete(params, "state")
	}
	if params["tag_id"] == -1 {
		delete(params, "tag_id")
	}
}

func (s *ArticleService) GetArticle(id int) (*models.Article, error) {
	return s.Model.GetArticle(id)
}

func (s *ArticleService) AddArticle(article article.AddArticle) (int, error) {
	s.Model.TagID = article.TagId
	s.Model.Title = article.Title
	s.Model.Desc = article.Desc
	s.Model.Content = article.Content
	s.Model.CreatedBy = article.CreatedBy
	s.Model.State = article.State
	return s.Model.Create()
}

func (s *ArticleService) EditArticle(tag tag.EditTag) error {

	return s.Model.Update()
}

func (s *ArticleService) Delete(id int) error {
	articleDetail, err := s.Model.GetArticle(id)
	if err != nil {
		return err
	}
	if articleDetail.ID == 0 {
		return errors.New("文章记录不存在")
	}

	if err := s.Model.Delete(articleDetail.ID); err != nil {
		return errors.New("删除文章失败")
	}
	return nil
}
