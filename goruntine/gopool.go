package goruntine

type GoPool struct {
	n int
	c chan struct{}
}

func NewGoPool(n int) *GoPool {
	return &GoPool{
		n: n,
		c: make(chan struct{}, n),
	}
}
func (g *GoPool) add() {
	g.c <- struct{}
}
func (g *GoPool) pop() {
	<-g.c
}
func (g *GoPool) Run(f func()) {
	g.add()
	f()
	g.pop()
}
