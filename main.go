package main

import (
	"fmt"
	"net"
	"tcp-chat/logger"
	"tcp-chat/models"
	"tcp-chat/server"
	"time"
)

func main() {
	server := server.NewServer()
	newConn := make(chan net.Conn)
	newMsg := make(chan models.Message, 1)
	abortedConn := make(chan net.Conn)
	host, port, err := net.SplitHostPort(server.Listener.Addr().String())

	if err != nil {
		panic(err)
	}

	fmt.Printf("Listening on host : %s, port: %s\n", host, port)
	n := server.GetRoomCount()
	println(n)

	for {
		if server.GetRoomCount() > n {
			n = server.GetRoomCount()
			println(n)
		}
		go server.NewConnection(newConn)
		select {
		case conn := <-newConn:
			go server.HandleConn(conn, newMsg, abortedConn)
		case abrtConn := <-abortedConn:
			user, _ := server.Connections.LoadAndDelete(abrtConn)
			server.Connections.Delete(abrtConn)
			str := fmt.Sprintf("%v left the chat", user.(*models.User).Name)
			msg := models.Message{Content: []byte(str), Sender: "Server", Time: time.Now(), Type: "serverMsg", Room: user.(*models.User).CurrentRoom,From: user.(*models.User).Conn}
			newMsg <- msg
		case msg := <-newMsg:
			go server.BroadcastMsg(msg)
			logger.MsgLogger(msg)
		}
	}
}
