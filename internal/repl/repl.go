package repl

import (
	"net"
	"log"
	"bufio"
	"os"
	"bytes"
	"fmt"
)

type Repl struct {
	logger *log.Logger
}

func NewRepl(logger *log.Logger) Repl{
	return Repl{logger: logger}
}

func(r *Repl) Dial() {
	serverConn, err := net.Dial("tcp", ":6969")
	if err != nil {
		panic(err)
	}

	go func(c net.Conn) {
		str := ""
		for {
			buffer := make([]byte, 8)
			n, err := c.Read(buffer)
			if err != nil {
				break
			}

			buffer = buffer[:n]
			if i := bytes.IndexByte(buffer, '\n'); i != -1 {
				str += string(buffer)
				fmt.Print(str)
				buffer = buffer[i + 1:]
				str = ""
			}

			str += string(buffer)
		}

		if len(str) > 0 {
			fmt.Print(str)
		}
	}(serverConn)	

	for {
		reader := bufio.NewReader(os.Stdin)
		for {
			line, err := (*reader).ReadString('\n')
			if err != nil {
				break
			}

			serverConn.Write([]byte(line))
		}
	}
}
