package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/xianlinbox/simple_bank/db/sqlc"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func  (server *ApiServer) CreateAccount(c *gin.Context) {
	var req createAccountRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.AddAccount(c, db.AddAccountParams{
		Owner: req.Owner,
		Balance: 0,
		Currency: req.Currency,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, account)
}

func (server *ApiServer) ListAccounts(c *gin.Context) {

	accouns, err := server.store.ListAccounts(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return		
	}

	c.JSON(http.StatusOK, accouns)
}

func (server *ApiServer) GetAccount(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"),10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(c, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return		
	}

	c.JSON(http.StatusOK, account)
}