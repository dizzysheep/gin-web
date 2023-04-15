package dao

import (
	"gin-web/global"
	"gin-web/internal/dto"
	"gin-web/internal/models"
	"github.com/jinzhu/gorm"
)

type TagDao struct {
	*gorm.DB
}

func NewTagDao() *TagDao {
	return &TagDao{global.BlogDB}
}

func (dao *TagDao) CreateTag(tag *models.Tag) (int, error) {
	err := dao.Create(tag).Error
	if err != nil {
		return 0, err
	}

	return tag.ID, nil
}

func (dao *TagDao) ExistTagByName(name string) bool {
	tag := &models.Tag{}
	err := dao.Where("name = ?", name).First(tag).Error
	if err != nil {
		return false
	}

	return true
}

func (dao *TagDao) GetTagByPage(req dto.SearchTagReqDTO, pageNo, pageSize int) (total int, list []models.Tag, err error) {
	dao.SetCondition(req)
	offset := (pageNo - 1) * pageSize
	err = dao.Find(&list).Limit(pageSize).Offset(offset).Count(&total).Error
	if err != nil {
		return
	}

	return
}

// SetCondition ---设置查询条件
func (dao *TagDao) SetCondition(req dto.SearchTagReqDTO) {
	if req.Name != "" {
		dao.DB = dao.Where("name = ?", req.Name)
	}

	if req.State != nil {
		dao.DB = dao.Where("state = ?", req.State)
	}
}
