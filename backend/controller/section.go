package controller

import (
	"github.com/gin-gonic/gin"

	"zendea/builder"
	"zendea/cache"
)

type SectionController struct {
	BaseController
}

func (c *SectionController) List(ctx *gin.Context) {
	//sections := service.SectionService.GetSections()
	sections := cache.SectionCache.GetAll()

	c.Success(ctx, builder.BuildSections(sections))
}