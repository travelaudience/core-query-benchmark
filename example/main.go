package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/travelaudience/sqlbench"
)

func main() {
	config := flag.String("config", "example/config.json", "config file")
	flag.Parse()

	b, e := sqlbench.New(*config)
	if e != nil {
		log.Panicln(e)
	}
	wait := b.Start()
	<-wait

	fmt.Println("Finished")
}
