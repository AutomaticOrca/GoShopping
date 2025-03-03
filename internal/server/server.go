package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/AutomaticOrca/GoShopping/configs"
	db "github.com/AutomaticOrca/GoShopping/internal/database/sqlc"
	"github.com/AutomaticOrca/GoShopping/internal/worker"
	"github.com/AutomaticOrca/GoShopping/pkg/token"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Server struct {
	config          configs.Config
	store           db.Store
	taskDistributor worker.TaskDistributor
	tokenMaker      token.Maker
	router          *gin.Engine
	httpServer      *http.Server
}

func NewServer(config configs.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		log.Fatalf("Failed to create token maker: %v", err)
		return nil, err
	}

	server := &Server{
		config:          config,
		store:           store,
		taskDistributor: taskDistributor,
		tokenMaker:      tokenMaker,
	}

	server.setupRouter()

	server.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", config.ServerAddress, config.Port),
		Handler: server.router,
	}

	log.Info("Server initialized successfully")

	return server, nil
}

func (server *Server) Start() error {
	log.Infof("Starting HTTP server at %s:%s", server.config.ServerAddress, server.config.Port)

	go func() {
		if err := server.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	return nil
}

func (server *Server) Shutdown() error {
	log.Info("Shutting down HTTP server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.httpServer.Shutdown(ctx); err != nil {
		log.Errorf("Server forced to shutdown: %v", err)
		return err
	}

	log.Info("HTTP server gracefully stopped. ;)")
	return nil
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
