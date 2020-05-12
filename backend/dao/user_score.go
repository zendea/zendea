package dao

import (
	"zendea/model"
	"zendea/util/sqlcnd"
)

var UserScoreDao = newUserScoreDao()

func newUserScoreDao() *userScoreDao {
	return &userScoreDao{}
}

type userScoreDao struct {
}

func (d *userScoreDao) Get(id int64) *model.UserScore {
	ret := &model.UserScore{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *userScoreDao) Take(where ...interface{}) *model.UserScore {
	ret := &model.UserScore{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (d *userScoreDao) Find(cnd *sqlcnd.SqlCnd) (list []model.UserScore) {
	cnd.Find(db, &list)
	return
}

func (d *userScoreDao) FindOne(cnd *sqlcnd.SqlCnd) *model.UserScore {
	ret := &model.UserScore{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (d *userScoreDao) List(cnd *sqlcnd.SqlCnd) (list []model.UserScore, paging *sqlcnd.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.UserScore{})

	paging = &sqlcnd.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (d *userScoreDao) Create(t *model.UserScore) (err error) {
	err = db.Create(t).Error
	return
}

func (d *userScoreDao) Update(t *model.UserScore) (err error) {
	err = db.Save(t).Error
	return
}

func (d *userScoreDao) Updates(id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.UserScore{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (d *userScoreDao) UpdateColumn(id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.UserScore{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (d *userScoreDao) Delete(id int64) {
	db.Delete(&model.UserScore{}, "id = ?", id)
}
