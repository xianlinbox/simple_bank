package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	security "github.com/xianlinbox/simple_bank/api/security"
	"github.com/xianlinbox/simple_bank/async_worker"
	db "github.com/xianlinbox/simple_bank/db/sqlc"
)

type ApiServer struct {
	store       db.Store
	db_conn     *pgx.Conn
	router      *gin.Engine
	tokenMaker  security.Maker
	distributor async_worker.Distributor
}

func NewServer(store db.Store, db_conn *pgx.Conn, tokenMaker security.Maker, distributor async_worker.Distributor) *ApiServer {

	router := gin.Default()
	server := &ApiServer{
		store:       store,
		db_conn:     db_conn,
		tokenMaker:  tokenMaker,
		distributor: distributor,
	}
	if engine, ok := binding.Validator.Engine().(*validator.Validate); ok {
		engine.RegisterValidation("positiveAccountID", validateAccountID)
	}
	router.POST("/users", server.CreateUser)
	router.POST("/user/login", server.Login)
	auth_group := router.Group("/")
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
