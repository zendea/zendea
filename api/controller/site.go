package controller

import (
	"github.com/gin-gonic/gin"

	"zendea/cache"
)

type SiteController struct {
	BaseController
}

func (c *SiteController) Stat(ctx *gin.Context) {
	data := make(map[string]interface{})
	data["userCount"] = cache.StatCache.GetUserCount()
	data["topicCount"] = cache.StatCache.GetTopicCount()
	data["commentCount"] = cache.StatCache.GetCommentCount()

	c.Success(ctx, data)
}
