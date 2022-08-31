package service

import (
	"errors"

	"zendea/dao"
	"zendea/model"
	"zendea/util"
	"zendea/util/sqlcnd"
)

var FavoriteService = newFavoriteService()

func newFavoriteService() *favoriteService {
	return &favoriteService{}
}

type favoriteService struct {
}

func (s *favoriteService) Get(id int64) *model.Favorite {
	return dao.FavoriteDao.Get(id)
}

func (s *favoriteService) Take(where ...interface{}) *model.Favorite {
	return dao.FavoriteDao.Take(where...)
}

func (s *favoriteService) Find(cnd *sqlcnd.SqlCnd) []model.Favorite {
	return dao.FavoriteDao.Find(cnd)
}

func (s *favoriteService) FindOne(cnd *sqlcnd.SqlCnd) *model.Favorite {
	return dao.FavoriteDao.FindOne(cnd)
}

func (s *favoriteService) List(cnd *sqlcnd.SqlCnd) (list []model.Favorite, paging *sqlcnd.Paging) {
	return dao.FavoriteDao.List(cnd)
}

func (s *favoriteService) Create(t *model.Favorite) error {
	return dao.FavoriteDao.Create(t)
}

func (s *favoriteService) Update(t *model.Favorite) error {
	return dao.FavoriteDao.Update(t)
}

func (s *favoriteService) Updates(id int64, columns map[string]interface{}) error {
	return dao.FavoriteDao.Updates(id, columns)
}

func (s *favoriteService) UpdateColumn(id int64, name string, value interface{}) error {
	return dao.FavoriteDao.UpdateColumn(id, name, value)
}

func (s *favoriteService) Delete(id int64) {
	dao.FavoriteDao.Delete(id)
}

func (s *favoriteService) GetBy(userId int64, entityType string, entityId int64) *model.Favorite {
	return dao.FavoriteDao.Take("user_id = ? and entity_type = ? and entity_id = ?",
		userId, entityType, entityId)
}

// 收藏文章
func (s *favoriteService) AddArticleFavorite(userId, articleId int64) error {
	article := dao.ArticleDao.Get(articleId)
	if article == nil || article.Status != model.StatusOk {
		return errors.New("收藏的文章不存在")
	}
	return s.addFavorite(userId, model.EntityTypeArticle, articleId)
}

// 收藏主题
func (s *favoriteService) AddTopicFavorite(userId, topicId int64) error {
	topic := dao.TopicDao.Get(topicId)
	if topic == nil || topic.Status != model.StatusOk {
		return errors.New("收藏的话题不存在")
	}
	return s.addFavorite(userId, model.EntityTypeTopic, topicId)
}

func (s *favoriteService) addFavorite(userId int64, entityType string, entityId int64) error {
	temp := s.GetBy(userId, entityType, entityId)
	if temp != nil { // 已经收藏
		return nil
	}
	return dao.FavoriteDao.Create(&model.Favorite{
		UserId:     userId,
		EntityType: entityType,
		EntityId:   entityId,
		CreateTime: util.NowTimestamp(),
	})
}
