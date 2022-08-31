package admin

import (
	"github.com/gin-gonic/gin"
	"strconv"

	"zendea/cache"
	"zendea/controller"
	"zendea/form"
	"zendea/service"
	"zendea/util"
	"zendea/util/sqlcnd"
)

// SectionController section controller
type SectionController struct {
	controller.BaseController
}

// Show show section
func (c *SectionController) Show(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		section := service.SectionService.Get(gDto.ID)
		if section == nil {
			c.Fail(ctx, util.NewErrorMsg("Section not found, id="+strconv.FormatInt(gDto.ID, 10)))
			return
		}
		c.Success(ctx, section)
	}
}

// Store create a section
func (c *SectionController) Store(ctx *gin.Context) {
	var sectionForm form.SectionCreateForm
	if !c.BindAndValidate(ctx, &sectionForm) {
		return
	}
	section, err := service.SectionService.Create(sectionForm)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	cache.SectionCache.Invalidate(section.ID)

	c.Success(ctx, section)
}

// Update update a section
func (c *SectionController) Update(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if !c.BindAndValidate(ctx, &gDto) {
		return
	}
	section := service.SectionService.Get(gDto.ID)
	if section == nil {
		c.Fail(ctx, util.NewErrorMsg("Section not found, id="+strconv.FormatInt(gDto.ID, 10)))
		return
	}

	var sectionForm form.SectionUpdateForm
	if !c.BindAndValidate(ctx, &sectionForm) {
		return
	}
	sectionForm.ID = gDto.ID
	err := service.SectionService.Update(sectionForm)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	cache.SectionCache.Invalidate(section.ID)

	c.Success(ctx, section)
}

// Delete delete section
func (c *SectionController) Delete(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if !c.BindAndValidate(ctx, &gDto) {
		return
	}
	service.SectionService.Delete(gDto.ID)
	cache.SectionCache.Invalidate(gDto.ID)
	c.Success(ctx, nil)
}

// List list sections
func (c *SectionController) List(ctx *gin.Context) {
	page := form.FormValueIntDefault(ctx, "page", 1)
	limit := form.FormValueIntDefault(ctx, "limit", 20)
	name := ctx.Request.FormValue("name")

	conditions := sqlcnd.NewSqlCnd()
	if len(name) > 0 {
		conditions.Like("name", name)
	}
	list, paging := service.SectionService.List(conditions.Page(page, limit).Asc("sort_no"))

	c.Success(ctx, &sqlcnd.PageResult{Results: list, Page: paging})
}

// Nodes list nodes
func (c *SectionController) Nodes(ctx *gin.Context) {
	page := form.FormValueIntDefault(ctx, "page", 1)
	limit := form.FormValueIntDefault(ctx, "limit", 20)

	var gDto form.GeneralGetDto
	if !c.BindAndValidate(ctx, &gDto) {
		return
	}
	conditions := sqlcnd.NewSqlCnd()
	conditions.Eq("section_id", gDto.ID)

	list, paging := service.NodeService.List(conditions.Page(page, limit).Asc("sort_no"))

	c.Success(ctx, &sqlcnd.PageResult{Results: list, Page: paging})
}
