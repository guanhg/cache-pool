package pool

// 协程池

import (
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"
)

// 带缓存的协程池，缓存会被垃圾回收
type Pool struct {
	size  int           // 协程数
	count atomic.Uint32 // 任务总数

	pc     chan Task
	cache  *poolCache // 任务缓存
	mux    sync.Mutex
	closed bool

	wg        sync.WaitGroup
	startOnce sync.Once
	closeOnce sync.Once
}

func NewPool(gSize int) *Pool {
	return &Pool{
		size: gSize,
		pc:   make(chan Task, gSize),
	}
}

func (p *Pool) Start() {
	// task execute
	p.startOnce.Do(func() {
		for i := 0; i < p.size; i++ {
			go func() {
				for task := range p.pc {
					task.Execute()
					p.count.Add(1)
					p.wg.Done()
				}
			}()
		}
	})
}

func (p *Pool) AddTask(t Task) error {
	if p.closed {
		return fmt.Errorf("Pool is closed")
	}
	p.wg.Add(1)
	addCache(p.pc, t, p.cache, unsafe.Pointer(&p.mux))
	return nil
}

func (p *Pool) Wait() error {
	p.wg.Wait()
	return nil
}

func (p *Pool) Close() {
	p.closeOnce.Do(
		func() {
			p.closed = true
			close(p.pc)
		})
}
