package lrucache_test

import (
	"testing"

	"codeberg.org/AnthonyLonsMax/lrucache"
)

func TestAddElements(t *testing.T) {
	lru := lrucache.New[int, int](100)
	for i := range 100 {
		lru.Put(i, 0)
	}
	lru.Put(500, 0)
	lru.Delete(500)
	if lru.Contains(500) {
		t.Fatalf("invalid lru cache implementation")
	}
}
