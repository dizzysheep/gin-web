package dao

import (
	"gin-web/global"
	"gin-web/internal/models"
	"github.com/jinzhu/gorm"
)

type AuthDao struct {
	*gorm.DB
}

func NewAuthDao() *AuthDao {
	return &AuthDao{global.BlogDB}
}

func (dao *AuthDao) GetInfoByUserName(username string) (*models.Auth, error) {
	var auth models.Auth
	err := dao.Where("username = ?", username).First(&auth).Error
	if err != nil {
		return &auth, err
	}
	return &auth, nil
}
