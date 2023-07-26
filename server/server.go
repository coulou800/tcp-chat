package server

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"tcp-chat/cmd"
	"tcp-chat/models"
	"time"
)

type Server struct {
	Listener    net.Listener
	Connections *sync.Map
	Rooms       *sync.Map
}

func (s *Server) HandleConn(conn net.Conn, msgChan chan models.Message, abrtConnChan chan net.Conn) {
	fmt.Fprint(conn, "\x1b[2J\x1b[H")
	conn.Write([]byte("Enter your name: "))
	var user models.User
	buff := make([]byte, 1024)
	len, err := conn.Read(buff)
	if err != nil {
		return
	}
	user.Name = string(buff[:len-1])
	user.CurrentRoom = s.GetRoom(0)
	user.Conn = conn
	fmt.Fprint(conn, "\x1b[1A")
	s.JoinRoom(conn, 0, &user, msgChan)

	s.Connections.Store(conn, &user)
	buff = make([]byte, 1024)
	fmt.Fprint(conn, " \x1b[38;5;51mtype here\x1b[5m\u21B3\x1b[25m\x1b[0m\t")

	for {
		len, err := conn.Read(buff)
		if err != nil {
			abrtConnChan <- conn
			conn.Close()
			return
		} else {
			req := models.Request{Input: buff[:len], From: conn}

			if cmd, err := cmd.ParseCmd(req); err != nil {
				var msg = models.Message{
					Content: buff[:len],
					Sender:  user.Name,
					Time:    time.Now(),
					Type:    models.MSGTYPE_USER_ROOM_MSG,
					Room:    user.CurrentRoom,
					From:    conn,
				}
				msgChan <- msg
			} else {
				s.ExecCmd(cmd, &user, msgChan)
			}

		}

	}
}

func NewServer() (s *Server) {
	addr := "0.0.0.0:8000"
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	s = &Server{
		Listener:    l,
		Connections: &sync.Map{},
		Rooms:       &sync.Map{},
	}
	s.InitGlobalRoom()
	return s
}

func (s *Server) NewConnection(ch chan net.Conn) {
	conn, err := s.Listener.Accept()
	if err != nil {
		log.Fatal("error accepting connection: ", err)
	}

	ch <- conn
}

func (s *Server) BroadcastMsg(msg models.Message) {
	if len(string(msg.Content)) != 0 {

		switch msg.Type {

		case models.MSGTYPE_SVER_JL_NOTIF:
			s.Connections.Range(func(conn, user any) bool {
				if msg.From.RemoteAddr() != conn.(net.Conn).RemoteAddr() && msg.Room.Number == user.(*models.User).CurrentRoom.Number {
					fmt.Fprint(conn.(net.Conn), msg.Prepare())
					fmt.Fprint(conn.(net.Conn), " \x1b[38;5;51mtype here\x1b[5m\u21B3\x1b[25m\x1b[0m\t")
				}
				return true
			})

		case models.MSGTYPE_SVER_RCREATED_NOTIF:
			s.Connections.Range(func(conn, user any) bool {
				if msg.From.RemoteAddr() == conn.(net.Conn).RemoteAddr() && msg.Room.Number == user.(*models.User).CurrentRoom.Number {

					fmt.Fprint(conn.(net.Conn), msg.Prepare())
					fmt.Fprint(conn.(net.Conn), " \x1b[38;5;51mtype here\x1b[5m\u21B3\x1b[25m\x1b[0m\t")

				}
				return true
			})
		case models.MSGTYPE_USER_ROOM_MSG:
			s.Connections.Range(func(conn, user any) bool {
				if msg.Sender != user.(*models.User).Name && msg.Room.Number == user.(*models.User).CurrentRoom.Number {
					fmt.Fprint(conn.(net.Conn), msg.Prepare())
				} else if msg.Sender == user.(*models.User).Name && msg.Room.Number == user.(*models.User).CurrentRoom.Number {
					tmp := msg.Sender
					msg.Sender = "You"
					fmt.Fprint(conn.(net.Conn), msg.Prepare())
					msg.Sender = tmp

				}
				fmt.Fprint(conn.(net.Conn), " \x1b[38;5;51mtype here\x1b[5m\u21B3\x1b[25m\x1b[0m\t")

				return true
			})

		}
	}

}

