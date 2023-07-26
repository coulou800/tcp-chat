package cmd

import (
	"fmt"
	"net"
	"strings"
	"tcp-chat/models"
)

type Cmd struct {
	From net.Conn
	Name string
	Arg  string
}

const (
	CMD_JOIN   = "\\join"
	CMD_EXIT   = "\\exit"
	CMD_CREATE = "\\create"
	CMD_PRIV   = "\\priv"
	CMD_WHERE = "\\where"
)

func ParseCmd(req models.Request) (cmd Cmd, err error) {
	w := strings.Split(string(req.Input), ":")
	if len(w) != 3 || w[0] != "cmd" || !isCmd(w[1]) {
		err = fmt.Errorf("not a valid command, %v", req.Input)
		return cmd, err
	}
	cmd.From = req.From
	cmd.Name = w[1]
	cmd.Arg = w[2]

	return cmd, nil
}

func isCmd(s string) bool {
	return s == CMD_CREATE || s == CMD_EXIT || s == CMD_JOIN || s == CMD_PRIV || s == CMD_WHERE
}
