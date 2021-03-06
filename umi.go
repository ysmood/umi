package umi

import (
	"time"
)

// Cache ...
type Cache struct {
	mem *memCache
	now int64
}

// New ...
func New(opts *Options) *Cache {
	opts = defaultOptions(opts)

	c := &Cache{
		now: time.Now().UnixNano(),
		mem: &memCache{
			list:      &memList{},
			maxSize:   uintptr(opts.MaxMemSize),
			dict:      make(map[string]*Item),
			onEvicted: opts.OnEvicted,
		},
	}

	if opts.GCSpan > 0 && opts.GCSize > 0 {
		go c.gcWorker(opts.GCSpan, opts.GCSize, opts.TTL)
	}

	return c
}

// Size ...
func (c *Cache) Size() uint64 {
	return uint64(c.mem.size)
}

// Count ...
func (c *Cache) Count() int {
	return c.mem.list.len
}

// Set the val parameter could be `umi.IItem`, which will overwrite
// the default behavior.
func (c *Cache) Set(key string, val interface{}) *Item {
	c.mem.Lock()
	item := c.mem.set(key, val, c.now)
	c.mem.Unlock()

	return item
}

// Del ...
func (c *Cache) Del(key string) {
	c.mem.Lock()
	c.mem.del(c.mem.dict[key])
	c.mem.Unlock()
}

// Purge ...
func (c *Cache) Purge() {
	c.mem.Lock()
	c.mem.purge()
	c.mem.Unlock()
}

// Get ...
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mem.RLock()

	item, has := c.mem.dict[key]

	if has {
		c.mem.list.promote(item, c.now)
		v := item.value
		c.mem.RUnlock()
		return v, has
	}

	c.mem.RUnlock()
	return nil, has
}

// Peek it wont' affect the promotion
func (c *Cache) Peek(key string) (interface{}, bool) {
	c.mem.RLock()

	item, has := c.mem.dict[key]

	if has {
		c.mem.RUnlock()
		return item.value, has
	}

	c.mem.RUnlock()
	return nil, has
}

// Keys all keys from head to tail
func (c *Cache) Keys() []string {
	c.mem.list.RLock()

	head := c.mem.list.head
	arr := make([]string, c.mem.list.len)
	i := 0

	for head != nil {
		arr[i] = head.Key()
		i++
		head = head.next
	}

	c.mem.list.RUnlock()

	return arr
}

// Values all values from head to tail
func (c *Cache) Values() []interface{} {
	c.mem.list.RLock()

	head := c.mem.list.head
	arr := make([]interface{}, c.mem.list.len)
	i := 0

	for head != nil {
		arr[i] = head.Value()
		i++
		head = head.next
	}

	c.mem.list.RUnlock()

	return arr
}

// Slice from head to tail
func (c *Cache) Slice(begin int, end int) []*Item {
	c.mem.list.RLock()

	item := c.mem.list.head
	l := end - begin
	items := make([]*Item, l)

	for i := 0; i < l; i++ {
		items[i] = item

		if item != nil {
			item = item.next
		}
	}

	c.mem.list.RUnlock()

	return items
}

// Items all items from head to tail
func (c *Cache) Items() []*Item {
	c.mem.list.RLock()

	items := make([]*Item, c.mem.list.len)

	item := c.mem.list.head
	i := 0

	for item != nil {
		items[i] = item
		item = item.next
		i++
	}

	c.mem.list.RUnlock()

	return items
}
