package main

import (
	"errors"
	"fmt"
	"net"
)

type Message struct {
	from    string
	message string
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

type Room struct {
	name              string
	activeConnections int
	conns             []ClientConnection
	users             []User
	board             []Message
}

type Roomer interface {
	Post(m Message) error
	Join(u User) error
}

func (r *Room) Join(u User) error {
	if r.activeConnections < MAX_CLIENT_CONNECTIONS-1 {
		r.users = append(r.users, u)
		r.activeConnections++
		u.present = true
		return nil
	}
	return errors.New("room full")
}

func (r *Room) Post(m Message) {
	r.board = append(r.board, m)
}

func (r *Room) HandleMessage(code SRVMSG, msg Message) {
	switch code {
	case JOIN:
		fmt.Printf("Handle Join\n")
		newUser := User{id: msg.message, present: false}
		eJoin := r.Join(newUser)
		if eJoin != nil {
			fmt.Println("Handle join error full room")
		}
		break
	case WRITE:
		fmt.Printf("Handle WRITE\n")
		break
	default:
		fmt.Printf("Handle UNKNOWN (%d)", int(code))
	}
}

type SRVMSG int

const (
	ACK   SRVMSG = 200
	JOIN  SRVMSG = 420
	LEAVE SRVMSG = 419
	WRITE SRVMSG = 69
)
