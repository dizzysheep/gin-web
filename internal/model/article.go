package model

type Article struct {
	AuditModel
	CategoryID  int64  `gorm:"column:category_id" json:"category_id"`
	Title       string `gorm:"column:title" json:"title"`
	Slug        string `gorm:"column:slug" json:"slug"`
	Desc        string `gorm:"column:desc" json:"desc"`
	Cover       string `gorm:"column:cover" json:"cover"`
	ContentMd   string `gorm:"column:content_md" json:"content_md"`
	ContentHtml string `gorm:"column:content_html" json:"content_html"`
	ViewCount   int64  `gorm:"column:view_count" json:"view_count"`
	IsTop       int8   `gorm:"column:is_top" json:"is_top"`
	IsDraft     int8   `gorm:"column:is_draft" json:"is_draft"`
	PublishedOn uint32 `gorm:"column:published_on" json:"published_on"`
}

func (m *Article) TableName() string {
	return "blog_article"
}

// ArticleTag 文章-标签关联（物理删除）
type ArticleTag struct {
	ID        int64 `gorm:"column:id;primaryKey"`
	ArticleID int64 `gorm:"column:article_id"`
	TagID     int64 `gorm:"column:tag_id"`
}

func (m *ArticleTag) TableName() string {
	return "blog_article_tag"
}

// ArchiveRow 归档统计行（按月分组）
type ArchiveRow struct {
	Month string `gorm:"column:month" json:"month"`
	Total int64  `gorm:"column:total" json:"total"`
}
