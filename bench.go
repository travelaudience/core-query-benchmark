package sqlbench

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
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

	var m sync.Map
	done := make(chan bool)
	var count int32
	for _, query := range b.config.Queries {
		m.Store(query.Name, 1)
		atomic.AddInt32(&count, 1)
		go func(q Query) {
			b.benchmarkQuery(q)
			m.Store(q.Name, 0)
			atomic.AddInt32(&count, -1)
			if atomic.LoadInt32(&count) == 0 {
				done <- true
			}
		}(query)
	}

	for c := true; c == true; {
		select {
		case <-time.After(time.Second * 5):
			var l []string
			var s string
			m.Range(func(k, v interface{}) bool {
				if v.(int) == 1 {
					l = append(l, k.(string))
				}
				return true
			})
			sort.Strings(l)
			for _, v := range l {
				s += v + ", "
			}

			log.Println("still running: ", s)
		case <-done:
			c = false
		}
	}

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
				log.Panicln(err)
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
					err := b.run(q.Query)
					if err != nil {
						log.Panicln("error while running", q.Name, err)
					}
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
			b.runLog.runs[q.Name] = Stats{min, sum / n, max, std(report, sum/n, n), pct(report, 95)}
			b.Unlock()
			fmt.Println("done", q.Name)
			return
		}
	}
}

func pct(r []float64, pct float64) float64 {
	sort.Float64s(r)
	k := math.Ceil(float64(len(r)) * pct / 100)
	return r[int(k-1)]
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
	if b.config.Logs.Csv != "" {
		fmt.Println("into", b.config.Logs.Csv)
		f, err := os.OpenFile(b.config.Logs.Csv, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			log.Panicln(err)
		}
		defer f.Close()

		var header []string
		var tags []string
		for _, t := range b.tags {
			tags = append(tags, t.Value)
			header = append(header, t.Name)
		}
		header = append(header, []string{"query", "Min", "Avg", "Max", "Stdv", "Pct95"}...)

		w := csv.NewWriter(f)
		stat, _ := f.Stat()
		if stat.Size() == 0 {
			w.Write(header)
		}

		for k, v := range b.runLog.runs {
			var line []string
			line = append(line, tags...)
			line = append(line, k)
			line = append(line, strconv.FormatFloat(v.Min, 'f', 0, 64))
			line = append(line, strconv.FormatFloat(v.Avg, 'f', 0, 64))
			line = append(line, strconv.FormatFloat(v.Max, 'f', 0, 64))
			line = append(line, strconv.FormatFloat(v.Stdv, 'f', 0, 64))
			line = append(line, strconv.FormatFloat(v.Pct95, 'f', 0, 64))
			w.Write(line)
		}

		w.Flush()
	}
}
