package goruntine

import "sync"

type GoPool struct {
	n  int
	c  chan struct{}
	wg sync.WaitGroup
}

func NewGoPool(n int, worker int) *GoPool {
	p := &GoPool{
		n: n,
		c: make(chan struct{}, n),
	}
	p.wg.Add(worker)
	return p
}
func (g *GoPool) add() {
	g.c <- struct{}{}
}
func (g *GoPool) pop() {
	<-g.c
}
func (g *GoPool) Run(f func()) {
	g.add()
	f()
	g.pop()
}
func (g *GoPool) Done() {
	g.wg.Done()
}
func (g *GoPool) Wait() {
	g.wg.Wait()
}
