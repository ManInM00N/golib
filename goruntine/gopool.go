package goruntine

type GoPool struct {
	n int
	c chan any
}

func NewGoPool(n int) *GoPool {
	return &GoPool{
		n: n,
		c: make(chan any, n),
	}
}
func (g *GoPool) add() {
	g.c <- nil
}
func (g *GoPool) pop() {
	<-g.c
}
func (g *GoPool) Run(f func()) {
	g.add()
	f()
	g.pop()
}
