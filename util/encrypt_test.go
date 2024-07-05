package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncryptPassword(t *testing.T) {

	hashedPassword, err := EncryptPassword("123456")
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	err = CheckPassword(hashedPassword, "123456")
	require.NoError(t, err)
}