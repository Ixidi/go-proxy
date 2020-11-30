package command

import (
	"strings"
)

type Manager interface {
	Match(name string) Command
	Create(name string) Creator
	Remove(name string)
}

type manager struct {
	commands []Command
}

func NewManager() Manager {
	return &manager{}
}

func (m *manager) Match(name string) Command {
	low := strings.ToLower(name)
	for _, c := range m.commands {
		if c.Name() == low {
			return c
		}

		for _, a := range c.Aliases() {
			if a == low {
				return c
			}
		}
	}

	return nil
}

func (m *manager) Create(name string) Creator {
	return &creator{
		name:       name,
		commands:   &m.commands,
		registered: false,
		desc:       "Default desc.",
		aliases:    nil,
		executor: func(sender Sender, args []string) {
			sender.Message("Default executor.")
		},
	}
}

func (m *manager) Remove(name string) {
	low := strings.ToLower(name)
	for i, c := range m.commands {
		if c.Name() == low {
			m.commands = append(m.commands[:i], m.commands[i+1:]...)
			return
		}
	}
}
