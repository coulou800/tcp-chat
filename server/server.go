package server

import (
	"fmt"
	"log"
	"net"
	"sync"
	"tcp-chat/models"
	"time"
)

type Server struct {
	Listener    net.Listener
	Connections *sync.Map
}

func (s *Server) HandleConn(conn net.Conn, msgChan chan models.Message, abrtConnChan chan net.Conn) {
	fmt.Fprint(conn,"\x1b[2J\x1b[H")
	conn.Write([]byte("Enter your name: "))
	var name string
	fmt.Fscanln(conn, &name)

	s.Connections.Store(conn, name)
	buff := make([]byte, 1024)
	str := fmt.Sprintf("%s has joigned the chat\n", name)
	serverMsg := models.Message{
		Content: []byte(str),
		Sender:  "Server",
		Time:    time.Now(),
		Type: "serverMsg",
	}
	msgChan <- serverMsg

	for {
		len, err := conn.Read(buff)
		if err != nil {
			abrtConnChan <- conn
			conn.Close()
			return
		} else {
			var msg = models.Message{
				Content: buff[:len],
				Sender:  name,
				Time:    time.Now(),
				Type: "userMsg",
			}

			msgChan <- msg
		}

	}
}

func NewServer() *Server {
	addr := "0.0.0.0:8000"
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	return &Server{
		Listener:    l,
		Connections: &sync.Map{},
	}
}

func (s *Server) NewConnection(ch chan net.Conn) {
	conn, err := s.Listener.Accept()
	if err != nil {
		log.Fatal("error accepting connection: ", err)
	}

	ch <- conn
}

func (s *Server) BroadcastMsg(msg models.Message) {
	if string(msg.Content) != "" {

		s.Connections.Range(func(conn, name any) bool {
			if msg.Sender != name {
				
				fmt.Fprint(conn.(net.Conn), msg.Prepare())
			}
			return true
		})
	}
}
