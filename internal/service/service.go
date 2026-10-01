package service

import (
	"gin-web/internal/service/article"
	"gin-web/internal/service/auth"
	"gin-web/internal/service/category"
	"gin-web/internal/service/comment"
	"gin-web/internal/service/link"
	"gin-web/internal/service/option"
	"gin-web/internal/service/tag"
	"github.com/google/wire"
)

type Services struct {
	Article  article.ArticleService
	Auth     auth.AuthService
	Category category.CategoryService
	Comment  comment.CommentService
	Link     link.LinkService
	Option   option.OptionService
	Tag      tag.TagService
}

var ProviderSet = wire.NewSet(
	article.NewArticleService,
	auth.NewAuthService,
	category.NewCategoryService,
	comment.NewCommentService,
	link.NewLinkService,
	option.NewOptionService,
	tag.NewTagService,
)
