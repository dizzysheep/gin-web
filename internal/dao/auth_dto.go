package dao

import (
	"gin-web/internal/models"
	"gorm.io/gorm"
)

type AuthDao struct {
	*gorm.DB
}

func NewAuthDao(db *gorm.DB) *AuthDao {
	return &AuthDao{db}
}

func (dao *AuthDao) GetInfoByUserName(username string) (*models.Auth, error) {
	var auth models.Auth
	err := dao.Where("username = ?", username).First(&auth).Error
	if err != nil {
		return &auth, err
	}
	return &auth, nil
}
