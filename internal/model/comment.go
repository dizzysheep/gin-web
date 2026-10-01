package model

// 评论状态
const (
	CommentStatePending int8 = 0 // 待审核
	CommentStatePass    int8 = 1 // 通过
	CommentStateReject  int8 = 2 // 拒绝
)

type Comment struct {
	AuditModel
	ArticleID int64  `gorm:"column:article_id" json:"article_id"`
	ParentID  int64  `gorm:"column:parent_id" json:"parent_id"`
	Nickname  string `gorm:"column:nickname" json:"nickname"`
	Email     string `gorm:"column:email" json:"-"`
	Content   string `gorm:"column:content" json:"content"`
	IsAdmin   int8   `gorm:"column:is_admin" json:"is_admin"`
	IP        string `gorm:"column:ip" json:"-"`
}

func (m *Comment) TableName() string {
	return "blog_comment"
}
