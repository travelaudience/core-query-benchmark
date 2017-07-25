package sqlbench

import (
	"time"
)

// Stats is collections data we collect for each query run
type Stats struct {
	// Minimum runtime
	Min, Max, Avg, Stdv float64
	XTags               []struct {
		Name  string
		Value string
	}
}

// Log of total exection and also queries benchmarks
type Log struct {
	Query map[string]Stats
	Runs  map[string]([]Stats)
}

func (b *Bench) start() {

	// signal the end
	go func() {
		time.Sleep(time.Second)
		b.wait <- true
	}()
}
