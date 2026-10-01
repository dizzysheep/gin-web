package model

type Tag struct {
	AuditModel
	Name string `gorm:"column:name" json:"name"`
}

func (m *Tag) TableName() string {
	return "blog_tag"
}
