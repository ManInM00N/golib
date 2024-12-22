package Map

import "sync"

type ChannelMap[K comparable, V any] struct {
	lock chan any
	mp   map[K]V
}

func (mm *ChannelMap[K, V]) Lock() {
	mm.lock <- nil
}

func (mm *ChannelMap[K, V]) UnLock() {
	<-mm.lock
}

func (mm *ChannelMap[K, V]) Delete(k K) V {
	mm.Lock()
	defer mm.UnLock()
	v := mm.mp[k]
	delete(mm.mp, k)
	return v
}

func (mm *ChannelMap[K, V]) Get(k K) (any, bool) {
	mm.Lock()
	defer mm.UnLock()
	v, has := mm.mp[k]
	return v, has
}

func (mm *ChannelMap[K, V]) Set(k K, v V) {
	mm.Lock()
	defer mm.UnLock()
	mm.mp[k] = v
}

func NewChannelMap[T comparable, V any]() *ChannelMap[T, V] {
	tmp := &ChannelMap[T, V]{
		lock: make(chan any, 1),
		mp:   make(map[T]V),
	}
	return tmp
}

type RWMap[K comparable, V any] struct {
	lock sync.RWMutex
	mp   map[K]V
}

func NewRWMap[K comparable, V any]() *RWMap[K, V] {
	return &RWMap[K, V]{
		mp: make(map[K]V),
	}
}

func (mm *RWMap[K, V]) Set(k K, v V) {
	mm.lock.Lock()
	defer mm.lock.Unlock()
	mm.mp[k] = v
}

func (mm *RWMap[K, V]) Get(k K) (V, bool) {
	mm.lock.RLock()
	defer mm.lock.RUnlock()
	v, has := mm.mp[k]
	return v, has
}

func (mm *RWMap[K, V]) Delete(k K) V {
	mm.lock.Lock()
	defer mm.lock.Unlock()
	v := mm.mp[k]
	delete(mm.mp, k)
	return v
}
