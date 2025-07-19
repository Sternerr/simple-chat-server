package main

import (
	"os"
	"log"
	"fmt"

	. "github.com/sternerr/termtalk/internal/server"
	. "github.com/sternerr/termtalk/internal/repl"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("please provide a cli command")
		os.Exit(0)
	}


	cmd := os.Args[1]
	switch cmd {
	case "server":
		logger := getLogger("/tmp/termtalk_server.txt")
		(*logger).Println("started server")

		server := NewServer(logger)
		server.Listen()
	case "client":
		logger := getLogger("/tmp/termtalk_client.txt")
		(*logger).Println("started client")
		repl := NewRepl(logger)
		repl.Dial()
	default:
		fmt.Println("invalid command")
		os.Exit(0)
	}
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("could not open file")
	}

	return log.New(logfile, "[termtalk]", log.Ldate|log.Ltime|log.Lshortfile)
}
