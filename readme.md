- [x] Thread safe

- [x] All the basic operations' complexity should be O(1).

- [x] The max size of the cache is byte based, not count based.

- [x] High performance TTL with a GC.

- [x] The algorithm used to replace cache is a variation of LRU.

- [x] With a rate to throw away promotions.

- [ ] The caches are grouped into two types: memory cache and file cache. Works like swap,
  but in a much more efficient way

# Quick Start

For more examples, see `umi_test.go`.

```go
c := umi.New(nil)

anyRandomData := struct {
    A string
    B int
}{
    A: "test",
    B: 10,
}

c.Set("a", anyRandomData)

v, _ := c.Get("a")

memorySize := c.Size()

fmt.Println(v, memorySize)
```


# Benchmark

`go test -bench . -benchmem`

The `get` performance is 4x faster than the `https://github.com/hashicorp/golang-lru`.
The `set` is a little slower. It's because `umi`'s data struct contains extra info
such as TTL and byte size. This trade-off for more functionalities is acceptable.

```
BenchmarkSet-8              	 2000000	       900 ns/op	     160 B/op	       4 allocs/op
BenchmarkGet-8              	50000000	        33.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkPeek-8             	50000000	        26.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkGetRate0-8         	30000000	        46.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkGetRate10000-8     	50000000	        33.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkSet_golang_lru-8   	 3000000	       653 ns/op	     105 B/op	       4 allocs/op
BenchmarkGet_golang_lru-8   	10000000	       133 ns/op	       0 B/op	       0 allocs/op
```

