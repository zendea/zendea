package builder

import (
	"html/template"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"

	"zendea/cache"
	"zendea/model"
	"zendea/service"
	"zendea/util"
	"zendea/util/log"
	"zendea/util/avatar"
	"zendea/util/markdown"
	"zendea/util/strtrim"
	"zendea/util/urls"
)

func BuildUserDefaultIfNull(id int64) *model.UserInfo {
	user := cache.UserCache.Get(id)
	if user == nil {
		user = &model.User{}
		user.ID = id
		user.Username = util.SqlNullString(strconv.FormatInt(id, 10))
		user.Avatar = avatar.DefaultAvatar
		user.CreateTime = util.NowTimestamp()
	}
	return BuildUser(user)
}

func BuildUserById(id int64) *model.UserInfo {
	user := cache.UserCache.Get(id)
	return BuildUser(user)
}

func BuildUser(user *model.User) *model.UserInfo {
	if user == nil {
		return nil
	}
	a := user.Avatar
	if len(a) == 0 {
		a = avatar.DefaultAvatar
	}
	levelName := "普通用户"
	if user.Level == model.UserLevelAdmin {
		levelName = "管理员"
	}
	ret := &model.UserInfo{
		Id:           user.ID,
		Username:     user.Username.String,
		Nickname:     user.Nickname,
		Avatar:       a,
		Email:        user.Email.String,
		Level:        user.Level,
		LevelName:    levelName,
		Website:      user.Website,
		Description:  user.Description,
		TopicCount:   user.TopicCount,
		CommentCount: user.CommentCount,
		PasswordSet:  len(user.Password) > 0,
		Status:       user.Status,
		CreateTime:   user.CreateTime,
	}
	if user.Status == model.StatusDeleted {
		ret.Username = "blacklist"
		ret.Nickname = "黑名单用户"
		ret.Avatar = avatar.DefaultAvatar
		ret.Email = ""
		ret.Website = ""
		ret.Description = ""
	} else {
		ret.Score = cache.UserCache.GetScore(user.ID)
	}
	return ret
}

func BuildUsers(users []model.User) []model.UserInfo {
	if len(users) == 0 {
		return nil
	}
	var responses []model.UserInfo
	for _, user := range users {
		item := BuildUser(&user)
		if item != nil {
			responses = append(responses, *item)
		}
	}
	return responses
}

func BuildArticle(article *model.Article) *model.ArticleResponse {
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

	rsp.User = BuildUserDefaultIfNull(article.UserId)

	tagIds := cache.ArticleTagCache.Get(article.ID)
	tags := cache.TagCache.GetList(tagIds)
	rsp.Tags = BuildTags(tags)

	if article.ContentType == model.ContentTypeMarkdown {
		mr := markdown.NewMd(markdown.MdWithTOC()).Run(article.Content)
		rsp.Content = template.HTML(BuildHtmlContent(mr.ContentHtml))
		rsp.Toc = template.HTML(mr.TocHtml)
		if len(rsp.Summary) == 0 {
			rsp.Summary = mr.SummaryText
		}
	} else {
		rsp.Content = template.HTML(BuildHtmlContent(article.Content))
		if len(rsp.Summary) == 0 {
			rsp.Summary = strtrim.GetTextSummary(article.Content, 256)
		}
	}

	return rsp
}

func BuildArticles(articles []model.Article) []model.ArticleResponse {
	if articles == nil || len(articles) == 0 {
		return nil
	}
	var responses []model.ArticleResponse
	for _, article := range articles {
		responses = append(responses, *BuildArticle(&article))
	}
	return responses
}

