package server

import (
	"github.com/Ixidi/flaming/conn"
	"github.com/Ixidi/flaming/conn/packet"
	"github.com/Ixidi/flaming/premium"
	uuid "github.com/satori/go.uuid"
)

func (s *Server) handleLogin(c conn.Connection) (Player, error) {
	var loginStart packet.LoginStart
	_, err := requestPacket(packet.LoginStartId, &loginStart, c, -1)
	if err != nil {
		c.Disconnect(err.Error())
		return Player{}, err
	}
	username := loginStart.Username

	if s.onlineMode {
		isPremium, err := premium.FetchProfile(string(loginStart.Username))
		if err != nil {
			return Player{}, err
		}

		if isPremium {
			if err := s.processEncryptedLogin(c); err != nil {
				return Player{}, err
			}
		}
	}

	//TODO check online player
	setCompression := packet.SetCompression{Threshold: conn.VarInt(s.threshold)}
	p, err := conn.Pack(packet.SetCompressionId, &setCompression)
	if err != nil {
		return Player{}, err
	}

	err = c.WritePacket(p)
	if err != nil {
		return Player{}, err
	}
	c.SetThreshold(s.threshold)

	u := uuid.NewV5(uuid.Nil, string(username))

	loginSuccess := packet.LoginSuccess{
		Uuid:     conn.UUID(u),
		Username: username,
	}

	p, err = conn.Pack(packet.LoginSuccessId, &loginSuccess)
	if err != nil {
		return Player{}, err
	}

	err = c.WritePacket(p)
	if err != nil {
		return Player{}, err
	}

	entityId := s.nextEntityId
	s.nextEntityId++

	player := Player{
		entityId:        entityId,
		conn:            c,
		name:            string(username),
		uuid:            u,
		incomingPackets: make(chan conn.Packet, 10),
		outgoingPackets: make(chan conn.Packet, 10),
	}

	go s.startPlayerTask(&player)
	return player, nil
}
