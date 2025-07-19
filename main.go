package main

import (
	"os"
	"log"
	"fmt"
)

func main() {
	logger := getLogger("/tmp/termtalk_log.txt")
	(*logger).Println("I will get the job done!")

	if len(os.Args) <= 1 {
		(*logger).Println("no cli command provided")
		fmt.Println("please provide a cli command")
		os.Exit(0)
	}


	cmd := os.Args[1]
	switch cmd {
	case "server":
		(*logger).Println("started server")
	case "client":
		(*logger).Println("started client")
	default:
		(*logger).Println("invalid command")
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
