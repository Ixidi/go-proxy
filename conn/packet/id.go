package packet

const (
	HandshakeId             = 0x00
	RequestId               = 0x00
	ResponseId              = 0x00
	PingPongId              = 0x01
	LoginStartId            = 0x00
	SetCompressionId        = 0x03
	LoginSuccessId          = 0x02
	LoginDisconnectId       = 0x00
	EncryptionRequestId     = 0x01
	EncryptionResponseId    = 0x01
	SpawnPositionId         = 0x42
	PlayerPositionAndLookId = 0x34
	JoinGameId              = 0x24
)
