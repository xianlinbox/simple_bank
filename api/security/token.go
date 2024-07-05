package security

import (
	"fmt"
	"time"

	uuid "github.com/google/uuid"
)

type Maker interface {
	GenerateToken(username string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}

type Payload struct {
	Id			 uuid.UUID `json:"id"`
	Username string `json:"username"`
	IssuedAt time.Time `json:"issued_at"`
	Expiration time.Time 	`json:"expiration"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Payload{
		Id: tokenID,
		Username: username,
		IssuedAt: time.Now(),
		Expiration: time.Now().Add(duration),
	}, nil
}

type TokenExpireError struct{
	expiredTime time.Time
}


func (e *TokenExpireError) Error() string {
	return fmt.Sprintf("token is expired at %v", e.expiredTime.String())
}