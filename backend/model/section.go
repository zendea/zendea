package model

// 话题节点
type Section struct {
	Model
	Name       string `gorm:"size:32;unique" json:"name" form:"name"`        // 名称
	SortNo     int    `gorm:"index:idx_sort_no" json:"sortNo" form:"sortNo"` // 排序编号
	CreateTime int64  `json:"createTime" form:"createTime"`                  // 创建时间
}
