package foree_service

import "time"

type CacheItem[T any] struct {
	item      T
	expiredAt time.Time
}
