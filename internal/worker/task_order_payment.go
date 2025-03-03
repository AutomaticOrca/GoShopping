package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	TaskOrderCancel    = "task:order_cancel"
	TaskPaymentTimeout = "task:payment_timeout"

	OrderCancelDelay    = 30 * time.Minute
	PaymentTimeoutDelay = 15 * time.Minute
)

type PayloadOrderCancel struct {
	OrderID int32 `json:"order_id"`
}

type PayloadPaymentTimeout struct {
	PaymentID string `json:"payment_id"`
}

func (d *RedisTaskDistributor) DistributeTaskOrderCancel(ctx context.Context, payload *PayloadOrderCancel, delay time.Duration) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Errorf("Failed to marshal task payload: %v", err)
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskOrderCancel, jsonPayload, asynq.ProcessIn(delay))
	_, err = d.client.EnqueueContext(ctx, task)
	if err != nil {
		log.Errorf("Failed to enqueue task: %v", err)
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.WithFields(log.Fields{
		"task":     TaskOrderCancel,
		"order_id": payload.OrderID,
	}).Info("Task enqueued")
	return nil
}

func (d *RedisTaskDistributor) DistributeTaskPaymentTimeout(ctx context.Context, payload *PayloadPaymentTimeout, delay time.Duration) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Errorf("Failed to marshal task payload: %v", err)
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskPaymentTimeout, jsonPayload, asynq.ProcessIn(delay))
	_, err = d.client.EnqueueContext(ctx, task)
	if err != nil {
		log.Errorf("Failed to enqueue task: %v", err)
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.WithFields(log.Fields{
		"task":       TaskPaymentTimeout,
		"payment_id": payload.PaymentID,
	}).Info("Task enqueued")
	return nil
}
