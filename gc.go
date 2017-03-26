package umi

import (
	"time"
)

func (c *Cache) gcWorker(span time.Duration, gcSize int, ttl time.Duration) {
	for {
		time.Sleep(span)

		l := len(c.mem.dict)

		var items []*Item
		if l > gcSize {
			items = c.Slice(l-gcSize, l)
		} else {
			items = c.Slice(0, l)
			gcSize = l
		}

		var item *Item
		for i := 0; i < gcSize; i++ {
			item = items[i]

			aliveable, ok := item.value.(Aliveable)
			var alive bool
			if ok {
				alive = aliveable.Alive()
			} else {
				alive = (c.now - item.time) < int64(ttl)
			}

			if !alive {
				c.mem.Lock()
				c.mem.del(item)
				c.mem.Unlock()
			}
		}

		c.now = time.Now().UnixNano()
	}
}
