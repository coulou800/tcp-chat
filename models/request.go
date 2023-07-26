package models

import "net"

type Request struct{
	Input []byte
	From net.Conn
}