package form

// TopicCreateForm topic create form
type TopicCreateForm struct {
	UserID    int64  //非表单赋值
	Title     string `form:"title" json:"title" binding:"required"`
	Content   string `form:"content" json:"content" binding:"required"`
	NodeID    int64  `form:"nodeId" json:"nodeId" binding: "required"`
	Tags      string `form:"tags" json:"tags"`
	ImageList string `form:"imageList" json:"imageList"`
}

// TopicUpdateForm topic update form
type TopicUpdateForm struct {
	ID        int64  //非表单赋值
	Title     string `form:"title" json:"title" binding:"required"`
	Content   string `form:"content" json:"content" binding:"required"`
	NodeID    int64  `form:"nodeId" json:"nodeId" binding: "required"`
	Tags      string `form:"tags" json:"tags"`
}