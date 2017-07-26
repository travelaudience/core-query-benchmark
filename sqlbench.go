/*
Package sqlbench can run a set of queries against a postgresql database and collect execution statistics.
*/
package sqlbench

import (
	"encoding/json"
	"io/ioutil"
)

// New will return a Bench structure that can be used to control the benchmark.
func New(configFile string) (Bench, error) {
	b := Bench{}
	var err error
	if b.config, err = config(configFile); err != nil {
		return b, err
	}

	b.runner = &sqlRunner{}
	return b, nil
}

// Start will started the benchmark and immediately returns.
// The returned channel can be used for waiting until benchmark is finish.
func (b *Bench) Start() chan bool {
	b.wait = make(chan bool)
	b.log.runs = make(map[string]Stats)

	go b.start()

	return b.wait
}

func config(fn string) (Config, error) {
	c := Config{}
	dat, err := ioutil.ReadFile(fn)

	if err != nil {
		return c, err
	}

	if err := json.Unmarshal(dat, &c); err != nil {
		return c, err
	}

	return c, nil
}
