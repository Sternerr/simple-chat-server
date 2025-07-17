package main

import (
	"log"
	"os"

	"github.com/sternerr/termtalk/internal/server"
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

	default:
		log.Fatalf("unkonw command: %s", cmd)
	}
}
