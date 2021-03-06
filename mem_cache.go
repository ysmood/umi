package umi

import (
	"sync"
)

type memCache struct {
	sync.RWMutex

	maxSize uintptr
	size    uintptr

	list *memList

	// for fast read
	dict map[string]*Item

	onEvicted func(*Item)
}

type memList struct {
	sync.RWMutex

	len int

	// for fast insertion & deletion
	head *Item
	tail *Item
}

func (list *memList) add(item *Item) {
	list.Lock()

	list.len++

	if list.head == nil {
		list.head = item
		list.tail = item
		list.Unlock()
		return
	}

	item.next = list.head
	list.head.prev = item

	list.head = item

	list.Unlock()
}

func (list *memList) del(item *Item) {
	list.Lock()

	list.len--

	// if head
	if item.prev == nil {
		list.head = item.next
	} else {
		item.prev.next = item.next
	}

	// if tail
	if item.next == nil {
		list.tail = item.prev
	} else {
		item.next.prev = item.prev
	}

	item.next = nil
	item.prev = nil

	list.Unlock()
}

/*
	c is the target
	a -> b -> c -> d
	convert to
	a -> c -> b -> d
*/
func (list *memList) promote(item *Item, now int64) {
	list.Lock()

	c := item

	b := c.prev
	if b == nil {
		list.Unlock()
		return
	}

	a := b.prev
	d := c.next

	if a == nil {
		list.head = c
	} else {
		a.next = c
	}
	c.prev = a
	c.next = b
	b.prev = c
	b.next = d
	if d == nil {
		list.tail = b
	} else {
		d.prev = b
	}

	item.time = now

	list.Unlock()
}

func (mem *memCache) set(key string, val interface{}, now int64) *Item {
	item, has := mem.dict[key]

	// if the content already exists, replace it with the new one
	if has {
		oldVal := item.value
		oldSize := item.size
		item.value = val
		item.updateSize()

		if item.size > mem.maxSize {
			item.value = oldVal
			item.size = oldSize
			return nil
		}

		size := mem.size - oldSize + item.size
		if size > mem.maxSize {
			mem.free(size - mem.maxSize)
		}
		mem.size = mem.size - oldSize + item.size
		mem.list.promote(item, now)
	} else {
		item = newItem(key, val, now)

		if item.size > mem.maxSize {
			return nil
		}

		size := mem.size + item.size
		if size > mem.maxSize {
			mem.free(size - mem.maxSize)
		}

		mem.size += item.size
		mem.list.add(item)
		mem.dict[key] = item
	}

	return item
}

func (mem *memCache) del(item *Item) {
	if item == nil {
		return
	}

	delete(mem.dict, item.key)
	mem.list.del(item)
	mem.size -= item.size
}

func (mem *memCache) delTail() {
	tail := mem.list.tail
	mem.del(tail)
	if mem.onEvicted != nil {
		mem.onEvicted(tail)
	}
}

func (mem *memCache) purge() {
	for _, v := range mem.dict {
		mem.del(v)
	}
}

// free multiple items until the freed size reaches the specified size
func (mem *memCache) free(size uintptr) {
	var freedSize uintptr

	for freedSize < size {
		freedSize += mem.list.tail.size
		mem.delTail()
	}
}
