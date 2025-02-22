package server

import "github.com/gin-gonic/gin"

func (server *Server) setupRouter() {
	router := gin.Default()
	router.GET("/_healthz", Healthz)
	router.POST("/users/register", server.createUser)
	router.POST("/users/login", server.loginUser)

	server.router = router
}
