package models

type Auth struct {
	ID       int    `gorm:"primary_key" json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (m *Auth) TableName() string {
	return "blog_auth"
}
