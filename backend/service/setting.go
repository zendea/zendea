package service

import (
	"errors"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/tidwall/gjson"

	"zendea/cache"
	"zendea/dao"
	"zendea/model"
	"zendea/util"
	"zendea/util/log"
	"zendea/util/sqlcnd"
)

var SettingService = newSettingService()

func newSettingService() *settingService {
	return &settingService{}
}

type settingService struct {
}

func (s *settingService) Get(id int64) *model.Setting {
	return dao.SettingDao.Get(id)
}

func (s *settingService) Take(where ...interface{}) *model.Setting {
	return dao.SettingDao.Take(where...)
}

func (s *settingService) Find(cnd *sqlcnd.SqlCnd) []model.Setting {
	return dao.SettingDao.Find(cnd)
}

func (s *settingService) FindOne(cnd *sqlcnd.SqlCnd) *model.Setting {
	return dao.SettingDao.FindOne(cnd)
}

func (s *settingService) List(cnd *sqlcnd.SqlCnd) (list []model.Setting, paging *sqlcnd.Paging) {
	return dao.SettingDao.List(cnd)
}

func (s *settingService) GetAll() []model.Setting {
	return dao.SettingDao.Find(sqlcnd.NewSqlCnd().Asc("id"))
}

func (s *settingService) SetAll(configStr string) error {
	json := gjson.Parse(configStr)
	configs, ok := json.Value().(map[string]interface{})
	if !ok {
		return errors.New("配置数据格式错误")
	}
	return dao.Tx(func(tx *gorm.DB) error {
		for k := range configs {
			v := json.Get(k).String()
			if err := s.setSingle(tx, k, v, "", ""); err != nil {
				return err
			}
		}
		return nil
	})
}

// 设置配置，如果配置不存在，那么创建
func (s *settingService) Set(key, value, name, description string) error {
	return dao.Tx(func(tx *gorm.DB) error {
		if err := s.setSingle(tx, key, value, name, description); err != nil {
			return err
		}
		return nil
	})
}

func (s *settingService) setSingle(db *gorm.DB, key, value, name, description string) error {
	if len(key) == 0 {
		return errors.New("sys config key is null")
	}
	sysConfig := dao.SettingDao.GetByKey(key)
	if sysConfig == nil {
		sysConfig = &model.Setting{
			CreateTime: util.NowTimestamp(),
		}
	}
	sysConfig.Key = key
	sysConfig.Value = value
	sysConfig.UpdateTime = util.NowTimestamp()

	if len(name) > 0 {
		sysConfig.Name = name
	}
	if len(description) > 0 {
		sysConfig.Description = description
	}

	var err error
	if sysConfig.ID > 0 {
		err = dao.SettingDao.Update(sysConfig)
	} else {
		err = dao.SettingDao.Create(sysConfig)
	}
	if err != nil {
		return err
	} else {
		cache.SettingCache.Invalidate(key)
		return nil
	}
}

func (s *settingService) GetSetting() *model.ConfigData {
	var (
		siteTitle        = cache.SettingCache.GetValue(model.SettingSiteTitle)
		siteDescription  = cache.SettingCache.GetValue(model.SettingSiteDescription)
		siteKeywords     = cache.SettingCache.GetValue(model.SettingSiteKeywords)
		siteNavs         = cache.SettingCache.GetValue(model.SettingSiteNavs)
		siteTips         = cache.SettingCache.GetValue(model.SettingSiteTips)
		siteNotification = cache.SettingCache.GetValue(model.SettingSiteNotification)
		recommendTags    = cache.SettingCache.GetValue(model.SettingRecommendTags)
		scoreConfigStr   = cache.SettingCache.GetValue(model.SettingScoreConfig)
		defaultNodeIdStr = cache.SettingCache.GetValue(model.SettingDefaultNodeId)
		siteIndexHtml    = cache.SettingCache.GetValue(model.SettingSiteIndexHtml)

		statUserCountStr    = cache.SettingCache.GetValue(model.SettingStatUserCount)
		statTopicCountStr   = cache.SettingCache.GetValue(model.SettingStatTopicCount)
		statCommentCountStr = cache.SettingCache.GetValue(model.SettingStatCommentCount)
	)

	var siteKeywordsArr []string
	if err := util.ParseJson(siteKeywords, &siteKeywordsArr); err != nil {
		log.Warn("站点关键词数据错误")
	}

	var siteNavsArr []model.SiteNav
	if err := util.ParseJson(siteNavs, &siteNavsArr); err != nil {
		log.Warn("站点导航数据错误")
	}

	var siteTipsArr []model.SiteTip
	if err := util.ParseJson(siteTips, &siteTipsArr); err != nil {
		log.Warn("小贴士数据错误")
	}

	var recommendTagsArr []string
	if err := util.ParseJson(recommendTags, &recommendTagsArr); err != nil {
		log.Warn("推荐标签数据错误")
	}

	var scoreConfig model.ScoreConfig
	if err := util.ParseJson(scoreConfigStr, &scoreConfig); err != nil {
		log.Warn("积分配置错误")
	}

	var defaultNodeId, _ = strconv.ParseInt(defaultNodeIdStr, 10, 64)
	var statUserCount, _ = strconv.ParseInt(statUserCountStr, 10, 64)
	var statTopicCount, _ = strconv.ParseInt(statTopicCountStr, 10, 64)
	var statCommentCount, _ = strconv.ParseInt(statCommentCountStr, 10, 64)

	return &model.ConfigData{
		SiteTitle:        siteTitle,
		SiteDescription:  siteDescription,
		SiteKeywords:     siteKeywordsArr,
		SiteNavs:         siteNavsArr,
		SiteTips:         siteTipsArr,
		SiteNotification: siteNotification,
		SiteIndexHtml:    siteIndexHtml,
		RecommendTags:    recommendTagsArr,
		ScoreConfig:      scoreConfig,
		DefaultNodeId:    defaultNodeId,
		StatUserCount:    statUserCount,
		StatTopicCount:   statTopicCount,
		StatCommentCount: statCommentCount,
	}
}
