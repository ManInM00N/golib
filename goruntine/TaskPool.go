package goruntine

import (
	"context"
	"sync"

	. "github.com/ManInM00N/go-tool/heap"
)

type TaskPool struct {
	n             int
	c             chan task
	wg            sync.WaitGroup
	queue         *PriorityQueue[task]
	cond          *sync.Cond
	maxConcurrent int
	mu            sync.Mutex
	sem           chan struct{}
	running       bool
}
type task struct {
	val   int
	Inner func()
	ctx   context.Context
	Info  interface{}
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

func NewTaskPool(worker int, maxConcurrent int) *TaskPool {
	temp := &TaskPool{
		n:   worker,
		c:   make(chan task),
		sem: make(chan struct{}, maxConcurrent),
		queue: NewPriorityQueue[task](func(i, j task) bool {
			return i.val < j.val
		}),
	}
	return temp
}

func (P *TaskPool) Run() {
	P.running = true
	for i := 0; i < P.n; i++ {
		go P.worker()
	}
}
func (p *TaskPool) worker() {
	for {
		p.mu.Lock()
		for p.queue.Len() == 0 && p.running {
			p.cond.Wait()
		}
		if !p.running && p.queue.Len() == 0 {
			p.mu.Unlock()
			return
		}
		t := p.queue.Pop().(task)
		p.mu.Unlock()

		p.sem <- struct{}{}
		p.wg.Add(1)
		go func(t task) {
			defer func() {
				<-p.sem
				p.wg.Done()
			}()
			t.Inner()
		}(t)
	}
}

func (g *TaskPool) AddWorker(num int) {
	g.wg.Add(num)
}
func (g *TaskPool) NewTask(fn func(), info interface{}, val int) (task, context.CancelFunc) {
	t, cancel := NewTask(val)
	t.Inner = fn
	t.Info = info
	return t, cancel
}
func (p *TaskPool) Add(t task) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.queue.Push(t)
	p.cond.Signal()
}

func (p *TaskPool) Stop() {
	p.mu.Lock()
	p.running = false
	p.mu.Unlock()
	p.cond.Broadcast()
	p.wg.Wait()
}
func (g *TaskPool) Pop() {
	<-g.c
}
func (g *TaskPool) Done() {
	g.wg.Done()
}
func (g *TaskPool) Wait() {
	g.wg.Wait()
}
func (g *TaskPool) Close() {
	close(g.c)
}
