package service

import (
	"gin-web/internal/service/article"
	"gin-web/internal/service/auth"
	"gin-web/internal/service/tag"
	"github.com/google/wire"
)

type Services struct {
	Article article.ArticleService
	Tag     tag.TagService
	Auth    auth.AuthService
}

var ProviderSet = wire.NewSet(
	article.NewArticleService,
	tag.NewTagService,
	auth.NewAuthService,
)
