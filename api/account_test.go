package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/xianlinbox/simple_bank/api/security"
	db "github.com/xianlinbox/simple_bank/db/sqlc"
	mockdb "github.com/xianlinbox/simple_bank/db/sqlc/mock"
)

func TestCreateAccountAPI(t *testing.T) {

	newAccount := db.Account{
		Owner: "owner",
		Currency: "USD",
	}

	params := db.AddAccountParams{
		Owner: newAccount.Owner,
		Balance: 0,
		Currency: newAccount.Currency,
	}
	controller := gomock.NewController(t)
	defer controller.Finish()
	store := mockdb.NewMockStore(controller)
	store.EXPECT().AddAccount(gomock.Any(), gomock.Eq(params)).Times(1)

	server := NewServer(store, &security.PasetoTokenMaker{})
	recorder := httptest.NewRecorder()
	requestData, err := json.Marshal(newAccount)
	require.NoError(t, err)
	request, err := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(requestData))	
	require.NoError(t, err)
	server.router.ServeHTTP(recorder, request)
}