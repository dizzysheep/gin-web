package tag

type EditTag struct {
	Id         int    `json:"id" binding:"required"`
	Name       string `json:"name" binding:"required"`
	State      int    `json:"state"  binding:"gte=0"`
	ModifiedBy string `json:"modified_by"  binding:"required"`
}
