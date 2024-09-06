package foree_service

import "time"

type CacheItem[T any] struct {
	item      T
	createdAt time.Time
}
