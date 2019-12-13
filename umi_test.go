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

func TestSet(t *testing.T) {
	c := umi.New(nil)

	c.Set("a", 1)
	c.Set("b", 2)
	c.Set("c", 3)

	assert.Equal(t, []interface{}{3, 2, 1}, c.Values())
}

func TestBasic(t *testing.T) {
	c := umi.New(nil)

	item := c.Set("a", int32(10))
	now := time.Now().UnixNano()
	assert.True(t, item.Time() < now && item.Time() > now-int64(time.Second))
	assert.Equal(t, 69, int(item.Size()))

	v, _ := c.Get("a")

	assert.Equal(t, 10, int(v.(int32)))
}

type testType struct {
	size uintptr
}

func (t *testType) Size() uintptr {
	return t.size
}

func (t *testType) Alive() bool {
	return true
}

var _ umi.Sizable = &testType{}
var _ umi.Aliveable = &testType{}

func TestSizable(t *testing.T) {
	c := umi.New(nil)

	c.Set("a", &testType{1000})

	assert.Equal(t, 1065, int(c.Size()))
}

func TestPeek(t *testing.T) {
	c := umi.New(nil)

	c.Set("a", int32(10))

	v, _ := c.Peek("a")

	assert.Equal(t, 10, int(v.(int32)))

	_, has := c.Peek("b")
	assert.False(t, has)
}

func TestDelHead(t *testing.T) {
	c := umi.New(nil)

	c.Set("a", 1)
	c.Set("b", 2)
	c.Set("c", 3)

	c.Del("c")

	values := c.Values()

	assert.Equal(t, []interface{}{2, 1}, values)
	assert.Equal(t, 2, c.Count())
	assert.Equal(t, []string{"b", "a"}, c.Keys())
}

func TestDelMiddle(t *testing.T) {
	c := umi.New(nil)

	c.Set("a", 1)
	c.Set("b", 2)
	c.Set("c", 3)

	c.Del("b")

	values := c.Values()

	assert.Equal(t, []interface{}{3, 1}, values)
}

func TestDelTail(t *testing.T) {
	c := umi.New(nil)

	c.Set("a", 1)
	c.Set("b", 2)
	c.Set("c", 3)

	c.Del("a")

	values := c.Values()

	assert.Equal(t, []interface{}{3, 2}, values)
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
	assert.Equal(t, uint64(0), c.Size())
}

func TestSize(t *testing.T) {
	c := umi.New(nil)

	c.Set("1", 1)
	c.Set("2", 2)
	c.Set("3", 3)

	assert.Equal(t, 219, int(c.Size()))
}

func TestUpdateSize(t *testing.T) {
	c := umi.New(nil)

	c.Set("1", "a")
	c.Set("1", "ab")

	assert.Equal(t, 83, int(c.Size()))
}

func TestReduceSize(t *testing.T) {
	c := umi.New(nil)

	c.Set("1", "a")
	c.Set("2", "b")
	c.Del("2")

	assert.Equal(t, 82, int(c.Size()))
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

func TestOverflowNotSet(t *testing.T) {
	c := umi.New(&umi.Options{
		MaxMemSize: 200,
	})

	large := make([]byte, 200)

	c.Set("1", 1)
	res := c.Set("2", large)
	assert.Nil(t, res)

	c.Set("1", 1)
	res = c.Set("1", large)
	assert.Nil(t, res)

	large = make([]byte, 100)

	c.Set("1", 1)
	c.Set("2", 1)
	c.Set("1", large)
	assert.Equal(t, 1, c.Count())

	c.Set("1", 1)
	c.Set("2", 1)
	c.Set("3", large)
	assert.Equal(t, 1, c.Count())
}

func TestOnEvicted(t *testing.T) {
	c := umi.New(&umi.Options{
		MaxMemSize: 100,
		OnEvicted: func(item *umi.Item) {
			assert.Equal(t, "1", item.Key())
		},
	})

	c.Set("1", 1)
	c.Set("2", 2)
}

func TestRace(t *testing.T) {
	c := umi.New(&umi.Options{
		TTL:    time.Microsecond * 5,
		GCSpan: time.Microsecond * 1,
	})

	vs := strings.Split("花に嵐のたとえもあるさ さよならだけが人生", "")
	l := len(vs)

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				time.Sleep(time.Nanosecond * 10)
				operator := rand.Int() % 8
				k := vs[rand.Int()%l]
				v := vs[rand.Int()%l]

				switch operator {
				case 0:
					fallthrough
				case 1:
					c.Set(k, v)
				case 2:
					fallthrough
				case 3:
					fallthrough
				case 4:
					c.Get(k)
				case 5:
					c.Del(k)
				case 6:
					items := c.Items()

					for _, item := range items {
						if item == nil {
							fmt.Println(items)
							panic("shouldn't be nil")
						}
					}
				case 7:
					c.Purge()
				}
			}
		}()
	}

	time.Sleep(time.Second * 3)
}
