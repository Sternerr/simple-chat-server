package repl

import (
	"net"
	"log"
	"bufio"
	"os"
	"bytes"
	"fmt"

	"github.com/sternerr/termtalk/internal/protocol"
	. "github.com/sternerr/termtalk/pkg/models"
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
	
	(*r).sendHandshake(serverConn)

	go func(c net.Conn) {
		reader := bufio.NewReader(os.Stdin)
		for {
			line, err := (*reader).ReadString('\n')
			if err != nil {
				break
			}

			req, err := protocol.EncodeMessage(Message{
				Type: MessageTypeText,
				From: "A",
				Message: string(line),
			})

			c.Write(append([]byte(req), '\n'))
		}
	}(serverConn)

	for res := range (*r).processByteStream(serverConn) {
		msg, err := protocol.DecodeMessage([]byte(res))
		if err != nil {
			(*(*r).logger).Println(err.Error())
		}

		switch msg.Type {
		case MessageTypeHandshakeDeny:
			fmt.Println(protocol.FormatMessage(msg))
			break
		case MessageTypeHandshakeAccept:
			fmt.Println("connected to server")
			continue
		case MessageTypeText:
			fmt.Print(protocol.FormatMessage(msg))
			break
		}
	}

}

func (r *Repl) sendHandshake(serverConn net.Conn) {
	req, err := protocol.EncodeMessage(Message{
		Type: MessageTypeHandshake,
		From: "A",
	})
	if err != nil {
		(*(*r).logger).Println(err.Error())
	}
	
	serverConn.Write(append(req, '\n'))
	(*(*r).logger).Println("sent handshake")
}

func(r *Repl) processByteStream(serverConn net.Conn) <-chan string {
	out := make(chan string)

	go func(c net.Conn) {
		defer close(out)
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
				out <- str
				buffer = buffer[i + 1:]
				str = ""
			}

			str += string(buffer)
		}

		if len(str) > 0 {
			out <- str
		}
	}(serverConn)

	return out
}
