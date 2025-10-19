package goruntine

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	. "github.com/ManInM00N/go-tool/heap"
	"github.com/teris-io/shortid"
)

type TaskPoolOption func(*TaskPool)

type TaskPool struct {
	n             int
	c             chan Task
	wg            sync.WaitGroup
	queue         *PriorityQueue[Task]
	cond          *sync.Cond
	maxConcurrent int
	mu            sync.Mutex
	sem           chan struct{}
	running       bool
	workersNum    atomic.Int32
	workers       map[string]*worker
}

type Task struct {
	inQueueTime time.Time       `json:"in_queue_time"`
	val         int             `json:"val"`
	inner       func()          `json:"-"`
	ctx         context.Context `json:"-"`
	info        *interface{}    `json:"info"`
	status      int             `json:"status"` // 0未执行 1执行中 2完成 -1取消
}

const (
	TaskStatusPending   = 0
	TaskStatusRunning   = 1
	TaskStatusCompleted = 2
	TaskStatusCanceled  = -1
)

// Getter 方法 - 提供只读访问
func (t *Task) GetVal() int {
	return t.val
}

func (t *Task) GetInfo() *interface{} {
	return t.info
}

func (t *Task) GetStatus() int {
	return t.status
}

func (t *Task) GetContext() context.Context {
	return t.ctx
}

func (t *Task) setStatus(status int) {
	t.status = status
}

func (t *Task) Cancel() {
	t.status = TaskStatusCanceled
	if t.ctx != nil {
		t.Cancel()
	}
}

func NewTask(val int) (Task, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	task := Task{
		val:    val,
		inner:  func() {},
		ctx:    ctx,
		status: 0,
	}
	return task, cancel
}

// NewTaskPool 创建任务池
// workers: worker 数量
// maxConcurrent: 最大并发数
// opts: heap.Option[task] 用于配置优先队列（可以使用 heap.WithLessFunc, heap.WithEqualFunc 等）
func NewTaskPool(workers int, maxConcurrent int, opts ...Option[Task]) *TaskPool {
	defaultOpts := []Option[Task]{
		WithLowestValueFirst(),
		WithTaskEqualityByVal(),
	}

	allOpts := append(defaultOpts, opts...)

	temp := &TaskPool{
		n:       workers,
		c:       make(chan Task),
		sem:     make(chan struct{}, maxConcurrent),
		cond:    sync.NewCond(&sync.Mutex{}),
		workers: make(map[string]*worker),
		queue:   NewPriorityQueueWithOptions(allOpts...),
		mu:      sync.Mutex{},
	}

	return temp
}

func (p *TaskPool) SetFunc(opt Option[Task]) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.queue.SetFunc(opt)
	return true
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
			status: WorkerStatusIdle,
		}

		p.workersNum.Add(1)
		p.workers[w.id] = w
		go w.Run()
	}
}

func (g *TaskPool) AddCount(num int) {
	g.wg.Add(num)
}

func (g *TaskPool) NewTask(fn func(), info interface{}, val int) (Task, context.CancelFunc) {
	t, cancel := NewTask(val)
	t.inner = fn
	t.info = &info
	return t, cancel
}

// RemoveTask 删除满足条件的等待任务
// 使用自定义的删除函数
//
//	pool.RemoveTaskByDeleteFunc(func(item Task) bool {
//			info := (*item.GetInfo()).(*ab)
//			return info.x == 2 && info.y < 2
//	})
func (p *TaskPool) RemoveTaskByDeleteFunc(deleteFunc func(item Task) bool) int {
	p.mu.Lock()
	defer p.mu.Unlock()
	num := p.queue.RemoveByDeleteFunc(deleteFunc)
	p.AddCount(-num)
	return num
}

// RemoveTask 删除满足条件的等待任务
func (p *TaskPool) RemoveEqualTask(item Task) int {
	p.mu.Lock()
	defer p.mu.Unlock()
	num := p.queue.RemoveEqual(item)
	p.AddCount(-num)
	return num
}

func (p *TaskPool) Add(t Task) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.AddCount(1)
	t.inQueueTime = time.Now()
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

func (p *TaskPool) GetTaskStatistic() ([]Task, map[string]interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	stats := make(map[string]interface{})
	arr := p.queue.Items()
	for k, v := range p.workers {
		t := map[string]interface{}{
			"status": v.GetStatus(),
		}
		if v.task != nil {
			t["task"] = v.GetTaskInfo()
		} else {
			t["task"] = "No task"
		}
		stats[k] = t
	}
	return arr, stats
}
