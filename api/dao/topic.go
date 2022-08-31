package dao

import (
	"zendea/model"
	"zendea/util/sqlcnd"
)

var TopicDao = newTopicDao()

func newTopicDao() *topicDao {
	return &topicDao{}
}

type topicDao struct {
}

func (d *topicDao) Get(id int64) *model.Topic {
	ret := &model.Topic{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *topicDao) Take(where ...interface{}) *model.Topic {
	ret := &model.Topic{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (d *topicDao) Find(cnd *sqlcnd.SqlCnd) (list []model.Topic) {
	cnd.Find(db, &list)
	return
}

func (d *topicDao) FindOne(cnd *sqlcnd.SqlCnd) *model.Topic {
	ret := &model.Topic{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (d *topicDao) List(cnd *sqlcnd.SqlCnd) (list []model.Topic, paging *sqlcnd.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.Topic{})

	paging = &sqlcnd.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (d *topicDao) Count(cnd *sqlcnd.SqlCnd) int {
	return cnd.Count(db, &model.Topic{})
}

func (d *topicDao) Create(t *model.Topic) (err error) {
	err = db.Create(t).Error
	return
}

func (d *topicDao) Update(t *model.Topic) (err error) {
	err = db.Save(t).Error
	return
}

func (d *topicDao) Updates(id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Topic{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (d *topicDao) UpdateColumn(id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Topic{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (d *topicDao) Delete(id int64) {
	db.Delete(&model.Topic{}, "id = ?", id)
}
