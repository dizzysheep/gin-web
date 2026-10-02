package model

type Option struct {
	ID         int64  `gorm:"column:id;primaryKey"`
	Key        string `gorm:"column:key" json:"key"`
	Value      string `gorm:"column:value" json:"value"`
	Remark     string `gorm:"column:remark" json:"remark"`
	ModifiedOn uint32 `gorm:"column:modified_on;autoUpdateTime" json:"modified_on"`
	ModifiedBy string `gorm:"column:modified_by" json:"modified_by"`
	DeletedOn  uint32 `gorm:"column:deleted_on" json:"-"`
	State      int8   `gorm:"column:state" json:"state"`
}

func (m *Option) TableName() string {
	return "blog_option"
}
