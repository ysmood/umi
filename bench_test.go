package umi_test

import (
	"testing"

	"strconv"

	golang_lru "github.com/hashicorp/golang-lru"
	"github.com/ysmood/umi"
)

func BenchmarkSet(b *testing.B) {
	c := umi.New(nil)
	val := 10

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		c.Set(strconv.FormatInt(int64(n), 10), &val)
	}
}

func BenchmarkGet(b *testing.B) {
	c := umi.New(nil)

	key := string(make([]byte, 100))
	val := 10

	c.Set(key, &val)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		c.Get(key)
	}
}
func BenchmarkPeek(b *testing.B) {
	c := umi.New(nil)

	key := string(make([]byte, 100))
	val := 10

	c.Set(key, &val)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		c.Peek(key)
	}
}

func BenchmarkGetRate0(b *testing.B) {
	c := umi.New(&umi.Options{
		PromoteRate: -1,
	})

	key := string(make([]byte, 100))
	val := 10

	c.Set(key, &val)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		c.Get(key)
	}
}

func BenchmarkGetRate10000(b *testing.B) {
	c := umi.New(&umi.Options{
		PromoteRate: 10000,
	})

	key := string(make([]byte, 100))
	val := 10

	c.Set(key, &val)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		c.Get(key)
	}
}

func BenchmarkSet_golang_lru(b *testing.B) {
	c, _ := golang_lru.New(1000)

	val := 10

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		c.Add(strconv.FormatInt(int64(n), 10), &val)
	}
}

func BenchmarkGet_golang_lru(b *testing.B) {
	c, _ := golang_lru.New(1000)

	key := string(make([]byte, 100))
	val := 10

	c.Add(key, &val)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		c.Get(key)
	}
}

func BenchmarkParallel(b *testing.B) {
	c := umi.New(nil)

	i := 0
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if i%10000 == 0 {
				c.Set("a", i)
			} else {
				c.Get("a")
			}
			i++
		}
	})
}

func BenchmarkParallel_golang_lru(b *testing.B) {
	c, _ := golang_lru.New(1000)

	i := 0
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if i%10000 == 0 {
				c.Add("a", i)
			} else {
				c.Get("a")
			}
			i++
		}
	})
}
