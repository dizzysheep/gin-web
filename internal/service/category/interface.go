package category

import (
	"context"
	"gin-web/dto"
)

type CategoryService interface {
	List(ctx context.Context, reqDTO *dto.ListCategoryReqDTO) (*dto.ListCategoryRespDTO, error)
	Add(ctx context.Context, reqDTO *dto.AddCategoryReqDTO) error
	Edit(ctx context.Context, reqDTO *dto.EditCategoryReqDTO) error
	Del(ctx context.Context, reqDTO *dto.IDReqDTO) error
}
