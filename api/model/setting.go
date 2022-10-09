package model

// 系统配置
type Setting struct {
	Model
	Key         string `gorm:"not null;size:128;unique" json:"key" form:"key"` // 配置key
	Value       string `gorm:"type:text" json:"value" form:"value"`            // 配置值
	Name        string `gorm:"not null;size:32" json:"name" form:"name"`       // 配置名称
	Description string `gorm:"size:128" json:"description" form:"description"` // 配置描述
	CreateTime  int64  `gorm:"not null" json:"createTime" form:"createTime"`   // 创建时间
	UpdateTime  int64  `gorm:"not null" json:"updateTime" form:"updateTime"`   // 更新时间
}

// 站点导航
type SiteNav struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

// 小贴士
type SiteTip struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// 积分配置
type ScoreConfig struct {
	PostTopicScore   int `json:"postTopicScore"`   // 发帖获得积分
	PostCommentScore int `json:"postCommentScore"` // 跟帖获得积分
}

// 配置返回结构体
type ConfigData struct {
	SiteTitle        string      `json:"siteTitle"`
	SiteDescription  string      `json:"siteDescription"`
	SiteKeywords     []string    `json:"siteKeywords"`
	SiteNavs         []SiteNav   `json:"siteNavs"`
	SiteTips         []SiteTip   `json:"siteTips"`
	SiteNotification string      `json:"siteNotification"`
	SiteIndexHtml    string      `json:"siteIndexHtml"`
	RecommendTags    []string    `json:"recommendTags"`
	ScoreConfig      ScoreConfig `json:"scoreConfig"`
	DefaultNodeId    int64       `json:"defaultNodeId"`
}

type AppData struct {
	Name           string `json:"name"`
	Version        string `json:"version"`
	UserLevelAdmin int    `json:"user_level_admin"`
}
