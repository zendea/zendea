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

// NodeController node controller
type NodeController struct {
	controller.BaseController
}

// Show show node
func (c *NodeController) Show(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		node := service.NodeService.Get(gDto.ID)
		if node == nil {
			c.Fail(ctx, util.NewErrorMsg("Node not found, id="+strconv.FormatInt(gDto.ID, 10)))
			return
		}
		c.Success(ctx, node)
	}
}

// Store create a node
func (c *NodeController) Store(ctx *gin.Context) {
	var nodeForm form.NodeCreateForm
	if !c.BindAndValidate(ctx, &nodeForm) {
		return
	}
	node, err := service.NodeService.Create(nodeForm)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	c.Success(ctx, node)
}

// Update update a node
func (c *NodeController) Update(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if !c.BindAndValidate(ctx, &gDto) {
		return
	}
	node := service.NodeService.Get(gDto.ID)
	if node == nil {
		c.Fail(ctx, util.NewErrorMsg("Node not found, id="+strconv.FormatInt(gDto.ID, 10)))
		return
	}

	var nodeForm form.NodeUpdateForm
	if !c.BindAndValidate(ctx, &nodeForm) {
		return
	}
	nodeForm.ID = gDto.ID
	err := service.NodeService.Update(nodeForm)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	c.Success(ctx, node)
}

// Delete delete node
func (c *NodeController) Delete(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if !c.BindAndValidate(ctx, &gDto) {
		return
	}
	service.NodeService.Delete(gDto.ID)
	c.Success(ctx, nil)
}

// List list nodes
func (c *NodeController) List(ctx *gin.Context) {
	page := form.FormValueIntDefault(ctx, "page", 1)
	limit := form.FormValueIntDefault(ctx, "limit", 20)
	id := ctx.Request.FormValue("id")
	name := ctx.Request.FormValue("name")
	sectionId := ctx.Request.FormValue("sectionId")

	conditions := sqlcnd.NewSqlCnd()
	if len(id) > 0 {
		conditions.Eq("id", id)
	}
	if len(sectionId) > 0 {
		conditions.Eq("section_id", sectionId)
	}
	if len(name) > 0 {
		conditions.Like("name", name)
	}
	list, paging := service.NodeService.List(conditions.Page(page, limit).Asc("section_id").Asc("sort_no"))
	var results []map[string]interface{}
	for _, node := range list {
		item := util.StructToMap(node)
		item["section"] = cache.SectionCache.Get(node.SectionID)
		results = append(results, item)
	}

	c.Success(ctx, &sqlcnd.PageResult{Results: results, Page: paging})
}
