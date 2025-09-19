package goruntine

import "time"

type worker struct {
	id     string
	pool   *TaskPool
	task   *task
	status int //0空闲 1工作 -1停止
}

func (p *worker) Run() {
	p.pool.workersNum.Add(1)
	for {
		p.pool.mu.Lock()

		if !p.pool.running && p.pool.queue.Len() == 0 {
			p.status = -1
			p.pool.mu.Unlock()
			break
		}
		if p.pool.queue.Len() == 0 {
			p.pool.mu.Unlock()
			time.Sleep(time.Second)
			continue
		}
		p.status = 1
		t := p.pool.queue.Pop().(task)
		p.pool.mu.Unlock()
		p.task = &t
		p.pool.sem <- struct{}{}
		t.SetStatus(1)
		defer func() {
			<-p.pool.sem
			p.pool.Done()
		}()
		t.Inner()
		if t.GetStatus() != -1 {
			t.SetStatus(2)
		} else {
			p.status = 0
		}
		p.task = nil
		// go func(t task) {
		// 	t.SetStatus(1)
		// 	defer func() {
		// 		<-p.pool.sem
		// 		p.pool.Done()
		// 	}()
		// 	t.Inner()
		// 	if t.GetStatus() != -1 {
		// 		t.SetStatus(2)
		// 	}
		// }(t)
	}
	p.pool.workersNum.Add(-1)
}

func (p *worker) Stop() {
	p.status = -1
}

func (p *worker) IsRunning() bool {
	return p.status == 1
}

func (p *worker) IsIdle() bool {
	return p.status == 0
}

func (p *worker) IsStop() bool {
	return p.status == -1
}

func (p *worker) GetID() string {
	return p.id
}

func (p *worker) GetPool() *TaskPool {
	return p.pool
}

func (p *worker) GetTaskInfo() any {
	if p.task.Info == nil {
		return nil
	}
	return *(p.task.Info)
}
func (p *worker) UpdateTaskInfo(info *any) {
	p.task.Info = info
}
func (p *worker) GetStatus() int {
	return p.status
}
