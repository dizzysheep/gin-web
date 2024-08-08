package model

type Tag struct {
	Model
	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int8   `json:"state"`
}

func (m *Tag) TableName() string {
	return "blog_tag"
}
