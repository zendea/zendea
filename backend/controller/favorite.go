package controller

import (
	"github.com/gin-gonic/gin"

	"zendea/util"
	"zendea/form"
	"zendea/service"
)

type FavoriteController struct {
	BaseController
}


// GetFavorited 是否收藏了
func (c *FavoriteController) GetFavorited(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	entityType := ctx.Request.FormValue("entityType")
	entityID := form.FormValueInt64Default(ctx, "entityId", 0)

	data := map[string]interface{}{}
	if user == nil || len(entityType) == 0 || entityID <= 0 {
		data["favorited"] = false
	} else {
		tmp := service.FavoriteService.GetBy(user.ID, entityType, entityID)
		data["favorited"] = tmp != nil
	}
	c.Success(ctx, data)
}

// Delete 取消收藏
func (c *FavoriteController) Delete(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	if user == nil {
		c.Fail(ctx, util.ErrorNotLogin)
		return
	}

	entityType := ctx.Request.FormValue("entityType")
	entityID := form.FormValueInt64Default(ctx, "entityId", 0)

	tmp := service.FavoriteService.GetBy(user.ID, entityType, entityID)
	if tmp != nil {
		service.FavoriteService.Delete(tmp.ID)
	}
	c.Success(ctx, nil)
}