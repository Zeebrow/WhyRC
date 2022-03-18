package main

import (
	"log"
	"os"
)

func NewLogger(l string) Log {
	f, err := os.OpenFile("logs/"+l+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file '%s': %s\n", l, err)
	}
	logger := log.New(f, "["+l+"]: ", log.LstdFlags)
	return Log{loggerName: l, logger: logger}
}

type Log struct {
	loggerName string
	logger     *log.Logger
}

func (l *Log) Printf(m string, a ...interface{}) func(m string, a ...interface{}) {
	return func(m string, a ...interface{}) {
		l.logger.Printf(m, a)
	}
}
func (l *Log) Println(t string) func(m string) {
	return func(m string) {
		l.logger.Println(m)
	}
}

func (l *Log) Fatalf(m string, a ...interface{}) func(m string, a ...interface{}) {
	return func(m string, a ...interface{}) {
		l.logger.Fatalf(m, a)
	}
}

func (l *Log) Fatalln(t string) func(m string) {
	return func(m string) {
		l.logger.Fatalln(m)
	}
}
