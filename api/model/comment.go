package model

// 评论
type Comment struct {
	Model
	UserId      int64  `gorm:"index:idx_comment_user_id;not null" json:"userId" form:"userId"`             // 用户编号
	EntityType  string `gorm:"index:idx_comment_entity_type;not null" json:"entityType" form:"entityType"` // 被评论实体类型
	EntityId    int64  `gorm:"index:idx_comment_entity_id;not null" json:"entityId" form:"entityId"`       // 被评论实体编号
	Content     string `gorm:"type:text;not null" json:"content" form:"content"`                           // 内容
	ContentType string `gorm:"type:varchar(32);not null" json:"contentType" form:"contentType"`            // 内容类型：markdown、html
	QuoteId     int64  `gorm:"not null"  json:"quoteId" form:"quoteId"`                                    // 引用的评论编号
	Status      int    `gorm:"int;index:idx_comment_status" json:"status" form:"status"`                   // 状态：0：待审核、1：审核通过、2：审核失败、3：已发布
	CreateTime  int64  `json:"createTime" form:"createTime"`                                               // 创建时间
}
