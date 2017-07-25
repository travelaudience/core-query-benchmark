package sqlbench

import (
	"fmt"
	"time"
)

func (b *Bench) start() {

	for _, q := range b.config.Queries {
		b.benchmarkQuery(q)
	}

	// signal the end
	go func() {
		time.Sleep(time.Second)
		b.wait <- true
	}()
}

func (b *Bench) benchmarkQuery(q Query) {
	reports := make(chan float64)
	done := make(chan bool)
	var count int

	go func() {
		ticker := time.NewTicker(time.Millisecond * time.Duration(q.Frequency))
		for range ticker.C {
			for i := 0; i < q.Parallel; i++ {
				go func() {
					t := time.Now().UnixNano()
					b.run(b.config.Db, q.Query)
					reports <- float64(time.Now().UnixNano()-t) / 1000000
				}()
			}
			count++
			if count >= q.Count {
				ticker.Stop()
				done <- true
				return
			}
		}
	}()

	<-done

	fmt.Println("Ticker stopped")
}
