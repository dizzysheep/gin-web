package article

import (
	"context"
	"gin-web/dto"
	"gin-web/internal/dao"
)

type articleService struct {
	daos *dao.Daos
}

func NewArticleService(daos *dao.Daos) ArticleService {
	return &articleService{daos}
}

func (s *articleService) List(ctx context.Context, reqDTO *dto.ListArticleReqDTO) (*dto.ListArticleRespDTO, error) {
	return nil, nil
}
