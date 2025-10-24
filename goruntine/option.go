package goruntine

import (
	. "github.com/ManInM00N/go-tool/heap"
)

// WithTaskPriority 设置任务优先级比较函数的便捷包装
// 使用示例：
//
//	pool := NewTaskPool(3, 5, WithTaskPriority(func(a, b Task) bool {
//	    return a.GetVal() > b.GetVal() // 值大的优先
//	}))
func WithTaskPriority(lessF func(a, b Task) bool) Option[Task] {
	return WithLessFunc(lessF)
}

// WithTaskEquality 设置任务相等判断函数的便捷包装
// 使用示例：
//
//	pool := NewTaskPool(3, 5, WithTaskEquality(func(a, b Task) bool {
//	    // 自定义相等逻辑
//	    return a.GetVal() == b.GetVal()
//	}))
func WithTaskEquality(equalF func(a, b Task) bool) Option[Task] {
	return WithEqualFunc(equalF)
}

// WithTaskEqualityByPriority 按 Priority 字段判断任务是否相等
func WithTaskEqualityByPriority() Option[Task] {
	return WithEqualFunc(func(a, b Task) bool {
		return a.GetPriority() == b.GetPriority()
	})
}

// WithTaskEqualityByInfo 根据任务的 info 字段进行相等判断(默认 info 需实现 comparable 接口)
func WithTaskEqualityByInfo[T comparable]() Option[Task] {
	return WithEqualFunc(func(a, b Task) bool {
		infoA, okA := (*a.GetInfo()).(T)
		infoB, okB := (*b.GetInfo()).(T)
		if !okA || !okB {
			return false
		}
		return infoA == infoB
	})
}

// WithTaskEqualityByInfoFunc 使用自定义函数比较 Info（更灵活）
// 使用示例：
//
//	pool := NewTaskPool(3, 5, WithTaskEqualityByInfoFunc(func(a, b interface{}) bool {
//	    infoA := a.(MyInfo)
//	    infoB := b.(MyInfo)
//	    return infoA.ID == infoB.ID
//	}))
func WithTaskEqualityByInfoFunc(equalF func(a, b interface{}) bool) Option[Task] {
	return WithEqualFunc[Task](func(a, b Task) bool {
		infoA := a.GetInfo()
		infoB := b.GetInfo()

		if infoA == nil || infoB == nil {
			return infoA == infoB
		}

		return equalF(*infoA, *infoB)
	})
}

// WithHighestPriorityFirst 值大的任务优先执行（降序）
func WithHighestPriorityFirst() Option[Task] {
	return WithLessFunc(func(a, b Task) bool {
		return a.GetPriority() > b.GetPriority()
	})
}

// WithLowestPriorityFirst 值小的任务优先执行（升序，默认）
func WithLowestPriorityFirst() Option[Task] {
	return WithLessFunc(func(a, b Task) bool {
		return a.GetPriority() < b.GetPriority()
	})
}

// WithFIFO 先进先出（需要使用时间戳作为 val）
func WithFIFO() Option[Task] {
	return WithLessFunc(func(a, b Task) bool {
		return a.inQueueTime.Before(b.inQueueTime)
	})
}

// WithLIFO 后进先出（需要使用时间戳作为 val）
func WithLIFO() Option[Task] {
	return WithLessFunc(func(a, b Task) bool {
		return a.inQueueTime.After(b.inQueueTime)
	})
}
