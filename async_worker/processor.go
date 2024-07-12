package async_worker

import (
	"context"

	"github.com/hibiken/asynq"
	db "github.com/xianlinbox/simple_bank/db/sqlc"
)

type TaskProcessor interface {
	DistributeSendVerificationEmailTask(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) *RedisTaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{})
	return &RedisTaskProcessor{
		server: server,
		store:  store,
	}
}
