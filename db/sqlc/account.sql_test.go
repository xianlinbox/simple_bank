package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddAccount(t *testing.T) {
	args := AddAccountParams{
		Owner:    "owner",
		Balance:  100,
		Currency: "USD",
	}
	account, err := testQueries.AddAccount(context.Background(), args)
	require.Nil(t, err)
	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Balance, account.Balance)
	require.Equal(t, args.Currency, account.Currency)
}