package dto

type CreateTagReqDTO struct {
	Name      string `json:"name" binding:"required"`
	State     int    `json:"state"  binding:"gte=0"`
	CreatedBy string `json:"created_by"  binding:"required"`
}

type UpdateTagReqDTO struct {
	Id         int    `json:"id" binding:"required"`
	Name       string `json:"name" binding:"required"`
	State      int    `json:"state"  binding:"gte=0"`
	ModifiedBy string `json:"modified_by"  binding:"required"`
}

type SearchTagReqDTO struct {
	Name  string `form:"name" `
	State *int   `form:"state"`
}
