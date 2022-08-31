package dao

import (
	"zendea/model"
	"zendea/util/sqlcnd"
)

var SectionDao = newSectionDao()

func newSectionDao() *sectionDao {
	return &sectionDao{}
}

type sectionDao struct {
}

func (d *sectionDao) Get(id int64) *model.Section {
	ret := &model.Section{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *sectionDao) Take(where ...interface{}) *model.Section {
	ret := &model.Section{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (d *sectionDao) Find(cnd *sqlcnd.SqlCnd) (list []model.Section) {
	cnd.Find(db, &list)
	return
}

func (d *sectionDao) FindOne(cnd *sqlcnd.SqlCnd) *model.Section {
	ret := &model.Section{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (d *sectionDao) List(cnd *sqlcnd.SqlCnd) (list []model.Section, paging *sqlcnd.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.Section{})

	paging = &sqlcnd.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (d *sectionDao) Count(cnd *sqlcnd.SqlCnd) int {
	return cnd.Count(db, &model.Section{})
}

func (d *sectionDao) Create(t *model.Section) (err error) {
	err = db.Create(t).Error
	return
}

func (d *sectionDao) Update(t *model.Section) (err error) {
	err = db.Save(t).Error
	return
}

func (d *sectionDao) Updates(id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Section{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (d *sectionDao) UpdateColumn(id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Section{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (d *sectionDao) Delete(id int64) {
	db.Delete(&model.Section{}, "id = ?", id)
}
