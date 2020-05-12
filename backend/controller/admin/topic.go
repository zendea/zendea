package admin

import (
	"strconv"
	"github.com/gin-gonic/gin"

	"zendea/builder"
	"zendea/controller"
	"zendea/util"
	"zendea/util/markdown"
	"zendea/util/sqlcnd"
	"zendea/form"
	"zendea/service"
)

// TopicController topic controller
type TopicController struct {
	controller.BaseController
}

// Show show topic
func (c *TopicController) Show(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		topic := service.TopicService.Get(gDto.ID)
		if topic == nil {
			c.Fail(ctx, util.NewErrorMsg("Topic not found, id=" + strconv.FormatInt(gDto.ID, 10)))
			return
		}
		c.Success(ctx, topic)
	}
}

// Update update a topic
func (c *TopicController) Update(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if !c.BindAndValidate(ctx, &gDto) {
		return
	}
	topic := service.TopicService.Get(gDto.ID)
	if topic == nil {
		c.Fail(ctx, util.NewErrorMsg("Topic not found, id=" + strconv.FormatInt(gDto.ID, 10)))
		return
	}
	
	var topicForm form.TopicUpdateForm
	if !c.BindAndValidate(ctx, &topicForm) {
		return
	}
	topicForm.ID = gDto.ID
	err := service.TopicService.Update(topicForm)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	c.Success(ctx, topic)
}

// Delete delete topic
func (c *TopicController) Delete(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if !c.BindAndValidate(ctx, &gDto) {
		return
	}
	service.TopicService.Delete(gDto.ID)
	c.Success(ctx, nil)
}

// Undelete delete topic
func (c *TopicController) Undelete(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if !c.BindAndValidate(ctx, &gDto) {
		return
	}
	service.TopicService.Undelete(gDto.ID)
	c.Success(ctx, nil)
}

// List list topics
func (c *TopicController) List(ctx *gin.Context) {
	page := form.FormValueIntDefault(ctx, "page", 1)
	limit := form.FormValueIntDefault(ctx, "limit", 20)
	id := ctx.Request.FormValue("id")
	userID := ctx.Request.FormValue("user_id")
	status := ctx.Request.FormValue("status")
	recommend := ctx.Request.FormValue("recommend")
	title := ctx.Request.FormValue("title")

	conditions := sqlcnd.NewSqlCnd()
	if len(id) > 0 {
		conditions.Eq("id", id)
	}
	if len(userID) > 0 {
		conditions.Eq("user_id", userID)
	}
	if len(status) > 0 {
		conditions.Eq("status", status)
	}
	if len(recommend) > 0 {
		conditions.Eq("recommend", recommend)
	}
	if len(title) > 0 {
		conditions.Like("title", title)
	}

	list, paging := service.TopicService.List(conditions.Page(page, limit).Desc("id"))

	var results []map[string]interface{}
	for _, topic := range list {
		result := util.StructToMap(topic, "content")
		result["user"] = builder.BuildUserDefaultIfNull(topic.UserId)
		result["node"] = service.NodeService.Get(topic.NodeId)
		result["tags"] = builder.BuildTags(service.TopicService.GetTopicTags(topic.ID))
		// 简介
		mr := markdown.NewMd().Run(topic.Content)
		result["summary"] = mr.SummaryText

		results = append(results, result)
	}

	c.Success(ctx, &sqlcnd.PageResult{Results: results, Page: paging})
}

// Recommend 推荐
func (c *TopicController) Recommend(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if !c.BindAndValidate(ctx, &gDto) {
		return
	}
	err := service.TopicService.SetRecommend(gDto.ID, true)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	c.Success(ctx, nil)
}

// Unrecommend 取消推荐
func (c *TopicController) Unrecommend(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if !c.BindAndValidate(ctx, &gDto) {
		return
	}
	err := service.TopicService.SetRecommend(gDto.ID, false)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	c.Success(ctx, nil)
}
