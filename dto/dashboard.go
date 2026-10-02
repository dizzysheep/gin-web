package dto

// DashboardOverviewRespDTO contains the counters shown in the admin dashboard.
type DashboardOverviewRespDTO struct {
	ArticleCount        int64 `json:"article_count"`
	CategoryCount       int64 `json:"category_count"`
	TagCount            int64 `json:"tag_count"`
	PendingCommentCount int64 `json:"pending_comment_count"`
}

type DashboardTrendItem struct {
	Date  string `json:"date"`
	Label string `json:"label"`
	Total int64  `json:"total"`
}

// DashboardTrendRespDTO always returns one item for each requested day,
// including days with no published articles.
type DashboardTrendRespDTO struct {
	Days  int                   `json:"days"`
	Items []*DashboardTrendItem `json:"items"`
}
