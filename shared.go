package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
)

type SRVMSG int

const (
	ACK   SRVMSG = 200
	JOIN  SRVMSG = 420
	LEAVE SRVMSG = 86
	WRITE SRVMSG = 69
)

type Message struct {
	Code    int    `json:"code"`
	From    string `json:"from"`
	Message string `json:"message"`
}
type MessageHandler interface {
	HandleMessage(code SRVMSG)
}

type ClientConnection struct {
	name string      // name entered when joining (42069)
	ch   chan string // channel for sending a client messages from server
	conn *net.Conn
}

type User struct {
	id      string
	present bool
}

type BoardMessage struct {
	from    string
	message string
}

func (b *Board) NewMessage(bm BoardMessage) error {
	f, err := os.OpenFile(b.filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0744)
	if err != nil {
		log.Printf("Could not open file '%s' for writing board messages.\n", b.filename)
		return err
	}
	fmt.Fprintf(f, "%s> %s", bm.from, bm.message)
	if err != nil {
		log.Fatal("Surely this will never happen Kappa")
		return err
	}
	b.messages = append(b.messages, bm)
	return nil
}

type Board struct {
	messages []BoardMessage
	filename string // why not a pointer to a file? because I'm dumb.
}

func NewBoard(fname string) Board {
	var m []BoardMessage
	return Board{messages: m, filename: fname}
}

type Room struct {
	name              string
	activeConnections int
	conns             []ClientConnection
	users             []User
	board             *Board
}

func (r *Room) Join(u User) error {
	if r.activeConnections < MAX_CLIENT_CONNECTIONS-1 {
		r.users = append(r.users, u)
		r.activeConnections++
		u.present = true
		return nil
	}
	r.board.NewMessage(BoardMessage{from: "__server__", message: u.id + "has joined."})
	return errors.New("room full")
}

func (r *Room) Post(m Message) {
	r.board.NewMessage(BoardMessage{from: m.From, message: m.Message})
}

func (r *Room) HandleMessage(code SRVMSG, msg Message) {
	switch code {
	case JOIN:
		fmt.Printf("Handle Join\n")
		newUser := User{id: msg.Message, present: false}
		eJoin := r.Join(newUser)
		if eJoin != nil {
			fmt.Println("Handle join error full room")
		}
		break
	case WRITE:
		fmt.Printf("Handle WRITE\n")
		break
	case LEAVE:
		fmt.Printf("Handle LEAVE\n")
		break
	default:
		fmt.Printf("Handle UNKNOWN (%d)", int(code))
	}
}
