package handler

import (
	"github.com/google/wire"
)

type Handlers struct {
	Common   *CommonHandler
	Health   *HealthHandler
	Auth     *AuthHandler
	Article  *ArticleHandler
	Category *CategoryHandler
	Tag      *TagHandler
	Comment  *CommentHandler
	Link     *LinkHandler
	Upload   *UploadHandler
}

var ProviderSet = wire.NewSet(
	NewHealthHandler,
	NewCommonHandler,
	NewAuthHandler,
	NewArticleHandler,
	NewCategoryHandler,
	NewTagHandler,
	NewCommentHandler,
	NewLinkHandler,
	NewUploadHandler,
)
