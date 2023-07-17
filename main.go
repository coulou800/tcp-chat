package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"tcp-chat/models"
	"tcp-chat/server"
	"time"

)

func main() {

	logFile, _ := os.Create("gb_msg.log")
	log.SetOutput(logFile)
	server := server.NewServer()
	newConn := make(chan net.Conn)
	newMsg := make(chan models.Message)
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
			go server.HandleConn(conn, newMsg, abortedConn)
		case abrtConn := <-abortedConn:
			name, _ := server.Connections.LoadAndDelete(abrtConn)
			server.Connections.Delete(abrtConn)
			str := fmt.Sprintf("%v left the chat", name)
			msg := models.Message{Content: []byte(str), Sender: "Server", Time: time.Now(),Type: "serverMsg"}
			go server.BroadcastMsg(msg)
		case msg := <-newMsg:
			go server.BroadcastMsg(msg)
		}
	}
}
