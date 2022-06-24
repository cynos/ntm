package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheRedis struct {
	client *redis.Client
}

func NewCacheRedis(rd *redis.Client) *CacheRedis {
	return &CacheRedis{client: rd}
}

func (c *CacheRedis) Get(key string) interface{} {
	x, err := c.client.Get(context.TODO(), key).Result()
	if err != nil {
		return nil
	}
	return x
}

func (c *CacheRedis) Put(key string, val interface{}, expr uint64) error {
	expiration, err := time.ParseDuration(fmt.Sprintf("%ds", expr))
	if err != nil {
		return err
	}
	return c.client.Set(context.TODO(), key, val, expiration).Err()
}

func (c *CacheRedis) Delete(key string) error {
	return c.client.Del(context.TODO(), key).Err()
}

func (c *CacheRedis) Flush() error {
	return c.client.FlushAll(context.TODO()).Err()
}

func (c *CacheRedis) IsExist(key string) bool {
	x, err := c.client.Exists(context.TODO(), key).Result()
	if err != nil {
		return false
	}
	return x > 0
}
