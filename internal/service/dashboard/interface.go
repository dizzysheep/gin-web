package dashboard

import (
	"context"

	"gin-web/dto"
)

type DashboardService interface {
	Overview(ctx context.Context) (*dto.DashboardOverviewRespDTO, error)
	Trend(ctx context.Context, days int) (*dto.DashboardTrendRespDTO, error)
}
