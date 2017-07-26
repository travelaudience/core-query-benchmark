package sqlbench

import (
	"fmt"
	"sync"
	"time"
)

func (b *Bench) start() {
	b.xTags()

	all := sync.WaitGroup{}
	for _, q := range b.config.Queries {
		all.Add(1)
		go func() {
			b.benchmarkQuery(q)
			all.Done()
		}()
	}
	all.Wait()

	// signal the end
	go func() {
		time.Sleep(time.Second)
		b.wait <- true
	}()
}

func (b *Bench) xTags() {
	var tags []Tag
	for i, t := range b.config.Tags {
		switch {
		case t.Value == "timestamp":
			tags = append(tags, Tag{b.config.Tags[i].Name, time.Now().Format("2006-01-02 15:04:05")})
		default:
			tags = append(tags, Tag{b.config.Tags[i].Name, "test"})
		}
	}
	b.tags = tags
}

func (b *Bench) benchmarkQuery(q Query) {
	var report []float64
	runTime := make(chan float64)
	done := make(chan bool)

	var count int
	go func() {
		ticker := time.NewTicker(time.Millisecond * time.Duration(q.Frequency))
		for {
			println("next")
			all := sync.WaitGroup{}
			for i := 0; i < q.Parallel; i++ {
				all.Add(1)
				go func() {
					t := time.Now().UnixNano()
					b.run(b.config.Db, q.Query)
					runTime <- float64(time.Now().UnixNano()-t) / 1000000
					all.Done()
				}()
			}
			all.Wait()
			count++
			if count >= q.Count {
				ticker.Stop()
				done <- true
				return
			}
			<-ticker.C
		}
	}()

	for {
		select {
		case r := <-runTime:
			report = append(report, r)
			fmt.Println(r)
		case <-done:
			fmt.Println(report)
			fmt.Println("save the results")
			return
		}
	}
}
