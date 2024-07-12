package async_worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type Distributor interface {
	DistributeSendVerificationEmailTask(
		ctx context.Context,
		payload *SendVerificationEmailTaskPayload,
		opts ...asynq.Option) error
}

type RedisDistributor struct {
	client *asynq.Client
}

func NewRedisDistributor(redisOpt *asynq.RedisClientOpt) *RedisDistributor {
	client := asynq.NewClient(redisOpt)
	return &RedisDistributor{client: client}
}
