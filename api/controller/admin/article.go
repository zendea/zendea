package admin

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"

	"zendea/builder"
	"zendea/cache"
	"zendea/controller"
	"zendea/form"
	"zendea/model"
	"zendea/service"
	"zendea/util"
	"zendea/util/markdown"
	"zendea/util/sqlcnd"
	"zendea/util/strtrim"
)

// ArticleController article controller
type ArticleController struct {
	controller.BaseController
}

// Show show article
func (c *ArticleController) Show(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		article := service.ArticleService.Get(gDto.ID)
		if article == nil {
			c.Fail(ctx, util.NewErrorMsg("Article not found, id="+strconv.FormatInt(gDto.ID, 10)))
			return
		}
		c.Success(ctx, article)
	}
}

// Update update a article
func (c *ArticleController) Update(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if !c.BindAndValidate(ctx, &gDto) {
		return
	}
	article := service.ArticleService.Get(gDto.ID)
	if article == nil {
		c.Fail(ctx, util.NewErrorMsg("Article not found, id="+strconv.FormatInt(gDto.ID, 10)))
		return
	}

	var articleForm form.ArticleUpdateForm
	if !c.BindAndValidate(ctx, &articleForm) {
		return
	}
	articleForm.ID = gDto.ID
	err := service.ArticleService.Update(articleForm)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	c.Success(ctx, article)
}

// Delete delete article
func (c *ArticleController) Delete(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if !c.BindAndValidate(ctx, &gDto) {
		return
	}
	service.ArticleService.Delete(gDto.ID)
	c.Success(ctx, nil)
}

// List list articles
func (c *ArticleController) List(ctx *gin.Context) {
	page := form.FormValueIntDefault(ctx, "page", 1)
	limit := form.FormValueIntDefault(ctx, "limit", 20)
	name := ctx.Request.FormValue("name")

	conditions := sqlcnd.NewSqlCnd()
	if len(name) > 0 {
		conditions.Like("name", name)
	}
	list, paging := service.ArticleService.List(conditions.Page(page, limit).Desc("id"))

	var results []map[string]interface{}
	for _, article := range list {
		item := util.StructToMap(article, "content")
		item["user"] = builder.BuildUserDefaultIfNull(article.UserId)

		// 简介
		if article.ContentType == model.ContentTypeMarkdown {
			mr := markdown.NewMd().Run(article.Content)
			if len(article.Summary) == 0 {
				item["summary"] = mr.SummaryText
			}
		} else {
			if len(article.Summary) == 0 {
				doc, err := goquery.NewDocumentFromReader(strings.NewReader(article.Content))
				if err != nil {
					item["summary"] = strtrim.GetTextSummary(doc.Text(), 256)
				}
			}
		}
		// 标签
		tagIds := cache.ArticleTagCache.Get(article.ID)
		tags := cache.TagCache.GetList(tagIds)
		item["tags"] = builder.BuildTags(tags)

		results = append(results, item)
	}

	c.Success(ctx, &sqlcnd.PageResult{Results: results, Page: paging})
}
