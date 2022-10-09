package controller

import (
	"github.com/gin-gonic/gin"

	"zendea/convert"
	"zendea/cache"
	"zendea/form"
	"zendea/model"
	"zendea/service"
	"zendea/util"
	"zendea/util/sqlcnd"
)

type TopicController struct {
	BaseController
}

// Show 话题详情
func (c *TopicController) Show(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		topic := service.TopicService.Get(gDto.ID)
		if topic == nil || topic.Status != model.StatusOk {
			c.Fail(ctx, util.ErrorTopicNotFound)
			return
		}
		service.TopicService.IncrViewCount(topic.ID) // 增加浏览量
		c.Success(ctx, convert.ToTopic(topic))
	}
}

// List 帖子列表
func (c *TopicController) List(ctx *gin.Context) {
	page := form.FormValueIntDefault(ctx, "page", 1)

	topics, paging := service.TopicService.List(sqlcnd.NewSqlCnd().
		Eq("status", model.StatusOk).
		Page(page, 20).Desc("last_comment_time"))

	data := map[string]interface{}{}
	data["results"] = convert.ToSimpleTopics(topics)
	data["page"] = paging
	c.Success(ctx, data)
}

// Store 发表帖子
func (c *TopicController) Store(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	var topicForm form.TopicCreateForm
	if c.BindAndValidate(ctx, &topicForm) {
		topicForm.UserID = user.ID
		topic, err := service.TopicService.Create(topicForm)
		if err != nil {
			c.Fail(ctx, util.FromError(err))
			return
		}
		c.Success(ctx, convert.ToSimpleTopic(topic))
	}
}

// Edit 为编辑话题准备数据
func (c *TopicController) Edit(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		topic := service.TopicService.Get(gDto.ID)
		if topic == nil || topic.Status != model.StatusOk {
			c.Fail(ctx, util.NewErrorMsg("话题不存在或已被删除"))
			return
		}
		if topic.UserId != user.ID {
			c.Fail(ctx, util.NewErrorMsg("无权限"))
			return
		}

		tags := service.TopicService.GetTopicTags(topic.ID)
		var tagNames []string
		if len(tags) > 0 {
			for _, tag := range tags {
				tagNames = append(tagNames, tag.Name)
			}
		}

		c.Success(ctx, gin.H{
			"topicId": topic.ID,
			"nodeId":  topic.NodeId,
			"title":   topic.Title,
			"content": topic.Content,
			"tags":    tagNames,
		})
	}
}

// Update 更新话题
func (c *TopicController) Update(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	var gDto form.GeneralGetDto
	if !c.BindAndValidate(ctx, &gDto) {
		c.Fail(ctx, util.ErrorTopicNotFound)
		return
	}

	topic := service.TopicService.Get(gDto.ID)
	if topic == nil || topic.Status == model.StatusDeleted {
		c.Fail(ctx, util.ErrorTopicNotFound)
		return
	}

	if topic.UserId != user.ID {
		c.Fail(ctx, util.NewErrorMsg("无权限"))
		return
	}

	var topicForm form.TopicUpdateForm
	if c.BindAndValidate(ctx, &topicForm) {
		topicForm.ID = topic.ID
		err := service.TopicService.Update(topicForm)
		if err != nil {
			c.Fail(ctx, util.FromError(err))
			return
		}
		c.Success(ctx, convert.ToSimpleTopic(topic))
	}
}

// GetRecentLikes 点赞用户
func (c *TopicController) GetRecentLikes(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		topicLikes := service.TopicLikeService.Recent(gDto.ID, 10)
		var users []model.UserInfo
		for _, topicLike := range topicLikes {
			userInfo := convert.ToUserById(topicLike.UserId)
			if userInfo != nil {
				users = append(users, *userInfo)
			}
		}
		c.Success(ctx, users)
	}
}

// 精华帖子
func (c *TopicController) GetTopicsExcellent(ctx *gin.Context) {
	topics := cache.TopicCache.GetRecommendTopics()

	var odd, even []model.Topic
	for i, topic := range topics {
		if i%2 == 1 {
			odd = append(odd, topic)
		} else {
			even = append(even, topic)
		}
	}

	data := make(map[string]interface{})
	data["odd"] = convert.ToSimpleTopics(odd)
	data["even"] = convert.ToSimpleTopics(even)

	c.Success(ctx, data)
}

