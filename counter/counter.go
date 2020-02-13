package counter

import (
	"sync"
	"sync/atomic"
)

type Counter struct {
	mu           sync.Mutex
	statisticMap sync.Map //string int32  key count
}

func NewCounter() *Counter {
	return &Counter{statisticMap: sync.Map{}, mu: sync.Mutex{}}
}

var zero int32 = 0

func (c *Counter) Get(key string) *int32 {
	if key == "" {
		zero = 0
		return &zero
	}
	defer c.mu.Unlock()
	c.mu.Lock()
	if v, ok := c.statisticMap.Load(key); ok {
		return v.(*int32)
	}
	var i int32
	c.statisticMap.Store(key, &i)
	return &i
}

func (c *Counter) Count(key string) int32 {
	return *c.Get(key)
}

func (c *Counter) Inc(key string) {
	atomic.AddInt32(c.Get(key), 1)
}

func (c *Counter) Statistic() (result map[string]int32) {
	result = make(map[string]int32)
	var total int32 = 0
	c.statisticMap.Range(func(key, value interface{}) bool {
		v := *value.(*int32)
		result[key.(string)] = v
		total = total + v
		return true
	})
	result["total"] = total
	return
}

func (c *Counter) Clear() {
	defer c.mu.Unlock()
	c.mu.Lock()
	c.statisticMap.Range(func(key, value interface{}) bool {
		c.statisticMap.Delete(key)
		return true
	})
}
