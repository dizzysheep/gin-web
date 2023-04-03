package article

type AddArticle struct {
	TagId     int    `json:"tag_id" binding:"required,gt=0"`
	Title     string `json:"title" binding:"required"`
	Desc      string `json:"desc" binding:"required"`
	Content   string `json:"content" binding:"required"`
	CreatedBy string `json:"created_by" binding:"required"`
	State     int    `json:"state" binding:"gte=0"`
}
