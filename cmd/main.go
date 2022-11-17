package main

import (
	"flag"
	"fmt"

	"github.com/Zeebrow/whyrc/client"
	"github.com/Zeebrow/whyrc/server"
)

func main() {

	mode := flag.String("m", "client", "Run in either 'client' or 'server' mode")
	flag.Parse()
	if *mode == "client" {
		client.RunClient()
	} else if *mode == "server" {
		server.RunServer()
	} else {
		fmt.Printf("Invalid mode option provided: '%s'\n", *mode)
	}

}
