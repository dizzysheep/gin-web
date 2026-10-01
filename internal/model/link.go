package model

type Link struct {
	AuditModel
	Name        string `gorm:"column:name" json:"name"`
	URL         string `gorm:"column:url" json:"url"`
	Logo        string `gorm:"column:logo" json:"logo"`
	Description string `gorm:"column:description" json:"description"`
	Sort        int64  `gorm:"column:sort" json:"sort"`
}

func (m *Link) TableName() string {
	return "blog_link"
}
