package server

import (
	"github.com/Ixidi/flaming/conn"
	"github.com/Ixidi/flaming/conn/packet"
	"github.com/Ixidi/flaming/misc"
)

func (s *Server) handleStatus(c conn.Connection) error {
	var request packet.Request
	p, err := requestPacket(packet.RequestId, &request, c, 0)
	if err != nil {
		c.Disconnect(err.Error())
		return err
	}

	status := misc.Status{
		Version: misc.StatusVersion{
			Name:     s.version,
			Protocol: s.protocol,
		},
		Players: misc.StatusPlayers{
			Max:     0,
			Online:  0,
			Players: make([]misc.StatusPlayer, 0),
		},
		Description: misc.StatusDescription{Text: s.Motd},
	}
	json, err := status.ToJson()
	if err != nil {
		c.Disconnect("server error")
		return err
	}

	response := packet.Response{Json: conn.String(json)}
	p, err = conn.Pack(packet.ResponseId, &response)
	if err != nil {
		c.Disconnect("server error")
		return err
	}

	err = c.WritePacket(p)
	if err != nil {
		c.Disconnect("server error")
		return err
	}

	var pingPong packet.PingPong
	p, err = requestPacket(packet.PingPongId, &pingPong, c, -1)
	if err != nil {
		c.Disconnect(err.Error())
		return err
	}

	err = c.WritePacket(p)
	if err != nil {
		c.Disconnect("server error")
		return err
	}

	c.Disconnect("end")
	return nil
}
