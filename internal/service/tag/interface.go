package tag

import (
	"context"
	"gin-web/dto"
)

type TagService interface {
	List(ctx context.Context, reqDTO *dto.ListTagReqDTO) (*dto.ListTagRespDTO, error)
}
