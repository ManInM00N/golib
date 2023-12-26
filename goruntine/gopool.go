package goruntine

import "sync"

type GoPool struct {
	n  int
	c  chan struct{}
	wg sync.WaitGroup
}

func NewGoPool(n int, worker int) *GoPool {
	wg.Add(worker)
	return &GoPool{
		n: n,
		c: make(chan struct{}, n),
	}

}
func (g *GoPool) add() {
	g.c <- struct{}{}
}
func (g *GoPool) pop() {
	<-g.c
}
func (g *GoPool) Run(f func()) {
	g.add()
	defer g.wg.Done()
	f()
	g.pop()
}
func (g *GoPool) Wait() {
	g.wg.Wait()
}
