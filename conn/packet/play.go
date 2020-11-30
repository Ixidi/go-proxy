package packet

import "github.com/Ixidi/flaming/conn"

type SpawnPosition struct {
	Position conn.Position
}

type PlayerPositionAndLook struct {
	X          conn.Double
	Y          conn.Double
	Z          conn.Double
	Yaw        conn.Float
	Pitch      conn.Float
	Flags      conn.Byte
	TeleportId conn.VarInt
}
