package model

type User struct {
	ID        int64  `gorm:"column:id;primaryKey" json:"id"`
	Username  string `gorm:"column:username" json:"username"`
	Password  string `gorm:"column:password" json:"-"`
	Nickname  string `gorm:"column:nickname" json:"nickname"`
	Avatar    string `gorm:"column:avatar" json:"avatar"`
	Email     string `gorm:"column:email" json:"email"`
	Role      int8   `gorm:"column:role" json:"role"`
	LastLogin uint32 `gorm:"column:last_login" json:"last_login"`
}

func (m *User) TableName() string {
	return "blog_user"
}
