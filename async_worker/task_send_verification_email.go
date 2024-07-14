package async_worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	TASK_SEND_VERIFICATION_EMAIL = "send_verification_email"
)

type SendVerificationEmailTaskPayload struct {
	Username string `json:"username"`
}

func (distributor *RedisDistributor) DistributeSendVerificationEmailTask(
	ctx context.Context,
	payload *SendVerificationEmailTaskPayload,
	opts ...asynq.Option) error {
	mashalledPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("could not marshal payload: %w", err)
	}

	task := asynq.NewTask(TASK_SEND_VERIFICATION_EMAIL, mashalledPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("could not enqueue task: %w", err)
	}
	log.Info().Str("type", task.Type()).Str("queue", info.Queue).Msg("task enqueued")
	return nil
}

func (processor *RedisTaskProcessor) HandleSendVerificationEmailTask(ctx context.Context, task *asynq.Task) error {
	var payload SendVerificationEmailTaskPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		log.Error().Err(err).Msg("could not unmarshal payload")
		return err
	}
	_, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		log.Error().Err(err).Str("username", payload.Username).Msg("could not get user")
		return err
	}
	log.Info().Str("type", task.Type()).Msgf("task processed: %v", task)
	return nil
}
