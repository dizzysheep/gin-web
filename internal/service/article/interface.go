package article

import (
	"context"
	"gin-web/dto"
)

type ArticleService interface {
	List(ctx context.Context, reqDTO *dto.ListArticleReqDTO) (*dto.ListArticleRespDTO, error)
	Add(ctx context.Context, reqDTO *dto.AddArticleReqDTO) error
}
