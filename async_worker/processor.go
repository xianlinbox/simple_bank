package async_worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	db "github.com/xianlinbox/simple_bank/db/sqlc"
)

type TaskProcessor interface {
	Start() error
	HandleSendVerificationEmailTask(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) *RedisTaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Error().Err(err).Str("type", task.Type()).Msg("task processing error")
		}),
	})
	return &RedisTaskProcessor{
		server: server,
		store:  store,
	}
}
func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TASK_SEND_VERIFICATION_EMAIL, processor.HandleSendVerificationEmailTask)
	if err := processor.server.Start(mux); err != nil {
		log.Error().Msgf("could not start processor: %v", err)
		return err
	}
	log.Info().Msg("Task processor started")
	return nil
}
