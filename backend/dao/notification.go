package dao

import (
	"zendea/model"
	"zendea/util/sqlcnd"
)

var NotificationDao = newNotificationDao()

func newNotificationDao() *notificationDao {
	return &notificationDao{}
}

type notificationDao struct {
}

func (d *notificationDao) Get(id int64) *model.Notification {
	ret := &model.Notification{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *notificationDao) Take(where ...interface{}) *model.Notification {
	ret := &model.Notification{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (d *notificationDao) Find(cnd *sqlcnd.SqlCnd) (list []model.Notification) {
	cnd.Find(db, &list)
	return
}

func (d *notificationDao) FindOne(cnd *sqlcnd.SqlCnd) *model.Notification {
	ret := &model.Notification{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (d *notificationDao) List(cnd *sqlcnd.SqlCnd) (list []model.Notification, paging *sqlcnd.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.Notification{})

	paging = &sqlcnd.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (d *notificationDao) Create(t *model.Notification) (err error) {
	err = db.Create(t).Error
	return
}

func (d *notificationDao) Update(t *model.Notification) (err error) {
	err = db.Save(t).Error
	return
}

func (d *notificationDao) Updates(id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Notification{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (d *notificationDao) UpdateColumn(id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Notification{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (s *notificationDao) GetUnReadCount(userId int64) (count int64) {
	db.Model(&model.Notification{}).Where("user_id = ? and status = ?", userId, model.NotificationStatusUnread).Count(&count)
	return
}

func (d *notificationDao) UpdateStatusBatch(userId int64) (err error) {
	err = db.Model(&model.Notification{}).Where("user_id = ? and status = ?", userId, model.NotificationStatusUnread).Updates(model.Notification{Status: model.NotificationStatusReaded}).Error
	return
}

func (d *notificationDao) Delete(id int64) {
	db.Delete(&model.Notification{}, "id = ?", id)
}
