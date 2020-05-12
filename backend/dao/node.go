package dao

import (
	"zendea/model"
	"zendea/util/sqlcnd"
)

var NodeDao = newNodeDao()

func newNodeDao() *nodeDao {
	return &nodeDao{}
}

type nodeDao struct {
}

func (d *nodeDao) Get(id int64) *model.Node {
	ret := &model.Node{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *nodeDao) Take(where ...interface{}) *model.Node {
	ret := &model.Node{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (d *nodeDao) Find(cnd *sqlcnd.SqlCnd) (list []model.Node) {
	cnd.Find(db, &list)
	return
}

func (d *nodeDao) FindOne(cnd *sqlcnd.SqlCnd) *model.Node {
	ret := &model.Node{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (d *nodeDao) List(cnd *sqlcnd.SqlCnd) (list []model.Node, paging *sqlcnd.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.Node{})

	paging = &sqlcnd.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (d *nodeDao) Create(t *model.Node) (err error) {
	err = db.Create(t).Error
	return
}

func (d *nodeDao) Update(t *model.Node) (err error) {
	err = db.Save(t).Error
	return
}

func (d *nodeDao) Updates(id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Node{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (d *nodeDao) UpdateColumn(id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Node{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (d *nodeDao) Delete(id int64) {
	db.Delete(&model.Node{}, "id = ?", id)
}
