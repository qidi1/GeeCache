package lru

import (
	"container/list"
	"errors"
	"log"
)

// Cache is a LRU cache. It is not safe for concurrent access.
type Cache struct {
	maxBytes   int64
	nBytes     int64
	totalEntry int
	newList    *list.List
	oldList    *list.List
	cache      map[string]*list.Element
	// optional and executed when an entry is purged.
	OnEvicted func(key string, value Value)
}
type Value interface {
	Len() int
}
type Entry struct {
	key      string
	value    Value
	IsNew    bool
	IsVisted bool
}

func New(maxBytes int64, OnEvicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes:   maxBytes,
		nBytes:     0,
		totalEntry: 0,
		newList:    list.New(),
		oldList:    list.New(),
		cache:      make(map[string]*list.Element),
		OnEvicted:  OnEvicted,
	}
}
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, right := c.cache[key]; right {
		entry := ele.Value.(*Entry)
		ok = true
		if entry.IsNew {
			c.newList.MoveToFront(ele)
			value = entry.value
		} else {
			if entry.IsVisted {
				c.oldList.Remove(ele)
				c.newList.PushFront(ele)
				entry.IsNew = true
				value = entry.value
			} else {
				c.oldList.MoveToFront(ele)
				entry.IsVisted = true
				value = entry.value
			}
		}
	}
	return
}
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		kv := ele.Value.(*Entry)
		oldLen := int64(kv.value.Len())
		if kv.IsNew {
			c.newList.MoveToFront(ele)
			ele.Value = value
		} else {
			if !kv.IsVisted {
				c.oldList.MoveToFront(ele)
				kv.IsVisted = true
				ele.Value = value
			} else {
				c.oldList.Remove(ele)
				c.newList.PushFront(ele)
				kv.IsNew = true
			}
		}
		c.nBytes -= oldLen - int64(value.Len())
	} else {
		kv := &Entry{
			key:      key,
			value:    value,
			IsNew:    false,
			IsVisted: false,
		}
		ele := c.oldList.PushFront(kv)
		c.cache[key] = ele
		c.totalEntry += 1
		c.nBytes += int64(len(key)+2) + int64(value.Len())
	}
	c.reorganize()
	for c.maxBytes != 0 && c.nBytes > c.maxBytes {
		if err := c.RemoveOldest(); err != nil {
			log.Panic(err)
			break
		}
	}
}
func (c *Cache) RemoveOldest() error {
	var list *list.List
	if c.oldList.Len() != 0 {
		list = c.oldList
	} else if c.newList.Len() != 0 {
		list = c.newList
	} else {
		return errors.New("there is no element in the list,something wrong happend")
	}
	ele := list.Back()
	list.Remove(ele)
	kv := ele.Value.(*Entry)
	delete(c.cache, kv.key)
	c.nBytes -= int64(len(kv.key)+2) + int64(kv.value.Len())
	c.totalEntry--
	if c.OnEvicted != nil {
		c.OnEvicted(kv.key, kv.value)
	}
	return nil
}
func (c *Cache) reorganize() {
	for c.newList.Len() > int(float64(c.totalEntry)*0.625) {
		ele := c.newList.Back()
		c.newList.Remove(ele)
		c.oldList.PushFront(ele)
		kv := ele.Value.(*Entry)
		kv.IsNew = false
	}
}
func (c *Cache) Len() int {
	return c.totalEntry
}
func (c *Cache) removeForSpace() {
	//TODO 移除出足够的空间
}
func (c *Cache) setOnEvicted(OnEvicted func(key string, value Value)) {
	c.OnEvicted = OnEvicted
}
