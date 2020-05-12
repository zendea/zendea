package model

// 友链
type Link struct {
	Model
	Url        string `gorm:"not null;type:text" json:"url" form:"url"`     // 链接
	Title      string `gorm:"not null;size:128" json:"title" form:"title"`  // 标题
	Summary    string `gorm:"size:1024" json:"summary" form:"summary"`      // 站点描述
	Logo       string `gorm:"type:text" json:"logo" form:"logo"`            // LOGO
	Status     int    `gorm:"not null" json:"status" form:"status"`         // 状态
	CreateTime int64  `gorm:"not null" json:"createTime" form:"createTime"` // 创建时间
}
