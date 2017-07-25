package main

import "github.com/travelaudience/sqlbench"
import "fmt"
import "log"

func main() {
	b, e := sqlbench.New("example/config.json")
	if e != nil {
		log.Fatal(e)
	}
	wait := b.Start()
	<-wait

	fmt.Println("Finished")
}
