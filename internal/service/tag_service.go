package service

import (
	"errors"
	"gin-web/internal/dao"
	"gin-web/internal/dto"
	"gin-web/internal/models"
)

type TagService struct {
	Model *models.Tag
}

func NewTagService() *TagService {
	return &TagService{
		Model: models.NewTag(),
	}
}

func (s *TagService) GetListByPage(req dto.SearchTagReqDTO, pageNo, pageSize int) (total int, list []models.Tag, err error) {
	return dao.NewTagDao().GetTagByPage(req, pageNo, pageSize)
}

// CreateTag ---创建标签tag---
func (s *TagService) CreateTag(req dto.CreateTagReqDTO) (int, error) {
	tagDao := dao.NewTagDao()
	if tagDao.ExistTagByName(req.Name) {
		return 0, errors.New("标签已存在")
	}

	return tagDao.CreateTag(&models.Tag{
		Name:      req.Name,
		CreatedBy: req.CreatedBy,
		State:     req.State,
	})
}

func (s *TagService) EditTag(tag dto.EditTagReqDTO) error {
	s.Model.FindById(tag.Id)
	s.Model.Name = tag.Name
	s.Model.State = tag.State
	s.Model.ModifiedBy = tag.ModifiedBy
	if s.Model.ID <= 0 {
		return errors.New("标签不存在")
	}
	return s.Model.Update()
}
