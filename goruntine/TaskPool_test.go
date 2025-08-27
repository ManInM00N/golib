package goruntine

import (
	"os"
	"testing"
	"time"
)

type ab struct {
	x, y int
}

func TestPriorityQueue(t *testing.T) {
	pool := NewTaskPool(5, 8)
	// defer pool.Stop()
	pool.Run()
	for i := 0; i < 10000; i++ {
		task, _ := pool.NewTask(func() {
			for i := 0; i < 1000; i++ {
				i += i % 2
			}
		}, nil, i)
		pool.Add(task)
	}
	pool.Wait()
}

func TestPriorityQueue_Priority(t *testing.T) {
	pool := NewTaskPool(5, 3)
	// defer pool.Stop()
	pool.Run()
	for i := 0; i < 10; i++ {
		task, _ := pool.NewTask(func() {
			t.Error("Task executed", i)
			time.Sleep((time.Duration)((i+3)/3) * time.Second)
		}, nil, i)
		pool.Add(task)
	}
	pool.Wait()
}

func TestMain(m *testing.M) {
	code := m.Run()

	// 退出测试
	os.Exit(code)
}
