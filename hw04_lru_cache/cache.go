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
	wg       sync.WaitGroup
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	_, ok := c.Get(key)
	defer c.wg.Done()
	c.wg.Add(1)
	c.mu.Lock()
	defer c.mu.Unlock()
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
	defer c.wg.Done()
	c.wg.Add(1)
	c.mu.Lock()
	defer c.mu.Unlock()
	v, exist := c.items[key]
	if !exist {
		return nil, false
	}
	c.items[key] = c.queue.PushFront(v.Value)
	return v.Value, true
}

func (c *lruCache) Clear() {
	defer c.wg.Done()
	c.wg.Add(1)
	c.mu.Lock()
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
