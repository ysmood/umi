# Umi 海

[![Build Status](https://travis-ci.org/ysmood/umi.svg)](https://travis-ci.org/ysmood/umi)

Umi is a high performance lightweight memory-disk combined cache lib.

Different from other libs, Umi will automatically calculate the memory size
of random data structure for you, and limit the number of items by the total bytes of them,
not just the count of them.

## Road Map

- [x] Thread safe

- [x] All the basic operations' complexity should be O(1).

- [x] The max size of the cache is byte based, not count based.

- [x] High performance TTL with a GC.

- [x] The algorithm used to replace cache is a variation of LRU.

- [x] With a rate to throw away promotions.

- [ ] When memory is draining. The cache will be queued into a disk based LRU. Works like swap,
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

# FAQ

- Is the auto-calculated byte size of item safe?

  I will try my best to make it close to the real size, but because of the GC nature of
  golang, the safety is not guaranteed. So in case you have a more precise way to calculate
  the size, you can implement the `Sizable` interface of each item.


# Benchmark

`go test -bench . -benchmem`

The `get` performance is 4x faster than the `https://github.com/hashicorp/golang-lru`.
The `set` is a little slower. It's because Umi's data struct contains extra info
such as TTL and byte size. This trade-off for more functionalities is acceptable.

Umi's faster performance is because it uses two read-write locks for the
internal `map` and `linked-list`, the more atomic lock time make the total lock time smaller.
Besides, Umi doesn't promote on each `get` operation, it promotes by chance.

How umi optimizes locks:

```
       umi: | ops1 | map-lock | ops2 | list-lock | ops3 |

golang_lru: | -------------- write-lock --------------- |

      time: ----------------------------------------------->
```

```
BenchmarkSet-8                   	 1000000	      1080 ns/op	     260 B/op	       4 allocs/op
BenchmarkGet-8                   	30000000	        34.4 ns/op	       1 B/op	       0 allocs/op
BenchmarkPeek-8                  	50000000	        26.8 ns/op	       1 B/op	       0 allocs/op
BenchmarkGetRate0-8              	20000000	        66.4 ns/op	       3 B/op	       0 allocs/op
BenchmarkGetRate10000-8          	50000000	        36.0 ns/op	       1 B/op	       0 allocs/op
BenchmarkSet_golang_lru-8        	 2000000	       774 ns/op	     137 B/op	       4 allocs/op
BenchmarkGet_golang_lru-8        	10000000	       158 ns/op	       7 B/op	       0 allocs/op
BenchmarkParallel-8              	20000000	        76.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkParallel_golang_lru-8   	 5000000	       230 ns/op	       9 B/op	       0 allocs/op
```

