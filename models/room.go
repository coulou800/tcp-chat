package models

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Room struct {
	Number  int
	Member  []User
	LogFile string
}

func (r *Room) CreateLogFile() error {
	filename := fmt.Sprintf("msg-logs/%v_msg.log", r.Number)
	_, err := os.Create(filename)
	if err != nil {
		return err
	}
	r.LogFile = filename
	return nil
}

func (r *Room) LoadMsg(u *User) ([]byte, error) {
	filepath := fmt.Sprintf("msg-logs/%v_msg.log", r.Number)
	file, err := os.Open(filepath)

	var str string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		host := strings.Split(line, " ")[0]

		if host == u.Conn.RemoteAddr().String() {
			i := strings.Index(line,u.Name)
			str += string(line[len(host)+1:i])+ "You"+string(line[i+len(u.Name):])+"\n"
			continue
		}
		str += line[len(host)+1:]+"\n"
	}

	return []byte(str), err
}
