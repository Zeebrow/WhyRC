package main

import (
	"log"
	"os"
)

type ServerLog struct {
}
type ServerLogger interface {
	Debug(msg string)
}

func (*ServerLog) Debug(s string) {
	log.Printf("[server]: %s\n", s)
}

func initLog() *os.File {
	logFilePath := "logs/server.log"

	f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("Could not open logfile '%s' for writing\n", logFilePath)
	}
	log.SetOutput(f)
	return f
}

func initClientLog() *os.File {
	logFilePath := "logs/client.log"

	f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("Could not open logfile '%s' for writing\n", logFilePath)
	}
	log.SetOutput(f)
	return f
}
