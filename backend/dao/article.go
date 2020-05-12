package dao

import (
	"zendea/model"
	"zendea/util/sqlcnd"
)

var ArticleDao = newArticleDao()

func newArticleDao() *articleDao {
	return &articleDao{}
}

type articleDao struct {
}

// Get returns article by given ID.
func (d *articleDao) Get(id int64) *model.Article {
	ret := &model.Article{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *articleDao) Find(cnd *sqlcnd.SqlCnd) (list []model.Article) {
	cnd.Find(db, &list)
	return
}

func (d *articleDao) List(cnd *sqlcnd.SqlCnd) (list []model.Article, paging *sqlcnd.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.Article{})

	paging = &sqlcnd.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (d *articleDao) Create(t *model.Article) (err error) {
	err = db.Create(t).Error
	return
}

func (d *articleDao) Update(t *model.Article) (err error) {
	err = db.Save(t).Error
	return
}

func (d *articleDao) Updates(id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Article{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (d *articleDao) UpdateColumn(id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Article{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (d *articleDao) Delete(id int64) {
	db.Delete(&model.Article{}, "id = ?", id)
}
