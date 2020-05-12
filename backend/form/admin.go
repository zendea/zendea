package form

// SectionUpdateForm section update form
type SectionUpdateForm struct {
	ID     int64  //非表单赋值
	Name   string `form:"name" json:"name" binding:"required"`
	SortNo int    `form:"sortNo" json:"sortNo"`
}

// SectionCreateForm section create form
type SectionCreateForm struct {
	Name   string `form:"name" json:"name" binding:"required"`
	SortNo int    `form:"sortNo" json:"sortNo"`
}

// NodeUpdateForm node update form
type NodeUpdateForm struct {
	ID          int64     //非表单赋值
	SectionID   int64  `form:"sectionId" json:"sectionId" binding:"required"`
	Name        string `form:"name" json:"name" binding:"required"`
	Description string `form:"description" json:"description"`
	SortNo      int    `form:"sortNo" json:"sortNo"`
	Status      int    `form:"status" json:"status"`
}

// NodeCreateForm node create form
type NodeCreateForm struct {
	Name   string `form:"name" json:"name" binding:"required"`
	Description string `form:"description" json:"description"`
	SortNo      int    `form:"sortNo" json:"sortNo"`
	Status      int    `form:"status" json:"status"`
}

// LinkUpdateForm link update form
type LinkUpdateForm struct {
	ID      int64     //非表单赋值
	Title   string `form:"title" json:"title" binding:"required"`
	URL     string `form:"url" json:"url" binding:"required"`
	Logo    string `form:"logo" json:"logo"`
	Summary string `form:"summary" json:"summary"`
	Status  int    `form:"status" json:"status"`
}

// LinkCreateForm node create form
type LinkCreateForm struct {
	Title   string `form:"title" json:"title" binding:"required"`
	URL     string `form:"url" json:"url" binding:"required"`
	Logo    string `form:"logo" json:"logo"`
	Summary string `form:"summary" json:"summary"`
	Status  int    `form:"status" json:"status"`
}

// CommentUpdateForm comment update form
type CommentUpdateForm struct {
	ID     	   int64  //非表单赋值
	Content    string `form:"content" json:"content" binding:"required,min=3"`
	Status     int    `form:"status" json:"status"`
}
