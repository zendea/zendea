package dao

import (
	"zendea/model"
	"zendea/util/sqlcnd"
)

var TopicLikeDao = newTopicLikeDao()

func newTopicLikeDao() *topicLikeDao {
	return &topicLikeDao{}
}

type topicLikeDao struct {
}

func (d *topicLikeDao) Get(id int64) *model.TopicLike {
	ret := &model.TopicLike{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *topicLikeDao) Take(where ...interface{}) *model.TopicLike {
	ret := &model.TopicLike{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (d *topicLikeDao) Find(cnd *sqlcnd.SqlCnd) (list []model.TopicLike) {
	cnd.Find(db, &list)
	return
}

func (d *topicLikeDao) FindOne(cnd *sqlcnd.SqlCnd) *model.TopicLike {
	ret := &model.TopicLike{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (d *topicLikeDao) List(cnd *sqlcnd.SqlCnd) (list []model.TopicLike, paging *sqlcnd.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.TopicLike{})

	paging = &sqlcnd.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (d *topicLikeDao) Create(t *model.TopicLike) (err error) {
	err = db.Create(t).Error
	return
}

func (d *topicLikeDao) Update(t *model.TopicLike) (err error) {
	err = db.Save(t).Error
	return
}

func (d *topicLikeDao) Updates(id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.TopicLike{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (d *topicLikeDao) UpdateColumn(id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.TopicLike{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (d *topicLikeDao) Delete(id int64) {
	db.Delete(&model.TopicLike{}, "id = ?", id)
}
