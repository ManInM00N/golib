package goruntine

import (
	"fmt"
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
	pool := NewTaskPool(4, 8, WithFIFO())
	// defer pool.Stop()
	for i := 0; i < 10; i++ {
		val := i
		task, _ := pool.NewTaskWithCost(func() {
			fmt.Println("Task executed", val)
			time.Sleep(time.Second * 3)
		}, &ab{i / 4, i % 4}, 1, 3)
		pool.Add(task)
		time.Sleep(time.Millisecond * 200)
	}
	pool.RemoveTaskByDeleteFunc(func(item Task) bool {
		info := (*item.GetInfo()).(*ab)
		return info.x == 2 && info.y < 2
	})
	a, _ := pool.GetTaskStatistic()
	for _, v := range a {
		fmt.Println(v.Info, v.InQueueTime)
	}
	pool.Run()
	pool.Wait()
	t.Error(111)
	/*
		expect output(RANDOM ORDER):
		Task executed 0
		Task executed 1
		Task executed 2
		Task executed 3
		// must be [0,3]
		Task executed 4
		Task executed 5
		Task executed 6
		Task executed 7
		// [8,9] removed
	*/
}

func TestWeightTask(t *testing.T) {
	pool := NewTaskPool(4, 10, WithFIFO())
	// defer pool.Stop()
	pool.Run()
	for i := 0; i < 20; i++ {
		weight := (i % 5) + 1
		val := i
		task, _ := pool.NewTaskWithCost(func() {
			fmt.Println("Task executed", val, "with weight", weight)
			arr, wk := pool.GetTaskStatistic()
			for _, v := range arr {
				fmt.Println("  In queue:", v.Info, "Weight:", v.Weight, "InQueueTime:")
			}
			for id, w := range wk {
				fmt.Println("  Worker:", id, "Status:", w)
			}
			time.Sleep(time.Second * 2)
		}, i, 1, weight)
		pool.Add(task)
	}
	pool.Wait()
	t.Error("All tasks completed")

}

func TestMain(m *testing.M) {
	code := m.Run()

	// 退出测试
	os.Exit(code)
}
