package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xianlinbox/simple_bank/api/security"
	db "github.com/xianlinbox/simple_bank/db/sqlc"
)

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func  (server *ApiServer) CreateAccount(c *gin.Context) {
	var req createAccountRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	auth_payload :=c.MustGet(AUTH_KEY).(*security.Payload)
	account, err := server.store.AddAccount(c, db.AddAccountParams{
		Owner: auth_payload.Username,
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
	auth_payload:= c.MustGet(AUTH_KEY).(*security.Payload)
	params :=db.GetAccountsByOwnerParams{
		Owner: auth_payload.Username,
		Limit: 50,
	}

	accounts, err := server.store.GetAccountsByOwner(c, params)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return		
	}

	c.JSON(http.StatusOK, accounts)
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