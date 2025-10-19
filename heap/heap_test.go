package heap

import (
	"os"
	"testing"
)

type ab struct {
	x, y int
}

func (xx *ab) Less(yy ab) bool {
	if xx.x == yy.x {
		return xx.y > yy.y
	}
	return xx.x < yy.x
}

func TestPriorityQueue(t *testing.T) {
	q := NewPriorityQueueWithOptions(
		WithLessFunc(func(x, y *ab) bool {
			return x.Less(*y)
		}),
		WithEqualFunc(func(x, y *ab) bool {
			return x.x == y.x && x.y == y.y
		}))
	q.Push(&ab{1, 2})
	q.Push(&ab{3, 2})
	q.Push(&ab{2, 2})
	q.Push(&ab{2, 2})
	q.Push(&ab{5, 2})
	q.Push(&ab{2, 4})
	q.Fix(&ab{1, 99}, 4)

	for q.Len() > 0 {
		tmp := q.Top()
		t.Error(tmp.x, tmp.y, *tmp == ab{1, 99})

		if tmp.y == 99 {
			q.RemoveByDeleteFunc(func(item *ab) bool {
				return item.x == 2 && item.y == 2
			})
			q.RemoveEqual(&ab{2, 2})
		}
		q.Pop()
	}
	/* expect
	1 99
	1 2
	2 4
	3 2

	*/
}

func TestMain(m *testing.M) {
	code := m.Run()

	// 退出测试
	os.Exit(code)
}
