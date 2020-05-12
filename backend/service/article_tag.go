package service

import (
	"zendea/dao"
	"zendea/model"
	"zendea/util/sqlcnd"
)

var ArticleTagService = newArticleTagService()

func newArticleTagService() *articleTagService {
	return &articleTagService{}
}

type articleTagService struct {
}

func (s *articleTagService) Get(id int64) *model.ArticleTag {
	return dao.ArticleTagDao.Get(id)
}

func (s *articleTagService) Take(where ...interface{}) *model.ArticleTag {
	return dao.ArticleTagDao.Take(where...)
}

func (s *articleTagService) Find(cnd *sqlcnd.SqlCnd) []model.ArticleTag {
	return dao.ArticleTagDao.Find(cnd)
}

func (s *articleTagService) List(cnd *sqlcnd.SqlCnd) (list []model.ArticleTag, paging *sqlcnd.Paging) {
	return dao.ArticleTagDao.List(cnd)
}

func (s *articleTagService) Create(t *model.ArticleTag) error {
	return dao.ArticleTagDao.Create(t)
}

func (s *articleTagService) Update(t *model.ArticleTag) error {
	return dao.ArticleTagDao.Update(t)
}

func (s *articleTagService) Updates(id int64, columns map[string]interface{}) error {
	return dao.ArticleTagDao.Updates(id, columns)
}

func (s *articleTagService) UpdateColumn(id int64, name string, value interface{}) error {
	return dao.ArticleTagDao.UpdateColumn(id, name, value)
}

func (s *articleTagService) DeleteByArticleId(topicId int64) {
	dao.DB().Model(model.ArticleTag{}).Where("article_id = ?", topicId).UpdateColumn("status", model.StatusDeleted)
}
