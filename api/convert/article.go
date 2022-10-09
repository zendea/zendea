package convert

import (
	"html/template"

	"zendea/cache"
	"zendea/model"
	"zendea/util/markdown"
	"zendea/util/strtrim"
)

func ToArticle(article *model.Article) *model.ArticleResponse {
	if article == nil {
		return nil
	}

	rsp := &model.ArticleResponse{}
	rsp.ArticleId = article.ID
	rsp.Title = article.Title
	rsp.Summary = article.Summary
	rsp.Share = article.Share
	rsp.SourceUrl = article.SourceUrl
	rsp.ViewCount = article.ViewCount
	rsp.CreateTime = article.CreateTime

	rsp.User = ToUserDefaultIfNull(article.UserId)

	tagIds := cache.ArticleTagCache.Get(article.ID)
	tags := cache.TagCache.GetList(tagIds)
	rsp.Tags = ToTags(tags)

	if article.ContentType == model.ContentTypeMarkdown {
		mr := markdown.NewMd(markdown.MdWithTOC()).Run(article.Content)
		rsp.Content = template.HTML(ToHtmlContent(mr.ContentHtml))
		rsp.Toc = template.HTML(mr.TocHtml)
		if len(rsp.Summary) == 0 {
			rsp.Summary = mr.SummaryText
		}
	} else {
		rsp.Content = template.HTML(ToHtmlContent(article.Content))
		if len(rsp.Summary) == 0 {
			rsp.Summary = strtrim.GetTextSummary(article.Content, 256)
		}
	}

	return rsp
}

func ToArticles(articles []model.Article) []model.ArticleResponse {
	if articles == nil || len(articles) == 0 {
		return nil
	}
	var responses []model.ArticleResponse
	for _, article := range articles {
		responses = append(responses, *ToArticle(&article))
	}
	return responses
}

func ToSimpleArticle(article *model.Article) *model.ArticleSimpleResponse {
	if article == nil {
		return nil
	}

	rsp := &model.ArticleSimpleResponse{}
	rsp.ArticleId = article.ID
	rsp.Title = article.Title
	rsp.Summary = article.Summary
	rsp.Share = article.Share
	rsp.SourceUrl = article.SourceUrl
	rsp.ViewCount = article.ViewCount
	rsp.CreateTime = article.CreateTime

	rsp.User = ToUserDefaultIfNull(article.UserId)

	tagIds := cache.ArticleTagCache.Get(article.ID)
	tags := cache.TagCache.GetList(tagIds)
	rsp.Tags = ToTags(tags)

	if article.ContentType == model.ContentTypeMarkdown {
		if len(rsp.Summary) == 0 {
			mr := markdown.NewMd(markdown.MdWithTOC()).Run(article.Content)
			rsp.Summary = mr.SummaryText
		}
	} else {
		if len(rsp.Summary) == 0 {
			rsp.Summary = strtrim.GetTextSummary(strtrim.GetHtmlText(article.Content), 256)
		}
	}

	return rsp
}

func ToSimpleArticles(articles []model.Article) []model.ArticleSimpleResponse {
	if articles == nil || len(articles) == 0 {
		return nil
	}
	var responses []model.ArticleSimpleResponse
	for _, article := range articles {
		responses = append(responses, *ToSimpleArticle(&article))
	}
	return responses
}
