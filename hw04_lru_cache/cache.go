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
	}
	currentQueueLen := c.queue.Len()
	if currentQueueLen == c.capacity {
		c.queue.Remove(c.queue.Back())
	}
	c.items[key] = c.queue.PushFront(value)
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	v, exist := c.items[key]
	if !exist {
		return nil, false
	}
	c.items[key] = c.queue.PushFront(v.Value)
	return v.Value, true
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
