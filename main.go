package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	api "github.com/xianlinbox/simple_bank/api"
	db "github.com/xianlinbox/simple_bank/db/sqlc"
)


const (
	dbSource = "postgresql://root:admin@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = ":8080"
)
func main() {
	pgConn, err := pgx.Connect(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	err = api.NewServer(db.New(pgConn)).Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}