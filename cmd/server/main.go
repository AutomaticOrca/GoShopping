package main

import (
	"context"
	"fmt"
	"github.com/AutomaticOrca/GoShopping/configs"
	db "github.com/AutomaticOrca/GoShopping/internal/database/sqlc"
	"github.com/AutomaticOrca/GoShopping/internal/server"
	"github.com/AutomaticOrca/GoShopping/internal/worker"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	log.Print("Starting...")
	config, err := configs.LoadConfig(".")

	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", config.DbUsername, config.DbPassword, config.DbPort, config.DbName)

	connPool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}

	store := db.NewStore(connPool)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	waitGroup, ctx := errgroup.WithContext(ctx)

	runTaskProcessor(ctx, waitGroup, config, redisOpt, store)

	runGinServer(ctx, waitGroup, config, store, taskDistributor)

	err = waitGroup.Wait()
	if err != nil {
		log.Fatalf("Error from wait group: %v", err)
	}
}

func runTaskProcessor(
	ctx context.Context,
	waitGroup *errgroup.Group,
	config configs.Config,
	redisOpt asynq.RedisClientOpt,
	store db.Store,
) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store)

	waitGroup.Go(func() error {
		log.Info("Starting task processor...")
		err := taskProcessor.Start()
		if err != nil {
			log.Fatalf("Failed to start task processor: %v", err)
		}
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info("Gracefully shutting down task processor.")

		taskProcessor.Shutdown()
		log.Info("task processor is stopped.")

		return nil
	})
}

func runGinServer(
	ctx context.Context,
	waitGroup *errgroup.Group,
	config configs.Config,
	store db.Store,
	taskDistributor worker.TaskDistributor,
) {
	s, err := server.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatalf("Failed to create new server: %v", err)
	}

	httpServerAddress := fmt.Sprintf("%s:%s", config.ServerAddress, config.Port)

	waitGroup.Go(func() error {
		log.Infof("Starting HTTP server at %s", httpServerAddress)
		return s.Start()
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info("Gracefully shutting down HTTP server...")
		return s.Shutdown()
	})
}