func (s *Server) InitGlobalRoom() {

	gb := models.Room{
		Number: 0,
	}
	gb.CreateLogFile()
	s.Rooms.Store(gb.Number, gb)
}

func (s *Server) ExecCmd(c cmd.Cmd, u *models.User, msgChan chan models.Message) {
	switch c.Name {
	case cmd.CMD_CREATE:
		s.CreateRoom(c.From, msgChan, u)
	case cmd.CMD_JOIN:
		n, _ := strconv.Atoi(strings.TrimSpace(c.Arg))
		s.JoinRoom(c.From, n, u, msgChan)
	}
}

func (s *Server) CreateRoom(conn net.Conn, msgChan chan models.Message, u *models.User) {
	r := models.Room{
		Number: rand.Intn(100),
	}
	r.CreateLogFile()
	s.Rooms.Store(r.Number, r)
	fmt.Println(s.GetRoomCount())
	str := fmt.Sprintf("Room #%v has been created\n", r.Number)
	file, _ := os.Open(r.LogFile)
	msg := models.Message{
		Sender:  "Server",
		Type:    models.MSGTYPE_SVER_RCREATED_NOTIF,
		Content: []byte(str),
		Room:    u.CurrentRoom,
		Time:    time.Now(),
		From:    conn,
	}
	fmt.Fprintf(file, ":8000 %s", msg.Prepare())

	s.JoinRoom(conn, r.Number, u, msgChan)
}

func (s *Server) GetRoom(n int) models.Room {
	room, _ := s.Rooms.Load(n)
	return room.(models.Room)
}

func (s *Server) JoinRoom(conn net.Conn, roomNumber int, u *models.User, msgChan chan models.Message) {
	if roomNumber != u.CurrentRoom.Number {

		str := fmt.Sprintf("%v left the chat", u.Name)
		msg := models.Message{
			Content: []byte(str),
			Sender:  "Server",
			Type:    models.MSGTYPE_SVER_JL_NOTIF,
			Room:    u.CurrentRoom,
			Time:    time.Now(),
			From:    conn,
		}
		msgChan <- msg

	}

	room := s.GetRoom(roomNumber)
	u.CurrentRoom = room
	fmt.Fprint(conn, "\x1b[2J\x1b[H")
	allMsg, err := room.LoadMsg(u)
	if err != nil {
		fmt.Fprintf(conn, "\x1b[48;5;142m\x1b[38;5;236m\u231A %s\u2B1F\x1b[38;5;26m\x1b[48;5;236m\x1b[1m Server\033[0m: \x1b[38;5;246m\x1b[3mCouldn't load previous message\n\033[0m", time.Now().Format("2006.01.02 15:04:05"))
	} else {
		conn.Write(allMsg)
	}
	fmt.Fprintf(conn, "\x1b[48;5;142m\x1b[38;5;236m\u231A %s\u2B1F\x1b[38;5;26m\x1b[48;5;236m\x1b[1m Server\033[0m: \x1b[38;5;246m\x1b[3mYou have joigned room #%d \n\033[0m", time.Now().Format("2006.01.02 15:04:05"), room.Number)

	str := fmt.Sprintf("%v has joigned the chat\n", u.Name)
	msg := models.Message{
		Content: []byte(str),
		Sender:  "Server",
		Type:    models.MSGTYPE_SVER_JL_NOTIF,
		Room:    u.CurrentRoom,
		Time:    time.Now(),
		From:    conn,
	}
	msgChan <- msg
}

func (s *Server) GetRoomCount() int {
	var c int
	s.Rooms.Range(func(key, value any) bool {
		c++
		return true
	})
	return c
}
