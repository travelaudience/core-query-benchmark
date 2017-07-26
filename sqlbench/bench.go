package sqlbench

import (
	"fmt"
	"log"
	"math"
	"sync"
	"time"
)

// Start will started the benchmark and immediately returns.
// The returned channel can be used for waiting until benchmark is finish.
func (b *Bench) Start() chan bool {
	b.wait = make(chan bool)
	b.runLog.runs = make(map[string]Stats)

	go b.start()

	return b.wait
}

func (b *Bench) start() {
	defer func() { b.wait <- true }()
	b.tag()

	all := sync.WaitGroup{}
	for _, q := range b.config.Queries {
		all.Add(1)
		go func() {
			b.benchmarkQuery(q)
			all.Done()
		}()
	}
	all.Wait()

	b.save()
}

func (b *Bench) tag() {
	var tags []Tag
	for i, t := range b.config.Tags {
		switch {
		case t.Value == "timestamp":
			tags = append(tags, Tag{b.config.Tags[i].Name, time.Now().Format("2006-01-02 15:04:05")})
		default:
			value, err := b.runner.tag(b.config.Tags[i].Value)
			if err != nil {
				log.Println(err)
				return
			}
			tags = append(tags, Tag{b.config.Tags[i].Name, value})
		}
		fmt.Println("Tag:", t.Name, ":", tags[i].Value)
	}
	b.runLog.tags = tags
}

func (b *Bench) benchmarkQuery(q Query) {
	fmt.Println("running", q.Name)
	var report []float64
	runTime := make(chan float64)
	done := make(chan bool)

	var count int
	go func() {
		ticker := time.NewTicker(time.Millisecond * time.Duration(q.Frequency))
		for {
			all := sync.WaitGroup{}
			for i := 0; i < q.Parallel; i++ {
				all.Add(1)
				go func() {
					t := time.Now().UnixNano()
					b.run(q.Query)
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

	var min, max, sum, n float64
	for {
		select {
		case r := <-runTime:
			report = append(report, r)
			switch {
			case r < min || n == 0:
				min = r
			case r > max:
				max = r
			}
			sum += r
			n++
		case <-done:
			b.Lock()
			b.runLog.runs[q.Name] = Stats{min, sum / n, max, std(report, sum/n, n)}
			b.Unlock()
			fmt.Println("done", q.Name)
			return
		}
	}
}

func std(r []float64, avg float64, n float64) float64 {
	var sum float64
	for _, f := range r {
		sum += (f - avg) * (f - avg)
	}
	sum = sum / n
	return math.Sqrt(sum)
}

func (b *Bench) save() {
	fmt.Println(b.runLog)
	if b.config.Logs.Csv != "" {
		fmt.Println("saving into", b.config.Logs.Csv)
	}
	if b.config.Logs.Datadog != "" {
		fmt.Println("send data to", b.config.Logs.Datadog)
	}
}
