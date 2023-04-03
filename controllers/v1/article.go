package v1

import (
	"gin-web/core/ginc"
	"gin-web/service"
	"gin-web/validators/article"
	"github.com/gin-gonic/gin"
)

// GetArticle 获取单个文章
func GetArticle(c *gin.Context) {
	id := ginc.GetInt(c, "id")
	if id <= 0 {
		ginc.Fail(c, "ID必须大于0")
		return
	}
	articleDetail, err := service.NewArticleService().GetArticle(id)
	if err != nil {
		ginc.Fail(c, "文章信息不存在")
		return
	}
	ginc.Ok(c, articleDetail)
}

// GetArticles 获取多个文章
func GetArticles(c *gin.Context) {
	params := map[string]interface{}{
		"state":  ginc.GetInt(c, "state", -1),
		"tag_id": ginc.GetInt(c, "tag_id", -1),
	}
	total, list, err := service.NewArticleService().GetListByPage(params, ginc.GetPage(c), ginc.GetPageSize(c))
	if err != nil {
		ginc.Fail(c, "获取文章列表失败")
		return
	}
	data := map[string]interface{}{
		"total": total,
		"list":  list,
	}
	ginc.Ok(c, data)
}

// AddArticle 新增文章
func AddArticle(c *gin.Context) {
	var articleData article.AddArticle
	if err := c.ShouldBindJSON(&articleData); err != nil {
		ginc.Fail(c, err.Error())
		return
	}
	id, err := service.NewArticleService().AddArticle(articleData)
	if err != nil {
		ginc.Fail(c, err.Error())
		return
	}
	ginc.Ok(c, map[string]interface{}{"id": id})
}

//修改文章
func EditArticle(c *gin.Context) {
}

//删除文章
func DeleteArticle(c *gin.Context) {
	id := ginc.GetInt(c, "id")
	if id <= 0 {
		ginc.Fail(c, "ID必须大于0")
		return
	}
	err := service.NewArticleService().Delete(id)
	if err != nil {
		ginc.Fail(c, err.Error())
		return
	}
	ginc.Ok(c, map[string]interface{}{})
}
