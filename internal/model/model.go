package model

type Model struct {
	ID        int64 `gorm:"column:id"`
	CreatedAt int64 `gorm:"column:create_time"`
	UpdatedAt int64 `gorm:"column:update_time"`
}
