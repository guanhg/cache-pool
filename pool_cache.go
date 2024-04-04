package pool

import (
	"sync"
	"time"
	"unsafe"
)

const _cacheSize = 10

type poolCache struct {
	count  int
	fulled chan struct{}
	tasks  []Task
}

func addCache(c chan Task, t Task, cache *poolCache, lock unsafe.Pointer) *poolCache {
	var ncache *poolCache
	var isfull = cache != nil && cache.count >= _cacheSize

	mux := (*sync.Mutex)(lock)
	mux.Lock()
	defer mux.Unlock()

	if cache == nil || isfull {
		ncache = &poolCache{
			count:  1,
			fulled: make(chan struct{}),
			tasks:  []Task{t},
		}
		if isfull {
			cache.full()
		}
		cache = ncache
		go ncache.await(c, lock)
		return ncache
	}

	cache.count++
	cache.tasks = append(cache.tasks, t)
	return cache
}

func (pc *poolCache) full() {
	close(pc.fulled)
}

func (pc *poolCache) await(c chan Task, lock unsafe.Pointer) {
	mux := (*sync.Mutex)(lock)

	select {
	case <-pc.fulled:
		for _, t := range pc.tasks {
			c <- t
		}
	case <-time.After(time.Second):
		var ts []Task
		if len(pc.tasks) > 0 {
			mux.Lock()
			ts = append(ts, pc.tasks...)
			pc.tasks = pc.tasks[:0]
			mux.Unlock()
		}

		for _, t := range ts {
			c <- t
		}
	}
}
