package packet

import "github.com/Ixidi/flaming/conn"

type Handshake struct {
	Protocol conn.VarInt
	Address  conn.String
	Port     conn.UShort
	State    conn.VarInt
}
