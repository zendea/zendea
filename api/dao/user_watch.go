package dao

import (
	"zendea/model"
	"zendea/util/sqlcnd"
)

var UserWatchDao = newUserWatchDao()

func newUserWatchDao() *userWatchDao {
	return &userWatchDao{}
}

type userWatchDao struct {
}

func (d *userWatchDao) Get(id int64) *model.UserWatch {
	ret := &model.UserWatch{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *userWatchDao) Take(where ...interface{}) *model.UserWatch {
	ret := &model.UserWatch{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (d *userWatchDao) Find(cnd *sqlcnd.SqlCnd) (list []model.UserWatch) {
	cnd.Find(db, &list)
	return
}

func (d *userWatchDao) FindOne(cnd *sqlcnd.SqlCnd) *model.UserWatch {
	ret := &model.UserWatch{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (d *userWatchDao) List(cnd *sqlcnd.SqlCnd) (list []model.UserWatch, paging *sqlcnd.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.UserWatch{})

	paging = &sqlcnd.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (d *userWatchDao) Create(t *model.UserWatch) (err error) {
	err = db.Create(t).Error
	return
}

func (d *userWatchDao) Update(t *model.UserWatch) (err error) {
	err = db.Save(t).Error
	return
}

func (d *userWatchDao) Updates(id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.UserWatch{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (d *userWatchDao) UpdateColumn(id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.UserWatch{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (d *userWatchDao) Delete(id int64) {
	db.Delete(&model.UserWatch{}, "id = ?", id)
}
