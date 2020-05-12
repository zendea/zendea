package service

import (
	"time"

	"github.com/spf13/viper"

	"zendea/dao"
	"zendea/model"
	"zendea/util"
	"zendea/util/log"
	"zendea/util/sitemap"
	"zendea/util/sqlcnd"
	"zendea/util/urls"
)

const (
	sitemapPath = "sitemap"
)

var SitemapService = newSitemapService()

func newSitemapService() *sitemapService {
	return &sitemapService{}
}

type sitemapService struct {
	building bool // is in building
}

func (s *sitemapService) Get(id int64) *model.Sitemap {
	return dao.SitemapDao.Get(id)
}

func (s *sitemapService) Take(where ...interface{}) *model.Sitemap {
	return dao.SitemapDao.Take(where...)
}

func (s *sitemapService) Find(cnd *sqlcnd.SqlCnd) []model.Sitemap {
	return dao.SitemapDao.Find(cnd)
}

func (s *sitemapService) FindOne(cnd *sqlcnd.SqlCnd) *model.Sitemap {
	return dao.SitemapDao.FindOne(cnd)
}

func (s *sitemapService) List(cnd *sqlcnd.SqlCnd) (list []model.Sitemap, paging *sqlcnd.Paging) {
	return dao.SitemapDao.List(cnd)
}

func (s *sitemapService) Create(t *model.Sitemap) error {
	return dao.SitemapDao.Create(t)
}

func (s *sitemapService) Update(t *model.Sitemap) error {
	return dao.SitemapDao.Update(t)
}

func (s *sitemapService) Updates(id int64, columns map[string]interface{}) error {
	return dao.SitemapDao.Updates(id, columns)
}

func (s *sitemapService) UpdateColumn(id int64, name string, value interface{}) error {
	return dao.SitemapDao.UpdateColumn(id, name, value)
}

func (s *sitemapService) Delete(id int64) {
	dao.SitemapDao.Delete(id)
}

func (s *sitemapService) GenerateToday() {
	if s.building {
		log.Info("sitemap is in building")
		return
	}

	s.building = true
	defer func() {
		s.building = false
	}()

	dateFrom := util.WithTimeAsStartOfDay(time.Now())
	dateTo := dateFrom.Add(time.Hour * 24)

	s.GenerateMisc()
	s.GenerateUser()
	s.Generate(util.Timestamp(dateFrom), util.Timestamp(dateTo))
}

func (s *sitemapService) Generate(dateFrom, dateTo int64) {
	sitemapName := "sitemap-" + util.TimeFormat(util.TimeFromTimestamp(dateFrom), util.FMT_DATE)
	sm := sitemap.NewGenerator(viper.GetString("uploader.oss.host"), sitemapPath, sitemapName, func(sm *sitemap.Generator, sitemapLoc string) {
		s.AddSitemapIndex(sm, sitemapLoc)
	})

	// topics
	TopicService.ScanDesc(dateFrom, dateTo, func(topics []model.Topic) {
		for _, topic := range topics {
			if topic.Status == model.StatusOk {
				sm.AddURL(sitemap.URL{
					Loc:        urls.TopicUrl(topic.ID),
					Lastmod:    util.TimeFromTimestamp(topic.LastCommentTime),
					Changefreq: sitemap.ChangefreqDaily,
					Priority:   "0.8",
				})
			}
		}
	})

	// articles
	ArticleService.ScanDesc(dateFrom, dateTo, func(articles []model.Article) {
		for _, article := range articles {
			if article.Status == model.StatusOk {
				sm.AddURL(sitemap.URL{
					Loc:        urls.ArticleUrl(article.ID),
					Lastmod:    util.TimeFromTimestamp(article.UpdateTime),
					Changefreq: sitemap.ChangefreqWeekly,
					Priority:   "0.6",
				})
			}
		}
	})

	sm.Finalize()
}

func (s *sitemapService) GenerateMisc() {
	sm := sitemap.NewGenerator(viper.GetString("uploader.oss.host"), sitemapPath, "sitemap-misc", func(sm *sitemap.Generator, sitemapLoc string) {
		s.AddSitemapIndex(sm, sitemapLoc)
	})
	sm.AddURL(sitemap.URL{
		Loc:        urls.AbsUrl("/"),
		Lastmod:    time.Now(),
		Changefreq: sitemap.ChangefreqDaily,
		Priority:   "1.0",
	})
	sm.AddURL(sitemap.URL{
		Loc:        urls.AbsUrl("/topics"),
		Lastmod:    time.Now(),
		Changefreq: sitemap.ChangefreqDaily,
		Priority:   "1.0",
	})
	sm.AddURL(sitemap.URL{
		Loc:        urls.AbsUrl("/articles"),
		Lastmod:    time.Now(),
		Changefreq: sitemap.ChangefreqDaily,
		Priority:   "1.0",
	})
	sm.AddURL(sitemap.URL{
		Loc:        urls.AbsUrl("/projects"),
		Lastmod:    time.Now(),
		Changefreq: sitemap.ChangefreqDaily,
		Priority:   "1.0",
	})

	TagService.Scan(func(tags []model.Tag) bool { 
		for _, tag := range tags {
			tagUrl := urls.TagArticlesUrl(tag.ID)

			sm.AddURL(sitemap.URL{
				Loc:        tagUrl,
				Lastmod:    time.Now(),
				Changefreq: sitemap.ChangefreqDaily,
				Priority:   "0.6",
			})
		}
		return true
	})

	sm.Finalize()
}

func (s *sitemapService) GenerateUser() {
	sm := sitemap.NewGenerator(viper.GetString("uploader.oss.host"), sitemapPath, "sitemap-user", func(sm *sitemap.Generator, sitemapLoc string) {
		s.AddSitemapIndex(sm, sitemapLoc)
	})
	UserService.Scan(func(users []model.User) {
		for _, user := range users {
			sm.AddURL(sitemap.URL{
				Loc:        urls.UserUrl(user.ID),
				Lastmod:    time.Now(),
				Changefreq: sitemap.ChangefreqWeekly,
				Priority:   "0.6",
			})
		}
	})

	sm.Finalize()
}

func (s *sitemapService) AddSitemapIndex(sm *sitemap.Generator, sitemapLoc string) {
	locName := util.MD5(sitemapLoc)
	t := s.FindOne(sqlcnd.NewSqlCnd().Eq("loc_name", locName))
	if t == nil {
		_ = s.Create(&model.Sitemap{
			Model:      model.Model{},
			Loc:        sitemapLoc,
			Lastmod:    util.NowTimestamp(),
			LocName:    locName,
			CreateTime: util.NowTimestamp(),
		})
	} else {
		t.Lastmod = util.NowTimestamp()
		_ = s.Update(t)
	}

	go func() {
		s.GenerateSitemapIndex(sm)
	}()
}

func (s *sitemapService) GenerateSitemapIndex(sm *sitemap.Generator) {
	sitemaps := s.Find(sqlcnd.NewSqlCnd().Desc("id"))

	if len(sitemaps) == 0 {
		return
	}

	var sitemapLocs []sitemap.IndexURL
	for _, s := range sitemaps {
		sitemapLocs = append(sitemapLocs, sitemap.IndexURL{
			Loc:     s.Loc,
			Lastmod: util.TimeFromTimestamp(s.Lastmod),
		})
	}
	sm.WriteIndex(sitemapLocs)
}
