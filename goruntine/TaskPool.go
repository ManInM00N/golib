package goruntine

import (
	"context"
	"sync"
	"sync/atomic"

	. "github.com/ManInM00N/go-tool/heap"
	"github.com/teris-io/shortid"
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
	workersNum    atomic.Int32
	workers       map[string]*worker
}

type task struct {
	val    int
	Inner  func()
	ctx    context.Context
	Info   *interface{}
	status int
}

func (t *task) GetStatus() int {
	return t.status
}

func (t *task) SetStatus(status int) {
	t.status = status
}

func (t *task) Cancel() {
	t.status = -1
	t.ctx.Done()
}

func NewTask(val int) (task, context.CancelFunc) {
	var task task
	task.val = val
	task.Inner = func() {
	}
	ctx, cancel := context.WithCancel(context.Background())
	task.ctx = ctx
	return task, cancel
}

func NewTaskPool(workers int, maxConcurrent int) *TaskPool {
	temp := &TaskPool{
		n:       workers,
		c:       make(chan task),
		sem:     make(chan struct{}, maxConcurrent),
		cond:    sync.NewCond(&sync.Mutex{}),
		workers: make(map[string]*worker),
		queue: NewPriorityQueue[task](func(i, j task) bool {
			return i.val > j.val
		}),
		mu: sync.Mutex{},
	}
	return temp
}

func (P *TaskPool) Run() {
	P.running = true
	P.AddWorker(P.n)
}

func (p *TaskPool) AddWorker(num int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i := 0; i < num; i++ {
		id, _ := shortid.Generate()
		w := &worker{
			id:     id,
			pool:   p,
			task:   nil,
			status: 0,
		}

		p.workersNum.Add(1)
		p.workers[w.id] = w
		go w.Run()
	}
}

func (g *TaskPool) AddCount(num int) {
	g.wg.Add(num)
}
func (g *TaskPool) NewTask(fn func(), info interface{}, val int) (task, context.CancelFunc) {
	t, cancel := NewTask(val)
	t.Inner = fn
	t.Info = &info
	return t, cancel
}
func (p *TaskPool) Add(t task) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.AddCount(1)
	p.queue.Push(t)
	p.cond.Signal()
}

func (p *TaskPool) Stop() {
	p.mu.Lock()
	p.running = false
	p.mu.Unlock()
	p.wg.Wait()
}
func (p *TaskPool) Pop() {
}
func (p *TaskPool) Done() {
	p.wg.Done()
}
func (p *TaskPool) Wait() {
	p.wg.Wait()
}
func (p *TaskPool) Close() {
	close(p.c)
}

func (p *TaskPool) Lock() {
	p.mu.Lock()
}

func (p *TaskPool) Unlock() {
	p.mu.Unlock()
}

func (p *TaskPool) GetWorkers() map[string]*worker {
	return p.workers
}

func (p *TaskPool) GetTaskStatistic() ([]task, map[string]interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	stats := make(map[string]interface{})
	arr := p.queue.Items()
	for k, v := range p.workers {
		stats[k] = map[string]interface{}{
			"status": v.GetStatus(),
			"task":   v.GetTaskInfo(),
		}
	}
	return arr, stats
}
