package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mu       sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.Get(key)
	if ok {
		first := c.queue.Front()
		first.Value = value
		c.items[key] = first
		return true
	} else {
		currentQueueLen := c.queue.Len()
		if currentQueueLen == c.capacity {
			c.queue.Remove(c.queue.Back())
		}
		c.items[key] = c.queue.PushFront(value)
		return false
	}

}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	if v, ok := c.items[key]; !ok {
		return nil, ok
	} else {
		c.items[key] = c.queue.PushFront(v.Value)
		return v.Value, true
	}

}

func (c *lruCache) Clear() {
	c = NewCache(c.capacity).(*lruCache)
}

type cacheItem struct {
	key   string
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
