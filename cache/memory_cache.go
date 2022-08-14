package cache

import (
	"sync"
	"time"
)

//-------memory cached map with a expiring time

type mapEntry[T any] struct {
	last_visit time.Time
	data       T
}

type MemoryCache[T any] interface {
	Set(int, T)
	Get(int) (T, bool)
	Delete(int)
}

type cacheMap[T any] struct {
	mutex           sync.RWMutex
	mp              map[int]*mapEntry[T]
	cacheExpireTime time.Duration
	maxCheckLength  int
}

func NewMemoryCache[T any](expire time.Duration, check_length int) MemoryCache[T] {
	return &cacheMap[T]{sync.RWMutex{}, make(map[int]*mapEntry[T]), expire, check_length}
}

func (cm *cacheMap[T]) Set(key int, val T) {
	cm.mutex.Lock()
	cm.mp[key] = &mapEntry[T]{time.Now(), val}
	if len(cm.mp) >= cm.maxCheckLength {
		current := time.Now().Add(-cm.cacheExpireTime)
		for i := range cm.mp {
			if cm.mp[i].last_visit.Before(current) {
				delete(cm.mp, i)
			}
		}
	}
	cm.mutex.Unlock()
}

func (cm *cacheMap[T]) Get(key int) (T, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	a, ok := cm.mp[key]
	if !ok {
		return *new(T), false
	}
	a.last_visit = time.Now()
	return a.data, true
}

func (cm *cacheMap[T]) Delete(key int) {
	cm.mutex.Lock()
	delete(cm.mp, key)
	cm.mutex.Unlock()
}