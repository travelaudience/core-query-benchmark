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

	b.runner = &sqlRunner{dsn: b.config.Db}
	b.runner.init()
	return b, nil
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
