package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/Zeebrow/whyrc/shared"
)

var ClientWG sync.WaitGroup
var USERNAME string

// func letsRead(conn net.Conn, readerChan chan<- string) {
func getServerResponse(conn net.Conn, readerChan chan<- shared.Message) {
	defer conn.Close() //Close connection as soon as server responds
	var err error = nil
	var receivedMessage shared.Message
	connReader := bufio.NewReader(conn)
	var recBytes []byte
	recBytes, err = connReader.ReadBytes('\n')
	if err != io.EOF {
		log.Printf("---> error reading: %s", err)
	}
	log.Printf("Read %d bytes to unmarshal.\n", len(recBytes))
	err = json.Unmarshal(recBytes, &receivedMessage)

	readerChan <- receivedMessage
	conn.Close()

	if err == io.EOF {
		log.Println("EOF reached")
	} else if err != nil {
		log.Printf("Error reading continuously: %s\n", err)
		conn.Close()
	}
}

func sendMessage(conn net.Conn, serverResponse chan<- shared.Message, name string, msg string) {
	var err error = nil
	readCh := make(chan shared.Message, 1)
	go getServerResponse(conn, readCh) //receives server response and writes it to channel

	writer := bufio.NewWriter(conn)
	_, err = writer.WriteString(name + " " + msg) //dat delimiter, doe
	if err != nil {
		log.Printf("failed to write to client: %s", err)
	}
	writer.Flush()
	writeout := <-readCh
	serverResponse <- writeout
	// log.Println("Closing connection")
}

func connect(name string, msg string, respChan chan shared.Message) {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalln("Error connecting to server")
	}

	ClientWG.Add(1)
	go sendMessage(conn, respChan, name, msg)
	log.Printf("received response from server: %v\n", <-respChan)
}

func srvrMsg(name string, msg string, respChan chan shared.Message) {
	connect(name, msg+"\n", respChan) //Only open a new connection when we want to say something
}

// func join(clientMsgChan chan string) func(ch chan string) {
func join(clientMsgChan chan shared.Message, serverMsgChan chan shared.Message) func() {
	var n string = USERNAME
	fmt.Println("Hi! You found the chat!")
	fmt.Printf("What's your name?>")
	name := bufio.NewReader(os.Stdin)
	n, _ = name.ReadString('\n')
	n = strings.ReplaceAll(n, "\n", "")
	fmt.Printf("Welcome, %s! Type your message and press 'Enter' to send. Ctrl+C exits.\n\n", n)

	srvrMsg("42069", n, serverMsgChan)
	// log.Printf("server message response to joining: %v\n", <-serverMsgChan)
	return func() {
		log.Println()
		log.Println("Awaiting client message input")
		fmt.Printf("%s>", n)
		reader := bufio.NewReader(os.Stdin)
		rb, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("error reading message from stdin: %s\n", err)
		}
		connect(n, rb, clientMsgChan) //Only open a new connection when we want to say something
		clientMsgChan <- shared.Message{Code: 200, From: n, Message: "HAI"}
	}
}

func RunClient() {
	f, err := os.OpenFile("logs/client.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("Could not open logfile")
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(f)
	log.Println()
	log.Println("Starting client")

	clientSendMsgChan := make(chan shared.Message, 1)
	clientReceiveSrvrMsgChan := make(chan shared.Message, 1)

	yMsg := join(clientSendMsgChan, clientReceiveSrvrMsgChan)
	for {
		go yMsg() // disconnect to reset name
		select {
		case c := <-clientSendMsgChan:
			log.Printf("Client message chan %v\n", c)
		case s := <-clientReceiveSrvrMsgChan:
			log.Printf("Got message from server: %v\n", s)
		}
	}
}
