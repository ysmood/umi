# Umi

[![GoDoc](https://godoc.org/github.com/ysmood/umi?status.svg)](http://godoc.org/github.com/ysmood/umi)
[![Build Status](https://travis-ci.org/ysmood/umi.svg)](https://travis-ci.org/ysmood/umi)

Umi is a high performance LRU cache lib.

Different from other libs, Umi will automatically calculate the memory size
of random data structure for you, and limit the number of items by the total bytes of them,
not just the count of them.

## Road Map

- [x] Thread safe

- [x] All the basic operations' complexity should be O(1).

- [x] The max size of the cache is byte based, not count based.

- [x] High performance TTL with a GC.

- [x] The algorithm used to replace cache is a variation of LRU.

## Quick Start

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

## FAQ

- Is the auto-calculated byte size of item safe?

  I will try my best to make it close to the real size, but because of the GC nature of
  golang, the safety is not guaranteed. So in case you have a more precise way to calculate
  the size, you can implement the `Sizable` interface of each item.

## Benchmark

`go test -bench . -benchmem`

The `get` performance is faster than the [golang-lru](https://github.com/hashicorp/golang-lru).
The `set` is slower because Umi's data struct contains extra info
to calculate such as TTL and byte size. This trade-off for more functionalities is acceptable.

Umi's faster performance is because it uses two read-write locks for the
internal `map` and `linked-list`, the more atomic lock time make the total lock time smaller.

How umi optimizes locks:

```txt
       umi: | ops1 | map-lock | ops2 | list-lock | ops3 |

golang_lru: | -------------- write-lock --------------- |

      time: ----------------------------------------------->
```

```txt
BenchmarkSet-6              	 3000000	       421 ns/op	     156 B/op	       2 allocs/op
BenchmarkSet_golang_lru-6   	 5000000	       332 ns/op	     105 B/op	       4 allocs/op
BenchmarkGet-6              	30000000	        41.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkGet_golang_lru-6   	20000000	        64.8 ns/op	      16 B/op	       1 allocs/op
BenchmarkPeek-6             	100000000	        23.4 ns/op	       0 B/op	       0 allocs/op
```
