package tag

import (
	"context"
	"gin-web/dto"
)

type TagService interface {
	List(ctx context.Context, reqDTO *dto.ListTagReqDTO) (*dto.ListTagRespDTO, error)
	Add(ctx context.Context, reqDTO *dto.AddTagReqDTO) error
	Edit(ctx context.Context, reqDTO *dto.EditTagReqDTO) error
	Del(ctx context.Context, reqDTO *dto.IDReqDTO) error
}
