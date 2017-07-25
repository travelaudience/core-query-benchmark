package sqlbench

import (
	"encoding/json"
	"io/ioutil"
)

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
