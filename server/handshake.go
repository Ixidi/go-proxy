package server

import (
	"github.com/Ixidi/flaming/conn"
	"github.com/Ixidi/flaming/conn/packet"
	log "github.com/sirupsen/logrus"
	"net"
	"runtime"
)

func (s *Server) handleHandshake(c net.Conn) {
	log.WithFields(log.Fields{
		"goroutine": runtime.NumGoroutine(),
		"address":   c.RemoteAddr().String(),
	}).Debug("Handled new connection.")

	connection := conn.NewConnection(c)
	p, err := connection.ReadPacket()
	if err != nil {
		connection.Disconnect("")
		return
	}

	var handshake packet.Handshake
	err = p.Unpack(&handshake)

	if uint(handshake.Port) != s.port || p.Id() != packet.HandshakeId {
		connection.Disconnect("bad handshake packet")
		return
	}

	nextState := conn.ConnectionState(handshake.State)
	connection.UpdateState(nextState)
	if nextState == conn.StatusState {
		err := s.handleStatus(connection)
		if err != nil {
			log.WithFields(log.Fields{
				"goroutine": runtime.NumGoroutine(),
				"address":   c.RemoteAddr().String(),
				"error":     err.Error(),
			}).Debug("Error while handling status.")
		} else {
			log.WithFields(log.Fields{
				"goroutine": runtime.NumGoroutine(),
				"address":   c.RemoteAddr().String(),
			}).Debug("Status connection closed successfully.")
		}
	} else if nextState == conn.LoginState {
		player, err := s.handleLogin(connection)
		if err != nil {
			log.WithFields(log.Fields{
				"goroutine": runtime.NumGoroutine(),
				"address":   c.RemoteAddr().String(),
				"error":     err.Error(),
			}).Debug("Error while handling login.")
		} else {
			log.WithFields(log.Fields{
				"goroutine": runtime.NumGoroutine(),
				"address":   c.RemoteAddr().String(),
				"player":    player.name,
				"uuid":      player.uuid.String(),
			}).Info("Player has logged in successfully.")

			if err := s.spawnPlayer(&player); err != nil {
				log.WithFields(log.Fields{
					"goroutine": runtime.NumGoroutine(),
					"address":   c.RemoteAddr().String(),
					"error":     err.Error(),
					"player":    player.name,
					"uuid":      player.uuid.String(),
				}).Debug("Error while spawning.")
			}
		}
	} else {
		connection.Disconnect("wrong next state")
	}

}
