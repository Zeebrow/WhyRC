package server

import (
	"errors"
	"fmt"

	"github.com/Zeebrow/whyrc/shared"
)

type Room struct {
	name              string
	activeConnections int
	conns             []ClientConnection
	users             []user
	board             *Board
}

func (r *Room) Join(u user) error {

	if r.activeConnections < MAX_CLIENT_CONNECTIONS-1 {
		r.users = append(r.users, u)
		r.activeConnections++
		u.present = true
		return nil
	}
	r.board.NewMessage(shared.Message{From: "__server__", Message: u.id + "has joined."})
	return errors.New("room full")
}

func (r *Room) Post(m shared.Message) {
	r.board.NewMessage(shared.Message{From: m.From, Message: m.Message})
}

func (r *Room) HandleMessage(code SRVMSG, msg shared.Message) {
	switch code {
	case JOIN:
		fmt.Printf("Handle Join\n")
		newUser := user{id: msg.Message, present: false}
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
