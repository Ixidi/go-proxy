package packet

import "github.com/Ixidi/flaming/conn"

type Request struct{}

type Response struct {
	Json conn.String
}

type PingPong struct {
	Payload conn.Long
}
