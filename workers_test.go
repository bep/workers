package workers

import (
	"context"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
)

func TestPara(t *testing.T) {
	if runtime.NumCPU() < 2 {
		t.Skipf("skip para test, CPU count is %d", runtime.NumCPU())
	}

	c := qt.New(t)

	c.Run("Order", func(c *qt.C) {
		n := 500
		ints := make([]int, n)
		for i := 0; i < n; i++ {
			ints[i] = i
		}

		p := New(4)
		r, _ := p.Start(context.Background())

		var result []int
		var mu sync.Mutex
		for i := 0; i < n; i++ {
			i := i
			r.Run(func() error {
				mu.Lock()
				defer mu.Unlock()
				result = append(result, i)
				return nil
			})
		}

		c.Assert(r.Wait(), qt.IsNil)
		c.Assert(result, qt.HasLen, len(ints))
		c.Assert(sort.IntsAreSorted(result), qt.Equals, false, qt.Commentf("Para does not seem to be parallel"))
		sort.Ints(result)
		c.Assert(result, qt.DeepEquals, ints)
	})

	c.Run("Time", func(c *qt.C) {
		const n = 100

		p := New(5)
		r, _ := p.Start(context.Background())

		start := time.Now()

		var counter int64

		for i := 0; i < n; i++ {
			r.Run(func() error {
				atomic.AddInt64(&counter, 1)
				time.Sleep(1 * time.Millisecond)
				return nil
			})
		}

		c.Assert(r.Wait(), qt.IsNil)
		c.Assert(counter, qt.Equals, int64(n))

		since := time.Since(start)
		limit := n / 2 * time.Millisecond
		c.Assert(since < limit, qt.Equals, true, qt.Commentf("%s >= %s", since, limit))
	})
}
