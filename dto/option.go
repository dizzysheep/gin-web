package dto

import (
	"gin-web/app/ext"
	"gin-web/internal/model"
	"github.com/gin-gonic/gin"
	"time"
)

// OptionVO 站点配置项
type OptionVO struct {
	Key    string `json:"key" example:"site_title"`
	Value  string `json:"value" example:"gin-web 个人博客"`
	Remark string `json:"remark" example:"站点标题"`
}

// GetOptionsRespDTO 站点配置
type GetOptionsRespDTO struct {
	Options map[string]string
}

type GetOptionsResponse struct {
	Options map[string]string `json:"options"`
}

func (g *GetOptionsRespDTO) ToVO() *GetOptionsResponse {
	if g.Options == nil {
		g.Options = map[string]string{}
	}
	return &GetOptionsResponse{Options: g.Options}
}

// SaveOptionsRequest 批量保存站点配置
type SaveOptionsRequest struct {
	Options []OptionItem `json:"options" binding:"required,min=1,max=100,dive"`
}

type OptionItem struct {
	Key   string `json:"key" binding:"required,max=100" example:"site_title"`
	Value string `json:"value" binding:"max=2000" example:"gin-web 个人博客"`
}

func SaveOptionsReqToDTO(c *gin.Context) (*SaveOptionsReqDTO, error) {
	var req SaveOptionsRequest
	if err := c.ShouldBind(&req); err != nil {
		return nil, err
	}

	items := make([]*model.Option, 0, len(req.Options))
	now := uint32(time.Now().Unix())
	for _, opt := range req.Options {
		items = append(items, &model.Option{
			Key:        opt.Key,
			Value:      opt.Value,
			ModifiedOn: now,
			ModifiedBy: ext.GetUsername(c),
			State:      model.StateEnabled,
		})
	}

	return &SaveOptionsReqDTO{Options: items}, nil
}

type SaveOptionsReqDTO struct {
	Options []*model.Option
}
