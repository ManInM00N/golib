package goruntine

import "time"

const (
	WorkerStatusIdle    = 0
	WorkerStatusRunning = 1
	WorkerStatusStop    = -1
)

type worker struct {
	id     string
	pool   *TaskPool
	task   *Task
	status int //0空闲 1工作 -1停止
}

func (p *worker) Run() {
	p.pool.workersNum.Add(1)
	for {
		p.pool.mu.Lock()

		if !p.pool.running && p.pool.queue.Len() == 0 {
			p.status = WorkerStatusStop
			p.pool.mu.Unlock()
			break
		}
		if p.pool.queue.Len() == 0 {
			p.pool.mu.Unlock()
			time.Sleep(time.Second)
			continue
		}
		weight := p.pool.queue.Top().weight
		if weight > int(p.pool.costSum.Load()) {
			p.pool.mu.Unlock()
			time.Sleep(time.Millisecond * 100)
			continue
		}
		p.pool.costSum.Add(int32(weight))
		t := p.pool.queue.Pop().(Task)
		p.status = WorkerStatusRunning
		p.pool.mu.Unlock()
		p.task = &t
		p.pool.sem <- struct{}{}
		t.setStatus(TaskStatusRunning)
		t.inner()
		p.pool.costSum.Add(-int32(weight))
		<-p.pool.sem
		p.pool.Done()
		if t.GetStatus() != TaskStatusCanceled {
			t.setStatus(TaskStatusCompleted)
		} else {
			p.status = WorkerStatusIdle
		}
		p.task = nil
	}
	p.pool.workersNum.Add(-1)
}

func (p *worker) Stop() {
	p.status = WorkerStatusStop
}

func (p *worker) IsRunning() bool {
	return p.status == WorkerStatusRunning
}

func (p *worker) IsIdle() bool {
	return p.status == WorkerStatusIdle
}

func (p *worker) IsStop() bool {
	return p.status == WorkerStatusStop
}

func (p *worker) GetID() string {
	return p.id
}

func (p *worker) GetPool() *TaskPool {
	return p.pool
}

func (p *worker) GetTaskInfo() any {
	if p.task.GetInfo() == nil {
		return nil
	}
	return *(p.task.GetInfo())
}
func (p *worker) UpdateTaskInfo(info *any) {
	p.task.info = info
}
func (p *worker) GetStatus() int {
	return p.status
}
