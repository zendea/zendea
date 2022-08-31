package form

// CommentCreateForm comment create form
type CommentCreateForm struct {
	UserID      int64  //非表单赋值
	EntityType  string `form:"entityType" json:"entityType" binding:"required"`
	EntityID    int64  `form:"entityId" json:"entityId" binding:"required"`
	Content     string `form:"content" json:"content" binding:"required,min=3"`
	QuoteID     int64  `form:"quoteId" json:"quoteId"`
	ContentType string `form:"contentType"`
}
