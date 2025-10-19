package goruntine

import "sync"

type SingleTask func()
type GoPool struct {
	n  int
	c  chan SingleTask
	wg sync.WaitGroup
}

func NewGoPool(n int, worker int) *GoPool {
	p := &GoPool{
		n: n,
		c: make(chan SingleTask, n),
	}
	p.AddWorker(worker)
	return p
}
func (g *GoPool) AddWorker(num int) {
	g.wg.Add(num)
}
func (g *GoPool) Add(SingleTask SingleTask) {
	g.c <- SingleTask
}
func (g *GoPool) Pop() {
	<-g.c
}
func (g *GoPool) Run() {
	go func() {
		defer g.Done()
		for SingleTask := range g.c {
			SingleTask()
		}
	}()
}
func (g *GoPool) Done() {
	g.wg.Done()
}
func (g *GoPool) Wait() {
	g.wg.Wait()
}
func (g *GoPool) Close() {
	close(g.c)
}
