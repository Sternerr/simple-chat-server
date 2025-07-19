package server

import (
	"net"
	"bytes"
	"log"

	"github.com/sternerr/termtalk/internal/protocol"
	. "github.com/sternerr/termtalk/pkg/models"
)

type Server struct {
	listener net.Listener
	logger *log.Logger
}

func NewServer(logger *log.Logger) Server {
	return Server{logger: logger}
}

func(s *Server) Listen() {
	listener, err := net.Listen("tcp", ":6969")
	if err != nil {
		(*(*s).logger).Printf("could not start listener: %s\n", err.Error())
		panic(err)
	}
	(*(*s).logger).Println("Listening on 0.0.0.0:6969")
	(*s).listener = listener
	
	(*s).accept()
}

func(s *Server) accept() {
	for {
		clientConn, err := (*s).listener.Accept()
		if err != nil {
			(*(*s).logger).Printf("client connection rejected: %s\n", err.Error())
		}
		(*(*s).logger).Printf("client connected: %s\n", clientConn)
	
		go (*s).handleConnection(clientConn)
	}
}

func(s *Server) handleConnection(clientConn net.Conn) {
	defer clientConn.Close()

	for req := range (*s).processByteStream(clientConn) {
		msg, err := protocol.DecodeMessage([]byte(req))
		if err != nil {
			(*(*s).logger).Println(err.Error())
			(*s).denyHandshake(clientConn)
			break
		}

		switch msg.Type {
		case MessageTypeHandshake:
			break
		default:
			(*s).denyHandshake(clientConn)
			break
		}
	}
}

func(s *Server) denyHandshake(clientConn net.Conn) {
	(*(*s).logger).Println("denied handshake: %s", clientConn)
	res, err := protocol.EncodeMessage(Message{
		Type: MessageTypeHandshakeDeny,
		From: "server",
	})
	if err != nil {
		(*(*s).logger).Println(err.Error())
	}

	clientConn.Write(res)
}

func(s *Server) processByteStream(clientConn net.Conn) <-chan string {
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
	}(clientConn)

	return out
}
