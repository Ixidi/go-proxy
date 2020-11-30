package server

import (
	"errors"
	"fmt"
	"github.com/Ixidi/flaming/command"
	"github.com/Ixidi/flaming/conn"
	log "github.com/sirupsen/logrus"
	"net"
	"testing/iotest"
	"time"
)

const (
	protocol = 754
	version  = "1.16.4"
)

type Server struct {
	port           uint
	protocol       int
	version        string
	threshold      int
	onlineMode     bool
	nextEntityId   int
	players        []*Player
	Motd           string
	ConsoleSender  command.Sender
	CommandManager command.Manager
}

func StartServer(config Config) *Server {
	port := config.Port
	l, err := net.Listen("tcp", fmt.Sprint(":", port))
	if err != nil {
		panic(err)
	}
	log.WithField("port", port).Info("Listening for connections...")
	s := Server{
		port:           port,
		protocol:       protocol,
		version:        version,
		threshold:      config.Threshold,
		onlineMode:     config.OnlineMode,
		nextEntityId:   0,
		Motd:           "Test",
		ConsoleSender:  &command.ConsoleSender{},
		CommandManager: command.NewManager(),
	}
	go func() {
		for {
			a, err := l.Accept()
			if err != nil {
				log.WithField("error", err.Error()).Debug("Failed to connect")
				continue
			}

			err = a.SetReadDeadline(time.Now().Add(time.Second * 3))
			if err != nil {
				log.WithField("error", err.Error()).Debug("Failed to set deadline.")
				continue
			}

			err = a.SetWriteDeadline(time.Now().Add(time.Second * 15))
			if err != nil {
				log.WithField("error", err.Error()).Debug("Failed to set deadline.")
				continue
			}
			go s.handleHandshake(a)
		}
	}()

	return &s
}

func (s *Server) startPlayerTask(player *Player) {
	s.players = append(s.players, player)

	var writingErr, readingErr error
	go player.startReading(&readingErr)
	go player.startWriting(&writingErr)

	for writingErr == nil && readingErr == nil {
		time.Sleep(time.Millisecond * 1000)
	}

	panic(readingErr) //sprawdzic czy nieaktywnosc jesli tak to kick i wiadomosc

	for i, p := range s.players {
		if p.entityId == player.entityId {
			if writingErr == iotest.ErrTimeout || readingErr == iotest.ErrTimeout {
				println("timeout!!!")
			}
			s.players = append(s.players[:i], s.players[i+1:]...)

			log.WithFields(log.Fields{
				"player": player.name,
				"uuid":   player.uuid.String(),
			}).Info("Player has disconnected.")
			return
		}
	}
}

func requestPacket(id conn.VarInt, v interface{}, c conn.Connection, dataSize int) (conn.Packet, error) {
	p, err := c.ReadPacket()
	if err != nil {
		return nil, err
	}
	if p.Id() != id {
		return nil, errors.New("expected another packet")
	}
	if dataSize >= 0 && len(p.Data()) != dataSize {
		return nil, errors.New("wrong packet size")
	}

	return p, p.Unpack(v)
}
