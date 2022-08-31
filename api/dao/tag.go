package dao

import (
	"errors"
	"strings"

	"zendea/model"
	"zendea/util"
	"zendea/util/sqlcnd"
)

var TagDao = newTagDao()

func newTagDao() *tagDao {
	return &tagDao{}
}

type tagDao struct {
}

func (d *tagDao) Get(id int64) *model.Tag {
	ret := &model.Tag{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *tagDao) Take(where ...interface{}) *model.Tag {
	ret := &model.Tag{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (d *tagDao) Find(cnd *sqlcnd.SqlCnd) (list []model.Tag) {
	cnd.Find(db, &list)
	return
}

func (d *tagDao) FindOne(cnd *sqlcnd.SqlCnd) *model.Tag {
	ret := &model.Tag{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (d *tagDao) List(cnd *sqlcnd.SqlCnd) (list []model.Tag, paging *sqlcnd.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.Tag{})

	paging = &sqlcnd.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (d *tagDao) Create(t *model.Tag) (err error) {
	err = db.Create(t).Error
	return
}

func (d *tagDao) Update(t *model.Tag) (err error) {
	err = db.Save(t).Error
	return
}

func (d *tagDao) Updates(id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Tag{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (d *tagDao) UpdateColumn(id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Tag{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (d *tagDao) Delete(id int64) {
	db.Delete(&model.Tag{}, "id = ?", id)
}

func (d *tagDao) GetTagInIds(tagIds []int64) []model.Tag {
	if len(tagIds) == 0 {
		return nil
	}
	var tags []model.Tag
	db.Where("id in (?)", tagIds).Find(&tags)
	return tags
}

func (d *tagDao) GetByName(name string) *model.Tag {
	if len(name) == 0 {
		return nil
	}
	return d.Take("name = ?", name)
}

func (d *tagDao) GetOrCreate(name string) (*model.Tag, error) {
	if len(name) == 0 {
		return nil, errors.New("标签为空")
	}
	tag := d.GetByName(name)
	if tag != nil {
		return tag, nil
	} else {
		tag = &model.Tag{
			Name:       name,
			Status:     model.StatusOk,
			CreateTime: util.NowTimestamp(),
			UpdateTime: util.NowTimestamp(),
		}
		err := d.Create(tag)
		if err != nil {
			return nil, err
		}
		return tag, nil
	}
}

func (d *tagDao) GetOrCreates(tags []string) (tagIDs []int64) {
	for _, tagName := range tags {
		tagName = strings.TrimSpace(tagName)
		tag, err := d.GetOrCreate(tagName)
		if err == nil {
			tagIDs = append(tagIDs, tag.ID)
		}
	}
	return
}
