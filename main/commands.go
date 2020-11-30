package main

import (
	"bufio"
	"fmt"
	"github.com/Ixidi/flaming/command"
	"github.com/Ixidi/flaming/command/executor"
	"github.com/Ixidi/flaming/server"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

func registerCommands(s *server.Server) error {
	manager := s.CommandManager

	motd := executor.Motd{S: s}
	creators := []command.Creator{
		manager.
			Create("motd").
			SetDesc("Manage current motd.").
			SetExecutor(motd.Executor),
	}

	for _, creator := range creators {
		err := creator.Register()
		if err != nil {
			return err
		}
	}

	return nil
}

func listenForCommands(s *server.Server) {
	reader := bufio.NewReader(os.Stdin)
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			log.Error(err)
		}
		cmdArray := strings.Split(strings.Trim(str, "\n\r"), " ")
		cmdName := cmdArray[0]
		cmd := s.CommandManager.Match(cmdName)
		if cmd == nil {
			s.ConsoleSender.Message(fmt.Sprintf("Unknown command '%s'.", cmdName))
			continue
		}

		var args []string
		if len(cmdArray) > 1 {
			args = cmdArray[1:]
		}
		cmd.Execute(s.ConsoleSender, args)
	}
}
