package client

import (
	"net"
	"log"
	"bytes"

	. "github.com/sternerr/termtalk/pkg/models"
	"github.com/sternerr/termtalk/internal/protocol"
	tea "github.com/charmbracelet/bubbletea"
)

type Connected bool

type Client struct {
	Conn net.Conn
	MsgChan chan tea.Msg
	logger *log.Logger
}

func NewClient(logger *log.Logger) Client {
	return Client{
		MsgChan: make(chan tea.Msg, 10),
		logger: logger,
	}
}

func(c *Client) Dial(host, port string) {
	conn, err := net.Dial("tcp", host + ":" + port)
	if err != nil {
		panic(err)
	}
	(*(*c).logger).Println("[info] server dialed")
	(*c).Conn = conn
	(*c).MsgChan <- Connected(true)

	go func() {
		for res := range (*c).processByteStream() {
			msg, err := protocol.DecodeMessage([]byte(res))
			if err != nil {
				(*(*c).logger).Println(err.Error())
			}

			switch msg.Type {
			case MessageTypeHandshakeDeny:
				c.MsgChan <- MessageTypeHandshakeDeny
			
			case MessageTypeHandshakeAccept:
				c.MsgChan <- MessageTypeHandshakeAccept
			
			case MessageTypeText:
				c.MsgChan <- msg
			}
		}
	}()
}

func(c *Client) processByteStream() <-chan string {
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
				str += string(buffer[:i])
				out <- str

				buffer = buffer[i+1:]
				str = ""
			}

			str += string(buffer)
		}

		if len(str) > 0 {
			out <- str
		}
	}((*c).Conn)

	return out
}
func (c *Client) SendHandshake(user User) {
	req, err := protocol.EncodeMessage(Message{
		Type: MessageTypeHandshake,
		From: user.Username,
	})
	if err != nil {
		(*(*c).logger).Println(err.Error())
	}
	
	(*(*c).logger).Println("[info] handshake sent")
	(*c).Conn.Write(append(req, '\n'))
}
