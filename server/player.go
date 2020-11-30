package server

import (
	"fmt"
	"github.com/Ixidi/flaming/conn"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"runtime"
)

type Player struct {
	entityId        int
	conn            conn.Connection
	name            string
	uuid            uuid.UUID
	incomingPackets chan conn.Packet
	outgoingPackets chan conn.Packet
}

func (p *Player) Message(message string) {
	panic("implement me")
}

func (p *Player) Name() string {
	return p.name
}

func (p *Player) Uuid() uuid.UUID {
	return p.uuid
}

func (p *Player) Kick(reason string) {
	p.conn.Disconnect(reason)
}

func (p *Player) startReading(e *error) {
	for {
		packet, err := p.conn.ReadPacket()
		if err != nil {
			*e = err
			return
		}

		select {
		case p.incomingPackets <- packet:
			log.WithFields(log.Fields{
				"player":    p.name,
				"uuid":      p.uuid.String(),
				"packet_id": fmt.Sprintf("0x%x", packet.Id()),
				"goroutine": runtime.NumGoroutine(),
			}).Debug("Received packet from playing player.")
		default:
			log.WithFields(log.Fields{
				"player":    p.name,
				"uuid":      p.uuid.String(),
				"packet_id": fmt.Sprintf("0x%x", packet.Id()),
				"goroutine": runtime.NumGoroutine(),
			}).Warn("Received packet from playing player, but it was discarded because buffer is full.")
		}

	}
}

func (p *Player) startWriting(e *error) {
	for {
		packet := <-p.outgoingPackets
		err := p.conn.WritePacket(packet)
		if err != nil {
			*e = err
			return
		}

		log.WithFields(log.Fields{
			"player":    p.name,
			"uuid":      p.uuid.String(),
			"packet_id": fmt.Sprintf("0x%x", packet.Id()),
			"goroutine": runtime.NumGoroutine(),
		}).Debug("Sent packet to playing player.")
	}
}
