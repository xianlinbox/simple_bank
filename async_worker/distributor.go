package async_worker

import "github.com/hibiken/asynq"

type Distributor interface{}

type RedisDistributor struct{
	client *asynq.Client
}

func NewRedisDistributor(redisOpt *asynq.RedisClientOpt) *RedisDistributor {
	client := asynq.NewClient(redisOpt)
	return &RedisDistributor{client: client}
}