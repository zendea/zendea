package convert

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"

	"zendea/model"
	"zendea/service"
	"zendea/util"
	"zendea/util/avatar"
	"zendea/util/strtrim"
	"zendea/util/urls"
)

func ToFavorites(favorites []model.Favorite) []model.FavoriteResponse {
	if favorites == nil || len(favorites) == 0 {
		return nil
	}
	var responses []model.FavoriteResponse
	for _, favorite := range favorites {
		responses = append(responses, *ToFavorite(&favorite))
	}
	return responses
}

func ToNotification(notification *model.Notification) *model.NotificationResponse {
	if notification == nil {
		return nil
	}

	detailUrl := ""
	icon := ""
	if notification.Type == model.MsgTypeComment {
		entityType := gjson.Get(notification.ExtraData, "entityType")
		entityId := gjson.Get(notification.ExtraData, "entityId")
		if entityType.String() == model.EntityTypeArticle {
			detailUrl = urls.ArticleUrl(entityId.Int())
		} else if entityType.String() == model.EntityTypeTopic {
			detailUrl = urls.TopicUrl(entityId.Int())
		}
		icon = "comment"
	} else if notification.Type == model.MsgTypeTopicLike {
		entityId := gjson.Get(notification.ExtraData, "entityId")
		detailUrl = urls.TopicUrl(entityId.Int())
		icon = "heart"
	} else if notification.Type == model.MsgTypeUserWatch {
		entityId := gjson.Get(notification.ExtraData, "entityId")
		detailUrl = urls.UserUrl(entityId.Int())
		icon = "eye"
	}
	from := ToUserDefaultIfNull(notification.FromId)
	if notification.FromId <= 0 {
		from.Nickname = "系统通知"
		from.Avatar = avatar.DefaultAvatar
	}

	return &model.NotificationResponse{
		MessageId:    notification.ID,
		From:         from,
		UserId:       notification.UserId,
		Content:      notification.Content,
		QuoteContent: notification.QuoteContent,
		Type:         notification.Type,
		Icon:         icon,
		DetailUrl:    detailUrl,
		ExtraData:    notification.ExtraData,
		Status:       notification.Status,
		CreateTime:   notification.CreateTime,
	}
}

func ToFavorite(favorite *model.Favorite) *model.FavoriteResponse {
	rsp := &model.FavoriteResponse{}
	rsp.FavoriteId = favorite.ID
	rsp.EntityType = favorite.EntityType
	rsp.CreateTime = favorite.CreateTime

	if favorite.EntityType == model.EntityTypeArticle {
		article := service.ArticleService.Get(favorite.EntityId)
		if article == nil || article.Status != model.StatusOk {
			rsp.Deleted = true
		} else {
			rsp.Url = urls.ArticleUrl(article.ID)
			rsp.User = ToUserById(article.UserId)
			rsp.Title = article.Title
			if article.ContentType == model.ContentTypeMarkdown {
				rsp.Content = util.GetMarkdownSummary(article.Content)
			} else {
				doc, err := goquery.NewDocumentFromReader(strings.NewReader(article.Content))
				if err == nil {
					text := doc.Text()
					rsp.Content = strtrim.GetTextSummary(text, 256)
				}
			}
		}
	} else {
		topic := service.TopicService.Get(favorite.EntityId)
		if topic == nil || topic.Status != model.StatusOk {
			rsp.Deleted = true
		} else {
			rsp.Url = urls.TopicUrl(topic.ID)
			rsp.User = ToUserById(topic.UserId)
			rsp.Title = topic.Title
			rsp.Content = util.GetMarkdownSummary(topic.Content)
		}
	}
	return rsp
}

func ToNotifications(notifications []model.Notification) []model.NotificationResponse {
	if len(notifications) == 0 {
		return nil
	}
	var responses []model.NotificationResponse
	for _, notification := range notifications {
		responses = append(responses, *ToNotification(&notification))
	}
	return responses
}

func ToHtmlContent(htmlContent string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return htmlContent
	}

	doc.Find("a").Each(func(i int, selection *goquery.Selection) {
		href := selection.AttrOr("href", "")

		if len(href) == 0 {
			return
		}

		// 不是内部链接
		if !urls.IsInternalUrl(href) {
			selection.SetAttr("target", "_blank")
			selection.SetAttr("rel", "external nofollow") // 标记站外链接，搜索引擎爬虫不传递权重值
		}

		// 如果是锚链接
		if urls.IsAnchor(href) {
			selection.ReplaceWithHtml(selection.Text())
		}

		// 如果a标签没有title，那么设置title
		title := selection.AttrOr("title", "")
		if len(title) == 0 {
			selection.SetAttr("title", selection.Text())
		}
	})

	// 处理图片
	doc.Find("img").Each(func(i int, selection *goquery.Selection) {
		src := selection.AttrOr("src", "")
		// 处理第三方图片
		if strings.Contains(src, "qpic.cn") {
			src = util.ParseUrl("/api/img/proxy").AddQuery("url", src).BuildStr()
			// selection.SetAttr("src", src)
		}

		// 处理lazyload
		selection.SetAttr("data-src", src)
		selection.RemoveAttr("src")
	})

	html, err := doc.Find("body").Html()
	if err != nil {
		return htmlContent
	}
	return html
}
