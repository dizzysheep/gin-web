package handler

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gin-web/app/response"
	"gin-web/core/config"
	"gin-web/dto"
	"gin-web/internal/service"
	"github.com/gin-gonic/gin"
)

var errInvalidSlug = errors.New("不合法的 slug")

type ArticleHandler struct {
	service *service.Services
}

func NewArticleHandler(service *service.Services) *ArticleHandler {
	return &ArticleHandler{service}
}

// List godoc
// @Summary 公开-文章列表
// @Description 仅返回已发布且启用的文章，支持分类/标签/关键词筛选
// @Tags 文章
// @Produce json
// @Param page query int false "页码(从1起，默认1)"
// @Param page_size query int false "每页数量(默认10，最大100)"
// @Param category_id query int false "分类ID"
// @Param tag_id query int false "标签ID"
// @Param keyword query string false "关键词(标题/简述/内容)"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/article [get]
func (h *ArticleHandler) List(c *gin.Context) {
	reqDTO, err := dto.ListArticleReqToDTO(c)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	respDTO, err := h.service.Article.List(c.Request.Context(), reqDTO)
	if err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, respDTO.ToVO())
}

// Detail godoc
// @Summary 公开-文章详情
// @Description 按ID查询已发布文章，并累计浏览量
// @Tags 文章
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/article/{id} [get]
func (h *ArticleHandler) Detail(c *gin.Context) {
	reqDTO, err := dto.IDReqDTOFromRequest(c)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	respDTO, err := h.service.Article.Detail(c.Request.Context(), reqDTO)
	if err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, respDTO.ToVO())
}

// DetailBySlug godoc
// @Summary 公开-按slug查询文章详情
// @Tags 文章
// @Produce json
// @Param slug path string true "文章slug"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/article/slug/{slug} [get]
func (h *ArticleHandler) DetailBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		response.BadRequest(c, errInvalidSlug)
		return
	}

	respDTO, err := h.service.Article.DetailBySlug(c.Request.Context(), &dto.ArticleSlugReqDTO{Slug: slug})
	if err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, respDTO.ToVO())
}

// Archive godoc
// @Summary 公开-文章归档
// @Description 按年月分组统计已发布文章
// @Tags 文章
// @Produce json
// @Success 200 {object} response.Response
// @Router /v1/article/archive [get]
func (h *ArticleHandler) Archive(c *gin.Context) {
	respDTO, err := h.service.Article.Archive(c.Request.Context())
	if err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, respDTO.ToVO())
}

type rssFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Title         string    `xml:"title"`
	Link          string    `xml:"link"`
	Description   string    `xml:"description"`
	LastBuildDate string    `xml:"lastBuildDate"`
	Items         []rssItem `xml:"item"`
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	GUID        string `xml:"guid"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// RSS godoc
// @Summary RSS订阅
// @Description 输出所有已发布文章的 RSS 2.0 文档
// @Tags 文章
// @Produce xml
// @Success 200 {string} string "RSS 2.0 XML文档"
// @Router /rss.xml [get]
func (h *ArticleHandler) RSS(c *gin.Context) {
	articles, err := h.service.Article.AllPublished(c.Request.Context())
	if err != nil {
		response.FailErr(c, err)
		return
	}
	base := publicBaseURL(c)
	title := config.GetString("site.title")
	if title == "" {
		title = config.AppName
	}
	description := config.GetString("site.description")
	items := make([]rssItem, 0, len(articles.Items))
	for _, item := range articles.Items {
		if item == nil || item.PO == nil {
			continue
		}
		link := fmt.Sprintf("%s/api/v1/article/%d", base, item.PO.ID)
		items = append(items, rssItem{
			Title: item.PO.Title, Link: link, GUID: link,
			Description: item.PO.Desc, PubDate: time.Unix(int64(item.PO.PublishedOn), 0).UTC().Format(time.RFC1123Z),
		})
	}
	feed, err := xml.Marshal(rssFeed{
		Version: "2.0",
		Channel: rssChannel{Title: title, Link: base, Description: description, LastBuildDate: time.Now().UTC().Format(time.RFC1123Z), Items: items},
	})
	if err != nil {
		response.FailErr(c, err)
		return
	}
	c.Data(http.StatusOK, "application/rss+xml; charset=utf-8", append([]byte(xml.Header), feed...))
}

type sitemap struct {
	XMLName xml.Name       `xml:"urlset"`
	XMLNS   string         `xml:"xmlns,attr"`
	URLs    []sitemapEntry `xml:"url"`
}

type sitemapEntry struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod,omitempty"`
}

