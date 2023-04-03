package tag

type AddTag struct {
	Name      string `form:"name" json:"name" binding:"required"`
	State     int    `form:"state"  json:"state"  binding:"gte=0"`
	CreatedBy string `form:"created_by"  json:"created_by"  binding:"required"`
}
