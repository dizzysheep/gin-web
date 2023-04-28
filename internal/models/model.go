package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

type Model struct {
	ID        int       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `gorm:"create_time"`
	UpdatedAt time.Time `gorm:"update_time"`
}
