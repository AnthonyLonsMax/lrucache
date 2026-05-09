# lrucache

Generic thread-safe LRU (Least Recently Used) cache for Go.

## Usage

```go
import "codeberg.org/AnthonyLonsMax/lrucache"

c := lrucache.New[string, int](100) // capacity

c.Put("a", 1)
c.Put("b", 2)

v, ok := c.Get("a")    // 1, true
c.Contains("a")        // true
c.Update("a", 99)
c.Delete("b")

// iterate
for k, v := range c.All() {
    fmt.Println(k, v)
}
```

When the cache exceeds capacity, the least recently used entry is evicted on insert.

> **Tip:** Don't call mutating methods (`Get`, `Put`, `Delete`, etc.) inside an `All()` loop — it will deadlock. Collect keys first if you need to modify while iterating.

## License

MIT
