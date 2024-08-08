package common

import "gorm.io/gorm"

type Pagination struct {
	Offset   int
	PageSize int
}

func (p *Pagination) Paginate() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(p.Offset).Limit(p.PageSize)
	}
}
