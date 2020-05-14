package model

import (
	"html/template"
)

type UserInfo struct {
	Id           int64    `json:"id"`
	Username     string   `json:"username"`
	Email        string   `json:"email"`
	Nickname     string   `json:"nickname"`
	Avatar       string   `json:"avatar"`
	Level        int      `json:"level"`
	LevelName    string   `json:"levelName"`
	Website      string   `json:"website"`
	Description  string   `json:"description"`
	Score        int      `json:"score"`        // 积分
	TopicCount   int      `json:"topicCount"`   // 话题数量
	CommentCount int      `json:"commentCount"` // 跟帖数量
	PasswordSet  bool     `json:"passwordSet"`  // 密码已设置
	Status       int      `json:"status"`
	CreateTime   int64    `json:"createTime"`
}

type TagResponse struct {
	TagId   int64  `json:"tagId"`
	TagName string `json:"tagName"`
}

type ArticleSimpleResponse struct {
	ArticleId  int64          `json:"articleId"`
	User       *UserInfo      `json:"user"`
	Tags       *[]TagResponse `json:"tags"`
	Title      string         `json:"title"`
	Summary    string         `json:"summary"`
	Share      bool           `json:"share"`
	SourceUrl  string         `json:"sourceUrl"`
	ViewCount  int64          `json:"viewCount"`
	CreateTime int64          `json:"createTime"`
}

type ArticleResponse struct {
	ArticleSimpleResponse
	Content template.HTML `json:"content"`
	Toc     template.HTML `json:"toc"`
}

type SectionResponse struct {
	SectionId int64           `json:"sectionId"`
	Name      string          `json:"name"`
	Nodes     *[]NodeResponse `json:"nodes"`
}

type NodeResponse struct {
	NodeId      int64  `json:"nodeId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TopicCount  int64  `json:"topicCount"`
}

// 帖子列表返回实体
type TopicSimpleResponse struct {
	TopicId         int64          `json:"topicId"`
	Type            int            `json:"type"`
	User            *UserInfo      `json:"user"`
	Node            *NodeResponse  `json:"node"`
	Tags            *[]TagResponse `json:"tags"`
	Title           string         `json:"title"`
	ImageList       *[]string      `json:"imageList"`
	LastCommentUser *UserInfo      `json:"lastCommentUser"`
	LastCommentTime int64          `json:"lastCommentTime"`
	ViewCount       int64          `json:"viewCount"`
	CommentCount    int64          `json:"commentCount"`
	LikeCount       int64          `json:"likeCount"`
	CreateTime      int64          `json:"createTime"`
}

// 帖子详情返回实体
type TopicResponse struct {
	TopicSimpleResponse
	Content template.HTML `json:"content"`
	Toc     template.HTML `json:"toc"`
}

// 回帖详情返回实体
type CommentResponse struct {
	CommentId    int64            `json:"commentId"`
	User         *UserInfo        `json:"user"`
	EntityType   string           `json:"entityType"`
	EntityId     int64            `json:"entityId"`
	Content      template.HTML    `json:"content"`
	QuoteId      int64            `json:"quoteId"`
	Quote        *CommentResponse `json:"quote"`
	QuoteContent template.HTML    `json:"quoteContent"`
	Status       int              `json:"status"`
	CreateTime   int64            `json:"createTime"`
}

type FavoriteResponse struct {
	FavoriteId int64     `json:"favoriteId"`
	EntityType string    `json:"entityType"`
	EntityId   int64     `json:"entityId"`
	Deleted    bool      `json:"deleted"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	User       *UserInfo `json:"user"`
	Url        string    `json:"url"`
	CreateTime int64     `json:"createTime"`
}

// 消息
type NotificationResponse struct {
	MessageId    int64     `json:"messageId"`
	From         *UserInfo `json:"from"`    // 消息发送人
	UserId       int64     `json:"userId"`  // 消息接收人编号
	Content      string    `json:"content"` // 消息内容
	QuoteContent string    `json:"quoteContent"`
	Type         int       `json:"type"`
	Icon         string    `json:"icon"`
	DetailUrl    string    `json:"detailUrl"` // 消息详情url
	ExtraData    string    `json:"extraData"`
	Status       int       `json:"status"`
	CreateTime   int64     `json:"createTime"`
}
