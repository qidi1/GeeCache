package GeeCache

import (
	"geecache/lru"
	"sync"
)

type cache struct {
	m          sync.Mutex
	l          *lru.Cache
	cacheBytes int64
}

func (c *cache) add(key string, value ByteView) {
	c.m.Lock()
	defer c.m.Unlock()
	if c.l == nil {
		c.l = lru.New(c.cacheBytes, nil)
	}
	c.l.Add(key, value)
}
func (c *cache) get(key string) (value ByteView, ok bool) {
	c.m.Lock()
	defer c.m.Unlock()
	ok = false
	if c.l == nil {
		return
	}
	if v, right := c.l.Get(key); right {
		value = v.(ByteView)
		ok = true
	}
	return value, ok
}
