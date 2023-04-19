package models

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var cacheInstance *cache.Cache

func InitCache() {
	cacheInstance = cache.New(10*time.Second, 1*time.Minute)
}

func InitCacheForTest() {
	cacheInstance = cache.New(1*time.Second, 1*time.Minute)
}