// 推荐帖子
func (c *TopicController) GetTopicsRecommend(ctx *gin.Context) {
	page := form.FormValueIntDefault(ctx, "page", 1)

	topics, paging := service.TopicService.List(sqlcnd.NewSqlCnd().
		Eq("recommend", true).
		Eq("status", model.StatusOk).
		Page(page, 20).Desc("last_comment_time"))

	data := map[string]interface{}{}
	data["results"] = convert.ToSimpleTopics(topics)
	data["page"] = paging
	c.Success(ctx, data)
}

// 最新发布帖子列表
func (c *TopicController) GetTopicsLast(ctx *gin.Context) {
	page := form.FormValueIntDefault(ctx, "page", 1)

	topics, paging := service.TopicService.List(sqlcnd.NewSqlCnd().
		Eq("status", model.StatusOk).
		Page(page, 20).Desc("id"))

	data := map[string]interface{}{}
	data["results"] = convert.ToSimpleTopics(topics)
	data["page"] = paging
	c.Success(ctx, data)
}

// 无人问津帖子列表
func (c *TopicController) GetTopicsNoreply(ctx *gin.Context) {
	page := form.FormValueIntDefault(ctx, "page", 1)

	topics, paging := service.TopicService.List(sqlcnd.NewSqlCnd().
		Eq("status", model.StatusOk).
		Eq("comment_count", 0).
		Page(page, 20).Desc("last_comment_time"))

	data := map[string]interface{}{}
	data["results"] = convert.ToSimpleTopics(topics)
	data["page"] = paging
	c.Success(ctx, data)
}

// 节点帖子列表
func (c *TopicController) GetNodeTopics(ctx *gin.Context) {
	page := form.FormValueIntDefault(ctx, "page", 1)
	nodeId := form.FormValueInt64Default(ctx, "nodeId", 0)

	topics, paging := service.TopicService.List(sqlcnd.NewSqlCnd().
		Eq("node_id", nodeId).
		Eq("status", model.StatusOk).
		Page(page, 20).Desc("last_comment_time"))

	data := map[string]interface{}{}
	data["results"] = convert.ToSimpleTopics(topics)
	data["page"] = paging

	c.Success(ctx, data)
}

// 标签帖子列表
func (c *TopicController) GetTagTopics(ctx *gin.Context) {
	page := form.FormValueIntDefault(ctx, "page", 1)
	tagId, err := form.FormValueInt64(ctx, "tagId")
	if err != nil {
		c.Fail(ctx, util.ErrorTagNotFound)
		return
	}
	topics, paging := service.TopicService.GetTagTopics(tagId, page)

	data := map[string]interface{}{}
	data["results"] = convert.ToSimpleTopics(topics)
	data["page"] = paging
	c.Success(ctx, data)
}

// GetUserRecent 用户最近的帖子
func (c *TopicController) GetUserRecent(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		topics := service.TopicService.Find(sqlcnd.NewSqlCnd().Where("user_id = ? and status = ?",
			gDto.ID, model.StatusOk).Desc("id").Limit(10))
		c.Success(ctx, convert.ToSimpleTopics(topics))
	}
}

// GetUserTopics 用户的帖子
func (c *TopicController) GetUserTopics(ctx *gin.Context) {
	page := form.FormValueIntDefault(ctx, "page", 1)
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		topics, paging := service.TopicService.List(sqlcnd.NewSqlCnd().
			Eq("user_id", gDto.ID).
			Eq("status", model.StatusOk).
			Page(page, 20).Desc("id"))

		c.Success(ctx, gin.H{
			"results": convert.ToSimpleTopics(topics),
			"page":    paging,
		})
	}
}

// Like 点赞
func (c *TopicController) Like(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		err := service.TopicLikeService.Like(user.ID, gDto.ID)
		if err != nil {
			c.Fail(ctx, util.FromError(err))
			return
		}
		c.Success(ctx, nil)
	}
}

// Favorite 收藏话题
func (c *TopicController) Favorite(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	var gDto form.GeneralGetDto
	if !c.BindAndValidate(ctx, &gDto) {
		c.Fail(ctx, util.ErrorTopicNotFound)
		return
	}
	err := service.FavoriteService.AddTopicFavorite(user.ID, gDto.ID)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	c.Success(ctx, nil)
}
