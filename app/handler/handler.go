package handler

import (
	"github.com/google/wire"
)

type Handlers struct {
	Common  *CommonHandler
	Health  *HealthHandler
	Tag     *TagHandler
	Auth    *AuthHandler
	Article *ArticleHandler
}

var ProviderSet = wire.NewSet(
	NewHealthHandler,
	NewCommonHandler,
	NewTagHandler,
	NewAuthHandler,
	NewArticleHandler,
)
