package main

import (
	"fmt"
	"log"
	"net"
)

func RunServer2() {
	port := "8080"
	receiveNewMsg := make(chan string, 1)
	ln, err := net.Listen("tcp4", ":"+port)
	if err != nil {
		log.Printf("%s\n", err)
	}
	fmt.Printf("Listening on %s...\n", ln.Addr().String())
	for {
		ServerWG.Add(1)

		go handleRawConnection(ln, receiveNewMsg) // allow clients to say something
		fmt.Printf("a")
		fmt.Println(<-receiveNewMsg)
	}
}
