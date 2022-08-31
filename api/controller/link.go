package controller

import (
	"github.com/gin-gonic/gin"

	"zendea/form"
	"zendea/model"
	"zendea/service"
	"zendea/util/sqlcnd"
)

type LinkController struct {
	BaseController
}

// List 列表
func (c *LinkController) List(ctx *gin.Context) {
	page := form.FormValueIntDefault(ctx, "page", 1)

	links, paging := service.LinkService.List(sqlcnd.NewSqlCnd().
		Eq("status", model.StatusOk).Page(page, 20).Asc("id"))

	var results []map[string]interface{}
	for _, v := range links {
		results = append(results, c.buildLink(v))
	}
	c.Success(ctx, gin.H{
		"results": results,
		"paging":  paging,
	})
}

// 前10个链接
func (c *LinkController) GetToplinks(ctx *gin.Context) {
	links := service.LinkService.Find(sqlcnd.NewSqlCnd().
		Eq("status", model.StatusOk).Limit(10).Asc("id"))

	var results []map[string]interface{}
	for _, v := range links {
		results = append(results, c.buildLink(v))
	}
	c.Success(ctx, results)
}

func (c *LinkController) buildLink(link model.Link) map[string]interface{} {
	return map[string]interface{}{
		"linkId":     link.ID,
		"url":        link.Url,
		"title":      link.Title,
		"summary":    link.Summary,
		"logo":       link.Logo,
		"createTime": link.CreateTime,
	}
}
