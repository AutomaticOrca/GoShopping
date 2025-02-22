package server

import "github.com/gin-gonic/gin"

func (server *Server) setupRouter() {
	router := gin.Default()
	router.GET("/_healthz", Healthz)

	userRoutes := router.Group("/users")
	{
		userRoutes.POST("/register", server.createUser)
		userRoutes.POST("/login", server.loginUser)
	}

	categoryRoutes := router.Group("/categories")
	{
		categoryRoutes.POST("/", server.createCategory)
		categoryRoutes.GET("/:id", server.getCategory)
		categoryRoutes.GET("/", server.listCategories)
		categoryRoutes.GET("/parent/:parent_id", server.getSubcategories)
		categoryRoutes.PUT("/", server.updateCategory)
		categoryRoutes.DELETE("/:id", server.deleteCategory)
	}

	server.router = router
}
