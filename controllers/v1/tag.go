package v1

import (
	"gin-web/core/ginc"
	"gin-web/models"
	"gin-web/service"
	"gin-web/validators/tag"
	"github.com/gin-gonic/gin"
)

// GetTags 获取多个文章标签
func GetTags(c *gin.Context) {
	name := ginc.GetString(c, "name")
	state := ginc.GetInt(c, "state", -1)
	maps := make(map[string]interface{})
	if name != "" {
		maps["name"] = name
	}
	if state != -1 {
		maps["state"] = state
	}
	total, list, err := service.NewTagService().GetListByPage(maps, ginc.GetPage(c), ginc.GetPageSize(c))
	if err != nil {
		ginc.Fail(c, err.Error())
		return
	}

	data := map[string]interface{}{"total": total, "list": list}
	ginc.Ok(c, data)
}

// AddTag
// @Summary 新增文章标签
// @Produce  json
// @Param name query string true "Name"
// @Param state query int false "State"
// @Param created_by query int false "CreatedBy"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags [post]
func AddTag(c *gin.Context) {
	var tagData tag.AddTag
	if err := c.ShouldBind(&tagData); err != nil {
		ginc.Fail(c, err.Error())
		return
	}
	id, err := service.NewTagService().AddTag(tagData)
	if err != nil {
		ginc.Fail(c, err.Error())
		return
	}
	data := map[string]interface{}{"id": id}
	ginc.Ok(c, data)
}

// EditTag 修改文章标签
func EditTag(c *gin.Context) {
	var tagData tag.EditTag
	if err := c.ShouldBindJSON(&tagData); err != nil {
		ginc.Fail(c, err.Error())
		return
	}
	err := service.NewTagService().EditTag(tagData)
	if err != nil {
		ginc.Fail(c, err.Error())
		return
	}
	ginc.Ok(c, nil)
}

// DeleteTag 删除文章标签
func DeleteTag(c *gin.Context) {
	id := ginc.GetInt(c, "id")
	if id == 0 {
		ginc.Fail(c, "id参数不合法")
		return
	}
	models.NewTag().DeleteTag(id)
	ginc.Ok(c, nil)
}
