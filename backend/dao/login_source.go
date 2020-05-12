package dao

import (
	"zendea/model"
	"zendea/util/sqlcnd"
)

var LoginSourceDao = newLoginSourceDao()

func newLoginSourceDao() *loginSourceDao {
	return &loginSourceDao{}
}

type loginSourceDao struct {
}

func (d *loginSourceDao) Get(id int64) *model.LoginSource {
	ret := &model.LoginSource{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *loginSourceDao) Take(where ...interface{}) *model.LoginSource {
	ret := &model.LoginSource{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (d *loginSourceDao) List(cnd *sqlcnd.SqlCnd) (list []model.LoginSource, paging *sqlcnd.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.LoginSource{})

	paging = &sqlcnd.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (d *loginSourceDao) Create(t *model.LoginSource) (err error) {
	err = db.Create(t).Error
	return
}

func (d *loginSourceDao) Update(t *model.LoginSource) (err error) {
	err = db.Save(t).Error
	return
}

func (d *loginSourceDao) Updates(id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.LoginSource{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (d *loginSourceDao) UpdateColumn(id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.LoginSource{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (d *loginSourceDao) Delete(id int64) {
	db.Delete(&model.LoginSource{}, "id = ?", id)
}
