package controller

import (
	"github.com/gin-gonic/gin"

	"zendea/builder"
	"zendea/service"
)

type SectionController struct {
	BaseController
}

func (c *SectionController) List(ctx *gin.Context) {
	sections := service.SectionService.GetSections()
	c.Success(ctx, builder.BuildSections(sections))
}