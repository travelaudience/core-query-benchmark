package main

import (
	"fmt"
	"log"

	"github.com/travelaudience/core-sqlbench/sqlbench"
)

func main() {

	b, e := sqlbench.New("example/config.json")
	if e != nil {
		log.Fatal(e)
	}
	wait := b.Start()
	<-wait

	fmt.Println("Finished")
}
