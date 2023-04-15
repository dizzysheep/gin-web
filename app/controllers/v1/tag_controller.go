package v1

import (
	"gin-web/core/ginc"
	"gin-web/core/result"
	"gin-web/internal/dto"
	"gin-web/internal/models"
	"gin-web/internal/service"
	"github.com/gin-gonic/gin"
)

type TagController struct {
}

func NewTagController() *TagController {
	return &TagController{}
}

func (tc *TagController) Router(router *gin.RouterGroup) {
	router.GET("/tags", tc.GetTags)             //获取标签列表
	router.POST("/tags/add", tc.CreateTag)      //新建标签
	router.POST("/tags/edit", tc.EditTag)       //更新指定标签
	router.DELETE("/tags/delete", tc.DeleteTag) //删除文章
}

// GetTags 获取多个文章标签
func (tc *TagController) GetTags(c *gin.Context) {
	var req dto.SearchTagReqDTO
	if err := c.ShouldBind(&req); err != nil {
		ginc.Fail(c, err.Error())
		return
	}

	page, size := ginc.GetPage(c), ginc.GetPageSize(c)
	total, list, err := service.NewTagService().GetListByPage(req, page, size)
	if err != nil {
		ginc.Fail(c, err.Error())
		return
	}

	pageInfo := &result.PageInfo{
		Total: total,
		List:  list,
		Page:  page,
		Size:  size,
	}
	ginc.Page(c, pageInfo)
}

// CreateTag
// @Summary 新增文章标签
// @Produce  json
// @Param name query string true "Name"
// @Param state query int false "State"
// @Param created_by query int false "CreatedBy"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags [post]
func (tc *TagController) CreateTag(c *gin.Context) {
	var req dto.CreateTagReqDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		ginc.Fail(c, err.Error())
		return
	}
	id, err := service.NewTagService().CreateTag(req)
	if err != nil {
		ginc.Fail(c, err.Error())
		return
	}
	data := map[string]interface{}{"id": id}
	ginc.Ok(c, data)
}

// EditTag 修改文章标签
func (tc *TagController) EditTag(c *gin.Context) {
	var req dto.EditTagReqDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		ginc.Fail(c, err.Error())
		return
	}
	err := service.NewTagService().EditTag(req)
	if err != nil {
		ginc.Fail(c, err.Error())
		return
	}
	ginc.Ok(c, nil)
}

// DeleteTag 删除文章标签
func (tc *TagController) DeleteTag(c *gin.Context) {
	id := ginc.GetInt(c, "id")
	if id == 0 {
		ginc.Fail(c, "id参数不合法")
		return
	}
	models.NewTag().DeleteTag(id)
	ginc.Ok(c, nil)
}
