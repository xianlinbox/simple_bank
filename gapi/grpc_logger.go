package gapi

import (
	"context"
	"time"

	"github.com/rs/zerolog"

	log "github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func GrpcLogger(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error){
	log.Info().Msgf("gRPC method: %s", info.FullMethod)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	startTime := time.Now()
	r, e :=handler(ctx, req)
	duration := time.Since(startTime)
	log.Info().Str("method", info.FullMethod).Dur("duration", duration).Msg("gRPC call")
	return r, e
}