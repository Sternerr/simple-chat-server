package client

import (
	"net"
	"fmt"
	"bytes"
	"bufio"
	"os"
	"errors"
	"strings"
	"log"

	"github.com/sternerr/termtalk/internal/protocol"
	. "github.com/sternerr/termtalk/pkg/models"
)

type config struct {
	username string
}

type Client struct {
	cfg config
	conn net.Conn
}

func (c *Client) Connect() error {
	conn, err := net.Dial("tcp", ":6969")
	if err != nil {
		return err
	}
	c.conn = conn

	c.handleConnection()

	return nil
}

func (c *Client) handleConnection() {
	defer c.conn.Close()
	
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			
			msg, err := protocol.EncodeMessage(Message{
				Type: MessageTypeChat,
				From: c.cfg.username,
				Message: string(line),
			})
			
			c.conn.Write(msg)
		}
	}()
	
	for msg := range c.processBytestream(c.conn) {
		msg, err := protocol.DecodeMessage([]byte(msg))
		if err != nil {
			log.Fatal("Error: ", err.Error())
		}
		fmt.Printf("[%s] %s", msg.From, msg.Message)
	}
}

func (c *Client) processBytestream(conn net.Conn) <-chan string {
	out := make(chan string)
	
	go func() {
		defer close(out)

		str := ""
		for {
			buffer := make([]byte, 8)
			n, err := conn.Read(buffer)
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
	}()

	return out
}

func NewClient(username string) Client {
	return Client{
		cfg: config{
			username: username,
		},
	}
}

func AskForUsername() (string, error) {
	fmt.Printf("Enter your username: ")
	reader := bufio.NewReader(os.Stdin)
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", errors.New("could not read username from input")
	}

	username = strings.TrimRight(username, "\r\n")
	return username, nil
}
