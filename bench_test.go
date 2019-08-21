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

func BenchmarkSet_golang_lru(b *testing.B) {
	c, _ := golang_lru.New(1000)

	val := 10

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		c.Add(strconv.FormatInt(int64(n), 10), &val)
	}
}

func BenchmarkGet(b *testing.B) {
	c := umi.New(nil)

	key := "key"
	val := 10

	c.Set(key, &val)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		c.Get(key)
	}
}

func BenchmarkGet_golang_lru(b *testing.B) {
	c, _ := golang_lru.New(1000)

	key := "key"
	val := 10

	c.Add(key, &val)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		c.Get(key)
	}
}

func BenchmarkPeek(b *testing.B) {
	c := umi.New(nil)

	key := "key"
	val := 10

	c.Set(key, &val)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		c.Peek(key)
	}
}
