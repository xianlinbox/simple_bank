package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/xianlinbox/simple_bank/db/sqlc"
)

type ApiServer struct {
	store *db.Queries
	router *gin.Engine
}

func NewServer(store *db.Queries) *ApiServer {	
	router := gin.Default()
	server := &ApiServer{
		store: store,
	}
	router.POST("/accounts", server.CreateAccount)
	router.GET("/accounts", server.ListAccounts)
	server.router = router
	return server
}

func (server *ApiServer) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}