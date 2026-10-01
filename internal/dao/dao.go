package dao

import (
	"gin-web/internal/dao/article"
	"gin-web/internal/dao/category"
	"gin-web/internal/dao/comment"
	"gin-web/internal/dao/link"
	"gin-web/internal/dao/option"
	"gin-web/internal/dao/tag"
	"gin-web/internal/dao/user"
	"github.com/google/wire"
)

type Daos struct {
	Article  article.ArticleDao
	Category category.CategoryDao
	Comment  comment.CommentDao
	Link     link.LinkDao
	Option   option.OptionDao
	Tag      tag.TagDao
	User     user.UserDao
}

var ProviderSet = wire.NewSet(
	article.NewArticleDao,
	category.NewCategoryDao,
	comment.NewCommentDao,
	link.NewLinkDao,
	option.NewOptionDao,
	tag.NewTagDao,
	user.NewUserDao,
)
