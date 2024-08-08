package dao

import (
	"gin-web/internal/dao/article"
	"gin-web/internal/dao/auth"
	"gin-web/internal/dao/tag"
	"github.com/google/wire"
)

type Daos struct {
	Article article.ArticleDao
	Auth    auth.AuthDao
	Tag     tag.TagDao
}

var ProviderSet = wire.NewSet(
	article.NewArticleDao,
	tag.NewTagDao,
	auth.NewAuthDao,
)
