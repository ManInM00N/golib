package Map

import (
	"strconv"
	"sync"
	"testing"
)

func TestChannelMap(t *testing.T) {
	mm := NewChannelMap[string, int]()
	var wg sync.WaitGroup
	wg.Add(100)
	for i := range 100 {
		if i&1 == 1 {
			go func() {
				mm.Set(strconv.FormatInt(int64(i/2), 10), i)
				wg.Done()
			}()
		} else {
			go func() {
				mm.Get(strconv.FormatInt(int64(i/2), 10))
				wg.Done()
			}()
		}
	}
	wg.Wait()
	if len(mm.mp) != 50 {
		t.Error("Add error")
	}
	wg.Add(50)
	for i := range 50 {
		go func() {
			mm.Delete(strconv.FormatInt(int64(i), 10))
			wg.Done()
		}()
	}
	wg.Wait()
	if len(mm.mp) != 0 {
		t.Error("Delete error")
	}
}

func TestRWMap(t *testing.T) {
	mm := NewRWMap[string, int]()
	var wg sync.WaitGroup
	wg.Add(100)
	for i := range 100 {
		if i&1 == 1 {
			go func() {
				mm.Set(strconv.FormatInt(int64(i/2), 10), i)
				wg.Done()
			}()
		} else {
			go func() {
				mm.Get(strconv.FormatInt(int64(i/2), 10))
				wg.Done()
			}()
		}
	}
	wg.Wait()
	if len(mm.mp) != 50 {
		t.Error("Add error")
	}
	wg.Add(50)
	for i := range 50 {
		go func() {
			mm.Delete(strconv.FormatInt(int64(i), 10))
			wg.Done()
		}()
	}
	wg.Wait()
	if len(mm.mp) != 0 {
		t.Error("Delete error")
	}
}

func BenchmarkRWMap(b *testing.B) {
	b.StopTimer()
	mm := NewRWMap[int, int]()
	var wg sync.WaitGroup
	wg.Add(b.N)
	b.StartTimer()
	for i := range b.N {
		if i&1 == 1 {
			go func() {
				mm.Set(i/2, i)
				wg.Done()
			}()
		} else {
			go func() {
				mm.Get(i / 2)
				wg.Done()
			}()
		}
	}
	wg.Wait()
}

func BenchmarkChannelMap(b *testing.B) {
	b.StopTimer()
	mm := NewChannelMap[int, int]()
	var wg sync.WaitGroup
	wg.Add(b.N)
	b.StartTimer()
	for i := range b.N {
		if i&1 == 1 {
			go func() {
				mm.Set(i/2, i)
				wg.Done()
			}()
		} else {
			go func() {
				mm.Get(i / 2)
				wg.Done()
			}()
		}
	}
	wg.Wait()
}
