package server

import (
	"net"
	"bytes"
	"log"
	"fmt"

	"github.com/sternerr/termtalk/internal/protocol"
	. "github.com/sternerr/termtalk/pkg/models"
)

type Server struct {
	listener net.Listener
	logger *log.Logger
	users map[net.Conn]User
	history CircularBuffer
}

func NewServer(logger *log.Logger) Server {
	return Server{
		logger: logger,
		users: make(map[net.Conn]User),
		history: NewCircularBuffer(50),
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
			continue
		}

		switch msg.Type {
		case MessageTypeHandshake:
			(*(*s).logger).Printf("recieved handshake from: %s", clientConn)
			if protocol.IsValidHandshake(msg) {
				(*s).acceptHandshake(clientConn)
				(*s).users[clientConn] = User{Username: msg.From, Conn: clientConn}
				
				(*s).sendBroadcast(fmt.Sprintf("%s connected", msg.From))
				(*s).sendHistory(clientConn)
				defer (*s).sendBroadcast(fmt.Sprintf("%s disconnected", msg.From))
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

			(*(*s).logger).Printf("recieved text from: %s", msg)
			(*s).history.Add(msg)
			(*s).sendMessage([]byte(req), clientConn)
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

func(s *Server) sendHistory(clientConn net.Conn) {
	for _, msg := range (*s).history.GetAll() {
		res, err := protocol.EncodeMessage(msg)
		if err != nil {
			(*(*s).logger).Println(err.Error())
			continue
		}

		clientConn.Write(append(res, '\n'))
	}
}

func(s *Server) sendMessage(msg []byte, exclude net.Conn) {
	for _, u := range (*s).users {
		if u.Conn != exclude {
			u.Conn.Write(append(msg, '\n'))
		}
	}
}

func(s *Server) sendBroadcast(msg string) {
	(*(*s).logger).Printf("broadcast: %s", msg)
	res, err := protocol.EncodeMessage(Message{
		Type: MessageTypeText,
		From: "server",
		Message: msg + "\n",
	})
	if err != nil {
		(*(*s).logger).Println(err.Error())
	}

	for c, _ := range (*s).users {
		c.Write(append(res, '\n'))
	}
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
				str += string(buffer[:i])
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
