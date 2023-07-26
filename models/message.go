package models

import (
	"bytes"
	"fmt"
	"net"
	"time"
)

type Message struct {
	Content []byte
	Sender  string
	From    net.Conn
	Time    time.Time
	Room    Room
	Type    string
}

const (
	MSGTYPE_SVER_JL_NOTIF       = "JL_NOTIF"
	MSGTYPE_SVER_RCREATED_NOTIF = "RCREATED_NOTIF"
	MSGTYPE_USER_ROOM_MSG       = "ROOM_MSG"
	MSGTYPE_USER_PRIV_MSG       = "PRIV_MSG"
)

func (msg *Message) Prepare() string {
	var str string
	msg.Content = bytes.TrimSpace(msg.Content)
	switch msg.Type {
	case MSGTYPE_SVER_JL_NOTIF:
		str = fmt.Sprintf("\x1b[48;5;142m\x1b[38;5;236m\u231A %s\u2B1F\x1b[38;5;26m\x1b[48;5;236m\x1b[1m %s\033[0m: \x1b[38;5;246m\x1b[3m%s\033[0m\n", msg.Time.Format("2006.01.02 15:04:05"), msg.Sender, msg.Content)
	case MSGTYPE_SVER_RCREATED_NOTIF:
		str = fmt.Sprintf("\x1b[48;5;142m\x1b[38;5;236m\u231A %s\u2B1F\x1b[38;5;26m\x1b[48;5;236m\x1b[1m %s\033[0m: \x1b[38;5;246m\x1b[3m%s\033[0m\n", msg.Time.Format("2006.01.02 15:04:05"), msg.Sender, msg.Content)
	case MSGTYPE_USER_ROOM_MSG:
		if msg.Sender == "You" {
			str = fmt.Sprintf("\x1b[1A\x1b[48;5;142m\x1b[38;5;236m\u231A %s\u2B1F\x1b[38;5;41m\x1b[48;5;236m\x1b[1m %s\033[0m: \x1b[38;5;45m%s\033[0m\n", msg.Time.Format("2006.01.02 15:04:05"), msg.Sender, msg.Content)
		} else {
			str = fmt.Sprintf("\x1b[48;5;142m\x1b[38;5;236m\u231A %s\u2B1F\x1b[38;5;142m\x1b[48;5;236m\x1b[1m %s\033[0m: \x1b[38;5;193m%s\033[0m\n", msg.Time.Format("2006.01.02 15:04:05"), msg.Sender, msg.Content)
		}
	}
	return str
}
