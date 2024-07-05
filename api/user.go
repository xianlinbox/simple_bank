package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/xianlinbox/simple_bank/db/sqlc"
	util "github.com/xianlinbox/simple_bank/util"
)

type AddUserRequest struct {
	Username string `json:"Username" binding:"required"`
	Password string `json:"Password" binding:"required" validate:"min=8"`
	FullName string `json:"FullName" binding:"required"`
	Email    string `json:"Email" binding:"required" validate:"email"`
}


func  (server *ApiServer) CreateUser(c *gin.Context) {
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
	user, err := server.store.AddUser(c, db.AddUserParams{
		Username: req.Username,
		Password: hashedPassword,
		FullName: req.FullName,
		Email: req.Email,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, user)
}