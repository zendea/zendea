package dao

import (
	"zendea/model"
	"zendea/util"
	"zendea/util/sqlcnd"
)

var TopicTagDao = newTopicTagDao()

func newTopicTagDao() *topicTagDao {
	return &topicTagDao{}
}

type topicTagDao struct {
}

func (d *topicTagDao) Get(id int64) *model.TopicTag {
	ret := &model.TopicTag{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *topicTagDao) Take(where ...interface{}) *model.TopicTag {
	ret := &model.TopicTag{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (d *topicTagDao) Find(cnd *sqlcnd.SqlCnd) (list []model.TopicTag) {
	cnd.Find(db, &list)
	return
}

func (d *topicTagDao) FindOne(cnd *sqlcnd.SqlCnd) *model.TopicTag {
	ret := &model.TopicTag{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (d *topicTagDao) List(cnd *sqlcnd.SqlCnd) (list []model.TopicTag, paging *sqlcnd.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.TopicTag{})

	paging = &sqlcnd.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (d *topicTagDao) Create(t *model.TopicTag) (err error) {
	err = db.Create(t).Error
	return
}

func (d *topicTagDao) Update(t *model.TopicTag) (err error) {
	err = db.Save(t).Error
	return
}

func (d *topicTagDao) Updates(id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.TopicTag{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (d *topicTagDao) UpdateColumn(id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.TopicTag{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (d *topicTagDao) Delete(id int64) {
	db.Delete(&model.TopicTag{}, "id = ?", id)
}

func (d *topicTagDao) AddTopicTags(topicId int64, tagIds []int64) {
	if topicId <= 0 || len(tagIds) == 0 {
		return
	}
	for _, tagId := range tagIds {
		_ = d.Create(&model.TopicTag{
			TopicId:    topicId,
			TagId:      tagId,
			CreateTime: util.NowTimestamp(),
		})
	}
}

func (d *topicTagDao) DeleteTopicTags(topicId int64) {
	if topicId <= 0 {
		return
	}
	db.Where("topic_id = ?", topicId).Delete(model.TopicTag{})
}
