package main

import (
	"context"
	"log"
	"net"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5"
	api "github.com/xianlinbox/simple_bank/api"
	"github.com/xianlinbox/simple_bank/api/security"
	"github.com/xianlinbox/simple_bank/async_worker"
	db "github.com/xianlinbox/simple_bank/db/sqlc"
	"github.com/xianlinbox/simple_bank/gapi"
	"github.com/xianlinbox/simple_bank/proto_code"
	util "github.com/xianlinbox/simple_bank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}
	pgConn, err := pgx.Connect(context.Background(), config.DBSource)
	store := db.New(pgConn)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	tokenMaker, err := security.NewPasetoTokenMaker(config.SymmetricKey)
	if err != nil {
		log.Fatal("can't create token maker: ", err)
	}

	redisClientopt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	distributor := async_worker.NewRedisDistributor(&redisClientopt)

	go runApiServer(store, pgConn, tokenMaker, config, distributor)
	runGrpcServer(store, tokenMaker, config)
}

func runGrpcServer(store *db.Queries, tokenMaker security.Maker, config util.Config) {
	server := gapi.NewServer(store, tokenMaker)
	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	proto_code.RegisterUsersServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal("cannot create listener: ", err)
	}
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}

func runApiServer(store *db.Queries, db_conn *pgx.Conn, tokenMaker security.Maker, config util.Config, distributor async_worker.Distributor) {
	err := api.NewServer(store, db_conn, tokenMaker, distributor).Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
