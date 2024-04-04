package pool

import (
	"fmt"
	"math/rand"
	"time"
)

type Task interface {
	Execute() error
}

var _ Task = PoolTask{}

type PoolTask struct {
	Name   string
	Handle func() error
}

func (t PoolTask) Execute() error {
	if t.Handle == nil {
		return fmt.Errorf("PoolTask.Handle is nil")
	}
	return t.Handle()
}

func (t PoolTask) GetName() string {
	return t.Name
}

var _ Task = (*nonTask)(nil)

type nonTask struct{}

func (nt *nonTask) Execute() error {
	return nil
}

func GetNonTask() Task {
	return &nonTask{}
}

// for test
type InvalidTask struct {
	Name string
}

func (it *InvalidTask) Execute() error {
	rn := float32(rand.Intn(4)) / 4
	time.Sleep(time.Duration(rn*1000) * time.Millisecond)
	if rn < 0.5 {
		return fmt.Errorf("testHandle: %f", rn)
	}
	return nil
}
