package server

import (
	"fmt"
	"github.com/AutomaticOrca/GoShopping/configs"
	db "github.com/AutomaticOrca/GoShopping/internal/database/sqlc"
	"github.com/AutomaticOrca/GoShopping/pkg/token"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config     configs.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config configs.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
