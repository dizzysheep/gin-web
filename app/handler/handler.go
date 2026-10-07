package handler

import (
	"github.com/google/wire"
)

type Handlers struct {
	Agent     *AgentHandler
	Common    *CommonHandler
	Health    *HealthHandler
	Auth      *AuthHandler
	Article   *ArticleHandler
	Category  *CategoryHandler
	Tag       *TagHandler
	Comment   *CommentHandler
	Dashboard *DashboardHandler
	Link      *LinkHandler
	Server    *ServerHandler
	Upload    *UploadHandler
}

var ProviderSet = wire.NewSet(
	NewAgentHandler,
	NewHealthHandler,
	NewCommonHandler,
	NewAuthHandler,
	NewArticleHandler,
	NewCategoryHandler,
	NewTagHandler,
	NewCommentHandler,
	NewDashboardHandler,
	NewLinkHandler,
	NewServerHandler,
	NewUploadHandler,
)
