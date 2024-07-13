package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xianlinbox/simple_bank/async_worker"
	db "github.com/xianlinbox/simple_bank/db/sqlc"
	util "github.com/xianlinbox/simple_bank/util"

	log "github.com/rs/zerolog/log"
)

type AddUserRequest struct {
	Username string `json:"Username" binding:"required"`
	Password string `json:"Password" binding:"required" validate:"min=8"`
	FullName string `json:"FullName" binding:"required"`
	Email    string `json:"Email" binding:"required" validate:"email"`
}

func (server *ApiServer) CreateUser(c *gin.Context) {
	var req AddUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	hashedPassword, err := util.EncryptPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	addUserParams := db.AddUserParams{
		Username: req.Username,
		Password: hashedPassword,
		FullName: req.FullName,
		Email:    req.Email,
	}
	args := db.CreateUserTxParams{
		AddUserParams: addUserParams,
		AfterCreate: func(user db.User) error {
			payload := async_worker.SendVerificationEmailTaskPayload{
				Username: user.Username,
			}
			err := server.distributor.DistributeSendVerificationEmailTask(c, &payload, nil)
			if err != nil {
				log.Error().Err(err).Msg("Failed to send verification email, user")
			}
			return nil
		},
	}
	user, err := server.store.CreateUserTx(c, server.db_conn, args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, user)
}

type LoginRequest struct {
	Username string `json:"Username" binding:"required"`
	Password string `json:"Password" binding:"required" validate:"min=8"`
}
type LoginResponse struct {
	User        string `json:"user"`
	AccessToken string `json:"access_token"`
}

func (server *ApiServer) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.store.GetUser(c, req.Username)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(user.Password, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	newToken, err := server.tokenMaker.GenerateToken(user.Username, time.Minute*15)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	response := LoginResponse{
		User:        user.Username,
		AccessToken: newToken,
	}
	c.JSON(http.StatusOK, response)
}
