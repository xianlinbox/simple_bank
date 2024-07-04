package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/xianlinbox/simple_bank/db/sqlc"
)

type ApiServer struct {
	store db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *ApiServer {	
	router := gin.Default()
	server := &ApiServer{
		store: store,
	}
	if engine,ok := binding.Validator.Engine().(*validator.Validate); ok {
		engine.RegisterValidation("positiveAccountID", validateAccountID)
	}
	router.POST("/accounts", server.CreateAccount)
	router.GET("/accounts", server.ListAccounts)
	router.GET("/accounts/:id", server.GetAccount)
	server.router = router
	return server
}

func (server *ApiServer) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}