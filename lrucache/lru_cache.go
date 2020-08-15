package lrucache

import (
	"github.com/arazmj/gocache/cache"
	"github.com/arazmj/gocache/dlinklist"
	"github.com/arazmj/gocache/stats"
	"github.com/inhies/go-bytesize"
	"sync"
)

//LRUCache data structure
type LRUCache struct {
	sync.RWMutex
	cache    map[string]*dlinklist.Node
	linklist *dlinklist.DLinkedList
	capacity bytesize.ByteSize
	size     bytesize.ByteSize
}

// NewCache LRUCache constructor
func NewCache(capacity bytesize.ByteSize) cache.ICache {
	return &LRUCache{
		cache:    map[string]*dlinklist.Node{},
		linklist: dlinklist.NewLinkedList(),
		capacity: capacity,
		size:     0,
	}
}

// Get returns the value for the key
func (c *LRUCache) Get(key string) (value string, ok bool) {
	defer c.Unlock()
	c.Lock()
	if value, ok := c.cache[key]; ok {
		stats.Hits.Inc()
		node := value
		c.linklist.RemoveNode(node)
		c.linklist.AddNode(node)
		return node.Value, true
	}
	stats.Miss.Inc()
	return "", false
}

// Put updates or insert a new entry, evicts the old entry
// if cache size is larger than capacity
func (c *LRUCache) Put(key string, value string) (created bool) {
	defer c.Unlock()
	c.Lock()
	if v, ok := c.cache[key]; ok {
		node := v
		c.linklist.RemoveNode(node)
		c.linklist.AddNode(node)
		node.Value = value
		created = false
	} else {
		node := &dlinklist.Node{Key: key, Value: value}
		c.linklist.AddNode(node)
		c.cache[key] = node
		stats.Adds.Inc()
		c.size += bytesize.ByteSize(len(value))
		for c.size > c.capacity {
			stats.Deletes.Inc()
			tail := c.linklist.PopTail()
			c.size -= bytesize.ByteSize(len(tail.Value))
			delete(c.cache, tail.Key)
		}
		created = true
	}
	return created
}

// HasKey indicates the key exists or not
func (c *LRUCache) HasKey(key string) bool {
	defer c.RUnlock()
	c.RLock()
	_, ok := c.cache[key]
	return ok
}