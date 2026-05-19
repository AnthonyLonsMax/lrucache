// Package lrucache provides a generic thread safe implementation of the
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

// Cache provides the lruchache implementation.
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
// element it's moved on front of the cache.
func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.rw.Lock()
	defer c.rw.Unlock()

	e, ok := c.mapKey[key]

	var zero V

	if !ok {
		return zero, false
	}

	c.elements.MoveToFront(e)
	ele, ok := e.Value.(*entry[K, V])

	if !ok {
		return zero, false
	}

	return ele.value, true
}

// Contains checks the key existence on the cache if exists the
// element it's moved on front of the cache.
func (c *Cache[K, V]) Contains(key K) bool {
	c.rw.Lock()
	defer c.rw.Unlock()

	e, ok := c.mapKey[key]

	if !ok {
		return false
	}

	c.elements.MoveToFront(e)

	return true
}

// Put insert a new entry in the cache if the key it's already found an update it's made
// element it's moved on front of the cache.
func (c *Cache[K, V]) Put(key K, value V) {
	c.rw.Lock()
	defer c.rw.Unlock()

	if e, ok := c.mapKey[key]; ok {
		ele, ok := e.Value.(*entry[K, V])
		if !ok {
			panic("invalid type cohercion")
		}

		ele.value = value

		c.elements.MoveToFront(e)

		return
	}

	e := c.elements.PushFront(&entry[K, V]{key: key, value: value})
	c.mapKey[key] = e

	if c.elements.Len() > c.capacity {
		e = c.elements.Back()
		c.elements.Remove(e)
		ele, ok := e.Value.(*entry[K, V])

		if !ok {
			panic("invalid type cohercion")
		}

		delete(c.mapKey, ele.key)
	}
}

// Update an entry if exists the element
// it's moved on front of the cache.
func (c *Cache[K, V]) Update(key K, value V) {
	c.rw.Lock()
	defer c.rw.Unlock()

	if e, ok := c.mapKey[key]; ok {
		ele, ok := e.Value.(*entry[K, V])

		if !ok {
			panic("invalid type cohercion")
		}

		ele.value = value

		c.elements.MoveToFront(e)
	}
}

// Delete removes an entry base on the key.
func (c *Cache[K, V]) Delete(key K) {
	c.rw.Lock()
	defer c.rw.Unlock()

	if e, ok := c.mapKey[key]; ok {
		c.elements.Remove(e)
		delete(c.mapKey, key)
	}
}

// All returns an iterator over the entries in the cache.
func (c *Cache[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		c.rw.RLock()
		defer c.rw.RUnlock()

		for e := c.elements.Front(); e != nil; e = e.Next() {
			kv, ok := e.Value.(*entry[K, V])
			if !ok {
				panic("invalid type cohercion")
			}

			if !yield(kv.key, kv.value) {
				return
			}
		}
	}
}
