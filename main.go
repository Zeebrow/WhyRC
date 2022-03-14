package main

import (
	"flag"
	"fmt"
)

func main() {

	mode := flag.String("m", "client", "Run in either 'client' or 'server' mode")
	flag.Parse()
	if *mode == "client" {
		RunClient()
	} else if *mode == "server" {
		RunServer()
	} else {
		fmt.Printf("Invalid mode option provided: '%s'\n", *mode)
	}

}
