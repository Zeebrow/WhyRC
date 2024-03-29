package server

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

	"github.com/Zeebrow/whyrc/shared"
)

type user struct {
	id      string
	present bool
}

type SRVMSG int

type ClientConnection struct {
	Name       string // name entered when joining (42069)
	ChToServer chan string
	ChToClient chan string
	Conn       *net.Conn
}

const (
	ACK   SRVMSG = 200
	JOIN  SRVMSG = 420
	LEAVE SRVMSG = 86
	WRITE SRVMSG = 69
)

func (b *Board) NewMessage(m shared.Message) error {
	fmt.Fprintf(b.writeTo, "%s> %s", m.From, m.Message)
	b.messages = append(b.messages, m)
	return nil
}

func NewClientConnection(n string) ClientConnection {
	chToServer := make(chan string)
	chToClient := make(chan string)
	return ClientConnection{Name: n, ChToServer: chToServer, ChToClient: chToClient}
}

func writeToClient(conn net.Conn, code int, msg string) {
	var respMsg shared.Message
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
		room.Post(shared.Message{From: sender.(string), Message: message})
		writeToClient(conn, 200, "You posted a message!")
		break
	case int:
		// "control" message for server (user has joined = 42069)
		fmt.Printf("A new user '%s' joined (code %d)\n", message, sender)
		for i := 0; i < MAX_CLIENT_CONNECTIONS; i++ {
			if clients[i].Name == "nobody" {
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
	board    [MAX_ROOM_MESSAGES]shared.Message
	users    [MAX_CLIENT_CONNECTIONS]user
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

	var board Board
	board = NewBoard("stdout")
	var gathering []user
	room = Room{name: "lolcats", users: gathering, board: &board}

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
		fmt.Println(<-receiveNewMsg)
	}
}
