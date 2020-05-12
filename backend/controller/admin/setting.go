package admin

import (
	"github.com/gin-gonic/gin"

	"zendea/controller"
	"zendea/util"
	// "zendea/form"
	"zendea/service"
)

// SettingController setting controller
type SettingController struct {
	controller.BaseController
}

// List list settings
func (c *SettingController) List(ctx *gin.Context) {
	c.Success(ctx, service.SettingService.GetSetting())
}

// Store store settings
func (c *SettingController) Store(ctx *gin.Context) {
	config := ctx.Request.FormValue("config")
	if err := service.SettingService.SetAll(config); err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}

	c.Success(ctx, nil)
}
