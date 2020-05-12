package dao

import (
	"zendea/model"
	"zendea/util/sqlcnd"
)

var UserDao = newUserDao()

func newUserDao() *userDao {
	return &userDao{}
}

type userDao struct {
}

func (d *userDao) Get(id int64) *model.User {
	ret := &model.User{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *userDao) Take(where ...interface{}) *model.User {
	ret := &model.User{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (d *userDao) Find(cnd *sqlcnd.SqlCnd) (list []model.User) {
	cnd.Find(db, &list)
	return
}

func (d *userDao) FindOne(cnd *sqlcnd.SqlCnd) *model.User {
	ret := &model.User{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (d *userDao) List(cnd *sqlcnd.SqlCnd) (list []model.User, paging *sqlcnd.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.User{})

	paging = &sqlcnd.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (d *userDao) Count(cnd *sqlcnd.SqlCnd) int {
	return cnd.Count(db, &model.User{})
}

func (d *userDao) Create(t *model.User) (err error) {
	err = db.Create(t).Error
	return
}

func (d *userDao) Update(t *model.User) (err error) {
	err = db.Save(t).Error
	return
}

func (d *userDao) Updates(id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.User{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (d *userDao) UpdateColumn(id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.User{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (d *userDao) Delete(id int64) {
	db.Delete(&model.User{}, "id = ?", id)
}

func (d *userDao) GetByEmail(email string) *model.User {
	return d.Take("email = ?", email)
}

func (d *userDao) GetByUsername(username string) *model.User {
	return d.Take("username = ?", username)
}
