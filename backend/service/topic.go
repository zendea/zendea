package service

import (
	"math"
	"path"
	"time"
	"errors"

	"github.com/gorilla/feeds"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"

	"zendea/cache"
	"zendea/dao"
	"zendea/model"
	"zendea/form"
	"zendea/util"
	"zendea/util/log"
	"zendea/util/sqlcnd"
	"zendea/util/urls"
)

type ScanTopicCallback func(topics []model.Topic)

var TopicService = newTopicService()

func newTopicService() *topicService {
	return &topicService{}
}

type topicService struct{}

func (s *topicService) Get(id int64) *model.Topic {
	return dao.TopicDao.Get(id)
}

func (s *topicService) Find(cnd *sqlcnd.SqlCnd) []model.Topic {
	return dao.TopicDao.Find(cnd)
}

func (s *topicService) List(cnd *sqlcnd.SqlCnd) (list []model.Topic, paging *sqlcnd.Paging) {
	return dao.TopicDao.List(cnd)
}

func (s *topicService) Count(cnd *sqlcnd.SqlCnd) int {
	return dao.TopicDao.Count(cnd)
}

// 删除
func (s *topicService) Delete(id int64) error {
	err := dao.TopicDao.UpdateColumn(id, "status", model.StatusDeleted)
	if err == nil {
		// 删掉标签文章
		TopicTagService.DeleteByTopicId(id)
	}
	return err
}

func (s *topicService) Update(dto form.TopicUpdateForm) error {
	node := dao.NodeDao.Get(dto.NodeID)
	if node == nil || node.Status != model.StatusOk {
		return util.NewErrorMsg("节点不存在")
	}
	err := dao.Tx(func(tx *gorm.DB) error {
		err := dao.TopicDao.Updates(dto.ID, map[string]interface{}{
			"node_id":     dto.NodeID,
			"title":       dto.Title,
			"content":     dto.Content,
			"update_time": util.NowTimestamp(),
		})
		if err != nil {
			return err
		}
		tagIds := dao.TagDao.GetOrCreates(util.ParseTagsToArray(dto.Tags)) // 创建文章对应标签
		dao.TopicTagDao.DeleteTopicTags(dto.ID)                           // 先删掉所有的标签
		dao.TopicTagDao.AddTopicTags(dto.ID, tagIds)                      // 然后重新添加标签
		return nil
	})

	return err
}

// 取消删除
func (s *topicService) Undelete(id int64) error {
	err := dao.TopicDao.UpdateColumn(id, "status", model.StatusOk)
	if err == nil {
		// 删掉标签文章
		TopicTagService.UndeleteByTopicId(id)
	}
	return err
}

// 发表话题
func (s *topicService) Create(dto form.TopicCreateForm) (*model.Topic, error) {
	nodeID := dto.NodeID
	if nodeID <= 0 {
		nodeID = SettingService.GetSetting().DefaultNodeId
		if nodeID <= 0 {
			return nil, errors.New("请配置默认节点")
		}
	}
	node := dao.NodeDao.Get(nodeID)
	if node == nil || node.Status != model.StatusOk {
		return nil, errors.New("节点不存在")
	}

	now := util.NowTimestamp()
	topic := &model.Topic{
		Type:            model.TopicTypeNormal,
		UserId:          dto.UserID,
		NodeId:          nodeID,
		Title:           dto.Title,
		Content:         dto.Content,
		ImageList:       dto.ImageList,
		Status:          model.StatusOk,
		LastCommentTime: now,
		CreateTime:      now,
	}

	err := dao.Tx(func(tx *gorm.DB) error {
		tagIds := dao.TagDao.GetOrCreates(util.ParseTagsToArray(dto.Tags))
		err := dao.TopicDao.Create(topic)
		if err != nil {
			return err
		}

		dao.TopicTagDao.AddTopicTags(topic.ID, tagIds)
		return nil
	})
	if err == nil {
		// 节点话题计数
		NodeService.IncrTopicCount(nodeID)
		// 用户话题计数
		UserService.IncrTopicCount(dto.UserID)
		// 获得积分
		UserScoreService.IncrementPostTopicScore(topic)
	}
	return topic, err
}

// 推荐
func (s *topicService) SetRecommend(topicId int64, recommend bool) error {
	return dao.TopicDao.UpdateColumn(topicId, "recommend", recommend)
}

// 话题的标签
func (s *topicService) GetTopicTags(topicId int64) []model.Tag {
	topicTags := dao.TopicTagDao.Find(sqlcnd.NewSqlCnd().Where("topic_id = ?", topicId))

	var tagIds []int64
	for _, topicTag := range topicTags {
		tagIds = append(tagIds, topicTag.TagId)
	}
	return cache.TagCache.GetList(tagIds)
}

