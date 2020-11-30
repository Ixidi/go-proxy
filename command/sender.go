package command

import (
	log "github.com/sirupsen/logrus"
)

type Sender interface {
	Message(message string)
}

type ConsoleSender struct{}

func (c *ConsoleSender) Message(message string) {
	log.Info(message)
}
