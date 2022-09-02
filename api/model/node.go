package model

// 话题节点
type Node struct {
	Model
	Name        string `gorm:"size:32;unique" json:"name" form:"name"`                           // 名称
	Description string `json:"description" form:"description"`                                   // 描述
	SortNo      int    `gorm:"index:idx_sort_no" json:"sortNo" form:"sortNo"`                    // 排序编号
	Status      int    `gorm:"not null" json:"status" form:"status"`                             // 状态
	TopicCount  int64  `gorm:"not null" json:"topicCount" form:"topicCount"`                     // 主题数量
	CreateTime  int64  `json:"createTime" form:"createTime"`                                     // 创建时间
}
