package packet

import "github.com/Ixidi/flaming/conn"

type LoginStart struct {
	Username conn.String
}

type SetCompression struct {
	Threshold conn.VarInt
}

type LoginSuccess struct {
	Uuid     conn.UUID
	Username conn.String
}

type LoginDisconnect struct {
	Reason conn.String
}

type EncryptionRequest struct {
	ServerId    conn.String
	PublicKey   conn.ByteArray
	VerifyToken conn.ByteArray
}

type EncryptionResponse struct {
	SharedSecret conn.ByteArray
	VerifyToken  conn.ByteArray
}
