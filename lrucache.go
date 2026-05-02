// Provides a generic thread safe implementation of the
// Least Recently Used algorithm for a cache with cache eviction policies
// which means if insert operations exceds the cache size the less frecuent
// item it's removed.
package lrucache

import (
	"container/list"
	"iter"
	"sync"
)

type entry[K comparable, V any] struct {
	key   K
	value V
}

type Cache[K comparable, V any] struct {
	capacity int

	rw       *sync.RWMutex
	mapKey   map[K]*list.Element
	elements *list.List
}

// New creates a new LRUCache with the specified capacity.
func New[K comparable, V any](capacity int) *Cache[K, V] {
	return &Cache[K, V]{
		capacity: capacity,
		mapKey:   map[K]*list.Element{},
		elements: list.New(),
		rw:       &sync.RWMutex{},
	}
}

// Get return the value of the entry base on the key after that the
// element it's moved on front of the cache.If no key are found returns nil.
func (c *Cache[K, V]) Get(key K) *V {
	c.rw.RLock()
	e, ok := c.mapKey[key]
	c.rw.RUnlock()
	if ok {
		c.rw.Lock()
		defer c.rw.Unlock()
		kv := e.Value.(*entry[K, V])
		c.elements.Remove(e)
		c.elements.PushFront(e)
		return &kv.value
	}
	return nil
}

// Contains checks the key existence on the cache if exists the
// element it's moved on front of the cache.
func (c *Cache[K, V]) Contains(key K) bool {
	c.rw.RLocker()
	e, ok := c.mapKey[key]
	c.rw.RUnlock()
	if ok {
		c.rw.Lock()
		defer c.rw.Unlock()
		c.elements.Remove(e)
		c.elements.PushFront(e)
		return true
	}
	return false
}

// Put insert a new entry in the cache if the key it's already found an update it's made
// element it's moved on front of the cache.
func (c *Cache[K, V]) Put(key K, value V) {
	c.rw.RLock()
	e, ok := c.mapKey[key]
	c.rw.RUnlock()
	if ok { // If the key it's already found for avoid duplicates keys
		c.rw.Lock()
		defer c.rw.Unlock()
		kv := e.Value.(*entry[K, V])
		kv.value = value
		c.elements.Remove(e)
		c.elements.PushFront(e)
		return
	}
	kv := &entry[K, V]{key: key, value: value}
	c.rw.Lock()
	e = c.elements.PushFront(kv)
	c.mapKey[key] = e
	c.rw.Unlock()
	c.rw.RLock()
	legth := c.elements.Len()
	c.rw.RUnlock()
	if legth > c.capacity {
		c.rw.Lock()
		defer c.rw.Lock()
		e = c.elements.Back()
		c.elements.Remove(e)
		kv = e.Value.(*entry[K, V])
		delete(c.mapKey, kv.key)
	}
}

// Update an entry if exists the element
// it's moved on front of the cache.
func (c *Cache[K, V]) Update(key K, value V) {
	c.rw.RLock()
	e, ok := c.mapKey[key]
	c.rw.RUnlock()
	if ok {
		c.rw.Lock()
		defer c.rw.Unlock()
		kv := e.Value.(*entry[K, V])
		kv.value = value
		c.elements.Remove(e)
		c.elements.MoveToFront(e)
	}
}

// Delete removes an entry base on the key.
func (c *Cache[K, V]) Delete(key K) {
	c.rw.RLock()
	e, ok := c.mapKey[key]
	c.rw.RUnlock()
	if ok {
		c.rw.Lock()
		defer c.rw.Unlock()
		kv := e.Value.(*entry[K, V])
		c.elements.Remove(e)
		delete(c.mapKey, kv.key)
	}
}

// Iterator over the entries in the cache
func (c *Cache[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for e := c.elements.Front(); e != nil; e = e.Next() {
			kv := e.Value.(*entry[K, V])
			if !yield(kv.key, kv.value) {
				return
			}
		}
	}
}
