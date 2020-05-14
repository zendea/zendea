package service

import (
	"sync"

	"zendea/cache"
	"zendea/dao"
	"zendea/model"
	"zendea/util"
	"zendea/util/log"
	"zendea/util/email"
	"zendea/util/sqlcnd"
	"zendea/util/urls"
)

var NotificationService = newNotificationService()

func newNotificationService() *notificationService {
	return &notificationService{
		notificationsChan: make(chan *model.Notification),
	}
}

type notificationService struct {
	notificationsChan        chan *model.Notification
	notificationsConsumeOnce sync.Once
}

func (s *notificationService) Get(id int64) *model.Notification {
	return dao.NotificationDao.Get(id)
}

func (s *notificationService) Take(where ...interface{}) *model.Notification {
	return dao.NotificationDao.Take(where...)
}

func (s *notificationService) Find(cnd *sqlcnd.SqlCnd) []model.Notification {
	return dao.NotificationDao.Find(cnd)
}

func (s *notificationService) FindOne(cnd *sqlcnd.SqlCnd) *model.Notification {
	return dao.NotificationDao.FindOne(cnd)
}

func (s *notificationService) List(cnd *sqlcnd.SqlCnd) (list []model.Notification, paging *sqlcnd.Paging) {
	return dao.NotificationDao.List(cnd)
}

func (s *notificationService) Create(t *model.Notification) error {
	return dao.NotificationDao.Create(t)
}

func (s *notificationService) Update(t *model.Notification) error {
	return dao.NotificationDao.Update(t)
}

func (s *notificationService) Updates(id int64, columns map[string]interface{}) error {
	return dao.NotificationDao.Updates(id, columns)
}

func (s *notificationService) UpdateColumn(id int64, name string, value interface{}) error {
	return dao.NotificationDao.UpdateColumn(id, name, value)
}

func (s *notificationService) Delete(id int64) {
	dao.NotificationDao.Delete(id)
}

// 获取未读消息数量
func (s *notificationService) GetUnReadCount(userId int64) (count int64) {
	return dao.NotificationDao.GetUnReadCount(userId)
}

// 将所有消息标记为已读
func (s *notificationService) MarkRead(userId int64) error {
	return dao.NotificationDao.UpdateStatusBatch(userId)
}

// 内容被点赞
func (s *notificationService) SendTopicLikeNotification(topicLike *model.TopicLike) {
	user := cache.UserCache.Get(topicLike.UserId)

	var (
		fromId       = topicLike.UserId // 消息发送人
		authorId     int64              // 点赞者编号
		content      string             // 消息内容
		quoteContent string             // 引用内容
	)
	topic := dao.TopicDao.Get(topicLike.TopicId)
	if topic != nil {
		authorId = topic.UserId
		content = user.Username.String + " 点赞了你的话题：" + topic.Title
		quoteContent = ""
	}

	if authorId <= 0 {
		return
	}
	// 给帖子作者发消息
	s.Produce(fromId, authorId, content, quoteContent, model.MsgTypeTopicLike, map[string]interface{}{
		"entityType":  model.EntityTypeTopic,
		"entityId":    topic.ID,
		"topicLikeId": topicLike.ID,
	})
}

