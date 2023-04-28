package v1

import (
	"gin-web/core/ginc"
	"gin-web/internal/dto"
	"gin-web/internal/service"
	"github.com/gin-gonic/gin"
)

type ArticleController struct {
}

func NewArticleController() *ArticleController {
	return &ArticleController{}
}

func (ac *ArticleController) Router(router *gin.RouterGroup) {
	router.GET("/article/:id", ac.GetArticle)          //获取单个文章
	router.GET("/articles", ac.GetArticles)            //获取多个文章
	router.POST("/article/create", ac.CreateArticle)   //新建文章
	router.POST("/article/edit", ac.EditArticle)       //更新文章
	router.DELETE("/article/delete", ac.DeleteArticle) //删除文章
}

// GetArticle 获取单个文章
func (ac *ArticleController) GetArticle(c *gin.Context) {
	id := ginc.GetInt(c, "id")
	if id <= 0 {
		ginc.Fail(c, "ID必须大于0")
		return
	}
	articleDetail, err := service.NewArticleService(c).GetArticle(id)
	if err != nil {
		ginc.Fail(c, "文章信息不存在")
		return
	}
	ginc.Ok(c, articleDetail)
}

// GetArticles 获取多个文章
func (ac *ArticleController) GetArticles(c *gin.Context) {
	var req dto.SearchArticleReqDTO
	if err := c.ShouldBind(&req); err != nil {
		ginc.Fail(c, err.Error())
		return
	}

	page, size := ginc.GetPage(c), ginc.GetPageSize(c)
	total, list, err := service.NewArticleService(c).GetListByPage(req, page, size)
	if err != nil {
		ginc.Fail(c, "获取文章列表失败")
		return
	}

	pageInfo := &ginc.PageInfo{
		Total: total,
		List:  list,
		Page:  page,
		Size:  size,
	}
	ginc.Ok(c, pageInfo)
}

// CreateArticle 新增文章
func (ac *ArticleController) CreateArticle(c *gin.Context) {
	var req dto.CreateArticleReqDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		ginc.Fail(c, err.Error())
		return
	}
	id, err := service.NewArticleService(c).CreateArticle(req)
	if err != nil {
		ginc.Fail(c, err.Error())
		return
	}
	ginc.Ok(c, map[string]interface{}{"id": id})
}

//修改文章
func (ac *ArticleController) EditArticle(c *gin.Context) {
}

// DeleteArticle ---删除文章
func (ac *ArticleController) DeleteArticle(c *gin.Context) {
	id := ginc.GetInt(c, "id")
	if id <= 0 {
		ginc.Fail(c, "ID必须大于0")
		return
	}

	err := service.NewArticleService(c).DeleteArticle(id)
	if err != nil {
		ginc.Fail(c, err.Error())
		return
	}
	ginc.Ok(c, map[string]interface{}{})
}
