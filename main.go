package main

import (
	"os"
	"log"
	"fmt"
	"strings"

	"github.com/sternerr/termtalk/internal/server"
	"github.com/sternerr/termtalk/internal/repl"
	"github.com/sternerr/termtalk/internal/tui"
	
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("please provide a cli command")
		os.Exit(0)
	}

	cmd := os.Args[1]

	port, err := parseFlag("--port", os.Args[2:])	
	if err != nil {
		fmt.Println("error: ", err.Error())
		os.Exit(1)
	}
	switch cmd {
	case "server":
		logger := getLogger("/tmp/termtalk_server.txt")
		(*logger).Println("started server")

		server := server.NewServer(logger)
		server.Listen(port)
	case "client":
		logger := getLogger("/tmp/termtalk_client.txt")

		mode, err := parseFlag("--mode", os.Args[2:])	
		if err != nil {
			fmt.Println("error: ", err.Error())
			os.Exit(1)
		}

		host, err := parseFlag("--host", os.Args[2:])	
		if err != nil {
			fmt.Println("error: ", err.Error())
			os.Exit(1)
		}

		switch strings.ToLower(mode) {
		case "repl":
			(*logger).Printf("starting client with mode %s\n", mode)
			client := repl.NewRepl(logger)
			client = repl.PromptUsername(client, logger)
			client.Dial(host, port)
		case "tui":
			(*logger).Printf("[info] starting client with mode %s\n", mode)
			c := tui.NewTUI(logger)
			c.Client.Dial(host, port)
			p := tea.NewProgram(c, tea.WithAltScreen())
			if err := p.Start(); err != nil {
				fmt.Printf("Error starting program: %v\n", err)
				os.Exit(1)
			}

		default:
			(*logger).Printf("client mode %s do not exists\n", mode)
			fmt.Printf("client mode %s do not exists\n", mode)
		}
	default:
		fmt.Println("invalid command")
		os.Exit(0)
	}
}

func parseFlag(flag string, args []string) (string, error) {
	for i := 0; i < len(args); i += 1 {
		if args[i] == flag {
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") {
				return strings.ToLower(args[i+1]), nil
			}

			return "", fmt.Errorf("%s flag requires a value", flag)
			return "", fmt.Errorf("%s flag requires a value", flag)
		}
	}

	return "", fmt.Errorf("%s flag is missing", flag)
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("could not open file")
	}

	return log.New(logfile, "[termtalk]", log.Ldate|log.Ltime|log.Lshortfile)
}
