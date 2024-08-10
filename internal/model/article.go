package model

type Article struct {
	Model
	TagID      int64  `gorm:"column:tag_id" gorm:"index"`
	Title      string `gorm:"column:title"`
	Desc       string `gorm:"column:desc"`
	Content    string `gorm:"column:content"`
	State      int8   `gorm:"column:state"`
	CreateUser string `gorm:"column:create_user"`
	UpdateUser string `gorm:"column:update_user"`
	Tag        *Tag
}

func (m *Article) TableName() string {
	return "blog_article"
}
