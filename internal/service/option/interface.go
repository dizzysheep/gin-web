package option

import (
	"context"
	"gin-web/dto"
)

type OptionService interface {
	Get(ctx context.Context) (*dto.GetOptionsRespDTO, error)
	Value(ctx context.Context, key string) string
	Save(ctx context.Context, reqDTO *dto.SaveOptionsReqDTO) error
}
