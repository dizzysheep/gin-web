package service

import (
	"gin-web/internal/agent"
	"gin-web/internal/service/article"
	"gin-web/internal/service/auth"
	"gin-web/internal/service/category"
	"gin-web/internal/service/comment"
	"gin-web/internal/service/dashboard"
	"gin-web/internal/service/link"
	"gin-web/internal/service/option"
	"gin-web/internal/service/server"
	"gin-web/internal/service/tag"
	"github.com/google/wire"
)

type Services struct {
	Agent     *agent.Service
	Article   article.ArticleService
	Auth      auth.AuthService
	Category  category.CategoryService
	Comment   comment.CommentService
	Dashboard dashboard.DashboardService
	Link      link.LinkService
	Option    option.OptionService
	Server    server.ServerService
	Tag       tag.TagService
}

var ProviderSet = wire.NewSet(
	agent.NewServiceFromConfig,
	article.NewArticleService,
	auth.NewAuthService,
	category.NewCategoryService,
	comment.NewCommentService,
	dashboard.NewDashboardService,
	link.NewLinkService,
	option.NewOptionService,
	server.NewServerService,
	tag.NewTagService,
)
