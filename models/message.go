package models

import (
	"fmt"
	"time"
)

type Message struct {
	Content []byte
	Sender  string
	Time    time.Time
	Room    Room
	Type    string
}

func (msg *Message) Prepare() string {
	var str string
	switch msg.Type {
	case "serverMsg":
		str =  fmt.Sprintf("\x1b[48;5;142m\x1b[38;5;236m\u231A %s\u2B1F\x1b[38;5;26m\x1b[48;5;236m\x1b[1m %s\033[0m: \x1b[38;5;246m\x1b[3m%s\n\033[0m", msg.Time.Format("2006.01.02 15:04:05"), msg.Sender, msg.Content)
	case "userMsg":
		str =  fmt.Sprintf("\x1b[48;5;142m\x1b[38;5;236m\u231A %s\u2B1F\x1b[38;5;142m\x1b[48;5;236m\x1b[1m %s\033[0m: \x1b[38;5;193m%s\n\033[0m", msg.Time.Format("2006.01.02 15:04:05"), msg.Sender, msg.Content)

	}
	return str
}
