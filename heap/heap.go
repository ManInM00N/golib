package heap

import (
	"errors"
	"slices"
	"sync"
)

type PriorityQueue[T Interface] struct {
	items  []T               `json:"items"`
	lessF  func(a, b T) bool `json:"-"`
	equalF func(a, b T) bool `json:"-"`
	lk     sync.RWMutex      `json:"-"`
}

func NewPriorityQueueWithOptions[T Interface](opts ...Option[T]) *PriorityQueue[T] {
	pq := &PriorityQueue[T]{
		items: []T{},
	}
	for _, opt := range opts {
		opt(pq)
	}
	if pq.lessF == nil {
		panic(ErrRequiredLess)
	}
	pq.Init()
	return pq
}

func NewPriorityQueue[T Interface](lessF func(a, b T) bool) *PriorityQueue[T] {
	pq := &PriorityQueue[T]{
		items: []T{},
		lessF: lessF,
	}
	pq.Init()
	return pq
}

/*
use PriorityQueue:
q:=NewPriorityQueue(...)
q.Push({})
it:=q.Pop()
it:=q.Top()
*/
var (
	ErrEmpty        = errors.New("priority queue is empty")
	ErrOutOfIndex   = errors.New("Out of index")
	ErrRequiredLess = errors.New("less function is required")
)

func (pq *PriorityQueue[T]) Len() int {
	pq.lk.RLock()
	defer pq.lk.RUnlock()
	return len(pq.items)
}
func (pq *PriorityQueue[T]) less(i, j int) bool { return pq.lessF(pq.items[i], pq.items[j]) }
func (pq *PriorityQueue[T]) swap(i, j int)      { pq.items[i], pq.items[j] = pq.items[j], pq.items[i] }
func (pq *PriorityQueue[T]) Push(x Interface) {
	pq.lk.Lock()
	defer pq.lk.Unlock()
	pq.items = append(pq.items, x.(T))
	pq.up(x, len(pq.items)-1)
}
func (pq *PriorityQueue[T]) Pop() any {
	pq.lk.Lock()
	defer pq.lk.Unlock()
	old := pq.items
	n := len(old) - 1
	if n == -1 {
		panic(ErrEmpty)
	}
	pq.swap(0, n)
	pq.down(0, n)
	item := old[n]
	pq.items = old[0:n]
	return item
}
func (pq *PriorityQueue[T]) Top() T {
	pq.lk.Lock()
	defer pq.lk.Unlock()
	if len(pq.items) == 0 {
		panic(ErrEmpty)
	}
	return pq.items[0]
}

func (pq *PriorityQueue[T]) SetFunc(opt Option[T]) {
	pq.lk.Lock()
	defer pq.lk.Unlock()
	opt(pq)
}

type Interface interface {
}

// RemoveEqual 删除满足条件的相同的元素，返回删除的数量或-1（当没有设置相等判断函数时）
func (pq *PriorityQueue[T]) RemoveEqual(item T) int {
	pq.lk.Lock()
	defer pq.lk.Unlock()
	if pq.equalF == nil {
		return -1
	}
	num := len(pq.items)
	// 删除满足条件的任务
	pq.items = slices.DeleteFunc(pq.items, func(i T) bool {
		return pq.equalF(item, i)
	})
	return num - len(pq.items)
}

// RemoveEqual 删除满足条件的元素，返回删除的数量或-1（当没有设置相等判断函数时）
func (pq *PriorityQueue[T]) RemoveByDeleteFunc(deleteFunc func(item T) bool) int {
	pq.lk.Lock()
	defer pq.lk.Unlock()
	num := len(pq.items)
	// 删除满足条件的任务
	pq.items = slices.DeleteFunc(pq.items, func(i T) bool {
		return deleteFunc(i)
	})
	return num - len(pq.items)
}

func (pq *PriorityQueue[T]) Init() {
	// heapify
	pq.lk.RLock()
	defer pq.lk.RUnlock()
	n := len(pq.items)
	for i := n/2 - 1; i >= 0; i-- {
		pq.down(i, n)
	}
}
func (pq *PriorityQueue[T]) Remove(h Interface, i int) any {
	pq.lk.Lock()
	defer pq.lk.Unlock()
	n := pq.Len() - 1
	if n != i {
		pq.swap(i, n)
		if !pq.down(i, n) {
			pq.up(h, i)
		}
	}
	return pq.Pop()
}

func (pq *PriorityQueue[T]) Fix(h Interface, i int) {
	pq.lk.Lock()
	defer pq.lk.Unlock()
	if i >= len(pq.items) {
		panic(ErrOutOfIndex)
	}
	pq.items[i] = h.(T)
	if !pq.down(i, len(pq.items)) {
		pq.up(h, i)
	}
}

func (pq *PriorityQueue[T]) up(h Interface, j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !pq.less(j, i) {
			break
		}
		pq.swap(i, j)
		j = i
	}
}

func (pq *PriorityQueue[T]) down(i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && pq.less(j2, j1) {
			j = j2 // = 2*i + 2  // right child
		}
		if !pq.less(j, i) {
			break
		}
		pq.swap(i, j)
		i = j
	}
	return i > i0
}

func (pq *PriorityQueue[T]) Clear() {
	pq.lk.Lock()
	defer pq.lk.Unlock()
	pq.items = []T{}
}

func (pq *PriorityQueue[T]) Items() []T {
	pq.lk.RLock()
	defer pq.lk.RUnlock()
	tmp := slices.Clone(pq.items)
	return tmp
}
