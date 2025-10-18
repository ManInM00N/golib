package heap

type Option[T Interface] func(pq *PriorityQueue[T])

// WithLessFunc 设置比较函数
func WithLessFunc[T Interface](lessF func(a, b T) bool) Option[T] {
	return func(pq *PriorityQueue[T]) {
		pq.lessF = lessF
	}
}

// WithEqualFunc 设置相等判断函数
func WithEqualFunc[T Interface](equalF func(a, b T) bool) Option[T] {
	return func(pq *PriorityQueue[T]) {
		pq.equalF = equalF
	}
}

// WithInitialItems 设置初始元素
func WithInitialItems[T Interface](items []T) Option[T] {
	return func(pq *PriorityQueue[T]) {
		pq.items = items
	}
}
