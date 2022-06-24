package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type CacheApp struct {
	client *cache.Cache
}

func NewCacheApp(rd *cache.Cache) *CacheApp {
	return &CacheApp{client: rd}
}

func (c *CacheApp) Get(key string) interface{} {
	x, exist := c.client.Get(key)
	if !exist {
		return nil
	}
	return x
}

func (c *CacheApp) Put(key string, val interface{}, expr uint64) error {
	c.client.Set(key, val, time.Duration(expr)*time.Second)
	return nil
}

func (c *CacheApp) Delete(key string) error {
	c.client.Delete(key)
	return nil
}

func (c *CacheApp) Flush() error {
	c.client.Flush()
	return nil
}

func (c *CacheApp) IsExist(key string) bool {
	if _, exist := c.client.Get(key); !exist {
		return false
	}
	return true
}
