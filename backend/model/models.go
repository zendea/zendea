package model

var Models = []interface{}{
	&User{}, &Tag{}, &Article{}, &ArticleTag{}, &Comment{}, &Favorite{},
	&Topic{}, &Section{}, &Node{}, &TopicTag{}, &TopicLike{}, &Notification{}, &Setting{}, &Link{},
	&LoginSource{}, &Sitemap{}, &UserScore{}, &UserScoreLog{},
}

type Model struct {
	ID int64 `gorm:"PRIMARY_KEY;AUTO_INCREMENT" json:"id" form:"id"`
}

const (
	StatusOk      = 0 // 正常
	StatusDeleted = 1 // 删除
	StatusPending = 2 // 待审核

	UserLevelGeneral  = 0  // 普通用户
	UserLevelAdmin    = 10 // 管理员

	ContentTypeHtml     = "html"
	ContentTypeMarkdown = "markdown"

	EntityTypeArticle = "article"
	EntityTypeTopic   = "topic"
	EntityTypeComment = "comment"

	NotificationStatusUnread = 0 // 消息未读
	NotificationStatusReaded = 1 // 消息已读

	MsgTypeComment = 0   // 回复消息
	MsgTypeTopicLike = 1 // 话题点赞

	LoginSourceTypeGithub = "github"
	LoginSourceTypeGitee  = "gitee"
	LoginSourceTypeQQ     = "qq"

	ScoreTypeIncr = 0 // 积分+
	ScoreTypeDecr = 1 // 积分-

	TopicTypeNormal  = 0 // 普通帖子
	TopicTypeTwitter = 1 // 推文
)
