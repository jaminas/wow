package cache

import (
	"sync"
	"time"
)

type InMemoryCache struct {
	dataMap map[int]inMemoryValue
	lock    *sync.Mutex
}

type inMemoryValue struct {
	SetTime    int64
	Expiration int64
}

func InitInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		dataMap: make(map[int]inMemoryValue, 0),
		lock:    &sync.Mutex{},
	}
}

func (c *InMemoryCache) Add(key int, expiration int64) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.dataMap[key] = inMemoryValue{
		SetTime:    time.Now().Unix(),
		Expiration: expiration,
	}
	return nil
}

func (c *InMemoryCache) Get(key int) (bool, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	value, ok := c.dataMap[key]
	if ok && time.Now().Unix()-value.SetTime > value.Expiration {
		return false, nil
	}
	return ok, nil
}

// Delete - delete key from cache
func (c *InMemoryCache) Delete(key int) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.dataMap, key)
}
