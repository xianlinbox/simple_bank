package db

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type Store interface {
	Querier
	CreateUserTx(c *gin.Context, con *pgx.Conn, params CreateUserTxParams) (*User, error)
}
