package article

import (
	"context"
	"gin-web/dto"
)

type ArticleService interface {
	List(ctx context.Context, reqDTO *dto.ListArticleReqDTO) (*dto.ListArticleRespDTO, error)
	AdminList(ctx context.Context, reqDTO *dto.AdminListArticleReqDTO) (*dto.ListArticleRespDTO, error)
	Detail(ctx context.Context, reqDTO *dto.IDReqDTO) (*dto.ArticleDetailRespDTO, error)
	DetailBySlug(ctx context.Context, reqDTO *dto.ArticleSlugReqDTO) (*dto.ArticleDetailRespDTO, error)
	AdminDetail(ctx context.Context, reqDTO *dto.IDReqDTO) (*dto.ArticleDetailRespDTO, error)
	Add(ctx context.Context, reqDTO *dto.AddArticleReqDTO) error
	Edit(ctx context.Context, reqDTO *dto.EditArticleReqDTO) error
	Del(ctx context.Context, reqDTO *dto.IDReqDTO) error
	Publish(ctx context.Context, reqDTO *dto.PublishArticleReqDTO) error
	State(ctx context.Context, reqDTO *dto.StateArticleReqDTO) error
	Top(ctx context.Context, reqDTO *dto.TopArticleReqDTO) error
	Archive(ctx context.Context) (*dto.ArchiveRespDTO, error)
	AllPublished(ctx context.Context) (*dto.ListArticleRespDTO, error)
}
