package umi

import (
	"github.com/ysmood/umi/lib"
)

const itemBaseSize = 8 + 16 + 8 + 8 + 8

// Sizable ...
type Sizable interface {
	Size() uintptr
}

// Aliveable ...
type Aliveable interface {
	Alive() bool
}

// Item ...
type Item struct {
	value interface{}
	size  uintptr
	key   string
	time  int64
	next  *Item
	prev  *Item
}

func newItem(key string, value interface{}, time int64) *Item {
	sizable, ok := value.(Sizable)

	var size uintptr
	if ok {
		size = itemBaseSize + lib.Size(key) + sizable.Size()
	} else {
		size = itemBaseSize + lib.Size(key) + lib.Size(value)
	}

	return &Item{
		value: value,
		size:  size,
		time:  time,
		key:   key,
	}
}

// Key ...
func (item *Item) Key() string {
	return item.key
}

// Value ...
func (item *Item) Value() interface{} {
	return item.value
}

// Time ...
func (item *Item) Time() int64 {
	return item.time
}
