package crypt

import (
	"crypto/aes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

type KeyPair struct {
	Public  []byte
	Private []byte
}

func Encrypt(key []byte, data []byte) (encrypted []byte, err error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	c.Encrypt(encrypted, data)
	return
}

func Decrypt(key []byte, data []byte) (decrypted []byte, err error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	c.Decrypt(decrypted, data)
	return
}

/*
func EncodePublicKeyDER(key rsa.PublicKey) ([]byte, error) {
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&key)
	if err != nil {
		return nil , err
	}

	publicKey := &pem.Block{
		Type:    "PUBLIC KEY",
		Bytes:   asn1Bytes,
	}

	var publicKeyBuff bytes.Buffer
	err = pem.Encode(&publicKeyBuff, publicKey)
	if err != nil {
		return nil, err
	}

	return publicKeyBuff.Bytes(), nil
}

func EncodePrivateKeyDER(key rsa.PrivateKey) ([]byte, error) {
	asn1Bytes, err := x509.MarshalPKCS8PrivateKey(&key)
	if err != nil {
		return nil , err
	}

	publicKey := &pem.Block{
		Type:    "PRIVATE KEY",
		Bytes:   asn1Bytes,
	}

	var privateKeyBuff bytes.Buffer
	err = pem.Encode(&privateKeyBuff, publicKey)
	if err != nil {
		return nil, err
	}

	return privateKeyBuff.Bytes(), nil
}*/

func GenerateKeyPair() (KeyPair, error) {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return KeyPair{}, err
	}

	privateBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return KeyPair{}, err
	}

	private := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: privateBytes,
		},
	)

	pubASN1, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		return KeyPair{}, err
	}

	public := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: pubASN1,
		},
	)

	return KeyPair{
		Public:  public,
		Private: private,
	}, nil
}
