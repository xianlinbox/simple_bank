package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	db "github.com/xianlinbox/simple_bank/db/sqlc"
	mockdb "github.com/xianlinbox/simple_bank/db/sqlc/mock"
	"github.com/xianlinbox/simple_bank/util"
)

type eqCreateUserParamsMatcher struct {
	arg db.AddUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.AddUserParams)
	if !ok {
		return false
	}
	errr := util.CheckPassword(arg.Password,e.password)
	if errr != nil {
		return false
	}
	return arg.Username == e.arg.Username && arg.FullName == e.arg.FullName && arg.Email == e.arg.Email
}
func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches user: %v with password %v", e.arg, e.password)
}

func eqCreateUserParams(arg db.AddUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg,password}
}

func TestCreateUserAPI(t *testing.T) {

	user := db.User{
		Username:  "TestUser",
		FullName : "Test User",
		Email: "a@b.com",
		Password: "12345678",
	}
	params := db.AddUserParams{
		Username: user.Username,
		FullName: user.FullName,
		Password: user.Password,
		Email: user.Email,
	}

	controller := gomock.NewController(t)
	defer controller.Finish()
	store := mockdb.NewMockStore(controller)
	store.EXPECT().AddUser(gomock.Any(), eqCreateUserParams(params, user.Password)).Times(1).Return(user, nil)

	server := NewServer(store)
	recorder := httptest.NewRecorder()
	requestData, err := json.Marshal(user)
	require.NoError(t, err)
	request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(requestData))	
	require.NoError(t, err)
	server.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)
	response,err:=io.ReadAll(recorder.Body)
	require.NoError(t,err)
	var userInResponse db.User
	err =json.Unmarshal(response,&userInResponse)
	fmt.Println(userInResponse)
	require.NoError(t,err)
	require.Equal(t,user,userInResponse)
}