package models

import "gin-web/global"

type Tag struct {
	Model
	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

func NewTag() *Tag {
	return &Tag{}
}

func (m *Tag) TableName() string {
	return "blog_tag"
}

func (m *Tag) GetListByPage(params interface{}, pageNo int, pageSize int) (total int, tags []Tag, err error) {
	offset := (pageNo - 1) * pageSize
	err = global.BlogDB.Model(m).Where(params).Offset(offset).Limit(pageSize).Find(&tags).Count(&total).Error
	if err != nil {
		return total, tags, err
	}

	return total, tags, nil
}

func (m *Tag) FindById(id int) {
	global.BlogDB.Model(m).Where("id = ?", id).First(m)
}

func (m *Tag) FindByWhere(params interface{}) bool {
	global.BlogDB.Model(m).Where(params).First(m)
	if m.ID > 0 {
		return true
	}
	return false
}

func (m *Tag) Create() (int, error) {
	err := global.BlogDB.Create(m).Error
	if err != nil {
		return 0, err
	}

	return m.ID, nil
}

func (m *Tag) Update() error {
	return global.BlogDB.Model(m).Save(m).Error
}

func (m *Tag) DeleteTag(id int) bool {
	global.BlogDB.Where("id = ?", id).Delete(&Tag{})
	return true
}