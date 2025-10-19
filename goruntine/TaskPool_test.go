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

func TestRemoveTask(t *testing.T) {
	pool := NewTaskPool(4, 4, WithFIFO())
	// defer pool.Stop()
	for i := 0; i < 10; i++ {
		val := i
		task, _ := pool.NewTask(func() {
			t.Error("Task executed", val)
			time.Sleep(2 * time.Second)
		}, &ab{i / 4, i % 4}, 1)
		pool.Add(task)
	}
	pool.RemoveTaskByDeleteFunc(func(item Task) bool {
		info := (*item.GetInfo()).(*ab)
		return info.x == 2 && info.y < 2
	})

	pool.Run()
	pool.Wait()
	/*
		expect output(RANDOM ORDER):
		Task executed 0
		Task executed 1
		Task executed 2
		Task executed 3
		Task executed 4
		Task executed 5
		Task executed 6
		Task executed 7
	*/
}

func TestMain(m *testing.M) {
	code := m.Run()

	// 退出测试
	os.Exit(code)
}
