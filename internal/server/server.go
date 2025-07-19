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
	users map[net.Conn]User
}

func NewServer(logger *log.Logger) Server {
	return Server{
		logger: logger,
		users: make(map[net.Conn]User),
	}
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

		go (*s).handleConnection(clientConn)
	}
}

func(s *Server) handleConnection(clientConn net.Conn) {
	defer func() {
		clientConn.Close()
		delete(s.users, clientConn)
	}()

	for req := range (*s).processByteStream(clientConn) {
		msg, err := protocol.DecodeMessage([]byte(req))
		if err != nil {
			(*(*s).logger).Printf("%s: %s", err.Error(), req)
			return
		}

		switch msg.Type {
		case MessageTypeHandshake:
			(*(*s).logger).Printf("recieved handshake from: %s", clientConn)
			if protocol.IsValidHandshake(msg) {
				(*s).acceptHandshake(clientConn)
				(*s).users[clientConn] = User{Username: msg.From}
				break
			} else {
				(*s).denyHandshake("unknown 'from' in handshake", clientConn)
				return
			}
			
		case MessageTypeText:
			_, exists := (*s).users[clientConn]
			if !exists {
				(*s).denyHandshake("No handshake established", clientConn)
				return
			}

			(*(*s).logger).Printf("recieved text from: %s", clientConn)
			clientConn.Write([]byte(req))
		default:
			return
		}
	}
}

func(s *Server) acceptHandshake(clientConn net.Conn) {
	(*(*s).logger).Printf("accepted handshake from: %s\n", clientConn)
	res, err := protocol.EncodeMessage(Message{
		Type: MessageTypeHandshakeAccept,
		From: "server",
	})
	if err != nil {
		(*(*s).logger).Println(err.Error())
	}

	clientConn.Write(append(res, '\n'))
}

func(s *Server) denyHandshake(msg string, clientConn net.Conn) {
	(*(*s).logger).Printf("denied handshake from: %s\n", clientConn)
	res, err := protocol.EncodeMessage(Message{
		Type: MessageTypeHandshakeDeny,
		From: "server",
		Message: msg,
	})
	if err != nil {
		(*(*s).logger).Println(err.Error())
	}

	clientConn.Write(append(res, '\n'))
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
