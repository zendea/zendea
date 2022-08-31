package dao

import (
	"zendea/model"
	"zendea/util/sqlcnd"
)

var SitemapDao = newSitemapDao()

func newSitemapDao() *sitemapDao {
	return &sitemapDao{}
}

type sitemapDao struct {
}

func (d *sitemapDao) Get(id int64) *model.Sitemap {
	ret := &model.Sitemap{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *sitemapDao) Take(where ...interface{}) *model.Sitemap {
	ret := &model.Sitemap{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (d *sitemapDao) Find(cnd *sqlcnd.SqlCnd) (list []model.Sitemap) {
	cnd.Find(db, &list)
	return
}

func (d *sitemapDao) FindOne(cnd *sqlcnd.SqlCnd) *model.Sitemap {
	ret := &model.Sitemap{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (d *sitemapDao) List(cnd *sqlcnd.SqlCnd) (list []model.Sitemap, paging *sqlcnd.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.Sitemap{})

	paging = &sqlcnd.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (d *sitemapDao) Create(t *model.Sitemap) (err error) {
	err = db.Create(t).Error
	return
}

func (d *sitemapDao) Update(t *model.Sitemap) (err error) {
	err = db.Save(t).Error
	return
}

func (d *sitemapDao) Updates(id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Sitemap{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (d *sitemapDao) UpdateColumn(id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Sitemap{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (d *sitemapDao) Delete(id int64) {
	db.Delete(&model.Sitemap{}, "id = ?", id)
}
