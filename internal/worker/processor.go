package worker

import (
	"context"
	"encoding/json"
	"fmt"
	db "github.com/AutomaticOrca/GoShopping/internal/database/sqlc"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

type TaskProcessor interface {
	Start() error
	Shutdown()
	ProcessTaskOrderCancel(ctx context.Context, task *asynq.Task) error
	ProcessTaskPaymentTimeout(ctx context.Context, task *asynq.Task) error
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) TaskProcessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				logrus.WithFields(logrus.Fields{
					"task_type": task.Type(),
					"payload":   string(task.Payload()),
				}).Errorf("Process task failed: %v", err)
			}),
		},
	)

	return &RedisTaskProcessor{
		server: server,
		store:  store,
	}
}

func (p *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskOrderCancel, p.ProcessTaskOrderCancel)
	mux.HandleFunc(TaskPaymentTimeout, p.ProcessTaskPaymentTimeout)
	return p.server.Start(mux)
}

func (p *RedisTaskProcessor) ProcessTaskOrderCancel(ctx context.Context, task *asynq.Task) error {
	var payload PayloadOrderCancel
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		logrus.Errorf("Invalid payload: %v", err)
		return fmt.Errorf("invalid payload: %w", err)
	}

	err := p.store.CancelOrder(ctx, payload.OrderID)
	if err != nil {
		logrus.Errorf("Failed to cancel order: %v", err)
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"order_id": payload.OrderID,
	}).Info("Order cancelled")
	return nil
}

func (p *RedisTaskProcessor) ProcessTaskPaymentTimeout(ctx context.Context, task *asynq.Task) error {
	var payload PayloadPaymentTimeout
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		logrus.Errorf("Invalid payload: %v", err)
		return fmt.Errorf("invalid payload: %w", err)
	}

	paymentID, err := db.StringToInt32(payload.PaymentID)
	if err != nil {
		return fmt.Errorf("invalid PaymentID: %s", payload.PaymentID)
	}

	err = p.store.CancelPayment(ctx, paymentID)
	if err != nil {
		logrus.Errorf("Failed to cancel payment: %v", err)
		return fmt.Errorf("failed to cancel payment: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"payment_id": payload.PaymentID,
	}).Info("Payment cancelled")
	return nil
}

func (p *RedisTaskProcessor) Shutdown() {
	p.server.Shutdown()
}
