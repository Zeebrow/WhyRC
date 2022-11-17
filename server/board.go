package server

import (
	"log"
	"os"

	"github.com/Zeebrow/whyrc/shared"
)

type Board struct {
	name     string
	messages []shared.Message
	writeTo  *os.File
}

func NewBoard(fname string) Board {
	var m []shared.Message
	if fname == "stdout" {
		return Board{messages: m, writeTo: os.Stdout}
	}
	of, err := os.OpenFile(fname, os.O_WRONLY, 0o600)
	if err != nil {
		log.Fatalf("Could not open file '%s' for writing\n", fname)
	}
	return Board{messages: m, writeTo: of}
}