// 指定标签下话题列表
func (s *topicService) GetTagTopics(tagId int64, page int) (topics []model.Topic, paging *sqlcnd.Paging) {
	topicTags, paging := dao.TopicTagDao.List(sqlcnd.NewSqlCnd().
		Eq("tag_id", tagId).
		Eq("status", model.StatusOk).
		Page(page, 20).Desc("last_comment_time"))
	if len(topicTags) > 0 {
		var topicIds []int64
		for _, topicTag := range topicTags {
			topicIds = append(topicIds, topicTag.TopicId)
		}

		topicsMap := s.GetTopicInIds(topicIds)
		if topicsMap != nil {
			for _, topicTag := range topicTags {
				if topic, found := topicsMap[topicTag.TopicId]; found {
					topics = append(topics, topic)
				}
			}
		}
	}
	return
}

// GetTopicInIds 根据编号批量获取主题
func (s *topicService) GetTopicInIds(topicIds []int64) map[int64]model.Topic {
	if len(topicIds) == 0 {
		return nil
	}
	var topics []model.Topic
	dao.DB().Where("id in (?)", topicIds).Find(&topics)

	topicsMap := make(map[int64]model.Topic, len(topics))
	for _, topic := range topics {
		topicsMap[topic.ID] = topic
	}
	return topicsMap
}

// 浏览数+1
func (s *topicService) IncrViewCount(topicId int64) {
	dao.DB().Model(&model.Topic{}).Where("id = ?", topicId).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))
}

// 当帖子被评论的时候，更新最后回复时间、回复数量+1
func (s *topicService) OnComment(topicId, lastCommentUserId, lastCommentTime int64) {
	dao.Tx(func(tx *gorm.DB) error {
		if err := dao.DB().Model(&model.Topic{}).Where("id = ?", topicId).Updates(map[string]interface{}{"comment_count": gorm.Expr("comment_count + ?", 1), "last_comment_user_id": lastCommentUserId, "lastCommentTime": lastCommentTime}).Error; err != nil {
			return err
		}
		if err := dao.DB().Model(&model.TopicTag{}).Where("topic_id = ?", topicId).Updates(map[string]interface{}{"last_comment_time": lastCommentTime}).Error; err != nil {
			return err
		}
		return nil
	})
}

// rss
func (s *topicService) GenerateRss() {
	topics := dao.TopicDao.Find(sqlcnd.NewSqlCnd().Where("status = ?", model.StatusOk).Desc("id").Limit(1000))

	var items []*feeds.Item
	for _, topic := range topics {
		topicUrl := urls.TopicUrl(topic.ID)
		user := cache.UserCache.Get(topic.UserId)
		if user == nil {
			continue
		}
		item := &feeds.Item{
			Title:       topic.Title,
			Link:        &feeds.Link{Href: topicUrl},
			Description: util.GetMarkdownSummary(topic.Content),
			Author:      &feeds.Author{Name: user.Avatar, Email: user.Email.String},
			Created:     util.TimeFromTimestamp(topic.CreateTime),
		}
		items = append(items, item)
	}
	siteTitle := cache.SettingCache.GetValue(model.SettingSiteTitle)
	siteDescription := cache.SettingCache.GetValue(model.SettingSiteDescription)
	feed := &feeds.Feed{
		Title:       siteTitle,
		Link:        &feeds.Link{Href: viper.GetString("base.baseUrl")},
		Description: siteDescription,
		Author:      &feeds.Author{Name: siteTitle},
		Created:     time.Now(),
		Items:       items,
	}
	atom, err := feed.ToAtom()
	if err != nil {
		log.Error(err.Error())
	} else {
		_ = util.WriteString(path.Join(viper.GetString("base.static_path"), "topic_atom.xml"), atom, false)
	}

	rss, err := feed.ToRss()
	if err != nil {
		log.Error(err.Error())
	} else {
		_ = util.WriteString(path.Join(viper.GetString("base.static_path"), "topic_rss.xml"), rss, false)
	}
}

// 倒序扫描
func (s *topicService) ScanDesc(dateFrom, dateTo int64, cb ScanTopicCallback) {
	var cursor int64 = math.MaxInt64
	for {
		list := dao.TopicDao.Find(sqlcnd.NewSqlCnd().Lt("id", cursor).
			Gte("create_time", dateFrom).Lt("create_time", dateTo).Desc("id").Limit(1000))
		if list == nil || len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].ID
		cb(list)
	}
}
