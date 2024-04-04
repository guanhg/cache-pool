package pool

import (
	"testing"
)

func TestPool(t *testing.T) {
	p := NewPool(5)
	defer p.Close()

	p.Start()
	gs, size := 50, 13
	for j := 0; j < gs; j++ {
		go func() {
			for i := 0; i < size; i++ {
				p.AddTask(GetNonTask())
			}
		}()
	}

	if err := p.Wait(); err != nil {
		t.Error(err)
	}

	count := int(p.count.Load())
	if count != size*gs {
		t.Errorf("count: %d, except: %d", count, size*gs)
	}

	for i := 0; i < 10; i++ {
		invalidTask := InvalidTask{Name: "test task"}
		p.AddTask(&invalidTask)
	}

	if err := p.Wait(); err != nil {
		t.Error(err)
	}
}
