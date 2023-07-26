package logger

import (
	"fmt"
	"os"
	"tcp-chat/models"
)

func MsgLogger(msg models.Message) error {
	file, err := os.OpenFile(msg.Room.LogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	str := fmt.Sprintf("%v %v", msg.From.RemoteAddr().String(),msg.Prepare())
	file.WriteString(str)
	return nil
}
