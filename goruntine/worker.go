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
		if p.pool.queue.Len() == 0 && p.pool.running {
			p.status = 0
			p.pool.mu.Unlock()
			time.Sleep(time.Second)
			continue
		}
		if !p.pool.running && p.pool.queue.Len() == 0 {
			p.status = -1
			p.pool.mu.Unlock()
			break
		}
		p.status = 1
		t := p.pool.queue.Pop().(task)
		p.pool.mu.Unlock()

		p.pool.sem <- struct{}{}
		go func(t task) {
			defer func() {
				<-p.pool.sem
				p.pool.Done()
			}()
			t.Inner()
		}(t)
	}
	p.pool.workersNum.Add(-1)
}
