package model

type Auth struct {
	ID       int64  `gorm:"primary_key" json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (m *Auth) TableName() string {
	return "blog_auth"
}
