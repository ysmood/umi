package umi_test

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"runtime"

	"strings"

	"github.com/stretchr/testify/assert"
	"github.com/ysmood/umi"
)

func mapValues(items []*umi.Item) []interface{} {
	l := len(items)
	arr := make([]interface{}, l)
	for i := 0; i < l; i++ {
		arr[i] = items[i].Value()
	}

	return arr
}

func TestBasic(t *testing.T) {
	c := umi.New(nil)

	c.Set("a", int32(10))

	v, _ := c.Get("a")

	assert.Equal(t, 10, int(v.(int32)))
}

func TestPeek(t *testing.T) {
	c := umi.New(nil)

	c.Set("a", int32(10))

	v, _ := c.Peek("a")

	assert.Equal(t, 10, int(v.(int32)))
}

func TestDel(t *testing.T) {
	c := umi.New(nil)

	c.Set("a", int32(10))

	c.Del("a")

	_, has := c.Get("a")

	assert.Equal(t, false, has)
}

func TestPromote(t *testing.T) {
	c := umi.New(&umi.Options{
		PromoteRate: -1,
	})

	for i := int64(0); i < 10; i++ {
		c.Set(strconv.FormatInt(i, 10), int(i))
	}

	c.Get("2")
	c.Get("7")
	c.Get("4")
	c.Get("4")

	arr := c.Values()
	assert.Equal(t, []interface{}{9, 7, 8, 4, 6, 5, 2, 3, 1, 0}, arr)
}

func TestSlice(t *testing.T) {
	c := umi.New(nil)

	for i := int64(0); i < 5; i++ {
		c.Set(strconv.FormatInt(i, 10), int(i))
	}

	items := mapValues(c.Slice(2, 4))

	assert.Equal(t, []interface{}{4, 3}, items)
}

func TestSliceOutOfRange(t *testing.T) {
	c := umi.New(nil)

	items := c.Slice(0, 2)

	assert.Equal(t, make([]*umi.Item, 2), items)
}

func TestItems(t *testing.T) {
	c := umi.New(nil)

	for i := int64(0); i < 5; i++ {
		c.Set(strconv.FormatInt(i, 10), int(i))
	}

	items := c.Values()

	assert.Equal(t, []interface{}{4, 3, 2, 1, 0}, items)
}

func TestPurge(t *testing.T) {
	c := umi.New(nil)

	for i := int64(0); i < 5; i++ {
		c.Set(strconv.FormatInt(i, 10), int(i))
	}

	c.Purge()

	arr := c.Values()

	assert.Equal(t, []interface{}{}, arr)
}

func TestPromoteUntilHead(t *testing.T) {
	c := umi.New(&umi.Options{
		PromoteRate: -1,
	})

	c.Set("1", 1)
	c.Set("2", 2)
	c.Set("3", 3)

	c.Get("1")
	c.Get("1")
	c.Get("1")
	c.Get("1")
	c.Get("1")

	arr := c.Values()
	assert.Equal(t, []interface{}{1, 3, 2}, arr)
}

func TestSize(t *testing.T) {
	c := umi.New(nil)

	c.Set("1", 1)
	c.Set("2", 2)
	c.Set("3", 3)

	assert.Equal(t, 219, int(c.Size()))
}

func TestOverflow(t *testing.T) {
	c := umi.New(&umi.Options{
		MaxMemSize: 100,
	})

	c.Set("1", 1)
	c.Set("2", 2)
	c.Set("3", 3)
	c.Set("4", 4)

	arr := c.Values()

	assert.Equal(t, 73, int(c.Size()))
	assert.Equal(t, []interface{}{4}, arr)
}

func TestRace(t *testing.T) {
	c := umi.New(&umi.Options{
		TTL:    time.Microsecond * 5,
		GCSpan: time.Microsecond * 1,
	})

	vs := strings.Split("time.Sleep(time.Nanosecond * 10", "")
	l := len(vs)

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				time.Sleep(time.Nanosecond * 10)
				operator := rand.Int() % 4
				k := vs[rand.Int()%l]
				v := vs[rand.Int()%l]

				switch operator {
				case 0:
					c.Set(k, v)
				case 1:
					c.Get(k)
				case 2:
					c.Del(k)
				case 3:
					items := c.Items()

					for _, item := range items {
						if item == nil {
							fmt.Println(items)
							panic("shouldn't be nil")
						}
					}
				}
			}
		}()
	}

	time.Sleep(time.Second * 3)
}
