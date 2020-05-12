package cache

import (
	"time"

	"github.com/goburrow/cache"

	"zendea/dao"
	"zendea/model"
)

type settingCache struct {
	cache cache.LoadingCache
}

var SettingCache = newSettingCache()

func newSettingCache() *settingCache {
	return &settingCache{
		cache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				value = dao.SettingDao.GetByKey(key.(string))
				return
			},
			cache.WithMaximumSize(1000),
			cache.WithExpireAfterAccess(30*time.Minute),
		),
	}
}

func (c *settingCache) Get(key string) *model.Setting {
	val, err := c.cache.Get(key)
	if err != nil {
		return nil
	}
	if val != nil {
		return val.(*model.Setting)
	}
	return nil
}

func (c *settingCache) GetValue(key string) string {
	sysConfig := c.Get(key)
	if sysConfig == nil {
		return ""
	}
	return sysConfig.Value
}

func (c *settingCache) Invalidate(key string) {
	c.cache.Invalidate(key)
}
