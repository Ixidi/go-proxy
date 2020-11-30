package main

import (
	"github.com/Ixidi/flaming/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors:   false,
		TimestampFormat: "2006-01-02\t15:04:05",
		FullTimestamp:   true,
		ForceColors:     true,
	})
	log.SetLevel(log.DebugLevel)
	s := server.StartServer(server.DefaultConfig())

	err := registerCommands(s)
	if err != nil {
		panic(err)
	}
	listenForCommands(s)
}
