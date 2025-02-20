package server

import (
	"context"
	"fmt"
	"github.com/AutomaticOrca/GoShopping/configs"
	db "github.com/AutomaticOrca/GoShopping/internal/database/sqlc"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	Config *configs.Config
	HTTP   *http.Server
	Router *mux.Router
	DB     *pgxpool.Pool
	Store  db.Store
}

func (s *Server) Create(ctx context.Context, config *configs.Config) error {
	s.Config = config
	s.Router = mux.NewRouter()
	s.HTTP = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", s.Config.ServerAddress, s.Config.Port),
		Handler: s.Router,
	}

	s.setupRoutes()

	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", config.DbUsername, config.DbPassword, config.DbPort, config.DbName)

	dbPool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Errorf("Failed to connect to database: %v", err)
		return fmt.Errorf("database connection failed: %w", err)
	}

	if err := dbPool.Ping(ctx); err != nil {
		log.Errorf("Failed to ping database: %v", err)
		return fmt.Errorf("database ping failed: %w", err)
	}

	s.DB = dbPool
	s.Store = db.NewStore(dbPool)

	log.Info("Database connection established successfully.")
	return nil
}

func (s *Server) Serve(ctx context.Context) error {
	idleConnClosed := make(chan struct{})

	go func(ctx context.Context, s *Server) {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

		<-stop
		log.Info("\n Shutdown signal received")

		if err := s.HTTP.Shutdown(ctx); err != nil {
			log.Errorf("Error shutting down server: %v", err)
		}

		if s.DB != nil {
			s.DB.Close()
			log.Info("Database connection closed")
		}

		close(idleConnClosed)
	}(ctx, s)

	log.Infof("Server is running on %s", s.HTTP.Addr)

	if err := s.HTTP.ListenAndServe(); err != http.ErrServerClosed {
		log.Println("Server error: ", err)
	}

	<-idleConnClosed
	log.Info("Server stopped gracefully.")
	return nil
}
