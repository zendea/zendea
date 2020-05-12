package admin

import (
	"strconv"
	"github.com/gin-gonic/gin"

	"zendea/controller"
	"zendea/util"
	"zendea/util/sqlcnd"
	"zendea/form"
	"zendea/service"
)

// LinkController link controller
type LinkController struct {
	controller.BaseController
}

// Show show link
func (c *LinkController) Show(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		link := service.LinkService.Get(gDto.ID)
		if link == nil {
			c.Fail(ctx, util.NewErrorMsg("Link not found, id=" + strconv.FormatInt(gDto.ID, 10)))
			return
		}
		c.Success(ctx, link)
	}
}

// Store create a link
func (c *LinkController) Store(ctx *gin.Context) {
	var linkForm form.LinkCreateForm
	if !c.BindAndValidate(ctx, &linkForm) {
		return
	}
	link, err := service.LinkService.Create(linkForm)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	c.Success(ctx, link)
}

// Update update a link
func (c *LinkController) Update(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if !c.BindAndValidate(ctx, &gDto) {
		return
	}
	link := service.LinkService.Get(gDto.ID)
	if link == nil {
		c.Fail(ctx, util.NewErrorMsg("Link not found, id=" + strconv.FormatInt(gDto.ID, 10)))
		return
	}
	
	var linkForm form.LinkUpdateForm
	if !c.BindAndValidate(ctx, &linkForm) {
		return
	}
	linkForm.ID = gDto.ID
	err := service.LinkService.Update(linkForm)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	c.Success(ctx, link)
}

// Delete delete link
func (c *LinkController) Delete(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if !c.BindAndValidate(ctx, &gDto) {
		return
	}
	service.LinkService.Delete(gDto.ID)
	c.Success(ctx, nil)
}

// List list links
func (c *LinkController) List(ctx *gin.Context) {
	page := form.FormValueIntDefault(ctx, "page", 1)
	limit := form.FormValueIntDefault(ctx, "limit", 20)
	name := ctx.Request.FormValue("name")

	conditions := sqlcnd.NewSqlCnd()
	if len(name) > 0 {
		conditions.Like("name", name)
	}
	list, paging := service.LinkService.List(conditions.Page(page, limit).Desc("id"))

	c.Success(ctx, &sqlcnd.PageResult{Results: list, Page: paging})
}
