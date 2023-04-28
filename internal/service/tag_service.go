package service

import (
	"errors"
	"gin-web/internal/dao"
	"gin-web/internal/dto"
	"gin-web/internal/models"
	"github.com/gin-gonic/gin"
)

type TagService struct {
	dao *dao.TagDao
}

func NewTagService(ctx *gin.Context) *TagService {
	return &TagService{dao.NewTagDao(GetBlogDB(ctx))}
}

func (svc *TagService) GetListByPage(req dto.SearchTagReqDTO, pageNo, pageSize int) (total int64, list []models.Tag, err error) {
	return svc.dao.GetTagByPage(req, pageNo, pageSize)
}

// CreateTag ---创建标签tag---
func (svc *TagService) CreateTag(req dto.CreateTagReqDTO) (int, error) {
	if svc.dao.ExistTagByName(req.Name) {
		return 0, errors.New("标签已存在")
	}

	return svc.dao.CreateTag(&models.Tag{
		Name:      req.Name,
		CreatedBy: req.CreatedBy,
		State:     req.State,
	})
}

func (svc *TagService) UpdateTag(tag dto.UpdateTagReqDTO) error {
	return svc.dao.UpdateInfoById(tag.Id, map[string]interface{}{
		"name":        tag.Name,
		"state":       tag.State,
		"modified_by": tag.ModifiedBy,
	})
}

func (svc *TagService) DeleteTag(id int) error {
	return svc.dao.Delete(id).Error

}
