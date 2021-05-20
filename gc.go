package umi

import (
	"time"
)

func (c *Cache) gcWorker(span time.Duration, gcSize int, ttl time.Duration) {
	for {
		time.Sleep(span)

		c.mem.Lock()

		for item, count := c.mem.list.tail, 0; item != nil && count < gcSize; count++ {
			aliveable, ok := item.value.(Aliveable)
			var alive bool
			if ok {
				alive = aliveable.Alive()
			} else {
				alive = (c.now - item.time) < int64(ttl)
			}

			if alive {
				break
			}

			c.mem.del(item)

			item = item.prev
		}

		c.now = time.Now().UnixNano()

		c.mem.Unlock()
	}
}
