/*
Package sqlbench can run a set of queries against a postgresql database and collect execution statistics.
*/
package sqlbench

// Config specifies the setup for benchmarks
type Config struct {
	Xtags []struct {
		// Name can be any string
		Name string `json:"name"`
		// Value can be only a query which will result in a numeric value. Alternatively it can be `epoch` or `datetime`.
		Value string `json:"value"`
	} `json:"xtags"`
	Queries []struct {
		// A name for the query
		Name string `json:"0"`
		// Running frequency in millisecond
		Frequency int `json:"1"`
		// Number of parallel runs for this query
		Parallel int `json:"2"`
		// Query to run
		Query string `json:"3"`
	} `json:"queries"`
	// PostgreSQL database DSN
	Db   string `json:"db"`
	Logs struct {
		// If set it will append the logs to this csv file.
		Csv string `json:"csv"`
		// If set it will send the results to datadog
		Datadog string `json:"datadog"`
	} `json:"logs"`
}

func pick() {

}
