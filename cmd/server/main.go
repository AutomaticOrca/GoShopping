package main

import (
	"context"
	"github.com/AutomaticOrca/GoShopping/configs"
	"github.com/AutomaticOrca/GoShopping/internal/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Print("Starting...")
	config, err := configs.LoadConfig(".")

	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	srv := &server.Server{}

	ctx := context.Background()
	if err := srv.Create(ctx, &config); err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	if err := srv.Serve(ctx); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
