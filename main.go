package main

import (
	"log"
	"os"

	"github.com/sternerr/termtalk/internal/server"
	"github.com/sternerr/termtalk/internal/client"
)

func main() {
	cmd := os.Args[1]
	switch cmd {
	case "server":
		srv, err := server.NewServer("0.0.0.0", "6969")
		if err != nil {
			log.Fatal(err.Error())
		}
		srv.Listen()
		break
	
	case "client":
		username, err := client.AskForUsername()
		if err != nil {
			log.Fatal("Error: ", err.Error())	
		}

		cl := client.NewClient(username)	
		err = cl.Connect()
		if err != nil {
			log.Fatal("Error: ", err.Error())	
		}

	default:
		log.Fatalf("unkonw command: %s", cmd)
	}
}
