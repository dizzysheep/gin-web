package models

type Article struct {
	Model
	TagID      int    `json:"tag_id" gorm:"index"`
	Tag        Tag    `json:"tag"`
	Title      string `json:"title"`
	Desc       string `json:"desc"`
	Content    string `json:"content"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

func NewArticle() *Article {
	return &Article{}
}

func (m *Article) TableName() string {
	return "blog_article"
}

func (m *Article) Create() (int, error) {
	err := db.Create(m).Error
	if err != nil {
		return 0, err
	}

	return m.ID, nil
}

func (m *Article) GetArticle(id int) (*Article, error) {
	articleOne := Article{}
	db.Where("id = ?", id).First(&articleOne)
	db.Model(articleOne).Related(&articleOne.Tag)
	return &articleOne, nil
}

func (m *Article) Update() error {
	return db.Save(m).Error
}

func (m *Article) Delete(id int) error {
	return db.Where("id = ?", id).Delete(m).Error
}

func (m *Article) GetListByPage(params interface{}, pageNo int, pageSize int) (total int, articles []Article, err error) {
	offset := (pageNo - 1) * pageSize
	err = db.Preload("Tag").Where(params).Offset(offset).Limit(pageSize).Find(&articles).Count(&total).Error
	if err != nil {
		return total, articles, err
	}

	return total, articles, nil
}
