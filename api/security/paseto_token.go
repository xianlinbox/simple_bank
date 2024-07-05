package security

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoTokenMaker struct {
	paseto *paseto.V2
	symmetricKey []byte
}

func NewPasetoTokenMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid symmetric key size")
	}
	return &PasetoTokenMaker{
		paseto: paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}, nil
}

func (maker *PasetoTokenMaker) GenerateToken(username string, duration time.Duration) (string, error){
	payload, err := NewPayload(username, duration)

	if err != nil {
		return "", err
	}
	return maker.paseto.Encrypt(maker.symmetricKey, &payload, nil)
}

func (maker *PasetoTokenMaker) VerifyToken(token string) (*Payload, error) {
	var payload Payload
	err:=maker.paseto.Decrypt(token, maker.symmetricKey, &payload, nil)
	if err != nil{
		return nil, err
	}

	if time.Now().After(payload.Expiration){
		return nil, &TokenExpireError{ expiredTime: payload.Expiration}
	}
	return &payload, nil
}