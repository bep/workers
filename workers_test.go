package workers

import (
	"context"
	"runtime"
	"sort"
	"sync"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestPara(t *testing.T) {
	if runtime.NumCPU() < 2 {
		t.Skipf("skip para test, CPU count is %d", runtime.NumCPU())
	}

	c := qt.New(t)

	n := 1024
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
}
