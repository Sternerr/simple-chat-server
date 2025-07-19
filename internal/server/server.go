package server

import (
	"net"
	"bytes"
	"log"
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
					c.Write([]byte(str))
					buffer = buffer[i + 1:]
					str = ""
				}

				str += string(buffer)
			}

			if len(str) > 0 {
				c.Write([]byte(str))
			}
		}(clientConn)
	}
}
