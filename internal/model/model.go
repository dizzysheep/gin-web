package model

// 状态常量（state 字段通用取值）
const (
	StateDisabled int8 = 0
	StateEnabled  int8 = 1
)

// Model 通用字段（对应 doc/sql/db.sql 中 created_on/modified_on/deleted_on/state 约定）
// 软删除为手动过滤：查询统一追加 deleted_on = 0，删除为 UPDATE deleted_on
// CreatedOn/ModifiedOn 由 GORM 自动维护（uint32 映射为秒级 Unix 时间戳），业务层无需手动赋值
type Model struct {
	ID         int64  `gorm:"column:id;primaryKey" json:"id"`
	CreatedOn  uint32 `gorm:"column:created_on;autoCreateTime" json:"created_on"`
	ModifiedOn uint32 `gorm:"column:modified_on;autoUpdateTime" json:"modified_on"`
	DeletedOn  uint32 `gorm:"column:deleted_on" json:"-"`
	State      int8   `gorm:"column:state" json:"state"`
}

// AuditModel 带审计人字段的通用字段
type AuditModel struct {
	Model
	CreatedBy  string `gorm:"column:created_by" json:"created_by"`
	ModifiedBy string `gorm:"column:modified_by" json:"modified_by"`
}
