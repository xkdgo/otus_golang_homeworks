package hw04lrucache

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
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	if _, keyInCache := c.Get(key); keyInCache {
		firstListI := c.queue.Front()
		cachI := firstListI.Value.(*cacheItem)
		cachI.value = value
		firstListI.Value = cachI
		c.items[key] = firstListI
		return true
	}
	if c.capacity == c.queue.Len() {
		deleteCandidate := c.queue.Back().Value.(*cacheItem)
		delete(c.items, Key(deleteCandidate.key))
		c.queue.Remove(c.queue.Back())
	}
	newItem := &cacheItem{key: string(key), value: value}
	c.items[key] = c.queue.PushFront(newItem)
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	listI, exist := c.items[key]
	if !exist {
		return nil, false
	}
	c.queue.MoveToFront(listI)
	c.items[key] = c.queue.Front()
	cachI := listI.Value.(*cacheItem)
	return cachI.value, true
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
