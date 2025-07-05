package cargo

import (
	"sync"
)

type KeyValue[K comparable, V any] interface {
	Has(key K) bool
	Get(key K) (V, bool)
	Set(key K, value V)
	Del(key K)
	Clear()
	Map() map[K]V
	Clone() KeyValue[K, V]
}

type Collection[K comparable, V any] struct {
	dic   map[K]V
	mutex sync.RWMutex
	copy  func(V) V
}

func NewCollection[K comparable, V any](copy func(V) V) *Collection[K, V] {
	if copy == nil {
		copy = func(v V) V { return v }
	}
	return &Collection[K, V]{
		dic:   make(map[K]V),
		mutex: sync.RWMutex{},
		copy:  copy,
	}
}

func (c *Collection[K, V]) Has(key K) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	_, ok := c.dic[key]
	return ok
}

func (c *Collection[K, V]) Get(key K) (V, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	value, ok := c.dic[key]
	return c.copy(value), ok
}

func (c *Collection[K, V]) Set(key K, value V) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.dic[key] = value
}

func (c *Collection[K, V]) Del(key K) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.dic, key)
}

func (c *Collection[K, V]) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.dic = make(map[K]V)
}

func (c *Collection[K, V]) Map() map[K]V {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	dic := make(map[K]V, len(c.dic))
	for key, value := range c.dic {
		dic[key] = c.copy(value)
	}
	return dic
}

func (c *Collection[K, V]) Clone() KeyValue[K, V] {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	clone := NewCollection[K](c.copy)
	for key, value := range c.dic {
		clone.Set(key, c.copy(value))
	}
	return clone
}
