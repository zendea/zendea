package service

import (
	"math"
	"path"
	"time"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/gorilla/feeds"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"

	"zendea/cache"
	"zendea/dao"
	"zendea/form"
	"zendea/model"
	"zendea/util"
	"zendea/util/log"
	"zendea/util/sqlcnd"
	"zendea/util/urls"
)

type ScanArticleCallback func(articles []model.Article)

var ArticleService = newArticleService()

func newArticleService() *articleService {
	return &articleService{}
}

type articleService struct {
}

func (s *articleService) Get(id int64) *model.Article {
	return dao.ArticleDao.Get(id)
}

func (s *articleService) Find(cnd *sqlcnd.SqlCnd) []model.Article {
	return dao.ArticleDao.Find(cnd)
}

func (s *articleService) List(cnd *sqlcnd.SqlCnd) (list []model.Article, paging *sqlcnd.Paging) {
	return dao.ArticleDao.List(cnd)
}

// Create 发表文章
func (s *articleService) Create(dto form.ArticleCreateForm) (*model.Article, error) {
	article := &model.Article{
		UserId:      dto.UserID,
		Title:       dto.Title,
		Summary:     dto.Summary,
		Content:     dto.Content,
		ContentType: model.ContentTypeMarkdown,
		Status:      model.StatusOk,
		Share:       false,
		SourceUrl:   "",
		CreateTime:  util.NowTimestamp(),
		UpdateTime:  util.NowTimestamp(),
	}

	err := dao.Tx(dao.DB(), func(tx *gorm.DB) error {
		tagIDs := dao.TagDao.GetOrCreates(util.ParseTagsToArray(dto.Tags))
		err := dao.ArticleDao.Create(article)
		if err != nil {
			return err
		}
		dao.ArticleTagDao.AddArticleTags(article.ID, tagIDs)
		return nil
	})
	return article, err
}

// Update 编辑文章
func (s *articleService) Update(dto form.ArticleUpdateForm) error {
	err := dao.Tx(dao.DB(), func(tx *gorm.DB) error {
		err := dao.ArticleDao.Updates(dto.ID, map[string]interface{}{
			"title":       dto.Title,
			"content":     dto.Content,
			"update_time": util.NowTimestamp(),
		})
		if err != nil {
			return err
		}
		tagIds := dao.TagDao.GetOrCreates(util.ParseTagsToArray(dto.Tags)) // 创建文章对应标签
		dao.ArticleTagDao.DeleteArticleTags(dto.ID)                        // 先删掉所有的标签
		dao.ArticleTagDao.AddArticleTags(dto.ID, tagIds)                   // 然后重新添加标签
		return nil
	})
	cache.ArticleTagCache.Invalidate(dto.ID)
	return err
}

func (s *articleService) Delete(id int64) error {
	err := dao.ArticleDao.UpdateColumn(id, "status", model.StatusDeleted)
	if err == nil {
		// 删掉标签文章
		ArticleTagService.DeleteByArticleId(id)
	}
	return err
}

// 根据文章编号批量获取文章
func (s *articleService) GetArticleInIds(articleIds []int64) []model.Article {
	if len(articleIds) == 0 {
		return nil
	}
	var articles []model.Article
	dao.DB().Where("id in (?)", articleIds).Find(&articles)
	return articles
}

// 文章列表
func (s *articleService) GetArticles(cursor int64) (articles []model.Article, nextCursor int64) {
	cnd := sqlcnd.NewSqlCnd().Eq("status", model.StatusOk).Desc("id").Limit(20)
	if cursor > 0 {
		cnd.Lt("id", cursor)
	}
	articles = dao.ArticleDao.Find(cnd)
	if len(articles) > 0 {
		nextCursor = articles[len(articles)-1].ID
	} else {
		nextCursor = cursor
	}
	return
}

// GetArticleTags 获取文章对应的标签
func (s *articleService) GetArticleTags(articleId int64) []model.Tag {
	articleTags := dao.ArticleTagDao.Find(sqlcnd.NewSqlCnd().Where("article_id = ?", articleId))
	var tagIds []int64
	for _, articleTag := range articleTags {
		tagIds = append(tagIds, articleTag.TagId)
	}
	return cache.TagCache.GetList(tagIds)
}

