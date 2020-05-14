package service

import (
	"errors"

	"github.com/jinzhu/gorm"

	"zendea/dao"
	"zendea/model"
	"zendea/util"
	"zendea/util/sqlcnd"
)

var TopicLikeService = newTopicLikeService()

func newTopicLikeService() *topicLikeService {
	return &topicLikeService{}
}

type topicLikeService struct {
}

func (s *topicLikeService) Get(id int64) *model.TopicLike {
	return dao.TopicLikeDao.Get(id)
}

func (s *topicLikeService) Take(where ...interface{}) *model.TopicLike {
	return dao.TopicLikeDao.Take(where...)
}

func (s *topicLikeService) Find(cnd *sqlcnd.SqlCnd) []model.TopicLike {
	return dao.TopicLikeDao.Find(cnd)
}

func (s *topicLikeService) FindOne(cnd *sqlcnd.SqlCnd) *model.TopicLike {
	return dao.TopicLikeDao.FindOne(cnd)
}

func (s *topicLikeService) List(cnd *sqlcnd.SqlCnd) (list []model.TopicLike, paging *sqlcnd.Paging) {
	return dao.TopicLikeDao.List(cnd)
}

func (s *topicLikeService) Create(t *model.TopicLike) error {
	return dao.TopicLikeDao.Create(t)
}

func (s *topicLikeService) Update(t *model.TopicLike) error {
	return dao.TopicLikeDao.Update(t)
}

func (s *topicLikeService) Updates(id int64, columns map[string]interface{}) error {
	return dao.TopicLikeDao.Updates(id, columns)
}

func (s *topicLikeService) UpdateColumn(id int64, name string, value interface{}) error {
	return dao.TopicLikeDao.UpdateColumn(id, name, value)
}

func (s *topicLikeService) Delete(id int64) {
	dao.TopicLikeDao.Delete(id)
}

// 统计数量
func (s *topicLikeService) Count(topicId int64) int64 {
	var count int64 = 0
	dao.DB().Model(&model.TopicLike{}).Where("topic_id = ?", topicId).Count(&count)
	return count
}

// 最近点赞
func (s *topicLikeService) Recent(topicId int64, count int) []model.TopicLike {
	return s.Find(sqlcnd.NewSqlCnd().Eq("topic_id", topicId).Desc("id").Limit(count))
}

func (s *topicLikeService) Like(userId int64, topicId int64) error {
	topic := dao.TopicDao.Get(topicId)
	if topic == nil || topic.Status != model.StatusOk {
		return errors.New("话题不存在")
	}

	// 判断是否已经点赞了
	topicLike := dao.TopicLikeDao.Take("user_id = ? and topic_id = ?", userId, topicId)
	if topicLike != nil {
		return errors.New("已点赞")
	}

	return dao.Tx(func(tx *gorm.DB) error {
		// 点赞
		topicLike := &model.TopicLike{
			UserId:     userId,
			TopicId:    topicId,
			CreateTime: util.NowTimestamp(),
		}
		err := dao.TopicLikeDao.Create(topicLike)
		if err != nil {
			return err
		}
		// 发送点赞通知
		NotificationService.SendTopicLikeNotification(topicLike)

		return dao.DB().Model(&topic).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error
	})
}
