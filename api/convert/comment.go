package convert

import (
	"html/template"

	"zendea/model"
	"zendea/service"
	"zendea/util/markdown"
)

func ToComments(comments []model.Comment) []model.CommentResponse {
	var ret []model.CommentResponse
	for _, comment := range comments {
		ret = append(ret, *ToComment(comment))
	}
	return ret
}

func ToComment(comment model.Comment) *model.CommentResponse {
	return _buildComment(&comment, true)
}

func _buildComment(comment *model.Comment, buildQuote bool) *model.CommentResponse {
	if comment == nil {
		return nil
	}

	ret := &model.CommentResponse{
		CommentId:  comment.ID,
		User:       ToUserDefaultIfNull(comment.UserId),
		EntityType: comment.EntityType,
		EntityId:   comment.EntityId,
		QuoteId:    comment.QuoteId,
		Status:     comment.Status,
		CreateTime: comment.CreateTime,
	}

	if comment.ContentType == model.ContentTypeMarkdown {
		markdownResult := markdown.NewMd().Run(comment.Content)
		ret.Content = template.HTML(ToHtmlContent(markdownResult.ContentHtml))
	} else {
		ret.Content = template.HTML(ToHtmlContent(comment.Content))
	}

	if buildQuote && comment.QuoteId > 0 {
		quote := _buildComment(service.CommentService.Get(comment.QuoteId), false)
		if quote != nil {
			ret.Quote = quote
			ret.QuoteContent = template.HTML(quote.User.Nickname+"：") + quote.Content
		}
	}
	return ret
}
