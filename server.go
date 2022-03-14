package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

type roomMsg struct {
	OP   string
	Text string
}
type messageSender interface {
	Send()
}

func (m *roomMsg) Send() {

}

func response(conn net.Conn) {
	resp := bufio.NewWriter(conn)
	mb, merr := resp.WriteString("200")
	if merr != nil {
		log.Fatal("error writing response to client")
	}
	log.Printf("Wrote %d bytes\n", mb)
	resp.Flush()
	log.Println("Response sent.")
	conn.Close()
}

func sendMessage(conn net.Conn, sender string, text string, ch chan string) {
	resp := bufio.NewWriter(conn)
	mb, merr := resp.WriteString(sender + " " + text)
	if merr != nil {
		log.Fatal("error writing response to client")
	}
	log.Printf("Wrote %d bytes\n", mb)
	resp.Flush()
	log.Println("Response sent.")
	conn.Close()
}

func splitSenderMessage(s string) (sender string, message string) {
	split := strings.Split(s, " ")
	sender = split[0]
	message = strings.Join(split[1:], " ")
	return
}

// func handleConn(conn net.Conn, ch chan string) {
func handleConn(ln net.Listener, incomingMsg chan<- string) {
	conn, err := ln.Accept()
	log.Printf("Connected to client: %s\n", conn.RemoteAddr().String())
	if err != nil {
		conn.Close()
		log.Fatal("Failed to bind to port " + conn.LocalAddr().String())
	}
	// defer conn.Close()
	defer ServerWG.Done()

	reader := bufio.NewReader(conn)
	s, err := reader.ReadString('\n')
	if err == io.EOF {
		log.Printf("%s\n", err)
	} else if err != nil {
		log.Fatal("Error reading from conn")
	}
	log.Printf("Read %d bytes: %s\n", len(s), s)
	incomingMsg <- s
	response(conn)
	log.Println("Done")
}

var ServerWG sync.WaitGroup

var roomText [512]roomMsg

func RunServer() {
	defer initLog().Close()

	msgCount := 0
	port := "8080"
	receiveNewMsg := make(chan string, 1)
	ln, err := net.Listen("tcp4", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Listening on %s...\n", ln.Addr().String())
	for {
		ServerWG.Add(1)
		go handleConn(ln, receiveNewMsg)
		sender, message := splitSenderMessage(<-receiveNewMsg) // blocks here
		log.Printf("Message from '%s': %s\n", sender, message)
		roomText[0] = roomMsg{OP: sender, Text: message}
		msgCount++
	}
}
