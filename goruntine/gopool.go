package goruntine

import "sync"

type GoPool struct {
	n  int
	c  chan func
	wg sync.WaitGroup
}
type Task func()
func NewGoPool(n int, worker int) *GoPool {
	p := &GoPool{
		n: n,
		c: make(chan Task, n),
	}
	p.wg.Add(worker)
	return p
}
func (g *GoPool) Add(task Task ) {
	g.c <- task
}
func (g *GoPool) Pop() {
	<-g.c
}
func (g *GoPool) Run() {
	go func(){
		defer g.Done()
		for task:=range g.c{
			task()
		}
	}()
}
func (g *GoPool) Done() {
	g.wg.Done()
}
func (g *GoPool) Wait() {
	g.wg.Wait()
}
func (g *GoPool)Close(){
	close(c)
}