package link

import (
	"context"
	"gin-web/dto"
)

type LinkService interface {
	List(ctx context.Context, reqDTO *dto.ListLinkReqDTO) (*dto.ListLinkRespDTO, error)
	Add(ctx context.Context, reqDTO *dto.AddLinkReqDTO) error
	Edit(ctx context.Context, reqDTO *dto.EditLinkReqDTO) error
	Del(ctx context.Context, reqDTO *dto.IDReqDTO) error
}