// Sitemap godoc
// @Summary Sitemap站点地图
// @Description 输出已发布文章的规范URL供爬虫抓取
// @Tags 文章
// @Produce xml
// @Success 200 {string} string "Sitemap XML文档"
// @Router /sitemap.xml [get]
func (h *ArticleHandler) Sitemap(c *gin.Context) {
	articles, err := h.service.Article.AllPublished(c.Request.Context())
	if err != nil {
		response.FailErr(c, err)
		return
	}
	base := publicBaseURL(c)
	urls := make([]sitemapEntry, 0, len(articles.Items))
	for _, item := range articles.Items {
		if item == nil || item.PO == nil {
			continue
		}
		urls = append(urls, sitemapEntry{
			Loc:     fmt.Sprintf("%s/api/v1/article/%d", base, item.PO.ID),
			LastMod: time.Unix(int64(item.PO.ModifiedOn), 0).UTC().Format("2006-01-02"),
		})
	}
	data, err := xml.Marshal(sitemap{XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9", URLs: urls})
	if err != nil {
		response.FailErr(c, err)
		return
	}
	c.Data(http.StatusOK, "application/xml; charset=utf-8", append([]byte(xml.Header), data...))
}

func publicBaseURL(c *gin.Context) string {
	if base := strings.TrimRight(config.GetString("site.baseURL"), "/"); base != "" {
		return base
	}
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + c.Request.Host
}

// AdminList godoc
// @Summary 管理-文章列表
// @Description 含草稿，支持状态/草稿/分类/标签/关键词筛选
// @Tags 文章
// @Produce json
// @Param page query int false "页码(从1起，默认1)"
// @Param page_size query int false "每页数量(默认10，最大100)"
// @Param category_id query int false "分类ID"
// @Param tag_id query int false "标签ID"
// @Param keyword query string false "关键词"
// @Param state query int false "状态 0禁用 1启用"
// @Param is_draft query int false "是否草稿 0否 1是"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/admin/article [get]
func (h *ArticleHandler) AdminList(c *gin.Context) {
	reqDTO, err := dto.AdminListArticleReqToDTO(c)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	respDTO, err := h.service.Article.AdminList(c.Request.Context(), reqDTO)
	if err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, respDTO.ToVO())
}

// AdminDetail godoc
// @Summary 管理-文章详情
// @Description 含草稿，不计浏览量
// @Tags 文章
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/admin/article/{id} [get]
func (h *ArticleHandler) AdminDetail(c *gin.Context) {
	reqDTO, err := dto.IDReqDTOFromRequest(c)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	respDTO, err := h.service.Article.AdminDetail(c.Request.Context(), reqDTO)
	if err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, respDTO.ToVO())
}

// Add godoc
// @Summary 管理-新建文章
// @Tags 文章
// @Accept json
// @Produce json
// @Param request body dto.AddArticleRequest true "文章内容"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/article [post]
func (h *ArticleHandler) Add(c *gin.Context) {
	reqDTO, err := dto.AddArticleReqToDTO(c)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	if err := h.service.Article.Add(c.Request.Context(), reqDTO); err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, "success")
}

// Edit godoc
// @Summary 管理-编辑文章
// @Tags 文章
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Param request body dto.EditArticleRequest true "文章内容"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/article/{id} [patch]
func (h *ArticleHandler) Edit(c *gin.Context) {
	reqDTO, err := dto.EditArticleReqToDTO(c)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	if err := h.service.Article.Edit(c.Request.Context(), reqDTO); err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, "success")
}

// Del godoc
// @Summary 管理-删除文章
// @Tags 文章
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/article/{id} [delete]
func (h *ArticleHandler) Del(c *gin.Context) {
	reqDTO, err := dto.IDReqDTOFromRequest(c)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	if err := h.service.Article.Del(c.Request.Context(), reqDTO); err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, "success")
}

// Publish godoc
// @Summary 管理-发布/下架文章
// @Tags 文章
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Param request body dto.PublishArticleRequest true "发布状态"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/article/{id}/publish [patch]
func (h *ArticleHandler) Publish(c *gin.Context) {
	reqDTO, err := dto.PublishArticleReqToDTO(c)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	if err := h.service.Article.Publish(c.Request.Context(), reqDTO); err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, "success")
}

// State godoc
// @Summary 管理-启用/禁用文章
// @Tags 文章
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Param request body dto.StateArticleRequest true "状态"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/article/{id}/state [patch]
func (h *ArticleHandler) State(c *gin.Context) {
	reqDTO, err := dto.StateArticleReqToDTO(c)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	if err := h.service.Article.State(c.Request.Context(), reqDTO); err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, "success")
}

// Top godoc
// @Summary 管理-置顶/取消置顶文章
// @Tags 文章
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Param request body dto.TopArticleRequest true "置顶状态"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/article/{id}/top [patch]
func (h *ArticleHandler) Top(c *gin.Context) {
	reqDTO, err := dto.TopArticleReqToDTO(c)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	if err := h.service.Article.Top(c.Request.Context(), reqDTO); err != nil {
		response.FailErr(c, err)
		return
	}

	response.Ok(c, "success")
}
