package executor

import (
	"fmt"
	"github.com/Ixidi/flaming/command"
	"github.com/Ixidi/flaming/server"
	"strings"
)

type Motd struct {
	S *server.Server
}

func (e *Motd) Executor(sender command.Sender, args []string) {
	if len(args) == 0 {
		sender.Message(fmt.Sprintf("Current motd is: '%s'. Type 'motd <new motd>' to change it.", e.S.Motd))
		return
	}

	newMotd := strings.Join(args, " ")
	e.S.Motd = newMotd
	sender.Message(fmt.Sprintf("Motd has been set to '%s'.", newMotd))
}
