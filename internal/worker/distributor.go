package worker

import (
	"context"
	"github.com/hibiken/asynq"
	"time"
)

type TaskDistributor interface {
	DistributeTaskOrderCancel(ctx context.Context, payload *PayloadOrderCancel, delay time.Duration) error
	DistributeTaskPaymentTimeout(ctx context.Context, payload *PayloadPaymentTimeout, delay time.Duration) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt)
	return &RedisTaskDistributor{client: client}
}
