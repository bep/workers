[![Tests](https://github.com/bep/workers/workflows/Test/badge.svg)](https://github.com/bep/workers/actions?query=workflow%3ATest)
[![Go Report Card](https://goreportcard.com/badge/github.com/bep/workers)](https://goreportcard.com/report/github.com/bep/workers)
[![GoDoc](https://godoc.org/github.com/bep/workers?status.svg)](https://godoc.org/github.com/bep/workers)


A simple Go library to set up tasks to be executed in parallel.

```go
package main

import (
	"context"
	"log"

	"github.com/bep/workers"
)

func main() {
	// Max 4 tasks to be executed in parallel.
	w := workers.New(4)
	r, _ := w.Start(context.Background())

	r.Run(func() error {
		return nil
	})

	// ... run more tasks.

	if err := r.Wait(); err != nil {
		log.Fatal(err)
	}
}
```