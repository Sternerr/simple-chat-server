package client

import (
	"net"
	"fmt"
	"bytes"
	"bufio"
	"os"
)

type Client struct {
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
			
			c.conn.Write([]byte(line))
		}
	}()
	
	for msg := range c.processBytestream(c.conn) {
		fmt.Printf(msg)
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

func NewClient() Client{
	return Client{}
}
