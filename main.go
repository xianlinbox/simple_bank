package main

import (
	"context"
	"net"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
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
		log.Error().Msgf("failed to load config: %+v", err)
	}
	pgConn, err := pgx.Connect(context.Background(), config.DBSource)
	store := db.New(pgConn)
	if err != nil {
		log.Error().Msgf("cannot connect to db: %+v", err)
	}
	tokenMaker, err := security.NewPasetoTokenMaker(config.SymmetricKey)
	if err != nil {
		log.Error().Msgf("can't create token maker: %+v", err)
	}

	redisClientopt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	go runApiServer(store, pgConn, tokenMaker, config, redisClientopt)
	go runTaskProcessor(store, redisClientopt)
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
		log.Error().Msgf("cannot create listener: %+v", err)
	}
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Error().Msgf("cannot start server: %+v", err)
	}
}

func runApiServer(store *db.Queries, db_conn *pgx.Conn, tokenMaker security.Maker, config util.Config, redisClientopt asynq.RedisClientOpt) {
	distributor := async_worker.NewRedisDistributor(&redisClientopt)
	err := api.NewServer(store, db_conn, tokenMaker, distributor).Start(config.ServerAddress)
	if err != nil {
		log.Error().Msgf("cannot start server: %+v", err)
	}
}

func runTaskProcessor(store db.Store, redisClientopt asynq.RedisClientOpt) {
	processor := async_worker.NewRedisTaskProcessor(redisClientopt, store)
	processor.Start()
}
