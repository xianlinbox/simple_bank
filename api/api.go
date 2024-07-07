package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	security "github.com/xianlinbox/simple_bank/api/security"
	db "github.com/xianlinbox/simple_bank/db/sqlc"
)

type ApiServer struct {
	store db.Store
	router *gin.Engine
	tokenMaker security.Maker
}

func NewServer(store db.Store, tokenMaker security.Maker) *ApiServer {	

	router := gin.Default()
	server := &ApiServer{
		store: store,
		tokenMaker: tokenMaker,
	}
	if engine,ok := binding.Validator.Engine().(*validator.Validate); ok {
		engine.RegisterValidation("positiveAccountID", validateAccountID)
	}
	router.POST("/users", server.CreateUser)
	router.POST("/user/login", server.Login)
	auth_group :=router.Group("/")
	auth_group.Use(authMiddleware(server.tokenMaker))
	auth_group.POST("/accounts", server.CreateAccount)
	auth_group.GET("/accounts", server.ListAccounts)
	auth_group.GET("/accounts/:id", server.GetAccount)
	
	server.router = router
	return server
}

func (server *ApiServer) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}