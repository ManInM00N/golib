package goruntine

import (
	"container/heap"

	. "github.com/ManInM00N/go-tool/heap"
)

type TaskPool struct {
	n     int
	c     chan Task
	queue PriorityQueue[Task]
}

func NewTaskPool(worker int) *TaskPool {
	temp := &TaskPool{
		n: worker,
		c: make(chan Task),
	}
	return temp
}

func (P *TaskPool) Run() {
	defer close(P.c)
	for {
		select {
		case v := <-P.c:
			{
				heap.Push(&P.queue, v)
			}
		default:
			{
			}
		}
	}
}
