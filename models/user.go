package models

import "net"

type User struct {
	Name string
	Rooms []Room
	CurrentRoom Room
	Conn net.Conn
}