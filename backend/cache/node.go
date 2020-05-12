package cache

import (
	"time"

	"github.com/goburrow/cache"

	"zendea/dao"
	"zendea/model"
)

type nodeCache struct {
	cache cache.LoadingCache
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
