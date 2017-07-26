package main

import (
	"fmt"
	"log"

	"github.com/travelaudience/sqlbench"
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
