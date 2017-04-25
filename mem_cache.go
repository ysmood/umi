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
}

type memList struct {
	sync.RWMutex

	len int

	// for fast insertion & deletion
	head *Item
	tail *Item
}

func (l *memList) add(item *Item) {
	l.Lock()

	l.len++

	if l.head == nil {
		l.head = item
		l.tail = l.head
		l.Unlock()
		return
	}

	item.next = l.head
	l.head.prev = item

	l.head = l.head.prev

	l.Unlock()
}

func (l *memList) del(item *Item) {
	l.Lock()

	l.len--

	// if head
	if item.prev == nil {
		l.head = item.next
	} else {
		item.prev.next = nil
	}

	// if tail
	if item.next == nil {
		l.tail = item.prev
	} else {
		item.next.prev = item.prev
		item.next = nil
	}

	item.prev = nil

	l.Unlock()
}

/*
	c is the target
	a -> b -> c -> d
	convert to
	a -> c -> b -> d
*/
func (l *memList) promote(item *Item, now int64) {
	l.Lock()

	c := item

	b := c.prev
	if b == nil {
		l.Unlock()
		return
	}

	a := b.prev
	d := c.next

	if a == nil {
		l.head = c
	} else {
		a.next = c
	}
	c.prev = a
	c.next = b
	b.prev = c
	b.next = d
	if d == nil {
		l.tail = b
	} else {
		d.prev = b
	}

	item.time = now

	l.Unlock()
}

func (m *memCache) delTail() {
	m.del(m.list.tail)
}

func (m *memCache) del(item *Item) {
	if item == nil {
		return
	}

	delete(m.dict, item.key)
	m.list.del(item)
	m.size -= item.size
}

// free multiple items until the freed size reaches the specified size
func (m *memCache) free(size uintptr) bool {
	var freedSize uintptr

	for freedSize < size {
		if m.list.tail == nil {
			// if after all items are freed, the space is still not enough
			return false
		}
		freedSize += m.list.tail.size
		m.delTail()
	}

	return true
}

func (m *memCache) set(key string, item *Item) {
	size := m.size + item.size

	if size > m.maxSize {
		if !m.free(size - m.maxSize) {
			return
		}
	}

	oldItem, has := m.dict[key]

	// if the content already exists, replace it with the new one
	if has {
		m.del(oldItem)
	}
	m.list.add(item)

	m.dict[key] = item

	m.size += item.size
}
