package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

func NewClientConnection(n string) ClientConnection {
	ch := make(chan string)
	return ClientConnection{name: n, ch: ch}
}

func writeToClient(conn net.Conn, code int, msg string) {
	var respMsg Message
	var err error
	respMsg.Code = code
	respMsg.From = "server"
	respMsg.Message = "~_~ " + msg
	resp := bufio.NewWriter(conn)

	log.Printf("%v\n", respMsg)
	msgBytes, err := json.Marshal(respMsg)
	if err != nil {
		log.Printf("error marshaling json: %s", err)
	}
	log.Println(string(msgBytes))

	mb, merr := resp.Write(msgBytes)
	if merr != nil {
		log.Fatalln("error writing response to client")
	}
	resp.Flush()
	log.Printf("Wrote %d bytes to client\n", mb)
	conn.Close()
}

func splitSenderMessage(s string) (sender interface{}, message string) {
	split := strings.Split(s, " ")
	sender = split[0]
	message = strings.Join(split[1:], " ")
	message = strings.Trim(message, "\n")
	_sender, err := strconv.Atoi(sender.(string))
	if err == nil {
		return _sender, message
	} else {
		return sender, message
	}
}

func handleRawConnection(ln net.Listener, incomingMsg chan<- string) {
	conn, err := ln.Accept() //blocks until client dials
	// log.Printf("Connected to client: %s\n", conn.RemoteAddr().String())
	if err != nil {
		conn.Close()
		log.Println("??? " + conn.LocalAddr().String())
	}
	defer ServerWG.Done()

	reader := bufio.NewReader(conn)
	s, err := reader.ReadString('\n')
	if err == io.EOF {
		log.Printf("%s\n", err)
	} else if err != nil {
		log.Fatalln("Error reading from conn")
	}
	log.Printf("Read %d bytes: %s\n", len(s), s)
	sender, message := splitSenderMessage(s) // wait for any 1 client to say something

	/* decide how to handle this client's connection */
	switch sender.(type) {
	case string:
		// someone has chatted
		msgCount++
		fmt.Printf("User '%s' sent a message (%d/%d) '%s'\n", sender, msgCount, MAX_ROOM_MESSAGES, message)
		room.Post(Message{From: sender.(string), Message: message})
		writeToClient(conn, 200, "You posted a message!")
		break
	case int:
		// "control" message for server (user has joined = 42069)
		fmt.Printf("A new user '%s' joined (code %d)\n", message, sender)
		for i := 0; i < MAX_CLIENT_CONNECTIONS; i++ {
			if clients[i].name == "nobody" {
				clients[i] = NewClientConnection(message)
				fmt.Printf("Registered new user '%s' (%d/%d)\n", message, i+1, MAX_CLIENT_CONNECTIONS)
				writeToClient(conn, 200, "Gotcha fam! xD Have fun in chat!!!")
				break
			}
		}
		break
	}

	incomingMsg <- s // phone home to caller to inform of message
	log.Println("Done")
}

const (
	MAX_CLIENT_CONNECTIONS = 5
	MAX_ROOM_MESSAGES      = 512
)

var (
	ServerWG sync.WaitGroup
	room     Room
	board    [MAX_ROOM_MESSAGES]Message
	users    [MAX_CLIENT_CONNECTIONS]User
	clients  [MAX_CLIENT_CONNECTIONS]ClientConnection
	msgCount int
)

func RunServer() {
	f, err := os.OpenFile("logs/server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("Could not open logfile")
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(f)
	log.Println("Staring server.")

	//initialize clients
	for i := 0; i < MAX_CLIENT_CONNECTIONS; i++ {
		clients[i] = NewClientConnection("nobody")
	}

	var board []Message
	var gathering []User
	room = Room{name: "lolcats", users: gathering, board: board}

	msgCount = 0
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
