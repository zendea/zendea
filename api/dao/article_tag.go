package dao

import (
	"zendea/model"
	"zendea/util"
	"zendea/util/sqlcnd"
)

var ArticleTagDao = newArticleTagDao()

func newArticleTagDao() *articleTagDao {
	return &articleTagDao{}
}

type articleTagDao struct {
}

func (d *articleTagDao) Get(id int64) *model.ArticleTag {
	ret := &model.ArticleTag{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *articleTagDao) Take(where ...interface{}) *model.ArticleTag {
	ret := &model.ArticleTag{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (d *articleTagDao) Find(cnd *sqlcnd.SqlCnd) (list []model.ArticleTag) {
	cnd.Find(db, &list)
	return
}

func (d *articleTagDao) List(cnd *sqlcnd.SqlCnd) (list []model.ArticleTag, paging *sqlcnd.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.ArticleTag{})

	paging = &sqlcnd.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (d *articleTagDao) Create(t *model.ArticleTag) (err error) {
	err = db.Create(t).Error
	return
}

func (d *articleTagDao) Update(t *model.ArticleTag) (err error) {
	err = db.Save(t).Error
	return
}

func (d *articleTagDao) Updates(id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.ArticleTag{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (d *articleTagDao) UpdateColumn(id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.ArticleTag{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (d *articleTagDao) Delete(id int64) {
	db.Delete(&model.ArticleTag{}, "id = ?", id)
}

func (d *articleTagDao) AddArticleTags(articleId int64, tagIds []int64) {
	if articleId <= 0 || len(tagIds) == 0 {
		return
	}

	for _, tagId := range tagIds {
		_ = d.Create(&model.ArticleTag{
			ArticleId:  articleId,
			TagId:      tagId,
			CreateTime: util.NowTimestamp(),
		})
	}
}

func (d *articleTagDao) DeleteArticleTags(articleId int64) {
	if articleId <= 0 {
		return
	}
	db.Where("article_id = ?", articleId).Delete(model.ArticleTag{})
}

func (d *articleTagDao) FindByArticleId(articleId int64) []model.ArticleTag {
	return d.Find(sqlcnd.NewSqlCnd().Where("article_id = ?", articleId))
}
