package service

import (
	"errors"
	"gin-web/models"
	"gin-web/validators/tag"
)

type TagService struct {
	Model *models.Tag
}

func NewTagService() *TagService {
	return &TagService{
		Model: models.NewTag(),
	}
}

func (s *TagService) GetListByPage(params map[string]interface{}, pageNo, pageSize int) (total int, list []models.Tag, err error) {
	total, list, err = s.Model.GetListByPage(params, pageNo, pageSize)
	if err != nil {
		return total, list, err
	}
	return total, list, nil
}

func (s *TagService) AddTag(tag tag.AddTag) (int, error) {
	if s.Model.FindByWhere(map[string]interface{}{"name": tag.Name}) {
		return 0, errors.New("标签已存在")
	}
	s.Model.Name = tag.Name
	s.Model.State = tag.State
	s.Model.CreatedBy = tag.CreatedBy
	id, err := s.Model.Create()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *TagService) EditTag(tag tag.EditTag) error {
	s.Model.FindById(tag.Id)
	s.Model.Name = tag.Name
	s.Model.State = tag.State
	s.Model.ModifiedBy = tag.ModifiedBy
	if s.Model.ID <= 0 {
		return errors.New("标签不存在")
	}
	return s.Model.Update()
}
