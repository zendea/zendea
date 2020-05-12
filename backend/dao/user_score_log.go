package dao

import (
	"zendea/model"
	"zendea/util/sqlcnd"
)

var UserScoreLogDao = newUserScoreLogDao()

func newUserScoreLogDao() *userScoreLogDao {
	return &userScoreLogDao{}
}

type userScoreLogDao struct {
}

func (d *userScoreLogDao) Get(id int64) *model.UserScoreLog {
	ret := &model.UserScoreLog{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *userScoreLogDao) Take(where ...interface{}) *model.UserScoreLog {
	ret := &model.UserScoreLog{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (d *userScoreLogDao) Find(cnd *sqlcnd.SqlCnd) (list []model.UserScoreLog) {
	cnd.Find(db, &list)
	return
}

func (d *userScoreLogDao) FindOne(cnd *sqlcnd.SqlCnd) *model.UserScoreLog {
	ret := &model.UserScoreLog{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (d *userScoreLogDao) List(cnd *sqlcnd.SqlCnd) (list []model.UserScoreLog, paging *sqlcnd.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.UserScoreLog{})

	paging = &sqlcnd.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (d *userScoreLogDao) Create(t *model.UserScoreLog) (err error) {
	err = db.Create(t).Error
	return
}

func (d *userScoreLogDao) Update(t *model.UserScoreLog) (err error) {
	err = db.Save(t).Error
	return
}

func (d *userScoreLogDao) Updates(id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.UserScoreLog{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (d *userScoreLogDao) UpdateColumn(id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.UserScoreLog{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (d *userScoreLogDao) Delete(id int64) {
	db.Delete(&model.UserScoreLog{}, "id = ?", id)
}
