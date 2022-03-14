package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

var ClientWG sync.WaitGroup
var USERNAME string = ""

func letsRead(conn net.Conn, readerChan chan<- string) {
	defer conn.Close() //Close connection as soon as server responds
	var err error = nil
	connReader, err := bufio.NewReader(conn).ReadString('\n')
	if err == io.EOF {
		readerChan <- connReader
	} else if err != nil {
		log.Printf("Error reading continuously: %s\n", err)
		conn.Close()
	}
}

func sendText(conn net.Conn, writerChan chan<- string, name string, msg string) {
	var err error = nil
	readCh := make(chan string, 1)
	go letsRead(conn, readCh) //receives server response and writes it to channel

	writer := bufio.NewWriter(conn)
	_, err = writer.WriteString(name + " " + msg) //dat delimiter, doe
	if err != nil {
		log.Printf("failed to write to client: %s", err)
	}
	writer.Flush()
	log.Println("Flushed")

	writeout := <-readCh
	writerChan <- writeout
	log.Println("Closing connection")
}

func startClient(conn net.Conn, responseChan chan<- string, readerChan chan string, name string, msg string) {
	log.Printf("started client at %s\n", conn.RemoteAddr().String())
	ch := make(chan string, 1)

	ClientWG.Add(1)
	go sendText(conn, ch, name, msg)
	resp := <-ch
	readerChan <- resp
}

type Message struct {
	Name string
	Msg  string
}

func connect(name string, msg string) {
	log.Println("Creating new connection")
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Error connecting to server")
	}
	log.Printf("Connected to server %s\n", conn.RemoteAddr().String())

	readerChan := make(chan string, 1)
	writerChan := make(chan string, 1)

	ClientWG.Add(1)
	go startClient(conn, readerChan, writerChan, name, msg)
	log.Printf("chan writer: %s\n", <-writerChan)
}
func srvrMsg(name string, msg string) {
	connect(name, msg+"\n") //Only open a new connection when we want to say something
}

func msg(name string) {
	fmt.Printf("%s>", name)
	reader := bufio.NewReader(os.Stdin)
	rb, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("error reading message from stdin: %s\n", err)
	}
	connect(name, rb) //Only open a new connection when we want to say something
}

func newServerMessageConnection(sm chan string) {
	// plan is to use this func to allow server to write to client

}
func handleServerMessage(m string) {
	log.Printf("Server message: %s\n", m)
}
func listenForServerMessages(fromServerChan chan<- string) {
	smChan := make(chan string)
	defer newServerMessageConnection(smChan)
	resp := <-smChan //blocks here?
	handleServerMessage(resp)
	fromServerChan <- resp
}

func join() func() {
	var n string = USERNAME
	fmt.Println("Hi! You found the chat!")
	fmt.Printf("What's your name?>")
	name := bufio.NewReader(os.Stdin)
	n, _ = name.ReadString('\n')
	n = strings.ReplaceAll(n, "\n", "")
	//tod omatch name regex
	fmt.Printf("Welcome, %s! Type your message and press 'Enter' to send. Ctrl+C exits.\n\n", n)

	srvrMsg("42069", n)
	serverResponse := make(chan string, 1)
	go listenForServerMessages(serverResponse)
	// Leave open connection for server messages?
	// go letsRead()
	return func() {
		fmt.Printf("%s>", n)
		reader := bufio.NewReader(os.Stdin)
		rb, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("error reading message from stdin: %s\n", err)
		}
		connect(n, rb) //Only open a new connection when we want to say something
	}
}
func RunClient() {
	defer initClientLog().Close()
	yMsg := join()
	for {
		yMsg() // disconnect to reset name
		// msg(n)
	}
}
