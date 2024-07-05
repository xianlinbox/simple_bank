package security

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPasetoToken(t *testing.T) {
	PasetoTokenMaker, err := NewPasetoTokenMaker("thisisasslongsstestingsskeyaaaaa")
	require.NoError(t, err)
	token,err:=PasetoTokenMaker.GenerateToken("test", time.Minute*15)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := PasetoTokenMaker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.Equal(t, "test", payload.Username)
}

func TestPasetoTokenExpired(t *testing.T) {
	PasetoTokenMaker, err := NewPasetoTokenMaker("this is a long testing keyaaaaaa")
	require.NoError(t, err)
	token,err:=PasetoTokenMaker.GenerateToken("test", -time.Minute*15)	
	require.NoError(t, err)
	_, err = PasetoTokenMaker.VerifyToken(token)
	require.ErrorContains(t, err, "token is expired")
}

