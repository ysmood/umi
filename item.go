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

func newItem(key string, value interface{}, time int64) (item *Item) {
	item = &Item{
		value: value,
		time:  time,
		key:   key,
	}
	item.updateSize()
	return
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

// updateSize ...
func (item *Item) updateSize() {
	sizable, ok := item.value.(Sizable)

	var size uintptr
	if ok {
		size = itemBaseSize + lib.Size(item.key) + sizable.Size()
	} else {
		size = itemBaseSize + lib.Size(item.key) + lib.Size(item.value)
	}
	item.size = size
}

// Size ...
func (item *Item) Size() uint64 {
	return uint64(item.size)
}
