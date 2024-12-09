package heap

import (
	"container/heap"
	"sync"
)

type PriorityQueue[T any] struct {
	items []T
	lessF func(a, b T) bool
	wg    sync.Pool
}

func NewPriorityQueue[T any](lessF func(a, b T) bool) *PriorityQueue[T] {
	pq := &PriorityQueue[T]{
		items: []T{},
		lessF: lessF,
	}
	heap.Init(pq)
	return pq
}

/*
use PriorityQueue:
q:=NewPriro...
heap.Push(q,{})
it:=heap.Pop(q)
*/
func (pq PriorityQueue[T]) Len() int           { return len(pq.items) }
func (pq PriorityQueue[T]) Less(i, j int) bool { return pq.lessF(pq.items[i], pq.items[j]) }
func (pq PriorityQueue[T]) Swap(i, j int)      { pq.items[i], pq.items[j] = pq.items[j], pq.items[i] }
func (pq *PriorityQueue[T]) Push(x any)        { pq.items = append(pq.items, x.(T)) }
func (pq *PriorityQueue[T]) Pop() any {
	old := pq.items
	n := len(old)
	if n == 0 {
		return nil
	}
	item := old[n-1]
	pq.items = old[0 : n-1]
	return item
}
