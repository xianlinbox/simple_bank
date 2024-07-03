package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	api "github.com/xianlinbox/simple_bank/api"
	db "github.com/xianlinbox/simple_bank/db/sqlc"
	util "github.com/xianlinbox/simple_bank/util"
)
func main() {
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("failed to load config: ", err)
	}

	pgConn, err := pgx.Connect(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	err = api.NewServer(db.New(pgConn)).Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}