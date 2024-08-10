package model

type Tag struct {
	Model
	Name       string `gorm:"column:name"`
	State      int8   `gorm:"column:state"`
	CreateUser string `gorm:"column:create_user"`
	UpdateUser string `gorm:"column:update_user"`
}

func (m *Tag) TableName() string {
	return "blog_tag"
}
