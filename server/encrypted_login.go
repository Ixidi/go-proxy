package server

import (
	"fmt"
	"github.com/Ixidi/flaming/conn"
	"github.com/Ixidi/flaming/conn/packet"
	"github.com/Ixidi/flaming/util/crypt"
	"math/rand"
)

func (s *Server) processEncryptedLogin(c conn.Connection) error {
	pair, err := crypt.GenerateKeyPair()
	if err != nil {
		return err
	}

	token := make([]byte, 4)
	rand.Read(token)

	encryptionRequest := packet.EncryptionRequest{
		ServerId:    "",
		PublicKey:   pair.Public,
		VerifyToken: token,
	}

	p, err := conn.Pack(packet.EncryptionRequestId, &encryptionRequest)
	if err != nil {
		return err
	}

	err = c.WritePacket(p)
	if err != nil {
		return err
	}

	var encryptionResponse packet.EncryptionResponse
	p, err = requestPacket(packet.EncryptionResponseId, &encryptionResponse, c, -1)
	if err != nil {
		return err
	}

	sharedSecret, err := crypt.Decrypt(pair.Private, encryptionResponse.SharedSecret)
	if err != nil {
		return err
	}

	verifyToken, err := crypt.Decrypt(pair.Private, encryptionResponse.VerifyToken)
	if err != nil {
		return err
	}

	fmt.Println(sharedSecret)
	fmt.Println(token)
	fmt.Println(verifyToken)
	return nil
}
