package cache

import (
	"time"
	"zendea/dao"
	"zendea/util/sqlcnd"
	"github.com/goburrow/cache"
)

type statCache struct {
	userCountCache    cache.LoadingCache
	topicCountCache   cache.LoadingCache
	commentCountCache cache.LoadingCache
}

var StatCache = newStatCache()

func newStatCache() *statCache {
	return &statCache{
		userCountCache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				value = dao.UserDao.Count(sqlcnd.NewSqlCnd())
				return
			},
			cache.WithMaximumSize(10),
			cache.WithExpireAfterAccess(1*time.Hour),
		),
		topicCountCache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				value = dao.TopicDao.Count(sqlcnd.NewSqlCnd())
				return
			},
			cache.WithMaximumSize(10),
			cache.WithExpireAfterAccess(30*time.Minute),
		),
		commentCountCache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				value = dao.CommentDao.Count(sqlcnd.NewSqlCnd())
				return
			},
			cache.WithMaximumSize(10),
			cache.WithExpireAfterAccess(15*time.Minute),
		),
	}
}

func (c *statCache) GetUserCount() int {
	val, err := c.userCountCache.Get("data")
	if err != nil {
		return 0
	}
	return val.(int)
}

func (c *statCache) GetTopicCount() int {
	val, err := c.topicCountCache.Get("data")
	if err != nil {
		return 0
	}
	return val.(int)
}

func (c *statCache) GetCommentCount() int {
	val, err := c.commentCountCache.Get("data")
	if err != nil {
		return 0
	}
	return val.(int)
}