// GetTagArticles 标签文章列表
func (s *articleService) GetTagArticles(tagId int64, cursor int64) (articles []model.Article, nextCursor int64) {
	cnd := sqlcnd.NewSqlCnd().Eq("tag_id", tagId).Eq("status", model.StatusOk).Desc("id").Limit(20)
	if cursor > 0 {
		cnd.Lt("id", cursor)
	}
	nextCursor = cursor
	articleTags := dao.ArticleTagDao.Find(cnd)
	if len(articleTags) > 0 {
		var articleIds []int64
		for _, articleTag := range articleTags {
			articleIds = append(articleIds, articleTag.ArticleId)
			nextCursor = articleTag.ID
		}
		articles = s.GetArticleInIds(articleIds)
	}
	return
}

// 相关文章
func (s *articleService) GetRelatedArticles(articleId int64) []model.Article {
	tagIds := cache.ArticleTagCache.Get(articleId)
	if len(tagIds) == 0 {
		return nil
	}
	var articleTags []model.ArticleTag
	dao.DB().Where("tag_id in (?)", tagIds).Limit(30).Find(&articleTags)

	set := hashset.New()
	if len(articleTags) > 0 {
		for _, articleTag := range articleTags {
			set.Add(articleTag.ArticleId)
		}
	}

	var articleIds []int64
	for i, articleId := range set.Values() {
		if i < 10 {
			articleIds = append(articleIds, articleId.(int64))
		}
	}

	return s.GetArticleInIds(articleIds)
}

// 最新文章
func (s *articleService) GetUserNewestArticles(userId int64) []model.Article {
	return dao.ArticleDao.Find(sqlcnd.NewSqlCnd().Where("user_id = ? and status = ?",
		userId, model.StatusOk).Desc("id").Limit(10))
}

// 倒序扫描
func (s *articleService) ScanDesc(dateFrom, dateTo int64, cb ScanArticleCallback) {
	var cursor int64 = math.MaxInt64
	for {
		list := dao.ArticleDao.Find(sqlcnd.NewSqlCnd("id", "status", "create_time", "update_time").
			Lt("id", cursor).Gte("create_time", dateFrom).Lt("create_time", dateTo).Desc("id").Limit(1000))
		if list == nil || len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].ID
		cb(list)
	}
}

// rss
func (s *articleService) GenerateRss() {
	articles := dao.ArticleDao.Find(sqlcnd.NewSqlCnd().Where("status = ?", model.StatusOk).Desc("id").Limit(1000))

	var items []*feeds.Item
	for _, article := range articles {
		articleUrl := urls.ArticleUrl(article.ID)
		user := cache.UserCache.Get(article.UserId)
		if user == nil {
			continue
		}
		description := ""
		if article.ContentType == model.ContentTypeMarkdown {
			description = util.GetMarkdownSummary(article.Content)
		} else {
			description = util.GetHtmlSummary(article.Content)
		}
		item := &feeds.Item{
			Title:       article.Title,
			Link:        &feeds.Link{Href: articleUrl},
			Description: description,
			Author:      &feeds.Author{Name: user.Avatar, Email: user.Email.String},
			Created:     util.TimeFromTimestamp(article.CreateTime),
		}
		items = append(items, item)
	}

	siteTitle := cache.SettingCache.GetValue(model.SettingSiteTitle)
	siteDescription := cache.SettingCache.GetValue(model.SettingSiteDescription)
	feed := &feeds.Feed{
		Title:       siteTitle,
		Link:        &feeds.Link{Href: viper.GetString("base.url")},
		Description: siteDescription,
		Author:      &feeds.Author{Name: siteTitle},
		Created:     time.Now(),
		Items:       items,
	}
	atom, err := feed.ToAtom()
	if err != nil {
		log.Error(err.Error())
	} else {
		_ = util.WriteString(path.Join(viper.GetString("base.static_path"), "atom.xml"), atom, false)
	}

	rss, err := feed.ToRss()
	if err != nil {
		log.Error(err.Error())
	} else {
		_ = util.WriteString(path.Join(viper.GetString("base.static_path"), "rss.xml"), rss, false)
	}
}
