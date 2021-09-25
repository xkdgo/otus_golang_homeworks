package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type cacheItem struct {
	key   string
	value interface{}
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
		firstCacheI := c.queue.Front()
		val, _ := firstCacheI.Value.(*cacheItem)
		val.value = value
		firstCacheI.Value = val
		c.items[key] = firstCacheI
		return true
	}
	currentQueueLen := c.queue.Len()
	if currentQueueLen == c.capacity {
		deleteCandidate := c.queue.Back().Value.(*cacheItem)
		delete(c.items, Key(deleteCandidate.key))
		c.queue.Remove(c.queue.Back())
	}
	newItem := new(cacheItem)
	newItem.key = string(key)
	newItem.value = value
	c.items[key] = c.queue.PushFront(newItem)
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	defer c.wg.Done()
	c.wg.Add(1)
	c.mu.Lock()
	defer c.mu.Unlock()
	listI, exist := c.items[key]
	if !exist {
		return nil, false
	}
	cacheI := listI.Value
	c.items[key] = c.queue.PushFront(cacheI)
	val, _ := cacheI.(*cacheItem)
	return val.value, true
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
