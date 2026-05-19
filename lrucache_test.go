package lrucache_test

import (
	"testing"

	"github.com/AnthonyLonsMax/lrucache"
)

func TestAddElements(t *testing.T) {
	t.Parallel()

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

func TestEviction(t *testing.T) {
	t.Parallel()

	lru := lrucache.New[int, int](3)
	lru.Put(1, 10)
	lru.Put(2, 20)
	lru.Put(3, 30)
	lru.Put(4, 40) // should evict key 1

	if _, ok := lru.Get(1); ok {
		t.Fatalf("expected key 1 to be evicted")
	}
}

func TestGetMovesToFront(t *testing.T) {
	t.Parallel()

	lru := lrucache.New[int, int](3)
	lru.Put(1, 10)
	lru.Put(2, 20)
	lru.Put(3, 30)
	lru.Get(1)     // move 1 to front
	lru.Put(4, 40) // should evict 2 (LRU), not 1

	if _, ok := lru.Get(2); ok {
		t.Fatalf("expected key 2 to be evicted")
	}

	if v, ok := lru.Get(1); !ok || v != 10 {
		t.Fatalf("expected key 1 to exist with value 10")
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	lru := lrucache.New[int, int](3)
	lru.Put(1, 10)
	lru.Update(1, 99)
	v, ok := lru.Get(1)

	if !ok || v != 99 {
		t.Fatalf("expected value 99, got %d", v)
	}
}

func TestGetNonExistent(t *testing.T) {
	t.Parallel()

	lru := lrucache.New[int, int](3)
	_, ok := lru.Get(42)

	if ok {
		t.Fatalf("expected false for non-existent key")
	}
}
