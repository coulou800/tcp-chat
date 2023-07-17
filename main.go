package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type Server struct {
	Listener    net.Listener
	Connections *sync.Map
}

type Message struct {
	Content []byte
	Sender  string
	Time    time.Time
}

func main() {
	server := NewServer()
	newConn := make(chan net.Conn)
	newMsg := make(chan Message)
	abortedConn := make(chan net.Conn)
	host, port, err := net.SplitHostPort(server.Listener.Addr().String())

	if err != nil {
		panic(err)
	}

	fmt.Printf("Listening on host : %s, port: %s\n", host, port)

	for {

		go server.NewConnection(newConn)

		select {
		case conn := <-newConn:
			go server.HandleConn(conn, newMsg,abortedConn)
		case abrtConn := <- abortedConn:
			name,_ := server.Connections.LoadAndDelete(abrtConn)
			server.Connections.Delete(abrtConn)
			str := fmt.Sprintf("%v left the chat",name)
			newMsg <- Message{Content: []byte(str),Sender: "Server",Time: time.Now()}
		case msg := <-newMsg:
			go server.BroadcastMsg(msg)
		}
	}
}

func (s *Server) HandleConn(conn net.Conn, msgChan chan Message, abrtConnChan chan net.Conn) {
	conn.Write([]byte("Enter your name: "))
	var name string
	fmt.Fscanln(conn, &name)

	s.Connections.Store(conn, name)
	buff := make([]byte, 1024)
	str := fmt.Sprintf("%s has joigned the chat\n", name)
	serverMsg := Message{
		Content: []byte(str),
		Sender:  "Server",
		Time:    time.Now(),
	}
	msgChan <- serverMsg
	for {
		len, err := conn.Read(buff)
		if err != nil {
			abrtConnChan <- conn
			return
		} else{
			var msg = Message{
			Content: buff[:len],
			Sender:  name,
			Time:    time.Now(),
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
	// defer l.Close()
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

func (s *Server) BroadcastMsg(msg Message) {
	if string(msg.Content) != "" {

		s.Connections.Range(func(conn, name any) bool {
			if msg.Sender != name {
				fmt.Fprintf(conn.(net.Conn), "[%s] [%s]: %s", msg.Time.Format("2006.01.02 15:04:05"), msg.Sender, msg.Content)
			}
			return true
		})
	}
}
