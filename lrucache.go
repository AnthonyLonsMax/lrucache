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

func New[K comparable, V any](capacity int) *Cache[K, V] {
	return &Cache[K, V]{
		capacity: capacity,
		mapKey:   map[K]*list.Element{},
		elements: list.New(),
		rw:       &sync.RWMutex{},
	}
}

func (c *Cache[K, V]) Get(key K) *V {
	c.rw.Lock()
	defer c.rw.Unlock()
	if e, ok := c.mapKey[key]; ok {
		kv := e.Value.(*entry[K, V])
		c.elements.Remove(e)
		c.elements.PushFront(e)
		return &kv.value
	}
	return nil
}

func (c *Cache[K, V]) Contains(key K) bool {
	c.rw.Lock()
	defer c.rw.Unlock()
	if e, ok := c.mapKey[key]; ok {
		c.elements.Remove(e)
		c.elements.PushFront(e)
		return true
	}
	return false
}

func (c *Cache[K, V]) Put(key K, value V) {
	c.rw.Lock()
	defer c.rw.Unlock()
	kv := &entry[K, V]{key: key, value: value}
	if e, ok := c.mapKey[key]; ok { // If the element it's already found for avoid duplicates
		kv := e.Value.(*entry[K, V])
		kv.value = value
		c.elements.Remove(e)
		c.elements.PushFront(e)
		return
	}
	e := c.elements.PushFront(kv)
	c.mapKey[key] = e
	if c.elements.Len() > c.capacity {
		e = c.elements.Back()
		c.elements.Remove(e)
		kv = e.Value.(*entry[K, V])
		delete(c.mapKey, kv.key)
	}
}

func (c *Cache[K, V]) Update(key K, value V) {
	c.rw.Lock()
	defer c.rw.Unlock()
	if e, ok := c.mapKey[key]; ok {
		kv := e.Value.(*entry[K, V])
		kv.value = value
		c.elements.Remove(e)
		c.elements.MoveToFront(e)
	}
}

func (c *Cache[K, V]) Delete(key K) {
	c.rw.Lock()
	defer c.rw.Unlock()
	if e, ok := c.mapKey[key]; ok {
		kv := e.Value.(*entry[K, V])
		c.elements.Remove(e)
		delete(c.mapKey, kv.key)
	}
}

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
