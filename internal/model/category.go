package model

type Category struct {
	AuditModel
	Name     string `gorm:"column:name" json:"name"`
	ParentID int64  `gorm:"column:parent_id" json:"parent_id"`
	Sort     int64  `gorm:"column:sort" json:"sort"`
}

func (m *Category) TableName() string {
	return "blog_category"
}
