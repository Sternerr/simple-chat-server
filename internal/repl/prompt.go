package repl

import (
	"fmt"
	"bufio"
	"os"
	"log"
)

func PromptUsername(client Repl, logger *log.Logger) Repl {
	fmt.Printf("Please provide a display name: ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		username := scanner.Text()
		if len(username) < 4 {
			fmt.Printf("Username to short, please try again: ")
		} else { 
			(*logger).Printf("Username provided: %s\n", username)
			client.Cfg.Username = username
			break
		}
	}

	return client
}
