package cache

import (
	"time"

	"github.com/goburrow/cache"

	"zendea/dao"
	"zendea/model"
)

type sectionCache struct {
	cache cache.LoadingCache
}

var SectionCache = newSectionCache()

func newSectionCache() *sectionCache {
	return &sectionCache{
		cache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				value = dao.SectionDao.Get(key2Int64(key))
				return
			},
			cache.WithMaximumSize(1000),
			cache.WithExpireAfterAccess(30*time.Minute),
		),
	}
}

func (c *sectionCache) Get(sectionId int64) *model.Section {
	if sectionId <= 0 {
		return nil
	}
	val, err := c.cache.Get(sectionId)
	if err != nil {
		return nil
	}
	return val.(*model.Section)
}

func (c *sectionCache) Invalidate(sectionId int64) {
	c.cache.Invalidate(sectionId)
}
