package form

// ArticleCreateForm topic create form
type ArticleCreateForm struct {
	UserID  int64  //非表单赋值
	Title   string `form:"title" json:"title" binding:"required"`
	Summary string `form:"summary" json:"summary"`
	Content string `form:"content" json:"content" binding:"required"`
	Tags    string `form:"tags" json:"tags"`
}

// ArticleUpdateForm topic update form
type ArticleUpdateForm struct {
	ID      int64  //非表单赋值
	Title   string `form:"title" json:"title" binding:"required"`
	Summary string `form:"summary" json:"summary"`
	Content string `form:"content" json:"content" binding:"required"`
	Tags    string `form:"tags" json:"tags"`
}
