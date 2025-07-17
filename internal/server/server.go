package server

import (
	"net"
	"fmt"
	"log"
	"bytes"
)

type Server struct {
	listener net.Listener
	users []net.Conn
}

func (s *Server) Listen() {
	fmt.Printf("Listening on %s\n", s.listener.Addr())

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Fatal("Error: ", err.Error())
		}
		

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	
	for msg := range s.processBytestream(conn) {
		conn.Write([]byte(msg))
	}
}

func (s *Server) processBytestream(conn net.Conn) <-chan string {
	out := make(chan string)
	
	go func() {
		defer close(out)
		str := ""

		buffer := make([]byte, 8)
		for {
			n, err := conn.Read(buffer)
			if err != nil {
				break
			}

			buffer = buffer[:n]
			if i := bytes.IndexByte(buffer, '\n'); i != -1 {
				str += string(buffer)
				out <-str
				buffer = buffer[i+1:]
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

func NewServer(host, port string) (Server, error) {
	listener, err := net.Listen("tcp", host + ":" + port)
	if err != nil {
		return Server{}, err
	}

	return Server{
		listener: listener,
		users: []net.Conn{},
	}, nil
}
