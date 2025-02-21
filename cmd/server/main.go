package main

import (
	"context"
	"fmt"
	"github.com/AutomaticOrca/GoShopping/configs"
	db "github.com/AutomaticOrca/GoShopping/internal/database/sqlc"
	"github.com/AutomaticOrca/GoShopping/internal/server"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
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

	server, err := server.NewServer(config, store)
	if err != nil {
		log.Fatalf("Failed to create new server: %v", err)
	}

	httpServerAddress := fmt.Sprintf("%s:%s", config.ServerAddress, config.Port)
	err = server.Start(httpServerAddress)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
