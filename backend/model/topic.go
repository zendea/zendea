package model

// 话题
type Topic struct {
	Model
	Type              int    `gorm:"not null;index:idx_topic_type" json:"type" form:"type"`          // 类型
	NodeId            int64  `gorm:"not null;index:idx_node_id;" json:"nodeId" form:"nodeId"`        // 节点编号
	UserId            int64  `gorm:"not null;index:idx_topic_user_id;" json:"userId" form:"userId"`  // 用户
	Title             string `gorm:"size:128" json:"title" form:"title"`                             // 标题
	Content           string `gorm:"type:longtext" json:"content" form:"content"`                    // 内容
	ImageList         string `gorm:"type:longtext" json:"imageList" form:"imageList"`                // 图片
	Recommend         bool   `gorm:"not null;index:idx_recommend" json:"recommend" form:"recommend"` // 是否推荐
	ViewCount         int64  `gorm:"not null" json:"viewCount" form:"viewCount"`                     // 查看数量
	CommentCount      int64  `gorm:"not null" json:"commentCount" form:"commentCount"`               // 跟帖数量
	LikeCount         int64  `gorm:"not null" json:"likeCount" form:"likeCount"`                     // 点赞数量
	Status            int    `gorm:"index:idx_topic_status;" json:"status" form:"status"`
	LastCommentUserId int64  `gorm:"index:idx_topic_last_comment_user_id" json:"lastCommentUserId" form:"lastCommentUserId"` // 最后回复时间                            // 状态：0：正常、1：删除
	LastCommentTime   int64  `gorm:"index:idx_topic_last_comment_time" json:"lastCommentTime" form:"lastCommentTime"`        // 最后回复时间
	CreateTime        int64  `gorm:"index:idx_topic_create_time" json:"createTime" form:"createTime"`                        // 创建时间
	ExtraData         string `gorm:"type:text" json:"extraData" form:"extraData"`                                            // 扩展数据
}

// 主题标签
type TopicTag struct {
	Model
	TopicId         int64 `gorm:"not null;index:idx_topic_tag_topic_id;" json:"topicId" form:"topicId"`                // 主题编号
	TagId           int64 `gorm:"not null;index:idx_topic_tag_tag_id;" json:"tagId" form:"tagId"`                      // 标签编号
	Status          int64 `gorm:"not null;index:idx_topic_tag_status" json:"status" form:"status"`                     // 状态：正常、删除
	LastCommentTime int64 `gorm:"index:idx_topic_tag_last_comment_time" json:"lastCommentTime" form:"lastCommentTime"` // 最后回复时间
	CreateTime      int64 `json:"createTime" form:"createTime"`                                                        // 创建时间
}

// 话题点赞
type TopicLike struct {
	Model
	UserId     int64 `gorm:"not null;index:idx_topic_like_user_id;" json:"userId" form:"userId"`    // 用户
	TopicId    int64 `gorm:"not null;index:idx_topic_like_topic_id;" json:"topicId" form:"topicId"` // 主题编号
	CreateTime int64 `json:"createTime" form:"createTime"`                                          // 创建时间
}
