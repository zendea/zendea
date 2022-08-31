package cache

import (
	"github.com/goburrow/cache"

	"zendea/util/log"
)

func key2Int64(key cache.Key) int64 {
	return key.(int64)
}

func Setup() {
	log.Info("Cache setup")
}
