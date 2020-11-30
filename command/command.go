package command

import "errors"

type Command interface {
	Name() string
	Desc() string
	Aliases() []string
	Execute(sender Sender, args []string)
}

type command struct {
	name     string
	desc     string
	aliases  []string
	executor func(sender Sender, args []string)
}

func (c *command) Name() string {
	return c.name
}

func (c *command) Desc() string {
	return c.desc
}

func (c *command) Aliases() []string {
	return c.aliases
}

func (c *command) Execute(sender Sender, args []string) {
	c.executor(sender, args)
}

type Creator interface {
	SetDesc(desc string) Creator
	AddAlias(alias string) Creator
	SetExecutor(executor func(sender Sender, args []string)) Creator
	Register() error
}

type creator struct {
	name       string
	commands   *[]Command
	registered bool
	desc       string
	aliases    []string
	executor   func(sender Sender, args []string)
}

func (c *creator) SetDesc(desc string) Creator {
	c.desc = desc
	return c
}

func (c *creator) AddAlias(alias string) Creator {
	c.aliases = append(c.aliases, alias)
	return c
}

func (c *creator) SetExecutor(executor func(sender Sender, args []string)) Creator {
	c.executor = executor
	return c
}

func (c *creator) Register() error {
	for _, cmd := range *c.commands {
		if cmd.Name() == c.name {
			return errors.New("duplicated command name")
		}

		for _, a := range cmd.Aliases() {
			if a == c.name {
				return errors.New("duplicated command name")
			}

			for _, al := range c.aliases {
				if a == al {
					return errors.New("duplicated command name")
				}
			}
		}
	}

	command := command{
		name:     c.name,
		desc:     c.desc,
		aliases:  c.aliases,
		executor: c.executor,
	}
	*c.commands = append(*c.commands, &command)
	return nil
}
