package dao

import (
	"zendea/model"
	"zendea/util/sqlcnd"
)

var SettingDao = newSettingDao()

func newSettingDao() *settingDao {
	return &settingDao{}
}

type settingDao struct {
}

func (d *settingDao) Get(id int64) *model.Setting {
	ret := &model.Setting{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *settingDao) Take(where ...interface{}) *model.Setting {
	ret := &model.Setting{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (d *settingDao) Find(cnd *sqlcnd.SqlCnd) (list []model.Setting) {
	cnd.Find(db, &list)
	return
}

func (d *settingDao) FindOne(cnd *sqlcnd.SqlCnd) *model.Setting {
	ret := &model.Setting{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (d *settingDao) List(cnd *sqlcnd.SqlCnd) (list []model.Setting, paging *sqlcnd.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.Setting{})

	paging = &sqlcnd.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (d *settingDao) Create(t *model.Setting) (err error) {
	err = db.Create(t).Error
	return
}

func (d *settingDao) Update(t *model.Setting) (err error) {
	err = db.Save(t).Error
	return
}

func (d *settingDao) Updates(id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Setting{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (d *settingDao) UpdateColumn(id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Setting{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (d *settingDao) Delete(id int64) {
	db.Delete(&model.Setting{}, "id = ?", id)
}

func (d *settingDao) GetByKey(key string) *model.Setting {
	if len(key) == 0 {
		return nil
	}
	return d.Take("`key` = ?", key)
}
