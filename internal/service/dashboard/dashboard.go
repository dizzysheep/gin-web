package dashboard

import (
	"context"
	"strings"
	"time"

	"gin-web/core/xtime"
	"gin-web/dto"
	"gin-web/internal/dao"
	"gin-web/internal/dao/common"
	"gin-web/internal/model"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

const (
	defaultTrendDays = 7
	maxTrendDays     = 31
)

type dashboardService struct {
	daos *dao.Daos
}

func NewDashboardService(daos *dao.Daos) DashboardService {
	return &dashboardService{daos: daos}
}

// Overview returns the same public-content scope used by the dashboard cards:
// enabled, published articles/categories/tags and pending comments.
func (s *dashboardService) Overview(ctx context.Context) (*dto.DashboardOverviewRespDTO, error) {
	var (
		result dto.DashboardOverviewRespDTO
		eg     errgroup.Group
	)

	eg.Go(func() error {
		count, err := s.daos.Article.Count(ctx, common.GormConditions{
			&common.EqCond{Field: "state", Value: model.StateEnabled},
			&common.EqCond{Field: "is_draft", Value: 0},
		})
		if err != nil {
			return errors.Wrap(err, "count dashboard articles")
		}
		result.ArticleCount = count
		return nil
	})

	eg.Go(func() error {
		count, err := s.daos.Category.Count(ctx, common.GormConditions{
			&common.EqCond{Field: "state", Value: model.StateEnabled},
		})
		if err != nil {
			return errors.Wrap(err, "count dashboard categories")
		}
		result.CategoryCount = count
		return nil
	})

	eg.Go(func() error {
		count, err := s.daos.Tag.Count(ctx, common.GormConditions{
			&common.EqCond{Field: "state", Value: model.StateEnabled},
		})
		if err != nil {
			return errors.Wrap(err, "count dashboard tags")
		}
		result.TagCount = count
		return nil
	})

	eg.Go(func() error {
		count, err := s.daos.Comment.Count(ctx, common.GormConditions{
			&common.EqCond{Field: "state", Value: model.CommentStatePending},
		})
		if err != nil {
			return errors.Wrap(err, "count dashboard comments")
		}
		result.PendingCommentCount = count
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return &result, nil
}

// Trend returns a complete day series so the frontend can render zero-activity days.
func (s *dashboardService) Trend(ctx context.Context, days int) (*dto.DashboardTrendRespDTO, error) {
	if days <= 0 {
		days = defaultTrendDays
	}
	if days > maxTrendDays {
		days = maxTrendDays
	}

	now := time.Now()
	endDay := xtime.ToDayStart(now).AddDate(0, 0, 1)
	startDay := endDay.AddDate(0, 0, -days)
	rows, err := s.daos.Article.CountPublishedByDay(ctx, startDay.Unix(), endDay.Unix())
	if err != nil {
		return nil, errors.Wrap(err, "count dashboard article trend")
	}

	counts := make(map[string]int64, len(rows))
	for _, row := range rows {
		if row != nil {
			// Keep the lookup resilient to drivers returning a date or datetime
			// representation instead of the requested YYYY-MM-DD string.
			day := strings.TrimSpace(row.Day)
			if len(day) >= len(xtime.DATE_FMT) {
				day = day[:len(xtime.DATE_FMT)]
			}
			counts[day] = row.Total
		}
	}

	items := make([]*dto.DashboardTrendItem, 0, days)
	for i := 0; i < days; i++ {
		day := startDay.AddDate(0, 0, i)
		date := day.Format(xtime.DATE_FMT)
		items = append(items, &dto.DashboardTrendItem{
			Date:  date,
			Label: trendDayLabel(day, now),
			Total: counts[date],
		})
	}

	return &dto.DashboardTrendRespDTO{Days: days, Items: items}, nil
}

func trendDayLabel(day, now time.Time) string {
	if xtime.ToDayStart(day).Equal(xtime.ToDayStart(now)) {
		return "今天"
	}
	labels := [...]string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}
	return labels[day.Weekday()]
}
