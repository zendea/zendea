package cache

import (
	"time"

	"github.com/goburrow/cache"

	"zendea/dao"
	"zendea/model"
	"zendea/util/sqlcnd"
)

var (
	allNodesCacheKey = "all_nodes_cache"
)

type nodeCache struct {
	cache cache.LoadingCache
	allCache cache.LoadingCache
}

var NodeCache = newNodeCache()

func newNodeCache() *nodeCache {
	return &nodeCache{
		cache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				value = dao.NodeDao.Get(key2Int64(key))
				return
			},
			cache.WithMaximumSize(1000),
			cache.WithExpireAfterAccess(30*time.Minute),
		),
		allCache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				value = dao.NodeDao.Find(sqlcnd.NewSqlCnd().Eq("status", model.StatusOk).Asc("sort_no").Desc("id"))
				return
			},
			cache.WithMaximumSize(10),
			cache.WithRefreshAfterWrite(30*time.Minute),
		),
	}
}

func (c *nodeCache) Get(nodeId int64) *model.Node {
	if nodeId <= 0 {
		return nil
	}
	val, err := c.cache.Get(nodeId)
	if err != nil {
		return nil
	}
	return val.(*model.Node)
}

func (c *nodeCache) Invalidate(nodeId int64) {
	c.cache.Invalidate(nodeId)
}

func (c *nodeCache) GetAll() []model.Node {
	val, err := c.allCache.Get(allNodesCacheKey)
	if err != nil {
		return nil
	}
	if val != nil {
		return val.([]model.Node)
	}
	return nil
}

func (c *nodeCache) InvalidateAll() {
	c.allCache.Invalidate(allNodesCacheKey)
}
