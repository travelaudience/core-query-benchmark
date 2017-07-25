package main

import "github.com/travelaudience/sqlbench"
import "fmt"
import "log"

func main() {
	b, e := sqlbench.New("example/config.json")
	if e != nil {
		log.Fatal(e)
	}
	w := b.Start()
	<-w
	fmt.Println("Finished")
}
