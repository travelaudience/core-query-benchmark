package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/travelaudience/core-sqlbench/sqlbench"
)

func main() {
	config := flag.String("config", "example/config.json", "config file")
	flag.Parse()

	b, e := sqlbench.New(*config)
	if e != nil {
		log.Fatal(e)
	}
	wait := b.Start()
	<-wait

	fmt.Println("Finished")
}
