package sqlbench

import (
	"io/ioutil"
	"testing"
	"time"
)

const configTest = `
{
  "tags":[
    {"name":"time","value": "timestamp"} ,
    {"name":"rows", "value": "SELECT count(*) FROM  mytable"}
  ],

  "queries": [
     {"name":"bidding01", "frequency":1, "parallel":2, "count": 3, "query": "SELECT * from mytable m join another_table a on (m.id=a.mytable_id) "}
  ],
  "logs": {
  }
}
`

func getBenchTest(t *testing.T) *Bench {
	b := Bench{runner: &sqlRunnerTest{}}
	var e error
	if b.config, e = config(createJSONFile(configTest)); e != nil {
		t.Error("Expected a new object with no error", e)
		return nil
	}
	return &b
}

func TestBench_Start(t *testing.T) {
	b := getBenchTest(t)

	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Error("can't create the temp file")
	}
	b.config.Logs.Csv = tmpfile.Name()

	wait := b.Start()
	select {
	case <-wait:
	case <-time.After(time.Second):
		t.Error("Benchmark did not finish")
	}

	fStat, _ := tmpfile.Stat()
	if fStat.Size() == 0 {
		t.Error("No data in log file")
	}
}
