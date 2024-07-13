package db

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type CreateUserTxParams struct {
	AddUserParams AddUserParams
	AfterCreate   func(user User) error
}

func (store *Queries) CreateUserTx(c *gin.Context, con *pgx.Conn, params CreateUserTxParams) (*User, error) {
	tx, err := con.Begin(c)
	defer tx.Rollback(c)
	if err != nil {
		return nil, err
	}

	storeWithTx := store.WithTx(tx)
	user, err := storeWithTx.AddUser(c, params.AddUserParams)
	if err != nil {
		return nil, err
	}

	err = params.AfterCreate(user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
