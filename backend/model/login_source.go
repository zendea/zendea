package model

import (
	"database/sql"
)

type LoginSource struct {
  Model
  UserID     sql.NullInt64 `gorm:"unique_index:idx_user_id_target_type;" json:"userId" form:"userId"`                                  // 用户编号
  Avatar     string        `gorm:"size:1024" json:"avatar" form:"avatar"`                                                             // 头像
  Nickname   string        `gorm:"size:32" json:"nickname" form:"nickname"`                                                           // 昵称
  TargetType string        `gorm:"size:32;not null;unique_index:idx_user_id_target_type,idx_target;" json:"targetType" form:"targetType"` // 第三方类型
  TargetID   string        `gorm:"size:64;not null;unique_index:idx_target;" json:"targetId" form:"targetId"`                            // 第三方唯一标识，例如：openId,unionId
  ExtraData  string        `gorm:"type:longtext" json:"extraData" form:"extraData"`                                                   // 扩展数据
  CreateTime int64         `json:"createTime" form:"createTime"`                                                                      // 创建时间
  UpdateTime int64         `json:"updateTime" form:"updateTime"`                                                                      // 更新时间
}
