package goruntine

import (
	"container/heap"
	"context"

	. "github.com/ManInM00N/go-tool/heap"
)

type TaskPool struct {
	n     int
	c     chan task
	queue *PriorityQueue[task]
}
type task struct {
	val   int
	Inner func()
	ctx   context.Context
}

func (t *task) Cancel() {
}

func NewTask(val int) (task, context.CancelFunc) {
	var task task
	task.val = val
	task.Inner = func() {
		// do something
	}
	ctx, cancel := context.WithCancel(context.Background())
	task.ctx = ctx
	return task, cancel
}

func NewTaskPool(worker int) *TaskPool {
	temp := &TaskPool{
		n: worker,
		c: make(chan task),
		queue: NewPriorityQueue[task](func(i, j task) bool {
			return i.val < j.val
		}),
	}
	return temp
}

func (P *TaskPool) Run() {
	defer close(P.c)
	for {
		select {
		case v := <-P.c:
			{
				heap.Push(P.queue, v)
			}
		default:
			{
			}
		}
	}
}
