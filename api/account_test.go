package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/xianlinbox/simple_bank/api/security"
	db "github.com/xianlinbox/simple_bank/db/sqlc"
	mockdb "github.com/xianlinbox/simple_bank/db/sqlc/mock"
)

func TestCreateAccountAPI(t *testing.T) {

	newAccount := db.Account{
		Currency: "USD",
	}

	params := db.AddAccountParams{
		Owner:    "testuser",
		Balance:  0,
		Currency: newAccount.Currency,
	}
	controller := gomock.NewController(t)
	defer controller.Finish()
	store := mockdb.NewMockStore(controller)
	store.EXPECT().AddAccount(gomock.Any(), gomock.Eq(params)).Times(1)
	tokenMaker, err := security.NewPasetoTokenMaker(TEST_SYMMETRIC_KEY)
	require.NoError(t, err)
	server := NewServer(store, tokenMaker, nil)
	recorder := httptest.NewRecorder()
	requestData, err := json.Marshal(newAccount)
	require.NoError(t, err)
	request, err := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(requestData))
	require.NoError(t, err)
	token, err := server.tokenMaker.GenerateToken("testuser", time.Minute)
	require.NoError(t, err)
	request.Header.Set(authorizationHeader, fmt.Sprintf("%v %v", "Bearer", token))
	require.NoError(t, err)
	server.router.ServeHTTP(recorder, request)
}
