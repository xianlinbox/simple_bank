package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/xianlinbox/simple_bank/api/security"
	mockdb "github.com/xianlinbox/simple_bank/db/sqlc/mock"
)

const (
	TEST_SYMMETRIC_KEY = "p-S3cr3tT0p-S3cr3tT0p-T0p-S3cr3t"
)

func TestAuthMiddleware(t *testing.T) {
	testcases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker security.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker security.Maker) {
				token, err := tokenMaker.GenerateToken("test user", time.Minute)
				require.NoError(t, err)
				request.Header.Set(authorizationHeader, fmt.Sprintf("%v %v", authorizationType, token))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Unauthorized",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker security.Maker) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Unsupported Authorization Header",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker security.Maker) {
				token, err := tokenMaker.GenerateToken("test user", time.Minute)
				require.NoError(t, err)
				request.Header.Set(authorizationHeader, fmt.Sprintf("%v %v", "unsupoorted", token))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Invalid Token",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker security.Maker) {
				token, err := tokenMaker.GenerateToken("test user", -time.Minute)
				require.NoError(t, err)
				request.Header.Set(authorizationHeader, fmt.Sprintf("%v %v", authorizationType, token))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			store := mockdb.NewMockStore(controller)
			tokenMaker, err := security.NewPasetoTokenMaker(TEST_SYMMETRIC_KEY)
			require.NoError(t, err)
			server := NewServer(store, nil, tokenMaker, nil)
			server.router.GET("/auth", authMiddleware(tokenMaker), func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{})
			})

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/auth", nil)
			tc.setupAuth(t, request, tokenMaker)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}
