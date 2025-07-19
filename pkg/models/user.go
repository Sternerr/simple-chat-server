package models

import (
	"net"
)

type User struct {
	Username string
	Conn net.Conn
}
