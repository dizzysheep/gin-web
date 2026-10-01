package comment

import (
	"context"
	"gin-web/dto"
)

type CommentService interface {
	ListByArticle(ctx context.Context, reqDTO *dto.ListCommentReqDTO) (*dto.ListCommentRespDTO, error)
	Add(ctx context.Context, reqDTO *dto.AddCommentReqDTO) error
	AdminList(ctx context.Context, reqDTO *dto.AdminListCommentReqDTO) (*dto.AdminListCommentRespDTO, error)
	Audit(ctx context.Context, reqDTO *dto.AuditCommentReqDTO) error
	Del(ctx context.Context, reqDTO *dto.IDReqDTO) error
}
