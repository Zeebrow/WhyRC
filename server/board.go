package server

import "github.com/Zeebrow/whyrc/shared"

type Board struct {
	messages []shared.Message
	filename string // why not a pointer to a file? because I'm dumb.
}

func NewBoard(fname string) Board {
	var m []shared.Message
	return Board{messages: m, filename: fname}
}
