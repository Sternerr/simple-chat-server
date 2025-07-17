package server

import (
	"net"
	"fmt"
	"log"
	"bytes"

	"github.com/sternerr/termtalk/internal/protocol"
	. "github.com/sternerr/termtalk/pkg/models"
)

type Server struct {
	listener net.Listener
	users []User
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
		msg, err := protocol.DecodeMessage([]byte(msg))
		if err != nil {
			log.Fatal(err.Error())
			break
		}


		switch msg.Type {
		case MessageTypeHandshake:
			user := User{ Conn: conn, Username: msg.From, }
			s.users = append(s.users, user)

			msg, err := protocol.EncodeMessage(Message{
				Type: MessageTypeChat,
				From: "server",
				Message: fmt.Sprintf("%s Connected\n", user.Username),
			})
			if err != nil {
				log.Fatal(err.Error())
			}

			s.sendMessage(msg, conn)
			break
		case MessageTypeChat:
			msg, err := protocol.EncodeMessage(msg)
			if err != nil {
				break
			}

			s.sendMessage(msg, conn)
			break
		default:
			fmt.Println("Invalid type")	
			break
		}
	}
}

func (s *Server) processBytestream(conn net.Conn) <-chan string {
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

func (s *Server) sendMessage(msg []byte, exclude net.Conn) {
	for _, u := range s.users {
		if u.Conn != exclude {
			u.Conn.Write(msg)
		}
	}
}

func NewServer(host, port string) (Server, error) {
	listener, err := net.Listen("tcp", host + ":" + port)
	if err != nil {
		return Server{}, err
	}

	return Server{
		listener: listener,
		users: []User{},
	}, nil
}