// 评论被回复消息
func (s *notificationService) SendCommentNotification(comment *model.Comment) {
	user := cache.UserCache.Get(comment.UserId)
	quote := s.getQuoteComment(comment.QuoteId)
	summary := util.GetMarkdownSummary(comment.Content)

	var (
		fromId       = comment.UserId // 消息发送人
		authorId     int64            // 帖子作者编号
		content      string           // 消息内容
		quoteContent string           // 引用内容
	)

	if comment.EntityType == model.EntityTypeArticle { // 文章被评论
		article := dao.ArticleDao.Get(comment.EntityId)
		if article != nil {
			authorId = article.UserId
			content = user.Username.String + " 回复了你的文章：" + summary
			quoteContent = "《" + article.Title + "》"
		}
	} else if comment.EntityType == model.EntityTypeTopic { // 话题被评论
		topic := dao.TopicDao.Get(comment.EntityId)
		if topic != nil {
			authorId = topic.UserId
			content = user.Username.String + " 回复了你的话题：" + summary
			quoteContent = "《" + topic.Title + "》"
		}
	}

	if authorId <= 0 {
		return
	}

	if quote != nil { // 回复跟帖
		if comment.UserId != authorId && quote.UserId != authorId { // 回复人和帖子作者不是同一个人，并且引用的用户不是帖子作者，需要给帖子作者也发送一下消息
			// 给帖子作者发消息
			s.Produce(fromId, authorId, content, quoteContent, model.MsgTypeComment, map[string]interface{}{
				"entityType": comment.EntityType,
				"entityId":   comment.EntityId,
				"commentId":  comment.ID,
				"quoteId":    comment.QuoteId,
			})
		}

		// 给被引用的人发消息
		s.Produce(fromId, quote.UserId, user.Username.String+" 回复了你的评论："+summary, util.GetMarkdownSummary(quote.Content), model.MsgTypeComment, map[string]interface{}{
			"entityType": comment.EntityType,
			"entityId":   comment.EntityId,
			"commentId":  comment.ID,
			"quoteId":    comment.QuoteId,
		})
	} else if comment.UserId != authorId { // 回复主贴，并且不是自己回复自己
		// 给帖子作者发消息
		s.Produce(fromId, authorId, content, quoteContent, model.MsgTypeComment, map[string]interface{}{
			"entityType": comment.EntityType,
			"entityId":   comment.EntityId,
			"commentId":  comment.ID,
			"quoteId":    comment.QuoteId,
		})
	}
}

func (s *notificationService) getQuoteComment(quoteId int64) *model.Comment {
	if quoteId <= 0 {
		return nil
	}
	return dao.CommentDao.Get(quoteId)
}

// 生产，将消息数据放入chan
func (s *notificationService) Produce(fromId, toId int64, content, quoteContent string, msgType int, extraDataMap map[string]interface{}) {
	to := cache.UserCache.Get(toId)
	if to == nil {
		return
	}

	s.Consume()

	var (
		extraData string
		err       error
	)
	if extraData, err = util.FormatJson(extraDataMap); err != nil {
		log.Error("格式化extraData错误")
	}
	s.notificationsChan <- &model.Notification{
		FromId:       fromId,
		UserId:       toId,
		Content:      content,
		QuoteContent: quoteContent,
		Type:         msgType,
		ExtraData:    extraData,
		Status:       model.NotificationStatusUnread,
		CreateTime:   util.NowTimestamp(),
	}
}

// 消费，消费chan中的消息
func (s *notificationService) Consume() {
	s.notificationsConsumeOnce.Do(func() {
		go func() {
			log.Info("开始消费系统消息...")
			for {
				msg := <-s.notificationsChan
				log.Info("处理消息：from=%s to=%s", msg.FromId, msg.UserId)

				if err := s.Create(msg); err != nil {
					log.Info("创建消息发生异常...")
				} else {
					s.SendEmailNotice(msg)
				}
			}
		}()
	})
}

// 发送邮件通知
func (s *notificationService) SendEmailNotice(notification *model.Notification) {
	user := cache.UserCache.Get(notification.UserId)
	if user != nil && len(user.Email.String) > 0 {
		siteTitle := cache.SettingCache.GetValue(model.SettingSiteTitle)
		emailTitle := siteTitle + " 新消息提醒"

		email.SendTemplateEmail(user.Email.String, emailTitle, emailTitle, notification.Content,
			notification.QuoteContent, urls.AbsUrl("/user/notifications"))
		log.Info("发送邮件...email=%s", user.Email)
	} else {
		log.Info("邮件未发送，没设置邮箱...")
	}
}
