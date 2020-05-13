package cache

import (
	"time"

	"github.com/goburrow/cache"

	"zendea/dao"
	"zendea/model"
	"zendea/util/sqlcnd"
)

var (
	allSectionsCacheKey = "all_sections_cache"
)

type sectionCache struct {
	cache cache.LoadingCache
	allCache cache.LoadingCache
	sectionNodesCache cache.LoadingCache
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
		allCache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				value = dao.SectionDao.Find(sqlcnd.NewSqlCnd().Asc("sort_no").Desc("id"))
				return
			},
			cache.WithMaximumSize(1000),
			cache.WithRefreshAfterWrite(30*time.Minute),
		),
		sectionNodesCache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				value = dao.NodeDao.Find(sqlcnd.NewSqlCnd().Where("section_id = ?", key2Int64(key)).Asc("sort_no").Desc("id"))
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

func (c *sectionCache) GetAll() []model.Section {
	val, err := c.allCache.Get(allSectionsCacheKey)
	if err != nil {
		return nil
	}
	if val != nil {
		return val.([]model.Section)
	}
	return nil
}

func (c *sectionCache) InvalidateAll() {
	c.allCache.Invalidate(allSectionsCacheKey)
}

func (c *sectionCache) GetSectionNodes(sectionId int64) []model.Node {
	val, err := c.sectionNodesCache.Get(sectionId)
	if err != nil {
		return nil
	}
	if val != nil {
		return val.([]model.Node)
	}
	return nil
}

func (c *sectionCache) InvalidateSectionNodes(sectionId int64) {
	c.sectionNodesCache.Invalidate(sectionId)
}

