package model

import (
	"database/sql"
)

type User struct {
	Model
	Username     sql.NullString `gorm:"size:32;unique;" json:"username" form:"username"`            // 用户名
	Email        sql.NullString `gorm:"size:128;unique;" json:"email" form:"email"`                 // 邮箱
	Nickname     string         `gorm:"size:16;" json:"nickname" form:"nickname"`                   // 昵称
	Avatar       string         `gorm:"type:text" json:"avatar" form:"avatar"`                      // 头像
	Password     string         `gorm:"size:512" json:"password" form:"password"`                   // 密码
	Website      string         `gorm:"size:1024" json:"website" form:"website"`                    // 个人主页
	Description  string         `gorm:"type:text" json:"description" form:"description"`            // 个人描述
	Status       int            `gorm:"index:idx_user_status;not null" json:"status" form:"status"` // 状态
	TopicCount   int            `gorm:"not null" json:"topicCount" form:"topicCount"`               // 帖子数量
	CommentCount int            `gorm:"not null" json:"commentCount" form:"commentCount"`           // 跟帖数量
	Level        int            `gorm:"not null" json:"level" form:"level"`                         // 用户等级
	CreateTime   int64          `json:"createTime" form:"createTime"`                               // 创建时间
	UpdateTime   int64          `json:"updateTime" form:"updateTime"`                               // 更新时间
}

// 用户积分
type UserScore struct {
	Model
	UserId     int64 `gorm:"unique;not null" json:"userId" form:"userId"` // 用户编号
	Score      int   `gorm:"not null" json:"score" form:"score"`          // 积分
	CreateTime int64 `json:"createTime" form:"createTime"`                // 创建时间
	UpdateTime int64 `json:"updateTime" form:"updateTime"`                // 更新时间
}

// 用户积分流水
type UserScoreLog struct {
	Model
	UserId      int64  `gorm:"not null;index:idx_user_score_log_user_id" json:"userId" form:"userId"`   // 用户编号
	SourceType  string `gorm:"not null;index:idx_user_score_score" json:"sourceType" form:"sourceType"` // 积分来源类型
	SourceId    string `gorm:"not null;index:idx_user_score_score" json:"sourceId" form:"sourceId"`     // 积分来源编号
	Description string `json:"description" form:"description"`                                          // 描述
	Type        int    `json:"type" form:"type"`                                                        // 类型(增加、减少)
	Score       int    `json:"score" form:"score"`                                                      // 积分
	CreateTime  int64  `json:"createTime" form:"createTime"`                                            // 创建时间
}