func BuildSimpleArticle(article *model.Article) *model.ArticleSimpleResponse {
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

	rsp.User = BuildUserDefaultIfNull(article.UserId)

	tagIds := cache.ArticleTagCache.Get(article.ID)
	tags := cache.TagCache.GetList(tagIds)
	rsp.Tags = BuildTags(tags)

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

func BuildSimpleArticles(articles []model.Article) []model.ArticleSimpleResponse {
	if articles == nil || len(articles) == 0 {
		return nil
	}
	var responses []model.ArticleSimpleResponse
	for _, article := range articles {
		responses = append(responses, *BuildSimpleArticle(&article))
	}
	return responses
}

func BuildSection(section *model.Section) *model.SectionResponse {
	if section == nil {
		return nil
	}
	rsp := &model.SectionResponse{}
	rsp.SectionId = section.ID
	rsp.Name = section.Name

	nodes := service.SectionService.GetSectionNodes(section.ID)
	rsp.Nodes = BuildNodes(nodes)

	return rsp
}

func BuildSections(sections []model.Section) *[]model.SectionResponse {
	if len(sections) == 0 {
		return nil
	}
	var ret []model.SectionResponse
	for _, section := range sections {
		ret = append(ret, *BuildSection(&section))
	}
	return &ret
}

func BuildNode(node *model.Node) *model.NodeResponse {
	if node == nil {
		return nil
	}
	return &model.NodeResponse{
		NodeId:      node.ID,
		Name:        node.Name,
		Description: node.Description,
		TopicCount:  node.TopicCount,
	}
}

//func BuildTags(tags []model.Tag) *[]model.TagResponse {
func BuildNodes(nodes []model.Node) *[]model.NodeResponse {
	if len(nodes) == 0 {
		return nil
	}
	var ret []model.NodeResponse
	for _, node := range nodes {
		ret = append(ret, *BuildNode(&node))
	}
	return &ret
}

func BuildTopic(topic *model.Topic) *model.TopicResponse {
	if topic == nil {
		return nil
	}

	rsp := &model.TopicResponse{}

	rsp.TopicId = topic.ID
	rsp.Type = topic.Type
	rsp.Title = topic.Title
	rsp.User = BuildUserDefaultIfNull(topic.UserId)
	rsp.LastCommentTime = topic.LastCommentTime
	rsp.CreateTime = topic.CreateTime
	rsp.ViewCount = topic.ViewCount
	rsp.CommentCount = topic.CommentCount
	rsp.LikeCount = topic.LikeCount

	if topic.NodeId > 0 {
		node := service.NodeService.Get(topic.NodeId)
		rsp.Node = BuildNode(node)
	}

	tags := service.TopicService.GetTopicTags(topic.ID)
	rsp.Tags = BuildTags(tags)

	mr := markdown.NewMd(markdown.MdWithTOC()).Run(topic.Content)
	rsp.Content = template.HTML(BuildHtmlContent(mr.ContentHtml))
	rsp.Toc = template.HTML(mr.TocHtml)

	if len(topic.ImageList) > 0 {
		if err := util.ParseJson(topic.ImageList, &rsp.ImageList); err != nil {
			log.Error(err.Error())
		}
	}

	return rsp
}

func BuildSimpleTopic(topic *model.Topic) *model.TopicSimpleResponse {
	if topic == nil {
		return nil
	}

	rsp := &model.TopicSimpleResponse{}

	rsp.TopicId = topic.ID
	rsp.Type = topic.Type
	rsp.Title = topic.Title
	rsp.User = BuildUserDefaultIfNull(topic.UserId)
	rsp.LastCommentUser = BuildUserDefaultIfNull(topic.LastCommentUserId)
	rsp.LastCommentTime = topic.LastCommentTime
	rsp.CreateTime = topic.CreateTime
	rsp.ViewCount = topic.ViewCount
	rsp.CommentCount = topic.CommentCount
	rsp.LikeCount = topic.LikeCount

	if len(topic.ImageList) > 0 {
		if err := util.ParseJson(topic.ImageList, &rsp.ImageList); err != nil {
			log.Error(err.Error())
		}
	}

	if topic.NodeId > 0 {
		node := service.NodeService.Get(topic.NodeId)
		rsp.Node = BuildNode(node)
	}

	tags := service.TopicService.GetTopicTags(topic.ID)
	rsp.Tags = BuildTags(tags)
	return rsp
}

func BuildSimpleTopics(topics []model.Topic) []model.TopicSimpleResponse {
	if topics == nil || len(topics) == 0 {
		return nil
	}
	var responses []model.TopicSimpleResponse
	for _, topic := range topics {
		responses = append(responses, *BuildSimpleTopic(&topic))
	}
	return responses
}

func BuildComments(comments []model.Comment) []model.CommentResponse {
	var ret []model.CommentResponse
	for _, comment := range comments {
		ret = append(ret, *BuildComment(comment))
	}
	return ret
}

func BuildComment(comment model.Comment) *model.CommentResponse {
	return _buildComment(&comment, true)
}

func _buildComment(comment *model.Comment, buildQuote bool) *model.CommentResponse {
	if comment == nil {
		return nil
	}

	ret := &model.CommentResponse{
		CommentId:  comment.ID,
		User:       BuildUserDefaultIfNull(comment.UserId),
		EntityType: comment.EntityType,
		EntityId:   comment.EntityId,
		QuoteId:    comment.QuoteId,
		Status:     comment.Status,
		CreateTime: comment.CreateTime,
	}

	if comment.ContentType == model.ContentTypeMarkdown {
		markdownResult := markdown.NewMd().Run(comment.Content)
		ret.Content = template.HTML(BuildHtmlContent(markdownResult.ContentHtml))
	} else {
		ret.Content = template.HTML(BuildHtmlContent(comment.Content))
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

func BuildTag(tag *model.Tag) *model.TagResponse {
	if tag == nil {
		return nil
	}
	return &model.TagResponse{TagId: tag.ID, TagName: tag.Name}
}

func BuildTags(tags []model.Tag) *[]model.TagResponse {
	if len(tags) == 0 {
		return nil
	}
	var responses []model.TagResponse
	for _, tag := range tags {
		responses = append(responses, *BuildTag(&tag))
	}
	return &responses
}

func BuildFavorites(favorites []model.Favorite) []model.FavoriteResponse {
	if favorites == nil || len(favorites) == 0 {
		return nil
	}
	var responses []model.FavoriteResponse
	for _, favorite := range favorites {
		responses = append(responses, *BuildFavorite(&favorite))
	}
	return responses
}

func BuildNotification(notification *model.Notification) *model.NotificationResponse {
	if notification == nil {
		return nil
	}

	detailUrl := ""
	if notification.Type == model.MsgTypeComment {
		entityType := gjson.Get(notification.ExtraData, "entityType")
		entityId := gjson.Get(notification.ExtraData, "entityId")
		if entityType.String() == model.EntityTypeArticle {
			detailUrl = urls.ArticleUrl(entityId.Int())
		} else if entityType.String() == model.EntityTypeTopic {
			detailUrl = urls.TopicUrl(entityId.Int())
		}
	}
	from := BuildUserDefaultIfNull(notification.FromId)
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
		DetailUrl:    detailUrl,
		ExtraData:    notification.ExtraData,
		Status:       notification.Status,
		CreateTime:   notification.CreateTime,
	}
}

func BuildFavorite(favorite *model.Favorite) *model.FavoriteResponse {
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
			rsp.User = BuildUserById(article.UserId)
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
			rsp.User = BuildUserById(topic.UserId)
			rsp.Title = topic.Title
			rsp.Content = util.GetMarkdownSummary(topic.Content)
		}
	}
	return rsp
}

func BuildNotifications(notifications []model.Notification) []model.NotificationResponse {
	if len(notifications) == 0 {
		return nil
	}
	var responses []model.NotificationResponse
	for _, notification := range notifications {
		responses = append(responses, *BuildNotification(&notification))
	}
	return responses
}

func BuildHtmlContent(htmlContent string) string {
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
